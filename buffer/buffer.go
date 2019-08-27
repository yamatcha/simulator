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

type ResultData struct {
	MaxPacketNum    int
	AccessCount     int
	CurrentSecCount int
	CurrentTimeCount int
	BufMax          int
	PacketNumAll    int
	AccessPers      []int
	EndFlag         bool
}

func GetDuration(first time.Time, now time.Time) float64 {
	return now.Sub(first).Seconds()
}

func (buf Buffers) CheckBufferTime(bufList []string, currentTime time.Time, timeWidth float64, result ResultData) (Buffers, []string, ResultData) {
	i := 0
	for _, k := range bufList {
		if GetDuration(buf[k].FirstTime, currentTime) > timeWidth || result.EndFlag == true {
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

func CheckCurrentSec(std_time time.Time, currentTime time.Time, perSec float64, result ResultData) ResultData {
	if GetDuration(std_time, currentTime) > float64(result.CurrentSecCount+1)*perSec {
		result.AccessPers = append(result.AccessPers, 0)
		result.CurrentSecCount++
		return result
	}
	return result
}

func (buf Buffers) AppendBuffer(bufList []string, currentTime time.Time, fivetuple string, result ResultData) (Buffers, []string, ResultData) {
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

func (buf Buffers) CheckGlobalTime(bufList []string, std_time time.Time, currentTime time.Time, timeWidth float64, perSec float64, result ResultData) (Buffers, []string, ResultData) {
	duration := GetDuration(std_time, currentTime)
	if  duration > float64(result.CurrentSecCount+1)*timeWidth {
		result.AccessCount+=len(bufList)
		result.AccessPers[result.CurrentSecCount]+=len(bufList)
		bufList = []string{}
		buf = Buffers{}
	}
	if duration > float64(result.CurrentTimeCount+1)*perSec{
		result.AccessPers = append(result.AccessPers, 0)
		result.CurrentTimeCount++
	}
	return buf, bufList, result
}

func (buf Buffers) CheckGlobalTimeIdeal(bufList []string, std_time time.Time, currentTime time.Time, timeWidth float64, perSec float64, result ResultData) (Buffers, []string, ResultData) {
	duration := GetDuration(std_time, currentTime)
	if  duration > float64(result.CurrentSecCount+1)*timeWidth {
		result.AccessCount+=len(bufList)
		result.AccessPers[result.CurrentSecCount]+=len(bufList)
		bufList = []string{}
		buf = Buffers{}
	}
	if duration > float64(result.CurrentTimeCount+1)*perSec{
		result.AccessPers = append(result.AccessPers, 0)
		result.CurrentTimeCount++
	}
	return buf, bufList, result
}