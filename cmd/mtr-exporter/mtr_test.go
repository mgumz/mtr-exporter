package main

import (
	"strings"
	"testing"
)

func Test_MtrReportDecoding(t *testing.T) {

	body := `
	{
		"report": {
		  "mtr": {
			"src": "ssp-server-33x-prod-us-east4-00",
			"dst": "ya.ru",
			"tos": "0x0",
			"psize": "64",
			"bitpattern": "0x00",
			"tests": "10"
		  },
		  "hubs": [{
			"count": "1",
			"host": "190.98.141.22",
			"Loss%": 0.00,
			"Snt": 10,
			"Last": 2.58,
			"Avg": 6.20,
			"Best": 1.92,
			"Wrst": 23.79,
			"StDev": 7.16
		  },
		  {
			"count": "2",
			"host": "telecom-italia-et-15-0-17-1-0-grtwaseq5.net.telefonicaglobalsolutions.com",
			"Loss%": 0.00,
			"Snt": 10,
			"Last": 1.87,
			"Avg": 2.47,
			"Best": 1.77,
			"Wrst": 6.87,
			"StDev": 1.57
		  },
		  {
			"count": "3",
			"host": "ae9.stoccolma1.sto.seabone.net",
			"Loss%": 20.00,
			"Snt": 10,
			"Last": 270.76,
			"Avg": 273.15,
			"Best": 264.23,
			"Wrst": 291.61,
			"StDev": 8.19
		  }]
		}
	  }
	  [root@ssp-server-33x-prod-us-east4-00 ~]# mtr -j -G 1 -m 3 ya.ru
	  {
		"report": {
		  "mtr": {
			"src": "ssp-server-33x-prod-us-east4-00",
			"dst": "ya.ru",
			"tos": "0x0",
			"psize": "64",
			"bitpattern": "0x00",
			"tests": "10"
		  },
		  "hubs": [{
			"count": "1",
			"host": "190.98.141.22",
			"Loss%": 0.00,
			"Snt": 10,
			"Last": 2.17,
			"Avg": 2.47,
			"Best": 1.68,
			"Wrst": 3.81,
			"StDev": 0.66
		  },
		  {
			"count": "2",
			"host": "telecom-italia-et-15-0-17-1-0-grtwaseq5.net.telefonicaglobalsolutions.com",
			"Loss%": 0.00,
			"Snt": 10,
			"Last": 2.21,
			"Avg": 2.40,
			"Best": 2.04,
			"Wrst": 3.52,
			"StDev": 0.45
		  },
		  {
			"count": "3",
			"host": "ae9.stoccolma1.sto.seabone.net",
			"Loss%": 20.00,
			"Snt": 10,
			"Last": 270.54,
			"Avg": 267.74,
			"Best": 263.99,
			"Wrst": 271.40,
			"StDev": 2.78
		  }]
		}
	  }`

	report := &mtrReport{}
	if err := report.Decode(strings.NewReader(body)); err != nil {
		t.Fatalf("error decoding: %s\n%s", err, body)
	}

	if report.Mtr.Dst != "ya.ru" {
		t.Fatalf("error parsing mtr report: expected %q, got %q\n%v",
			"ya.ru",
			report.Mtr.Dst,
			report)
	}
}
