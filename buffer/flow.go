package buffer

import (
	"fmt"
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
	// FiveTuple string
	Firstime  time.Time
	TimeList  *[]float64
}

type Buffers map[string]Buffer

type Result_data struct {
	// NumAccess    int
	MaxPacketNum int
	AccessCount  int
	CsCount      int
	AccessPers   []int
	EndFlag      bool
}

func (buf Buffers) CheckBufferTime(bufList []string, currentTime time.Time, time_width float64, result Result_data) (Buffers, []string, Result_data) {
	i := 0
	for _, k := range bufList {
		if GetDuration(buf[k].Firstime, currentTime) > time_width || result.EndFlag == true {
			// fmt.Println(len(*(buf[k].TimeList)))
			result.AccessCount++
			fmt.Println(result.AccessPers)
			fmt.Println(len(result.AccessPers),result.CsCount)
			result.AccessPers[result.CsCount]++
			if result.MaxPacketNum < len(*(buf[k].TimeList)) {
				result.MaxPacketNum = len(*(buf[k].TimeList))
			}
			delete(buf, k)
			i++
			continue
		}
		return buf, bufList[i:],result
	}
	return buf, nil, result
}

func CheckSeconds(std_time time.Time, currentTime time.Time, time_width float64, result Result_data) Result_data {
	if GetDuration(std_time, currentTime) > float64(result.CsCount+1)*time_width {
		result.AccessPers = append(result.AccessPers,0)
		result.CsCount++
		fmt.Println("[",len(result.AccessPers),result.CsCount,"]")
		return result
	}
	return result
}

func (buf Buffers) AppendBuffer(bufList []string, currentTime time.Time, fivetuple string) (Buffers, []string) {
	_, ok := buf[fivetuple]
	if !ok {
		new_timelist := []float64{0.0}
		newbuf := Buffer{currentTime, &new_timelist}
		buf[fivetuple] = newbuf
		bufList = append(bufList, fivetuple)
	} else {
		b := buf[fivetuple]
		*(b.TimeList) = append(*(b.TimeList), GetDuration(b.Firstime, currentTime))
	}
	fmt.Println("{",bufList,"}")
	return buf, bufList
}
