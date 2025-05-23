package main

import (
	"fmt"
	"golang.org/x/exp/slices" // crutch for versions before 1.21, since then the slices module is built-in
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"
)

type player struct {
	dice         [][]int // Current decomposition of bones
	points       int     // Total points (current)
	name         string  // Player name (entered at the beginning of the game)
	deskPosition bool    // Displays in the player console from above or below
}

func main() {
	// entry point

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
		deskPosition: true, // 1
	}

	fmt.Println("Name of 1st player:")
	_, err := fmt.Scanf("%s\n", &player0.name)
	if err != nil {
		return
	}

	fmt.Println("Name of the 2nd player:")
	_, err = fmt.Scanf("%s\n", &player1.name)
	if err != nil {
		return
	}
	fmt.Print("\n")

	step := true // variable that determines which player is currently moving

	var curPlayer player   // we store a reference to the player who is currently making a move
	var otherPlayer player // link to the second player

	var number int // the number we get in each round

	var newState int // the number of the column in which the player wants to write the number

	var wg sync.WaitGroup // counter variable for waiting for goroutines

	for { // we start an infinite loop, from which we will exit only if diceIsFull() returns true
		clearConsole() // clearing the console

		step = !step // pass the turn to another player

		// we print the players' fields:
		printPlayerFields(&player0)
		fmt.Print("\n\n\n")
		printPlayerFields(&player1)

		// take the players' data depending on the move
		if step {
			curPlayer = player1
			otherPlayer = player0
		} else {
			curPlayer = player0
			otherPlayer = player1
		}

		fmt.Print("\n\n\n")

		// print the name of the player who moves and his new number:
		fmt.Printf("%s, your move!", curPlayer.name)
		fmt.Print("\n")
		number = getRandomDice()
		fmt.Printf("Your number: %d", number)
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
		wg.Add(4) // launch 4 goroutines
		go player0.calcPoints(&wg)
		go player1.calcPoints(&wg)

		chDiceIsFull := make(chan bool, 2)
		DiceIsFull := false
		// look to see if one of the players has a filled field
		go player0.diceIsFull(&wg, chDiceIsFull)
		go player1.diceIsFull(&wg, chDiceIsFull)
		wg.Wait() // wait for all goroutines to complete (zeroing the counter)

		// read results from channel and update boolean variable
		for i := 0; i <= 1; i++ {
			if res := <-chDiceIsFull; res {
				DiceIsFull = true
			}
		}

		if DiceIsFull {
			break // if someone's field is filled, then we exit the game
		}

	}

	clearConsole() // clearing the console before announcing the winner

	fmt.Printf("%s - %d\n", player0.name, player0.points)
	fmt.Printf("%s - %d\n", player1.name, player1.points)

	if player0.points > player1.points {
		fmt.Printf("Winner: %s\n", player0.name)
	} else if player1.points > player0.points {
		fmt.Printf("Winner: %s\n", player1.name)
	} else {
		fmt.Println("Draw!")
	}
}

func (p *player) calcPoints(wg *sync.WaitGroup) {
	defer wg.Done() // after goroutine execution we subtract one from the counter

	points := 0

	for column := range p.dice {
		sliceCopy := getUniqueElements(p.dice[column]) // necessary to eliminate repetitions in order to correctly calculate the multipliers
		for _, value := range sliceCopy {
			factor := count(value, p.dice[column])
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

// Recalculates dice after opponent's move
func (p *player) reCalcDice(column int, opponentColumnDice []int) {
	for idx, value := range p.dice[column] {
		if slices.Contains(opponentColumnDice, value) {
			p.dice[column][idx] = 0
		}
	}
}

// Gets a random number 1-6 (inclusive)
func getRandomDice() int {
	// update seed before use - otherwise everything comes down to one number
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(6) + 1
}

// Determines that the game is over -> the player has all fields inside the dice filled after the move
func (p *player) diceIsFull(wg *sync.WaitGroup, ch chan bool) {
	defer wg.Done()

	for column := range p.dice {
		for row := range p.dice[column] {
			if p.dice[column][row] == 0 {
				ch <- false
				return
			}
		}
	}
	ch <- true
}

// Get a list of writable columns
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

// Returns the position to write to the string - (0, 1, 2)
func getIndexToInsertDice(dice []int) int {
	var res int
	for idx, value := range dice {
		if value == 0 {
			res = idx
		}
	}
	return res
}

// Console clearing function
func clearConsole() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	err := c.Run()
	if err != nil {
		return
	}
}

// Returns the number of occurrences of a number in an array
func count(value int, arr []int) int {
	res := 0
	for _, arrVal := range arr {
		if arrVal == value {
			res += 1
		}
	}
	return res
}

// Return slice without duplicate elements
func getUniqueElements[T comparable](arr []T) []T {
	m, uniq := make(map[T]struct{}), make([]T, 0, len(arr))
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			m[v], uniq = struct{}{}, append(uniq, v)
		}
	}
	return uniq
}

func printPlayerFields(p *player) {
	fmt.Println(p.name)
	for column := range p.dice {
		for row := range p.dice[column] {
			var state int
			if p.deskPosition {
				state = p.dice[row][2-column] // "2 - column" it is necessary, since we output the matrix of the second player in the direction to the matrix of the first player
			} else {
				state = p.dice[row][column]
			}
			if state != 0 {
				fmt.Print(state)
			} else {
				fmt.Print(" ")
			}
			fmt.Print(" ")
		}
		fmt.Print("\n")
	}

	fmt.Print("\n")
	fmt.Printf("Score: %d", p.points)
}
