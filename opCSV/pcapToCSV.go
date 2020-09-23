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
	// pcapFile := []string{"./chicagoA/20140320-130000", "./chicagoA/20140320-130100", "./chicagoA/20140320-130200", "./chicagoA/20140320-130300", "./chicagoA/20140320-130400", "./chicagoA/20140320-130500", "./chicagoA/20140320-130600", "./chicagoA/20140320-130700", "./chicagoA/20140320-130800", "./chicagoA/20140320-130900", "./chicagoA/20140320-131000", "./chicagoA/20140320-131100", "./chicagoA/20140320-131200", "./chicagoA/20140320-131300", "./chicagoA/20140320-131400"}
	// open pcap file and initialize csv

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
			if fiveTuple != "" {
				currentTime = Info.GetTime(packet)
				if i < 5 {

					fmt.Println(startTime, currentTime, strconv.FormatFloat(Info.GetDuration(startTime, currentTime), 'f', 6, 64))
				}
				timeStamp := strconv.FormatFloat(Info.GetDuration(startTime, currentTime), 'f', 6, 64)
				writer.Write([]string{fiveTuple, timeStamp})
				i++
			}
		}
	}
	writer.Flush()
}
