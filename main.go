package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yamatcha/simulator/buffer"
	"github.com/yamatcha/simulator/opCSV"
	"github.com/yamatcha/simulator/simulation"

	"encoding/csv"
)

var (
	err error
)

const (
	perSec float64 = 1.0
)

func main() {

	var (
		csvpath      = flag.String("csvPath", "", "csv path")
		mode         = flag.Int("mode", 0, "simulator mode")
		timeWidth    = flag.Float64("timeWidth", 0, "time width")
		bufSize      = flag.Int("bufSize", 0, "the number of buffers")
		entrySize    = flag.Int("entrySize", 0, "the number of entries per buffer")
		protocol     = flag.String("protocol", "", "L3 Protocol UDP or TCP")
		selectedPort = flag.String("selectedPort", "", "the port targeted flow chunk buffer")
		pcapPath     = flag.String("pcapPath", "", "path of pcap converted to csv")
	)

	buf := buffer.Buffers{}
	bufList := []string{}
	result := buffer.ResultData{AccessPers: []int{0}}

	flag.Parse()
	params := buffer.Params{TimeWidth: *timeWidth, BufSize: *bufSize, EntrySize: *entrySize, Protocol: *protocol, SelectedPort: strings.Split(*selectedPort, ",")}

	params.PerSec = perSec

	// open csv
	var reader *csv.Reader
	if *mode != 7 {
		file, err := os.Open(*csvpath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader = csv.NewReader(file)
	}

	//select simulator
	switch *mode {
	case 0:
		buf, bufList, result = simulation.GlobalTimeBase(reader, buf, bufList, result, params, false)
	case 1:
		buf, bufList, result = simulation.GlobalTimeBase(reader, buf, bufList, result, params, true)
	case 2:
		params.Stupid = true
		buf, bufList, result = simulation.GlobalTimeBase(reader, buf, bufList, result, params, true)
		printAcccessPersAvg(result)
	case 3:
		buf, bufList, result = simulation.PreEval(reader, buf, bufList, result, params)
	case 4:
		fmt.Println(simulation.GetRtt(reader, buf, params))
	case 5:
		simulation.GetWindow(reader, buf, params)
	case 6:
		simulation.Protocol(reader, buf, bufList, result, params)
	case 7:
		opCSV.PcapToCSV(strings.Split(*pcapPath, ","))
	}

}
func printAcccessPersAvg(result buffer.ResultData) {
	sum := 0
	fmt.Println(result.AccessPers)
	for _, v := range result.AccessPers {
		fmt.Println(v)
		sum += v
	}
	fmt.Println(float64(sum) / 900.0)
}
