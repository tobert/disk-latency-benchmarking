package fiotools

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Latency struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Mean  float64 `json:"mean"`
	Stdev float64 `json:"stdev"`
	// JSON keys are strings, let the parser do its thing
	// then copy to Percentile after
	PercentileRaw map[string]float64 `json:"percentile"`
	Percentile    map[float64]float64
}

type JobStats struct {
	IoBytes   int      `json:"io_bytes"`
	Bandwidth float64  `json:"bw"`
	BwMin     float64  `json:"bw_min"`
	BwMax     float64  `json:"bw_max"`
	BwAgg     float64  `json:"bw_agg"`
	BwMean    float64  `json:"bw_mean"`
	BwStdev   float64  `json:"bw_dev"`
	Iops      int      `json:"iops"`
	Runtime   int      `json:"runtime"`
	Slat      *Latency `json:"slat"`
	Clat      *Latency `json:"clat"`
	Lat       *Latency `json:"lat"`
}

type ClientStat struct {
	Jobname         string             `json:"jobname"`
	Groupid         int                `json:"groupid"`
	Error           int                `json:"error"`
	Mixed           *JobStats          `json:"mixed"` // fio config dependent
	UsrCpu          float64            `json:"usr_cpu"`
	SysCpu          float64            `json:"sys_cpu"`
	ContextSwitches int                `json:"ctx"`
	MajorFaults     int                `json:"majf"`
	MinorFaults     int                `json:"minf"`
	IODepthLevelRaw map[string]float64 `json:"iodepth_level"`
	LatencyUsecRaw  map[string]float64 `json:"latency_us"`
	LatencyMsecRaw  map[string]float64 `json:"latency_ms"`
	// same deal as PercentileStr above
	IODepthLevel      map[float64]float64
	LatencyUsec       map[float64]float64
	LatencyMsec       map[float64]float64
	LatencyDepth      int     `json:"latency_depth"`
	LatencyTarget     int     `json:"latency_target"`
	LatencyPercentile float64 `json:"latency_percentile"`
	LatencyWindow     int     `json:"latency_window"`
	Hostname          string  `json:"hostname"`
	Port              int     `json:"port"`
}

type FioData struct {
	FioVersion    string        `json:"fio version"`
	HeaderGarbage string        `json:"garbage"`
	ClientStats   []ClientStat  `json:"client_stats"`
	DiskUtil      []interface{} `json:"disk_util"` // unused for now
}

func LoadFioJson(filename string) (fio_data FioData) {
	dataBytes, err := ioutil.ReadFile(filename)

	if os.IsNotExist(err) {
		log.Fatal("Could not read file %s: %s", filename, err)
	}

	// fio writes a bunch of crap out to the output file before the JSON
	// so for now do the easy thing and find the first { after a \n
	// and call it good enough
	offset := bytes.Index(dataBytes, []byte("\n{"))

	err = json.Unmarshal(dataBytes[offset:], &fio_data)
	if err != nil {
		log.Fatal("Could parse JSON: %s", err)
	}

	fio_data.HeaderGarbage = string(dataBytes[0:offset])

	// now go over the maps of string => float64 and fix them up to be float64 => float64
	for _, cs := range fio_data.ClientStats {
		cs.IODepthLevel = cleanKeys(cs.IODepthLevelRaw)
		cs.LatencyUsec = cleanKeys(cs.LatencyUsecRaw)
		cs.LatencyMsec = cleanKeys(cs.LatencyMsecRaw)

		// TODO: add similar checks for Read/Write once I know the final names
		cleanPercentiles(cs.Mixed)
	}

	return fio_data
}

func cleanPercentiles(in *JobStats) {
	in.Lat.Percentile = cleanKeys(in.Lat.PercentileRaw)
	in.Clat.Percentile = cleanKeys(in.Clat.PercentileRaw)
	in.Slat.Percentile = cleanKeys(in.Slat.PercentileRaw)
}

func cleanKeys(in map[string]float64) map[float64]float64 {
	if in == nil {
		return nil
	}

	out := make(map[float64]float64, len(in))

	for k, v := range in {
		// remove the >= fio puts in some of the keys
		cleaned := strings.TrimPrefix(k, ">=")
		fkey, _ := strconv.ParseFloat(cleaned, 64)
		out[fkey] = v
	}

	return out
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
