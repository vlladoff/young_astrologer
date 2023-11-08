package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vlladoff/young_astrologer/internal/http-server/handlers/astro"
	"github.com/vlladoff/young_astrologer/internal/lib/logger/sl"
	"io"
	"net/http"
	"net/url"
	"path"
)

func (w *Worker) fetchAPOD() error {
	const op = "worker.fetchAPOD"

	u, err := url.Parse(w.cfg.APODEndpoint)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	queryParams := url.Values{}
	queryParams.Add("api_key", w.cfg.APODAPIKey)

	u.RawQuery = queryParams.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
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

	var image, hdImage *bytes.Buffer
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
			imagesId, err = w.storage.SaveImages(image, hdImage)
			if err != nil {
				return err
			}
		}
	}

	if err := w.storage.SaveAstroData(data, imagesId); err != nil {
		return err
	}

	return nil
}
