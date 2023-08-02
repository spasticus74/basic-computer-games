package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type game struct {
	time     int
	score    [2]int //first value is opponents score, second is home
	shot     int
	z1       float64
	defense  float64
	opponent string
}

func NewGame() *game {
	g := game{time: 0, score: [2]int{0, 0}}
	g.defense = getDefenseChoice()
	g.opponent = getOpponentsName()
	return &g
}

// Print the curent score
func (g *game) printScore() {
	fmt.Printf("Score:  %d to %d\n", g.score[1], g.score[0])
}

// Add points to the score. Team can take 0 or 1, for opponent or Dartmouth, respectively
func (g *game) addPoints(team, points int) {
	if team < 0 || team > 1 {
		log.Fatal("invalid value for 'team'")
	}
	if points < 0 || points > 2 {
		log.Fatal("invalid value for 'points'")
	}

	g.score[team] += points
	g.printScore()
}

// Called when t = 92
func (g *game) twoMinuteWarning() {
	fmt.Println("   *** Two minutes left in the game ***")
}

func (g *game) getDartmouthBallChoice() {
	var shot int
	var err error
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nYour shot? [1-4, 0 to change defense]")
		scanner.Scan()

		shot, err = strconv.Atoi(scanner.Text())
		if (err == nil) && (shot >= 0) && (shot <= 4) {
			break
		}
	}
	g.shot = shot
}

// Called when the user enters 1 or 2 for their shot
func (g *game) dartmouthJumpShot() {
	g.time += 1

	if g.time == 50 {
		g.halftime()
	} else if g.time == 92 {
		g.twoMinuteWarning()
	}

	fmt.Println("Jump Shot.")

	// simulates chances of different possible outcomes
	if rand.Float64() > 0.341*g.defense/8 {
		if rand.Float64() > 0.682*g.defense/8 {
			if rand.Float64() > 0.782*g.defense/8 {
				if rand.Float64() > 0.843*g.defense/8 {
					print("Charging foul. Dartmouth loses ball.\n")
					g.opponentBall()
				} else {
					// player is fouled
					g.foulShots(1)
					g.opponentBall()
				}
			} else {
				if rand.Float64() > 0.5 {
					fmt.Printf("Shot is blocked. Ball controlled by %s\n", g.opponent)
					g.opponentBall()
				} else {
					fmt.Println("Shot is blocked. Ball controlled by Dartmouth.")
					g.dartmouthBall()
				}
			}
		} else {
			fmt.Println("Shot is off target.")
			if (g.defense / 6 * rand.Float64()) > 0.45 {
				print("Rebound to " + g.opponent + "\n")
				g.opponentBall()
			} else {
				fmt.Println("Dartmouth controls the rebound.")
				if rand.Float64() > 0.4 {
					if g.defense == 6 && rand.Float64() > 0.6 {
						fmt.Printf("Pass stolen by %s, easy lay up\n", g.opponent)
						g.addPoints(0, 2)
						g.dartmouthBall()
					} else {
						// ball is passed back to you
						g.ballPassedBack()
					}
				} else {
					print()
					g.dartmouthNonJumpShot()
				}
			}
		}
	} else {
		fmt.Println("Shot is good.")
		g.addPoints(1, 2)
		g.opponentBall()
	}

}

// Lay up, set shot, or defense change, called when the user enters 0, 3, or 4
func (g *game) dartmouthNonJumpShot() {
	g.time += 1
	if g.time == 50 {
		g.halftime()
	} else if g.time == 92 {
		g.twoMinuteWarning()
	}

	if g.shot == 4 {
		fmt.Println("Set shot.")
	} else if g.shot == 3 {
		fmt.Println("Lay up.")
	} else if g.shot == 0 {
		g.changeDefense()
	}

	// simulates different outcomes after a lay up or set shot
	if 7/g.defense*rand.Float64() > 0.4 {
		if 7/g.defense*rand.Float64() > 0.7 {
			if 7/g.defense*rand.Float64() > 0.875 {
				if 7/g.defense*rand.Float64() > 0.925 {
					fmt.Println("Charging foul. Dartmouth loses the ball.")
					g.opponentBall()
				} else {
					fmt.Printf("Shot blocked. %s's ball.\n", g.opponent)
					g.opponentBall()
				}
			} else {
				g.foulShots(1)
				g.opponentBall()
			}
		} else {
			fmt.Println("Shot is off the rim.")
			if rand.Float64() > (2.0 / 3.0) {
				fmt.Println("Dartmouth controls the rebound.")
				if rand.Float64() > 0.4 {
					fmt.Println("Ball passed back to you.")
					g.dartmouthBall()
				} else {
					g.dartmouthNonJumpShot()
				}
			} else {
				fmt.Printf("%s controls the rebound.\n", g.opponent)
				g.opponentBall()
			}
		}
	} else {
		fmt.Println("Shot is good. Two points.")
		g.addPoints(1, 2)
		g.opponentBall()
	}
}

// Plays out a Dartmouth posession, starting with your choice of shot
func (g *game) dartmouthBall() {
	g.getDartmouthBallChoice()

	if g.time < 100 || rand.Float32() < 0.5 {
		if g.shot == 1 || g.shot == 2 {
			g.dartmouthJumpShot()
		} else {
			g.dartmouthNonJumpShot()
		}
	} else {
		if g.score[0] != g.score[1] {
			fmt.Println("\n   ***** End Of Game *****")
			fmt.Printf("Final score: Dartmouth: %d %s: %d\n", g.score[1], g.opponent, g.score[0])
			return
		} else {
			fmt.Println("\n   ***** End Of Second Half *****")
			fmt.Println("Score at end of regulation time:")
			fmt.Printf("     Dartmouth: %d %s: %d\n", g.score[1], g.opponent, g.score[0])
			fmt.Println("Begin two minute overtime period")
			g.time = 93
			g.startPeriod()
		}
	}
}

func (g *game) ballPassedBack() {
	fmt.Println("Ball passed back to you.")
	g.dartmouthBall()
}

func (g *game) changeDefense() {
	var defense float64
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Set defensive alignment to? ")
		scanner.Scan()
		defense, err = strconv.ParseFloat(scanner.Text(), 32)
		if (err == nil) && ((defense == 6.0) || (defense == 6.5) || (defense == 7.0) || (defense == 7.5)) {
			break
		}
	}

	g.defense = defense
	g.dartmouthBall()
}

// Simulate two foul shots for a player and adds the points.
func (g *game) foulShots(team int) {
	fmt.Println("Shooter fouled. Two shots.")

	if rand.Float32() > 0.49 {
		if rand.Float32() > 0.75 {
			fmt.Println("Both shots missed")
		} else {
			fmt.Println("Shooter makes one shot and misses one.")
			g.addPoints(team, 1)
		}
	} else {
		fmt.Println("Shooter makes both shots.")
		g.addPoints(team, 2)
	}
}

// Called when t = 50, starts new period
func (g *game) halftime() {
	fmt.Println("\n   ***** End of first half *****")
	g.printScore()
	g.startPeriod()
}

// Simulate a center jump for posession at the beginning of a period
func (g *game) startPeriod() {
	fmt.Println("Center jump")
	if rand.Float32() > 0.6 {
		fmt.Println("Dartmouth controls the tap.")
		g.dartmouthBall()
	} else {
		fmt.Printf("%s controls the tap.\n", g.opponent)
		g.opponentBall()
	}
}

// Simulate opponents lay up or set shot.
func (g *game) opponentNonJumpshot() {
	if g.z1 > 3 {
		fmt.Println("Set shot.")
	} else {
		fmt.Println("Lay up")
	}

	if 7/g.defense*rand.Float64() > 0.413 {
		fmt.Println("Shot is missed.")
		if g.defense/6*rand.Float64() > 0.5 {
			fmt.Printf("%s controls the rebound.\n", g.opponent)
			if g.defense == 6 {
				if rand.Float64() > 0.75 {
					fmt.Println("Ball stolen. Easy lay up for Dartmouth.")
					g.addPoints(1, 2)
					g.opponentBall()
				} else {
					if rand.Float64() > 0.5 {
						fmt.Println()
						g.opponentNonJumpshot()
					} else {
						fmt.Printf("Pass back to %s guard.\n", g.opponent)
						g.opponentBall()
					}
				}
			} else {
				if rand.Float64() > 0.5 {
					fmt.Println()
					g.opponentNonJumpshot()
				} else {
					fmt.Printf("Pass back to %s guard.\n", g.opponent)
					g.opponentBall()
				}
			}
		} else {
			fmt.Println("Dartmouth controls the rebound.")
			g.dartmouthBall()
		}
	} else {
		fmt.Println("Shot is good.")
		g.addPoints(0, 2)
		g.dartmouthBall()
	}
}

// Simulate the opponents jumpshot
func (g *game) opponentJumpshot() {
	fmt.Println("Jump Shot.")

	if 8/g.defense*rand.Float64() > 0.35 {
		if 8/g.defense*rand.Float64() > 0.75 {
			if 8/g.defense*rand.Float64() > 0.9 {
				fmt.Println("Offensive foul. Dartmouth's ball.")
				g.dartmouthBall()
			} else {
				g.foulShots(0)
				g.dartmouthBall()
			}
		} else {
			fmt.Println("Shot is off the rim.")
			if g.defense/6*rand.Float64() > 0.5 {
				fmt.Printf("%s controls the rebound.\n", g.opponent)
				if g.defense == 6 {
					if rand.Float64() > 0.75 {
						fmt.Println("Ball stolen. Easy layup for Dartmouth.")
						g.addPoints(1, 2)
						g.opponentBall()
					} else {
						if rand.Float64() > 0.5 {
							fmt.Println()
							g.opponentNonJumpshot()
						} else {
							fmt.Printf("Pass back to %s guard.\n", g.opponent)
							g.opponentBall()
						}
					}
				} else {
					if rand.Float64() > 0.5 {
						g.opponentNonJumpshot()
					} else {
						fmt.Printf("Pass back to %s guard.\n", g.opponent)
						g.opponentBall()
					}
				}
			} else {
				fmt.Println("Dartmouth controls the rebound")
				fmt.Println()
				g.dartmouthBall()
			}
		}
	} else {
		fmt.Println("Shot is good.")
		g.addPoints(0, 2)
		g.dartmouthBall()
	}
}

// Simulate an opponents possesion. Randomly picks jump shot or lay up / set shot.
func (g *game) opponentBall() {
	g.time += 1
	if g.time == 50 {
		g.halftime()
	}
	g.z1 = 10/4*rand.Float64() + 1
	if g.z1 > 2 {
		g.opponentNonJumpshot()
	} else {
		g.opponentJumpshot()
	}
}

func printIntro() {
	fmt.Println("\t\t\t Basketball")
	fmt.Printf("\t Creative Computing  Morristown, New Jersey\n\n\n\n")
	fmt.Println("This is Dartmouth College basketball. ")
	fmt.Println("Υou will be Dartmouth captain and playmaker.")
	fmt.Println("Call shots as follows:")
	fmt.Println("1. Long (30ft.) Jump Shot;\n2. Short (15 ft.) Jump Shot;\n3. Lay up; 4. Set Shot")
	fmt.Println("Both teams will use the same defense. Call Defense as follows:")
	fmt.Println("6. Press; 6.5 Man-to-Man; 7. Zone; 7.5 None.")
	fmt.Println("To change defense, just type 0 as your next shot.")
}

func getOpponentsName() string {
	fmt.Println("\nChoose your opponent? ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func getDefenseChoice() float64 {
	var defense float64
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Set defensive alignment to? ")
		scanner.Scan()
		defense, err = strconv.ParseFloat(scanner.Text(), 32)
		if (err == nil) && ((defense == 6.0) || (defense == 6.5) || (defense == 7.0) || (defense == 7.5)) {
			break
		}
	}

	return defense
}
func main() {
	rand.Seed(time.Now().UnixNano())
	printIntro()
	g := NewGame()

	g.startPeriod()
}
