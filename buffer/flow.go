package buffer

import (
	"fmt"
	"github.com/google/gopacket"
	// "github.com/google/gopacket/layers"
	// "github.com/google/gopacket/pcap"
	// "net"
	// "io"
	// "log"
	// "reflect"
	"time"
)

type Buffer struct{
	FiveTuple string
	Firstime time.Time
	TimeList []float64
}

func Check_buf_time(buf []*Buffer,nowtime time.Time, time_width float64, cnt int ) int{
	for i,b := range buf{
		if GetDuration(b.Firstime,nowtime) > time_width{
			cnt++
			fmt.Println("[",len(buf),"]")
			fmt.Println(buf[i])
			fmt.Println(i)
			// if len(buf)==1 {
			// 	buf = []*Buffer{}
			// 	break
			// }
			buf = append(buf[:i],buf[i+1:]...)
		}
	}
	return len(buf)
}

func Search_buf(buf []*Buffer,fivetuple string)int{
	for i,b := range buf{
		if b.FiveTuple==fivetuple{
			return i
		}
	}
	return -1
}

func Append_buf(packet *gopacket.Packet, buf *[]*Buffer){
	fivetuple := GetFiveTuple(*packet)
	bnum := Search_buf(*buf,fivetuple)
	// fmt.Println(bnum)
	ptime := GetTime(*packet)
	if bnum==-1{
		new_timelist := []float64{0.0}
		newbuf:=Buffer{fivetuple,ptime,new_timelist}
		*buf = append(*buf,&newbuf)
		// for i,v := range *buf{
		// 	fmt.Println(i,*v)
		// }
		// fmt.Println(buf)
	}else{
		b := (*buf)[bnum]
		b.TimeList = append(b.TimeList,GetDuration(b.Firstime,ptime))
	}
}