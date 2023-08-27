package job

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mgumz/mtr-exporter/pkg/mtr"
)

const (
	integerBase int = 10
)

// ServeHTTP writes promtheues styled metrics about the last executed `mtr`
// run, see https://prometheus.io/docs/instrumenting/exposition_formats/#line-format
//
// NOTE: at the moment, no use of github.com/prometheus/client_golang/prometheus
// because overhead in size and complexity. once mtr-exporter requires features
// like push-gateway-export or graphite export or the like, we switch.
func (c *Collector) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Fprintln(w, "# HELP mtr_report_duration_seconds duration of last mtr run (in seconds)")
	fmt.Fprintln(w, "# TYPE mtr_report_duration_seconds gauge")
	fmt.Fprintln(w, "# HELP mtr_report_count_hubs number of hops visited in the last mtr run")
	fmt.Fprintln(w, "# TYPE mtr_report_count_hubs gauge")

	mtr.WriteMetricsHelpType(w)

	if len(c.jobs) == 0 {
		fmt.Fprintln(w, "# no mtr jobs defined (yet).")
		return
	}

	fmt.Fprintf(w, "# %d mtr jobs defined\n", len(c.jobs))

	for _, job := range c.jobs {

		if len(job.Report.Hubs) == 0 {
			continue
		}

		// the original job.Report might be changed in the
		// background by a successful run of mtr. copy (pointer to) the report
		// to have something safe to work on
		report := job.Report
		ts := job.Launched.UTC()
		d := job.Duration

		labels := report.Mtr.Labels()
		labels["mtr_exporter_job"] = job.Label
		tsMs := ts.UnixNano() / int64(time.Millisecond)

		fmt.Fprintf(w, "# mtr run %s: %s -- %s\n", job.Label, ts.Format(time.RFC3339Nano), job.CmdLine)

		l := labels2Prom(labels)

		// FIXME: remove deprecated metrics with mtr-exporter:0.4
		fmt.Fprintf(w, "mtr_report_duration_seconds{%s} %f %d\n",
			l, float64(d)/float64(time.Second), tsMs)
		fmt.Fprintln(w, "# deprecated metric name, use mtr_report_duration_seconds")
		fmt.Fprintf(w, "mtr_report_duration_ms_gauge{%s} %d %d\n",
			l, d/time.Millisecond, tsMs)
		fmt.Fprintf(w, "mtr_report_count_hubs{%s} %d %d\n",
			l, len(report.Hubs), tsMs)
		fmt.Fprintln(w, "# deprecated metric, use mtr_report_count_hubs")
		fmt.Fprintf(w, "mtr_report_count_hubs_gauge{%s} %d %d\n",
			l, len(report.Hubs), tsMs)

		for i, hub := range report.Hubs {
			labels["host"] = hub.Host
			labels["count"] = strconv.FormatInt(int64(hub.Count), integerBase)
			// mark last hub to have it easily identified
			if i < (len(report.Hubs) - 1) {
				hub.WriteMetrics(w, labels2Prom(labels), tsMs)
			} else {
				labels["last"] = "true"
				hub.WriteMetrics(w, labels2Prom(labels), tsMs)
				delete(labels, "last")
			}
		}
	}

}

func labels2Prom(labels map[string]string) string {
	sl := make(sort.StringSlice, 0, len(labels))
	for k, v := range labels {
		sl = append(sl, fmt.Sprintf("%s=%q", k, v))
	}
	sl.Sort()
	return strings.Join(sl, ",")
}
