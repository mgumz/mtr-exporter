package timeshift

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Mode int

const (
	None Mode = iota
	RandomDeviation
	RandomDelay
)

func NewSchedule(mode Mode, spec, timeshift string) (cron.Schedule, error) {

	sched, err := cron.ParseStandard(spec)
	if err != nil {
		return nil, err
	}

	switch mode {
	case None:
		return sched, nil

	case RandomDeviation:
		deviation, err := time.ParseDuration(timeshift)
		if err != nil {
			return nil, fmt.Errorf(errUnparseableDeviation, timeshift)
		}
		rds, err := NewRandomDeviationSchedule(sched, deviation)
		return rds, err

	case RandomDelay:
		delay, err := time.ParseDuration(timeshift)
		if err != nil {
			return nil, fmt.Errorf(errUnparseableDelay, timeshift)
		}
		rds, err := NewRandomDelaySchedule(sched, delay)
		return rds, err
	}

	return nil, fmt.Errorf(errUnknownMode, mode)
}
