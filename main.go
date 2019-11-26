package main

import (
	"fmt"
	// "github.com/google/gopacket"
	"os"
	// "github.com/google/gopacket/layers"
	"github.com/yamatcha/simulator/buffer"
	// "net"
	"io"
	// "log"
	// "time"
	// "sort"
	"flag"
	"strconv"
	// "reflect"
	"encoding/csv"
	"strings"
)

var (
	err error
)

const (
	// csvpath string  = "../csv/http.csv"
	// csvpath string = "../csv/201704122345.csv"
	// csvpath string  = "./opCSV/2019070314002.csv"
	// csvpath string ="./opCSV/sinet.csv"
	// csvpath string ="./opCSV/chicago.csv"
	perSec  float64 = 1.0
	// maxSec  int     = 900
)

func packetTimeBase(csvpath string,buf buffer.Buffers, bufList []string, result buffer.ResultData, params buffer.Params) (buffer.Buffers, []string, buffer.ResultData) {
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
		params.CurrentTime, _ = strconv.ParseFloat(line[1],64)
		buf, bufList, result = buf.CheckBufferTime(bufList, params, result)
		buf, bufList, result = buf.AppendBuffer(bufList, params, fiveTuple, result)
	}
	result.EndFlag = true
	buf, bufList, result = buf.CheckBufferTime(bufList, params, result)
	return buf, bufList, result
}

func globalTimeBase(csvpath string,buf buffer.Buffers, bufList []string, result buffer.ResultData, params buffer.Params, ideal bool) (buffer.Buffers, []string, buffer.ResultData) {
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
		params.CurrentTime, _ = strconv.ParseFloat(line[1],64)
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

func preEval(csvpath string,buf buffer.Buffers, bufList []string, result buffer.ResultData, params buffer.Params) (buffer.Buffers, []string, buffer.ResultData){
		// open csv
	file, err := os.Open(csvpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var line []string
	i := 0
	for i = 0; ; i++ {
		line, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1],64)
		buf, bufList, result = buf.AppendBuffer(bufList, params, fiveTuple, result)
	}
	result.EndFlag = true
	result.PacketNumAll = i
	return buf, bufList, result
}

func getRtt(csvpath string,buf buffer.Buffers, params buffer.Params) float64{
		// open csv
	file, err := os.Open(csvpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rttSum := 0.0
	rttCount := 0

	reader := csv.NewReader(file)
	var line []string
	i := 0
	for i = 0; ; i++ {
		line, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1],64)
		_,ok:=buf[fiveTuple]
		if !ok{
			// new_timelist := []float64{0.0}
			newbuf := buffer.Buffer{params.CurrentTime, 1}
			buf[fiveTuple] = newbuf
			list := strings.Split(fiveTuple," ")
			syn := strings.Join(append(append(list[2:4],list[0:2]...), list[4])," ")
			b, ok := buf[syn]
			// fmt.Println(buf)
			if ok{
				rttCount++
				rttSum+= (params.CurrentTime-b.FirstTime)
			}
		}
	}
	return rttSum/float64(rttCount)
}


// function for display result of list
func printAccessPers(result buffer.ResultData) {
	for i, v := range result.AccessPers {
		fmt.Println(float64(i+1)*(perSec), v)
	}
}

func printAcccessPersTotal(result buffer.ResultData, mode int) {
	sum := 0
	for _, v := range result.AccessPers {
		// fmt.Println(v)
		sum += v
	}
	if mode==0 || mode==1{
		fmt.Println(float64(sum)/900.0)
	}else{
		fmt.Println(float64(sum)/90.0)
	}
}

func printEntryNums(result buffer.ResultData,timeWidth float64) {
	for i, _ := range result.EntryNums[0] {
		fmt.Print(strconv.FormatFloat(float64(i)*timeWidth,'f',2,64) + " ")
		for j:=0;j<10;j++ {
			fmt.Print(strconv.Itoa(result.EntryNums[j][i])+ " ")
		}
		fmt.Println()
	}
}

func printPreEval(buf buffer.Buffers){
	for _, b := range buf{
		fmt.Println(b.Len)
	}
}

func main() {
	buf := buffer.Buffers{}
	bufList := []string{}
	result := buffer.ResultData{0, 0, 0, 0, 0, 0, []int{0}, make([][]int,10), false}
	params := buffer.Params{0.0, 0, 0, 0, 0, false}
	csvpaths := []string{"./opCSV/wide.csv","./opCSV/chicago.csv","./opCSV/sinet.csv"}

	// read time width and buffer size
	flag.Parse()
	csvmode, _ := strconv.Atoi(flag.Arg(0))
	mode, _ := strconv.Atoi(flag.Arg(1))
	params.TimeWidth, _ = strconv.ParseFloat(flag.Arg(2), 64)
	params.BufSize, _ = strconv.Atoi(flag.Arg(3))
	params.EntrySize, _ = strconv.Atoi(flag.Arg(4))

	params.PerSec = perSec
	csvpath := csvpaths[csvmode]
	// fmt.Println(csvpath)

	//select simulator
	if mode == 0 {
		// buf, bufList, result = packetTimeBase(buf, bufList, result, params)
	} else if mode == 1 {
		buf, bufList, result = globalTimeBase(csvpath,buf, bufList, result, params, false)
	} else if mode == 2 {
		//
		buf, bufList, result = globalTimeBase(csvpath,buf, bufList, result, params, true)
	} else if mode == 3 {
		params.Stupid = true
		buf, bufList, result = globalTimeBase(csvpath,buf, bufList, result, params, true)
	} else if mode == 4{
		buf, bufList, result = preEval(csvpath,buf,bufList,result,params)
		printPreEval(buf)
	} else if mode == 5{
		fmt.Println(getRtt(csvpath,buf,params))
	}

	// print result

	// fmt.Println(result.BufMax)
	// fmt.Println(result.EntryNum)

	// printAccessPers(result)
	// printAcccessPersTotal(result,csvmode)
	// printEntryNums(result, params.TimeWidth)
	printPreEval(buf)

	// fmt.Println(result.MaxPacketNum, result.AccessCount, float64(result.AccessCount)/float64(result.PacketNumAll),result.PacketNumAll)

}
