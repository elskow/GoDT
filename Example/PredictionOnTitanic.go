package main

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"os"
)

func main() {
	csvFile, err := os.Open("Data/Titanic/train.csv")

	if err != nil {
		panic(err)
	}

	train := dataframe.ReadCSV(csvFile)

	fmt.Println(train)
}
