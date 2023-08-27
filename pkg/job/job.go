package job

import (
	"bytes"
	"os/exec"
	"strings"
	"time"

	"github.com/mgumz/mtr-exporter/pkg/mtr"
)

type JobMeta struct {
	Report   mtr.Report
	Launched time.Time
	Duration time.Duration
	Schedule string
	Label    string
	CmdLine  string
}

type Job struct {
	JobMeta

	mtrBinary string
	args      []string
	cmdLine   string

	UpdateFn func(JobMeta) bool
}

func NewJob(mtr string, args []string, schedule string) *Job {
	extra := []string{
		"-j", // json output
	}
	args = append(extra, args...)
	cmd := exec.Command(mtr, args...)
	job := Job{
		args:      args,
		mtrBinary: mtr,
		cmdLine:   strings.Join(cmd.Args, " "),
	}
	job.JobMeta.Schedule = schedule
	job.JobMeta.CmdLine = job.cmdLine
	return &job
}

func (job *Job) Launch() error {

	// TODO: maybe use CommandContext to have an upper limit in the execution

	cmd := exec.Command(job.mtrBinary, job.args...)

	// launch mtr
	buf := bytes.Buffer{}
	cmd.Stdout = &buf
	launched := time.Now()
	if err := cmd.Run(); err != nil {
		return err
	}
	duration := time.Since(launched)

	// decode the report
	report := mtr.Report{}
	if err := report.Decode(&buf); err != nil {
		return err
	}

	// copy the report into the job
	job.Report = report
	job.Launched = launched
	job.Duration = duration

	if job.UpdateFn != nil {
		job.UpdateFn(job.JobMeta)
	}

	// done.
	return nil
}
