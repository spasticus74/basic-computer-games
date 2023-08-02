package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	MAXHISTORY     int = 9
	LOSINGBOOKSIZE int = 50
)

var (
	moveCount  int
	lossBook   [LOSINGBOOKSIZE]int
	gameNumber int
	n          int
)

type gameBoard [14]int

func printWelcome() {
	fmt.Println("                                  AWARI")
	fmt.Println("               CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY")
	fmt.Println()
	fmt.Println()
}

func (board *gameBoard) drawPit(line string, pit int) string {
	val := board[pit]
	line += " "
	if val < 10 {
		line += " "
	}
	line += strconv.Itoa(val) + " "
	return line
}

func (board *gameBoard) draw() {
	fmt.Println()

	// draw the computer's pits
	line := "   "
	for i := 12; i > 6; i-- {
		line = board.drawPit(line, i)
	}
	fmt.Println(line)

	// draw the side pits
	line = board.drawPit("", 13)
	line += "                      "
	line = board.drawPit(line, 6)
	fmt.Println(line)

	// draw the bottom (player) pits
	line = "   "
	for i := 0; i < 6; i++ {
		line = board.drawPit(line, i)
	}
	fmt.Printf("%s\n\n\n", line)
}

func (board *gameBoard) playGame() {
	// initialise the board
	for i := 0; i < 13; i++ {
		board[i] = 3
	}

	// clear the home pits
	board[6] = 0
	board[13] = 0

	moveCount = 0

	// clear the history for this game
	lossBook[gameNumber] = 0

	for {
		board.draw()

		// player's turn
		fmt.Println("YOUR MOVE")
		landingSpot, stillGoing, home := board.playerMove()
		if !stillGoing {
			break
		}
		if landingSpot == home {
			_, stillGoing, _ = board.playerMoveAgain()
		}
		if !stillGoing {
			break
		}

		// computer's turn
		fmt.Println("MY MOVE")
		landingSpot, stillGoing, home, msg := board.computerMove("")
		if !stillGoing {
			fmt.Println(msg)
			break
		}
		if landingSpot == home {
			_, stillGoing, _, msg = board.computerMove(msg + " , ")
		}
		if !stillGoing {
			fmt.Println(msg)
			break
		}
		print(msg)
	}
	board.gameOver()
}

func (board *gameBoard) computerMove(message string) (int, bool, int, string) {
	bestQuality := -99
	savedBoard := board
	selectedMove := -1

	for cm := 7; cm < 13; cm++ {
		if board[cm] == 0 {
			continue
		}

		board.doMove(cm, 13)

		bestPlayerMoveQuality := 0

		for humanMoveStart := 0; humanMoveStart < 6; humanMoveStart++ {
			if board[humanMoveStart] == 0 {
				continue
			}

			humanMoveEnd := board[humanMoveStart] + humanMoveStart
			thisPlayerMoveQuality := 0

			for humanMoveEnd > 13 {
				humanMoveEnd = humanMoveEnd - 14
				thisPlayerMoveQuality++
			}

			if (board[humanMoveEnd] == 0) && (humanMoveEnd != 6) && (humanMoveEnd != 13) {
				thisPlayerMoveQuality += board[12-humanMoveEnd]
			}

			if thisPlayerMoveQuality > bestPlayerMoveQuality {
				bestPlayerMoveQuality = thisPlayerMoveQuality
			}
		}
		computerMoveQuality := board[13] - board[6] - bestPlayerMoveQuality

		if moveCount < MAXHISTORY {
			moveDigit := cm
			if moveDigit > 6 {
				moveDigit = moveDigit - 7
			}

			for previousGameNumber := 0; previousGameNumber < gameNumber; previousGameNumber++ {
				if (lossBook[gameNumber]*6 + moveDigit) == int(float32((lossBook[previousGameNumber]/6.0)^(7.0-moveCount))+0.1) {
					computerMoveQuality = computerMoveQuality - 2
				}
			}
		}

		for i := 0; i < 14; i++ {
			board[i] = savedBoard[i]
		}

		if computerMoveQuality >= bestQuality {
			selectedMove = cm
			bestQuality = computerMoveQuality
		}
	}

	moveString := strconv.Itoa(42 + selectedMove)

	if len(message) > 0 {
		message += (", " + moveString)
	} else {
		message = moveString
	}

	moveNumber, stillGoing, home := board.executeMove(selectedMove, 13)
	return moveNumber, stillGoing, home, message
}

func (board *gameBoard) gameOver() {
	fmt.Print("\nGAME OVER\n")

	pitDiff := board[6] - board[13]
	if pitDiff < 0 {
		fmt.Printf("I WIN BY %d POINTS\n", pitDiff*-1)
	} else {
		n++

		if pitDiff == 0 {
			fmt.Println("DRAWN GAME")
		} else {
			fmt.Printf("YOU WIN BY %d POINTS\n", pitDiff)
		}
	}
}

func (board *gameBoard) doCapture(m, home int) {
	board[home] += board[12-m] + 1
	board[m] = 0
	board[12-m] = 0
}

func (board *gameBoard) doMove(m, home int) int {
	moveStones := board[m]
	board[m] = 0

	for s := moveStones; s > 0; s-- {
		m = m + 1
		if m > 13 {
			m = m - 14
		}
		board[m] += 1
	}
	if board[m] == 1 && (m != 6) && (m != 13) && (board[12-m] != 0) {
		board.doCapture(m, home)
	}
	return m
}

func (board *gameBoard) playerHasStones() bool {
	for i := 0; i < 6; i++ {
		if board[i] > 0 {
			return true
		}
	}
	return false
}

func (board *gameBoard) computerHasStones() bool {
	for i := 7; i < 13; i++ {
		if board[i] > 0 {
			return true
		}
	}
	return false
}

func (board *gameBoard) executeMove(move, home int) (int, bool, int) {
	moveDigit := move
	lastLocation := board.doMove(move, home)
	var stillGoing bool

	if moveDigit > 6 {
		moveDigit = moveDigit - 7
	}

	moveCount++

	if moveCount < MAXHISTORY {
		lossBook[gameNumber] = lossBook[gameNumber]*6 + moveDigit
	}

	if board.playerHasStones() && board.computerHasStones() {
		stillGoing = true
	} else {
		stillGoing = false
	}
	return lastLocation, stillGoing, home
}

func (board *gameBoard) playerMoveAgain() (int, bool, int) {
	fmt.Println("AGAIN")
	return board.playerMove()
}

func (board *gameBoard) playerMove() (int, bool, int) {
	scanner := bufio.NewScanner(os.Stdin)
	var m int
	for {
		fmt.Println("SELECT MOVE (1-6):")
		scanner.Scan()
		m, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("ILLEGAL MOVE")
			continue
		}
		m--
		if (m > 5) || (m < 0) || (board[m] == 0) {
			fmt.Println("ILLEGAL MOVE")
			continue
		}
		break
	}
	ending_spot, is_still_going, home := board.executeMove(m, 6)

	board.draw()

	return ending_spot, is_still_going, home
}

func main() {
	printWelcome()

	board := gameBoard{}
	lossBook = [LOSINGBOOKSIZE]int{}

	for {
		board.playGame()
	}
}
