package main

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tmw/exploring-tilt/internal/realtime"
	todoservice "github.com/tmw/exploring-tilt/internal/todo"
)

const (
	defaultServerPort     = "9090"
	defaultNatsServerAddr = "nats://nats:4222"
)

type config struct {
	ServerPort     string
	NatsServerAddr string
}

func getConfig() config {
	return config{
		ServerPort:     cmp.Or(os.Getenv("SERVER_PORT"), defaultServerPort),
		NatsServerAddr: cmp.Or(os.Getenv("NATS_SERVER_ADDR"), defaultNatsServerAddr),
	}
}

func run() error {
	config := getConfig()
	conn, err := todoservice.NewNatsConnection(config.NatsServerAddr)
	if err != nil {
		return fmt.Errorf("error connecting to nats: %w", err)
	}
	defer conn.Drain()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	subscriber, err := realtime.NewSubscriber(conn, "todos", "todos", logger)
	if err != nil {
		return fmt.Errorf("error setting up subscriber: %w", err)
	}

	go subscriber.Start(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", eventHandler(logger, subscriber))

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", config.ServerPort),
		Handler: mux,

		ReadTimeout:  30 * time.Second,
		WriteTimeout: 0,
		IdleTimeout:  10 * time.Minute,
	}

	go func() {
		logger.Info("server started", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil {
			logger.Error("error while binding server", slog.Any("error", err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received, shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("graceful shutdown complete, bye")

	return nil
}
func eventHandler(logger *slog.Logger, sub *realtime.Subscriber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// required headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		updateCh, unsubscribe := sub.Subscribe(r.Context(), "todos")
		defer unsubscribe()

		heartbeat := time.NewTicker(10 * time.Second)
		defer heartbeat.Stop()

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		for {
			select {
			case <-r.Context().Done():
				logger.Info("client hung up", slog.String("remoteAddr", r.RemoteAddr))
				return

			case t := <-heartbeat.C:
				msg := fmt.Sprintf("data: {\"kind\": \"heartbeat\", \"at\":\"%s\"}\n\n", t.Format(time.RFC3339))
				_, err := w.Write([]byte(msg))
				if err != nil {
					logger.Error("error while writing heartbeat to stream", slog.Any("error", err))
					return
				}

				flusher.Flush()

			case u := <-updateCh:
				msg := fmt.Sprintf("data: {\"kind\": \"update\", \"at\":\"%s\"}\n\n", u.When.Format(time.RFC3339))
				_, err := w.Write([]byte(msg))
				if err != nil {
					logger.Error("error while writing to stream", slog.Any("error", err))
					return
				}

				flusher.Flush()
			}
		}
	}
}

func main() {
	if err := run(); err != nil {
		slog.Error("error while running server", slog.Any("error", err))
		os.Exit(1)
	}
}
