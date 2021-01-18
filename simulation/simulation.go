package simulation

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/yamatcha/simulator/buffer"
)

func GlobalTimeBase(csvReader *csv.Reader, buf buffer.Buffers, bufOrderList []string, result buffer.ResultData, params buffer.Params, ideal bool) (buffer.Buffers, []string, buffer.ResultData) {
	for ; ; result.PacketNumAll++ {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)

		buf, bufOrderList, result = buf.CheckAck(fiveTuple, bufOrderList, params, result)
		if ideal == false {
			buf, bufOrderList, result = buf.CheckGlobalTimeWithUnlimitedBuffers(bufOrderList, params, result)
		} else {
			buf, bufOrderList, result = buf.CheckGlobalTime(bufOrderList, params, result)
		}
		buf, bufOrderList, result = buf.Append(bufOrderList, params, fiveTuple, result)

		sum := 0
		for _, v := range buf {
			sum += v.Len
		}

		if sum != result.PacketOfAllBuffers {
			panic(fmt.Errorf("sum of buffers:%d result.PacketOfAllBuffers %d\n", sum, result.PacketOfAllBuffers))
		}
	}
	result.EndFlag = true
	if ideal == false {
		buf, bufOrderList, result = buf.CheckGlobalTimeWithUnlimitedBuffers(bufOrderList, params, result)
	} else {
		buf, bufOrderList, result = buf.CheckGlobalTime(bufOrderList, params, result)
	}
	return buf, bufOrderList, result
}
