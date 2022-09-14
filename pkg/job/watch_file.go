package job

import (
	"bytes"
	"log"
	"os"

	"github.com/robfig/cron/v3"
)

func WatchJobsFile(name, mtrBin, watchSchedule string, collector *Collector) {

	watcher := &jobFileWatch{name: name,
		mtrBin:    mtrBin,
		scheduler: cron.New(),
		collector: collector,
	}

	watcher.Run()

	// check `name` according to `watchSchedule`
	scheduler := cron.New()
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

	jw.scheduler.Stop()

	// step 1: clean out currently scheduled jobs
	entries := []cron.EntryID{}
	for _, entry := range jw.scheduler.Entries() {
		entries = append(entries, entry.ID)
	}
	for _, entry := range entries {
		jw.scheduler.Remove(entry)
	}

	// step 2: unregister previous jobs
	for i := range jw.jobs {
		jw.collector.RemoveJob(jw.jobs[i].Label)
	}

	// step 3: launch current set of jobs and register them
	// at the collector
	n := 0
	for _, j := range jobs {
		if jw.collector.AddJob(j.JobMeta) {
			log.Printf("info: schedule %q to %q", j.Label, j.Schedule)
			j.UpdateFn = func(meta JobMeta) bool { return jw.collector.UpdateJob(meta) }
			if _, err := jw.scheduler.AddJob(j.Schedule, j); err != nil {
				log.Printf("error: unable to add %q to scheduler: %v", j.Label, err)
			}
			n++
		} else {
			log.Printf("error: unable to add job %q to collector", j.Label)
		}
	}
	if n > 0 {
		log.Printf("info: restart scheduler (%d)", n)
		jw.scheduler.Start()
	}

	jw.jobs = jobs
	jw.chksum = chksum
}
