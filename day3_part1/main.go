package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type SchematicSymbol struct {
	Value string
	Kind  string
}

type Coordinate struct {
	X, Y int
}

func ParseRow(rowNum int, row string, m map[Coordinate]*string) map[Coordinate]*string {
	cur := new(string)
	for col, rune := range row {
		s := string(rune)
		_, err := strconv.Atoi(s)
		if err != nil {
			// allocate a new pointer and store the symbol we found
			cur = new(string)
			*cur = s

		} else {
			// modify the current pointer as we discover more digits
			*cur += s
		}

		// fmt.Printf("(%d, %d) contains %s\n", col, rowNum, *cur)
		m[Coordinate{X: col, Y: rowNum}] = cur
		if err != nil {
			// immediately allocate a new pointer if the last string wasn't a digit
			cur = new(string)
		}
	}

	return m
}

func main() {
	f, err := os.ReadFile("./sample_input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(f), "\n")
	m := make(map[Coordinate]*string)
	var ans int
	for rowNum, line := range lines {
		m = ParseRow(rowNum, line, m)
	}

	for k, v := range m {
		fmt.Printf("coord: %+v, value: %+v\n", k, *v)
	}

	fmt.Println(ans)
}
