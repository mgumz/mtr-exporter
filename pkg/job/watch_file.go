package job

import (
	"bytes"
	"log"
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
		log.Printf("error: unable to launch watch-jobs scheduler: %v", err)
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

	log.Printf("info: starting to parse %q", jw.name)
	jobs, chksum, err := ParseJobFile(jw.name, jw.mtrBin)
	if err != nil {
		log.Printf("warning: parsing %q: %v", jw.name, err)
		return
	}
	log.Printf("info: done parsing %q: %d jobs (previous-sha256:%x|current-sha256:%x) ", jw.name, len(jobs), jw.chksum, chksum)

	if bytes.Equal(jw.chksum, chksum) {
		log.Printf("info: watched file is unchanged")
		return
	}

	log.Printf("info: watched file changed: replacing %d jobs", len(jobs))

	jobs.ReSchedule(jw.scheduler, jw.collector)

	jw.jobs = jobs
	jw.chksum = chksum
}
