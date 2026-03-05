package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ReggieJTech/SuperAgent/internal/agent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Build information (set via ldflags)
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func main() {
	// Command-line flags
	configFile := flag.String("config", "/etc/bigpanda-agent/config.yaml", "Path to configuration file")
	showVersion := flag.Bool("version", false, "Show version information")
	logLevel := flag.String("log-level", "", "Log level (trace, debug, info, warn, error, fatal)")
	validateConfig := flag.Bool("validate", false, "Validate configuration and exit")
	flag.Parse()

	// Show version and exit
	if *showVersion {
		fmt.Printf("BigPanda Super Agent\n")
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Build Time: %s\n", BuildTime)
		os.Exit(0)
	}

	// Initialize logger
	setupLogger(*logLevel)

	log.Info().
		Str("version", Version).
		Str("commit", GitCommit).
		Str("build_time", BuildTime).
		Msg("Starting BigPanda Super Agent")

	// Load configuration
	cfg, err := agent.LoadConfig(*configFile)
	if err != nil {
		log.Fatal().Err(err).Str("config_file", *configFile).Msg("Failed to load configuration")
	}

	// Validate config and exit if requested
	if *validateConfig {
		log.Info().Msg("Configuration is valid")
		os.Exit(0)
	}

	// Override log level from config if not specified on command line
	if *logLevel == "" && cfg.Logging.Level != "" {
		setLogLevel(cfg.Logging.Level)
	}

	// Create agent instance
	ag, err := agent.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create agent")
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Start agent
	if err := ag.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to start agent")
	}

	log.Info().Msg("Agent started successfully")

	// Wait for shutdown signal
	for {
		sig := <-sigChan
		log.Info().Str("signal", sig.String()).Msg("Received signal")

		switch sig {
		case syscall.SIGHUP:
			// Reload configuration
			log.Info().Msg("Reloading configuration")
			if err := ag.Reload(*configFile); err != nil {
				log.Error().Err(err).Msg("Failed to reload configuration")
			} else {
				log.Info().Msg("Configuration reloaded successfully")
			}

		case syscall.SIGINT, syscall.SIGTERM:
			// Graceful shutdown
			log.Info().Msg("Initiating graceful shutdown")

			// Create shutdown context with timeout
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer shutdownCancel()

			// Stop agent
			if err := ag.Stop(shutdownCtx); err != nil {
				log.Error().Err(err).Msg("Error during shutdown")
				os.Exit(1)
			}

			log.Info().Msg("Agent stopped successfully")
			return
		}
	}
}

// setupLogger initializes the global logger
func setupLogger(level string) {
	// Use console output for now (will be configurable later)
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()

	// Set log level
	if level != "" {
		setLogLevel(level)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// setLogLevel sets the global log level
func setLogLevel(level string) {
	switch level {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	default:
		log.Warn().Str("level", level).Msg("Unknown log level, using info")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
