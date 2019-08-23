package main

import (
	"fmt"
	// "github.com/google/gopacket"
	"os"
	// "github.com/google/gopacket/layers"
	"./buffer"
	// "net"
	"io"
	// "log"
	"time"
	// "sort"
	"flag"
	"strconv"
	// "reflect"
	"encoding/csv"
)

var (
	err error
)

const (
	// csvpath string  = "../csv/http.csv"
	// csvpath string = "../csv/201704122345.csv"
	csvpath string  = "../csv/201907031400.csv"
	perSec  float64 = 1.0
	maxSec  int     = 900
)

func packetTimeBase(buf buffer.Buffers, bufList []string, result buffer.Result_data) (buffer.Buffers, []string, buffer.Result_data) {
	var std_time time.Time
	var currentTime time.Time

	// read time width and buffer size
	flag.Parse()
	time_width, _ := strconv.ParseFloat(flag.Arg(0), 64)
	// bufsize, _ := strconv.ParseFloat(flag.Arg(1), 64)

	// open csv
	file, err := os.Open(csvpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var line []string

	//reading csv and do
	i := 0
	for i = 0; ; i++ {
		line, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		currentTime, _ = time.Parse(time.RFC3339Nano, line[1])
		if i == 0 {
			std_time = currentTime
		}
		buf, bufList, result = buf.CheckBufferTime(bufList, currentTime, time_width, result)
		if result.CurrentSecCount != int(float64(maxSec)/perSec) {
			result = buffer.CheckCurrentSec(std_time, currentTime, perSec, result)
		}
		buf, bufList, result = buf.AppendBuffer(bufList, currentTime, fiveTuple,result)
	}
	result.EndFlag = true
	buf, bufList, result = buf.CheckBufferTime(bufList, currentTime, time_width, result)
	return buf, bufList, result
}

func globalTimeBase(buf buffer.Buffers, bufList []string, result buffer.Result_data) (buffer.Buffers, []string, buffer.Result_data) {
	var std_time time.Time
	var currentTime time.Time

	// read time width and buffer size
	flag.Parse()
	time_width, _ := strconv.ParseFloat(flag.Arg(0), 64)
	// bufsize, _ := strconv.ParseFloat(flag.Arg(1), 64)

	// open csv
	file, err := os.Open(csvpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var line []string

	//reading csv and do
	i := 0
	for i = 0; ; i++ {
		line, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		currentTime, _ = time.Parse(time.RFC3339Nano, line[1])
		if i == 0 {
			std_time = currentTime
		}
		buf, bufList, result = buf.CheckGlobalTime(bufList, std_time, currentTime, time_width, result)
		buf, bufList, result = buf.AppendBuffer(bufList, currentTime, fiveTuple, result)
	}
	result.EndFlag = true
	result.PacketNumAll = i
	buf, bufList, result = buf.CheckGlobalTime(bufList, std_time, currentTime, time_width, result)
	return buf, bufList, result
}

func main() {
	buf := buffer.Buffers{}
	bufList := []string{}
	result := buffer.Result_data{0, 0, 0, 0, 0, []int{0}, false}

	// buf, bufList, result = packetTimeBase(buf,bufList,result)
	buf, bufList, result = globalTimeBase(buf, bufList, result)

	// print result

	fmt.Println(result.BufMax)
	// for i, v := range result.AccessPers {
	// 	fmt.Println(float64(i+1)*(perSec), v)
	// }
	// fmt.Println(result.MaxPacketNum, result.AccessCount, float64(result.AccessCount)/float64(result.PacketNumAll),result.PacketNumAll)

}
