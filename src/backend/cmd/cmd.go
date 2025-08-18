package cmd

import (
	"context"
	"fmt"
	"local/client"
	"local/config"
	"local/endpoint"
	"local/model"
	"local/repository"
	"local/service/common"
	"local/service/initial"
	httpTransoprt "local/transport/http"
	"log"
	http "net/http"
	"os"
	"os/signal"
	"syscall"
)

var ServiceName = "simple-chat-service"

func Run() {
	config.LoadConfig()
	ctx := context.Background()
	defer ctx.Done()

	// Initialize MySQL repository
	repo, err := repository.NewRepository()
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	initParams := &model.InitParams{
		ServiceName: ServiceName,
		Ctx:    ctx,
	}
	clt := client.NewClient(initParams)

	// Run the service
	svc := initial.NewService(&common.Params{
		Repo:   repo,
		Client: clt,
	})

	endpoints := endpoint.NewEndpoints(&svc)

	runServer(initParams, endpoints)
}

func runServer(initParams *model.InitParams, endpoints *endpoint.Endpoints) {
	svr := httpTransoprt.MakeHttpTransport(initParams, endpoints)
	log.Printf("HTTP server listening on %s", fmt.Sprintf(":%d", config.Config.HTTPPort))

	errCh := make(chan error, 1)

	go func() {
		errCh <- http.ListenAndServe(fmt.Sprintf(":%d", config.Config.HTTPPort), svr)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		log.Printf("Received shutdown signal. Stopping server...")
		os.Exit(0)
	case err := <-errCh:
		log.Printf("server stopped with error: %v", err)
		os.Exit(1)
	}
}
