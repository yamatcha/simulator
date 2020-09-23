package simulation

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"github.com/yamatcha/simulator/buffer"
)

func GlobalTimeBase(csvReader *csv.Reader, buf buffer.Buffers, bufList []string, result buffer.ResultData, params buffer.Params, ideal bool) (buffer.Buffers, []string, buffer.ResultData) {
	for ; ; result.PacketNumAll++ {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		if !FiveTupleContains(fiveTuple, params) {
			continue
		}
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)
		if ideal == false {
			buf, bufList, result = buf.CheckGlobalTimeWithUnlimitedBuffers(bufList, params, result)
		} else {
			buf, bufList, result = buf.CheckGlobalTime(bufList, params, result)
		}
		buf, bufList, result = buf.Append(bufList, params, fiveTuple, result)
		buf.CheckAck(fiveTuple, bufList, params, result)
	}
	result.EndFlag = true
	if ideal == false {
		buf, bufList, result = buf.CheckGlobalTimeWithUnlimitedBuffers(bufList, params, result)
	} else {
		buf, bufList, result = buf.CheckGlobalTime(bufList, params, result)
	}
	return buf, bufList, result
}

func FiveTupleContains(fiveTuple string, params buffer.Params) bool {
	if len(params.SelectedPort) == 0 {
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
