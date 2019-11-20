package Info

import (
	// "fmt"
	"os"
	"strings"
	// "reflect"
	"bufio"
	"encoding/csv"
	// "log"
)



func TxtToCSV() {
	fp, err := os.Open("./pcap/sinet/sinet_use.txt")
	failOnError(err)
	file, err := os.OpenFile("./sinet.csv", os.O_WRONLY|os.O_CREATE, 0600)
	failOnError(err)
	err = file.Truncate(0)
	failOnError(err)

	writer := csv.NewWriter(file)

	defer fp.Close()
	defer file.Close()

	// fmt.Println(fp.Name())

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		list := strings.Split(scanner.Text(), " ")
		if (len(list)==8){
			// fmt.Println(list)
			list = []string{strings.Join(list[1:6]," "), list[0]}
			// fmt.Println(list)
			writer.Write(list)
		}
	}
	writer.Flush()
}
