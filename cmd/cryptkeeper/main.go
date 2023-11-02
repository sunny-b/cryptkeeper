package main

import (
	"fmt"
	"os"

	"github.com/sunny-b/cryptkeeper/internal/commands"
	"github.com/sunny-b/cryptkeeper/internal/config"
	"github.com/sunny-b/cryptkeeper/internal/logger"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func main() {
	var rootCmd = &cobra.Command{Use: "cryptkeeper"}

	rootCmd.AddCommand(commands.Init)
	rootCmd.AddCommand(commands.Set)
	rootCmd.AddCommand(commands.Remove)
	rootCmd.AddCommand(commands.Export)
	rootCmd.AddCommand(commands.Decrypt)
	rootCmd.AddCommand(commands.Env)
	rootCmd.AddCommand(commands.Hook)
	rootCmd.AddCommand(commands.Verify)
	rootCmd.AddCommand(commands.Direnv)

	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = "info"
	}

	logLevel, err := log.ParseLevel(logLevelStr)
	if err != nil {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
	log.SetFormatter(&logger.CustomFormatter{})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	if err := config.ReadInConfig(); err != nil {
		log.Debugf("failed to read the config file: %s", err)
	}
}
