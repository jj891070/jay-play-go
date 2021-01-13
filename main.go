package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	// Open the file
	csvfile, err := os.Open("input.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(csvfile)
	r.Comma = '\t'
	inputDatas := []string{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for i := range record {
			inputDatas = append(inputDatas, strings.Split(record[i], ",")...)
		}
	}
	fmt.Println("inputDatas -> ", inputDatas)
	fmt.Println("len -> ", len(inputDatas))
}
