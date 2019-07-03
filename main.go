package main

import (
	"fmt"
	"github.com/google/gopacket"
	// "github.com/google/gopacket/layers"
	"./buffer"
	"github.com/google/gopacket/pcap"
	// "net"
	"io"
	"log"
	// "time"
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
	pcapFile string = "./pcap/201704122345.pcap"
	// pcapFile string = "./pcap/http.pcap"
)

func main() {
	// open pcap file and call FlowDivide
	handle, err = pcap.OpenOffline(pcapFile)
	buf := buffer.Buffers{}
	buflist := []string{}
	cnt := 0
	// fmt.Println(len(buflist)==0)
	// read time width
	flag.Parse()
	time_width, _ := strconv.ParseFloat(flag.Arg(0), 64)

	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	max := 0
	count := 0
	for {
		packet, err := packetSource.NextPacket()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error:", err)
			continue
		}
		nowtime := buffer.GetTime(packet)
		fivetuple := buffer.GetFiveTuple(packet)
		buffer.Check_buf_time(buf, &buflist, nowtime, time_width, &cnt, &max)
		buffer.Append_buf(&buf, &buflist, nowtime, fivetuple)
		// fmt.Println(buf)
		// for _,v := range buf{
		// 	fmt.Println(*v.TimeList)
		// }
		count++
	}
	fmt.Println(max, cnt, count)
}
