package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
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
		h.render500(w, r)
	}
}

func (h *Handler) HandleSearch(query string, w http.ResponseWriter, r *http.Request) {

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(data.ErrorCode)

	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := h.tmpl.ExecuteTemplate(w, "error", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
		}
		return
	}

	_, err := w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, data.ErrorDescription)))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
