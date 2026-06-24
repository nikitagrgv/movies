package web

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/nikitagrgv/movies/internal/media"
	"github.com/nikitagrgv/movies/internal/watch"
)

type Handler struct {
	tmpl  *template.Template
	media *media.Service
	watch *watch.Service
}

func NewHandler(tmpl *template.Template, media *media.Service, watch *watch.Service) *Handler {
	return &Handler{tmpl: tmpl, media: media, watch: watch}
}

func LoadTemplates(cacheVersion int) (*template.Template, error) {
	funcMap := template.FuncMap{
		"static": func(relPath string) string {
			return ResolveStaticAssetPath(cacheVersion, relPath)
		},
	}

	tmpl, err := template.
		New("").
		Funcs(funcMap).
		ParseFS(Assets, "templates/*.html", "templates/partials/*.html")

	return tmpl, err
}

func (h *Handler) showMain(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, r, "main", nil)
}

func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("s")
	searchType := r.URL.Query().Get("type")
	pageStr := r.URL.Query().Get("p")

	mtype, err := media.ParseMediaType(searchType)
	if err != nil {
		h.render400(w, r)
		return
	}

	var page = 1
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			h.render400(w, r)
			return
		}
		page = p
	}

	var result media.SearchResult

	switch mtype {
	case media.MovieType:
		result, err = h.media.SearchMovies(r.Context(), query, page)
		if err != nil {
			h.render500(w, r, err.Error())
			return
		}
	case media.TvShowType:
		result, err = h.media.SearchTvShows(r.Context(), query, page)
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

func (h *Handler) handleMovie(idStr string, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.render400(w, r)
		return
	}

	var serverId = ""
	if r.URL.Query().Has("srv") {
		srv := r.URL.Query().Get("srv")
		serverId = srv
	}

	movie, err := h.media.GetMovie(r.Context(), id)
	if err != nil {
		h.render500(w, r, err.Error())
		return
	}

	servers, err := h.watch.GetServers(r.Context())
	if err != nil {
		h.render500(w, r, err.Error())
		return
	}

	if len(servers) == 0 {
		h.render500(w, r, "no servers found")
		return
	}

	if serverId == "" {
		serverId = servers[0].ID
	}

	watchURL, err := h.watch.GetMovieWatchLink(r.Context(), serverId, id)
	if err != nil {
		h.render500(w, r, err.Error())
		return
	}

	var serverViews []WatchServerView
	for _, s := range servers {
		serverViews = append(serverViews, WatchServerView{
			Name: s.Name,
			ID:   s.ID,
		})
	}

	data := MovieView{
		ID:            movie.ID,
		Title:         movie.Title,
		Overview:      movie.Overview,
		PosterURL:     movie.PosterURL,
		ReleaseYear:   movie.ReleaseYear,
		CurrentServer: serverId,
		WatchURL:      watchURL,
		Servers:       serverViews,
	}

	h.renderTemplate(w, r, "movie", data)
}

func (h *Handler) handleTvShow(idStr string, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.render400(w, r)
		return
	}

	var seasonNum = 1
	var episodeNum = 1
	var serverId = ""
	if r.URL.Query().Has("s") {
		s, err := strconv.Atoi(r.URL.Query().Get("s"))
		if err != nil {
			h.render400(w, r)
		}
		seasonNum = s
	}
	if r.URL.Query().Has("e") {
		e, err := strconv.Atoi(r.URL.Query().Get("e"))
		if err != nil {
			h.render400(w, r)
		}
		episodeNum = e
	}
	if r.URL.Query().Has("srv") {
		srv := r.URL.Query().Get("srv")
		serverId = srv
	}

	tvShow, err := h.media.GetTvShow(r.Context(), id)
	if err != nil {
		h.render500(w, r, err.Error())
		return
	}

	season, err := h.media.GetTvShowSeason(r.Context(), id, seasonNum)
	if err != nil {
		h.render500(w, r, err.Error())
		return
	}

	servers, err := h.watch.GetServers(r.Context())
	if err != nil {
		h.render500(w, r, err.Error())
		return
	}

	if len(servers) == 0 {
		h.render500(w, r, "no servers found")
		return
	}

	if serverId == "" {
		serverId = servers[0].ID
	}

	watchURL, err := h.watch.GetTvShowWatchLink(r.Context(), serverId, id, seasonNum, episodeNum)
	if err != nil {
		h.render500(w, r, err.Error())
		return
	}

	var seasonViews []SeasonView
	for i := range tvShow.TotalSeasons {
		seasonNumber := i + 1
		seasonViews = append(seasonViews, SeasonView{
			SeasonNumber: seasonNumber,
			Name:         "S" + strconv.Itoa(seasonNumber),
			EpisodeCount: episodeNum})
	}

	var episodeViews []EpisodeView
	for i := range len(season.Episodes) {
		episodeNumber := i + 1
		var date string
		if season.Episodes[i].Date.IsZero() {
			date = "Unknown"
		} else {
			date = season.Episodes[i].Date.Format("January 2, 2006")
		}
		episodeViews = append(episodeViews, EpisodeView{
			EpisodeNumber: episodeNumber,
			Name:          season.Episodes[i].Name,
			Date:          date,
			IsAvailable:   true, // TODO: #implement
		})
	}

	var serverViews []WatchServerView
	for _, s := range servers {
		serverViews = append(serverViews, WatchServerView{
			Name: s.Name,
			ID:   s.ID,
		})
	}

	data := TvShowView{
		ID:             tvShow.ID,
		Title:          tvShow.Title,
		Overview:       tvShow.Overview,
		PosterURL:      tvShow.PosterURL,
		ReleaseYear:    tvShow.ReleaseYear,
		CurrentSeason:  seasonNum,
		CurrentEpisode: episodeNum,
		CurrentServer:  serverId,
		WatchURL:       watchURL,
		Seasons:        seasonViews,
		Episodes:       episodeViews,
		Servers:        serverViews,
	}

	h.renderTemplate(w, r, "tv", data)
}

func (h *Handler) showNotFound(w http.ResponseWriter, r *http.Request) {
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
