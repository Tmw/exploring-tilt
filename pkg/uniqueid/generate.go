package uniqueid

import (
	"crypto/rand"
	"encoding/hex"
)

func Generate(length int) string {
	buf := make([]byte, length)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}
