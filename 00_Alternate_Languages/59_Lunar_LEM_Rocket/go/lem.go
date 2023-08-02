package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	PAGE_WIDTH = 64

	B  = 750.0
	M0 = 7.45
	R0 = 926.0
	V0 = 1.29
)

type Simulation struct {
	A  float64
	A1 float64
	A3 float64
	C  float64
	F  float64
	F1 float64
	G3 float64
	G5 float64
	H0 float64
	M  float64
	M1 float64
	M2 float64
	N  float64
	P  float64
	R  float64
	R1 float64
	R3 float64
	S  float64
	T  int
	T1 float64
	Z  float64

	ShortUnit string
	LongUnit  string
}

func NewSim() Simulation {
	s := Simulation{}

	s.A = -3.425
	s.A1 = 8.84361e-04
	s.F1 = 5.25
	s.H0 = 60
	s.M = 7.45
	s.M1 = M0
	s.N = 1 // initialised twice in lem.bas. bug?
	s.R = R0 + s.H0

	return s
}

func (s *Simulation) setUnits(u int) {
	if u == 0 {
		s.G3 = 0.592
		s.G5 = 6080
		s.ShortUnit = "FEET"
		s.LongUnit = "N. MILES"
		s.Z = 6080
	} else {
		s.G3 = 3.6
		s.G5 = 1000
		s.LongUnit = "KILOMETRES"
		s.ShortUnit = "METRES"
		s.Z = 1852.8
	}
}

func printTitle() {
	printCentered("LEM")
	printCentered("CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY")
}

func hasFlownBefore() bool {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\nLUNAR LANDING SIMULATION")
	fmt.Println("\nHAVE YOU FLOWN AN APOLLO/LEM MISSION BEFORE")

	for {
		fmt.Print(" (YES OR NO) ")
		scanner.Scan()

		if strings.ToUpper(scanner.Text())[0:1] == "Y" {
			return true
		} else if strings.ToUpper(scanner.Text())[0:1] == "N" {
			return false
		}

		fmt.Println("JUST ANSWER THE QUESTION, PLEASE, ")
	}
}

func main() {
	s := NewSim()

	printTitle()

	units := 1
	if hasFlownBefore() {
		units = getInt("\nINPUT MEASUREMENT OPTION NUMBER ")
	} else {
		units = getInt("\nWHICH SYSTEM OF MEASUREMENT DO YOU PREFER?\n 1=METRIC     0=ENGLISH\nENTER THE APPROPRIATE NUMBER ")
	}
	s.setUnits(units)

	fmt.Println(s)

}

func printCentered(s string) {
	div := (PAGE_WIDTH - len(s)) / 2

	fmt.Println(strings.Repeat(" ", div) + s)
}

func getInt(prompt string) int {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(prompt)
		scanner.Scan()

		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("INVALID INPUT, TRY AGAIN")
			continue
		}

		return val
	}
}
