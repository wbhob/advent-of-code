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

	Part1(string(file))
	Part2(string(file))
}

const NilChannel = '.'

func Part1(mapData string) {
	cityMap := NewMap(mapData)
	for i := range cityMap.data {
		cityMap.checkAntinode(i%cityMap.width, i/cityMap.width)
	}

	antinodes := 0
	for _, tile := range cityMap.data {
		if tile.isAntinode {
			antinodes++
		}
	}

	fmt.Println(antinodes)
}

func Part2(mapData string) {
	cityMap := NewMap(mapData)
	for i := range cityMap.data {
		cityMap.checkAntinodePart2(i%cityMap.width, i/cityMap.width)
	}

	antinodes := 0
	for _, tile := range cityMap.data {
		if tile.isAntinode {
			antinodes++
		}
	}

	fmt.Println(antinodes)
}

type Map struct {
	data     []*MapTile
	channels map[rune]bool
	width    int
	height   int
}

type MapTile struct {
	antennaChannel   rune
	isAntinode       bool
	antinodeChannels []rune
}

func NewMap(raw string) *Map {
	lines := strings.Split(raw, "\n")
	height := len(lines)
	width := len(lines[0])
	data := make([]*MapTile, height*width)
	channels := make(map[rune]bool)
	for y, line := range lines {
		for x, char := range line {
			data[y*width+x] = NewMapTile(char)
			channels[char] = true
		}
	}

	return &Map{
		data:     data,
		channels: channels,
		width:    width,
		height:   height,
	}
}

func NewMapTile(channel rune) *MapTile {
	return &MapTile{
		antennaChannel:   channel,
		isAntinode:       false,
		antinodeChannels: make([]rune, 0),
	}
}

func (m *Map) get(x, y int) *MapTile {
	return m.data[y*m.width+x]
}

func (m *Map) coordsOf(index int) (int, int) {
	return index % m.width, index / m.width
}

// checkAntinode checks for antinodes at a given point if it is an antenna
// for a channel. Antinodes are when there is 2 antennas with the same
// channel in the same line in any direction and one is twice as far from
// the point as the other.
func (m *Map) checkAntinode(x, y int) {
	channel := m.get(x, y).antennaChannel
	if channel == NilChannel {
		return
	}

	// Find all antennas with the same channel
	var antennas [][2]int
	for i, tile := range m.data {
		if tile.antennaChannel == channel {
			antennaX, antennaY := m.coordsOf(i)
			antennas = append(antennas, [2]int{antennaX, antennaY})
		}
	}

	// For each point in the grid
	for py := 0; py < m.height; py++ {
		for px := 0; px < m.width; px++ {
			// For each pair of antennas
			for i := 0; i < len(antennas); i++ {
				for j := i + 1; j < len(antennas); j++ {
					ant1 := antennas[i]
					ant2 := antennas[j]

					// Check if point is collinear with the two antennas
					if !isCollinear(px, py, ant1[0], ant1[1], ant2[0], ant2[1]) {
						continue
					}

					// Calculate distances from point to both antennas
					dist1 := manhattan(px, py, ant1[0], ant1[1])
					dist2 := manhattan(px, py, ant2[0], ant2[1])

					// If one distance is exactly twice the other, it's an antinode
					if (dist1 > 0 && dist2 > 0) && (dist1 == 2*dist2 || dist2 == 2*dist1) {
						m.get(px, py).isAntinode = true
					}
				}
			}
		}
	}
}

// isCollinear checks if three points lie on the same line (horizontal, vertical, or diagonal)
func isCollinear(x, y, x1, y1, x2, y2 int) bool {
	// Check horizontal
	if y == y1 && y == y2 {
		return true
	}
	// Check vertical
	if x == x1 && x == x2 {
		return true
	}
	// Check diagonal
	dx1 := x - x1
	dy1 := y - y1
	dx2 := x - x2
	dy2 := y - y2

	// Check if slopes are equal (avoiding division by zero)
	return dx1*dy2 == dx2*dy1
}

func manhattan(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (m *Map) checkAntinodePart2(x, y int) {
	channel := m.get(x, y).antennaChannel
	if channel == NilChannel {
		return
	}

	// Find all antennas with the same channel
	var antennas [][2]int
	for i, tile := range m.data {
		if tile.antennaChannel == channel {
			antennaX, antennaY := m.coordsOf(i)
			antennas = append(antennas, [2]int{antennaX, antennaY})
		}
	}

	// For each point in the grid
	for py := 0; py < m.height; py++ {
		for px := 0; px < m.width; px++ {
			// For each pair of antennas
			for i := 0; i < len(antennas); i++ {
				for j := i + 1; j < len(antennas); j++ {
					ant1 := antennas[i]
					ant2 := antennas[j]

					// If point is collinear with the two antennas, it's an antinode
					if isCollinear(px, py, ant1[0], ant1[1], ant2[0], ant2[1]) {
						m.get(px, py).isAntinode = true
					}
				}
			}
		}
	}
}
