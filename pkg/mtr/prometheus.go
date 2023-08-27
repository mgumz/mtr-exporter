package mtr

import (
	"fmt"
	"io"
)

func (hub *Hub) WriteMetrics(w io.Writer, labels string, ts int64) {
	// FIXME: remove deprecated metrics with mtr-exporter:0.4
	fmt.Fprintf(w, "mtr_report_snt{%s} %d %d\n", labels, hub.Snt, ts)
	fmt.Fprintln(w, "# deprecated metric name, use mtr_report_snt")
	fmt.Fprintf(w, "mtr_report_snt_gauge{%s} %d %d\n", labels, hub.Snt, ts)
	fmt.Fprintf(w, "mtr_report_loss{%s} %f %d\n", labels, hub.Loss, ts)
	fmt.Fprintln(w, "# deprecated metric name, use mtr_report_loss")
	fmt.Fprintf(w, "mtr_report_loss_gauge{%s} %f %d\n", labels, hub.Loss, ts)
	fmt.Fprintf(w, "mtr_report_best{%s} %f %d\n", labels, hub.Best, ts)
	fmt.Fprintln(w, "# deprecated metric name, use mtr_report_best")
	fmt.Fprintf(w, "mtr_report_best_gauge{%s} %f %d\n", labels, hub.Best, ts)
	fmt.Fprintf(w, "mtr_report_wrst{%s} %f %d\n", labels, hub.Wrst, ts)
	fmt.Fprintln(w, "# deprecated metric name, use mtr_report_wrst")
	fmt.Fprintf(w, "mtr_report_wrst_gauge{%s} %f %d\n", labels, hub.Wrst, ts)
	fmt.Fprintf(w, "mtr_report_avg{%s} %f %d\n", labels, hub.Avg, ts)
	fmt.Fprintln(w, "# deprecated metric name, use mtr_report_avg")
	fmt.Fprintf(w, "mtr_report_avg_gauge{%s} %f %d\n", labels, hub.Avg, ts)
	fmt.Fprintf(w, "mtr_report_last{%s} %f %d\n", labels, hub.Last, ts)
	fmt.Fprintln(w, "# deprecated metric name, use mtr_report_last")
	fmt.Fprintf(w, "mtr_report_last_gauge{%s} %f %d\n", labels, hub.Last, ts)
	fmt.Fprintf(w, "mtr_report_stdev{%s} %f %d\n", labels, hub.StDev, ts)
	fmt.Fprintln(w, "# deprecated metric name, use mtr_report_stdev")
	fmt.Fprintf(w, "mtr_report_stdev_gauge{%s} %f %d\n", labels, hub.StDev, ts)
}

func WriteMetricsHelpType(w io.Writer) {
	// FIXME: remove deprecated metrics with mtr-exporter:0.4
	fmt.Fprintln(w, "# HELP mtr_report_snt number of packets sent via mtr towards dst-host (see man mtr '-c')")
	fmt.Fprintln(w, "# TYPE mtr_report_snt gauge")
	fmt.Fprintln(w, "# HELP mtr_report_loss packet loss (percentage) of packets sent in last cycle")
	fmt.Fprintln(w, "# TYPE mtr_report_loss gauge")
	fmt.Fprintln(w, "# HELP mtr_report_best best rtt towards dst-host as observed in last cycle")
	fmt.Fprintln(w, "# TYPE mtr_report_best gauge")
	fmt.Fprintln(w, "# HELP mtr_report_wrst worst rtt towards dst-host as observed in last cycle")
	fmt.Fprintln(w, "# TYPE mtr_report_wrst gauge")
	fmt.Fprintln(w, "# HELP mtr_report_avg average rtt towards dst-host as observed in last cycle over all sent packets")
	fmt.Fprintln(w, "# TYPE mtr_report_avg gauge")
	fmt.Fprintln(w, "# HELP mtr_report_last last observed rtt of last cycle towards dst-host")
	fmt.Fprintln(w, "# TYPE mtr_report_last gauge")
	fmt.Fprintln(w, "# HELP mtr_report_stdev std-deviation of rtt of last cycle towards dst-host")
	fmt.Fprintln(w, "# TYPE mtr_report_stdev gauge")
}
