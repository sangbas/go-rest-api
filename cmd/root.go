package cmd

import (
	serve "github.com/go-rest-api/cmd/http-serve"
	wrapper "github.com/go-rest-api/pkg/config"
	"github.com/go-rest-api/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	config "github.com/spf13/viper"
)

var (
	rootCmd cobra.Command
)

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "../configurations", "the config path location")
	cobra.OnInitialize(onInitialize)
}

func onInitialize() {
	cfg, err := rootCmd.PersistentFlags().GetString("config")
	if err != nil {
		logrus.Fatal(err)
		return
	}
	wrapper.SetConfig(cfg)

	logLevel := config.GetString("log.level")
	if logLevel == "" {
		logLevel = "debug"
	}
	_ = logger.SetLevel(logLevel)

	// Set circuit breaker

}

// Execute serve
func Execute() {
	var (
		cmdServeHTTP = &cobra.Command{
			Use:   "http-serve",
			Short: "Listening HTTP server",
			Long:  "Service Listening HTTP server",
			Run: func(command *cobra.Command, args []string) {
				httpServer := serve.NewHTTPServer()
				httpServer.Serve(command, args)
			},
		}
	)

	rootCmd.AddCommand(cmdServeHTTP)
	_ = rootCmd.Execute()
}
