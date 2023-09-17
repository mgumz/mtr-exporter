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

	scheduler := cron.New(
		cron.WithLocation(time.UTC),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DiscardLogger),
		),
	)
	collector := job.NewCollector()
	collector.SetRenderDeprecatedMetrics(mtef.doRenderDeprecatedMetrics)

	if len(flag.Args()) > 0 {
		j := job.NewJob(mtef.mtrBin, flag.Args(), mtef.schedule)
		j.Label = mtef.jobLabel
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

	if mtef.jobFile != "" {
		if mtef.doWatchJobsFile != "" {
			log.Printf("info: watching %q at %q", mtef.jobFile, mtef.doWatchJobsFile)
			job.WatchJobsFile(mtef.jobFile, mtef.mtrBin, mtef.doWatchJobsFile, collector)
		} else {
			jobs, _, err := job.ParseJobFile(mtef.jobFile, mtef.mtrBin)
			if err != nil {
				log.Printf("error: parsing jobs file %q: %s", mtef.jobFile, err)
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

	log.Println("serving /metrics at", mtef.bindAddr, "...")
	log.Fatal(http.ListenAndServe(mtef.bindAddr, nil))
}
