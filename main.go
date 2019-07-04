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
	pcapFile string = "./pcap/201704122345.pcap"
	// pcapFile string = "./pcap/http.pcap"
)

func main() {
	// open pcap file and call FlowDivide
	handle, err = pcap.OpenOffline(pcapFile)
	buf := buffer.Buffers{}
	buflist := []string{}
	access_cnt := 0
	max := 0
	count := 0
	num_access :=0
	access_pers :=[]int{}
	var std_time time.Time
	cs_count:=0.0
	per_s:=1.0
	// fmt.Println(len(buflist)==0)
	// read time width
	flag.Parse()
	time_width, _ := strconv.ParseFloat(flag.Arg(0), 64)

	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
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
		if count==0{
			std_time = nowtime
		}
		buffer.Check_buf_time(buf, &buflist, nowtime, time_width, &access_cnt, &max, &num_access)
		cs_count=buffer.Check_seconds(std_time,nowtime,per_s,&num_access,&access_pers,cs_count)
		buffer.Append_buf(&buf, &buflist, nowtime, fivetuple)
		count++
	}

	// fmt.Println(max, access_cnt, count)
	// for i,v := range access_pers{
	// 	fmt.Println(float64(i+1)*(per_s),v)
	// }
}
