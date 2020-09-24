package http_serve

import (
	"context"
	"fmt"
	"github.com/go-rest-api/api"
	"github.com/go-rest-api/internal/movie/delivery/http"
	"github.com/go-rest-api/internal/movie/repository"
	"github.com/go-rest-api/internal/movie/service"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	config "github.com/spf13/viper"
	nethttp "net/http"
	"os"
	"os/signal"
	"time"

	healthCheckHandler "github.com/go-rest-api/internal/healthcheck/delivery/http"
	healthCheckRepository "github.com/go-rest-api/internal/healthcheck/repository"
	healthCheckService "github.com/go-rest-api/internal/healthcheck/service"
)

const (
	serviceName = "go-rest-api-http-serve"
	bannerInfo  = `
Name: %s
Port: %s
------------------------------------------------------------------------------
`
)

// env constants
const (
	development = "development"
	staging     = "staging"
	production  = "production"
)

// DsnFormat stands for database source name format
var DsnFormat = "%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=%s&loc=%s"

// Server as the http server
type Server struct {
	dbMaster *sqlx.DB
	dbSlave  *sqlx.DB
}

func NewHTTPServer() *Server {
	s := &Server{}
	dbMaster, err := s.buildMysqlClientMaster()
	if err != nil {
		panic(err)
	}

	dbSlave, err := s.buildMysqlClientSlave()
	if err != nil {
		panic(err)
	}

	s.dbMaster = dbMaster
	s.dbSlave = dbSlave

	return s
}

func (s *Server) buildMysqlClientMaster() (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("%s:%s@(%s:%s)/%s", config.GetString("database.master.user"),
		config.GetString("database.master.password"),
		config.GetString("database.master.host"),
		config.GetString("database.master.port"),
		config.GetString("database.master.name"),
	)
	db, err := sqlx.Connect("mysql", dataSource)
	if err != nil {
		return db, err
	}
	return db, nil
}

func (s *Server) buildMysqlClientSlave() (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("%s:%s@(%s:%s)/%s", config.GetString("database.slave.user"),
		config.GetString("database.slave.password"),
		config.GetString("database.slave.host"),
		config.GetString("database.slave.port"),
		config.GetString("database.slave.name"),
	)
	db, err := sqlx.Connect("mysql", dataSource)
	if err != nil {
		return db, err
	}
	return db, nil
}

// Serve listen and serve server
func (s *Server) Serve(cmd *cobra.Command, args []string) {
	// HealthCheck
	healthCheckRepo, err := healthCheckRepository.NewHealthCheckRepository(s.dbMaster, s.dbSlave)
	if err != nil {
		panic(err)
	}

	healthCheckService, err := healthCheckService.NewHealthCheckService(healthCheckRepo)
	if err != nil {
		panic(err)
	}

	healthCheckDelegate, err := healthCheckHandler.NewHealthCheckHandler(healthCheckService)
	if err != nil {
		panic(err)
	}

	movieRepo, err := repository.NewMovieRepository(s.dbMaster, s.dbSlave)
	if err != nil {
		panic(err)
	}

	movieService, err := service.NewMovieService(movieRepo)
	if err != nil {
		panic(err)
	}

	movieDelegate, err := http.NewMovieHandler(movieService)
	if err != nil {
		panic(err)
	}

	httpHandler := api.NewRoute(healthCheckDelegate, movieDelegate).GetHandler()
	server := &nethttp.Server{
		Addr:    fmt.Sprintf(":%d", config.GetInt("app.port")),
		Handler: httpHandler,
	}

	printBannerInfo(server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	gracefulTimeout := time.Duration(config.GetInt("app.graceful_timeout")) * time.Second
	if gracefulTimeout == 0 {
		gracefulTimeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Printf("HTTP server Shutdown: %v", err)
	}
	logger.Println("shutting down")
	os.Exit(0)
}

func printBannerInfo(port string) {
	fmt.Printf(bannerInfo, serviceName, port)
}
