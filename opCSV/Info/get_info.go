package Info

import (
	// "fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func GetFiveTuple(packet gopacket.Packet) string {
	// extract the factor of packet
	ip, _ := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	var info string
	// if ip address exist
	if ip != nil {
		//get value of protocol, sourceIP,destinationIP
		protocol := ip.Protocol.String()
		srcip := ip.SrcIP.String()
		dstip := ip.DstIP.String()

		// port is different between udp and tcp
		if ip.Protocol.String() == "UDP" {
			udp, _ := packet.Layer(layers.LayerTypeUDP).(*layers.UDP)
			if udp != nil {
				srcport := udp.SrcPort.String()
				dstport := udp.DstPort.String()
				info = srcip + " " + srcport + " " + dstip + " " + dstport + " " + protocol
			}
		}
		if ip.Protocol.String() == "TCP" {
			tcp, _ := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)
			if tcp != nil {
				srcport := tcp.SrcPort.String()
				dstport := tcp.DstPort.String()
				info = srcip + " " + srcport + " " + dstip + " " + dstport + " " + protocol
			}
		}
	}
	return info
}

func GetTime(packet gopacket.Packet) time.Time {
	meta := packet.Metadata()
	return meta.Timestamp
}

func GetDuration(first time.Time, now time.Time) float64 {
	return now.Sub(first).Seconds()
}
