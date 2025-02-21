package mqtt

import (
	. "crypto/rand"
	"encoding/hex"
	"golang.org/x/exp/rand"
	"time"
)

func randomClientID() string {
	b := make([]byte, 8)
	_, _ = Read(b)
	return "go_mqtt_" + hex.EncodeToString(b)
}
func generateRandom10DigitNumber() int64 {
	rand.Seed(uint64(time.Now().UnixNano()))
	return rand.Int63n(1e10)
}
