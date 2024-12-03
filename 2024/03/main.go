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

	statements := findMulStatements(string(file))

	total := 0

	for _, statement := range statements {
		total += parseMulStatement(statement)
	}

	fmt.Println(total)
}

func findMulStatements(text string) []string {
	statements := make([]string, 0)
	i := 0

	current := ""
	shouldDo := true
	for i < len(text)-3 {
		if i+7 < len(text) && text[i:i+7] == "don't()" {
			shouldDo = false
			i += 7
			continue
		} else if i+4 < len(text) && text[i:i+4] == "do()" {
			shouldDo = true
			i += 4
			continue
		}

		if i+4 < len(text) && text[i:i+4] == "mul(" {
			current = "mul("
			i += 4
			continue
		} else if current != "" {
			if i+1 < len(text) && text[i] == ')' {
				current += ")"
				if shouldDo {
					statements = append(statements, current)
				}
				current = ""
			} else {
				current += string(text[i])
			}
		}

		i++
	}

	return statements
}

func parseMulStatement(statement string) int {
	// format "mul(1,2)"
	statement = strings.TrimSpace(statement[4 : len(statement)-1])
	first, second := "", ""

	for i := 0; i < len(statement); i++ {
		if statement[i] == ',' {
			first = statement[:i]
			second = statement[i+1:]
			break
		}
	}

	a, err := strconv.Atoi(first)
	if err != nil {
		return 0
	}
	b, err := strconv.Atoi(second)
	if err != nil {
		return 0
	}

	return a * b
}
