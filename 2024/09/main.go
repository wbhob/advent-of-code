package main

import (
	"fmt"
	"os"
	"strconv"
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
	disk := NewDisk(diskData)
	disk.compress()
}

func Part2(diskData string) {
	disk := NewDisk(diskData)
	disk.compress2()
}

const FREE = -1

type Disk struct {
	data []int
}

func NewDisk(diskData string) *Disk {
	disk := &Disk{
		data: make([]int, 0, len(diskData)),
	}

	for i, char := range diskData {
		count := rtoi(char)
		id := i / 2
		for j := 0; j < count; j++ {
			if i%2 == 0 {
				disk.data = append(disk.data, id)
			} else {
				disk.data = append(disk.data, FREE)
			}
		}

	}

	return disk
}

func (d *Disk) compress() {
	data := d.data

	head := 0
	tail := len(data) - 1

	for head < tail {
		if data[head] != FREE {
			head++
			continue
		}

		if data[tail] == FREE {
			tail--
			continue
		}

		data[head] = data[tail]
		data[tail] = FREE
		head++
		tail--
	}

	checksum := 0
	for i, fileID := range data {
		if fileID == FREE {
			continue
		}

		checksum += i * fileID
	}

	fmt.Printf("checksum: %d\n", checksum)
}

func (d *Disk) compress2() int {
	data := d.data

	moved := make(map[int]bool)

	for blockID := sliceMax(data); blockID > 0; blockID-- {
		if moved[blockID] {
			continue
		}

		moved[blockID] = true

		block := d.getBlock(blockID)

		writeIndex := d.balloc(block.size)
		if writeIndex == -1 || writeIndex >= block.start {
			continue
		}

		for j := writeIndex; j < writeIndex+block.size; j++ {
			data[j] = blockID
		}

		for j := block.start; j < block.end(); j++ {
			data[j] = FREE
		}

		if blockID == 5200 {
			block := d.getBlock(blockID)
			writeIndex := d.balloc(block.size)
			fmt.Printf("Block 5200: size=%d, current_start=%d, found_write_index=%d\n",
				block.size, block.start, writeIndex)
		}

	}

	checksum := 0
	for i, fileID := range data {
		if fileID == FREE {
			fmt.Print(".")
			continue
		} else {
			fmt.Printf("%d", fileID)
		}

		checksum += i * fileID
	}

	fmt.Printf("checksum: %d\n", checksum)
	return checksum
}

type Block struct {
	start int
	size  int
}

func (b *Block) end() int {
	return b.start + b.size
}

func (d *Disk) getBlock(id int) *Block {
	start := 0
	for i := 0; i < len(d.data); i++ {
		if d.data[i] == id {
			start = i
			break
		}
	}

	size := 0
	for i := start; i < len(d.data); i++ {
		if d.data[i] == id {
			size++
		} else {
			break
		}
	}

	return &Block{
		start: start,
		size:  size,
	}
}

// block allocation
func (d *Disk) balloc(size int) int {
	i := 0
	blockSize := 0
	for i < len(d.data) && blockSize < size {
		if d.data[i] == FREE {
			blockSize++
		} else {
			blockSize = 0
		}
		i++
	}

	if blockSize < size {
		return -1
	}

	return i - blockSize
}

func sliceMax(data []int) int {
	bound := 0
	for _, id := range data {
		bound = max(bound, id)
	}
	return bound
}

func rtoi(r rune) int {
	i, err := strconv.Atoi(string(r))
	if err != nil {
		panic(err)
	}
	return i
}
