package static

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/nikitagrgv/movies/internal/domain"
)

// WatchService
// Template placeholders:
// {id} - Movie/TV Show ID
// {s} - TV show season
// {e} - TV show episode
type WatchService struct {
	ID                int
	Name              string
	MovieURLTemplate  string
	TvShowURLTemplate string
}

type WatchServiceProvider struct {
	services    []domain.WatchService
	servicesMap map[int]WatchService
}

func NewWatchServiceProvider(services []WatchService) (*WatchServiceProvider, error) {
	p := &WatchServiceProvider{}
	p.servicesMap = make(map[int]WatchService)
	for _, service := range services {
		if _, ok := p.servicesMap[service.ID]; ok {
			return nil, fmt.Errorf("duplicate service %d", service.ID)
		}
		p.servicesMap[service.ID] = service

		p.services = append(p.services, domain.WatchService{
			ID:   service.ID,
			Name: service.Name,
		})
	}

	return p, nil
}

func (p *WatchServiceProvider) GetServices(ctx context.Context) ([]domain.WatchService, error) {
	return p.services, nil
}

func (p *WatchServiceProvider) GetMovieWatchLink(ctx context.Context, serviceID, movieID int) (string, error) {
	service, ok := p.servicesMap[serviceID]
	if !ok {
		return "", fmt.Errorf("service %d not found", serviceID)
	}

	cooked := service.MovieURLTemplate
	cooked = strings.ReplaceAll(cooked, "{id}", strconv.Itoa(movieID))
	return cooked, nil
}

func (p *WatchServiceProvider) GetTvShowWatchLink(ctx context.Context, serviceID, tvID, season, episode int) (string, error) {
	service, ok := p.servicesMap[serviceID]
	if !ok {
		return "", fmt.Errorf("service %d not found", serviceID)
	}

	cooked := service.TvShowURLTemplate
	cooked = strings.ReplaceAll(cooked, "{id}", strconv.Itoa(tvID))
	cooked = strings.ReplaceAll(cooked, "{s}", strconv.Itoa(season))
	cooked = strings.ReplaceAll(cooked, "{e}", strconv.Itoa(episode))
	return cooked, nil
}
