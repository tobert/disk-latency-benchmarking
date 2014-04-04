package fiotools

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Latency struct {
	Min        float64         `json:"min"`
	Max        float64         `json:"max"`
	Mean       float64         `json:"mean"`
	Stdev      float64         `json:"stdev"`
	Percentile map[float64]int `json:"percentile"`
}

type JobStats struct {
	IoBytes   int     `json:"io_bytes"`
	Bandwidth float64 `json:"bw"`
	BwMin     float64 `json:"bw_min"`
	BwMax     float64 `json:"bw_max"`
	BwAgg     float64 `json:"bw_agg"`
	BwMean    float64 `json:"bw_mean"`
	BwStdev   float64 `json:"bw_dev"`
	Iops      int     `json:"iops"`
	Runtime   int     `json:"runtime"`
	Slat      Latency `json:"slat"`
	Clat      Latency `json:"clat"`
	Lat       Latency `json:"lat"`
}

type ClientStat struct {
	Jobname           string             `json:"jobname"`
	Groupid           int                `json:"groupid"`
	Error             int                `json:"error"`
	Mixed             JobStats           `json:"mixed"`
	UsrCpu            float64            `json:"usr_cpu"`
	SysCpu            float64            `json:"sys_cpu"`
	ContextSwitches   int                `json:"ctx"`
	MajorFaults       int                `json:"majf"`
	MinorFaults       int                `json:"minf"`
	IODepthLevel      map[string]float64 `json:"iodepth_level"`
	LatencyUsec       map[string]float64 `json:"latency_us"`
	LatencyMsec       map[string]float64 `json:"latency_ms"`
	LatencyDepth      int                `json:"latency_depth"`
	LatencyTarget     int                `json:"latency_target"`
	LatencyPercentile float64            `json:"latency_percentile"`
	LatencyWindow     int                `json:"latency_window"`
	Hostname          string             `json:"hostname"`
	Port              int                `json:"port"`
}

type FioData struct {
	FioVersion    string       `json:"fio version"`
	HeaderGarbage string       `json:"garbage"`
	ClientStats   []ClientStat `json:"client_stats"`
}

func LoadFioJson(filename string) (fdata FioData) {
	dataBytes, err := ioutil.ReadFile(filename)

	if os.IsNotExist(err) {
		log.Fatal("Could not read file %s: %s", filename, err)
	}

	// fio writes a bunch of crap out to the output file before the JSON
	// so for now do the easy thing and find the first { after a \n
	// and call it good enough
	offset := bytes.Index(dataBytes, []byte("\n{"))

	err = json.Unmarshal(dataBytes[offset:], &fdata)
	if err != nil {
		log.Fatal("Could parse JSON: %s", err)
	}

	fdata.HeaderGarbage = string(dataBytes[0:offset])

	return fdata
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
