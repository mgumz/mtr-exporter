package main

// *mtr-exporter* periodically executes *mtr* to a given host and provides the
// measured results as prometheus metrics.

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/mgumz/mtr-exporter/pkg/job"
	"github.com/mgumz/mtr-exporter/pkg/timeshift"

	"github.com/robfig/cron/v3"
)

func main() {

	mtef := newFlags()
	flag.Usage = usage
	flag.Parse()

	if mtef.doPrintVersion {
		printVersion()
		return
	}
	if mtef.doPrintLicense {
		printLicense()
		return
	}
	if mtef.doPrintUsage {
		flag.Usage()
		return
	}

	// logging
	setupLogging(mtef.logLevel, mtef.doTimeStampLogs)

	// and go …
	collector := job.NewCollector()

	jobs := job.Jobs{}

	if len(flag.Args()) > 0 {

		tsmode := timeshift.None
		if mtef.timeShift != "" {
			tsmode = timeshift.RandomDelay
		}
		j := job.NewJob(mtef.mtrBin, flag.Args(), mtef.schedule, tsmode, mtef.timeShift)
		j.Label = mtef.jobLabel
		jobs = append(jobs, j)
	}

	jobsAvailable := !jobs.Empty()

	if mtef.jobFile != "" {
		if mtef.doWatchJobsFile != "" {
			slog.Info("watching -jobs-file",
				"jobs.fileName", mtef.jobFile,
				"jobs.schedule", mtef.doWatchJobsFile)
			job.WatchJobsFile(mtef.jobFile, mtef.mtrBin, mtef.doWatchJobsFile, collector)
			jobsAvailable = true
		} else {
			jobsFromFile, _, err := job.ParseJobFile(mtef.jobFile, mtef.mtrBin)
			if err != nil {
				slog.Error("parsing jobs file failed",
					"jobs.fileName", mtef.jobFile,
					"error", err)
				os.Exit(1)
			}
			if !jobsFromFile.Empty() {
				jobs = append(jobs, jobsFromFile...)
				jobsAvailable = true
			}
		}
	}

	if !jobsAvailable {
		slog.Error("no mtr jobs defined - provide at least one via -file or via arguments")
		os.Exit(1)
	}

	scheduler := cron.New(
		cron.WithLocation(time.UTC),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DiscardLogger),
		),
	)

	if err := jobs.ReSchedule(scheduler, collector); err != nil {
		slog.Error("", "error", err)
		os.Exit(1)
	}

	http.Handle("/metrics", collector)
	http.HandleFunc("/health", mtrHealthPage)
	http.HandleFunc("/", mtrIndexPage)

	slog.Info("serving...",
		"http.path", "/metrics",
		"http.bindAddr", mtef.bindAddr)

	server := &http.Server{
		Addr:              mtef.bindAddr,
		ReadHeaderTimeout: 1 * time.Second,
	}

	maybeError := slog.Attr{}
	if err := server.ListenAndServe(); err != nil {
		maybeError.Key = "error"
		maybeError.Value = slog.StringValue(err.Error())
	}
	slog.Info("done.", maybeError)
}
