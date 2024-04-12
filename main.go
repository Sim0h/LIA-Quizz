package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func problemPuller(fileName string) ([]problem, error) {
	if fObj, err := os.Open(fileName); err == nil {
		csvR := csv.NewReader(fObj)
		if cLines, err := csvR.ReadAll(); err == nil {
			return problemParser(cLines), nil
		} else {
			return nil, fmt.Errorf("Error in reading data in CSV"+"format from %s file; %s", fileName, err.Error())
		}
	} else {
		return nil, fmt.Errorf("Error in opening the %s file; %s", fileName, err.Error())
	}
}

func main() {
	fName := flag.String("f", "quiz.csv", "path of csv file")
	timer := flag.Int("t", 30, "timer of quiz")
	flag.Parse()

	problems, err := problemPuller(*fName)

	if err != nil {
		exit(fmt.Sprintf("something went wrong: %s", err.Error()))
	}

	var wg sync.WaitGroup
	var correctAns int
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)

	for i, p := range problems {
		var answer string
		wg.Add(1)

		fmt.Printf("Question %d: %s\n", i+1, p.q)

		go func() {
			defer wg.Done()
			reader := bufio.NewReader(os.Stdin)
			line, _ := reader.ReadString('\n')
			answer = line[:len(line)-1]
			ansC <- answer
		}()

		select {
		case <-tObj.C:
			fmt.Println("Times up!")
			close(ansC)
			return
		case iAns := <-ansC:
			iAns = strings.TrimSpace(iAns)
			if iAns == p.a {
				correctAns++
			}
		}
	}

	wg.Wait()
	close(ansC)

	fmt.Printf("\nYour result is: %d out of %d\n", correctAns, len(problems))
	fmt.Printf("Press enter to exit.")
	fmt.Scanln()
}

func problemParser(lines [][]string) []problem {
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{q: lines[i][0], a: lines[i][1]}
	}
	return r
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
