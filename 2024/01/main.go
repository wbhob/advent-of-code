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
		panic(err)
	}

	lines := strings.Split(string(file), "\n")

	a, b := parseLists(lines)
	slices.Sort(a)
	slices.Sort(b)

	distance := 0

	for i := range a {
		distance += abs(a[i] - b[i])
	}

	fmt.Printf("distance: %d\n", distance)

	rightListCounts := make(map[int]int)
	for _, value := range b {
		rightListCounts[value]++
	}

	similarityScore := 0
	for _, value := range a {
		if rightListCounts[value] > 0 {
			similarityScore += value * rightListCounts[value]
		}
	}

	fmt.Printf("similarity score: %d\n", similarityScore)
}

func parseLists(lines []string) ([]int, []int) {
	a, b := make([]int, 0, 100), make([]int, 0, 100)

	for _, line := range lines {
		parts := strings.Split(line, "   ")
		if len(parts) != 2 {
			panic("invalid line")
		}
		first, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}
		second, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		a = append(a, first)
		b = append(b, second)
	}

	return a, b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
