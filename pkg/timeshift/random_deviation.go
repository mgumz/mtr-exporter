package timeshift

import (
	"time"

	"github.com/robfig/cron/v3"
)

// NOTE: special situation: in case timeshift.RandomDeviation scheduler
// is used, it might happen that the first execution is _before_ the
// original schedule. As a result, cron calculates the .Next run
// based on this schedule. Consequently, the timeshift.RandomDeviation
// schedule launches `job` 1+ too often:
//
// configured: sched ±deviation
//
//  |<-- deviation ------ | ------ deviation --> |
// now .a. .b. .c. .d.  sched ..
//
// There is a chance, that `job` gets executed on each .*. point in time.
// On .a., cron executes `job` (async) and calculates .Next(). The next
// regulara launch time is the same as previously. The deviation value is added
// or subtracted. If .b. is now also _before_ the regular planned launch.
//
// Real life observation:
//
// time=2025-05-07T19:41:45.000Z level=INFO msg="watching -jobs-file" fileName=/config/speedtest-jobs.txt schedule="@every 1m"
// time=2025-05-07T19:41:45.001Z level=INFO msg="watched jobFile has changed, scheduling jobs" fileName=/config/speedtest-jobs.txt numberJobs=1
// time=2025-05-07T19:41:45.001Z level=INFO msg="schedule job" job=example schedule="15 2 * * *" timeshift.mode=deviation timeshift.by=15m
// time=2025-05-07T19:41:45.001Z level=INFO msg="restart scheduler" njobs=1
// time=2025-05-07T19:41:45.001Z level=INFO msg="next run" job=abc duration=6h30m27s
// time=2025-05-07T19:41:45.001Z level=INFO msg="serving ..." path=/metrics bindAddr=:8080
// time=2025-05-08T02:12:12.333Z level=INFO msg=launching job=example cmdline="speedtest-go --json --unit decimal-bytes -s 00001"
// time=2025-05-08T02:12:37.761Z level=INFO msg="job done" job=example duration=25.42742238s
// time=2025-05-08T02:12:37.761Z level=INFO msg="next launch" job=example nextRun=12m46s
// time=2025-05-08T02:25:24.237Z level=INFO msg=launching job=example cmdline="speedtest-go --json --unit decimal-bytes -s 00001"
// time=2025-05-08T02:25:49.448Z level=INFO msg="job done" job=example duration=25.210635636s
// time=2025-05-08T02:25:49.448Z level=INFO msg="next launch" job=example nextRun=23h56m4s
//
// job `example` is scheduled to be executed ±15 minutes around 02:15 on
// each day. It got executed on 02:12:12 (effective
// timeshift.RandomDeviation of 3min). The speedtest run took ~ 25s. Since
// 02:12:12 is _before_ 02:15 … another run calculated and, randomly
// 02:25:24 was picked.
// To prevent this, we wait until the regular planned launch + deviation has passed.
// cron.SkipIfStillRunning() ensures that cron does not launch `job`
// "early"
//
// Not affected by this:
// * regular scheduler
// * timeshift.RandomDelay
//
// NOT in current use until a nice way of dealing with above situation can be
// found

type RandomDeviationSchedule struct {
	baseSchedule cron.Schedule
	deviation    time.Duration
}

func (rds *RandomDeviationSchedule) Range(t time.Time) (time.Time, time.Time) {
	// t = rds.baseSchedule.Next(t)
	return t.Add(-rds.deviation), t.Add(rds.deviation)
}

func (rds *RandomDeviationSchedule) UnshiftedNext(t time.Time) time.Time {
	return rds.baseSchedule.Next(t)
}

func NewRandomDeviationSchedule(base cron.Schedule, deviation time.Duration) (*RandomDeviationSchedule, error) {
	sched := RandomDeviationSchedule{
		baseSchedule: base,
		deviation:    deviation,
	}
	return &sched, nil
}

func (rds *RandomDeviationSchedule) Next(t time.Time) time.Time {

	nt := rds.baseSchedule.Next(t)

	if rds.deviation == time.Duration(0) {
		return nt
	}

	off := 2 * randDurationMax(rds.deviation)
	nt = nt.Add(-rds.deviation + off)

	return nt
}
