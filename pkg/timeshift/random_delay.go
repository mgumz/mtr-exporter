package timeshift

import (
	"time"

	"github.com/robfig/cron/v3"
)

// RandomDelaySchedule takes a base schedule and
// adds a random fraction of up to `delay`
// when RandomDelaySchedule.Next() is called
type RandomDelaySchedule struct {
	baseSchedule cron.Schedule
	delay        time.Duration
}

func NewRandomDelaySchedule(base cron.Schedule, delay time.Duration) (*RandomDelaySchedule, error) {
	sched := RandomDelaySchedule{
		baseSchedule: base,
		delay:        delay,
	}
	return &sched, nil
}

func (rds *RandomDelaySchedule) Next(t time.Time) time.Time {

	nt := rds.baseSchedule.Next(t)
	if rds.delay == time.Duration(0) {
		return nt
	}
	delay := randDurationMax(rds.delay)
	nt = nt.Add(delay)

	return nt
}
