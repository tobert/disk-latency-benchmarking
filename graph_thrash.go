package main

/* Draw a graph based on 'thrash.json' which is the output from fio run with
 * the fio-configs/thrash.fio config and the following command line:
 * sudo fio --output-format=json thrash.fio --output=thrash.json
 *
 * Run this file with: go run graph_thrash.go
 * It will create the file 'test.png' in the current working directory.
 */

import (
	"./src/fiotools"
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/plotutil"
	"code.google.com/p/plotinum/vg"
	"fmt"
	"math"
	"sort"
)

// TODO: flags

// uses the bargraph render instead of histogram because fio only
// provides the bucket averages
func main() {
	// TODO: throw this in a JSON file
	legend := map[string]string{
		"/dev/disk/by-path/pci-0000:03:00.0-sas-0x5000c5000d7f96d9-lun-0": "SAS",
		"/dev/disk/by-id/wwn-0x5000c500151229dd":                          "SATA",
		"/dev/disk/by-id/ata-Samsung_SSD_840_PRO_Series_S1ANNSADB05219A":  "SSD",
		"/dev/disk/by-id/md-uuid-6bb71ed6:e4410fc9:b27af0b7:0afe758d":     "MDRAID",
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Add(plotter.NewGrid())

	// TODO: command line arg
	fio_data := fiotools.LoadFioJson("thrash.json")

	color_idx := 0
	var ordered_keys []float64
	for _, cstat := range fio_data.ClientStats {
		// warning: only works with unified_rw_reporting enabled
		// TODO: this will all fall apart with split read / write output
		// Will refactor this after I generate some sample output.
		lat := cstat.Mixed.Clat.Percentile

		// fio histograms might have empty slots named "0.00"
		// that will flatten to one 0.00 that we now delete
		delete(lat, 0.00)

		// will get re-run on each iteration, no biggie
		// this var is used again in the main scope at the end!
		ordered_keys = make([]float64, len(lat))
		var i int = 0
		for k, _ := range lat {
			ordered_keys[i] = k
			i++
		}
		// sort the keys to be sure they're in order since they passed
		// through a JSON map and ordering could be screwed up
		sort.Float64s(ordered_keys)

		// copy into the plotinum type
		d := make(plotter.XYs, len(lat))
		for i, key := range ordered_keys {
			d[i].X = key
			d[i].Y = lat[key]
		}

		line, _ := plotter.NewLine(d)
		line.LineStyle.Width = vg.Points(1)
		// QUESTION: when will this turn ugly?
		line.LineStyle.Color = plotutil.Color(color_idx)
		color_idx++

		// attach to the plot, add to the legend
		p.Add(line)
		p.Legend.Add(legend[cstat.Jobname], line)
	}

	// TODO: flag and/or JSON config
	p.Legend.Top = true
	p.Title.Text = "SAS/SATA/SSD/RAID Latency Study"
	p.Y.Label.Text = "IO Latency in Microseconds"
	p.X.Label.Text = "Percentile of IOs"

	// replace the default marker with a custom one that prints the percentiles
	// this will break if fio results aren't mod 10.0 = 0.0
	p.X.Tick.Marker = func(min, max float64) []plot.Tick {
		ticks := make([]plot.Tick, len(ordered_keys))

		for i, val := range ordered_keys {
			if math.Mod(val, 10.0) != 0.0 {
				continue
			}
			ticks[i] = plot.Tick{Value: val, Label: fmt.Sprintf("%.0fth", val)}
		}
		return ticks
	}

	// TODO: flag
	if err := p.Save(10, 10, "test.png"); err != nil {
		panic(err)
	}
	fmt.Println("Wrote test.png.")
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
