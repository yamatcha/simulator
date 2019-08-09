package main

import (
	"fmt"
	// "github.com/google/gopacket"
	"os"
	// "github.com/google/gopacket/layers"
	"./buffer"
	// "net"
	"io"
	"log"
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
	csvpath string = "../csv/201907031400.csv"
	perSec  float64 = 1.0
	maxSec  int     = 900
)

func main() {
	// open csv and call FlowDivide
	buf := buffer.Buffers{}
	bufList := []string{}
	var std_time time.Time
	var currentTime time.Time
	result := buffer.Result_data{0, 0, 0, 0, []int{0}, false}

	// read time width and buffer size
	flag.Parse()
	time_width, _ := strconv.ParseFloat(flag.Arg(0), 64)
	// bufsize, _ := strconv.ParseFloat(flag.Arg(1), 64)

	file, err := os.Open(csvpath)
	if err != nil {
		panic(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var line []string

	i := 0
	for i = 0; ; i++ {
		line, err = reader.Read()
		if err==io.EOF{
			break
		}else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		currentTime, _ = time.Parse(time.RFC3339Nano, line[1])
		// fmt.Println("["+line[1]+"]")
		// fmt.Println(currentTime)
		if i == 0 {
			std_time = currentTime
		}
		buf, bufList, result = buf.CheckBufferTime(bufList, currentTime, time_width, result)
		if result.CurrentSecCount != int(float64(maxSec)/perSec) {
			result = buffer.CheckCurrentSec(std_time, currentTime, perSec, result)
		}
		buf, bufList = buf.AppendBuffer(bufList, currentTime, fiveTuple)
	}
	result.EndFlag = true
	buf, bufList, result = buf.CheckBufferTime(bufList, currentTime, time_width, result)

	// print result

	fmt.Println(result.BufMax)
	// for i,v := range result.AccessPers{
	// 	fmt.Println(float64(i+1)*(perSec),v)
	// }
	fmt.Println(result.MaxPacketNum, result.AccessCount, float64(result.AccessCount)/float64(i),i)
}
