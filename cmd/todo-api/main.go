package main

import (
	"cmp"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

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

func main() {
	config := getConfig()
	persistance := todosservice.NewFilePersitanceWithPath(config.StoragePath)
	svc, err := todosservice.NewWithPersistance(persistance)
	if err != nil {
		slog.Error("error setting up service", slog.Any("error", err))
		os.Exit(1)
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
		logger.Error("error while binding server", slog.Any("error", err))
	}
}
