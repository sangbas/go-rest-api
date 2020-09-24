package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-rest-api/internal/movie/entity"
	"github.com/go-rest-api/pkg/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type MovieRepositoryFactory interface {
	GetAllMovies(ctx context.Context) ([]entity.MovieRepo, error)
	GetMovie(ctx context.Context, movieId int64) (entity.MovieRepo, error)
	SaveMovie(ctx context.Context, movieRepo entity.MovieRepo) (entity.MovieRepo, error)
}

type MovieRepository struct {
	mysql mysql.BaseRepository
}

func NewMovieRepository(masterDB *sqlx.DB, slaveDB *sqlx.DB) (*MovieRepository, error) {
	if masterDB == nil {
		return nil, errors.New("the master DB connection is nil")
	}

	if slaveDB == nil {
		return nil, errors.New("the slave DB connection is nil")
	}

	m := &MovieRepository{}
	m.mysql.MasterDB = masterDB
	m.mysql.SlaveDB = slaveDB
	return m, nil
}

func (m *MovieRepository) GetAllMovies(ctx context.Context) ([]entity.MovieRepo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "")
	defer span.Finish()

	q := fmt.Sprintf("select id, name, duration, genre from movies")

	var movies []entity.MovieRepo

	err := m.mysql.FetchRows(ctx, q, &movies)
	if err != nil {
		return movies, err
	}

	return movies, nil
}

func (m *MovieRepository) GetMovie(ctx context.Context, movieId int64) (entity.MovieRepo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "")
	defer span.Finish()

	q := fmt.Sprintf("select id, name, duration, genre from movies where id = ?")

	var movie entity.MovieRepo

	err := m.mysql.FetchRow(ctx, q, &movie, movieId)
	if err != nil {
		return movie, err
	}

	return movie, nil

}

func (m *MovieRepository) SaveMovie(ctx context.Context, movieRepo entity.MovieRepo) (entity.MovieRepo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "")
	defer span.Finish()

	q := fmt.Sprintf("insert into movies (name, genre, duration) values (:name, :genre, :duration)")

	res, err := m.mysql.Exec(ctx, q, movieRepo)
	if err != nil {
		return movieRepo, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return movieRepo, err
	}
	movieRepo.ID = id

	return movieRepo, nil
}
