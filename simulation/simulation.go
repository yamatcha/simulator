package simulation

import (
	"encoding/csv"
	"io"
	"strconv"

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
