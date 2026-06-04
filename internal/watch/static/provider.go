package static

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/nikitagrgv/movies/internal/watch"
)

// WatchServerDescription
// Template placeholders:
// {id} - Movie/TV Show ID
// {s} - TV show season
// {e} - TV show episode
type WatchServerDescription struct {
	ID                string
	Name              string
	MovieURLTemplate  string
	TvShowURLTemplate string
}

type WatchServerProvider struct {
	servers    []watch.WatchServer
	serversMap map[string]WatchServerDescription
}

func NewWatchServerProvider(servers []WatchServerDescription) (*WatchServerProvider, error) {
	p := &WatchServerProvider{}
	p.serversMap = make(map[string]WatchServerDescription)
	for _, server := range servers {
		if _, ok := p.serversMap[server.ID]; ok {
			return nil, fmt.Errorf("duplicate server %s", server.ID)
		}
		p.serversMap[server.ID] = server

		p.servers = append(p.servers, watch.WatchServer{
			ID:   server.ID,
			Name: server.Name,
		})
	}

	return p, nil
}

func (p *WatchServerProvider) GetServers(ctx context.Context) ([]watch.WatchServer, error) {
	return p.servers, nil
}

func (p *WatchServerProvider) GetMovieWatchLink(ctx context.Context, serverID string, movieID int) (string, error) {
	server, ok := p.serversMap[serverID]
	if !ok {
		return "", fmt.Errorf("server %s not found", serverID)
	}

	cooked := server.MovieURLTemplate
	cooked = strings.ReplaceAll(cooked, "{id}", strconv.Itoa(movieID))
	return cooked, nil
}

func (p *WatchServerProvider) GetTvShowWatchLink(ctx context.Context, serverID string, tvID, season, episode int) (string, error) {
	server, ok := p.serversMap[serverID]
	if !ok {
		return "", fmt.Errorf("server %s not found", serverID)
	}

	cooked := server.TvShowURLTemplate
	cooked = strings.ReplaceAll(cooked, "{id}", strconv.Itoa(tvID))
	cooked = strings.ReplaceAll(cooked, "{s}", strconv.Itoa(season))
	cooked = strings.ReplaceAll(cooked, "{e}", strconv.Itoa(episode))
	return cooked, nil
}
