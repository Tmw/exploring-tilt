package realtime

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Subscriber struct {
	js      jetstream.JetStream
	kv      jetstream.KeyValue
	logger  *slog.Logger
	keyName string

	clients   map[chan Update]struct{}
	clientsMu sync.Mutex
}

type Update struct {
	When time.Time
}

func NewSubscriber(conn *nats.Conn, bucketName, keyName string, logger *slog.Logger) (*Subscriber, error) {
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, fmt.Errorf("error setting up jetstream from connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	kv, _ := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket: bucketName,
	})

	return &Subscriber{
		js:      js,
		kv:      kv,
		logger:  logger,
		keyName: keyName,
		clients: make(map[chan Update]struct{}),
	}, nil
}

func (s *Subscriber) fanout(u Update) int {
	numClientsInformed := 0
	for ch, _ := range s.clients {
		ch <- u
		numClientsInformed++
	}

	return numClientsInformed
}

func (s *Subscriber) Start(ctx context.Context) error {
	watcher, err := s.kv.Watch(ctx, s.keyName, jetstream.UpdatesOnly())
	if err != nil {
		return fmt.Errorf("error setting up watcher: %w", err)
	}
	defer watcher.Stop()

	for update := range watcher.Updates() {
		if update == nil {
			continue
		}

		numClientsInformed := s.fanout(Update{When: update.Created()})
		s.logger.Info(
			"received update",
			slog.Time("createdAt", update.Created()),
			slog.Int("numClientsInformed", numClientsInformed),
		)
	}

	s.logger.Info("subscription channel closed")
	return nil
}

func (s *Subscriber) Subscribe(ctx context.Context, keyName string) (channel <-chan Update, unsubFn func()) {
	clientChan := make(chan Update)
	s.clientsMu.Lock()
	s.clients[clientChan] = struct{}{}
	s.clientsMu.Unlock()

	s.logger.Info(
		"client subscribed",
		slog.Int("numSubscribers", len(s.clients)),
	)

	unsub := func() {
		s.clientsMu.Lock()
		delete(s.clients, clientChan)
		s.clientsMu.Unlock()

		s.logger.Info(
			"client unsubscribed",
			slog.Int("numSubscribers", len(s.clients)),
		)
	}

	return clientChan, unsub
}
