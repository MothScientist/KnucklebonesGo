package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"slices"
	"sync"
	"time"
)

type player struct {
	dice          [][]int // Current decomposition of bones
	points        int     // Total points (current)
	name          string  // Player name (entered at the beginning of the game)
	desk_position bool    // Displays in the player console from above or below
}

func main() {
	// entry point

	player_0 := player{
		dice: [][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		points:        0,
		name:          "",
		desk_position: false, // 0 -> we bring out this player's field first
	}

	player_1 := player{
		dice: [][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		points:        0,
		name:          "",
		desk_position: true, // 1
	}

	fmt.Println("Name of 1st player:")
	fmt.Scanf("%s\n", &player_0.name)

	fmt.Println("Name of the 2nd player:")
	fmt.Scanf("%s\n", &player_1.name)
	fmt.Print("\n")

	step := true // variable that determines which player is currently moving

	var cur_player player   // we store a reference to the player who is currently making a move
	var other_player player // link to the second player

	var number int // the number we get in each round

	var new_state int // the number of the column in which the player wants to write the number

	var wg sync.WaitGroup // counter variable for waiting for goroutines

	for { // we start an infinite loop, from which we will exit only if diceIsFull() returns true
		clearConsole() // clearing the console

		step = !step // pass the turn to another player

		// we print the players' fields:
		printPlayerFields(&player_0)
		fmt.Print("\n\n\n")
		printPlayerFields(&player_1)

		// take the players' data depending on the move
		if step {
			cur_player = player_1
			other_player = player_0
		} else {
			cur_player = player_0
			other_player = player_1
		}

		fmt.Print("\n\n\n")

		// print the name of the player who moves and his new number:
		fmt.Printf("%s, your move!", cur_player.name)
		fmt.Print("\n")
		number = getRandomDice()
		fmt.Printf("Your number: %d", number)
		fmt.Print("\n")
		availableColumns := cur_player.getAvailableColumns()
		for {
			fmt.Print("In which column should this cube be placed [1-3]: ")
			fmt.Fscan(os.Stdin, &new_state)
			new_state -= 1 // array starts from zero
			if slices.Contains(availableColumns, new_state) {
				cur_player.dice[new_state][getIndexToInsertDice(cur_player.dice[new_state])] = number
				break
			} else if new_state < 1 || new_state > 3 {
				fmt.Println("Such column does not exist")
			} else {
				fmt.Println("This column is already filled in.")
			}
		}

		// we look to see if the user's bones are "knocked down"
		other_player.reCalcDice(new_state, cur_player.dice[new_state])

		// we recalculate points for both players after the move and "knocking out" the opponent's dice
		wg.Add(4) // launch 4 goroutines
		go player_0.calcPoints(&wg)
		go player_1.calcPoints(&wg)

		ch_dice_is_full := make(chan bool, 2)
		dice_is_full := false
		// look to see if one of the players has a filled field
		go player_0.diceIsFull(&wg, ch_dice_is_full)
		go player_1.diceIsFull(&wg, ch_dice_is_full)
		wg.Wait() // wait for all goroutines to complete (zeroing the counter)

		// read results from channel and update boolean variable
		for i := 0; i <= 1; i++ {
			if res := <- ch_dice_is_full; res {
				dice_is_full = true
		}
	}

		if dice_is_full {
			break // if someone's field is filled, then we exit the game
		}

	}

	clearConsole() // clearing the console before announcing the winner

	fmt.Printf("%s - %d\n", player_0.name, player_0.points)
	fmt.Printf("%s - %d\n", player_1.name, player_1.points)

	if player_0.points > player_1.points {
		fmt.Printf("Winner: %s\n", player_0.name)
	} else if player_1.points > player_0.points {
		fmt.Printf("Winner: %s\n", player_1.name)
	} else {
		fmt.Println("Draw!")
	}
}

func (p *player) calcPoints(wg *sync.WaitGroup) {
	defer wg.Done() // after goroutine execution we subtract one from the counter

	points := 0

	for column := range p.dice {
		slice_copy := getUniqueElements(p.dice[column]) // necessary to eliminate repetitions in order to correctly calculate the multipliers
		for _, value := range slice_copy {
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
	c.Run()
}

// Returns the number of occurrences of a number in an array
func count(value int, arr []int) int {
	res := 0
	for _, arr_val := range arr {
		if arr_val == value {
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
				if p.desk_position {
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
