package opCSV

import (
	"os"
	"strings"

	"bufio"
	"encoding/csv"
)

func TxtToCSV() {
	fp, err := os.Open("")
	failOnError(err)
	file, err := os.OpenFile("", os.O_WRONLY|os.O_CREATE, 0600)
	failOnError(err)
	err = file.Truncate(0)
	failOnError(err)

	writer := csv.NewWriter(file)

	defer fp.Close()
	defer file.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		list := strings.Split(scanner.Text(), " ")
		if len(list) == 8 {
			list = []string{strings.Join(list[1:6], " "), list[0]}
			writer.Write(list)
		}
	}
	writer.Flush()
}
