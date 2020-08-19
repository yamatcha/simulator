package simulation

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/yamatcha/simulator/buffer"
)

func GetRtt(csvReader *csv.Reader, buf buffer.Buffers, params buffer.Params) float64 {
	rttSum := 0.0
	rttCount := 0

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)

		_, ok := buf[fiveTuple]
		if !ok {
			newbuf := buffer.Buffer{params.CurrentTime, 1}
			buf[fiveTuple] = newbuf
			list := strings.Split(fiveTuple, " ")
			syn := strings.Join(append(append(list[2:4], list[0:2]...), list[4]), " ")
			b, ok := buf[syn]
			if ok {
				rttCount++
				rttSum += (params.CurrentTime - b.FirstTime)
			}
		}
	}
	return rttSum / float64(rttCount)
}

func GetWindow(csvReader *csv.Reader, buf buffer.Buffers, params buffer.Params) float64 {
	rttSum := 0
	rttCount := 0

	type flowWindow struct {
		sum   int
		count int
	}

	flowWindows := map[string]flowWindow{}

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)
		b, ok := buf[fiveTuple]
		if !ok {
			newbuf := buffer.Buffer{params.CurrentTime, 1}
			buf[fiveTuple] = newbuf
			list := strings.Split(fiveTuple, " ")
			syn := strings.Join(append(append(list[2:4], list[0:2]...), list[4]), " ")
			_, ok := buf[syn]
			// fmt.Println(buf)
			if ok {
				// rttCount++
				// rttSum+= buf[syn].Len
				f, ok := flowWindows[syn]
				if ok {
					flowWindows[syn] = flowWindow{f.sum + buf[syn].Len, f.count + 1}
				} else {
					flowWindows[syn] = flowWindow{buf[syn].Len, 1}
				}
				delete(buf, syn)
			}
		} else {
			newbuf := buffer.Buffer{b.FirstTime, b.Len + 1}
			buf[fiveTuple] = newbuf
		}
	}
	for _, window := range flowWindows {
		fmt.Println(float64(window.sum) / float64(window.count))
	}
	return float64(rttSum) / float64(rttCount)
}

func PreEval(csvReader *csv.Reader, buf buffer.Buffers, bufList []string, result buffer.ResultData, params buffer.Params) (buffer.Buffers, []string, buffer.ResultData) {
	for ; ; result.PacketNumAll++ {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)
		buf, bufList, result = buf.Append(bufList, params, fiveTuple, result)
	}
	result.EndFlag = true
	return buf, bufList, result
}

func Protocol(csvReader *csv.Reader, buf buffer.Buffers, bufList []string, result buffer.ResultData, params buffer.Params) (buffer.Buffers, []string, buffer.ResultData) {
	var protocolPort map[int]string = map[int]string{
		80: "HTTP", 443: "HTTPS", 25: "SMTP", 110: "POP3", 143: "IMAP4", 53: "DNS", 20: "FTP", 21: "FTP", 67: "DHCP", 68: "DHCP", 23: "TELNET", 179: "BGP",
	}
	tcp := make([]map[int]int, 450)
	udp := make([]map[int]int, 450)
	for num, _ := range protocolPort {
		tcp[num] = map[int]int{}
		udp[num] = map[int]int{}
	}
	for ; ; result.PacketNumAll++ {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)
		buf, bufList, result = buf.Append(bufList, params, fiveTuple, result)
	}
	for fiveTuple, buffer := range buf {
		ft := strings.Split(fiveTuple, " ")
		portA, _ := strconv.Atoi(strings.Split(ft[1], "(")[0])
		portB, _ := strconv.Atoi(strings.Split(ft[3], "(")[0])
		ports := []int{portA, portB}
		protocol := ft[4]
		if protocol == "TCP" {
			for _, port := range ports {
				_, ok := protocolPort[port]
				if ok {
					tcp[port][buffer.Len]++
				}
			}
		} else if protocol == "UDP" {
			for _, port := range ports {
				_, ok := protocolPort[port]
				if ok {
					udp[port][buffer.Len]++
				}
			}
		}

	}
	for num, name := range protocolPort {
		var keys []int
		for k := range tcp[num] {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		fmt.Print(name)
		for _, key := range keys {
			fmt.Printf(", (%d: %d)", tcp[num][key], tcp[num][key])
		}
		fmt.Println()
	}
	return buf, bufList, result
}
