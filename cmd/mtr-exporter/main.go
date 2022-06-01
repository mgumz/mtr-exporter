package main

// *mtr-exporter* periodically executes *mtr* to a given host and provides the
// measured results as prometheus metrics.

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/robfig/cron/v3"
)

func main() {

	log.SetFlags(0)

	mtrBin := flag.String("mtr", "mtr", "path to `mtr` binary")
	bind := flag.String("bind", ":8080", "bind address")
	schedule := flag.String("schedule", "@every 60s", "Schedule at which often `mtr` is launched")
	doPrintVersion := flag.Bool("version", false, "show version")
	doPrintUsage := flag.Bool("h", false, "show help")
	doTimeStampLogs := flag.Bool("tslogs", false, "use timestamps in logs")

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

	if len(flag.Args()) == 0 {
		log.Println("error: no mtr arguments given - at least the target host must be defined.")
		os.Exit(1)
		return
	}

	job := newMtrJob(*mtrBin, flag.Args())

	c := cron.New()
	c.AddFunc(*schedule, func() {
		log.Println("launching", job.cmdLine)
		if err := job.Launch(); err != nil {
			log.Println("failed:", err)
			return
		}
		for _, report := range job.Report {
			log.Println("done: ",
				len(report.Hubs), "hops in", job.Duration, ".")
		}

	})
	c.Start()

	http.Handle("/metrics", job)

	log.Println("serving /metrics at", *bind, "...")
	log.Fatal(http.ListenAndServe(*bind, nil))
}
