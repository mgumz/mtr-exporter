package job

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/blake2b"

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

		// in case the network does not provide any hubs between source and
		// destination, the number of report.Hubs is 0. this might happen
		// in VPN situations. to allow alert-systems to catch this situation,
		// mtr_report_min_loss is a metric
		minLoss := 100.0

		if !report.Empty() {
			writeMetricsForHubs(w, report, tsMs, labels, &minLoss)
		}

		l := labels2Prom(labels)

		for k, v := range job.Runs {
			fmt.Fprintf(w, "mtr_runs_total{%s%s} %d %d\n",
				l, fmt.Sprintf(",error=%q", k), v, tsMs)
		}

		fmt.Fprintf(w, "mtr_report_duration_seconds{%s} %f %d\n",
			l, float64(d)/float64(time.Second), tsMs)

		fmt.Fprintf(w, "mtr_report_count_hubs{%s} %d %d\n",
			l, len(report.Hubs), tsMs)

		fmt.Fprintf(w, "mtr_report_min_loss{%s} %f %d\n",
			l, minLoss, tsMs)
	}
}

func writeMetricsForHubs(w io.Writer, report mtr.Report, tsMs int64, labels map[string]string, minLoss *float64) {

	path := []string{}

	lh := report.HubsTotal() - 1
	for i, hub := range report.Hubs {

		*minLoss = min(*minLoss, hub.Loss)
		path = append(path, hub.Host)

		labels["host"] = hub.Host
		labels["count"] = strconv.FormatInt(int64(hub.Count), integerBase)
		labels["hop"] = hopLabel(i, lh)

		// "last" as label is redundant with `hop="last"`, but
		// also an existing label ("api contract") since mtr-exporter:0.1.0:
		// lets keep it for now.
		if i == lh {
			labels["last"] = "true"
			labels["path_id"] = strconv.FormatInt(pathId(path), integerBase)
		}

		labelsStr := labels2Prom(labels)

		hub.WriteMetrics(w, labelsStr, tsMs)
		if i == lh {
			fmt.Fprintf(w, "mtr_report_path_id{%s} %s %d\n", labelsStr, labels["path_id"], tsMs)
		}

		// map "labels" is modified and propagated back to the caller. this
		// will lead to a "last=true" label on metric mtr_report_min_loss
		// which is not desired.
		delete(labels, "last")
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

	switch {
	case (i == 0) && (last == 0):
		return "first_last"
	case i == last:
		return "last"
	case i == 0:
		return "first"
	default:
		return "intermediate"
	}
}

// calculates a "pathId" of the list of hosts.
// when the path to the destination changes, the pathId
// should change.
func pathId(hosts []string) int64 {

	path := strings.Join(hosts, " ")

	// motivation for blake2b:
	//
	// situation:
	// - `mtr` yields    `"host": "host.example.com"`
	// - `mtr -n` yields `"host": "192.0.2.1"`
	// - `mtr -b` yields `"host": "host.example.com (192.0.2.1)"`
	// - an unidentifyable host is represented as "???"
	// - host IPs can be IPv4 and IPv6
	//
	// so, although chksum(hosts) is not directly understandable
	// by the human eye, it will work for all of the above cases,
	// it will change when the path changes, it will change when
	// the reverse DNS lookup changes and it is reasonable "long" / "short".
	// blake2b is picked over sha256 because it can produce shorter
	// checksums (the chance of an "attack" is rather slim: to cause
	// a collision to "hide" a specific hop/host in a observed
	// path would require to craft a collision causing reverse DNS
	// entry for the host in question.

	// 8 byte (64 bit) digest. RFC7693 states 2^80 collision security for a 20
	// byte (160 bit) digest, 2^256 collision for a 64 byte ( 512 bit). so, 8
	// byte (64 bit) should be good for 2^32.

	hasher, _ := blake2b.New(8, nil)
	hasher.Write([]byte(path))

	pathId := int64(0)
	for i, n := range hasher.Sum(nil) {
		pathId = pathId | (int64(n) << (i * 8))
	}

	return pathId
}
