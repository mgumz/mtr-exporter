package mtr

import (
    "io"
    "fmt"
)

func (hub Hub) WriteMetrics(w io.Writer, labels string, ts int64) {
	fmt.Fprintf(w, "mtr_report_loss_gauge{%s} %f %d\n", labels, hub.Loss, ts)
	fmt.Fprintf(w, "mtr_report_snt_gauge{%s} %d %d\n", labels, hub.Snt, ts)
	fmt.Fprintf(w, "mtr_report_last_gauge{%s} %f %d\n", labels, hub.Last, ts)
	fmt.Fprintf(w, "mtr_report_avg_gauge{%s} %f %d\n", labels, hub.Avg, ts)
	fmt.Fprintf(w, "mtr_report_best_gauge{%s} %f %d\n", labels, hub.Best, ts)
	fmt.Fprintf(w, "mtr_report_wrst_gauge{%s} %f %d\n", labels, hub.Wrst, ts)
	fmt.Fprintf(w, "mtr_report_stdev_gauge{%s} %f %d\n", labels, hub.StDev, ts)
}

