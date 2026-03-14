package job

import (
	"errors"
	"log/slog"

	"github.com/mgumz/mtr-exporter/pkg/timeshift"
	"github.com/robfig/cron/v3"
)

type Jobs []*Job

func (jobs Jobs) Count() int { return len(jobs) }

func (jobs Jobs) Empty() bool { return len(jobs) == 0 }

func (jobs Jobs) CollectedReports() int {
	n := 0
	for _, job := range jobs {
		if !job.Report.Empty() {
			n++
		}
	}
	return n
}

func (jobs Jobs) ReSchedule(scheduler *cron.Cron, collector *Collector) error {

	scheduler.Stop()

	// step 1: clean out currently scheduled jobs
	entries := scheduler.Entries()
	ids := make([]cron.EntryID, len(entries))
	for i, entry := range entries {
		ids[i] = entry.ID
	}
	for _, id := range ids {
		scheduler.Remove(id)
	}

	// step 2: unregister previous jobs
	for i := range jobs {
		collector.RemoveJob(jobs[i].Label)
	}

	var err error

	// step 3: launch current set of jobs and register them
	// at the collector
	n := 0
	for _, j := range jobs {
		if !collector.AddJob(j.JobMeta) {
			slog.Error(errAddToCollector,
				"job.label", j.Label)
			if err != nil {
				err = errors.New(errGenericCollector)
			}
			continue
		}
		slog.Info(infoJobSchedule,
			"job.label", j.Label,
			"job.schedule", j.scheduler.spec,
			"timeshift", &j.Timeshift,
		)

		j.UpdateFn = func(meta JobMeta) bool { return collector.UpdateJob(meta) }

		s, err2 := timeshift.NewSchedule(j.Timeshift.Mode, j.scheduler.spec, j.Timeshift.Spec)

		if err2 != nil {
			slog.Error(errAddToScheduler,
				"job.label", j.Label,
				"error", err2)
			if err != nil {
				err = errors.New("schedule error")
			}
			continue
		}

		j.scheduler.entryID = scheduler.Schedule(s, j)
		j.scheduler.instance = scheduler
		n++
	}
	if n > 0 {
		slog.Info("restart scheduler", "status", "start")
		scheduler.Start()
		logSchedulerEntries(scheduler)
		slog.Info("restart scheduler", "status", "done")
	}

	return err
}
