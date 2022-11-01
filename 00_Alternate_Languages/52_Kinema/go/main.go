package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	GRAVITY           = 10.0
	EXPECTED_ACCURACY = 0.15
)

func printIntro() {
	fmt.Println("                                  KINEMA")
	fmt.Println("               CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY")
	fmt.Print("\n\n\n")
}

func doQuiz() {
	numCorrect := 0
	v0 := float64(rand.Intn(35) + 5)

	fmt.Printf("A BALL IS THROWN UPWARDS AT %0.2f METRES PER SECOND.\n\n", v0)

	// Q1
	A1 := (v0 * v0) / (2 * GRAVITY)
	numCorrect += askPlayer("HOW HIGH WILL IT GO (IN METRES)?", A1)

	// Q2
	A2 := 2.0 * v0 / GRAVITY
	numCorrect += askPlayer("HOW LONG UNTIL IT RETURNS (IN SECONDS)", A2)

	// Q3
	t := 1 + rand.Intn(2*int(v0))/GRAVITY
	A3 := v0 - GRAVITY*float64(t)
	numCorrect += askPlayer(fmt.Sprintf("WHAT WILL ITS VELOCITY BE AFTER %d SECONDS?", t), A3)

	fmt.Printf("\n%d right out of 3.\n", numCorrect)
	if numCorrect >= 2 {
		fmt.Println("  NOT BAD.")
	}
}

func askPlayer(prompt string, answer float64) (score int) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(prompt)

	for {
		scanner.Scan()
		response, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			fmt.Println("INVALID INPUT, TRY AGAIN")
			continue
		}

		if math.Abs((response-answer)/answer) < EXPECTED_ACCURACY {
			fmt.Println("CLOSE ENOUGH")
			score = 1
		} else {
			fmt.Println("NOT EVEN CLOSE...")
			score = 0
		}

		fmt.Printf("CORRECT ANSWER IS %0.2f\n\n", answer)
		return
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	printIntro()
	doQuiz()
}
