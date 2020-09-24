package http

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/go-rest-api/internal/movie/entity"
	"github.com/go-rest-api/internal/movie/service"
	"github.com/go-rest-api/pkg/response"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	nethttp "net/http"
	"strconv"
)

type MovieHandler struct {
	service service.MovieServiceFactory
}

func NewMovieHandler(service service.MovieServiceFactory) (*MovieHandler, error) {
	return &MovieHandler{
		service: service,
	}, nil
}

func (m *MovieHandler) GetAllMovies(w nethttp.ResponseWriter, r *nethttp.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "")
	defer span.Finish()

	movies, err := m.service.GetAllMovies(ctx)
	if err != nil {
		response.WriteAPIErrorMessage(w, response.APIInternalError)
		return
	}

	response.WriteAPIOKWithData(w, movies)
}

func (m *MovieHandler) GetMovie(w nethttp.ResponseWriter, r *nethttp.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "")
	defer span.Finish()

	params := mux.Vars(r)
	movieId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		response.WriteAPIError(w, response.APIErrorBadRequest, err)
		return
	}

	movie, err := m.service.GetMovie(ctx, movieId)
	if err != nil {
		response.WriteAPIErrorMessage(w, response.APIInternalError)
		return
	}

	response.WriteAPIOKWithData(w, movie)
}

func (m *MovieHandler) SaveMovie(w nethttp.ResponseWriter, r *nethttp.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "")
	defer span.Finish()

	var payload entity.MovieRepo
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		response.WriteAPIError(w, response.APIErrorBadRequest, err)
		return
	}

	isValid, err := govalidator.ValidateStruct(payload)
	if !isValid {
		response.WriteAPIError(w, response.APIErrorBadRequest, err)
		return
	}

	movie, err := m.service.SaveMovie(ctx, payload)
	if err != nil {
		response.WriteAPIErrorMessage(w, response.APIInternalError)
		return
	}

	response.WriteAPIOKWithData(w, movie)
}
