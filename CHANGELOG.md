
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
