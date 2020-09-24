package service

import (
	"context"
	"github.com/go-rest-api/internal/movie/entity"
	"github.com/go-rest-api/internal/movie/repository"
	"github.com/opentracing/opentracing-go"
)

type MovieServiceFactory interface {
	GetAllMovies(ctx context.Context) ([]entity.MovieResp, error)
	GetMovie(ctx context.Context, movieId int64) (entity.MovieResp, error)
	SaveMovie(ctx context.Context, movieRepo entity.MovieRepo) (entity.MovieRepo, error)
}

type MovieService struct {
	repo repository.MovieRepositoryFactory
}

func NewMovieService(repo repository.MovieRepositoryFactory) (*MovieService, error) {
	return &MovieService{
		repo: repo,
	}, nil
}

func (m *MovieService) GetAllMovies(ctx context.Context) ([]entity.MovieResp, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "")
	defer span.Finish()

	var movieResps []entity.MovieResp
	movieRepos, err := m.repo.GetAllMovies(ctx)
	if err != nil {
		return movieResps, err
	}

	for _, movieRepo := range movieRepos {
		movieResp := entity.MovieResp{
			ID:       movieRepo.ID,
			Name:     movieRepo.Name,
			Duration: movieRepo.Duration,
			Genre:    movieRepo.Genre,
		}
		movieResps = append(movieResps, movieResp)
	}

	return movieResps, nil
}

func (m *MovieService) GetMovie(ctx context.Context, movieId int64) (entity.MovieResp, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "")
	defer span.Finish()

	var movieResp entity.MovieResp
	movieRepo, err := m.repo.GetMovie(ctx, movieId)
	if err != nil {
		return movieResp, err
	}

	movieResp = entity.MovieResp{
		ID:       movieRepo.ID,
		Name:     movieRepo.Name,
		Duration: movieRepo.Duration,
		Genre:    movieRepo.Genre,
	}

	return movieResp, nil
}

func (m *MovieService) SaveMovie(ctx context.Context, movieRepo entity.MovieRepo) (entity.MovieRepo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "")
	defer span.Finish()

	movie, err := m.repo.SaveMovie(ctx, movieRepo)
	if err != nil {
		return movie, err
	}

	return movie, nil
}
