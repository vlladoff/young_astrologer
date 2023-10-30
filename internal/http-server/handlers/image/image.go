package image

import (
	"github.com/go-chi/render"
	"github.com/vlladoff/young_astrologer/internal/lib/logger/sl"
	"log/slog"
	"net/http"
)

type ImageGetter interface {
	GetImageData(imageName string) ([]byte, error)
}

func GetImage(log *slog.Logger, imageGetter ImageGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		imageName := r.URL.Path[1:]

		imageData, err := imageGetter.GetImageData(imageName)
		if err != nil {
			log.Error("failed to retrieve image", sl.Err(err))
			render.Status(r, http.StatusNotFound)

			return
		}

		w.Header().Set("Content-Type", "image/jpeg")

		_, err = w.Write(imageData)
		if err != nil {
			log.Error("failed to write image data to the response", sl.Err(err))
			render.Status(r, http.StatusNotFound)

			return
		}
	}
}
