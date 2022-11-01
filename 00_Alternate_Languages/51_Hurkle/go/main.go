package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	GRIDSIZE = 10
	ATTEMPTS = 5
)

func printIntro() {
	fmt.Println("                                 HURKLE")
	fmt.Println("               CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY")
	fmt.Print("\n\n\n")
	fmt.Printf("A HURKLE IS HIDING ON A %d BY %d GRID. HOMEBASE\n", GRIDSIZE, GRIDSIZE)
	fmt.Println("ON THE GRID IS POINT 0,0 IN THE SOUTHWEST CORNER,")
	fmt.Println("AND ANY POINT ON THE GRID IS DESIGNATED BY A")
	fmt.Println("PAIR OF WHOLE NUMBERS SEPERATED BY A COMMA. THE FIRST")
	fmt.Println("NUMBER IS THE HORIZONTAL POSITION AND THE SECOND NUMBER")
	fmt.Println("IS THE VERTICAL POSITION. YOU MUST TRY TO")
	fmt.Printf("GUESS THE HURKLE'S GRIDPOINT. YOU GET %d TRIES.\n", ATTEMPTS)
	fmt.Println("AFTER EACH TRY, I WILL TELL YOU THE APPROXIMATE")
	fmt.Println("DIRECTION TO GO TO LOOK FOR THE HURKLE.")
	fmt.Println()
}

func parseCoords(n int) (int, int) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("GUESS # %d\n", n)
		scanner.Scan()

		vals := strings.Split(scanner.Text(), ",")
		if len(vals) != 2 {
			fmt.Println("INVALID INPUT, PLEASE ENTER VALUES LIKE x,y")
			continue
		}

		// parse the first coord
		x, err := strconv.Atoi(strings.TrimSpace(vals[0]))
		if err != nil {
			fmt.Printf("INVALID INPUT '%s', PLEASE ENTER WHOLE NUMBERS ONLY\n", vals[0])
			continue
		}

		// parse the second coord
		y, err := strconv.Atoi(strings.TrimSpace(vals[1]))
		if err != nil {
			fmt.Printf("INVALID INPUT '%s', PLEASE ENTER WHOLE NUMBERS ONLY\n", vals[1])
			continue
		}

		return x, y
	}
}

func printHint(A, B, X, Y int) {
	hint := strings.Builder{}
	hint.WriteString("\nGO ")

	if Y < B {
		hint.WriteString("NORTH ")
	} else if Y > B {
		hint.WriteString("SOUTH ")
	}
	if X < A {
		hint.WriteString("EAST ")
	} else if X > A {
		hint.WriteString("WEST ")
	}

	fmt.Println(hint.String())
}

func main() {
	rand.Seed(time.Now().UnixNano())

	printIntro()

	for {
		// position the HURKLE
		hurkleX := int(GRIDSIZE * rand.Float64())
		hurkleY := int(GRIDSIZE * rand.Float64())

		// read user guesses
		for guess := 0; guess < ATTEMPTS; guess++ {
			guessX, guessY := parseCoords(guess)

			if math.Abs(float64(hurkleX-guessX)) == 0 && math.Abs(float64(hurkleY-guessY)) == 0 {
				fmt.Printf("\n YOU FOUND HIM IN %d GUESSES!\n", guess+1)
				break
			} else {
				printHint(hurkleX, hurkleY, guessX, guessY)
			}
		}

		fmt.Println("\n\nLET'S PLAY AGAIN, HURKLE IS HIDING.")
	}
}
