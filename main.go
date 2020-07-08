package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/yamatcha/simulator/buffer"

	"encoding/csv"
	"strings"
)

var (
	err error
)

const (
	perSec float64 = 1.0
	// maxSec  int     = 900
)

func globalTimeBase(csvReader *csv.Reader, buf buffer.Buffers, bufList []string, result buffer.ResultData, params buffer.Params, ideal bool) (buffer.Buffers, []string, buffer.ResultData) {
	for ; ; result.PacketNumAll++ {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)
		if ideal == false {
			buf, bufList, result = buf.CheckGlobalTimeWithUnlimitedBuffers(bufList, params, result)
		} else {
			buf, bufList, result = buf.CheckGlobalTime(bufList, params, result)
		}
		buf, bufList, result = buf.Append(bufList, params, fiveTuple, result)
	}
	result.EndFlag = true
	if ideal == false {
		buf, bufList, result = buf.CheckGlobalTimeWithUnlimitedBuffers(bufList, params, result)
	} else {
		buf, bufList, result = buf.CheckGlobalTime(bufList, params, result)
	}
	return buf, bufList, result
}

func preEval(csvReader *csv.Reader, buf buffer.Buffers, bufList []string, result buffer.ResultData, params buffer.Params) (buffer.Buffers, []string, buffer.ResultData) {
	for ; ; result.PacketNumAll++ {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)
		buf, bufList, result = buf.Append(bufList, params, fiveTuple, result)
	}
	result.EndFlag = true
	return buf, bufList, result
}

func getRtt(csvReader *csv.Reader, buf buffer.Buffers, params buffer.Params) float64 {
	rttSum := 0.0
	rttCount := 0

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)

		_, ok := buf[fiveTuple]
		if !ok {
			newbuf := buffer.Buffer{params.CurrentTime, 1}
			buf[fiveTuple] = newbuf
			list := strings.Split(fiveTuple, " ")
			syn := strings.Join(append(append(list[2:4], list[0:2]...), list[4]), " ")
			b, ok := buf[syn]
			if ok {
				rttCount++
				rttSum += (params.CurrentTime - b.FirstTime)
			}
		}
	}
	return rttSum / float64(rttCount)
}

func getWindow(csvReader *csv.Reader, buf buffer.Buffers, params buffer.Params) float64 {
	rttSum := 0
	rttCount := 0

	type flowWindow struct {
		sum   int
		count int
	}

	flowWindows := map[string]flowWindow{}

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)
		b, ok := buf[fiveTuple]
		if !ok {
			newbuf := buffer.Buffer{params.CurrentTime, 1}
			buf[fiveTuple] = newbuf
			list := strings.Split(fiveTuple, " ")
			syn := strings.Join(append(append(list[2:4], list[0:2]...), list[4]), " ")
			_, ok := buf[syn]
			// fmt.Println(buf)
			if ok {
				// rttCount++
				// rttSum+= buf[syn].Len
				f, ok := flowWindows[syn]
				if ok {
					flowWindows[syn] = flowWindow{f.sum + buf[syn].Len, f.count + 1}
				} else {
					flowWindows[syn] = flowWindow{buf[syn].Len, 1}
				}
				delete(buf, syn)
			}
		} else {
			newbuf := buffer.Buffer{b.FirstTime, b.Len + 1}
			buf[fiveTuple] = newbuf
		}
	}
	for _, window := range flowWindows {
		fmt.Println(float64(window.sum) / float64(window.count))
	}
	return float64(rttSum) / float64(rttCount)
}

func main() {

	var (
		csvpath   = *flag.String("path", "", "csv path")
		mode      = *flag.Int("mode", 0, "simulator mode")
		timeWidth = *flag.Float64("timeWidth", 0, "time width")
		bufSize   = *flag.Int("bufsize", 0, "the number of buffers")
		entrySize = *flag.Int("entrysize", 0, "the number of entries per buffer")
	)

	buf := buffer.Buffers{}
	bufList := []string{}
	result := buffer.ResultData{AccessPers: []int{0}}
	params := buffer.Params{TimeWidth: timeWidth, BufSize: bufSize, EntrySize: entrySize}

	// read time width and buffer size
	flag.Parse()

	params.PerSec = perSec

	// open csv
	file, err := os.Open(csvpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	//select simulator
	switch mode {
	case 0:
		buf, bufList, result = globalTimeBase(reader, buf, bufList, result, params, false)
	case 1:
		buf, bufList, result = globalTimeBase(reader, buf, bufList, result, params, true)
	case 2:
		params.Stupid = true
		buf, bufList, result = globalTimeBase(reader, buf, bufList, result, params, true)
	case 3:
		buf, bufList, result = preEval(reader, buf, bufList, result, params)
	case 4:
		fmt.Println(getRtt(reader, buf, params))
	case 5:
		getWindow(reader, buf, params)
	}

}
