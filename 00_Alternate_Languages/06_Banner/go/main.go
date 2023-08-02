package main

import (
	"fmt"
	"math"
	"strings"
)

var letters map[string][]int = map[string][]int{
	" ": {0, 0, 0, 0, 0, 0, 0},
	"A": {505, 37, 35, 34, 35, 37, 505},
	"G": {125, 131, 258, 258, 290, 163, 101},
	"E": {512, 274, 274, 274, 274, 258, 258},
	"T": {2, 2, 2, 512, 2, 2, 2},
	"W": {256, 257, 129, 65, 129, 257, 256},
	"L": {512, 257, 257, 257, 257, 257, 257},
	"S": {69, 139, 274, 274, 274, 163, 69},
	"O": {125, 131, 258, 258, 258, 131, 125},
	"N": {512, 7, 9, 17, 33, 193, 512},
	"F": {512, 18, 18, 18, 18, 2, 2},
	"K": {512, 17, 17, 41, 69, 131, 258},
	"B": {512, 274, 274, 274, 274, 274, 239},
	"D": {512, 258, 258, 258, 258, 131, 125},
	"H": {512, 17, 17, 17, 17, 17, 512},
	"M": {512, 7, 13, 25, 13, 7, 512},
	"?": {5, 3, 2, 354, 18, 11, 5},
	"U": {128, 129, 257, 257, 257, 129, 128},
	"R": {512, 18, 18, 50, 82, 146, 271},
	"P": {512, 18, 18, 18, 18, 18, 15},
	"Q": {125, 131, 258, 258, 322, 131, 381},
	"Y": {8, 9, 17, 481, 17, 9, 8},
	"V": {64, 65, 129, 257, 129, 65, 64},
	"X": {388, 69, 41, 17, 41, 69, 388},
	"Z": {386, 322, 290, 274, 266, 262, 260},
	"I": {258, 258, 258, 512, 258, 258, 258},
	"C": {125, 131, 258, 258, 258, 131, 69},
	"J": {65, 129, 257, 257, 257, 129, 128},
	"1": {0, 0, 261, 259, 512, 257, 257},
	"2": {261, 387, 322, 290, 274, 267, 261},
	"*": {69, 41, 17, 512, 17, 41, 69},
	"3": {66, 130, 258, 274, 266, 150, 100},
	"4": {33, 49, 41, 37, 35, 512, 33},
	"5": {160, 274, 274, 274, 274, 274, 226},
	"6": {194, 291, 293, 297, 305, 289, 193},
	"7": {258, 130, 66, 34, 18, 10, 8},
	"8": {69, 171, 274, 274, 274, 171, 69},
	"9": {263, 138, 74, 42, 26, 10, 7},
	"=": {41, 41, 41, 41, 41, 41, 41},
	"!": {1, 1, 1, 384, 1, 1, 1},
	"0": {57, 69, 131, 258, 131, 69, 57},
	".": {1, 1, 129, 449, 129, 1, 1},
}

func printBanner(horizontal, vertical, gl int, character, statement string) {
	f := make([]int, 7)
	j := make([]int, 9)

	for _, v := range statement {
		s := letters[string(v)]
		x_str := character
		if character == "ALL" {
			x_str = string(v)
		}
		if x_str == " " {
			for i := 0; i < 7*horizontal; i++ {
				fmt.Println()
			}
		} else {
			for u := 0; u < 7; u++ {
				for k := 8; k > -1; k-- {
					if int(math.Pow(2, float64(k))) >= s[u] {
						j[8-k] = 0
					} else {
						j[8-k] = 1
						s[u] = s[u] - int(math.Pow(2, float64(k)))
						if s[u] == 1 {
							f[u] = 8 - k
							break
						}
					}
				}
				for _tl := 1; _tl < (horizontal + 1); _tl++ {
					lineStr := getRepeatedString(" ", int((float64(63)-float64(4.5)*float64(vertical))*float64(gl)/float64(len(x_str))+1))

					for b := 0; b < f[u]+1; b++ {
						if j[b] == 0 {
							for z := 1; z < vertical+1; z++ {
								lineStr = lineStr + getRepeatedString(" ", len(x_str))
							}
						} else {
							lineStr = lineStr + getRepeatedString(x_str, vertical)
						}
					}
					fmt.Println(lineStr)
				}
			}
			fmt.Println(getRepeatedString("\n", 2*horizontal+1))
		}
	}
}

func getRepeatedString(s string, length int) string {
	ret := strings.Builder{}
	for x := 0; x < length; x++ {
		ret.WriteString(s)
	}
	return ret.String()
}

func main() {
	/*scanner := bufio.NewScanner(os.Stdin)
	vertical := 0
	horizontal := 0
	var err error

	for {
		fmt.Println("Horizontal ")
		scanner.Scan()
		horizontal, err = strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Please enter a number")
		} else if horizontal < 1 {
			fmt.Println("Horizontal must be greater than zero")
		} else {
			break
		}
	}
	for {
		fmt.Println("Vertical ")
		scanner.Scan()
		vertical, err = strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Please enter a number")
		} else if vertical < 1 {
			fmt.Println("Vertical must be greater than zero")
		} else {
			break
		}
	}

	gl := 0
	fmt.Println("Centered ")
	scanner.Scan()
	if strings.ToLower(scanner.Text()[0:1]) == "y" {
		gl = 1
	}

	fmt.Println("Character (type 'ALL' if you want character being printed) ")
	scanner.Scan()
	character := strings.ToUpper(scanner.Text())

	fmt.Println("Statement ")
	scanner.Scan()
	statement := strings.ToUpper(scanner.Text())

	// this means to prepare printer, just press Enter
	fmt.Println("Set page ")
	scanner.Scan()

	printBanner(horizontal, vertical, gl, character, statement)*/
	printBanner(2, 0, 0, "ALL", "HELLO THERE")
}
