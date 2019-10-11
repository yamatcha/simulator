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

func packetTimeBase(buf buffer.Buffers, bufList []string, result buffer.ResultData, timeWidth float64) (buffer.Buffers, []string, buffer.ResultData) {
	var startTime time.Time
	var currentTime time.Time

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
			startTime = currentTime
		}
		buf, bufList, result = buf.CheckBufferTime(bufList, currentTime, startTime, timeWidth, perSec, result)
		buf, bufList, result = buf.AppendBuffer(bufList, currentTime, fiveTuple,result)
	}
	result.EndFlag = true
	buf, bufList, result = buf.CheckBufferTime(bufList, currentTime, startTime, timeWidth, perSec, result)
	return buf, bufList, result
}

func globalTimeBase(buf buffer.Buffers, bufList []string, result buffer.ResultData, timeWidth float64, bufSize int, ideal bool) (buffer.Buffers, []string, buffer.ResultData) {
	var startTime time.Time
	var currentTime time.Time


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
			startTime = currentTime
		}
		if ideal==false{
			buf, bufList, result = buf.CheckGlobalTime(bufList, startTime, currentTime, timeWidth, perSec, result)
		}else {
			buf, bufList, result = buf.CheckGlobalTimeIdeal(bufList, startTime, currentTime, timeWidth, perSec, result,bufSize)
		}
		buf, bufList, result = buf.AppendBuffer(bufList, currentTime, fiveTuple, result)
	}
	result.EndFlag = true
	result.PacketNumAll = i
	if ideal==false{
		buf, bufList, result = buf.CheckGlobalTime(bufList, startTime, currentTime, timeWidth, perSec, result)
	}else {
		buf, bufList, result = buf.CheckGlobalTimeIdeal(bufList, startTime, currentTime, timeWidth, perSec, result,bufSize)
	}
	return buf, bufList, result
}

func printAccessPers(result buffer.ResultData){
	for i, v := range result.AccessPers {
		fmt.Println(float64(i+1)*(perSec), v)
	}
}

func main() {
	buf := buffer.Buffers{}
	bufList := []string{}
	result := buffer.ResultData{0, 0, 0, 0, 0, 0, []int{0}, false}


	// read time width and buffer size
	flag.Parse()
	mode, _ :=strconv.Atoi(flag.Arg(0))
	timeWidth, _ := strconv.ParseFloat(flag.Arg(1), 64)
	bufSize, _ := strconv.Atoi(flag.Arg(2))

	//select simulator
	if mode==0{
		buf, bufList, result = packetTimeBase(buf,bufList,result,timeWidth)
	}else if mode==1{
		buf, bufList, result = globalTimeBase(buf, bufList, result, timeWidth,bufSize, false)
	}else if mode==2{
		buf, bufList, result = globalTimeBase(buf, bufList, result, timeWidth, bufSize, true)
	}
	
	

	// print result

	// fmt.Println(result.BufMax)

	// printAccessPers(result)

	// fmt.Println(result.MaxPacketNum, result.AccessCount, float64(result.AccessCount)/float64(result.PacketNumAll),result.PacketNumAll)

}
