package job

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func WatchJobsFile(name, mtrBin, watchSchedule string, collector *Collector) {

	watcher := &jobFileWatch{
		name:   name,
		mtrBin: mtrBin,
		scheduler: cron.New(
			cron.WithLocation(time.UTC),
			cron.WithChain(
				cron.SkipIfStillRunning(cron.DiscardLogger),
			),
		),
		collector: collector,
	}

	watcher.Run()

	// check `name` according to `watchSchedule`
	scheduler := cron.New(
		cron.WithLocation(time.UTC),
	)

	if _, err := scheduler.AddJob(watchSchedule, watcher); err != nil {
		slog.Error(errLaunchJobWatch,
			"error", err)
		os.Exit(1)
	}
	scheduler.Start()
}

type jobFileWatch struct {
	name      string
	mtrBin    string
	scheduler *cron.Cron
	chksum    []byte
	jobs      Jobs
	collector *Collector
}

func (jw *jobFileWatch) Run() {

	slog.Info(infoStartingParse,
		"status", "starting",
		"job.file", jw.name)

	jobs, chksum, err := ParseJobFile(jw.name, jw.mtrBin)
	if err != nil {
		slog.Warn(
			"parsing failed",
			"job.file", jw.name,
			"error", err,
		)
		return
	}

	slog.Info(infoDoneParse,
		"status", "done",
		"njobs", len(jobs),
		"job.file", jw.name,
		"job.file.prev-sha256", fmt.Sprintf("%x", string(jw.chksum)),
		"job.file.sha256", fmt.Sprintf("%x", string(chksum)),
	)

	if bytes.Equal(jw.chksum, chksum) {
		slog.Info(infoJobFileUnchanged,
			"job.file", jw.name)
		return
	}

	slog.Info(infoJobFileChanged,
		"njobs", len(jobs),
		"job.file", jw.name)

	_ = jobs.ReSchedule(jw.scheduler, jw.collector)

	jw.jobs = jobs
	jw.chksum = chksum
}
