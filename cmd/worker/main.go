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
	"sync"
	"time"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	storage, err := postgresql.New(cfg.StorageDataSource)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	defer postgresql.Close(storage)

	s := gocron.NewScheduler(time.UTC)
	//job, err := s.Every(1).Second().Do(fetchAPOD, log, storage, cfg)
	job, err := s.Every(1).Day().At("00:00").Do(fetchAPOD, log, storage, cfg)

	log.Info("young astrologer worker started")

	if err != nil {
		log.Error("Error in job", job.GetName(), sl.Err(err))
	}

	s.StartBlocking()
}

func fetchAPOD(log *slog.Logger, storage *postgresql.Storage, cfg *config.Config) error {
	const op = "worker.fetchAPOD"

	log.Info("fetched apod data started")

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
	if data.Url != "" && data.HdUrl != "" {
		image, hdImage, err = getImages(data.Url, data.HdUrl)
	} else if data.Url != "" {
		image, err = downloadImage(data.Url)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, _ = image, hdImage
	if err := storage.SaveAstroData(data, image, hdImage); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("fetched apod data successful complete")

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
