package Info

import (
	"fmt"
	"os"
	"strings"
	// "reflect"
	"bufio"
	"encoding/csv"
	// "log"
)



func txtToCSV() {
	fp, err := os.Open("../../pcap/sinet_sort.txt")
	failOnError(err)
	file, err := os.OpenFile("./test.csv", os.O_WRONLY|os.O_CREATE, 0600)
	failOnError(err)
	err = file.Truncate(0)
	failOnError(err)

	writer := csv.NewWriter(file)

	defer fp.Close()
	defer file.Close()

	fmt.Println(fp.Name())

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		list := strings.Split(scanner.Text(), " ")
		list = []string{list[1], list[2], list[3], list[4], list[5], list[0]}
		writer.Write(list)
	}
	writer.Flush()
}
