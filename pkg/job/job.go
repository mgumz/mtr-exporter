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

	Runs map[string]int64
}

func (jm *JobMeta) DataAvailable() bool { return len(jm.Runs) > 0 }

type Job struct {
	JobMeta

	mtrBinary string
	args      []string
	cmdLine   string

	UpdateFn func(JobMeta) bool
}

func NewJob(mtr string, args []string, schedule string) *Job {
	extra := []string{
		"-j", // JSON output
	}
	args = append(extra, args...)
	job := Job{
		args:      args,
		mtrBinary: mtr,
		cmdLine:   strings.Join(append([]string{mtr}, args...), " "),
	}
	job.JobMeta.Runs = map[string]int64{}
	job.JobMeta.Schedule = schedule
	job.JobMeta.CmdLine = job.cmdLine
	return &job
}

func (job *Job) Launch() error {

	// TODO: maybe use CommandContext to have an upper limit in the execution

	cmd := exec.Command(job.mtrBinary, job.args...)

	// launch mtr
	bufStdout, bufStderr := bytes.Buffer{}, bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = &bufStdout, &bufStderr
	launched := time.Now()
	cmd.Run()
	duration := time.Since(launched)

	errMsg := normalizeMtrErrorMsg(bufStderr.String())
	if val, exists := job.Runs[errMsg]; exists {
		job.Runs[errMsg] = val + 1
	} else {
		job.Runs[errMsg] = 1
	}

	// decode the report
	report := mtr.Report{}
	if err := report.Decode(&bufStdout); err != nil {
		report.ErrorMsg = "error-decoding-mtr-json"
	} else {
		report.ErrorMsg = errMsg
	}

	// copy the report into the job
	job.JobMeta.Report = report
	job.JobMeta.Launched = launched
	job.JobMeta.Duration = duration

	if job.UpdateFn != nil {
		job.UpdateFn(job.JobMeta)
	}

	// done.
	return nil
}

func normalizeMtrErrorMsg(msg string) string {
	mf := func(r rune) rune {
		switch {
		case r == ' ' || r == ':':
			return r
		case r >= '0' && r <= '9':
			return r
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return 'a' + (r - 'A')
		}
		return '-'
	}
	return strings.Map(mf, strings.TrimSpace(msg))
}
