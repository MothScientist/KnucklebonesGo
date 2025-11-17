package main

import (
	"fmt"
	"slices"
)

type player struct {
	dice         [][]int // Current decomposition of bones
	points       int     // Total points (current)
	name         string  // Player name (entered at the beginning of the game)
	deskPosition bool    // Displays in the player console from above or below
}

// calcPoints player score calculation
func (p *player) calcPoints() {
	points := 0
	for column := range p.dice {
		sliceCopy := getUniqueElements(p.dice[column]) // necessary to eliminate repetitions in order to correctly calculate the multipliers
		for _, value := range sliceCopy {
			factor := countValuesInArray(value, p.dice[column])
			switch factor {
			case 3:
				points += (value * 3) * 3
			case 2:
				points += (value * 2) * 2
			case 1:
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

// diceIsFull determines that the game is over -> the player has all fields inside the dice filled after the move
func (p *player) diceIsFull() bool {
	for column := range p.dice {
		if slices.Contains(p.dice[column], 0) {
			return false
		}
	}
	return true
}

// getAvailableColumns get a list of writable columns
func (p *player) getAvailableColumns() (res []int) {
	for column := range p.dice {
		if slices.Contains(p.dice[column], 0) {
			res = append(res, column)
		}
	}
	return
}

// printPlayerFields prints the player's field according to his parameters to the console
func (p *player) printPlayerFields() {
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

func (p *player) scanNameOfPlayer() {
	_, err := fmt.Scanf("%s\n", &p.name)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
}