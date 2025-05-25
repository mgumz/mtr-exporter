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

	fmt.Fprintln(w, "# HELP mtr_runs_total number of mtr runs")
	fmt.Fprintln(w, "# TYPE mtr_runs_total counter")
	fmt.Fprintln(w, "# HELP mtr_report_duration_seconds duration of last mtr run (in seconds)")
	fmt.Fprintln(w, "# TYPE mtr_report_duration_seconds gauge")
	fmt.Fprintln(w, "# HELP mtr_report_count_hubs number of hops visited in the last mtr run")
	fmt.Fprintln(w, "# TYPE mtr_report_count_hubs gauge")
	fmt.Fprintln(w, "# HELP mtr_report_min_loss minimum packet loss (percentage float, 0 to 100) of all reported hubs")
	fmt.Fprintln(w, "# TYPE mtr_report_min_loss gauge")

	mtr.WriteMetricsHelpType(w)

	if len(c.jobs) == 0 {
		fmt.Fprintln(w, "# no mtr jobs defined (yet).")
		return
	}

	fmt.Fprintf(w, "# %d mtr jobs defined\n", len(c.jobs))

	for _, job := range c.jobs {

		if !job.DataAvailable() {
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

		errMsg := ""
		if report.ErrorMsg != "" {
			errMsg = fmt.Sprintf(" # (err: %q)", report.ErrorMsg)
		}
		fmt.Fprintf(w, "# mtr run %s: %s -- %s%s\n", job.Label, ts.Format(time.RFC3339Nano), job.CmdLine, errMsg)

		l := labels2Prom(labels)

		for k, v := range job.Runs {
			fmt.Fprintf(w, "mtr_runs_total{%s%s} %d %d\n",
				l, fmt.Sprintf(",error=%q", k), v, tsMs)
		}

		fmt.Fprintf(w, "mtr_report_duration_seconds{%s} %f %d\n",
			l, float64(d)/float64(time.Second), tsMs)

		fmt.Fprintf(w, "mtr_report_count_hubs{%s} %d %d\n",
			l, len(report.Hubs), tsMs)

		// in case the network does not provide any hubs between source and
		// destination, the number of report.Hubs is 0. this might happen
		// in VPN situations. to allow alert-systems to catch this situation,
		// mtr_report_min_loss is a metric
		minLoss := 100.0
		defer func() {
			fmt.Fprintf(w, "mtr_report_min_loss{%s} %f %d\n",
				l, minLoss, tsMs)
		}()

		if report.Empty() {
			continue
		}

		lh := report.HubsTotal() - 1
		for i, hub := range report.Hubs {

			minLoss = min(minLoss, hub.Loss)

			labels["host"] = hub.Host
			labels["count"] = strconv.FormatInt(int64(hub.Count), integerBase)
			labels["hop"] = hopLabel(i, lh)

			// "last" as label is redundant with `hop="last"`, but
			// also an existing label since mtr-exporter:0.1.0:
			// lets keep it for now.
			if i == lh {
				labels["last"] = "true"
			}

			hub.WriteMetrics(w, labels2Prom(labels), tsMs)

			delete(labels, "last")
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

func hopLabel(i, last int) string {

	if i == last {
		if i == 0 {
			return "first_last"
		}
		return "last"
	} else if i == 0 {
		return "first"
	}
	return "intermediate"
}
