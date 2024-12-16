package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	file, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err) // Terminate if the input file cannot be read, as further processing is impossible
	}

	Part1(string(file))
	Part2(string(file))
}

func Part1(diskData string) {
	fmt.Println("Starting Part 1...")
	topoMap := NewTopoMap(diskData)
	trailheads := topoMap.findTrailheads()

	result := RecurseResult{}

	scoreSum := 0
	for _, trailhead := range trailheads {
		r := topoMap.recurseTraverse(trailhead)
		result.dones += r.dones
		result.peaks += r.peaks
		uniquePeaks := make(map[int]bool)
		for _, peak := range r.peakCoords {
			uniquePeaks[peak] = true
		}

		scoreSum += len(uniquePeaks)
	}

	fmt.Println(result.dones, " dones", result.peaks, " peaks")

	fmt.Println(scoreSum, " score sum")
}

func Part2(diskData string) {
}

type TopoCoord int

func (t TopoCoord) isTrailhead() bool {
	return t == 0
}

type TopoMap struct {
	data   []TopoCoord
	width  int
	height int
}

func NewTopoMap(mapData string) *TopoMap {
	lines := strings.Split(mapData, "\n")
	width := len(lines[0])
	height := len(lines)

	data := make([]TopoCoord, width*height)

	i := 0
	for _, char := range mapData {
		if unicode.IsSpace(char) {
			continue
		}

		data[i] = TopoCoord(rtoi(char))
		i++
	}

	return &TopoMap{data: data, width: width, height: height}
}

func (t *TopoMap) get(x, y int) TopoCoord {
	return t.data[y*t.width+x]
}

func (t *TopoMap) coordsAt(index int) (x, y int) {
	return index % t.width, index / t.width
}

func (t *TopoMap) findTrailheads() []int {
	trailheads := make([]int, 0)
	for i, coord := range t.data {
		if coord.isTrailhead() {
			trailheads = append(trailheads, i)
		}
	}

	return trailheads
}

type RecurseResult struct {
	dones      int
	peaks      int
	peakCoords []int
}

func (t *TopoMap) recurseTraverse(index int) RecurseResult {
	search := t.findNext(index)

	if search.isDone {
		if search.isPeak {
			return RecurseResult{dones: 1, peaks: 1, peakCoords: search.peakCoords}
		}
		return RecurseResult{dones: 1, peaks: 0, peakCoords: []int{}}
	}

	result := RecurseResult{}

	if search.n {
		r := t.recurseTraverse(index - t.width)
		result.dones += r.dones
		result.peaks += r.peaks
		result.peakCoords = append(result.peakCoords, r.peakCoords...)
	}

	if search.e {
		r := t.recurseTraverse(index + 1)
		result.dones += r.dones
		result.peaks += r.peaks
		result.peakCoords = append(result.peakCoords, r.peakCoords...)
	}

	if search.s {
		r := t.recurseTraverse(index + t.width)
		result.dones += r.dones
		result.peaks += r.peaks
		result.peakCoords = append(result.peakCoords, r.peakCoords...)
	}

	if search.w {
		r := t.recurseTraverse(index - 1)
		result.dones += r.dones
		result.peaks += r.peaks
		result.peakCoords = append(result.peakCoords, r.peakCoords...)
	}

	return result
}

type NextSearchInfo struct {
	n, e, s, w, isDone, isPeak bool
	peakCoords                 []int
}

func (t *TopoMap) findNext(index int) NextSearchInfo {
	value := t.data[index]
	x, y := t.coordsAt(index)

	if value == 9 {
		return NextSearchInfo{isDone: true, isPeak: true, peakCoords: []int{index}}
	}

	result := NextSearchInfo{}

	if y > 0 {
		result.n = t.get(x, y-1) == value+1
	}

	if x < t.width-1 {
		result.e = t.get(x+1, y) == value+1
	}

	if y < t.height-1 {
		result.s = t.get(x, y+1) == value+1
	}

	if x > 0 {
		result.w = t.get(x-1, y) == value+1
	}

	if !result.n && !result.e && !result.s && !result.w {
		result.isDone = true
		result.isPeak = false
	}

	return result
}

func rtoi(r rune) int {
	i, err := strconv.Atoi(string(r))
	if err != nil {
		panic(err)
	}
	return i
}
