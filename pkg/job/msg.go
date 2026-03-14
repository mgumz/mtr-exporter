package job

const (
	errDecodeError       = "error-decoding-mtr-json"
	errAddToCollector    = "unable to add job to collector"
	errGenericCollector  = "collector error"
	errAddToScheduler    = "unable to add job to scheduler"
	errGenericSchedule   = "schedule error"
	errInvalidSchedule   = "invalid schedule in line %d"
	errInvalidJobLine    = "invalid jobLine %d: expect '<label> -- <schedule> -- <mtr-flags>'"
	errTimeshiftFormat   = "timeshift deviation %q, format error %s"
	errTimeshiftNegative = "timeshift deviation must not be negative %q"
	errLaunchJobWatch    = "unable to launch watch-jobs scheduler"
	errJobLaunchFailed   = "launch failed"
)

const (
	infoStartingParse    = "starting to parse"
	infoDoneParse        = "done parsing"
	infoJobFileUnchanged = "watched file is unchanged"
	infoJobFileChanged   = "watched file changed: replacing jobs"
	infoJobLaunching     = "launching job"
	infoJobDone          = "job done"
	infoJobSchedule      = "schedule job"
)
