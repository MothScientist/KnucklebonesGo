package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"slices"
	"time"
)

const (
	Reset   = "\033[0m" // color reset
	Green   = "\033[32m"
	Red     = "\033[31m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Cyan    = "\033[36m"
	Magenta = "\033[35m"
)

type player struct {
	dice         [][]int // Current decomposition of bones
	points       int     // Total points (current)
	name         string  // Player name (entered at the beginning of the game)
	deskPosition bool    // Displays in the player console from above or below
}

func main() {
	player0 := player{
		dice: [][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		points:       0,
		name:         "",
		deskPosition: false, // 0 -> we bring out this player's field first
	}

	player1 := player{
		dice: [][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		points:       0,
		name:         "",
		deskPosition: true,
	}

	fmt.Printf("%sName of the 1st player:%s\n", Yellow, Reset)

	_, err := fmt.Scanf("%s\n", &player0.name)
	if err != nil {
		return
	}

	fmt.Printf("%sName of the 2nd player:%s\n", Yellow, Reset)

	_, err = fmt.Scanf("%s\n", &player1.name)
	if err != nil {
		return
	}
	fmt.Print("\n")

	step := true // variable that determines which player is currently moving

	var (
		curPlayer   player // we store a reference to the player who is currently making a move
		otherPlayer player // link to the second player

		number   int // the number we get in each round
		newState int // the number of the column in which the player wants to write the number
	)

	for { // we start an infinite loop, from which we will exit only if diceIsFull() returns true
		clearConsole() // clearing the console

		step = !step // pass the turn to another player

		// we print the players' fields:
		printPlayerFields(&player0)
		fmt.Print("\n")
		printPlayerFields(&player1)

		// take the players' data depending on the move
		if step {
			curPlayer = player1
			otherPlayer = player0
		} else {
			curPlayer = player0
			otherPlayer = player1
		}

		fmt.Print("\n")

		// print the name of the player who moves and his new number:
		fmt.Printf("%s%s, your move!%s", Green, curPlayer.name, Reset)
		fmt.Print("\n")
		number = getRandomDice()
		fmt.Printf("%sYour number: %d%s", Red, number, Reset)
		fmt.Print("\n")
		availableColumns := curPlayer.getAvailableColumns()
		for {
			fmt.Print("In which column should this cube be placed [1-3]: ")
			_, err := fmt.Fscan(os.Stdin, &newState)
			if err != nil {
				return
			}
			newState -= 1 // array starts from zero
			if slices.Contains(availableColumns, newState) {
				curPlayer.dice[newState][getIndexToInsertDice(curPlayer.dice[newState])] = number
				break
			} else if newState < 1 || newState > 3 {
				fmt.Println("Such column does not exist")
			} else {
				fmt.Println("This column is already filled in.")
			}
		}

		// we look to see if the user's bones are "knocked down"
		otherPlayer.reCalcDice(newState, curPlayer.dice[newState])

		// we recalculate points for both players after the move and "knocking out" the opponent's dice
		player0.calcPoints()
		player1.calcPoints()

		// look to see if one of the players has a filled field
		diceIsFull := player0.diceIsFull()
		if !diceIsFull {
			diceIsFull = player1.diceIsFull()
		}

		if diceIsFull {
			break // if someone's field is filled, then we exit the game
		}

		// if necessary, we move the numbers down the board
		player0.dropDiceNumbers()
		player1.dropDiceNumbers()

	}

	clearConsole() // clearing the console before announcing the winner

	fmt.Printf("%s%s - %d\n", Blue, player0.name, player0.points)
	fmt.Printf("%s%s - %d%s\n", Blue, player1.name, player1.points, Reset)

	if player0.points > player1.points {
		fmt.Printf("%sWinner: %s%s\n", Green, player0.name, Reset)
	} else if player1.points > player0.points {
		fmt.Printf("%sWinner: %s%s\n", Green, player1.name, Reset)
	} else {
		fmt.Printf("%sDraw!%s\n", Green, Reset)
	}
}

// calcPoints player score calculation
func (p *player) calcPoints() {
	points := 0
	for column := range p.dice {
		sliceCopy := getUniqueElements(p.dice[column]) // necessary to eliminate repetitions in order to correctly calculate the multipliers
		for _, value := range sliceCopy {
			factor := countValuesInArray(value, p.dice[column])
			if factor == 3 {
				points += (value * 3) * 3
			} else if factor == 2 {
				points += (value * 2) * 2
			} else {
				points += value
			}
		}
	}
	p.points = points
}

// reCalcDice recalculates dice after opponent's move
func (p *player) reCalcDice(column int, opponentColumnDice []int) {
	for idx, value := range p.dice[column] {
		if slices.Contains(opponentColumnDice, value) {
			p.dice[column][idx] = 0
		}
	}
}

// getRandomDice gets a random number 1-6 (inclusive)
func getRandomDice() int {
	// update seed before use - otherwise everything comes down to one number
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(6) + 1
}

// diceIsFull determines that the game is over -> the player has all fields inside the dice filled after the move
func (p *player) diceIsFull() bool {
	for column := range p.dice {
		for row := range p.dice[column] {
			if p.dice[column][row] == 0 {
				return false
			}
		}
	}
	return true
}

// getAvailableColumns get a list of writable columns
func (p *player) getAvailableColumns() []int {
	var res []int // we write the available columns into this slice
	for column := range p.dice {
		for row := range p.dice[column] {
			if p.dice[column][row] == 0 {
				res = append(res, column)
				break
			}
		}
	}
	return res
}

// getIndexToInsertDice returns the position to write to the string - (0, 1, 2)
func getIndexToInsertDice(dice []int) int {
	var res int
	for idx, value := range dice {
		if value == 0 {
			res = idx
		}
	}
	return res
}

// clearConsole completely clears the console
func clearConsole() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	err := c.Run()
	if err != nil {
		return
	}
}

// countValuesInArray returns the number of occurrences of a number in an array
func countValuesInArray(value int, arr []int) int {
	res := 0
	for _, arrVal := range arr {
		if arrVal == value {
			res += 1
		}
	}
	return res
}

// getUniqueElements return slice without duplicate elements
func getUniqueElements[T comparable](arr []T) []T {
	m, uniq := make(map[T]struct{}), make([]T, 0, len(arr))
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			m[v], uniq = struct{}{}, append(uniq, v)
		}
	}
	return uniq
}

// printPlayerFields prints the player's field according to his parameters to the console
func printPlayerFields(p *player) {
	fmt.Printf("%s%s%s\n", Magenta, p.name, Cyan)
	for column := range p.dice {
		for row := range p.dice[column] {
			var state int
			if p.deskPosition {
				// "2 - column": since we output the matrix of the second player
				// in the direction to the matrix of the first player
				state = p.dice[row][2-column]
			} else {
				state = p.dice[row][column]
			}
			if state != 0 {
				fmt.Printf("%d", state)
			} else {
				fmt.Print(" ")
			}
			fmt.Print(" ")
		}
		fmt.Print("\n")
	}
	fmt.Printf("%s%sScore: %d%s\n", Reset, Magenta, p.points, Reset)
}

// dropDiceNumbers replaces empty bottom cells with dice
func (p *player) dropDiceNumbers() {
	for i := 0; i <= 2; i++ {
		p.dice[i] = removeZeros(p.dice[i])
	}
}

// removeZeros removes zeros from a slice
func removeZeros(column []int) []int {
	var nonZeros []int

	for _, value := range column {
		if value != 0 {
			nonZeros = append(nonZeros, value)
		}
	}

	if len(nonZeros) == 0 {
		return column
	} else {
		if len(nonZeros) < 3 {
			nonZeros = append(make([]int, 3-len(nonZeros)), nonZeros...) // Fill with zeros at the beginning
		}
		return nonZeros
	}
}
