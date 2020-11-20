package buffer

import (
	"container/heap"
	"fmt"
	"strings"
)

type Buffer struct {
	FirstTime float64
	Len       int
}

type Buffers map[string]Buffer

type ResultData struct {
	AccessCount        int
	NextAccessTime     int
	BufMax             int
	PacketNumAll       int
	PacketOfAllBuffers int
	AccessPerSecList   []int
	EndFlag            bool
}

type Params struct {
	CurrentTime  float64
	BufSize      int
	EntrySize    int
	TimeWidth    float64
	IsStupid     bool
	Protocol     string
	SelectedPort []string
	IsLimited    bool
	PerSec       float64
}

func (buf Buffers) Append(params Params, fivetuple string, result ResultData) (Buffers, ResultData) {
	_, ok := buf[fivetuple]
	result.PacketOfAllBuffers++
	if !ok {
		buf[fivetuple] = Buffer{params.CurrentTime, 1}
	} else {
		b := buf[fivetuple]
		buf[fivetuple] = Buffer{b.FirstTime, b.Len + 1}
		if len(buf) > result.BufMax {
			result.BufMax = len(buf)
		}
	}
	return buf, result
}

func Access(result ResultData, accessCount int) ResultData {
	result.AccessCount += accessCount
	result.AccessPerSecList[len(result.AccessPerSecList)-1] += accessCount
	return result
}

func batchProcessing(buf Buffers, params Params, result ResultData) (Buffers, ResultData) {
	var reducing int
	if params.IsStupid == false {
		reducing = buf.getSortedMap(params.BufSize)
	} else {
		reducing = buf.getStupidMap(params.BufSize)
	}
	if result.PacketOfAllBuffers < reducing {
		panic(fmt.Errorf("error: reducing(%d) is more than packet of all buffers(%d)%d", reducing, result.PacketOfAllBuffers))
	}

	accessCount := result.PacketOfAllBuffers - reducing
	result = Access(result, accessCount)
	buf = Buffers{}
	result.PacketOfAllBuffers = 0
	result.NextAccessTime = int(params.CurrentTime/params.TimeWidth) + 1

	return buf, result
}

func (buf Buffers) CheckAck(fiveTuple string, params Params, result ResultData) (Buffers, ResultData) {
	list := strings.Split(fiveTuple, " ")
	ack := strings.Join(append(append(list[2:4], list[0:2]...), list[4]), " ")
	_, ok := buf[ack]
	if ok {
		accessCount := int(float64(buf[ack].Len) / float64(params.EntrySize))
		result = Access(result, accessCount)
		result.PacketOfAllBuffers -= buf[ack].Len
		delete(buf, ack)
	}
	return buf, result
}

func (buf Buffers) CheckGlobalTime(params Params, result ResultData) (Buffers, ResultData) {
	if params.BufSize > len(buf) {
		params.BufSize = len(buf)
	}
	if params.CurrentTime > float64(result.NextAccessTime)*params.TimeWidth || result.EndFlag == true {
		buf, result = batchProcessing(buf, params, result)
	}
	if params.CurrentTime > float64(len(result.AccessPerSecList))*params.PerSec {
		result.AccessPerSecList = append(result.AccessPerSecList, 0)
	}
	return buf, result
}

func (buf Buffers) CheckGlobalTimeWithUnlimitedBuffers(params Params, result ResultData) (Buffers, ResultData) {
	if params.CurrentTime > float64(result.NextAccessTime)*params.TimeWidth || result.EndFlag == true {
		Access(result, len(buf))
		buf = Buffers{}
		result.NextAccessTime = int(params.CurrentTime)*100 + 1
	}
	if params.CurrentTime > float64(len(result.AccessPerSecList))*params.PerSec {
		result.AccessPerSecList = append(result.AccessPerSecList, 0)
	}
	return buf, result
}

//for get sorted map
type Entry struct {
	key   int
	value int
}
type EntryHeap []Entry

func (h EntryHeap) Len() int           { return len(h) }
func (h EntryHeap) Less(i, j int) bool { return h[i].key < h[j].key }
func (h EntryHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *EntryHeap) Push(x interface{}) {
	*h = append(*h, x.(Entry))
}
func (h *EntryHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

//using stupid simulation
func (buf Buffers) getStupidMap(bufSize int) int {
	sortedMap := &EntryHeap{}
	heap.Init(sortedMap)
	sum := 0
	for _, v := range buf {
		element := Entry{int(v.FirstTime * 1000000), v.Len}
		heap.Push(sortedMap, element)
		sum += v.Len
		if sortedMap.Len() > bufSize {
			poppedEntry := heap.Pop(sortedMap).(Entry)
			sum -= poppedEntry.value
		}
	}
	return sum
}

func (buf Buffers) getSortedMap(bufSize int) int {
	sortedMap := &EntryHeap{}
	heap.Init(sortedMap)
	sum := 0
	for _, v := range buf {
		element := Entry{v.Len, v.Len}
		heap.Push(sortedMap, element)
		sum += v.Len
		if sortedMap.Len() > bufSize {
			poppedEntry := heap.Pop(sortedMap).(Entry)
			sum -= poppedEntry.value
		}
	}
	return sum
}

func FiveTupleContains(fiveTuple string, params Params) bool {
	if params.SelectedPort[0] == "" {
		return true
	}
	List := strings.Split(fiveTuple, " ")
	if params.Protocol != List[4] {
		return false
	}
	for _, v := range params.SelectedPort {
		if v == strings.Split(List[1], "(")[0] || v == strings.Split(List[3], "(")[0] {
			return true
		}
	}
	return false
}
