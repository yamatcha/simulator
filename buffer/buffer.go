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
	len 		int
}

type Buffers map[string]Buffer

type ResultData struct {
	MaxPacketNum    int
	AccessCount     int
	CurrentTimeCount int
	BufMax          int
	PacketNumAll    int
	PacketOfAllBuffers int
	BiggestBufferFiveTuple string
	AccessPers      []int
	EndFlag         bool
}

func GetDuration(first time.Time, now time.Time) float64 {
	return now.Sub(first).Seconds()
}

func (buf Buffers) CheckBufferTime(bufList []string, currentTime time.Time, startTime time.Time, timeWidth float64, perSec float64, result ResultData) (Buffers, []string, ResultData) {
	i := 0
	for _, k := range bufList {
		if GetDuration(buf[k].FirstTime, currentTime) > timeWidth || result.EndFlag == true {
			result.AccessCount++
			currentSec := int(GetDuration(startTime, buf[k].FirstTime.Add(time.Duration(timeWidth*1000)))/perSec)
			if currentSec>=len(result.AccessPers){
				result.AccessPers = append(result.AccessPers, 0)				
			}
			result.AccessPers[currentSec]++
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

func (buf Buffers) AppendBuffer(bufList []string, currentTime time.Time, fivetuple string, result ResultData) (Buffers, []string, ResultData) {
	_, ok := buf[fivetuple]
	result.PacketOfAllBuffers++
	if !ok {
		new_timelist := []float64{0.0}
		newbuf := Buffer{currentTime, &new_timelist, 1}
		buf[fivetuple] = newbuf
		bufList = append(bufList, fivetuple)

		//for ideal simulator
		if buf[result.BiggestBufferFiveTuple].len<buf[fivetuple].len{
			result.BiggestBufferFiveTuple=fivetuple
		}
	} else {
		b := buf[fivetuple]
		*(b.TimeList) = append(*(b.TimeList), GetDuration(b.FirstTime, currentTime))
		b.len++
		if len(bufList) > result.BufMax {
			result.BufMax = len(bufList)
		}
	}
	return buf, bufList, result
}

//using in Global time base func

func (buf Buffers) CheckGlobalTime(bufList []string, firstTime time.Time, currentTime time.Time, timeWidth float64, perSec float64, result ResultData) (Buffers, []string, ResultData) {
	duration := GetDuration(firstTime, currentTime)
	if  duration > float64(result.CurrentTimeCount+1)*timeWidth {
		// fmt.Println(len(bufList))
		result.AccessCount+=len(bufList)
		result.AccessPers[len(result.AccessPers)-1]+=len(bufList)
		bufList = []string{}
		buf = Buffers{}
		result.CurrentTimeCount++
	}
	if duration > float64(len(result.AccessPers))*perSec{
		result.AccessPers = append(result.AccessPers, 0)
	}
	return buf, bufList, result
}

func (buf Buffers) CheckGlobalTimeIdeal(bufList []string, firstTime time.Time, currentTime time.Time, timeWidth float64, perSec float64, result ResultData) (Buffers, []string, ResultData) {
	duration := GetDuration(firstTime, currentTime)
	if  duration > float64(result.CurrentTimeCount+1)*timeWidth {
		accessCount := result.PacketOfAllBuffers-buf[result.BiggestBufferFiveTuple].len
		result.AccessCount+= accessCount
		result.AccessPers[len(result.AccessPers)]+= accessCount
		result.AccessCount+=1
		delete(buf,result.BiggestBufferFiveTuple)
		bufList = []string{}
		buf = Buffers{}
		result.BiggestBufferFiveTuple=""
		result.PacketOfAllBuffers=0
		result.CurrentTimeCount++
	}
	if duration > float64(len(result.AccessPers)+1)*perSec{
		result.AccessPers = append(result.AccessPers, 0)
		len(result.AccessPers)++
	}
	return buf, bufList, result
}