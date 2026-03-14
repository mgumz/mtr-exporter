package timeshift

import (
	"log"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

func TestRandomOffsetSchedule(t *testing.T) {

	baseSched, _ := cron.ParseStandard("@midnight")
	rds, _ := NewRandomDeviationSchedule(baseSched, time.Duration(1*time.Hour))

	n := time.Date(2025, 05, 01, 15, 16, 17, 00, time.UTC)

	for range 50 {
		nt := rds.Next(n)
		log.Printf("%v -> %v", n, nt)
	}

}
