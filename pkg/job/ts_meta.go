package job

import (
	"log/slog"

	"github.com/mgumz/mtr-exporter/pkg/timeshift"
)

type tsMeta struct {
	Mode timeshift.Mode
	Spec string
}

func (ts *tsMeta) LogValue() slog.Value {
	switch ts.Mode {
	case timeshift.RandomDelay:
		return slog.GroupValue(
			slog.String("mode", "delay"),
			slog.String("by", ts.Spec),
		)
	case timeshift.RandomDeviation:
		return slog.GroupValue(
			slog.String("mode", "deviation"),
			slog.String("by", ts.Spec),
		)
	}
	return slog.GroupValue(
		slog.String("mode", "none"),
	)
}
