package main

import "fmt"

func usage() {
	usage := `Usage: mtr-exporter [FLAGS] -- [MTR-FLAGS]

FLAGS:
-bind <bind-address>
	  bind address (default ":8080")
-h	  show help
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
	
Examples:

$> mtr-exporter -- example.com
# probe every minute "example.com"

$> mtr-exporter -- -n example.com
# probe every minute "example.com", do not resolve DNS

$> mtr-exporter -schedule "@every 30s" -- -G 1 -m 3 -I ven3 -n example.com
# probe every 30s "example.com", wait 1s for response, try a max of 3 hops,
# use interface "ven3", do not resolve DNS.`

	fmt.Println(usage)
}
