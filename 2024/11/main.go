package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	file, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err) // Terminate if the input file cannot be read, as further processing is impossible
	}

	// Part 1
	BlinkAlot(string(file), 25)
	BlinkAlot("125 17", 6)

	// Part 2
	BlinkAlot(string(file), 75)
}

var memoized = make(map[int][]int)

type NumberGroup struct {
	value int
	count int64
}

func BlinkAlot(diskData string, iterations int) {
	values := makeStones(strings.Split(diskData, " "))

	// Group similar numbers
	groups := make(map[int]int64)
	for _, v := range values {
		groups[v]++
	}

	fmt.Printf("Initial groups: %v\n", groups)
	for i := 0; i < iterations; i++ {
		newGroups := make(map[int]int64)

		for value, count := range groups {
			if value == 0 {
				newGroups[1] += count
			} else if countDigits(value)%2 == 0 {
				divisor := pow10(countDigits(value) / 2)
				left := value / divisor
				right := value % divisor
				if left == 0 {
					newGroups[right] += count
				} else {
					newGroups[left] += count
					newGroups[right] += count
				}
			} else {
				newGroups[value*2024] += count
			}
		}

		groups = newGroups
		var total int64
		for _, count := range groups {
			total += count
		}
		fmt.Printf("Iteration %d: %d numbers in %d groups\n", i, total, len(groups))
	}
}

func countDigits(n int) int {
	if n == 0 {
		return 1
	}
	count := 0
	for n > 0 {
		count++
		n /= 10
	}
	return count
}

func pow10(n int) int {
	result := 1
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

func BlinkParallel(diskData string) {
	values := makeStones(strings.Split(diskData, " "))

	numWorkers := runtime.GOMAXPROCS(0) // Use number of available CPUs
	for i := 0; i < 75; i++ {
		values = blinkStonesParallel(values, numWorkers)
	}

	fmt.Println(len(values))
}

func makeStones(s []string) []int {
	if len(s) == 0 {
		return nil
	}

	values := make([]int, 0, len(s))
	for _, numStr := range s {
		values = append(values, atoi(numStr))
	}

	return values
}

func blinkStones(s []int) []int {
	values := make([]int, 0, len(s)*2)

	for _, value := range s {
		values = append(values, blinkStone(value)...)
	}

	return values
}

// returns the new head and tail of the partial list
func blinkStone(s int) []int {
	if result, exists := memoized[s]; exists {
		return result
	}

	if s == 0 {
		return []int{1}
	}

	// Count digits
	temp := s
	digits := 0
	for temp > 0 {
		digits++
		temp /= 10
	}

	if digits%2 == 0 {
		// Split number mathematically
		divisor := 1
		for i := 0; i < digits/2; i++ {
			divisor *= 10
		}
		return []int{s / divisor, s % divisor}
	}

	return []int{s * 2024}
}

func blinkStonesParallel(s []int, numWorkers int) []int {
	if len(s) < numWorkers {
		return blinkStones(s) // Use sequential version for small inputs
	}

	// Create channel for results
	resultChan := make(chan []int, numWorkers)

	// Calculate chunk size for each worker
	chunkSize := len(s) / numWorkers

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = len(s) // Last worker takes any remainder
		}

		go func(chunk []int) {
			results := make([]int, 0, len(chunk)*2)
			for _, value := range chunk {
				results = append(results, blinkStone(value)...)
			}
			resultChan <- results
		}(s[start:end])
	}

	// Collect results
	finalResult := make([]int, 0, len(s)*2)
	for i := 0; i < numWorkers; i++ {
		chunk := <-resultChan
		finalResult = append(finalResult, chunk...)
	}

	return finalResult
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
