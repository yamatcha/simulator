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
	var line []string

	//reading csv and do
	i := 0
	for i = 0; ; i++ {
		line, err = csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)
		if ideal == false {
			buf, bufList, result = buf.CheckGlobalTime(bufList, params, result)
		} else {
			buf, bufList, result = buf.CheckGlobalTimeIdeal(bufList, params, result)
		}
		buf, bufList, result = buf.AppendBuffer(bufList, params, fiveTuple, result)
	}
	result.EndFlag = true
	result.PacketNumAll = i
	if ideal == false {
		buf, bufList, result = buf.CheckGlobalTime(bufList, params, result)
	} else {
		buf, bufList, result = buf.CheckGlobalTimeIdeal(bufList, params, result)
	}
	return buf, bufList, result
}

func preEval(csvReader *csv.Reader, buf buffer.Buffers, bufList []string, result buffer.ResultData, params buffer.Params) (buffer.Buffers, []string, buffer.ResultData) {
	var line []string
	for ; ; result.PacketNumAll++ {
		line, err = csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)
		buf, bufList, result = buf.AppendBuffer(bufList, params, fiveTuple, result)
	}
	result.EndFlag = true
	return buf, bufList, result
}

func getRtt(csvReader *csv.Reader, buf buffer.Buffers, params buffer.Params) float64 {
	rttSum := 0.0
	rttCount := 0

	var line []string
	i := 0
	for i = 0; ; i++ {
		line, err = csvReader.Read()
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

	var line []string
	i := 0
	for i = 0; ; i++ {
		line, err = csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)
		b, ok := buf[fiveTuple]
		if !ok {
			// new_timelist := []float64{0.0}
			newbuf := buffer.Buffer{params.CurrentTime, 1}
			buf[fiveTuple] = newbuf
			list := strings.Split(fiveTuple, " ")
			syn := strings.Join(append(append(list[2:4], list[0:2]...), list[4]), " ")
			_, ok := buf[syn]
			// fmt.Println(buf)
			if ok {
				// rttCount++
				// rttSum+= buf[syn].Len
				// fmt.Println(buf[syn].Len)
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
	buf := buffer.Buffers{}
	bufList := []string{}
	result := buffer.ResultData{0, 0, 0, 0, 0, 0, []int{0}, false}
	params := buffer.Params{}
	csvpaths := []string{"./opCSV/wide.csv", "./opCSV/chicago.csv", "./opCSV/sinet.csv"}

	// read time width and buffer size
	flag.Parse()
	csvmode, _ := strconv.Atoi(flag.Arg(0))
	mode, _ := strconv.Atoi(flag.Arg(1))
	params.TimeWidth, _ = strconv.ParseFloat(flag.Arg(2), 64)
	params.BufSize, _ = strconv.Atoi(flag.Arg(3))
	params.EntrySize, _ = strconv.Atoi(flag.Arg(4))

	params.PerSec = perSec
	csvpath := csvpaths[csvmode]

	// open csv
	file, err := os.Open(csvpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	//select simulator
	switch mode {
	case 1:
		buf, bufList, result = globalTimeBase(reader, buf, bufList, result, params, false)
	case 2:
		buf, bufList, result = globalTimeBase(reader, buf, bufList, result, params, true)
	case 3:
		params.Stupid = true
		buf, bufList, result = globalTimeBase(reader, buf, bufList, result, params, true)
	case 4:
		buf, bufList, result = preEval(reader, buf, bufList, result, params)
	case 5:
		fmt.Println(getRtt(reader, buf, params))
	case 6:
		// fmt.Println("result:",getWindow(reader,buf,params))
		getWindow(reader, buf, params)
	}

}
