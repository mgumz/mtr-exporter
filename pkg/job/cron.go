package job

import (
	"log"
)

// cron.v3 interface
func (job *Job) Run() {

	log.Printf("info: %q launching via %q", job.Label, job.cmdLine)
	if err := job.Launch(); err != nil {
		log.Printf("info: %q failed: %s", job.Label, err)
		return
	}
	log.Printf("info: %q done: %d hops in %s.", job.Label,
		len(job.Report.Hubs), job.Duration)
}
