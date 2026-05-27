package http

import (
	"bytes"
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
	h.renderTemplate(w, r, "main", nil)
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

	var result domain.SearchResult

	switch mtype {
	case domain.MovieType:
		result, err = h.search.SearchMovies(r.Context(), query, page)
		if err != nil {
			h.render500(w, r, err.Error())
			return
		}
	case domain.TvShowType:
		result, err = h.search.SearchTvShows(r.Context(), query, page)
		if err != nil {
			h.render500(w, r, err.Error())
			return
		}
	default:
		h.render400(w, r)
		return
	}

	var view []SearchItemView
	for _, m := range result.Items {
		item := SearchItemView{
			ID:          m.ID,
			Title:       m.Title,
			Overview:    m.Overview,
			PosterURL:   m.PosterURL,
			ReleaseYear: m.ReleaseYear,
		}
		view = append(view, item)
	}

	data := SearchView{
		SearchString: query,
		MediaType:    string(mtype),
		CurrentPage:  page,
		TotalPages:   result.TotalPages,
		PrevPage:     page - 1,
		NextPage:     page + 1,
		Items:        view,
	}

	h.renderTemplate(w, r, "search", data)
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

	data := MovieView{
		ID:          movie.ID,
		Title:       movie.Title,
		Overview:    movie.Overview,
		PosterURL:   movie.PosterURL,
		ReleaseYear: movie.ReleaseYear,
	}

	h.renderTemplate(w, r, "movie", data)
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

	data := TvShowView{
		ID:          tvShow.ID,
		Title:       tvShow.Title,
		Overview:    tvShow.Overview,
		PosterURL:   tvShow.PosterURL,
		ReleaseYear: tvShow.ReleaseYear,
	}

	h.renderTemplate(w, r, "tv", data)
}

func (h *Handler) ShowNotFound(w http.ResponseWriter, r *http.Request) {
	h.render404(w, r)
}

func (h *Handler) render400(w http.ResponseWriter, r *http.Request) {
	data := ErrorPageView{ErrorCode: http.StatusBadRequest, ErrorTitle: "Bad Request"}
	h.renderError(w, r, data)
}

func (h *Handler) render404(w http.ResponseWriter, r *http.Request) {
	data := ErrorPageView{ErrorCode: http.StatusNotFound, ErrorTitle: "Not Found"}
	h.renderError(w, r, data)
}

func (h *Handler) render500(w http.ResponseWriter, r *http.Request, description string) {
	data := ErrorPageView{
		ErrorCode:        http.StatusInternalServerError,
		ErrorTitle:       "Internal Error",
		ErrorDescription: description,
	}
	h.renderError(w, r, data)
}

func (h *Handler) renderTemplate(w http.ResponseWriter, r *http.Request, name string, data any) {
	var buf bytes.Buffer

	if err := h.tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		h.render500(w, r, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Printf("write error: %v", err)
	}
}

func (h *Handler) renderError(w http.ResponseWriter, r *http.Request, data ErrorPageView) {
	isHtml := strings.Contains(r.Header.Get("Accept"), "text/html")
	if isHtml {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(data.ErrorCode)

	if isHtml {
		var buf bytes.Buffer
		err := h.tmpl.ExecuteTemplate(&buf, "error", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
		}
		_, err = buf.WriteTo(w)
		if err != nil {
			log.Printf("error writing response: %s", err)
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
