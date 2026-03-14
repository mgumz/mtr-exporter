package job

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/shlex"

	"github.com/mgumz/mtr-exporter/pkg/timeshift"
)

// JobFile definition
//
// # comments, ignore everything after #
// ^space*$ - empty lines
// <label> -- <schedule> -- <mtr-flags>

func ParseJobs(r io.Reader, mtr string) (Jobs, error) {

	var err error
	var jobs = Jobs{}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	n := 0
	for scanner.Scan() {
		line := scanner.Text()
		n++
		job, err2 := parseJobLine(line, n, mtr)
		if err2 != nil {
			err = err2
			break
		}
		if job != nil {
			jobs = append(jobs, job)
		}
	}

	if err == nil {
		err = scanner.Err()
	}
	if err != nil {
		return Jobs{}, err
	}

	return jobs, nil
}

func parseJobLine(line string, lnr int, mtr string) (*Job, error) {

	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return nil, nil
	}

	if strings.HasPrefix(line, "#") {
		return nil, nil
	}

	const maxParts = 3

	parts := strings.SplitN(line, " -- ", maxParts)
	if len(parts) != maxParts {
		return nil, fmt.Errorf(errInvalidJobLine, lnr)
	}

	label, _ := parseLabel(strings.TrimSpace(parts[0]))
	schedule, tmode, tspec, err := parseSchedule(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, fmt.Errorf(errInvalidJobLine, lnr)
	}
	mtrArgs, _ := parseMtrArgs(strings.TrimSpace(parts[2]))

	job := NewJob(mtr, mtrArgs, schedule, tmode, tspec)
	job.Label = label

	return job, nil
}

func parseLabel(l string) (string, error) {
	return l, nil
}

// schedule s is either
// * "*/5 * * * * " - a regular cron expression
// * "@every 5m" - a constant delay expression
// optional, there is a random-delay at the end of the
// normal cron/delay expression:
// * "~30s"
func parseSchedule(s string) (string, timeshift.Mode, string, error) {

	tsMode := timeshift.None
	tsMarker := ""

	if i := strings.IndexAny(s, "~"); i > 0 {
		switch { //nolint:gocritic
		// NOTE: disabled for now, see timeshift.RandomDeviationScheduler for
		// the yet-to-be-solved issues
		// case strings.HasPrefix(s[i:], "±"):
		//	tsMode = timeshift.RandomDeviation
		//	tsMarker = "±"
		case strings.HasPrefix(s[i:], "~"):
			tsMode = timeshift.RandomDelay
			tsMarker = "~"
		}
	}

	if tsMode == timeshift.None {
		return strings.TrimSpace(s), tsMode, "", nil
	}

	// we don't check for "ok": we _know_ the marker is
	// in the string, the above code ensures that
	schedule, tsSpec, _ := strings.Cut(s, tsMarker)
	schedule = strings.TrimSpace(schedule)
	tsSpec = strings.TrimSpace(tsSpec)
	d, err := time.ParseDuration(tsSpec)
	if err != nil {
		return s, tsMode, "", fmt.Errorf(errTimeshiftFormat, tsSpec, err)
	}
	if d < 0 {
		return s, tsMode, "", fmt.Errorf(errTimeshiftNegative, tsSpec)
	}
	return schedule, tsMode, tsSpec, nil

}

func parseMtrArgs(s string) ([]string, error) {
	args, err := shlex.Split(s)
	return args, err
}
