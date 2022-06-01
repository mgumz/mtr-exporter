package main

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

// ServeHTTP writes promtheues styled metrics about the last executed `mtr`
// run, see https://prometheus.io/docs/instrumenting/exposition_formats/#line-format
//
// NOTE: at the moment, no use of github.com/prometheus/client_golang/prometheus
// because overhead in size and complexity. once mtr-exporter requires features
// like push-gateway-export or graphite export or the like, we switch.
func (job *mtrJob) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if job.Report == nil {
		fmt.Fprintln(w, "# no current mtr runs performed (yet).")
		return
	}

	// the original job.Report might be changed in the
	// background by a successful run of mtr. copy (pointer to) the report
	// to have something safe to work on
	job.Lock()
	reports := job.Report
	ts := job.Launched.UTC()
	d := job.Duration
	job.Unlock()

	for _, report := range reports {
		labels := report.Mtr.Labels()
		tsMs := ts.UnixNano() / int64(time.Millisecond)

		fmt.Fprintf(w, "# mtr run: %s\n", ts.Format(time.RFC3339Nano))
		fmt.Fprintf(w, "# cmdline: %s\n", job.cmdLine)
		fmt.Fprintf(w, "mtr_report_duration_ms_gauge{%s} %d %d\n",
			labels2Prom(labels), d/time.Millisecond, tsMs)
		fmt.Fprintf(w, "mtr_report_count_hubs_gauge{%s} %d %d\n",
			labels2Prom(labels), len(report.Hubs), tsMs)

		for i, hub := range report.Hubs {
			labels["host"] = hub.Host
			labels["count"] = hub.Count
			// mark last hub to have it easily identified
			if i < (len(report.Hubs) - 1) {
				hub.writeMetrics(w, labels2Prom(labels), tsMs)
			} else {
				labels["last"] = "true"
				hub.writeMetrics(w, labels2Prom(labels), tsMs)
				delete(labels, "last")
			}
		}
	}
}

func (hub mtrHub) writeMetrics(w io.Writer, labels string, ts int64) {
	fmt.Fprintf(w, "mtr_report_loss_gauge{%s} %f %d\n", labels, hub.Loss, ts)
	fmt.Fprintf(w, "mtr_report_snt_gauge{%s} %d %d\n", labels, hub.Snt, ts)
	fmt.Fprintf(w, "mtr_report_last_gauge{%s} %f %d\n", labels, hub.Last, ts)
	fmt.Fprintf(w, "mtr_report_avg_gauge{%s} %f %d\n", labels, hub.Avg, ts)
	fmt.Fprintf(w, "mtr_report_best_gauge{%s} %f %d\n", labels, hub.Best, ts)
	fmt.Fprintf(w, "mtr_report_wrst_gauge{%s} %f %d\n", labels, hub.Wrst, ts)
	fmt.Fprintf(w, "mtr_report_stdev_gauge{%s} %f %d\n", labels, hub.StDev, ts)
}

func labels2Prom(labels map[string]string) string {
	sl := make(sort.StringSlice, 0, len(labels))
	for k, v := range labels {
		sl = append(sl, fmt.Sprintf("%s=%q", k, v))
	}
	sl.Sort()
	return strings.Join(sl, ",")
}
