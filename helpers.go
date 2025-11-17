package main

import (
	"os"
	"os/exec"
	"math/rand"
	"time"
)

// getRandomDice gets a random number 1-6 (inclusive)
func getRandomDice() int {
	// update seed before use - otherwise everything comes down to one number
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(6) + 1
}

// getIndexToInsertDice returns the position to write to the string - (0, 1, 2)
func getIndexToInsertDice(dice []int) (res int) {
	for idx, value := range dice {
		if value == 0 {
			res = idx
		}
	}
	return
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