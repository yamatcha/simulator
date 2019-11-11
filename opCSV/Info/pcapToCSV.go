package Info



import (
	// "github.com/yamatcha/simulator/opCSV/Info"
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
	"strconv"
	// "strings"
)

var (
	handle *pcap.Handle
	err    error
)

const (
	// pcapFile string = "../pcap/201704122345.pcap"
	// pcapFile string = "../../pcap/201907031400.pcap"
	// pcapFile string  = "../pcap/http.pcap"
	perSec float64 = 1.0
	maxSec int     = 900
)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func PcapToCSV() {
	pcapFile := []string{"./chicago/20140320-130000","./chicago/20140320-130100","./chicago/20140320-130200","./chicago/20140320-130300","./chicago/20140320-130400","./chicago/20140320-130500","./chicago/20140320-130600","./chicago/20140320-130700","./chicago/20140320-130800","./chicago/20140320-130900","./chicago/20140320-131000","./chicago/20140320-131100","./chicago/20140320-131200","./chicago/20140320-131300","./chicago/20140320-131400","./chicago/20140320-131500"}
	// open pcap file and initialize csv

	// name := strings.Split(strings.Split(pcapFile, "/")[3], ".")[0]
	name := "chicago20140320-1300"
	file, err := os.OpenFile("./"+name+".csv", os.O_WRONLY|os.O_CREATE, 0600)
	failOnError(err)
	defer file.Close()

	err = file.Truncate(0) // initialize file (for after the second)
	failOnError(err)
	//read pcap and write csv
	writer := csv.NewWriter(file)
	var startTime time.Time
	for j := 0; j < len(pcapFile); j++ {
		handle, err = pcap.OpenOffline(pcapFile[j]+".pcap")
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
			if j == 0 && i==0 {
				fmt.Println("hoge")
				currentTime = GetTime(packet)
				startTime = currentTime
			}
			fiveTuple := GetFiveTuple(packet)
			if fiveTuple != "" {
				currentTime = GetTime(packet)
				if i < 5 {

					fmt.Println(startTime, currentTime, strconv.FormatFloat(GetDuration(startTime,currentTime), 'f', 6, 64))
				}
				timeStamp := strconv.FormatFloat(GetDuration(startTime,currentTime), 'f', 6, 64)
				writer.Write([]string{fiveTuple, timeStamp})
				i++
			}
		}
	}
	writer.Flush()
}
