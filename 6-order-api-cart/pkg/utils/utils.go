package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomCode() uint16 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return uint16(r.Intn(9000) + 1000)
}
