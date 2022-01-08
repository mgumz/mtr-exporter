package main

import (
	"encoding/json"
	"io"
	"strconv"
)

type mtrResult struct {
	Report mtrReport `json:"report"`
}

type mtrReport struct {
	Mtr  mtrMtr   `json:"mtr"`
	Hubs []mtrHub `json:"hubs"`
}

type mtrMtr struct {
	Src        string `json:"src"`
	Dst        string `json:"dst"`
	Tos        int64  `json:"tos"`
	PSize      string `json:"psize"`
	BitPattern string `json:"bitpattern"`
	Tests      int64  `json:"tests"`
}

type mtrHub struct {
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

func (report *mtrReport) Decode(r io.Reader) error {
	dec := json.NewDecoder(r)
	result := mtrResult{}
	if err := dec.Decode(&result); err != nil {
		return err
	}
	*report = result.Report
	return nil
}

func (mtr *mtrMtr) Labels() map[string]string {
	return map[string]string{
		"src":        mtr.Src,
		"dst":        mtr.Dst,
		"tos":        strconv.FormatInt(mtr.Tos, 10),
		"psize":      mtr.PSize,
		"bitpattern": mtr.BitPattern,
		"tests":      strconv.FormatInt(mtr.Tests, 10),
	}
}
