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
	MediaType   string `json:"media_type,omitempty"`
	Date        string `json:"date"`
	Explanation string `json:"explanation"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	HdUrl       string `json:"hdurl,omitempty"`
	OriginalUrl string `json:"original_url"`
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
