package todoservice

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/tmw/exploring-tilt/pkg/uniqueid"
)

const idLength = 16

type Persistence interface {
	Store([]Todo) error
	Load() ([]Todo, error)
}

type Service struct {
	items   []Todo
	itemsMu sync.RWMutex

	persistance Persistence
}

func NewWithPersistance(p Persistence) (*Service, error) {
	items, err := p.Load()
	if err != nil {
		return nil, fmt.Errorf("unable to read todos from persistance: %w", err)
	}

	return &Service{
		items:       items,
		persistance: p,
	}, nil
}

func (s *Service) flush() {
	if err := s.persistance.Store(s.items); err != nil {
		slog.Error("error flushing to persistance", slog.Any("error", err))
	}
}

type CreateTodoParams struct {
	Title       string     `json:"title"`
	CompletedAt *time.Time `json:"compeltedAt"`
}

func (s *Service) Create(params CreateTodoParams) []Todo {
	s.itemsMu.Lock()
	defer s.itemsMu.Unlock()
	s.items = append(s.items, Todo{
		ID:          uniqueid.Generate(idLength),
		Title:       params.Title,
		CompletedAt: params.CompletedAt,
		CreatedAt:   time.Now().UTC(),
	})

	s.flush()
	return s.items
}

type ToggleTodoParams struct {
	NewState bool `json:"newState"`
}

func now() *time.Time {
	n := time.Now()
	return &n
}

func (s *Service) Toggle(id string, params ToggleTodoParams) ([]Todo, error) {
	s.itemsMu.Lock()
	defer s.itemsMu.Unlock()

	found := false

	for idx := range s.items {
		item := s.items[idx]
		if item.ID == id {
			found = true
			if params.NewState == true {
				s.items[idx].CompletedAt = now()
			} else {
				s.items[idx].CompletedAt = nil
			}
			break
		}
	}

	if !found {
		return s.items, ErrItemNotFound
	}

	s.flush()
	return s.items, nil
}

func (s *Service) List() []Todo {
	s.itemsMu.RLock()
	defer s.itemsMu.RUnlock()
	return s.items
}

var (
	ErrItemNotFound = errors.New("item not found")
)

func (s *Service) Delete(id string) ([]Todo, error) {
	s.itemsMu.Lock()
	defer s.itemsMu.Unlock()

	idx := slices.IndexFunc(s.items, func(todo Todo) bool {
		return todo.ID == id
	})

	if idx == -1 {
		return s.items, ErrItemNotFound
	}

	s.items = append(s.items[0:idx], s.items[idx+1:]...)
	s.flush()
	return s.items, nil
}
