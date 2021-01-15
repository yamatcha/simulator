package Info

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
)

type FiveTuple struct {
	srcIp, srcPort, dstIp, dstPort, protocol string
}

func GetFiveTuple(packet gopacket.Packet) FiveTuple {
	dstIp := packet.NetworkLayer().NetworkFlow().Dst().String()
	srcIp := packet.NetworkLayer().NetworkFlow().Src().String()
	dstPort := packet.TransportLayer().TransportFlow().Dst().String()
	srcPort := packet.TransportLayer().TransportFlow().Src().String()
	protocol := packet.TransportLayer().LayerType().String()
	if dstIp == "" || srcIp == "" || dstPort == "" || srcPort == "" || protocol == "" {
		return FiveTuple{}
	}
	return FiveTuple{srcIp, srcPort, dstIp, dstPort, protocol}
}

func (f FiveTuple) String() string {
	return fmt.Sprintf("%s %s %s %s %s", f.srcIp, f.srcPort, f.dstIp, f.dstPort, f.protocol)
}

func (f FiveTuple) ForCacheSimulator() []string {
	return []string{f.srcIp, f.srcPort, f.dstIp, f.dstPort, f.protocol}
}

func GetTime(packet gopacket.Packet) time.Time {
	meta := packet.Metadata()
	return meta.Timestamp
}

func GetLen(packet gopacket.Packet) int {
	meta := packet.Metadata()
	return meta.Length
}

func GetDuration(first time.Time, now time.Time) float64 {
	return now.Sub(first).Seconds()
}
