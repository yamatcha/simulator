package opCSV

import (
	"fmt"
	"strings"

	"github.com/google/gopacket"
	"github.com/yamatcha/simulator/opCSV/Info"

	"io"
	"log"
	"time"

	"github.com/google/gopacket/pcap"

	"encoding/csv"
	"os"
	"strconv"
)

var (
	handle *pcap.Handle
	err    error
)

const (
	perSec float64 = 1.0
	maxSec int     = 900
)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func PcapToCSV(pcapFile []string) {
	name := strings.Split(strings.Split(pcapFile[0], "/")[3], ".")[0]
	file, err := os.OpenFile("./"+name+".csv", os.O_WRONLY|os.O_CREATE, 0600)
	failOnError(err)
	defer file.Close()

	err = file.Truncate(0) // initialize file (for after the second)
	failOnError(err)
	//read pcap and write csv
	writer := csv.NewWriter(file)
	var startTime time.Time
	for j := 0; j < len(pcapFile); j++ {
		handle, err = pcap.OpenOffline(pcapFile[j] + ".pcap")
		var currentTime time.Time
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()

		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for i := 0; ; {
			packet, err := packetSource.NextPacket()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Println("Error:", err)
				continue
			}
			if j == 0 && i == 0 {
				currentTime = Info.GetTime(packet)
				startTime = currentTime
			}
			fiveTuple := Info.GetFiveTuple(packet)
			if fiveTuple.String() != "" {
				currentTime = Info.GetTime(packet)
				if i < 5 {

					fmt.Println(startTime, currentTime, strconv.FormatFloat(Info.GetDuration(startTime, currentTime), 'f', 6, 64))
				}
				timeStamp := strconv.FormatFloat(Info.GetDuration(startTime, currentTime), 'f', 6, 64)
				writer.Write([]string{fiveTuple.String(), timeStamp})
				i++
			}
		}
	}
	writer.Flush()
}

func CsvGenForCacheSimulator(pcapFile []string) {
	name := strings.Split(strings.Split(pcapFile[0], "/")[3], ".")[0]
	file, err := os.OpenFile("./"+name+".csv", os.O_WRONLY|os.O_CREATE, 0600)
	failOnError(err)
	defer file.Close()

	err = file.Truncate(0)
	failOnError(err)
	writer := csv.NewWriter(file)
	var startTime time.Time
	for j := 0; j < len(pcapFile); j++ {
		handle, err = pcap.OpenOffline(pcapFile[j] + ".pcap")
		var currentTime time.Time
		failOnError(err)

		defer handle.Close()

		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for i := 0; ; {
			packet, err := packetSource.NextPacket()
			if err == io.EOF {
				break
			} else {
				failOnError(err)
			}
			if packet.TransportLayer() == nil || packet.NetworkLayer() == nil {
				continue
			}
			if j == 0 && i == 0 {
				currentTime = Info.GetTime(packet)
				startTime = currentTime
			}
			fiveTuple := Info.GetFiveTuple(packet)
			if fiveTuple.String() != "" {
				currentTime = Info.GetTime(packet)
				len := strconv.Itoa(Info.GetLen(packet))
				timeStamp := strconv.FormatFloat(Info.GetDuration(startTime, currentTime), 'f', 6, 64)
				writer.Write(append(append([]string{timeStamp}, len), fiveTuple.ForCacheSimulator()...))
				i++
			}
		}
	}
	writer.Flush()
}
