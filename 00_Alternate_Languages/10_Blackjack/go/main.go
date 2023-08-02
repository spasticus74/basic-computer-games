package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	NUMBER_OF_DECKS  = 3
	STARTING_BALANCE = 100
)

type PlayerType int8

const (
	PLAYER PlayerType = iota
	DEALER
)

func (s PlayerType) String() string {
	switch s {
	case PLAYER:
		return "Player"
	case DEALER:
		return "Dealer"
	}
	return "unknown"
}

type Action int8

const (
	STAND Action = iota
	HIT
	DOUBLEDOWN
	SPLIT
)

func (a Action) String() string {
	switch a {
	case STAND:
		return "Stand"
	case HIT:
		return "Hit!"
	case DOUBLEDOWN:
		return "Double down!"
	case SPLIT:
		return "Split"
	}
	return "unknown"
}

type Card struct {
	name  string
	value int
}

// HAND
type Hand struct {
	cards []*Card
}

func (h *Hand) addCard(c *Card) {
	h.cards = append(h.cards, c)
}

func (h *Hand) getTotal() int {
	total := 0
	haveAce := false
	for _, c := range h.cards {
		total += c.value
		if c.name == "ACE" {
			haveAce = true
		}
	}

	if total > 21 && haveAce {
		total -= 10
	}

	return total
}

func (h *Hand) discard(discard *Deck) {
	discard.cards = append(discard.cards, h.cards...)

	h.cards = nil
}

// DECK
type Deck struct {
	cards []*Card
}

func NewDeck() *Deck {
	d := Deck{}

	for n := 0; n < NUMBER_OF_DECKS; n++ {
		for c := 0; c < 14; c++ {
			for m := 0; m < 4; m++ {
				switch c {
				case 0:
					d.cards = append(d.cards, &Card{"ACE", 11})
				case 1:
					d.cards = append(d.cards, &Card{"2", 2})
				case 2:
					d.cards = append(d.cards, &Card{"3", 3})
				case 3:
					d.cards = append(d.cards, &Card{"4", 4})
				case 5:
					d.cards = append(d.cards, &Card{"5", 5})
				case 6:
					d.cards = append(d.cards, &Card{"6", 6})
				case 7:
					d.cards = append(d.cards, &Card{"7", 7})
				case 8:
					d.cards = append(d.cards, &Card{"8", 8})
				case 9:
					d.cards = append(d.cards, &Card{"9", 9})
				case 10:
					d.cards = append(d.cards, &Card{"10", 10})
				case 11:
					d.cards = append(d.cards, &Card{"JACK", 10})
				case 12:
					d.cards = append(d.cards, &Card{"QUEEN", 10})
				case 13:
					d.cards = append(d.cards, &Card{"KING", 10})
				}
			}
		}
	}

	d.shuffle()
	return &d
}

func (d *Deck) shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] })
}

func (d *Deck) drawCard(discard *Deck) *Card {
	if len(d.cards) == 0 {
		if len(discard.cards) > 0 {
			fmt.Println("deck is empty, shuffling")
			d.cards = discard.cards
			discard.cards = nil
			d.shuffle()
			return d.drawCard(discard)
		}
	}
	card := d.cards[len(d.cards)-1]
	d.cards = d.cards[:len(d.cards)-1]
	return card
}

// PLAYER
type Player struct {
	hand       *Hand
	balance    int
	bet        int
	wins       int
	playerType PlayerType
	index      int
}

func NewPlayer(t PlayerType, i int) *Player {
	p := Player{}

	p.hand = &Hand{}

	p.playerType = t
	p.balance = STARTING_BALANCE
	p.index = i

	return &p
}

func (p *Player) getName() string {
	return fmt.Sprintf("%s%d", p.playerType.String(), p.index)
}

func (p *Player) getBet() {
	if p.playerType == PLAYER {
		if p.balance < 1 {
			fmt.Printf("%s is out of money :(\n", p.getName())
			p.bet = 0
		} else {
			p.bet = getNumberFromUser("What is your bet", 1, p.balance)
		}
	}
}

func (p *Player) handToString(hideDealer bool) string {
	s := ""
	if !hideDealer || p.playerType == PLAYER {
		for i := 0; i < len(p.hand.cards); i++ {
			s += fmt.Sprintf("%s\t", p.hand.cards[i].name)
		}
		s += fmt.Sprintf("total points = %d", p.hand.getTotal())
		return s
	} else {
		for i, v := range p.hand.cards {
			if i == 0 {
				s += v.name
			} else {
				s += "\t*"
			}
		}
	}
	return s
}

func (p *Player) getAction() Action {
	if p.playerType == DEALER {
		if p.hand.getTotal() > 16 {
			return STAND
		} else {
			return HIT
		}
	}

	var validResults []string
	if len(p.hand.cards) > 2 {
		validResults = []string{"s", "h"}
	} else {
		validResults = []string{"s", "h", "d", "/"}
	}

	for {
		action := getCharFromUser("\tWhat is your play?", validResults)

		switch action {
		case "s":
			return STAND
		case "h":
			return HIT
		case "d":
			return DOUBLEDOWN
		case "/":
			return SPLIT
		default:
			fmt.Println("Invalid, try again")
		}
	}
}

// GAME
type Game struct {
	players     []Player
	playDeck    Deck
	discard     Deck
	gamesPlayed int
}

func NewGame(playerCount int) *Game {
	g := Game{}

	g.players = append(g.players, *NewPlayer(DEALER, 0))
	for i := 0; i < playerCount; i++ {
		g.players = append(g.players, *NewPlayer(PLAYER, i+1))
	}

	if getCharFromUser("Do you want instructions?", []string{"y", "n"}) == "y" {
		printInstructions()
	}
	fmt.Println()

	g.playDeck = *NewDeck()
	g.discard = Deck{}

	return &g
}

func (g *Game) printStats() {
	fmt.Print(g.statsAsString())
}

func (g *Game) statsAsString() string {
	s := "Scores:\n"
	for _, p := range g.players {
		s += fmt.Sprintf("%s Wins:\t%d", p.getName(), p.wins)
		if p.playerType == PLAYER {
			s += fmt.Sprintf("\t\tBalance:\t%d\t\tBet:\t%d", p.balance, p.bet)
		}
		s += "\n"
	}
	return s
}

func (g *Game) play() {
	game := g.gamesPlayed
	playerHandsMessage := ""

	// deal
	for d := 0; d < 2; d++ {
		for _, p := range g.players {
			p.hand.addCard(g.playDeck.drawCard(&g.discard))
		}
	}

	// get bets
	for _, p := range g.players {
		p.getBet()
	}
	scores := g.statsAsString()

	// play for each player
	for _, p := range g.players {
		for {
			clearScreen()
			printWelcome()
			fmt.Printf("\n\t\tGame %d\n", game)
			fmt.Println(scores)
			fmt.Println(playerHandsMessage)
			fmt.Printf("%s Hand:\t%s\n", p.getName(), p.handToString(true))

			if p.playerType == PLAYER && p.bet == 0 {
				break
			}

			// play this turn
			// first check for done
			score := p.hand.getTotal()
			if score >= 21 {
				if score == 21 {
					fmt.Println("\tBlackjack! (21 points)")
				} else {
					fmt.Printf("\tBust       (%d points)\n", score)
				}
				break
			}

			// get their move
			play := p.getAction()

			// process the play
			switch play {
			case STAND:
				fmt.Printf("\t%s\n", STAND.String())
			case HIT:
				fmt.Printf("\t%s\n", HIT.String())
				p.hand.addCard(g.playDeck.drawCard(&g.discard))
			case DOUBLEDOWN:
				fmt.Printf("\t%s\n", DOUBLEDOWN.String())
				if p.bet*2 < p.balance {
					p.bet *= 2
				} else {
					p.bet = p.balance
				}
				p.hand.addCard(g.playDeck.drawCard(&g.discard))
			}
		}
		playerHandsMessage += fmt.Sprintf("%s Hand:\t %s\n", p.getName(), p.handToString(true))
	}

	// determine the winner
	topScore := 0
	numWinners := 1

	for _, p := range g.players {
		if p.hand.getTotal() <= 21 {
			score := p.hand.getTotal()
			if score > topScore {
				topScore = score
				numWinners = 1
			} else if score == topScore {
				numWinners += 1
			}
		}
	}

	for _, p := range g.players {
		if p.hand.getTotal() == topScore {
			fmt.Printf("%s ", p.getName())
			p.wins += 1
			p.balance += p.bet
		} else {
			p.balance -= p.bet
		}

		p.hand.discard(&g.discard)
	}

	if numWinners > 1 {
		fmt.Printf("all tie with %d\n\n\n", topScore)
	} else {
		fmt.Printf("wins with %d\n", topScore)
	}

	g.gamesPlayed += 1
}

func main() {
	rand.Seed(time.Now().UnixNano())

	printWelcome()

	//g := NewGame(getNumberFromUser("How many players should there be?", 1, 7))
	g := NewGame(1)

	for {
		g.play()
		if getCharFromUser("Play again?", []string{"y", "n"}) == "n" {
			return
		}
	}
}

func printWelcome() {
	fmt.Println("                            BLACK JACK")
	fmt.Println("            CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY")
	fmt.Println()
}

func printInstructions() {
	fmt.Println("THIS IS THE GAME OF 21. AS MANY AS 7 PLAYERS MAY PLAY THE")
	fmt.Println("GAME. ON EACH DEAL, BETS WILL BE ASKED FOR, AND THE")
	fmt.Println("PLAYERS' BETS SHOULD BE TYPED IN. THE CARDS WILL THEN BE")
	fmt.Println("DEALT, AND EACH PLAYER IN TURN PLAYS HIS HAND. THE")
	fmt.Println("FIRST RESPONSE SHOULD BE EITHER 'D', INDICATING THAT THE")
	fmt.Println("PLAYER IS DOUBLING DOWN, 'S', INDICATING THAT HE IS")
	fmt.Println("STANDING, 'H', INDICATING HE WANTS ANOTHER CARD, OR '/',")
	fmt.Println("INDICATING THAT HE WANTS TO SPLIT HIS CARDS. AFTER THE")
	fmt.Println("INITIAL RESPONSE, ALL FURTHER RESPONSES SHOULD BE 'S' OR")
	fmt.Println("'H', UNLESS THE CARDS WERE SPLIT, IN WHICH CASE DOUBLING")
	fmt.Println("DOWN IS AGAIN PERMITTED. IN ORDER TO COLLECT FOR")
	fmt.Println("BLACKJACK, THE INITIAL RESPONSE SHOULD BE 'S'.")
	fmt.Println("NUMBER OF PLAYERS")
	fmt.Println()
	fmt.Println("NOTE:'/' (splitting) is not currently implemented, and does nothing")
	fmt.Println()
	fmt.Println("PRESS ENTER TO CONTINUE")
}

func getNumberFromUser(prompt string, min, max int) int {
	scanner := bufio.NewScanner(os.Stdin)
	input := min
	var err error

	for {
		fmt.Printf("%s (%d - %d)\n", prompt, min, max)
		scanner.Scan()
		input, err = strconv.Atoi(scanner.Text())
		//input, err = strconv.Atoi("2")
		if err == nil && input >= min && input <= max {
			return input
		} else {
			fmt.Println("Invalid input, please try again")
		}
	}
}

func getCharFromUser(prompt string, validResults []string) string {
	input := ""
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%s, [%s]\n", prompt, strings.Join(validResults, ","))
		scanner.Scan()
		input = scanner.Text()
		input = strings.ToLower(string(input[0]))
		for _, v := range validResults {
			if input == strings.ToLower(v) {
				return input
			}
		}
		fmt.Println("Invalid input, please try again")
	}
}

func clearScreen() {
	fmt.Println("\x1b[2J\x1b[0;0H")
}
