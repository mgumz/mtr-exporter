package job

import (
	"errors"
	"log"

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
	entries := []cron.EntryID{}
	for _, entry := range scheduler.Entries() {
		entries = append(entries, entry.ID)
	}
	for _, entry := range entries {
		scheduler.Remove(entry)
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
			log.Printf("error: unable to add job %q to collector", j.Label)
			if err != nil {
				err = errors.New("collector error")
			}
			continue
		}
		log.Printf("info: schedule %q to %q", j.Label, j.Schedule)
		j.UpdateFn = func(meta JobMeta) bool { return collector.UpdateJob(meta) }
		if _, err2 := scheduler.AddJob(j.Schedule, j); err2 != nil {
			log.Printf("error: unable to add %q to scheduler: %v", j.Label, err2)
			if err != nil {
				err = errors.New("schedule error")
			}
		}
		n++
	}
	if n > 0 {
		log.Printf("info: restart scheduler (%d)", n)
		scheduler.Start()
	}

	return err
}
