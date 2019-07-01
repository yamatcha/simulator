package buffer

import (
	// "fmt"
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
	len int
}

func Check_buf_time(buf []*Buffer,nowtime time.Time, time_width float64, cnt int ) int{
	for i,b := range buf{
		duration:= nowtime.Sub(b.Firstime).Seconds()
		if duration > time_width{
			cnt++
			buf = append(buf[:i],buf[(i+1):]...)
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

func Append_buf(packet gopacket.Packet, buf []*Buffer){
	fivetuple := GetFiveTuple(packet)
	bnum := Search_buf(buf,fivetuple)
	ptime := GetTime(packet)
	if bnum!=-1{
		newbuf:=Buffer{fivetuple,ptime,1}
		buf = append(buf,&newbuf)
	}else{
		buf[bnum].len++
	}
}