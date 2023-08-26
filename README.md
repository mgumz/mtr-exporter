# mtr-exporter

[![Go Report Card](https://goreportcard.com/badge/github.com/mgumz/mtr-exporter)](https://goreportcard.com/report/github.com/mgumz/mtr-exporter)

*mtr-exporter* periodically executes [mtr] to a given host and provides the
measured results as [prometheus] metrics.

Usually, [mtr] is producing the following output:

     HOST: src.example.com       Loss%   Snt   Last   Avg  Best  Wrst StDev
     1.|-- 127.0.0.1             0.0%     2    0.6   0.6   0.6   0.7   0.1
     2.|-- 127.0.0.2             0.0%     2    6.1  10.2   6.1  14.3   5.8
     3.|-- 127.0.0.3             0.0%     2   13.0  12.3  11.6  13.0   1.0
     4.|-- 127.0.0.4             0.0%     2    7.0   9.1   7.0  11.1   2.9
     5.|-- 127.0.0.5             0.0%     2   12.5  16.5  12.5  20.6   5.7
     6.|-- 127.0.0.6             0.0%     2   19.1  18.5  17.9  19.1   0.9
     7.|-- 127.0.0.7             0.0%     2   18.3  18.2  18.0  18.3   0.2
     8.|-- 127.0.0.8             0.0%     2   89.9  61.6  33.3  89.9  40.0
     9.|-- 127.0.0.9             0.0%     2   18.5  18.3  18.1  18.5   0.2
    10.|-- 127.0.0.10            0.0%     2   20.4  19.8  19.2  20.4   0.8

`mtr-exporter` exposes the measured values like this:

    # mtr run: 2020-03-08T16:37:05.000377Z
    # cmdline: /usr/local/sbin/mtr -j -c 2 -n example.com
    mtr_report_duration_ms_gauge{bitpattern="0x00",dst="example.com",psize="64",src="src.example.com",tests="2",tos="0x0"} 7179 1583685425000
    mtr_report_count_hubs_gauge{bitpattern="0x00",dst="example.com",psize="64",src="src.example.com",tests="2",tos="0x0"} 10 1583685425000
    mtr_report_loss_gauge{bitpattern="0x00",count="1",dst="example.com",host="127.0.0.1",psize="64",src="src.example.com",tests="2",tos="0x0"} 0.000000 1583685425000
    mtr_report_snt_gauge{bitpattern="0x00",count="1",dst="example.com",host="127.0.0.1",psize="64",src="src.example.com",tests="2",tos="0x0"} 2 1583685425000
    mtr_report_last_gauge{bitpattern="0x00",count="1",dst="example.com",host="127.0.0.1",psize="64",src="src.example.com",tests="2",tos="0x0"} 0.380000 1583685425000
    mtr_report_avg_gauge{bitpattern="0x00",count="1",dst="example.com",host="127.0.0.1",psize="64",src="src.example.com",tests="2",tos="0x0"} 0.480000 1583685425000
    mtr_report_best_gauge{bitpattern="0x00",count="1",dst="example.com",host="127.0.0.1",psize="64",src="src.example.com",tests="2",tos="0x0"} 0.380000 1583685425000
    mtr_report_wrst_gauge{bitpattern="0x00",count="1",dst="example.com",host="127.0.0.1",psize="64",src="src.example.com",tests="2",tos="0x0"} 0.570000 1583685425000
    mtr_report_stdev_gauge{bitpattern="0x00",count="1",dst="example.com",host="127.0.0.1",psize="64",src="src.example.com",tests="2",tos="0x0"} 0.130000 1583685425000

The last hop in the list of tested hosts contains the label `"last"="true"`.

When [prometheus] scrapes the data, you can visualise the observed values:

![MTR results in prometheus](./media/screenshot-2020-03-08+181019.9188279670.png "MTR 1")

![MTR results in prometheus](./media/screenshot-2020-03-08+181030.4810786850.png "MTR 1")

## Usage

    $> mtr-exporter [OPTS] -- [MTR-OPTS]

    mtr-exporter [FLAGS] -- [MTR-FLAGS]

    FLAGS:
    -bind <bind-address>
              bind address (default ":8080")
    -h        show help
    -mtr <path-to-binary>
              path to mtr binary (default "mtr")
    -schedule <schedule>
              schedule at which often mtr is launched (default "@every 60s")
              examples:
                @every <dur>  - example "@every 60s"
                @hourly       - run once per hour
                10 * * * *    - execute 10 minutes after the full hour
              see https://en.wikipedia.org/wiki/Cron
    -tslogs
              use timestamps in logs
    -version
              show version

    MTR-FLAGS:
            see "man mtr" for valid flags to mtr.

At `/metrics` the measured values of the last run are exposed.

Examples:

    $> mtr-exporter -- example.com
    # probe every minute "example.com"

    $> mtr-exporter -- -n example.com
    # probe every minute "example.com", do not resolve DNS

    $> mtr-exporter -schedule "@every 30s" -- -G 1 -m 3 -I ven3 -n example.com
    # probe every 30s "example.com", wait 1s for response, try a max of 3 hops,
    # use interface "ven3", do not resolve DNS.

## Requirements

Runtime:

* mtr-0.89 and newer (added --json support)

Build:

* golang-1.13 and newer

## Building

    $> git clone https://github.com/mgumz/mtr-exporter
    $> cd mtr-exporter
    $> make

One-off building and "installation":

    $> go get github.com/mgumz/mtr-exporter/cmd/mtr-exporter

## License

see LICENSE file

## Author(s)

* Mathias Gumz <mg@2hoch5.com>

[mtr]: https://www.bitwizard.nl/mtr/index.html
[prometheus]: https://prometheus.io
