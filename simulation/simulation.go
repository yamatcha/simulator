package simulation

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/yamatcha/simulator/buffer"
)

func GlobalTimeBase(csvReader *csv.Reader, buf buffer.Buffers, result buffer.ResultData, params buffer.Params) (buffer.Buffers, buffer.ResultData) {
	var count int
	for ; ; result.PacketNumAll++ {
		count++
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fiveTuple := line[0]
		params.CurrentTime, _ = strconv.ParseFloat(line[1], 64)

		buf, result = buf.CheckAck(fiveTuple, params, result)
		if params.IsLimited == false {
			buf, result = buf.CheckGlobalTimeWithUnlimitedBuffers(params, result)
		} else {
			buf, result = buf.CheckGlobalTime(params, result)
		}
		if !buffer.FiveTupleContains(fiveTuple, params) {
			buffer.Access(result, 1)
			return buf, result
		} else {
			buf, result = buf.Append(params, fiveTuple, result)
		}

		sum := 0
		for _, v := range buf {
			sum += v.Len
		}

		if sum != result.PacketOfAllBuffers {
			panic(fmt.Errorf("sum of buffers:%d result.PacketOfAllBuffers %d\n", sum, result.PacketOfAllBuffers))
		}
	}
	result.EndFlag = true
	if params.IsLimited == false {
		buf, result = buf.CheckGlobalTimeWithUnlimitedBuffers(params, result)
	} else {
		buf, result = buf.CheckGlobalTime(params, result)
	}
	// if buf != buffer.Buffers {
	// 	panic(fmt.Errorf("buf is not empty %v", buf))
	// }
	fmt.Println(buf)
	return buf, result
}
