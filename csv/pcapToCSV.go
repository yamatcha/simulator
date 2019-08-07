package main

import (
	"./Info"
	"fmt"
	"github.com/google/gopacket"
	// "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	// "net"
	"io"
	"log"
	"time"
	// "sort"
	// "flag"
	// "strconv"
	// "reflect"
	"encoding/csv"
	"os"
	"strings"
)

var (
	handle *pcap.Handle
	err    error
)

const (
	// pcapFile string = "../pcap/201704122345.pcap"
	pcapFile string = "../pcap/201907031400.pcap"
	// pcapFile string  = "../pcap/http.pcap"
	perSec   float64 = 1.0
	maxSec   int     = 900
)

func failOnError(err error) {
    if err != nil {
        log.Fatal("Error:", err)
    }
}

func main() {
	// open pcap file and initialize csv
	handle, err = pcap.OpenOffline(pcapFile)
	var currentTime string
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	name := strings.Split(strings.Split(pcapFile,"/")[2],".")[0]
	file, err := os.OpenFile("./"+name+".csv", os.O_WRONLY|os.O_CREATE, 0600)
	failOnError(err)
	defer file.Close()

    err = file.Truncate(0) // initialize file (for after the second)
    failOnError(err)
	//read pcap and write csv
	writer := csv.NewWriter(file)
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	i:=0
	for ; ;  {
		packet, err := packetSource.NextPacket()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error:", err)
			continue
		}
		currentTime = Info.GetTime(packet).Format(time.RFC3339Nano)
		fiveTuple := Info.GetFiveTuple(packet)
		if fiveTuple !=""{
			i++
			// fmt.Println(fiveTuple,currentTime)
			writer.Write([]string{fiveTuple,currentTime})
		}
	}
	fmt.Println(i)
	writer.Flush()
}
