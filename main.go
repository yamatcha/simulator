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
	result := buffer.ResultData{AccessPerSecList: []int{0}, NextAccessTime: 1}

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
		params.IsLimited = false
		buf, result = simulation.GlobalTimeBase(reader, buf, result, params)
	case 1:
		params.IsLimited = true
		buf, result = simulation.GlobalTimeBase(reader, buf, result, params)
		printAcccessPersAvg(result)
	case 2:
		params.IsLimited = true
		params.IsStupid = true
		buf, result = simulation.GlobalTimeBase(reader, buf, result, params)
		printAcccessPersAvg(result)
	case 3:
		buf, result = simulation.PreEval(reader, buf, result, params)
	case 4:
		fmt.Println(simulation.GetRtt(reader, buf, params))
	case 5:
		simulation.GetWindow(reader, buf, params)
	case 6:
		simulation.Protocol(reader, buf, result, params)
	case 7:
		opCSV.PcapToCSV(strings.Split(*pcapPath, ","))
	}

}
func printAcccessPersAvg(result buffer.ResultData) {
	sum := 0
	for _, v := range result.AccessPerSecList {
		//		fmt.Println(i, v)
		sum += v
	}
	fmt.Println(float64(sum) / 900.0)
}
