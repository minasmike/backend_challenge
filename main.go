package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
)

type Record struct {
	Name         string `csv:"Department Name"`
	Numeric      string `csv:"Number of sales"`
	Abbreviation string `csv:"Abbreviation"`
}

type outputRecord struct {
	Name    string `csv:"Department Name"`
	Numeric string `csv:"Total number of sales"`
}

func main() {
	//fmt.Println("Measuring time in GO")
	start := time.Now()
	if len(os.Args) != 2 {
		fmt.Println("Usage:\nimporter <filename.csv>")
		return
	}
	if err := run(os.Args[1]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	timeElapsed := time.Since(start)
	fmt.Printf("The program took %s\n", timeElapsed)
}

func run(fileName string) error {
	fileHandle, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer fileHandle.Close()

	// make channel
	c := make(chan Record)

	go func() { // start parsing the CSV file
		err = gocsv.UnmarshalToChan(fileHandle, c) // <---- here it is
		if err != nil {
			log.Fatal(err)
		}
	}()
	// do something with the records
	months := make(map[string]int)
	newRecord := []*outputRecord{}
	for r := range c {
		//fmt.Println(r.Name, r.Numeric)
		// interesting code here
		numnum, err := strconv.Atoi(r.Numeric)
		if err != nil {
			// ... handle error
			panic(err)
		}
		months[r.Name] += numnum

	}

	for key, value := range months {
		strValue := strconv.Itoa(value)
		newRecord = append(newRecord, &outputRecord{Name: key, Numeric: strValue})
	}
	csvContent, err := gocsv.MarshalString(&newRecord)
	if err != nil {
		panic(err)
	}
	inputFileName := os.Args[1]
	outputFileName := "output" + inputFileName
	outputFile, err := os.OpenFile(outputFileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	//fmt.Println(outputFile)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = gocsv.MarshalFile(&newRecord, outputFile)
	if err != nil {
		panic(err)
	}
	fmt.Println(csvContent)
	return nil
}
