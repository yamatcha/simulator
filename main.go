package main

import (
	// "fmt"
	"github.com/google/gopacket"
	// "github.com/google/gopacket/layers"
	"./buffer"
	"github.com/google/gopacket/pcap"
	// "net"
	"io"
	"log"
	"time"
	// "sort"
	"flag"
	"strconv"
	// "reflect"
)

var (
	handle *pcap.Handle
	err    error
)

const (
	// pcapFile string = "./pcap/201704122345.pcap"
	pcapFile string = "./pcap/201907031400.pcap"
	// pcapFile string  = "./pcap/http.pcap"
	perSec   float64 = 1.0
	maxSec   int     = 900
)

func main() {
	// open pcap file and call FlowDivide
	handle, err = pcap.OpenOffline(pcapFile)
	buf := buffer.Buffers{}
	bufList := []string{}
	var std_time time.Time
	var currentTime time.Time
	result := buffer.Result_data{0,0,0,[]int{0},false}
	// read time width
	flag.Parse()
	time_width, _ := strconv.ParseFloat(flag.Arg(0), 64)
	// bufsize, _ := strconv.ParseFloat(flag.Arg(1), 64)

	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	i:=0
	for i = 0; ; i++ {
		packet, err := packetSource.NextPacket()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error:", err)
			continue
		}
		currentTime = buffer.GetTime(packet)
		fiveTuple := buffer.GetFiveTuple(packet)
		if i == 0 {
			std_time = currentTime
		}
		buf, bufList, result = buf.CheckBufferTime(bufList, currentTime, time_width, result)
		if result.CsCount != int(float64(maxSec)/perSec) {
			result = buffer.CheckSeconds(std_time, currentTime, perSec, result)
		}
		buf, bufList = buf.AppendBuffer(bufList, currentTime, fiveTuple)
	}
	result.EndFlag = true
	buf, bufList, result = buf.CheckBufferTime(bufList, currentTime, time_width, result)
	// for i,v := range result.AccessPers{
	// 	fmt.Println(float64(i+1)*(perSec),v)
	// }
	// fmt.Println(result.MaxPacketNum, result.AccessCount, float64(result.AccessCount)/float64(i),i)
}
