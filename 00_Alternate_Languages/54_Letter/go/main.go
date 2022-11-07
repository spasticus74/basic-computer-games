package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func printIntro() {
	fmt.Println("                                 LETTER")
	fmt.Println("               CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY")
	fmt.Print("\n\n\n")
}

func printInstructions() {
	fmt.Println("LETTER GUESSING GAME")
	fmt.Println()
	fmt.Println("I'LL THINK OF A LETTER OF THE ALPHABET, A TO Z.")
	fmt.Println("TRY TO GUESS MY LETTER AND I'LL GIVE YOU CLUES")
	fmt.Println("AS TO HOW CLOSE YOU'RE GETTING TO MY LETTER.")
}

func playGame() {
	scanner := bufio.NewScanner(os.Stdin)

	// get the target
	alphabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	target := alphabet[rand.Intn(len(alphabet))]

	fmt.Println("\nO.K., I HAVE A LETTER.  START GUESSING.")
	fmt.Println()

	// play
	guessCount := 0
	for {
		fmt.Println("WHAT IS YOUR GUESS?")
		scanner.Scan()
		guessCount++

		fmt.Println()
		guess := strings.ToUpper(scanner.Text())
		if guess == target {
			fmt.Printf("YOU FOR IT IN %d GUESSES!!\n", guessCount)
			if guessCount > 5 {
				fmt.Println("BUT IT SHOULDN'T TAKE MORE THAN 5 GUESSES!")
			}
			fmt.Println("GOOD JOB !!!!!")

			fmt.Println("\nLET'S PLAY AGAIN.....")
			return
		} else if guess > target {
			fmt.Println("TOO HIGH. TRY A LOWER LETTER.")
		} else {
			fmt.Println("TOO LOW. TRY A HIGHER LETTER.")
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	printIntro()

	printInstructions()

	for {
		playGame()
	}
}
