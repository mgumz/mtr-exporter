## Changelog for mtr-exporter 0.6.0 (2025-06-05)

Features:

* Add `mtr_report_min_loss` metric: observed minimal packet loss for
  a given run, defaults to 100.

  Improves the situation for GH#33.

Maintenance:

* Bump base container image to Alpine:3.22
  (no update of mtr, still on 0.95)

## Changelog for mtr-exporter 0.5.1 (2025-02-26)

Bug Fixes:

* Fix improper placing of label "error" into the metrics. Fixes GH#32

## Changelog for mtr-exporter 0.5.0 (2025-02-13)

Features

* Add `mtr_runs_total` metric: successful runs and failed runs
  (which get label "error") are exposed. As a result, one can
  check intermediate (or permanent) failures.

  `sum by(mtr_exporter_job)(mtr_runs_total{})` - provides the
  absolute number of `mtr` runs for each job

  `sum by(mtr_exporter, error)(mtr_runs_total{error!=""})` - provides the
  number of failed runs

  The diff between these two is the amount of successful runs.

  This implements GH#29, GH#30

* Add small information page on "/", linking to /metrics and the project page.

Bug Fixes:

* Fix -watch-jobs - allow launching with zero given jobs, they might be
  added later (fixes GH#22)
* Fix picking up job file in a container / pod scenario (GH#20, thanks Clavin)
* Remove -flag.deprecatedMetrics from "-h" output (fixes GH#26)
* Fix logic bug for job(s) from command line
* Fix printing the version when -version is given.

Maintenance:

* Improve documentation (thanks Guillaume)

Contributors:

* Guillaume Berche - https://github.com/gberche-orange
* Clavianus Juneardo - https://github.com/clavinjune

## Changelog for mtr-exporter 0.4.0 (2024-11-25)

Features:

* add label 'hop' for the 'first', 'last' and
  'intermediate' hops towards the traced destination

Maintenance:

As announced with mtr-exporter 0.3.0, the following deprecated
metrics are removed:

| deprecated                   |
| ---------------------------- |
| mtr_report_duration_ms_gauge |
| mtr_report_count_hubs_gauge  |
| mtr_report_snt_gauge         |
| mtr_report_loss_gauge        |
| mtr_report_best_gauge        |
| mtr_report_wrst_gauge        |
| mtr_report_avg_gauge         |
| mtr_report_last_gauge        |
| mtr_report_stdev_gauge       |

Also, `-flag.deprecatedMetrics` is removed.


## Changelog for mtr-exporter 0.3.0 (2023-08-28)

Features:

* allow tracing multiple mtr destinations at once
  via mtr jobs file

* improve robustness on parsing mtr JSON output
  (<=mtr:0.93, >=mtr:0.94 differ)

* add -flag.deprecatedMetrics to render deprecated metrics.
  helps with transition time until deprecated metrics are
  gone

Improvements:

* implemented Prometheus Best Practices, as a result, some metrics
  are renamed, the old names are marked deprecated.

| deprecated                   | new                         |
| ---------------------------- | --------------------------- |
| mtr_report_duration_ms_gauge | mtr_report_duration_seconds |
| mtr_report_count_hubs_gauge  | mtr_report_count_hubs       |
| mtr_report_snt_gauge         | mtr_report_snt              |
| mtr_report_loss_gauge        | mtr_report_loss             |
| mtr_report_best_gauge        | mtr_report_best             |
| mtr_report_wrst_gauge        | mtr_report_wrst             |
| mtr_report_avg_gauge         | mtr_report_avg              |
| mtr_report_last_gauge        | mtr_report_last             |
| mtr_report_stdev_gauge       | mtr_report_stdev            |

Note: mtr_report_duration_seconds not only changed its name, but also
the unit: its **seconds** now, not **milliseconds** anymore.

With mtr-exporter:0.4.0 the deprecated ones will be removed.
Use -flag.deprecatedMetrics to have them exposed in the /metrics
response.

* internal: code adjustments, satisfying varios linters, code-quality
  reporters, vulnerability checkers etc
* internal: build system

## Changelog for mtr-exporter 0.2.0 (2022-07-15)

Features:

* Add /health endpoint for quick & cheap checking healthiness
  (Thanks Jakub Sokołowski)

## Changelog for mtr-exporter 0.1.0 (2022-01-08)

This is the initial version.
