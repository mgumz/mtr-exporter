package job

import (
	"fmt"
	"log"
)

// cron.v3 interface
func (job *Job) Run() {

	log.Printf("info: %q launching via %q", job.Label, job.cmdLine)
	if err := job.Launch(); err != nil {
		log.Printf("info: %q failed: %s", job.Label, err)
		return
	}

	errMsg := ""
	if job.Report.ErrorMsg != "" {
		errMsg = fmt.Sprintf("(err: %q)", job.Report.ErrorMsg)
	}
	log.Printf("info: %q done%s: %d hops in %s.", job.Label,
		errMsg, len(job.Report.Hubs), job.Duration)
}
