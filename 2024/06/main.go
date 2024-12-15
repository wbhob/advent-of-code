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
	grid := make([][]bool, len(lines))
	history := make([][]Direction, len(lines))
	start := Position{}

	for i, line := range lines {
		grid[i] = make([]bool, len(line))
		history[i] = make([]Direction, len(line))
		for j, char := range line {
			if char == '#' {
				grid[i][j] = true
			} else if char == '^' {
				start = Position{i, j}
			}
		}
	}

	state := &State{
		grid:      grid,
		history:   history,
		direction: North,
		pos:       start,
	}

	state.history[start.row][start.col] = North
	state.moveCount++

	loopCount := 0
	for i := range state.grid {
		for j := range state.grid[i] {
			if i == start.row && j == start.col {
				continue
			}

			if state.willLoopWithObstacleAt(Position{i, j}) {
				loopCount++
			}
		}
	}

	fmt.Println(state.moveCount)
	fmt.Println(state.turnCount)
	fmt.Println(loopCount)
}

type Direction int

const (
	None Direction = iota
	North
	East
	South
	West
)

type Position struct {
	row, col int
}

type State struct {
	grid      [][]bool
	history   [][]Direction
	direction Direction
	pos       Position
	moveCount int
	turnCount int
}

func (s *State) isDone() bool {
	return s.pos.row == 0 && s.direction == North ||
		s.pos.row == len(s.grid)-1 && s.direction == South ||
		s.pos.col == 0 && s.direction == West ||
		s.pos.col == len(s.grid[0])-1 && s.direction == East
}

func (s *State) next() {
	if s.isCollision() {
		s.turn()
	} else {
		s.move()
	}
}

func (s *State) turn() {
	s.turnCount++
	switch s.direction {
	case North:
		s.direction = East
	case East:
		s.direction = South
	case South:
		s.direction = West
	case West:
		s.direction = North
	}
}

func (s *State) move() {
	if s.history[s.pos.row][s.pos.col] == None {
		s.history[s.pos.row][s.pos.col] = s.direction
		s.moveCount++
	}

	switch s.direction {
	case North:
		s.pos.row--
	case East:
		s.pos.col++
	case South:
		s.pos.row++
	case West:
		s.pos.col--
	}
}

func (s *State) isCollision() bool {
	nextRow, nextCol := s.pos.row, s.pos.col
	switch s.direction {
	case North:
		nextRow = s.pos.row - 1
	case East:
		nextCol = s.pos.col + 1
	case South:
		nextRow = s.pos.row + 1
	case West:
		nextCol = s.pos.col - 1
	}

	// Check if next position would be out of bounds
	if nextRow < 0 || nextRow >= len(s.grid) ||
		nextCol < 0 || nextCol >= len(s.grid[0]) {
		return true
	}

	return s.grid[nextRow][nextCol]
}

func (s *State) willLoopWithObstacleAt(pos Position) bool {
	grid := make([][]bool, len(s.grid))
	for i := range grid {
		grid[i] = make([]bool, len(s.grid[i]))
		copy(grid[i], s.grid[i])
	}
	history := make([][]Direction, len(s.history))
	for i := range history {
		history[i] = make([]Direction, len(s.history[i]))
		copy(history[i], s.history[i])
	}

	grid[pos.row][pos.col] = true

	state := &State{
		grid:      grid,
		history:   history,
		direction: s.direction,
		pos:       s.pos,
	}

	for !state.isDone() {
		state.next()
		if state.history[state.pos.row][state.pos.col] == state.direction {
			return true
		}
	}

	return false
}
