package buffer

import (
	// "fmt"
	// "github.com/google/gopacket"
	// "github.com/google/gopacket/layers"
	// "github.com/google/gopacket/pcap"
	// "net"
	// "io"
	// "log"
	// "reflect"
	"time"
)

type Buffer struct {
	FirstTime time.Time
	TimeList  *[]float64
}

type Buffers map[string]Buffer

type Result_data struct {
	MaxPacketNum    int
	AccessCount     int
	CurrentSecCount int
	BufMax          int
	PacketNumAll    int
	AccessPers      []int
	EndFlag         bool
}

func GetDuration(first time.Time, now time.Time) float64 {
	return now.Sub(first).Seconds()
}

func (buf Buffers) CheckBufferTime(bufList []string, currentTime time.Time, time_width float64, result Result_data) (Buffers, []string, Result_data) {
	i := 0
	for _, k := range bufList {
		if GetDuration(buf[k].FirstTime, currentTime) > time_width || result.EndFlag == true {
			result.AccessCount++
			result.AccessPers[result.CurrentSecCount]++
			if result.MaxPacketNum < len(*(buf[k].TimeList)) {
				result.MaxPacketNum = len(*(buf[k].TimeList))
			}
			// fmt.Println(len(*(buf[k].TimeList)))
			delete(buf, k)
			i++
			continue
		}
		return buf, bufList[i:], result
	}
	return buf, nil, result
}

func CheckCurrentSec(std_time time.Time, currentTime time.Time, time_width float64, result Result_data) Result_data {
	if GetDuration(std_time, currentTime) > float64(result.CurrentSecCount+1)*time_width {
		result.AccessPers = append(result.AccessPers, 0)
		result.CurrentSecCount++
		return result
	}
	return result
}

func (buf Buffers) AppendBuffer(bufList []string, currentTime time.Time, fivetuple string, result Result_data) (Buffers, []string, Result_data) {
	_, ok := buf[fivetuple]
	if !ok {
		new_timelist := []float64{0.0}
		newbuf := Buffer{currentTime, &new_timelist}
		buf[fivetuple] = newbuf
		bufList = append(bufList, fivetuple)
	} else {
		b := buf[fivetuple]
		*(b.TimeList) = append(*(b.TimeList), GetDuration(b.FirstTime, currentTime))

		if len(bufList) > result.BufMax {
			result.BufMax = len(bufList)
		}
	}
	return buf, bufList, result
}

//using in Global time base func

func (buf Buffers) CheckGlobalTime(bufList []string, std_time time.Time, currentTime time.Time, time_width float64, result Result_data) (Buffers, []string, Result_data) {
	if GetDuration(std_time, currentTime) > float64(result.CurrentSecCount+1)*time_width {
		result.AccessPers = append(result.AccessPers, 0)
		result.CurrentSecCount++
		bufList = []string{}
		buf = Buffers{}
	}
	result.AccessCount++
	result.AccessPers[result.CurrentSecCount]++
	return buf, bufList, result
}
