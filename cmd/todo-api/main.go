package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	todosservice "github.com/tmw/exploring-tilt/internal/todo-service"
	"github.com/tmw/exploring-tilt/pkg/httphelper"
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
	mux := http.NewServeMux()
	config := getConfig()
	persistance := todosservice.NewFilePersitanceWithPath(config.StoragePath)
	svc, err := todosservice.NewWithPersistance(persistance)
	if err != nil {
		slog.Error("error setting up service", slog.Any("error", err))
		os.Exit(1)
	}

	mux.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		todos := svc.List()
		json.NewEncoder(w).Encode(todos)
	})

	mux.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		var params todosservice.CreateTodoParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			slog.Warn("error decoding JSON body", slog.Any("error", err))
			return
		}

		todos := svc.Create(params)
		json.NewEncoder(w).Encode(todos)
	})

	mux.HandleFunc("DELETE /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		todos, err := svc.Delete(id)

		if err != nil {
			if errors.Is(err, todosservice.ErrItemNotFound) {
				httphelper.RespondError(w, 404, "item not found")
				return
			}

			slog.Error("error deleting todo", slog.Any("error", err))
			httphelper.RespondError(w, 500, "internal server error")
			return
		}

		json.NewEncoder(w).Encode(todos)
	})

	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		version := "0.0.2"
		fmt.Fprintf(w, "{\"version\": \"%s\"}", version)
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", config.ServerPort),
		Handler: mux,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	slog.Info("server started", slog.String("addr", server.Addr))

	if err := server.ListenAndServe(); err != nil {
		slog.Error("error while binding server", slog.Any("error", err))
	}
}
