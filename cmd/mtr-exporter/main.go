package main

// *mtr-exporter* periodically executes *mtr* to a given host and provides the
// measured results as prometheus metrics.

import (
	"flag"
	"fmt"
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

	_, err := c.AddFunc(*schedule, func() {
		log.Println("launching", job.cmdLine)
		if err := job.Launch(); err != nil {
			log.Println("failed:", err)
			return
		}
		log.Println("done: ",
			len(job.Report.Hubs), "hops in", job.Duration, ".")
	})
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}

	c.Start()

	http.Handle("/metrics", job)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	log.Println("serving /metrics at", *bind, "...")
	log.Fatal(http.ListenAndServe(*bind, nil))
}
