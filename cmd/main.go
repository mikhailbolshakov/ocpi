package main

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/ocpi"
	"github.com/mikhailbolshakov/ocpi/bootstrap"
	"os"
	"os/signal"
	"syscall"
)

// @title OCPI
// @version 1.0
// @description ocpi integration API description
//
// @BasePath  /
func main() {
	// init context
	ctx := kit.NewRequestCtx().Empty().WithNewRequestId().ToContext(context.Background())

	// create a new service
	s := bootstrap.New()

	l := ocpi.L().Mth("main").Inf("created")

	// init service
	if err := s.Init(ctx); err != nil {
		l.E(err).St().Err("initialization")
		os.Exit(1)
	}

	l.Inf("initialized")

	// start listening
	if err := s.Start(ctx); err != nil {
		l.E(err).St().Err("listen")
		os.Exit(1)
	}

	l.Inf("listening")

	// handle app close
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	l.Inf("graceful shutdown")
	s.Close(ctx)
	os.Exit(0)
}
