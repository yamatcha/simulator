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
	CurrentTimeCount   int
	BufMax             int
	PacketNumAll       int
	PacketOfAllBuffers int
	AccessPers         []int
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

func (buf Buffers) Append(bufList []string, params Params, fivetuple string, result ResultData) (Buffers, []string, ResultData) {
	_, ok := buf[fivetuple]
	result.PacketOfAllBuffers++
	if !ok {
		buf[fivetuple] = Buffer{params.CurrentTime, 1}
		bufList = append(bufList, fivetuple)
	} else {
		b := buf[fivetuple]
		buf[fivetuple] = Buffer{b.FirstTime, b.Len + 1}
		if len(bufList) > result.BufMax {
			result.BufMax = len(bufList)
		}
	}
	return buf, bufList, result
}

func batchProcessing(buf Buffers, bufList []string, params Params, result ResultData) (Buffers, []string, Params, ResultData) {
	sortedMap := List{}
	if params.Stupid == false {
		sortedMap = buf.getSortedMap(params.BufSize)
	} else {
		sortedMap = buf.getStupidMap(bufList, params)
	}
	reducing := sortedMap.getListSum(params)
	if result.PacketOfAllBuffers < reducing {
		fmt.Println(result.PacketOfAllBuffers, reducing)
	}

	accessCount := result.PacketOfAllBuffers - reducing
	result.AccessCount += accessCount
	result.AccessPers[len(result.AccessPers)-1] += accessCount
	bufList = []string{}
	buf = Buffers{}
	result.PacketOfAllBuffers = 0
	result.CurrentTimeCount++

	return buf, bufList, params, result
}

func (buf Buffers) CheckAck(fiveTuple string, bufList []string, params Params, result ResultData) {
	list := strings.Split(fiveTuple, " ")
	ack := strings.Join(append(append(list[2:4], list[0:2]...), list[4]), " ")
	_, ok := buf[ack]
	if ok {
		accessCount := int(float64(buf[ack].Len) / float64(params.EntrySize))
		result.AccessCount += accessCount
		result.AccessPers[len(result.AccessPers)-1] += accessCount
		result.PacketOfAllBuffers -= buf[ack].Len
		bufList = deleteList(bufList, ack)
		delete(buf, ack)
	}
}

func (buf Buffers) EndProcessing(bufList []string, params Params, result ResultData) {
	buf, bufList, params, result = batchProcessing(buf, bufList, params, result)
}

func (buf Buffers) CheckGlobalTime(bufList []string, params Params, result ResultData) (Buffers, []string, ResultData) {
	if params.BufSize > len(buf) {
		params.BufSize = len(buf)
	}
	if params.CurrentTime > float64(result.CurrentTimeCount+1)*params.TimeWidth || result.EndFlag == true {
		buf, bufList, params, result = batchProcessing(buf, bufList, params, result)
	}
	fmt.Println(params.CurrentTime, float64(len(result.AccessPers))*params.PerSec)
	if params.CurrentTime > float64(len(result.AccessPers))*params.PerSec {
		result.AccessPers = append(result.AccessPers, 0)
	}
	return buf, bufList, result
}

func (buf Buffers) CheckGlobalTimeWithUnlimitedBuffers(bufList []string, params Params, result ResultData) (Buffers, []string, ResultData) {
	if params.CurrentTime > float64(result.CurrentTimeCount+1)*params.TimeWidth || result.EndFlag == true {
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
