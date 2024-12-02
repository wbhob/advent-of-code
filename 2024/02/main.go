package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err) // Terminate if the input file cannot be read, as further processing is impossible
	}

	lines := strings.Split(string(file), "\n")

	lists := parseLists(lines)
	count := 0

	for _, list := range lists {
		if len(list) == 0 {
			continue // Skip empty lists to avoid unnecessary processing
		}

		if !faultTolerance(list) {
			continue // Skip lists that cannot be made valid by removing one element
		}

		fmt.Println(list) // Output the valid list
		count++
	}

	fmt.Println(count) // Output the count of valid lists
}

func parseLists(lines []string) [][]int {
	lists := make([][]int, 0, 100) // Preallocate space for efficiency

	for _, line := range lines {
		parts := strings.Split(line, " ")
		ints := make([]int, 0, len(parts))
		for _, part := range parts {
			if part == "" {
				continue // Skip empty parts to avoid conversion errors
			}
			ints = append(ints, atoi(part))
		}
		lists = append(lists, ints)
	}

	return lists
}

func faultTolerance(list []int) bool {
	if isOk(list) {
		return true // List is already valid, no need to modify
	}

	for i := 0; i < len(list); i++ {
		sublist := make([]int, len(list)-1)
		copy(sublist, list[:i])
		copy(sublist[i:], list[i+1:])

		if isOk(sublist) {
			fmt.Println("removed index", i) // Indicate which index was removed to make the list valid
			return true
		}
	}

	return false // List cannot be made valid by removing a single element
}

func isOk(list []int) bool {
	return isOneDirection(list) && processList(list) // List must be both monotonic and have valid differences
}

func isOneDirection(list []int) bool {
	// Determine direction from first two numbers
	if len(list) < 2 {
		return true // A list with fewer than two elements is trivially monotonic
	}

	increasing := list[0] < list[1]

	// Check that all subsequent pairs follow the same direction
	for i := 1; i < len(list)-1; i++ {
		if increasing {
			if list[i] >= list[i+1] {
				return false // Breaks increasing order
			}
		} else {
			if list[i] <= list[i+1] {
				return false // Breaks decreasing order
			}
		}
	}
	return true
}

func processList(list []int) bool {
	for i := 0; i < len(list)-1; i++ {
		av := abs(list[i] - list[i+1])
		if av > 3 || av < 1 {
			return false // Differences must be within the range [1, 3] to be valid
		}
	}

	return true
}

func abs(a int) int {
	if a < 0 {
		return -a // Return positive value for negative input
	}
	return a
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err) // Terminate if conversion fails, as the input is expected to be valid integers
	}
	return i
}
