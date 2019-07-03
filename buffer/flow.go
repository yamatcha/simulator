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
	FiveTuple string
	Firstime  time.Time
	TimeList  *[]float64
}

type Buffers map[string]Buffer

func Check_buf_time(buf Buffers, buflist *[]string, nowtime time.Time, time_width float64, cnt *int, max *int) {

	for {
		if len(*buflist)==0{
			return
		}
		k := (*buflist)[0]
		if GetDuration(buf[k].Firstime, nowtime) > time_width {
			// fmt.Println("duration",GetDuration(buf[k].Firstime, nowtime))
			*cnt++
			// fmt.Println(len(buf))
			// for _,v := range buf{
			// 	fmt.Println(v)
			// }
			// fmt.Println(*buflist)
			// fmt.Println(k,buf[k])
			// fmt.Println("before:",*buf[k].TimeList)
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
