package mtr

import (
	"encoding/json"
	"strings"
	"testing"
)

func Test_MtrReportDecoding(t *testing.T) {
	body := `
	{
		"report": {
		  "mtr": {
			"src": "src.example.com",
			"dst": "dst.example.com",
			"tos": 0,
			"tests": 2,
			"psize": "64",
			"bitpattern": "0x00"
		  },
		  "hubs": [{
			"count": 1,
			"host": "127.0.0.1",
			"Loss%": 0.00,
			"Snt": 2,
			"Last": 0.54,
			"Avg": 0.54,
			"Best": 0.54,
			"Wrst": 0.54,
			"StDev": 0.00
		  },
		  {
			"count": 2,
			"host": "127.0.0.2",
			"Loss%": 50.00,
			"Snt": 2,
			"Last": 5.26,
			"Avg": 5.26,
			"Best": 5.26,
			"Wrst": 5.26,
			"StDev": 0.00
		  },
		  {
			"count": 3,
			"host": "127.0.0.3",
			"Loss%": 0.00,
			"Snt": 2,
			"Last": 5.97,
			"Avg": 6.62,
			"Best": 5.97,
			"Wrst": 7.26,
			"StDev": 0.91
		  },
		  {
			"count": 4,
			"host": "127.0.0.4",
			"Loss%": 0.00,
			"Snt": 2,
			"Last": 7.13,
			"Avg": 6.69,
			"Best": 6.24,
			"Wrst": 7.13,
			"StDev": 0.63
		  },
		  {
			"count": 5,
			"host": "127.0.0.5",
			"Loss%": 0.00,
			"Snt": 2,
			"Last": 13.12,
			"Avg": 15.43,
			"Best": 13.12,
			"Wrst": 17.73,
			"StDev": 3.26
		  },
		  {
			"count": 6,
			"host": "127.0.0.6",
			"Loss%": 0.00,
			"Snt": 2,
			"Last": 20.08,
			"Avg": 18.65,
			"Best": 17.22,
			"Wrst": 20.08,
			"StDev": 2.02
		  },
		  {
			"count": 7,
			"host": "127.0.0.7",
			"Loss%": 0.00,
			"Snt": 2,
			"Last": 19.00,
			"Avg": 18.15,
			"Best": 17.30,
			"Wrst": 19.00,
			"StDev": 1.20
		  },
		  {
			"count": 8,
			"host": "127.0.0.8",
			"Loss%": 0.00,
			"Snt": 2,
			"Last": 18.02,
			"Avg": 18.32,
			"Best": 18.02,
			"Wrst": 18.61,
			"StDev": 0.42
		  },
		  {
			"count": 9,
			"host": "127.0.0.9",
			"Loss%": 0.00,
			"Snt": 2,
			"Last": 19.15,
			"Avg": 19.05,
			"Best": 18.96,
			"Wrst": 19.15,
			"StDev": 0.13
		  },
		  {
			"count": 10,
			"host": "127.0.0.10",
			"Loss%": 0.00,
			"Snt": 2,
			"Last": 22.54,
			"Avg": 20.74,
			"Best": 18.94,
			"Wrst": 22.54,
			"StDev": 2.54
		  }]
		}
	}`

	report := &Report{}
	if err := report.Decode(strings.NewReader(body)); err != nil {
		t.Fatalf("error decoding: %s\n%s", err, body)
	}

	if report.Mtr.Dst != "dst.example.com" {
		t.Fatalf("error parsing mtr report: expected %q, got %q\n%v",
			"dst.example.com",
			report.Mtr.Dst,
			report)
	}
}

func Test_MtrEmptyHubs(t *testing.T) {
	body := `
	{
		"report": {
			"mtr": {
				"src": "example-src.test",
				"dst": "example-dst.invalid",
				"tos": 0,
				"tests": 10,
				"psize": "64",
				"bitpattern": "0x00"
			},
			"hubs": []
		}
	}`

	report := &Report{}
	if err := report.Decode(strings.NewReader(body)); err != nil {
		t.Fatalf("error decoding: %s\n%s", err, body)
	}

	if len(report.Hubs) != 0 {
		t.Fatalf("error: expected [] hubs, got %d", len(report.Hubs))
	}
}

func Test_MtrJSONDecoding(t *testing.T) {

	fixtures := []string{
		// <= mtr:0.93 format
		`{
			"src": "src.example.com",
			"dst": "dst.example.com",
			"tos": "0x0",
			"tests": 2,
			"psize": "64",
			"bitpattern": "0x00"
		}`,
		// >= mtr:0.94 format
		`{
			"src": "src.example.com",
			"dst": "dst.example.com",
			"tos": 64,
			"tests": 2,
			"psize": "64",
			"bitpattern": "0x00"
		}`,
	}

	for _, f := range fixtures {
		r := strings.NewReader(f)
		d := json.NewDecoder(r)
		m := Mtr{}
		err := d.Decode(&m)
		if err != nil {
			t.Fatalf("%s\n%s", f, err)
		}
	}
}
