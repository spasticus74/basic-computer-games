package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	START_POPULATION = 95
	START_GRAIN      = 2800
	START_HARVEST    = 3000
	START_YIELD      = 3
	START_IMMIGRANTS = 5
)

type game struct {
	scanner     bufio.Scanner
	deaths      int
	deathRate   float64
	year        int
	population  int
	grain       int
	harvest     int
	yield       int
	costPerAcre int
	acres       int
	immigrants  int
	plague      bool
	eatenByRats int
}

func NewGame() game {
	g := game{}
	g.scanner = *bufio.NewScanner(os.Stdin)
	g.population = START_POPULATION
	g.grain = START_GRAIN
	g.harvest = START_HARVEST
	g.yield = START_YIELD
	g.immigrants = START_IMMIGRANTS
	g.acres = START_HARVEST / START_YIELD

	return g
}

func (g *game) getInt(prompt string) int {
	fmt.Println(prompt)
	for {
		g.scanner.Scan()
		val, err := strconv.Atoi(g.scanner.Text())
		if err != nil {
			fmt.Println("INVALD INPUT, TRY AGAIN")
			continue
		}
		return val
	}
}

func (g *game) reportStats() {
	fmt.Println("POPULATION IS NOW", g.population)
	fmt.Println("THE CITY NOW OWNS", g.acres, "ACRES.")
	fmt.Println("YOU HARVESTED", g.yield, "BUSHELS PER ACRE.")
	fmt.Println("THE RATS ATE", g.eatenByRats, "BUSHELS.")
	fmt.Println("YOU NOW HAVE ", g.grain, "BUSHELS IN STORE.")
}

func (g *game) printInsufficientGrain() {
	fmt.Println("HAMMURABI:  THINK AGAIN.  YOU HAVE ONLY")
	fmt.Printf("%d BUSHELS OF GRAIN.  NOW THEN,\n", g.grain)
}

func (g *game) printInsufficientAcres() {
	fmt.Printf("HAMMURABI:  THINK AGAIN.  YOU HAVE ONLY %d ACRES.  NOW THEN,\n", g.acres)
}

func (g *game) summariseGame() {
	fmt.Printf("IN YOUR 10 YEAR TERM OF OFFICE, %f PERCENT OF THE\n", g.deathRate)
	fmt.Println("POPULATION STARVED PER YEAR ON THE AVERAGE, I.E. A TOTAL OF")
	fmt.Printf("%d PEOPLE DIED!!\n", g.deaths)
	fmt.Println("YOU STARTED WITH 10 ACRES PER PERSON AND ENDED WITH")
	landRate := float64(g.acres) / float64(g.population)
	fmt.Printf("%0.2f ACRES PER PERSON.\n", landRate)

	if g.deathRate > 33 || landRate < 7 {
		printNationalFink()
	} else if g.deathRate > 10 || landRate < 9 {
		fmt.Println("YOUR HEAVY-HANDED PERFORMANCE SMACKS OF NERO AND IVAN IV.")
		fmt.Println("THE PEOPLE (REMIANING) FIND YOU AN UNPLEASANT RULER, AND,")
		fmt.Println("FRANKLY, HATE YOUR GUTS!!")
	} else if g.deathRate > 3 || landRate < 10 {
		fmt.Println("YOUR PERFORMANCE COULD HAVE BEEN SOMEWHAT BETTER, BUT")
		fmt.Printf("REALLY WASN'T TOO BAD AT ALL. %d PEOPLE\n", int(float64(g.population)*0.8*rand.Float64()))
		fmt.Println("WOULD DEARLY LIKE TO SEE YOU ASSASSINATED BUT WE ALL HAVE OUR")
		fmt.Println("TRIVIAL PROBLEMS.")
	} else {
		fmt.Println("A FANTASTIC PERFORMANCE!!!  CHARLEMANGE, DISRAELI, AND")
		fmt.Println("JEFFERSON COMBINED COULD NOT HAVE DONE BETTER!")
		fmt.Println()
	}
}

func (g *game) play() {
	people := 0

	for g.year < 11 {
		fmt.Println("\n\n\nHAMMURABI:  I BEG TO REPORT TO YOU")
		g.year++
		fmt.Printf("IN YEAR %d, %d PEOPLE STARVED, %d CAME TO THE CITY\n", g.year, people, g.immigrants)

		if g.plague {
			g.population = int(g.population / 2)
			fmt.Println("A HORRIBLE PLAGUE STRUCK!  HALF THE PEOPLE DIED.")
		}

		g.reportStats()

		// recalculate the cost of land
		g.costPerAcre = int(rand.Float64()*10) + 17
		fmt.Println("LAND IS TRADING AT", g.costPerAcre, "BUSHELS PER ACRE.")

		// buy land?
	BUYLAND:
		newAcres := g.getInt("HOW MANY ACRES TO YOU WISH TO BUY?")
		if newAcres < 0 {
			printBadInput()
			break // REALLY?
		} else if g.costPerAcre*newAcres > g.grain { // more than we can afford
			g.printInsufficientGrain()
			goto BUYLAND // allow the player to try again
		} else {
			g.acres += newAcres
			g.grain -= (newAcres * g.costPerAcre)
		}

		// sell land?
	SELLLAND:
		sellLand := g.getInt("HOW MANY ACRES DO YOU WISH TO SELL? ")
		if sellLand < 0 {
			printBadInput()
			break // REALLY?
		} else if sellLand <= g.acres {
			g.acres -= sellLand
			g.grain += (g.costPerAcre * sellLand)
		} else { // tried to sell more than owned
			g.printInsufficientAcres()
			goto SELLLAND
		}
		fmt.Println()

		// feed people
	FEED:
		feed := g.getInt("HOW MANY BUSHELS DO YOU WISH TO FEED YOUR PEOPLE? ")
		if feed < 0 {
			printBadInput()
			break
		} else if feed > g.grain {
			g.printInsufficientGrain()
			goto FEED
		} else {
			g.grain -= feed
		}
		fmt.Println()

		// plant the fields
	PLANT:
		acresToPlant := g.getInt("HOW MANY ACRES DO YOU WISH TO PLANT WITH SEED? ")
		if acresToPlant < 0 {
			printBadInput()
			break
		} else if acresToPlant > 0 {
			if acresToPlant > g.acres { // trying to plant more acres than owned
				g.printInsufficientAcres()
				goto PLANT
			} else if int(acresToPlant/2) > g.grain { // enough grain for seed?
				g.printInsufficientGrain()
				goto PLANT
			} else if acresToPlant > (10 * g.population) { // enough people to tend the crops?
				fmt.Printf("BUT YOU HAVE ONLY %d PEOPLE TO TEND THE FIELDS!  NOW, THEN,", g.population)
				goto PLANT
			} else {
				g.grain -= int(acresToPlant / 2)
			}

		}

		// harvest
		g.yield = getRandInt()
		g.harvest = g.yield * g.acres
		g.eatenByRats = 0

		// rat infestation?
		if rand.Float64() > 0.5 {
			g.eatenByRats = g.grain / getRandInt()
		}

		// update grain storage
		g.grain = g.grain - g.eatenByRats + g.harvest

		// calculate immigration/births
		g.immigrants = int(getRandInt()*(20*g.acres+g.grain)/g.population/100 + 1)

		// plague?
		if rand.Float64() <= 0.15 {
			g.plague = true
		} else {
			g.plague = false
		}

		// starvation?
		fullPeople := int(feed / 20)
		if g.population > fullPeople {
			hungryPeople := g.population - fullPeople
			if float64(hungryPeople) > (0.45 * float64(g.population)) {
				fmt.Printf("YOU STARVED %d PEOPLE IN ONE YEAR!!!\n", hungryPeople)
				printNationalFink()
				break
			}
			g.deathRate = ((float64(g.year)-1)*g.deathRate + float64(hungryPeople)*100.0/float64(g.population)) / float64(g.year)
			g.population = fullPeople
			g.deaths += hungryPeople
		}
	}
}

func printTitle() {
	fmt.Println("                                HAMMURABI")
	fmt.Println("               CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY")
	fmt.Print("\n\n\n")
	fmt.Println("TRY YOUR HAND AT GOVERNING ANCIENT SUMERIA")
	fmt.Println("FOR A TEN-YEAR TERM OF OFFICE.")
	fmt.Println()
}

func printNationalFink() {
	fmt.Println("DUE TO THIS EXTREME MISMANAGEMENT YOU HAVE NOT ONLY")
	fmt.Println("BEEN IMPEACHED AND THROWN OUT OF OFFICE BUT YOU HAVE")
	fmt.Println("ALSO BEEN DECLARED NATIONAL FINK!!!!")
}

func getRandInt() int {
	return int(rand.Float64()*5) + 1
}

func printBadInput() {
	fmt.Println("\nHAMMURABI:  I CANNOT DO WHAT YOU WISH.")
	fmt.Println("GET YOURSELF ANOTHER STEWARD!!!!!")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	printTitle()

	g := NewGame()

	g.play()

	// did we survive the full term?
	if g.year == 10 {
		g.summariseGame()
	}

	fmt.Println("\nSO LONG FOR NOW.")
	fmt.Println()
}
