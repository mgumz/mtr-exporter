package job

import (
	"log/slog"
)

// cron.v3 interface
func (job *Job) Run() {

	slog.Info(infoJobLaunching,
		"job.label", job.Label,
		"job.cmd", job.cmdLine)

	if err := job.Launch(); err != nil {
		slog.Error(errJobLaunchFailed,
			"job.label", job.Label,
			"error", err)
		return
	}

	errMsg := slog.Attr{}
	if job.Report.ErrorMsg != "" {
		errMsg.Key = "error"
		errMsg.Value = slog.StringValue(job.Report.ErrorMsg)
	}

	slog.Info(infoJobDone,
		"job.label", job.Label,
		"nhops", len(job.Report.Hubs),
		"duration", job.Duration,
		errMsg)
}
