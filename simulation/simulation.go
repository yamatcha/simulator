package simulation

import (
	"encoding/csv"
	"io"
        "fmt"
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
		if ideal == false {
			buf, bufOrderList, result = buf.CheckGlobalTimeWithUnlimitedBuffers(bufOrderList, params, result)
		} else {
			buf, bufOrderList, result = buf.CheckGlobalTime(bufOrderList, params, result)
		}
                prev := buf
		buf, bufOrderList, result = buf.Append(bufOrderList, params, fiveTuple, result)

                if sum!=result.PacketOfAllBuffers{
                        panic(fmt.Errorf("%d %d %v %v\n",sum,result.PacketOfAllBuffers,buf,prev))
                }
		buf, bufOrderList,result=buf.CheckAck(fiveTuple, bufOrderList, params, result)
	}
	result.EndFlag = true
	if ideal == false {
		buf, bufOrderList, result = buf.CheckGlobalTimeWithUnlimitedBuffers(bufOrderList, params, result)
	} else {
		buf, bufOrderList, result = buf.CheckGlobalTime(bufOrderList, params, result)
	}
	return buf, bufOrderList, result
}
