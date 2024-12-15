package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	file, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err) // Terminate if the input file cannot be read, as further processing is impossible
	}

	lines := strings.Split(string(file), "\n")

	isRulesDone := false
	rules := make([][2]int, 0)
	sum := 0
	for _, line := range lines {
		if line == "" {
			isRulesDone = true
			continue
		}

		if !isRulesDone {
			rules = append(rules, parseRule(line))
			continue
		}

		sequenceStrings := strings.Split(line, ",")
		sequence := make([]int, len(sequenceStrings))
		for i, s := range sequenceStrings {
			sequence[i] = parseInt(s)
		}

		// PART 1
		// if !isValidSequence(sequence, rules) {
		// 	continue
		// }

		// PART 2
		if isValidSequence(sequence, rules) {
			continue
		}

		sequence = fixSequence(sequence, rules)

		middleDigit := sequence[len(sequence)/2]
		sum += middleDigit

	}

	fmt.Println(sum)
}

func fixSequence(sequence []int, rules [][2]int) []int {
	for !isValidSequence(sequence, rules) {
		for _, r := range rules {
			if !containsAll(sequence, r[0], r[1]) {
				continue
			}

			indexOfFirst := indexOf(sequence, r[0])
			indexOfSecond := indexOf(sequence, r[1])
			if indexOfFirst > indexOfSecond {
				sequence[indexOfFirst], sequence[indexOfSecond] = sequence[indexOfSecond], sequence[indexOfFirst]
			}
		}
	}

	return sequence
}

func indexOf(sequence []int, i int) int {
	for index, s := range sequence {
		if s == i {
			return index
		}
	}

	return -1
}

func isValidSequence(sequence []int, rule [][2]int) bool {
	for _, r := range rule {
		if !containsAll(sequence, r[0], r[1]) {
			continue
		}

		firstFound, secondFound := false, false
		for _, s := range sequence {
			if s == r[0] {
				firstFound = true
			}
			if s == r[1] {
				secondFound = true
			}
			if secondFound && !firstFound {
				return false
			}
		}
	}

	return true
}

func containsAll(sequence []int, ints ...int) bool {
	for _, i := range ints {
		if !slices.Contains(sequence, i) {
			return false
		}
	}

	return true
}

func parseRule(line string) [2]int {
	parts := strings.Split(line, "|")

	a := parseInt(parts[0])
	b := parseInt(parts[1])

	return [2]int{a, b}
}

func parseInt(s string) int {
	i, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		panic(err)
	}
	return i
}
