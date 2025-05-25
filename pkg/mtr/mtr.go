package mtr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

const (
	integerBase int = 10
)

type Result struct {
	Report Report `json:"report"`
}

type Report struct {
	Mtr      Mtr    `json:"mtr"`
	Hubs     []Hub  `json:"hubs"`
	ErrorMsg string // carrying the error message of mtr
}

type Mtr struct {
	Src        string    `json:"src"`
	Dst        string    `json:"dst"`
	Tos        MtrNumber `json:"tos"`
	PSize      string    `json:"psize"`
	BitPattern string    `json:"bitpattern"`
	Tests      MtrNumber `json:"tests"`
}

// MtrNumber helps with JSON flavours used by <= mtr:0.93 and
// => mtr:0.94 which changed the type of the field
// see https://github.com/traviscross/mtr/pull/355
type MtrNumber int64

func (n *MtrNumber) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return fmt.Errorf("empty MtrNumber")
	}
	b = bytes.TrimLeft(b, `"`)
	b = bytes.TrimRight(b, `"`)
	i, err := strconv.ParseInt(string(b), 0, 64)
	if err != nil {
		return err
	}
	*n = MtrNumber(i)
	return nil
}

type Hub struct {
	Count MtrNumber `json:"count"`
	Host  string    `json:"host"`
	Loss  float64   `json:"Loss%"`
	Snt   int64     `json:"Snt"`
	Last  float64   `json:"Last"`
	Avg   float64   `json:"Avg"`
	Best  float64   `json:"Best"`
	Wrst  float64   `json:"Wrst"`
	StDev float64   `json:"StDev"`
}

func (report *Report) Decode(r io.Reader) error {
	dec := json.NewDecoder(r)
	result := Result{}
	if err := dec.Decode(&result); err != nil {
		return err
	}
	*report = result.Report
	return nil
}

func (report *Report) Empty() bool {
	return len(report.Hubs) == 0
}

func (report *Report) HubsTotal() int { return len(report.Hubs) }

func (mtr *Mtr) Labels() map[string]string {
	return map[string]string{
		"src":        mtr.Src,
		"dst":        mtr.Dst,
		"tos":        strconv.FormatInt(int64(mtr.Tos), integerBase),
		"psize":      mtr.PSize,
		"bitpattern": mtr.BitPattern,
		"tests":      strconv.FormatInt(int64(mtr.Tests), integerBase),
	}
}
