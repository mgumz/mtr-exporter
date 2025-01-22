package main

import "fmt"

func usage() {

	const usage string = `Usage: mtr-exporter [FLAGS] -- [MTR-FLAGS]

FLAGS:
-bind       <bind-address>
            bind address (default ":8080")
-h
            show help
-jobs       <path-to-jobsfile>
            file describing multiple mtr-jobs. syntax is given below.
-label      <job-label>
            use <job-label> in prometheus-metrics (default: "mtr-exporter-cli")
-mtr        <path-to-binary>
            path to mtr binary (default: "mtr")
-schedule   <schedule>
            schedule at which often mtr is launched (default: "@every 60s")
            examples:
               @every <dur>  - example "@every 60s"
               @hourly       - run once per hour
               10 * * * *    - execute 10 minutes after the full hour
            see https://en.wikipedia.org/wiki/Cron
-tslogs
            use timestamps in logs
-watch-jobs <schedule>
            periodically watch the file defined via -jobs (default: "")
            if it has changed stop previously running mtr-jobs and apply
            all jobs defined in -jobs.
-version
            show version

MTR-FLAGS:
see "man mtr" for valid flags to mtr.

Examples:

$> mtr-exporter -- example.com
# probe every minute "example.com"

$> mtr-exporter -- -n example.com
# probe every minute "example.com", do not resolve DNS

$> mtr-exporter -schedule "@every 30s" -- -G 1 -m 3 -I ven3 -n example.com
# probe every 30s "example.com", wait 1s for response, try a max of 3 hops,
# use interface "ven3", do not resolve DNS.

Example Job File:

    # comments are ignored
    job1 -- @every 30s -- -I ven1 -n example.com
    job2 -- @every 30s -- -I ven2 -n example.com`

	fmt.Println(usage)
}
