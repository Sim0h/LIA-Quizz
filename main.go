package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	//1. input the name of the file
	fName := flag.String("f", "problems.csv", "path of the csv file")
	//2. skapa en timer
	timer := flag.Int("t", 30, "timer of quizz")
	flag.Parse()
	//3. Hämta frågor från CSV fil
	problems, err := grabQuestions(*fName)
	//4. felhantering
	if err != nil {
		closeApp(fmt.Sprintf("Something went wrong: %s", err.Error()))
	}

	correctAns := 0
	//6. initialisera timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)
	//7. loopa igenom frågorna och skriv ut dem. spara svar.

	for i, p := range problems {
		var answer string
		fmt.Printf("Question %d: %s\n", i+1, p.q)

		go func() {
			fmt.Scanf("%s", &answer)
			ansC <- answer

		}()

		select {
		case <-tObj.C:
			fmt.Println()
			break

		case iAns := <-ansC:
			if iAns == p.a {
				correctAns++
			}
			if i == len(problems)-1 {
				close(ansC)
			}
		}
	}

	fmt.Printf("Your result is %d out of %d\n", correctAns, len(problems))
	fmt.Printf("Press enter to Exit")
	<-ansC

}

func grabQuestions(fileName string) ([]problem, error) {
	if fObj, err := os.Open(fileName); err == nil {
		csvR := csv.NewReader(fObj)
		if cLines, err := csvR.ReadAll(); err == nil {
			return parsQuestions(cLines), nil
		} else {
			return nil, fmt.Errorf("Error in reading data in CSV"+"format from %s file; %s", fileName, err.Error())
		}

	} else {
		return nil, fmt.Errorf("error in opening %s file; %s", fileName, err.Error())
	}

}

type problem struct {
	//. struct som gör att att golang kan läsa CSV fil.
	q string
	a string
}

func parsQuestions(lines [][]string) []problem {
	// kollar csv text och parse för readability
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{q: lines[i][0], a: lines[i][1]}
	}

	return r
}

func closeApp(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
