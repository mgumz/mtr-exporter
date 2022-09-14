package job

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/google/shlex"
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

	parts := strings.SplitN(line, " -- ", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid jobLine %d: expect '<label> -- <schedule> -- <mtr-flags>'", lnr)
	}

	label, _ := parseLabel(strings.TrimSpace(parts[0]))
	schedule, _ := parseSchedule(strings.TrimSpace(parts[1]))
	mtrArgs, _ := parseMtrArgs(strings.TrimSpace(parts[2]))

	job := NewJob(mtr, mtrArgs, schedule)
	job.Label = label

	return job, nil
}

func parseLabel(l string) (string, error) {
	return l, nil
}

func parseSchedule(s string) (string, error) {
	return s, nil
}

func parseMtrArgs(s string) ([]string, error) {
	args, err := shlex.Split(s)
	return args, err
}
