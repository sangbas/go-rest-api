package api

import (
	"github.com/go-rest-api/internal/movie/delivery/http"
	"github.com/go-rest-api/pkg/response"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	nethttp "net/http"

	healthCheckHandler "github.com/go-rest-api/internal/healthcheck/delivery/http"
)

// Route http request pattern
type Route struct {
	healthCheckHandler *healthCheckHandler.HealthCheckHandler
	movieHandler       *http.MovieHandler
}

// NewRoute instances
func NewRoute(healthCheckHandler *healthCheckHandler.HealthCheckHandler, movieHandler *http.MovieHandler) *Route {
	return &Route{
		healthCheckHandler: healthCheckHandler,
		movieHandler:       movieHandler,
	}
}

// GetHandler Build the all router
func (r *Route) GetHandler() nethttp.Handler {
	router := mux.NewRouter()

	router.NotFoundHandler = nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		response.WriteAPIErrorMessage(w, response.APINotFoundHandler)
	})

	v1 := router.PathPrefix("/v1").Subrouter()

	healthCheck := v1.PathPrefix("/health").Subrouter()
	healthCheck.HandleFunc("/api", r.healthCheckHandler.API).Methods("GET")
	healthCheck.HandleFunc("/infrastructure", r.healthCheckHandler.Infrastructure).Methods("GET")

	movie := v1.PathPrefix("/movies").Subrouter()
	movie.HandleFunc("", r.movieHandler.GetAllMovies).Methods("GET")
	movie.HandleFunc("/{id:[0-9]+}", r.movieHandler.GetMovie).Methods("GET")
	movie.HandleFunc("", r.movieHandler.SaveMovie).Methods("POST")

	n := negroni.New()
	n.UseHandler(router)
	return n
}
