package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	FOREST_LAND      = 1000
	INITIAL_LAND     = 2000
	COST_OF_LIVING   = 100
	COST_OF_FUNERAL  = 9
	TERM_LENGTH      = 8
	POLLUTION_FACTOR = 25
)

type Game struct {
	rallods        int
	countrymen     int
	land           int
	foreignWorkers int
	yearsInOffice  int

	// previous year
	cropLossLastYear int

	// current year
	countrymenDeaths int
	pollutionDeaths  int
	populationChange int

	// current market
	plantingCost  int
	landPrice     int
	tourismIncome int

	scanner bufio.Scanner
}

func NewGame() Game {
	g := Game{}

	g.land = INITIAL_LAND

	g.rallods = rand.Intn(2001) + 59000

	g.countrymen = rand.Intn(21) + 490

	g.plantingCost = rand.Intn(6) + 10

	g.scanner = *bufio.NewScanner(os.Stdin)

	return g
}

func resume() Game {
	g := NewGame()

	for {
		years := getInt("HOW MANY YEARS HAD YOU BEEN IN OFFICE WHEN INTERRUPTED? ")
		if years < 0 {
			os.Exit(0)
		} else if years >= TERM_LENGTH {
			fmt.Printf("   COME ON, YOUR TERM IN OFFICE IS ONLY %d YEARS.\n", TERM_LENGTH)
		} else {
			g.yearsInOffice = years
			break
		}
	}

	funds := getInt("HOW MUCH DID YOU HAVE IN THE TREASURY? ")
	if funds < 0 {
		os.Exit(0)
	} else {
		g.rallods = funds
	}

	population := getInt("HOW MANY COUNTRYMEN? ")
	if population < 0 {
		os.Exit(0)
	} else {
		g.countrymen = population
	}

	workers := getInt("HOW MANY WORKERS? ")
	if workers < 0 {
		os.Exit(0)
	} else {
		g.foreignWorkers = workers
	}

	for {
		land := getInt("HOW MANY SQUARE MILES OF LAND? ")
		if land < 0 {
			os.Exit(0)
		} else if land > INITIAL_LAND {
			farmLand := INITIAL_LAND - FOREST_LAND
			fmt.Printf("   COME ON, YOU STARTED WITH %d SQ. MILES OF FARM LAND\n", farmLand)
			fmt.Printf("   AND %d SQ. MILES OF FOREST LAND.\n", FOREST_LAND)
		} else if land > FOREST_LAND {
			break
		}
	}

	return g
}

func (g *Game) setMarket() {
	g.landPrice = rand.Intn(11) + 95
	g.plantingCost = rand.Intn(6) + 10
}

func (g *Game) farmland() int {
	return g.land - FOREST_LAND
}

func (g *Game) settledPeople() int {
	return g.countrymen - g.populationChange
}

func (g *Game) sellLand(amount int) {
	if amount > g.farmland() {
		log.Fatal("invalid value supplied to 'sellLand'")
	}

	g.land -= amount
	g.rallods += (g.landPrice * amount)
}

func (g *Game) distributeRallods(amount int) {
	g.rallods -= amount
}

func (g *Game) spendPollutionControl(amount int) {
	g.rallods -= amount
}

func (g *Game) plantCrops(area int) {
	g.rallods -= (area * g.plantingCost)
}

func (g *Game) printStatus() {
	fmt.Printf("\n\nYOU NOW HAVE %d RALLODS in THE TREASURY.\n", g.rallods)
	fmt.Printf("%d COUNTRYMEN, ", g.countrymen)
	if g.foreignWorkers > 0 {
		fmt.Printf("%d FOREIGN WORKERS, ", g.foreignWorkers)
	}
	fmt.Printf("AND %d SQ. MILES OF LAND.\n", g.land)
	fmt.Printf("THIS YEAR INDUSTRY WILL BUY LAND FOR %d RALLODS PER SQUARE MILE.\n", g.landPrice)
	fmt.Printf("LAND CURRENTLY COSTS %d RALLODS PER SQUARE MILE TO PLANT.\n\n", g.plantingCost)
}

func (g *Game) handleDeaths(distributeRallods, pollutionControlSpent int) {
	// starvation
	starvedCountrymen := int(math.Max(0, float64(g.countrymen)-float64(distributeRallods)/COST_OF_LIVING))
	if starvedCountrymen > 0 {
		fmt.Printf("%d COUNTRYMEN DIED OF STARVATION\n", starvedCountrymen)
	}

	// pollution
	g.pollutionDeaths = int(rand.Float64() * float64(INITIAL_LAND-g.land))
	if pollutionControlSpent >= POLLUTION_FACTOR {
		g.pollutionDeaths = int(float64(g.pollutionDeaths) / (float64(pollutionControlSpent) / float64(POLLUTION_FACTOR)))
	}
	if g.pollutionDeaths > 0 {
		fmt.Printf("%d COUNTRYMEN DIED OF CARBON MONOXIDE AND DUST INHALATION\n", g.pollutionDeaths)
	}

	g.countrymenDeaths = starvedCountrymen + g.pollutionDeaths
	if g.countrymenDeaths > 0 {
		funeralCosts := g.countrymenDeaths * COST_OF_FUNERAL
		fmt.Printf("   YOU WERE FORCED TO SPEND %d RALLODS ON FUNERAL EXPENSES\n", funeralCosts)
		g.rallods -= funeralCosts

		if g.rallods < 0 {
			fmt.Println("   INSUFFICIENT RESERVES TO COVER COST - LAND WAS SOLD")
			g.land += int(float64(g.rallods) / float64(g.landPrice))
			g.rallods = 0
		}
		g.countrymen -= g.countrymenDeaths
	}
}

func (g *Game) handleTourism() {
	v1 := (g.settledPeople() * 22) + int(rand.Float64()*500.0)
	v2 := (INITIAL_LAND - g.land) * 15

	tourismEarnings := v1 - v2
	fmt.Printf("YOU MADE %d RALLODS FROM TOURIST TRADE.\n", tourismEarnings)

	if v2 != 0 && !(v1-v2 >= g.tourismIncome) {
		fmt.Print("   DECREASE BECAUSE ")
		reason := rand.Intn(10)
		switch {
		case reason <= 2:
			fmt.Println("FISH POPULATION HAS DWINDLED DUE TO WATER POLLUTION.")
		case reason <= 4:
			fmt.Println("AIR POLLUTION IS KILLING GAME BIRD POPULATION.")
		case reason <= 6:
			fmt.Println("MINERAL BATHS ARE BEING RUINED BY WATER POLLUTION.")
		case reason <= 8:
			fmt.Println("UNPLEASANT SMOG IS DISCOURAGING SUN BATHERS.")
		case reason <= 10:
			fmt.Println("HOTELS ARE LOOKING SHABBY DUE TO SMOG GRIT.")
		}
	}

	g.tourismIncome = int(math.Abs(float64(v1) - float64(v2)))
	g.rallods += g.tourismIncome
}

func (g *Game) handleHarvest(area int) {
	cropLoss := int(float64(INITIAL_LAND-g.land) * ((rand.Float64() + 1.5) / 2))

	if g.foreignWorkers != 0 {
		fmt.Printf("OF %d SQ. MILES PLANTED,", area)
	}
	if area <= cropLoss {
		cropLoss = area
	}
	harvested := area - cropLoss
	fmt.Printf(" YOU HARVESTED %d SQ. MILES OF CROPS.\n", harvested)

	if cropLoss != 0 {
		fmt.Print("   (DUE TO ")
		if (cropLoss - g.cropLossLastYear) > 2 {
			fmt.Print("INCREASED ")
		}
		fmt.Println("AIR AND WATER POLLUTION FROM FOREIGN INDUSTRY.)")
	}

	revenue := (area - cropLoss) * int(float64(g.landPrice)/2.0)
	fmt.Printf("MAKING %d RALLODS.\n", revenue)

	g.cropLossLastYear = cropLoss
	g.rallods += revenue
}

func (g *Game) handleForeignWorkers(soldToIndustry, distributedFunds, pollutionControlSpent int) {
	foreignWorkersInflux := 0
	if soldToIndustry != 0 {
		foreignWorkersInflux = soldToIndustry + int((rand.Float64()*10)-(rand.Float64()*20))
		if g.foreignWorkers <= 0 {
			foreignWorkersInflux += 20
		}
		fmt.Printf("%d FOREIGN WORKERS CAME TO THE COUNTRY AND\n", foreignWorkersInflux)
	}

	surplusDistributed := distributedFunds/COST_OF_LIVING - g.countrymen

	populationChange := int((float64(surplusDistributed) / 10) + (float64(pollutionControlSpent) / float64(POLLUTION_FACTOR)) - (float64(INITIAL_LAND-g.land) / 50) - (float64(g.countrymenDeaths) / 2))
	fmt.Printf("%d COUNTRYMEN ", int(math.Abs(float64(populationChange))))
	if populationChange < 0 {
		fmt.Print("LEFT ")
	} else {
		fmt.Print("CAME TO ")
	}
	fmt.Println("THE ISLAND")

	g.countrymen += populationChange
	g.foreignWorkers += foreignWorkersInflux
}

func (g *Game) handleExcessiveDeaths() {
	fmt.Printf("\n\n\n%d COUNTRYMEN DIED IN ONE YEAR!!!!!\n", g.countrymenDeaths)
	fmt.Println("\n\n\nDUE TO THIS EXTREME MISMANAGEMENT, YOU HAVE NOT ONLY")
	fmt.Println("BEEN IMPEACHED AND THROWN OUT OF OFFICE, BUT YOU")
	message := rand.Intn(10)
	if message <= 3 {
		print("ALSO HAD YOUR LEFT EYE GOUGED OUT!")
	} else if message <= 6 {
		print("HAVE ALSO GAINED A VERY BAD REPUTATION.")
	} else if message <= 10 {
		print("HAVE ALSO BEEN DECLARED NATIONAL FINK.")
	}
	os.Exit(0)
}

func (g *Game) handleThirdDied() {
	fmt.Println()
	fmt.Println()
	fmt.Println("OVER ONE THIRD OF THE POPULTATION HAS DIED SINCE YOU")
	fmt.Println("WERE ELECTED TO OFFICE. THE PEOPLE (REMAINING)")
	fmt.Println("HATE YOUR GUTS.")
	g.endGame()
}

func (g *Game) handleMoneyMismanagement() {
	fmt.Println()
	fmt.Println("MONEY WAS LEFT OVER IN THE TREASURY WHICH YOU DID")
	fmt.Println("NOT SPEND. AS A RESULT, SOME OF YOUR COUNTRYMEN DIED")
	fmt.Println("OF STARVATION. THE PUBLIC IS ENRAGED AND YOU HAVE")
	fmt.Println("BEEN FORCED TO EITHER RESIGN OR COMMIT SUICIDE.")
	fmt.Println("THE CHOICE IS YOURS.")
	fmt.Println("IF YOU CHOOSE THE LATTER, PLEASE TURN OFF YOUR COMPUTER")
	fmt.Println("BEFORE PROCEEDING.")
	os.Exit(0)
}

func (g *Game) handleExcessiveForeigners() {
	fmt.Println("\n\nTHE NUMBER OF FOREIGN WORKERS HAS EXCEEDED THE NUMBER")
	fmt.Println("OF COUNTRYMEN. AS A MINORITY, THEY HAVE REVOLTED AND")
	fmt.Println("TAKEN OVER THE COUNTRY.")
	g.endGame()
}

func (g *Game) endGame() {
	if rand.Float64() < 0.5 {
		fmt.Println("YOU HAVE BEEN ASSASSINATED.")
	} else {
		fmt.Println("YOU HAVE BEEN THROWN OUT OF OFFICE AND ARE NOW")
		fmt.Println("RESIDING IN PRISON.")
	}
	os.Exit(0)
}

func (g *Game) handleCongratulations() {
	fmt.Println("\n\nCONGRATULATIONS!!!!!!!!!!!!!!!!!!")
	fmt.Printf("YOU HAVE SUCCESFULLY COMPLETED YOUR %d YEAR TERM\n", TERM_LENGTH)
	fmt.Println("OF OFFICE. YOU WERE, OF COURSE, EXTREMELY LUCKY, BUT")
	fmt.Println("NEVERTHELESS, IT'S QUITE AN ACHIEVEMENT. GOODBYE AND GOOD")
	fmt.Println("LUCK - YOU'LL PROBABLY NEED IT IF YOU'RE THE TYPE THAT")
	fmt.Println("PLAYS THIS GAME.")
	os.Exit(0)
}

func (g *Game) getAreaToPlant() int {
	for {
		fmt.Println("HOW MANY SQUARE MILES DO YOU WISH TO PLANT? ")
		g.scanner.Scan()

		area, err := strconv.Atoi(g.scanner.Text())
		if err != nil || area < 0 {
			fmt.Println("INVALID INPUT")
			continue
		}

		if area > (g.countrymen * 2) {
			fmt.Println("   SORRY, BUT EACH COUNTRYMAN CAN ONLY PLANT 2 SQ. MILES.")
		} else if area > g.farmland() {
			fmt.Printf("   SORRY BUT YOU ONLY HAVE %d SQ. MILES OF FARMLAND.\n", g.farmland())
		} else if (area * g.plantingCost) > g.rallods {
			fmt.Printf("   THINK AGAIN. YOU'VE ONLY %d RALLODS LEFT IN THE TREASURY.\n", g.rallods)
		} else {
			return area
		}
	}
}

func (g *Game) getPollutionSpending() int {
	for {
		fmt.Println("HOW MANY RALLRODS DO YOU WISH TO SPEND ON POLLUTION CONTROL? ")
		g.scanner.Scan()
		rallods, err := strconv.Atoi(g.scanner.Text())
		if err != nil || rallods < 0 {
			fmt.Println("INVALID INPUT")
			continue
		}

		if rallods > g.rallods {
			fmt.Printf("   THINK AGAIN. YOU ONLY HAVE %d RALLODS REMAINING.\n", g.rallods)
		} else {
			return rallods
		}
	}
}

func (g *Game) getSellAmount() int {
	displayedFirstError := false
	for {
		fmt.Println("HOW MANY SQUARE MILES DO YOU WISH TO SELL TO INDUSTRY? ")
		g.scanner.Scan()

		amount, err := strconv.Atoi(g.scanner.Text())

		if err != nil || amount < 0 {
			fmt.Println("INVALID INPUT")
			continue
		}

		if amount > g.farmland() {
			if !displayedFirstError {
				fmt.Println("(FOREIGN INDUSTRY WILL ONLY BUY FARM LAND BECAUSE\nFOREST LAND IS UNECONOMICAL TO STRIP MINE DUE TO TREES,\nTHICKER TOP SOIL, ETC.)")
				displayedFirstError = true
			}
			fmt.Printf("***  THINK AGAIN, YOU ONLY HAVE %d SQUARE MILES OF FARM LAND.\n", g.farmland())
		} else {
			return amount
		}
	}
}

func (g *Game) getDistibutionAmount() int {
	for {
		fmt.Println("HOW MANY RALLODS WILL YOU DISTRIBUTE AMONG YOUR COUNTRYMEN? ")
		g.scanner.Scan()

		dole, err := strconv.Atoi(g.scanner.Text())
		if err != nil {
			fmt.Println("INVALID INPUT")
			continue
		}

		if dole < 0 {
			continue
		} else if dole > g.rallods {
			fmt.Printf("   THINK AGAIN. YOU'VE ONLY %d RALLODS IN THE TREASURY.\n", g.rallods)
		} else {
			return dole
		}
	}
}

func printIntro() {
	fmt.Println("                                  KING")
	fmt.Println("               CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY")
	fmt.Println()
	fmt.Println()
}

func printInstructions() {
	fmt.Println("\n\n\nCONGRATULATIONS! YOU'VE JUST BEEN ELECTED PREMIER OF SETATS")
	fmt.Println("DETINU, A SMALL COMMUNIST ISLAND 30 BY 70 MILES LONG. YOUR")
	fmt.Println("JOB IS TO DECIDE UPON THE COUNTRY'S BUDGET AND DISTRIBUTE")
	fmt.Println("MONEY TO YOUR COUNTRYMEN FROM THE COMMUNAL TREASURY.")
	fmt.Printf("THE MONEY SYSTEM IS RALLODS, AND EACH PERSON NEEDS %d\n", COST_OF_LIVING)
	fmt.Println("RALLODS PER YEAR TO SURVIVE. YOUR COUNTRY'S INCOME COMES")
	fmt.Println("FROM FARM PRODUCE AND TOURISTS VISITING YOUR MAGNIFICENT")
	fmt.Println("FORESTS, HUNTING, FISHING, ETC. HALF YOUR LAND IS FARM LAND")
	fmt.Println("WHICH ALSO HAS AN EXCELLENT MINERAL CONTENT AND MAY BE SOLD")
	fmt.Println("TO FOREIGN INDUSTRY (STRIP MINING) WHO IMPORT AND SUPPORT")
	fmt.Println("THEIR OWN WORKERS. CROPS COST BETWEEN 10 AND 15 RALLODS PER")
	fmt.Println("SQUARE MILE TO PLANT.")
	fmt.Println("YOUR GOAL IS TO COMPLETE YOUR {YEARS_IN_TERM} YEAR TERM OF OFFICE.")
	fmt.Println("GOOD LUCK!")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	g := NewGame()

	printIntro()

	showInstructions := getString("DO YOU WANT INSTRUCTIONS? ")

	if showInstructions == "AGAIN" {
		g = resume()
	} else if showInstructions[0:1] != "N" {
		printInstructions()
	}

	for {
		g.setMarket()
		g.printStatus()

		// user actions
		sellAmount := g.getSellAmount()
		g.sellLand(sellAmount)

		dole := g.getDistibutionAmount()
		g.distributeRallods(dole)

		plantArea := g.getAreaToPlant()
		g.plantCrops(plantArea)

		pollutionSpend := g.getPollutionSpending()
		g.spendPollutionControl(pollutionSpend)

		// FA-
		fmt.Println()
		g.handleDeaths(dole, pollutionSpend)

		g.handleForeignWorkers(sellAmount, dole, pollutionSpend)

		g.handleHarvest(plantArea)

		g.handleTourism()

		// FO
		if g.countrymenDeaths > 200 {
			g.handleExcessiveDeaths()
		}

		if g.countrymen < 343 {
			g.handleThirdDied()
		} else if ((g.rallods / 100) > 5) && ((g.countrymenDeaths - g.pollutionDeaths) >= 2) {
			g.handleMoneyMismanagement()
		}

		if g.foreignWorkers > g.countrymen {
			g.handleExcessiveForeigners()
		} else if TERM_LENGTH-1 == g.yearsInOffice {
			g.handleCongratulations()
		} else {
			g.yearsInOffice += 1
			g.countrymenDeaths = 0
		}
	}
}

func getString(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(prompt)
	scanner.Scan()

	return strings.ToUpper(scanner.Text())
}

func getInt(prompt string) int {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println(prompt)
		scanner.Scan()
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("INVALID INPUT, TRY AGAIN")
			continue
		}
		return val
	}
}
