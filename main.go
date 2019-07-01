package main

import (
	"fmt"
	"github.com/google/gopacket"
	// "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"./buffer"
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
	// pcapFile string = "./pcap/201704122345.pcap"
	pcapFile string = "./pcap/http.pcap"
)



func main() {
	// open pcap file and call FlowDivide
	handle, err = pcap.OpenOffline(pcapFile)
	buf := []*buffer.Buffer{}

	cnt:=0
	// read time width
	flag.Parse()
	time_width,_:=strconv.ParseFloat(flag.Arg(0),64)

	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		packet, err := packetSource.NextPacket()
		// fmt.Println(packet)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error:", err)
			continue
		}
		nowtime:=buffer.GetTime(packet)
		buffer.Check_buf_time(buf,nowtime,time_width,cnt)
		buffer.Append_buf(&packet,&buf)
		for i,v :=range buf{
			fmt.Println("hoge")
			fmt.Println(i,v)
		}
		fmt.Println()
	}
}