package timeshift

import (
	"crypto/rand"
	"math/big"
	"time"
)

func randInt64N(max int64) int64 {
	m := big.Int{}
	m.SetInt64(max)
	n, _ := rand.Int(rand.Reader, &m)
	return n.Int64()
}

func randDurationMax(max time.Duration) time.Duration {
	return time.Duration(randInt64N(int64(max)))
}
