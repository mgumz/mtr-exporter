package main

import (
	"bytes"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type mtrJob struct {
	Report   *mtrReport
	Launched time.Time
	Duration time.Duration

	mtrBinary string
	args      []string
	cmdLine   string

	sync.Mutex
}

func newMtrJob(mtr string,  args []string) *mtrJob {
	extra := []string{
		"-j", // json output
	}
	args = append(extra, args...)
	cmd := exec.Command(mtr, args...)

	return &mtrJob{
		args:      args,
		mtrBinary: mtr,
		cmdLine:   strings.Join(cmd.Args, " "),
	}
}

func (job *mtrJob) Launch() error {

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
	report := &mtrReport{}
	if err := report.Decode(&buf); err != nil {
		return err
	}

	// copy the report into the job
	job.Lock()
	job.Report = report
	job.Launched = launched
	job.Duration = duration
	job.Unlock()

	// done.
	return nil
}
