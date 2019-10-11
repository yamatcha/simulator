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
	"sort"
	"time"
)

type Buffer struct {
	FirstTime time.Time
	TimeList  []float64
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
			if result.MaxPacketNum < len((buf[k].TimeList)) {
				result.MaxPacketNum = len((buf[k].TimeList))
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
		newbuf := Buffer{currentTime, new_timelist, 1}
		buf[fivetuple] = newbuf
		bufList = append(bufList, fivetuple)
	} else {
		b := buf[fivetuple]
		newbuf := Buffer{b.FirstTime,(append((b.TimeList), GetDuration(b.FirstTime, currentTime))),b.len+1}
		// *(b.TimeList) = append(*(b.TimeList), GetDuration(b.FirstTime, currentTime))
		// b.len++
		buf[fivetuple]=newbuf
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
		fmt.Println(float64(result.CurrentTimeCount+1)*timeWidth ,len(bufList))
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

//for get sorted map

type Entry struct {
    name  string
    value int
}
type List []Entry

func (buf Buffers) getSortedMap(bufSize int) List{
	sortedMap:= List{}
	for k,v := range buf{
		element:=Entry{k,v.len}
		sortedMap=append(sortedMap,element)
	}
	sort.Sort(sort.Reverse(sortedMap))
	return sortedMap[:bufSize]
}

func (l List) Len() int {
    return len(l)
}

func (l List) Swap(i, j int) {
    l[i], l[j] = l[j], l[i]
}

func (l List) Less(i, j int) bool {
    if l[i].value == l[j].value {
        return (l[i].name < l[j].name)
    } else {
        return (l[i].value < l[j].value)
    }
}

func (l List) getListSum() int{
	sum:=0
	for _,v :=range l{
		sum+=v.value
	}
	return sum
}

// 

func (buf Buffers) CheckGlobalTimeIdeal(bufList []string, firstTime time.Time, currentTime time.Time, timeWidth float64, perSec float64, result ResultData, bufSize int) (Buffers, []string, ResultData) {
	duration := GetDuration(firstTime, currentTime)
	if bufSize>len(buf){
		bufSize=len(buf)
	}
	if  duration > float64(result.CurrentTimeCount+1)*timeWidth {
		sortedMap:=buf.getSortedMap(bufSize)
		bufferedPacket := sortedMap.getListSum()
		accessCount := result.PacketOfAllBuffers-bufferedPacket+bufSize
		fmt.Println(accessCount,len(bufList))
		result.AccessCount+= accessCount
		result.AccessPers[len(result.AccessPers)-1]+= accessCount
		result.AccessCount+=len(sortedMap)
		bufList = []string{}
		buf = Buffers{}
		result.PacketOfAllBuffers=0
		result.CurrentTimeCount++
	}
	if duration > float64(len(result.AccessPers))*perSec{
		result.AccessPers = append(result.AccessPers, 0)
	}
	return buf, bufList, result
}