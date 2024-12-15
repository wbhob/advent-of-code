package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Equation struct {
	TestValue int
	Values    []int
}

type Operator int

const (
	Add Operator = iota
	Multiply
	Concatenate
)

func main() {
	file, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err) // Terminate if the input file cannot be read, as further processing is impossible
	}

	// line format: testvalue: int (...int)
	lines := strings.Split(string(file), "\n")
	equations := makeEquations(lines)

	fmt.Println("Running concurrent version:")
	start := time.Now()
	testEquations(equations)
	fmt.Printf("Concurrent time: %v\n", time.Since(start))

	fmt.Println("\nRunning sequential version:")
	start = time.Now()
	testEquationsSequential(equations)
	fmt.Printf("Sequential time: %v\n", time.Since(start))
}

// The engineers just need the total calibration result,
// which is the sum of the test values from just the
// equations that could possibly be true.
func testEquations(equations []Equation) {
	sum := 0
	for _, equation := range equations {
		result := testEquation(equation)
		if result != 0 {
			sum += result
		}
	}

	fmt.Println(sum)
}

func testEquation(equation Equation) int {
	// Create a channel to receive results
	resultChan := make(chan int)
	// Create a channel to signal when we found a solution
	done := make(chan struct{})

	// Number of goroutines to use
	numWorkers := 3

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			operators := make([]Operator, len(equation.Values)-1)
			result := recursiveTestWorker(equation, operators, 0, workerID, numWorkers, done)
			resultChan <- result // Always send a result, even if it's 0
		}(i)
	}

	// Wait for all results, return first non-zero one found
	var finalResult int
	doneClose := false
	for i := 0; i < numWorkers; i++ {
		if result := <-resultChan; result != 0 && !doneClose {
			close(done) // Signal other goroutines to stop
			doneClose = true
			finalResult = result
		}
	}

	return finalResult
}

func recursiveTestWorker(equation Equation, operators []Operator, position int, workerID, numWorkers int, done chan struct{}) int {
	// Check if we should stop
	select {
	case <-done:
		return 0
	default:
	}

	// If we've filled all operator positions, test the equation
	if position == len(operators) {
		return testEquationWithPermutation(equation, operators)
	}

	// For the first position, each worker only tries their assigned operators
	startOp := 0
	endOp := 3
	if position == 0 {
		startOp = workerID
		endOp = startOp + 1
	}

	// Try each operator at the current position
	for op := Operator(startOp); op < Operator(endOp); op++ {
		operators[position] = op
		if result := recursiveTestWorker(equation, operators, position+1, workerID, numWorkers, done); result != 0 {
			return result
		}
	}

	return 0
}

func testEquationWithPermutation(equation Equation, permutation []Operator) int {
	// Make copies of the slices to work with
	values := make([]int, len(equation.Values))
	copy(values, equation.Values)

	result := values[0]

	for i := 0; i < len(permutation); i++ {
		switch permutation[i] {
		case Add:
			result += values[i+1]
		case Multiply:
			result *= values[i+1]
		case Concatenate:
			// Convert both numbers to strings, concatenate, then convert back to int
			resultStr := fmt.Sprintf("%d%d", result, values[i+1])
			result = atoi(resultStr)
		}
	}

	if result == equation.TestValue {
		return equation.TestValue
	}
	return 0
}

// permuteOperators generates all possible combinations of Add and Multiply operators
// for a given number of values. For n values, we need (n-1) operators between them.
// The function uses binary counting to generate all combinations efficiently.
func permuteOperators(valueCount int) [][]Operator {
	// We need one less operator than values (e.g., for 3 values, we need 2 operators: a ? b ? c)
	operatorCount := valueCount - 1

	// Calculate total number of possible combinations using bit shifting
	// For n operators, we need 2^n combinations (each position can be Add or Multiply)
	totalCombinations := 1 << operatorCount // Same as 2^numberOfOperators

	// Initialize slice to hold all permutations
	permutations := make([][]Operator, totalCombinations)

	// Iterate through all possible binary numbers from 0 to 2^numberOfOperators-1
	// Each number represents a unique combination of operators
	for i := 0; i < totalCombinations; i++ {
		// Create a new slice for this combination of operators
		operators := make([]Operator, operatorCount)

		// For each bit position in the current number i
		// A set bit (1) represents Multiply, an unset bit (0) represents Add
		for j := 0; j < operatorCount; j++ {
			// Check if the j-th bit is set in number i
			// (i & (1 << j)) performs bitwise AND with a number that has only the j-th bit set
			if (i & (1 << j)) != 0 {
				operators[j] = Multiply
			} else {
				operators[j] = Add
			}
		}
		permutations[i] = operators
	}

	return permutations
}

func makeEquations(lines []string) []Equation {
	equations := make([]Equation, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, ":")
		testValue := atoi(parts[0])
		values := strings.Split(parts[1], " ")
		valuesInt := make([]int, len(values))
		for j, value := range values {
			if trimmed := strings.TrimSpace(value); trimmed != "" {
				valuesInt[j] = atoi(trimmed)
			}
		}
		equations[i] = Equation{TestValue: testValue, Values: valuesInt}
	}
	return equations
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func testEquationsSequential(equations []Equation) {
	sum := 0
	for _, equation := range equations {
		result := testEquationSequential(equation)
		if result != 0 {
			sum += result
		}
	}
	fmt.Println("Sequential result:", sum)
}

func testEquationSequential(equation Equation) int {
	operators := make([]Operator, len(equation.Values)-1)
	return recursiveTest(equation, operators, 0)
}

func recursiveTest(equation Equation, operators []Operator, position int) int {
	// If we've filled all operator positions, test the equation
	if position == len(operators) {
		return testEquationWithPermutation(equation, operators)
	}

	// Try each operator at the current position
	for op := Operator(0); op < Operator(3); op++ {
		operators[position] = op
		if result := recursiveTest(equation, operators, position+1); result != 0 {
			return result
		}
	}

	return 0
}
