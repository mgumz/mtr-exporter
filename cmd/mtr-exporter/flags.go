package main

import "flag"

type mteFlags struct {
	mtrBin   string
	jobLabel string
	bindAddr string
	jobFile  string
	schedule string

	doWatchJobsFile string
	doPrintVersion  bool
	doPrintUsage    bool
	doTimeStampLogs bool
}

func newFlags() *mteFlags {

	mte := new(mteFlags)

	flag.StringVar(&mte.mtrBin, "mtr", "mtr", "path to `mtr` binary")
	flag.StringVar(&mte.jobLabel, "label", "mtr-exporter-cli", "job label")
	flag.StringVar(&mte.bindAddr, "bind", ":8080", "bind address")
	flag.StringVar(&mte.jobFile, "jobs", "", "file containing job definitions")
	flag.StringVar(&mte.schedule, "schedule", "@every 60s", "schedule at which often `mtr` is launched")
	flag.StringVar(&mte.doWatchJobsFile, "watch-jobs", "", "re-parse -jobs file to schedule")
	flag.BoolVar(&mte.doPrintVersion, "version", false, "show version")
	flag.BoolVar(&mte.doPrintUsage, "h", false, "show help")
	flag.BoolVar(&mte.doTimeStampLogs, "tslogs", false, "use timestamps in logs")

	return mte
}
