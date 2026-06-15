package uniqueid

import (
	ulid "github.com/oklog/ulid/v2"
)

func Generate(length int) string {
	id := ulid.Make()
	return id.String()
}
