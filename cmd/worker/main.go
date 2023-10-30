package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/vlladoff/young_astrologer/internal/config"
	"github.com/vlladoff/young_astrologer/internal/http-server/handlers/astro"
	"github.com/vlladoff/young_astrologer/internal/lib/logger/sl"
	"github.com/vlladoff/young_astrologer/internal/storage/postgresql"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
	"time"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	storage, err := postgresql.New(cfg.StorageDataSource)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	defer postgresql.Close(storage)

	s := gocron.NewScheduler(time.UTC)
	job, err := s.Every(1).Day().Do(fetchAPOD, storage, cfg)
	//job, err := s.Every(1).Day().At("00:00").Do(fetchAPOD, storage, cfg)

	log.Info("young astrologer worker started")

	if err != nil {
		log.Error("Error in job", job.GetName(), sl.Err(err))
	}

	job.RegisterEventListeners(
		gocron.BeforeJobRuns(func(jobName string) {
			log.Info("job started", jobName)
		}),
		gocron.WhenJobReturnsError(func(jobName string, err error) {
			log.Error("Error in job", jobName, sl.Err(err))
		}),
		gocron.WhenJobReturnsNoError(func(jobName string) {
			log.Info("job done", jobName)
		}),
	)

	s.StartBlocking()
}

func fetchAPOD(storage *postgresql.Storage, cfg *config.Config) error {
	const op = "worker.fetchAPOD"

	u, err := url.Parse(cfg.APODEndpoint)
	if err != nil {
		sl.Err(fmt.Errorf("%s: %w", op, err))
	}

	queryParams := url.Values{}
	queryParams.Add("api_key", cfg.APODAPIKey)

	u.RawQuery = queryParams.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		sl.Err(fmt.Errorf("%s: %w", op, resp.Status))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var data astro.AstroData
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var image, hdImage []byte
	var imagesId int64
	if data.MediaType == "image" {
		if data.Url != "" && data.HdUrl != "" {
			image, hdImage, err = getImages(data.Url, data.HdUrl)
		} else if data.Url != "" {
			image, err = downloadImage(data.Url)
		}
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if image != nil {
			parsedURL, _ := url.Parse(data.Url)
			data.Url = path.Base(parsedURL.Path)
		}
		if hdImage != nil {
			parsedURL, _ := url.Parse(data.HdUrl)
			data.HdUrl = path.Base(parsedURL.Path)
		}

		if image != nil || hdImage != nil {
			imagesId, err = storage.SaveImages(image, hdImage)
			if err != nil {
				return err
			}
		}
	}

	if err := storage.SaveAstroData(data, imagesId); err != nil {
		return err
	}

	return nil
}

func getImages(url, hdUrl string) ([]byte, []byte, error) {
	const op = "worker.getImages"

	var wg sync.WaitGroup
	imageCh := make(chan []byte)
	hdImageCh := make(chan []byte)
	var imageErr, hdImageErr error

	wg.Add(1)
	go getImage(url, &wg, imageCh, &imageErr)

	wg.Add(1)
	go getImage(hdUrl, &wg, hdImageCh, &hdImageErr)

	go func() {
		wg.Wait()
		close(imageCh)
		close(hdImageCh)
	}()

	imageData := <-imageCh
	hdImageData := <-hdImageCh

	if imageErr != nil || hdImageErr != nil {
		return nil, nil, fmt.Errorf("%s: %v, %v", op, imageErr, hdImageErr)
	}

	return imageData, hdImageData, nil
}

func getImage(url string, wg *sync.WaitGroup, ch chan []byte, returnedErr *error) {
	defer wg.Done()

	imageData, err := downloadImage(url)
	if err != nil {
		*returnedErr = err
		return
	}

	ch <- imageData
}

func downloadImage(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}
