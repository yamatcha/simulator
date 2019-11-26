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
	"math"
	// "time"
)

type Buffer struct {
	FirstTime float64
	// TimeList  []float64
	Len       int
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
	CurrentTime float64
	PerSec      float64
	BufSize     int
	EntrySize     int
	TimeWidth   float64
	Stupid      bool
}


func remove(strings []string, search string) []string {
    result := []string{}
    for _, v := range strings {
        if v != search {
            result = append(result, v)
        }
    }
    return result
}

func (buf Buffers) CheckBufferTime(bufList []string, params Params, result ResultData) (Buffers, []string, ResultData) {
	i := 0
	for _, k := range bufList {
		if (params.CurrentTime - buf[k].FirstTime) > params.TimeWidth || result.EndFlag == true {
			result.AccessCount++
			if int(params.CurrentTime) >= len(result.AccessPers) {
				result.AccessPers = append(result.AccessPers, 0)
			}
			result.AccessPers[len(result.AccessPers)-1]++
			if result.MaxPacketNum < buf[k].Len {
				result.MaxPacketNum = buf[k].Len
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
		// new_timelist := []float64{0.0}
		newbuf := Buffer{params.CurrentTime,  1}
		buf[fivetuple] = newbuf
		bufList = append(bufList, fivetuple)
	} else {
		b := buf[fivetuple]
		newbuf := Buffer{b.FirstTime,  b.Len + 1}
		// *(b.TimeList) = append(*(b.TimeList), GetDuration(b.FirstTime, params.CurrentTime))
		// b.Len++
		buf[fivetuple] = newbuf
		if len(bufList) > result.BufMax {
			result.BufMax = len(bufList)
		}
	}
	// if buf[fivetuple].Len >= params.EntrySize{// && params.EntrySize!=0{
	// 	result.AccessCount++
	// 	result.AccessPers[len(result.AccessPers)-1]++
	// 	result.PacketOfAllBuffers -= buf[fivetuple].Len
	// 	delete(buf,fivetuple)
	// 	bufList = remove(bufList,fivetuple)
	// 	if len(buf)!=len(bufList){
	// 		fmt.Println("error")
	// 	}
	// }
	return buf, bufList, result
}

//using in Global time base func

func (buf Buffers) CheckGlobalTime(bufList []string, params Params, result ResultData) (Buffers, []string, ResultData) {
	if params.CurrentTime > float64(result.CurrentTimeCount+1)*params.TimeWidth|| result.EndFlag == true  {
		// fmt.Println(float64(result.CurrentTimeCount+1)*params.TimeWidth ,len(bufList))
		// fmt.Println(duration,len(bufList))
		result.AccessCount += len(bufList)
		result.AccessPers[len(result.AccessPers)-1] += len(bufList)
		bufList = []string{}
		buf = Buffers{}
		result.CurrentTimeCount++
	}
	if params.CurrentTime > float64(len(result.AccessPers))*params.PerSec {
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
		element := Entry{k, v.Len}
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

func (l List) getListSum(params Params) int {
	sum := 0
	for _, v := range l {
		sum += v.value
	}
	// fmt.Println(float64(sum)/(float64(params.EntrySize-1)/float64(params.EntrySize)))
	// fmt.Println(params.EntrySize)
	return int(math.Ceil(float64(sum)*(float64(params.EntrySize-1)/float64(params.EntrySize))))
}

//using stupid simulation
func (buf Buffers) getStupidMap(bufList []string, params Params) List {
	sortedMap := List{}
	count := 0
	for _, k := range bufList {
		element := Entry{k, buf[k].Len}
		sortedMap = append(sortedMap, element)
		count++
		if count == params.BufSize {
			break
		}
	}
	return sortedMap
}

//

func (buf Buffers) CheckGlobalTimeIdeal(bufList []string, params Params, result ResultData) (Buffers, []string, ResultData) {
	if params.BufSize > len(buf) {
		params.BufSize = len(buf)
	}
	sortedMap := List{}
	if params.CurrentTime > float64(result.CurrentTimeCount+1)*params.TimeWidth || result.EndFlag == true {
		// fmt.Println(params.BufSize)
		if params.Stupid == false {
			sortedMap = buf.getSortedMap(params.BufSize)
		} else {
			sortedMap = buf.getStupidMap(bufList, params)
		}
		// save the 10 biggest number of entry
		// for i := 0; i < 10 && params.BufSize>=10; i++ {
		// 	result.EntryNums[i] = append(result.EntryNums[i],sortedMap[i].value)
		// }
		reducing := sortedMap.getListSum(params)
		if result.PacketOfAllBuffers< reducing{
			fmt.Println(result.PacketOfAllBuffers, reducing)
		}

		accessCount := result.PacketOfAllBuffers - reducing
		// fmt.Println(result.PacketOfAllBuffers == bufferedPacket , params.BufSize)
		// fmt.Println(accessCount,len(bufList))
		result.AccessCount += accessCount
		result.AccessPers[len(result.AccessPers)-1] += accessCount
		// result.AccessCount += len(sortedMap)
		bufList = []string{}
		buf = Buffers{}
		result.PacketOfAllBuffers = 0
		result.CurrentTimeCount++
	}
	if params.CurrentTime > float64(len(result.AccessPers))*params.PerSec {
		// fmt.Println(result.AccessPers[len(result.AccessPers)-1])
		result.AccessPers = append(result.AccessPers, 0)
	}
	return buf, bufList, result
}
