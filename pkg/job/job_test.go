package job

import (
	"testing"
)

func Test_NormalizeMtrErrorMsg(t *testing.T) {

	fixtures := []struct {
		msg        string
		normalized string
	}{
		{"mtr: Unexpected mtr-packet error", "mtr: unexpected mtr-packet error"},
	}

	for i := range fixtures {
		normalized := normalizeMtrErrorMsg(fixtures[i].msg)
		if normalized != fixtures[i].normalized {
			t.Fatalf("expected %q for %q, got %q", fixtures[i].normalized, fixtures[i].msg, normalized)
		}
	}
}
