package usecase

import (
	"context"

	"github.com/nikitagrgv/movies/internal/domain"
)

type GetMediaUsecase struct {
	getter     domain.MediaGetter
	noImageURL string
}

func NewGetMediaUsecase(getter domain.MediaGetter, noImageURL string) *GetMediaUsecase {
	return &GetMediaUsecase{getter: getter, noImageURL: noImageURL}
}

func (u *GetMediaUsecase) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
	media, err := u.getter.GetMovie(ctx, id)
	if err != nil {
		return domain.Movie{}, err
	}

	media.Base = normalizeMedia(media.Base, u.noImageURL)
	return media, nil
}

func (u *GetMediaUsecase) GetTvShow(ctx context.Context, id int) (domain.TvShow, error) {
	media, err := u.getter.GetTvShow(ctx, id)
	if err != nil {
		return domain.TvShow{}, err
	}

	media.Base = normalizeMedia(media.Base, u.noImageURL)
	return media, nil
}

func normalizeMedia(media domain.MediaBase, noImageURL string) domain.MediaBase {
	if media.PosterURL == "" {
		media.PosterURL = noImageURL
	}

	return media
}
