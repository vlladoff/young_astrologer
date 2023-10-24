package astro

import (
	"github.com/go-chi/render"
	"github.com/vlladoff/young_astrologer/internal/lib/api/response"
	"log/slog"
	"net/http"
)

type Request struct {
	Day string `json:"day"`
}

type Response struct {
	response.Response
	Data AstroData
}

type AstroData struct {
	Date        string `json:"date,omitempty"`
	Explanation string `json:"explanation,omitempty"`
	Hdurl       string `json:"hdurl,omitempty"`
	Title       string `json:"title,omitempty"`
	Url         string `json:"url,omitempty"`
	Image       string `json:"image,omitempty"`
}

type AstroGetter interface {
	GetAllAstroData() ([]AstroData, error)
	GetAstroDataByDay(date string) (AstroData, error)
}

func GetAll(log *slog.Logger, astroGetter AstroGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func GetByDay(log *slog.Logger, astroGetter AstroGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, data AstroData) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Data:     data,
	})
}
