package astro

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/vlladoff/young_astrologer/internal/lib/api/response"
	"github.com/vlladoff/young_astrologer/internal/lib/logger/sl"
	"io"
	"log/slog"
	"net/http"
)

type Request struct {
	Date string `json:"date"`
}

type Response struct {
	response.Response
	Data interface{} `json:"data,omitempty"`
}

type AstroData struct {
	MediaType   string `json:"media_type,omitempty"`
	Date        string `json:"date"`
	Explanation string `json:"explanation"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	HdUrl       string `json:"hdurl,omitempty"`
}

type AstroGetter interface {
	GetAllAstroData() ([]AstroData, error)
	GetAstroDataByDay(date string) (AstroData, error)
}

func GetAll(log *slog.Logger, astroGetter AstroGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.astro.GetAll"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		data, err := astroGetter.GetAllAstroData()
		if err != nil {
			log.Error("failed to get astro data", sl.Err(err))

			render.JSON(w, r, response.Error("failed to get astro data"))

			return
		}

		log.Info("request done")

		responseOK(w, r, data)
	}
}

func GetByDay(log *slog.Logger, astroGetter AstroGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.astro.GetByDay"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, response.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		data, err := astroGetter.GetAstroDataByDay(req.Date)
		if err != nil {
			log.Error("failed to get astro data", sl.Err(err))

			render.JSON(w, r, response.Error("failed to get astro data"))

			return
		}

		log.Info("request done")

		responseOK(w, r, []AstroData{data})
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, data []AstroData) {
	if len(data) > 1 {
		render.JSON(w, r, Response{
			Response: response.OK(),
			Data:     data,
		})
	} else {
		render.JSON(w, r, Response{
			Response: response.OK(),
			Data:     data[0],
		})
	}
}
