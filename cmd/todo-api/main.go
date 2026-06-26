package main

import (
	"cmp"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	todoservice "github.com/tmw/exploring-tilt/internal/todo"
	todosservice "github.com/tmw/exploring-tilt/internal/todo"
	todoapi "github.com/tmw/exploring-tilt/internal/todo-api"
)

const (
	defaultStoragePath = "/storage/items.db"
	defaultServerPort  = "9191"
)

type config struct {
	StoragePath string
	ServerPort  string
}

func getConfig() config {
	return config{
		StoragePath: cmp.Or(os.Getenv("STORAGE_PATH"), defaultStoragePath),
		ServerPort:  cmp.Or(os.Getenv("SERVER_PORT"), defaultServerPort),
	}
}

func run() error {
	config := getConfig()
	conn, err := todoservice.NewNatsConnection("nats://nats:4222")
	if err != nil {
		return fmt.Errorf("error connecting to nats: %w", err)
	}
	defer conn.Drain()

	persistance, err := todoservice.NewNatsPersistance(conn, "todos")
	if err != nil {
		panic(err)
	}
	svc, err := todosservice.NewWithPersistance(persistance)
	if err != nil {
		return fmt.Errorf("error setting up service: %w", err)
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	api := todoapi.New(svc, logger)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", config.ServerPort),
		Handler: api.Mux(),

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	logger.Info("server started", slog.String("addr", server.Addr))
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("error while binding server: %w", err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		slog.Error("error while running server", slog.Any("error", err))
		os.Exit(1)
	}
}
