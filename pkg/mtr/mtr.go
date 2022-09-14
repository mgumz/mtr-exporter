package mtr

import (
	"encoding/json"
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
	Mtr  Mtr   `json:"mtr"`
	Hubs []Hub `json:"hubs"`
}

type Mtr struct {
	Src        string `json:"src"`
	Dst        string `json:"dst"`
	Tos        int64  `json:"tos"`
	PSize      string `json:"psize"`
	BitPattern string `json:"bitpattern"`
	Tests      int64  `json:"tests"`
}

type Hub struct {
	Count int64   `json:"count"`
	Host  string  `json:"host"`
	Loss  float64 `json:"Loss%"`
	Snt   int64   `json:"Snt"`
	Last  float64 `json:"Last"`
	Avg   float64 `json:"Avg"`
	Best  float64 `json:"Best"`
	Wrst  float64 `json:"Wrst"`
	StDev float64 `json:"StDev"`
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

func (mtr *Mtr) Labels() map[string]string {
	return map[string]string{
		"src":        mtr.Src,
		"dst":        mtr.Dst,
		"tos":        strconv.FormatInt(mtr.Tos, integerBase),
		"psize":      mtr.PSize,
		"bitpattern": mtr.BitPattern,
		"tests":      strconv.FormatInt(mtr.Tests, integerBase),
	}
}
