package buffer

import (
	"fmt"
	"sort"
	"strings"
)

type Buffer struct {
	FirstTime float64
	Len       int
}

type Buffers map[string]Buffer

type ResultData struct {
	MaxPacketNum       int
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
	PerSec       float64
	BufSize      int
	EntrySize    int
	TimeWidth    float64
	Stupid       bool
	Protocol     string
	SelectedPort []string
}

func (buf Buffers) Append(bufOrderList []string, params Params, fivetuple string, result ResultData) (Buffers, []string, ResultData) {
	_, ok := buf[fivetuple]
	result.PacketOfAllBuffers++
	if !ok {
		buf[fivetuple] = Buffer{params.CurrentTime, 1}
		bufOrderList = append(bufOrderList, fivetuple)
	} else {
		b := buf[fivetuple]
		buf[fivetuple] = Buffer{b.FirstTime, b.Len + 1}
		if len(bufOrderList) > result.BufMax {
			result.BufMax = len(bufOrderList)
		}
	}
	return buf, bufOrderList, result
}

func batchProcessing(buf Buffers, bufOrderList []string, params Params, result ResultData) (Buffers, []string, Params, ResultData) {
	sortedMap := List{}
	if params.Stupid == false {
		sortedMap = buf.getSortedMap(params.BufSize)
	} else {
		sortedMap = buf.getStupidMap(bufOrderList, params)
	}
	reducing := sortedMap.getListSum(params)
	if result.PacketOfAllBuffers < reducing {
		panic(fmt.Errorf("error: reducing(%d) is more than packet of all buffers(%d)%d", reducing, result.PacketOfAllBuffers, len(sortedMap)))
	}

	accessCount := result.PacketOfAllBuffers - reducing
	result.AccessCount += accessCount
	result.AccessPerSecList[len(result.AccessPerSecList)-1] += accessCount
	bufOrderList = []string{}
	buf = Buffers{}
	result.PacketOfAllBuffers = 0
	result.NextAccessTime = int(params.CurrentTime/params.TimeWidth) + 1

	return buf, bufOrderList, params, result
}

func (buf Buffers) CheckAck(fiveTuple string, bufOrderList []string, params Params, result ResultData) (Buffers, []string, ResultData) {
	list := strings.Split(fiveTuple, " ")
	ack := strings.Join(append(append(list[2:4], list[0:2]...), list[4]), " ")
	_, ok := buf[ack]
	if ok {
		accessCount := int(float64(buf[ack].Len) / float64(params.EntrySize))
		result.AccessCount += accessCount
		result.AccessPerSecList[len(result.AccessPerSecList)-1] += accessCount
		result.PacketOfAllBuffers -= buf[ack].Len
		bufOrderList = deleteList(bufOrderList, ack)
		delete(buf, ack)
	}
	return buf, bufOrderList, result
}

func (buf Buffers) EndProcessing(bufOrderList []string, params Params, result ResultData) {
	buf, bufOrderList, params, result = batchProcessing(buf, bufOrderList, params, result)
}

func (buf Buffers) CheckGlobalTime(bufOrderList []string, params Params, result ResultData) (Buffers, []string, ResultData) {
	if params.BufSize > len(buf) {
		params.BufSize = len(buf)
	}
	if params.CurrentTime/params.TimeWidth > float64(result.NextAccessTime) || result.EndFlag == true {
		buf, bufOrderList, params, result = batchProcessing(buf, bufOrderList, params, result)
	}
	if params.CurrentTime > float64(len(result.AccessPerSecList))*params.PerSec {
		result.AccessPerSecList = append(result.AccessPerSecList, 0)
	}
	return buf, bufOrderList, result
}

func (buf Buffers) CheckGlobalTimeWithUnlimitedBuffers(bufOrderList []string, params Params, result ResultData) (Buffers, []string, ResultData) {
	if params.CurrentTime > float64(result.NextAccessTime)*params.TimeWidth || result.EndFlag == true {
		result.AccessCount += len(bufOrderList)
		result.AccessPerSecList[len(result.AccessPerSecList)-1] += len(bufOrderList)
		bufOrderList = []string{}
		buf = Buffers{}
		result.NextAccessTime = int(params.CurrentTime)*100 + 1
	}
	if params.CurrentTime > float64(len(result.AccessPerSecList))*params.PerSec {
		result.AccessPerSecList = append(result.AccessPerSecList, 0)
	}
	return buf, bufOrderList, result
}

//for get sorted map

type Entry struct {
	name  string
	value int
}
type List []Entry

//using stupid simulation
func (buf Buffers) getStupidMap(bufOrderList []string, params Params) List {
	sortedMap := List{}
	count := 0
	for _, k := range bufOrderList {
		element := Entry{k, buf[k].Len}
		if FiveTupleContains(k, params) {
			sortedMap = append(sortedMap, element)
			count++
		}
		if count == params.BufSize {
			break
		}
	}
	return sortedMap
}

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
		sum += int(float64(v.value) * (float64(params.EntrySize-1) / float64(params.EntrySize)))
	}
	return sum
}

func deleteList(l []string, s string) []string {
	for i, v := range l {
		if v == s {
			l := append(l[:i], l[i+1:]...)
			n := make([]string, len(s))
			copy(n, l)
			return l
		}
	}
	return l
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
