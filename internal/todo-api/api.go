package todoapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	todo "github.com/tmw/exploring-tilt/internal/todo"
	"github.com/tmw/exploring-tilt/pkg/httphelper"
	"github.com/tmw/exploring-tilt/pkg/middleware"
	"github.com/tmw/exploring-tilt/pkg/middleware/middlewares"
)

type Service interface {
	Create(params todo.CreateTodoParams) []todo.Todo
	List() []todo.Todo
	Delete(id string) ([]todo.Todo, error)
	Toggle(id string, params todo.ToggleTodoParams) ([]todo.Todo, error)
}

type Api struct {
	svc    Service
	logger *slog.Logger
}

func New(svc Service, logger *slog.Logger) *Api {
	return &Api{
		svc:    svc,
		logger: logger,
	}
}

func (a *Api) HandleGetTodos(w http.ResponseWriter, r *http.Request) {
	todos := a.svc.List()
	json.NewEncoder(w).Encode(todos)
}

func (a *Api) HandleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var params todo.CreateTodoParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		a.logger.Warn("error decoding JSON body", slog.Any("error", err))
		return
	}

	todos := a.svc.Create(params)
	json.NewEncoder(w).Encode(todos)
}

func (a *Api) HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	todos, err := a.svc.Delete(id)

	if err != nil {
		if errors.Is(err, todo.ErrItemNotFound) {
			httphelper.RespondError(w, 404, "item not found")
			return
		}

		a.logger.Error("error deleting todo", slog.Any("error", err))
		httphelper.RespondError(w, 500, "internal server error")
		return
	}

	json.NewEncoder(w).Encode(todos)
}

func (a *Api) HandleToggleTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var params todo.ToggleTodoParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		a.logger.Warn("error decoding JSON body", slog.Any("error", err))
		return
	}

	todos, err := a.svc.Toggle(id, params)
	if err != nil {
		if errors.Is(err, todo.ErrItemNotFound) {
			httphelper.RespondError(w, 404, "item not found")
			return
		}

		a.logger.Error("error toggling todo", slog.Any("error", err))
		httphelper.RespondError(w, 500, "internal server error")
		return
	}

	json.NewEncoder(w).Encode(todos)
}

func (a *Api) Mux() http.Handler {
	mux := http.NewServeMux()
	mw := middleware.New(
		middlewares.Cors(
			middlewares.CorsConfig{
				Origin: "http://localhost:9090",
				Methods: []string{
					http.MethodGet,
					http.MethodPatch,
					http.MethodPost,
					http.MethodDelete,
					http.MethodOptions,
				},
				Headers: []string{
					"Content-Type",
				},
			},
		),
		middlewares.JSON(),
	)

	mux.HandleFunc("GET /todos", a.HandleGetTodos)
	mux.HandleFunc("POST /todos", a.HandleCreateTodo)
	mux.HandleFunc("DELETE /todos/{id}", a.HandleDeleteTodo)
	mux.HandleFunc("PATCH /todos/{id}/status", a.HandleToggleTodo)

	return mw.Wrap(mux)
}
