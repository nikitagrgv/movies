package http

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/nikitagrgv/movies/internal/domain"
	"github.com/nikitagrgv/movies/internal/usecase"
)

type Handler struct {
	tmpl   *template.Template
	search *usecase.SearchMediaUsecase
	get    *usecase.GetMediaUsecase
}

func NewHandler(tmpl *template.Template, search *usecase.SearchMediaUsecase, get *usecase.GetMediaUsecase) *Handler {
	return &Handler{tmpl: tmpl, search: search, get: get}
}

func (h *Handler) ShowMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := h.tmpl.ExecuteTemplate(w, "main", nil)
	if err != nil {
		h.render500(w, r, err.Error())
	}
}

func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("s")
	searchType := r.URL.Query().Get("type")
	pageStr := r.URL.Query().Get("p")

	mtype, err := domain.ParseMediaType(searchType)
	if err != nil {
		h.render400(w, r)
		return
	}

	var page int = 1
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			h.render400(w, r)
			return
		}
		page = p
	}

	var totalPages int
	var result []domain.MediaBase

	switch mtype {
	case domain.MovieType:
		sr, err := h.search.SearchMovies(r.Context(), query, page)
		if err != nil {
			h.render500(w, r, err.Error())
			return
		}
		totalPages = sr.TotalPages
		for _, m := range sr.Movies {
			result = append(result, m.Base)
		}
	case domain.TvShowType:
		sr, err := h.search.SearchTvShows(r.Context(), query, page)
		if err != nil {
			h.render500(w, r, err.Error())
			return
		}
		for _, m := range sr.TvShows {
			result = append(result, m.Base)
		}
	default:
		h.render400(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := SearchPageData{
		SearchString: query,
		MediaType:    string(mtype),
		CurrentPage:  page,
		TotalPages:   totalPages,
		PrevPage:     page - 1,
		NextPage:     page + 1,
		Medias:       result,
	}
	err = h.tmpl.ExecuteTemplate(w, "search", data)
	if err != nil {
		h.render500(w, r, err.Error())
	}
}

func (h *Handler) HandleMovie(idStr string, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.render400(w, r)
		return
	}

	movie, err := h.get.GetMovie(r.Context(), id)
	if err != nil {
		h.render500(w, r, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := MovieView{
		ID:          movie.Base.ID,
		Title:       movie.Base.Title,
		Overview:    movie.Base.Overview,
		PosterURL:   movie.Base.PosterURL,
		ReleaseDate: movie.Base.ReleaseDate,
	}
	err = h.tmpl.ExecuteTemplate(w, "movie", data)
	if err != nil {
		h.render500(w, r, err.Error())
	}
}

func (h *Handler) HandleTvShow(idStr string, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.render400(w, r)
		return
	}

	tvShow, err := h.get.GetTvShow(r.Context(), id)
	if err != nil {
		h.render500(w, r, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := TvShowView{
		ID:          tvShow.Base.ID,
		Title:       tvShow.Base.Title,
		Overview:    tvShow.Base.Overview,
		PosterURL:   tvShow.Base.PosterURL,
		ReleaseDate: tvShow.Base.ReleaseDate,
	}
	err = h.tmpl.ExecuteTemplate(w, "tvshow", data)
	if err != nil {
		h.render500(w, r, err.Error())
	}
}

func (h *Handler) ShowNotFound(w http.ResponseWriter, r *http.Request) {
	h.render404(w, r)
}

func (h *Handler) render400(w http.ResponseWriter, r *http.Request) {
	data := ErrorPageData{ErrorCode: http.StatusBadRequest, ErrorTitle: "Bad Request"}
	h.renderError(w, r, data)
}

func (h *Handler) render404(w http.ResponseWriter, r *http.Request) {
	data := ErrorPageData{ErrorCode: http.StatusNotFound, ErrorTitle: "Not Found"}
	h.renderError(w, r, data)
}

func (h *Handler) render500(w http.ResponseWriter, r *http.Request, description string) {
	data := ErrorPageData{
		ErrorCode:        http.StatusInternalServerError,
		ErrorTitle:       "Internal Error",
		ErrorDescription: description,
	}
	h.renderError(w, r, data)
}

func (h *Handler) renderError(w http.ResponseWriter, r *http.Request, data ErrorPageData) {
	isHtml := strings.Contains(r.Header.Get("Accept"), "text/html")
	if isHtml {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(data.ErrorCode)

	if isHtml {
		err := h.tmpl.ExecuteTemplate(w, "error", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
		}
		return
	}

	b, _ := json.Marshal(map[string]string{
		"error": data.ErrorTitle,
	})
	_, err := w.Write(b)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
