package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
)

const totalQuestions = 5

type Question struct {
	question string
	answer   string
}

func main() {

	csvFilename, timeLimit := readArguments()
	f, err := openFile(csvfile)
	if err != nil {
		return
	}
	questions, err := readProblems(f)

	if err != nil {
		fmt.Println(err.error())
		return
	}
}

func readArguments(string, int) {

	csvfile := flag.String(
		"csvfile",
		"problems.csv",
		"The filename of the CSV file containing quiz Questions")
	timeLimit := flag.Int(
		"limit",
		20,
		"Time limit for questions")
	flag.Parse()
	return *csvfile, timeLimit

}

func readProblems(f io.Reader) ([]Question, error) {
	Problems, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	// kolla om det finns frågor i problems.csv
	numOfProblems := len(Problems)
	if numOfProblems == 0 {
		return nil, fmt.Errorf("No Questions in problems.csv")
	}

}
