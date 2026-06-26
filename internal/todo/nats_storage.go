package todoservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsPersistance struct {
	js jetstream.JetStream
	kv jetstream.KeyValue
}

var (
	keyName = "todos"
)

func NewNatsConnection(url string) (*nats.Conn, error) {
	return nats.Connect(url)
}

func NewNatsPersistance(conn *nats.Conn, bucketName string) (*NatsPersistance, error) {
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, fmt.Errorf("error getting jetstream from nats conn: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	kv, _ := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket: bucketName,
	})

	return &NatsPersistance{
		js: js,
		kv: kv,
	}, nil
}

func (p *NatsPersistance) Load() ([]Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	entry, err := p.kv.Get(ctx, keyName)
	if err != nil {
		// key not found means we never stored the key yet,
		// that's OK - just return empty slice.
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return []Todo{}, nil
		}

		// if that wasn't it, propagate.
		return nil, fmt.Errorf("error fetching todos: %w", err)
	}

	var out []Todo
	if err = json.Unmarshal(entry.Value(), &out); err != nil {
		return nil, fmt.Errorf("error unmarshalling: %w", err)
	}

	return out, nil
}

func (p *NatsPersistance) Store(items []Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payload, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("error marshalling todo payload: %w", err)
	}

	_, err = p.kv.Put(ctx, keyName, payload)
	if err != nil {
		return fmt.Errorf("error storing todos: %w", err)
	}

	return nil
}
