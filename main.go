package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yamatcha/simulator/buffer"
	"github.com/yamatcha/simulator/simulation"

	"encoding/csv"
)

var (
	err error
)

const (
	perSec float64 = 1.0
	// maxSec  int     = 900
)

func main() {

	var (
		csvpath      = flag.String("path", "", "csv path")
		mode         = flag.Int("mode", 0, "simulator mode")
		timeWidth    = flag.Float64("timeWidth", 0, "time width")
		bufSize      = flag.Int("bufsize", 0, "the number of buffers")
		entrySize    = flag.Int("entrysize", 0, "the number of entries per buffer")
		protocol     = flag.String("protocol", "", "L3 Protocol UDP or TCP")
		selectedPort = flag.String("selectedPort", "", "the port targeted flow chunk buffer")
	)

	buf := buffer.Buffers{}
	bufList := []string{}
	result := buffer.ResultData{AccessPers: []int{0}}
	params := buffer.Params{TimeWidth: *timeWidth, BufSize: *bufSize, EntrySize: *entrySize, Protocol: *protocol, SelectedPort: strings.Split(*selectedPort, ",")}

	// read time width and buffer size
	flag.Parse()

	params.PerSec = perSec

	// open csv
	file, err := os.Open(*csvpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

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
	}

}
func printAcccessPersAvg(result buffer.ResultData) {
	sum := 0
	for _, v := range result.AccessPers {
		// fmt.Println(v)
		sum += v
	}
	fmt.Println(float64(sum) / 900.0)
}
