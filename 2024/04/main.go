package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err) // Terminate if the input file cannot be read, as further processing is impossible
	}

	lines := strings.Split(string(file), "\n")

	fmt.Println(lines)

	sourceGrid := make([][]rune, len(lines))

	for i, line := range lines {
		sourceGrid[i] = make([]rune, len(line))
		for j, col := range line {
			sourceGrid[i][j] = col
		}
	}

	_, acrossCount := readAcross(sourceGrid)
	_, downCount := readDown(sourceGrid)
	_, diagonalCount := readDiagonal(sourceGrid)
	_, reverseDiagonalCount := readReverseDiagonal(sourceGrid)

	// for i, row := range acrossGrid {
	// 	for j := range row {
	// 		if !acrossGrid[i][j] && !downGrid[i][j] && !diagonalGrid[i][j] && !reverseDiagonalGrid[i][j] {
	// 			sourceGrid[i][j] = '.'
	// 		}
	// 	}
	// }

	fmt.Println(acrossCount, downCount, diagonalCount, reverseDiagonalCount)
	fmt.Println(acrossCount + downCount + diagonalCount + reverseDiagonalCount)

	masGrid, masCount := readCrossMAS(sourceGrid)

	for i, row := range masGrid {
		for j := range row {
			if !masGrid[i][j] {
				sourceGrid[i][j] = '.'
			}
		}
	}

	for _, row := range sourceGrid {
		fmt.Println(string(row))
	}

	fmt.Println(masCount)
}

func readAcross(sourceGrid [][]rune) ([][]bool, int) {
	grid := make([][]bool, len(sourceGrid))

	count := 0
	for i, row := range sourceGrid {
		grid[i] = make([]bool, len(row))

		for j := 0; j <= len(row)-4; j++ {
			if isXMAS(row[j],
				row[j+1],
				row[j+2],
				row[j+3]) {

				grid[i][j] = true
				grid[i][j+1] = true
				grid[i][j+2] = true
				grid[i][j+3] = true

				count++

			}
		}
	}

	return grid, count
}

func readDown(sourceGrid [][]rune) ([][]bool, int) {
	grid := make([][]bool, len(sourceGrid))

	for i := range grid {
		grid[i] = make([]bool, len(sourceGrid[i]))
	}

	count := 0
	for i := 0; i <= len(sourceGrid)-4; i++ {
		for j := range sourceGrid[i] {
			if isXMAS(sourceGrid[i][j],
				sourceGrid[i+1][j],
				sourceGrid[i+2][j],
				sourceGrid[i+3][j]) {

				grid[i][j] = true
				grid[i+1][j] = true
				grid[i+2][j] = true
				grid[i+3][j] = true

				count++
			}
		}
	}

	return grid, count
}

func readDiagonal(sourceGrid [][]rune) ([][]bool, int) {
	grid := make([][]bool, len(sourceGrid))

	for i := range grid {
		grid[i] = make([]bool, len(sourceGrid[i]))
	}

	count := 0

	for i := 0; i <= len(sourceGrid)-4; i++ {
		for j := 0; j <= len(sourceGrid[i])-4; j++ {
			if isXMAS(sourceGrid[i][j],
				sourceGrid[i+1][j+1],
				sourceGrid[i+2][j+2],
				sourceGrid[i+3][j+3]) {

				grid[i][j] = true
				grid[i+1][j+1] = true
				grid[i+2][j+2] = true
				grid[i+3][j+3] = true

				count++
			}
		}
	}

	return grid, count
}

func readReverseDiagonal(sourceGrid [][]rune) ([][]bool, int) {
	grid := make([][]bool, len(sourceGrid))

	for i := range grid {
		grid[i] = make([]bool, len(sourceGrid[i]))
	}

	count := 0

	for i := 0; i <= len(sourceGrid)-4; i++ {
		for j := len(sourceGrid[i]) - 1; j >= 3; j-- {
			if isXMAS(sourceGrid[i][j],
				sourceGrid[i+1][j-1],
				sourceGrid[i+2][j-2],
				sourceGrid[i+3][j-3]) {

				grid[i][j] = true
				grid[i+1][j-1] = true
				grid[i+2][j-2] = true
				grid[i+3][j-3] = true

				count++
			}
		}
	}

	return grid, count
}

func readCrossMAS(sourceGrid [][]rune) ([][]bool, int) {
	grid := make([][]bool, len(sourceGrid))

	for i := range grid {
		grid[i] = make([]bool, len(sourceGrid[i]))
	}

	count := 0
	for i := 0; i <= len(sourceGrid)-3; i++ {
		for j := 0; j <= len(sourceGrid[i])-3; j++ {
			if isMAS(sourceGrid[i][j],
				sourceGrid[i+1][j+1],
				sourceGrid[i+2][j+2]) &&
				isMAS(sourceGrid[i+2][j],
					sourceGrid[i+1][j+1],
					sourceGrid[i][j+2]) {

				grid[i][j] = true
				grid[i+1][j+1] = true
				grid[i+2][j+2] = true
				grid[i][j+2] = true
				grid[i+2][j] = true

				count++
			}
		}
	}

	return grid, count
}

func isMAS(m, a, s rune) bool {
	return m == 'M' && a == 'A' && s == 'S' ||
		m == 'S' && a == 'A' && s == 'M'
}

func isXMAS(x, m, a, s rune) bool {
	return (x == 'X' && m == 'M' && a == 'A' && s == 'S') ||
		(x == 'S' && m == 'A' && a == 'M' && s == 'X')
}
