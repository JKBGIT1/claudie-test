package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/berops/claudie/internal/envs"
	"github.com/berops/claudie/internal/healthcheck"
	"github.com/berops/claudie/internal/utils"
	inboundAdapters "github.com/berops/claudie/services/frontend/adapters/inbound"
	outboundAdapters "github.com/berops/claudie/services/frontend/adapters/outbound"
	"github.com/berops/claudie/services/frontend/domain/usecases"
)

const (
	// healthcheckPort is the port on which Kubernetes readiness and liveness probes send request
	// for performing health checks.
	healthcheckPort = 50058
)

func main() {
	utils.InitLog("frontend")

	if err := run(); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func run() error {
	contextBoxConnector := outboundAdapters.NewContextBoxConnector(envs.ContextBoxURL)
	err := contextBoxConnector.Connect()
	if err != nil {
		return err
	}

	usecaseContext, usecaseCancel := context.WithCancel(context.Background())
	usecases := &usecases.Usecases{
		ContextBox:    contextBoxConnector,
		SaveChannel:   make(chan *usecases.RawManifest),
		DeleteChannel: make(chan *usecases.RawManifest),
		Context:       usecaseContext,
	}

	secretWatcher, err := inboundAdapters.NewSecretWatcher(usecases)
	if err != nil {
		usecaseCancel()
		return err
	}

	// Start Kubernetes liveness and readiness probe responders
	healthcheck.NewClientHealthChecker(fmt.Sprint(healthcheckPort),
		func() error {
			if err := secretWatcher.PerformHealthCheck(); err != nil {
				return err
			}
			return contextBoxConnector.PerformHealthCheck()
		},
	).StartProbes()

	// usecases.ProcessManifestFiles() goroutine returns on usecases.Context cancels
	go usecases.ProcessManifestFiles()
	log.Info().Msgf("Frontend is ready to process input manifests")

	// usecases.WatchConfigs() goroutine returns on usecases.Context cancels
	go usecases.WatchConfigs()
	log.Info().Msgf("Frontend is ready to watch input manifest statuses")

	// secretWatcher.Monitor() goroutine returns on usecases.Context cancels
	go secretWatcher.Monitor()
	log.Info().Msgf("Frontend is watching for any new input manifest")

	// Cancel context for usecases functions to terminate goroutines.
	defer usecaseCancel()

	// Interrupt signal listener
	shutdownSignalChan := make(chan os.Signal, 1)
	signal.Notify(shutdownSignalChan, os.Interrupt, syscall.SIGTERM)
	sig := <-shutdownSignalChan

	log.Info().Msgf("Received program shutdown signal %v", sig)

	// Disconnect from context-box
	if err := contextBoxConnector.Disconnect(); err != nil {
		log.Err(err).Msgf("Failed to gracefully shutdown ContextBoxConnector")
	}
	defer signal.Stop(shutdownSignalChan)

	return fmt.Errorf("program interrupt signal")
}