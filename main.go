package main

import (
	"fmt"
	"os"
	"slices"
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

func main() {
	player0  := &player{
		dice: [][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		points:       0,
		name:         "",
		deskPosition: false, // 0 -> we bring out this player's field first
	}

	player1 := &player{
		dice: [][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		points:       0,
		name:         "",
		deskPosition: true,
	}

	clearConsole()

	fmt.Printf("%sName of the 1st player:%s\n", Yellow, Reset)
	player0.scanNameOfPlayer()

	fmt.Printf("%sName of the 2nd player:%s\n", Yellow, Reset)
	player1.scanNameOfPlayer()

	fmt.Print("\n")

	step := true // variable that determines which player is currently moving

	var (
		curPlayer   *player // we store a reference to the player who is currently making a move
		otherPlayer *player // link to the second player

		number   int // the number we get in each round
		newState int // the number of the column in which the player wants to write the number
	)

	for { // we start an infinite loop, from which we will exit only if diceIsFull() returns true
		clearConsole() // clearing the console

		step = !step // pass the turn to another player

		// we print the players' fields:
		player0.printPlayerFields()
		fmt.Print("\n")
		player1.printPlayerFields()

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

			if newState < 1 || newState > 3 {
				fmt.Println("Such column does not exist")
				continue
			}

			newState -= 1 // array starts from zero
			if slices.Contains(availableColumns, newState) {
				indexToInsert := getIndexToInsertDice(curPlayer.dice[newState])
				curPlayer.dice[newState][indexToInsert] = number
				break
			} else {
				fmt.Println("This column is already filled in.")
				continue
			}
		}

		// look to see if the user's bones are "knocked down"
		otherPlayer.reCalcDice(newState, curPlayer.dice[newState])

		// recalculate points for both players after the move and "knocking out" the opponent's dice
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
