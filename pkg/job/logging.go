package job

import (
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

func logSchedulerEntries(scheduler *cron.Cron) {

	entries := scheduler.Entries()

	slog.Info("active schedule",
		"njobs", len(entries),
	)

	for _, e := range entries {
		j, ok := e.Job.(*Job)
		if !ok {
			continue
		}
		now := time.Now().UTC()
		nrun := e.Next.Format(time.RFC3339)

		slog.Info("planned job",
			"time.when", nrun,
			"job.label", j.Label,
			"job.schedule", j.scheduler.spec,
			"time.in", e.Next.Sub(now).Truncate(time.Second),
			"timeshift", &j.Timeshift,
		)
	}
}
