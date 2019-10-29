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
	"sort"
	"time"
)

type Buffer struct {
	FirstTime time.Time
	TimeList  []float64
	len       int
}

type Buffers map[string]Buffer

type ResultData struct {
	MaxPacketNum       int
	AccessCount        int
	CurrentTimeCount   int
	BufMax             int
	PacketNumAll       int
	PacketOfAllBuffers int
	AccessPers         []int
	EntryNums          [][]int
	EndFlag            bool
}

type Params struct {
	FirstTime   time.Time
	CurrentTime time.Time
	PerSec      float64
	BufSize     int
	TimeWidth   float64
	Stupid      bool
}

func GetDuration(first time.Time, now time.Time) float64 {
	return now.Sub(first).Seconds()
}

func (buf Buffers) CheckBufferTime(bufList []string, params Params, result ResultData) (Buffers, []string, ResultData) {
	i := 0
	for _, k := range bufList {
		if GetDuration(buf[k].FirstTime, params.CurrentTime) > params.TimeWidth || result.EndFlag == true {
			result.AccessCount++
			currentSec := int(GetDuration(params.FirstTime, buf[k].FirstTime.Add(time.Duration(params.TimeWidth*1000))) / params.PerSec)
			if currentSec >= len(result.AccessPers) {
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

func (buf Buffers) AppendBuffer(bufList []string, params Params, fivetuple string, result ResultData) (Buffers, []string, ResultData) {
	_, ok := buf[fivetuple]
	result.PacketOfAllBuffers++
	if !ok {
		new_timelist := []float64{0.0}
		newbuf := Buffer{params.CurrentTime, new_timelist, 1}
		buf[fivetuple] = newbuf
		bufList = append(bufList, fivetuple)
	} else {
		b := buf[fivetuple]
		newbuf := Buffer{b.FirstTime, (append((b.TimeList), GetDuration(b.FirstTime, params.CurrentTime))), b.len + 1}
		// *(b.TimeList) = append(*(b.TimeList), GetDuration(b.FirstTime, params.CurrentTime))
		// b.len++
		buf[fivetuple] = newbuf
		if len(bufList) > result.BufMax {
			result.BufMax = len(bufList)
		}
	}
	return buf, bufList, result
}

//using in Global time base func

func (buf Buffers) CheckGlobalTime(bufList []string, params Params, result ResultData) (Buffers, []string, ResultData) {
	duration := GetDuration(params.FirstTime, params.CurrentTime)
	if duration > float64(result.CurrentTimeCount+1)*params.TimeWidth {
		// fmt.Println(float64(result.CurrentTimeCount+1)*params.TimeWidth ,len(bufList))
		// fmt.Println(duration,len(bufList))
		result.AccessCount += len(bufList)
		result.AccessPers[len(result.AccessPers)-1] += len(bufList)
		bufList = []string{}
		buf = Buffers{}
		result.CurrentTimeCount++
	}
	if duration > float64(len(result.AccessPers))*params.PerSec {
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

func (buf Buffers) getSortedMap(bufSize int) List {
	sortedMap := List{}
	for k, v := range buf {
		element := Entry{k, v.len}
		sortedMap = append(sortedMap, element)
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

func (l List) getListSum() int {
	sum := 0
	for _, v := range l {
		sum += v.value
	}
	return sum
}

//using stupid simulation
func (buf Buffers) getStupidMap(bufList []string, bufSize int) List {
	sortedMap := List{}
	count := 0
	for _, k := range bufList {
		element := Entry{k, buf[k].len}
		sortedMap = append(sortedMap, element)
		count++
		if count == bufSize {
			break
		}
	}
	return sortedMap
}

//

func (buf Buffers) CheckGlobalTimeIdeal(bufList []string, params Params, result ResultData) (Buffers, []string, ResultData) {
	duration := GetDuration(params.FirstTime, params.CurrentTime)
	if params.BufSize > len(buf) {
		params.BufSize = len(buf)
	}
	sortedMap := List{}
	if duration > float64(result.CurrentTimeCount+1)*params.TimeWidth {
		if params.Stupid == false {
			sortedMap = buf.getSortedMap(params.BufSize)
		} else {
			sortedMap = buf.getStupidMap(bufList, params.BufSize)
		}
		// save the 10 biggest number of entry
		for i := 0; i < 10; i++ {
			// fmt.Println(sortedMap)
			// result.EntryNums[len(result.EntryNums)-1][i] += sortedMap[i].value
			result.EntryNums[i] = append(result.EntryNums[i],sortedMap[i].value)
		}
		bufferedPacket := sortedMap.getListSum()
		accessCount := result.PacketOfAllBuffers - bufferedPacket + params.BufSize
		// fmt.Println(accessCount,len(bufList))
		result.AccessCount += accessCount
		result.AccessPers[len(result.AccessPers)-1] += accessCount
		result.AccessCount += len(sortedMap)
		bufList = []string{}
		buf = Buffers{}
		result.PacketOfAllBuffers = 0
		result.CurrentTimeCount++
	}
	if duration > float64(len(result.AccessPers))*params.PerSec {
		result.AccessPers = append(result.AccessPers, 0)
		// result.EntryNums = append(result.EntryNums, make([]int, 10))
	}
	return buf, bufList, result
}
