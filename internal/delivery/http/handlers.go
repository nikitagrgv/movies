package http

import (
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	tmpl *template.Template
}

func NewHandler(tmpl *template.Template) *Handler {
	return &Handler{tmpl: tmpl}
}

func (h *Handler) ShowMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := h.tmpl.ExecuteTemplate(w, "main", nil)
	if err != nil {
		h.render500(w)
	}
}

func (h *Handler) HandleSearch(query string, w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) ShowNotFound(w http.ResponseWriter, r *http.Request) {
	h.render404(w)
}

func (h *Handler) handleError(err error, w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) render404(w http.ResponseWriter) {
	data := ErrorPageData{ErrorCode: http.StatusNotFound, ErrorDescription: "Not Found"}
	h.renderError(w, data)
}

func (h *Handler) render500(w http.ResponseWriter) {
	data := ErrorPageData{ErrorCode: http.StatusInternalServerError, ErrorDescription: "Internal Error"}
	h.renderError(w, data)
}

func (h *Handler) renderError(w http.ResponseWriter, data ErrorPageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(data.ErrorCode)

	err := h.tmpl.ExecuteTemplate(w, "error", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
