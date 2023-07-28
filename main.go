package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"example.com/be_test/internal"
	"example.com/be_test/pkg/env"
	"example.com/be_test/pkg/logger"
)

func main() {
	opts := internal.Options{
		ListenAddressHTTP: env.MustGet("LISTEN_ADDRESS_HTTP"),
		Production:        env.GetBool("PRODUCTION", false),
		LogQuery:          env.GetBool("LOG_QUERY", false),
		DBURL:             env.MustGet("DB_URL"),
		JWTSignKey:        env.MustGet("JWT_SIGN_KEY"),
	}

	log := logger.New()

	ctx := context.Background()

	svc, err := internal.New(ctx, opts, log)
	if err != nil {
		log.Fatal(err)
	}

	err = svc.Run()
	if err != nil {
		log.Fatal(err)
	}

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	log.Debug("Waiting for the signal...")
	sig := <-signalChannel

	// Handle the signal
	log.Infof("Received signal: %v\n", sig)

	svc.Shutdown(ctx)

	os.Exit(0)
}
