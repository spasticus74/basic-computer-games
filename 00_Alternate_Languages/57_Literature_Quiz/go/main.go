package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Question struct {
	statement     string
	options       []string
	answer        int
	wrongResponse string
	rightResponse string
}

func NewQuestion(s, w, r string, a int, o []string) Question {
	q := Question{}
	q.options = make([]string, len(o))

	q.statement = s
	q.options = o
	q.answer = a
	q.wrongResponse = w
	q.rightResponse = r
	return q
}

func (q *Question) ask() bool {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(q.statement)
	for i, option := range q.options {
		fmt.Printf("%d. %s, ", i+1, option)
	}
	for {
		fmt.Println()
		scanner.Scan()
		val, err := strconv.Atoi(scanner.Text())
		if err != nil || val < 1 || val > len(q.options) {
			fmt.Println("INVALID INPUT. TRY AGAIN.")
			continue
		}

		if val == q.answer {
			fmt.Println(q.rightResponse)
			return true
		} else {
			fmt.Println(q.wrongResponse)
			return false
		}
	}
}

func printInstructions() {
	fmt.Println("TEST YOUR KNOWLEDGE OF CHILDREN'S LITERATURE.")
	fmt.Println()
	fmt.Println("THIS IS A MULTIPLE-CHOICE QUIZ.")
	fmt.Println("TYPE A 1, 2, 3, OR 4 AFTER THE QUESTION MARK.")
	fmt.Println()
	fmt.Println("GOOD LUCK!")
	fmt.Println()
	fmt.Println()
}

func main() {
	quiz := []Question{}
	quiz = append(quiz, NewQuestion("IN PINOCCHIO, WHAT WAS THE NAME OF THE CAT?", "SORRY...FIGARO WAS HIS NAME.", "VERY GOOD!  HERE'S ANOTHER.", 3, []string{"TIGGER", "CICERO", "FIGARO", "GUIPETTO"}))
	quiz = append(quiz, NewQuestion("FROM WHOSE GARDEN DID BUGS BUNNY STEAL THE CARROTS?", "TOO BAD...IT WAS ELMER FUDD'S GARDEN.", "PRETTY GOOD!", 2, []string{"MR. NIXON'S", "ELMER FUDD'S", "CLEM JUDD'S", "STROMBOLI'S"}))
	quiz = append(quiz, NewQuestion("IN THE WIZARD OF OS, DOROTHY'S DOG WAS NAMED?", "BACK TO THE BOOKS,...TOTO WAS HIS NAME.", "YEA!  YOU'RE A REAL LITERATURE GIANT.", 4, []string{"CICERO", "TRIXIA", "KING", "TOTO"}))
	quiz = append(quiz, NewQuestion("WHO WAS THE FAIR MAIDEN WHO ATE THE POISON APPLE?", "OH, COME ON NOW...IT WAS SNOW WHITE.", "GOOD MEMORY!", 3, []string{"SLEEPING BEAUTY", "CINDERELLA", "SNOW WHITE", "WENDY"}))

	printCentered("LITERATURE QUIZ", 64, " ")
	printCentered("CREATIVE COMPUTING MORRISTOWN, NEW JERSEY", 64, " ")
	fmt.Print("\n\n\n")

	printInstructions()

	score := 0

	for _, q := range quiz {
		if q.ask() {
			score++
		}
		fmt.Print("\n\n")
	}

	if score == len(quiz) {
		fmt.Println("WOW!  THAT'S SUPER!  YOU REALLY KNOW YOUR NURSERY")
		fmt.Println("YOUR NEXT QUIZ WILL BE ON 2ND CENTURY CHINESE")
		fmt.Println("LITERATURE (HA, HA, HA)")
	} else if score < len(quiz)/2 {
		fmt.Println("UGH.  THAT WAS DEFINITELY NOT TOO SWIFT.  BACK TO")
		fmt.Println("NURSERY SCHOOL FOR YOU, MY FRIEND.")
	} else {
		fmt.Println("NOT BAD, BUT YOU MIGHT SPEND A LITTLE MORE TIME")
		fmt.Println("READING THE NURSERY GREATS.")
	}
}

func printCentered(s string, n int, fill string) {
	div := (n / 2) - len(s)/2

	fmt.Println(strings.Repeat(fill, div) + s)
}
