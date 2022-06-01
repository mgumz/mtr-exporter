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

func newMtrJob(mtr string, args []string) *mtrJob {
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
	domains := []string{
		"us-east-bidder.mathtag.com",
		"33across-us-east.lb.indexww.com",
	}
	args1 := job.args
	args1 = append(args1, domains[0])

	args2 := job.args
	args2 = append(args2, domains[1])

	cmd1 := exec.Command(job.mtrBinary, args1...) // Будет работать если не передать домен через пробел
	cmd2 := exec.Command(job.mtrBinary, args2...)

	// launch mtr
	buf := bytes.Buffer{}
	cmd1.Stdout = &buf
	launched := time.Now()
	if err := cmd2.Run(); err != nil {
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
