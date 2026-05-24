package http

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/nikitagrgv/movies/internal/usecase"
)

type Handler struct {
	tmpl   *template.Template
	search *usecase.SearchMoviesUsecase
}

func NewHandler(tmpl *template.Template, search *usecase.SearchMoviesUsecase) *Handler {
	return &Handler{tmpl: tmpl, search: search}
}

func (h *Handler) ShowMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := h.tmpl.ExecuteTemplate(w, "main", nil)
	if err != nil {
		h.render500(w, r)
	}
}

func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("s")
	pageStr := r.URL.Query().Get("p")

	var page int = 1
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			h.render500(w, r)
			return
		}
		page = p
	}

	result, err := h.search.SearchMovies(r.Context(), query, page)
	if err != nil {
		h.render500(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := SearchPageData{
		SearchString: query,
		CurrentPage:  page,
		TotalPages:   result.TotalPages,
		PrevPage:     page - 1,
		NextPage:     page + 1,
		Movies:       result.Movies,
	}
	err = h.tmpl.ExecuteTemplate(w, "search", data)
	if err != nil {
		h.render500(w, r)
	}
}

func (h *Handler) ShowNotFound(w http.ResponseWriter, r *http.Request) {
	h.render404(w, r)
}

func (h *Handler) render404(w http.ResponseWriter, r *http.Request) {
	data := ErrorPageData{ErrorCode: http.StatusNotFound, ErrorDescription: "Not Found"}
	h.renderError(w, r, data)
}

func (h *Handler) render500(w http.ResponseWriter, r *http.Request) {
	data := ErrorPageData{ErrorCode: http.StatusInternalServerError, ErrorDescription: "Internal Error"}
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
		"error": data.ErrorDescription,
	})
	_, err := w.Write(b)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
