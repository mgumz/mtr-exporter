package main

// *mtr-exporter* periodically executes *mtr* to a given host and provides the
// measured results as prometheus metrics.

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mgumz/mtr-exporter/pkg/job"

	"github.com/robfig/cron/v3"
)

func main() {
	log.SetFlags(0)

	mtef := newFlags()
	flag.Usage = usage
	flag.Parse()

	if mtef.doPrintVersion {
		printVersion()
		return
	}
	if mtef.doPrintUsage {
		flag.Usage()
		return
	}
	if mtef.doTimeStampLogs {
		log.SetFlags(log.LstdFlags | log.LUTC)
	}

	collector := job.NewCollector()

	jobs := job.Jobs{}

	if len(flag.Args()) > 0 {
		j := job.NewJob(mtef.mtrBin, flag.Args(), mtef.schedule)
		j.Label = mtef.jobLabel
		jobs = append(jobs, j)
	}

	jobsAvailable := !jobs.Empty()

	if mtef.jobFile != "" {
		if mtef.doWatchJobsFile != "" {
			log.Printf("info: watching %q at %q", mtef.jobFile, mtef.doWatchJobsFile)
			job.WatchJobsFile(mtef.jobFile, mtef.mtrBin, mtef.doWatchJobsFile, collector)
			jobsAvailable = true
		} else {
			jobsFromFile, _, err := job.ParseJobFile(mtef.jobFile, mtef.mtrBin)
			if err != nil {
				log.Printf("error: parsing jobs file %q: %s", mtef.jobFile, err)
				os.Exit(1)
			}
			if !jobsFromFile.Empty() {
				jobs = append(jobs, jobsFromFile...)
				jobsAvailable = true
			}
		}
	}

	if !jobsAvailable {
		log.Println("error: no mtr jobs defined - provide at least one via -file or via arguments")
		os.Exit(1)
	}

	scheduler := cron.New(
		cron.WithLocation(time.UTC),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DiscardLogger),
		),
	)

	if err := jobs.ReSchedule(scheduler, collector); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}

	http.Handle("/metrics", collector)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `<!doctype html>
	<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>mtr-exporter</title>
	</head>
	<body>
		mtr-exporter - <a href="https://github.com/mgumz/mtr-exporter">https://github.com/mgumz/mtr-exporter<a><br>
		see <a href="/metrics">/metrics</a>.
	</body>`)
	})

	log.Println("serving /metrics at", mtef.bindAddr, "...")
	log.Fatal(http.ListenAndServe(mtef.bindAddr, nil))
}
