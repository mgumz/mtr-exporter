package main

// *mtr-exporter* periodically executes *mtr* to a given host and provides the
// measured results as prometheus metrics.

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mgumz/mtr-exporter/pkg/job"

	"github.com/robfig/cron/v3"
)

func main() {
	log.SetFlags(0)

	mtrBin := flag.String("mtr", "mtr", "path to `mtr` binary")
	jobLabel := flag.String("label", "mtr-exporter-cli", "job label")
	bind := flag.String("bind", ":8080", "bind address")
	jobFile := flag.String("jobs", "", "file containing job definitions")
	schedule := flag.String("schedule", "@every 60s", "schedule at which often `mtr` is launched")
	doWatchJobsFile := flag.String("watch-jobs", "", "re-parse -jobs file to schedule")
	doPrintVersion := flag.Bool("version", false, "show version")
	doPrintUsage := flag.Bool("h", false, "show help")
	doTimeStampLogs := flag.Bool("tslogs", false, "use timestamps in logs")
	doRenderDeprecatedMetrics := flag.Bool("flag.deprecatedMetrics", false, "show deprecated metrics")

	flag.Usage = usage
	flag.Parse()

	if *doPrintVersion {
		printVersion()
		return
	}
	if *doPrintUsage {
		flag.Usage()
		return
	}
	if *doTimeStampLogs {
		log.SetFlags(log.LstdFlags | log.LUTC)
	}

	scheduler := cron.New()
	collector := job.NewCollector()
	collector.SetRenderDeprecatedMetrics(*doRenderDeprecatedMetrics)

	if len(flag.Args()) > 0 {
		j := job.NewJob(*mtrBin, flag.Args(), *schedule)
		j.Label = *jobLabel
		if _, err := scheduler.AddJob(j.Schedule, j); err != nil {
			log.Printf("error: unable to add %q to scheduler: %v", j.Label, err)
			os.Exit(1)
		}
		if !collector.AddJob(j.JobMeta) {
			log.Printf("error: unable to add %q to collector", j.Label)
			os.Exit(1)
		}
		j.UpdateFn = func(meta job.JobMeta) bool { return collector.UpdateJob(meta) }
	}

	if *jobFile != "" {
		if *doWatchJobsFile != "" {
			log.Printf("info: watching %q at %q", *jobFile, *doWatchJobsFile)
			job.WatchJobsFile(*jobFile, *mtrBin, *doWatchJobsFile, collector)
		} else {
			jobs, _, err := job.ParseJobFile(*jobFile, *mtrBin)
			if err != nil {
				log.Printf("error: parsing jobs file %q: %s", *jobFile, err)
				os.Exit(1)
			}
			if jobs.Empty() {
				log.Println("error: no mtr jobs defined - provide at least one via -file or via arguments")
				os.Exit(1)
			}
			for _, j := range jobs {
				if collector.AddJob(j.JobMeta) {
					if _, err := scheduler.AddJob(j.Schedule, j); err != nil {
						log.Printf("error: unable to add %q to collector: %v", j.Label, err)
						os.Exit(1)
					}
					j.UpdateFn = func(meta job.JobMeta) bool { return collector.UpdateJob(meta) }
				} // FIXME: log failed addition to collector, most likely
				// due to duplicate label
			}
		}
	}

	scheduler.Start()

	http.Handle("/metrics", collector)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	log.Println("serving /metrics at", *bind, "...")
	log.Fatal(http.ListenAndServe(*bind, nil))
}
