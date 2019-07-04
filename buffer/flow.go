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
	FiveTuple string
	Firstime  time.Time
	TimeList  *[]float64
}

type Buffers map[string]Buffer

func Check_buf_time(buf Buffers, buflist *[]string, nowtime time.Time, time_width float64, cnt *int, max *int, num_access *int) {

	for {
		if len(*buflist)==0{
			return
		}
		k := (*buflist)[0]
		if GetDuration(buf[k].Firstime, nowtime) > time_width {
			fmt.Println(len(*(buf[k].TimeList)))
			*cnt++
			*num_access++
			if *max < len(*(buf[k].TimeList)) {
				*max = len(*(buf[k].TimeList))
			}
			delete(buf, k)
			*buflist = append((*buflist)[:0],(*buflist)[1:]...)
			continue
		}
		return
	}
}

func Check_seconds(std_time time.Time,nowtime time.Time,time_width float64, num_access *int, access_pers *[]int,cs_count float64) float64{
	if GetDuration(std_time,nowtime) > cs_count*time_width{
		*access_pers = append(*access_pers,*num_access)
		*num_access = 0
		return (cs_count+1.0)
	}
	return cs_count
}

func Check_last(buf Buffers, buflist *[]string, cnt *int, max *int) {
	for {
		if len(*buflist)==0{
			return
		}
		k := (*buflist)[0]
			*cnt++
			if *max < len(*(buf[k].TimeList)) {
				*max = len(*(buf[k].TimeList))
			}
			delete(buf, k)
			*buflist = append((*buflist)[:0],(*buflist)[1:]...)
	}
}


func Append_buf(buf *Buffers, buflist *[]string, nowtime time.Time, fivetuple string) {
	_, ok := (*buf)[fivetuple]
	if !ok {
		new_timelist := []float64{0.0}
		newbuf := Buffer{fivetuple, nowtime, &new_timelist}
		(*buf)[fivetuple] = newbuf
		*buflist = append(*buflist, fivetuple)
	} else {
		b := (*buf)[fivetuple]
		*(b.TimeList) = append(*(b.TimeList), GetDuration(b.Firstime, nowtime))
	}
}
