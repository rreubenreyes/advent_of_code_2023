package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Glyph struct {
	Schematic   *Schematic
	Value       string
	Kind        string
	Coordinates Point
}

// TODO:  reimplement this, too little sleep rn
func (g Glyph) Neighbors() (neighbors []Glyph) {
	for i := g.Coordinates.X; i < g.Coordinates.X+len(g.Value); i++ {
		neighbors = append(neighbors, g.Schematic.GlyphAt(i, g.Coordinates.Y))
	}

	return
}

func NewGlyph(s *Schematic, value string, coords Point) Glyph {
	var kind string
	if _, err := strconv.Atoi(value); err == nil {
		kind = "number"
	} else if value == "." {
		kind = "blank"
	} else {
		kind = "symbol"
	}

	return Glyph{
		Schematic:   s,
		Value:       value,
		Kind:        kind,
		Coordinates: coords,
	}
}

type Point struct {
	X, Y int
}

type Schematic struct {
	rows [][]Glyph
}

// TODO:  reimplement this with a map, too little sleep rn
func (s Schematic) GlyphAt(x, y int) Glyph {
	return s.rows[y][x]
}

func (sch *Schematic) AppendRow(row string) (glyphs []Glyph) {
	var cur string
	rowNum := len(sch.rows)
	for colNum, rune := range row {
		s := string(rune)
		_, err := strconv.Atoi(s)
		if err != nil {
			// string at col is not a digit
			if len(cur) > 0 {
				glyphs = append(glyphs, NewGlyph(sch, cur, Point{
					X: colNum - 1 - (len(cur) - 1),
					Y: rowNum,
				}))
				cur = ""
			}

			glyphs = append(glyphs, NewGlyph(sch, s, Point{
				X: colNum - (len(cur) - 1),
				Y: rowNum,
			}))
		} else {
			// string at col is a digit
			cur += s
		}

	}

	sch.rows = append(sch.rows, glyphs)

	return
}

func main() {
	f, err := os.ReadFile("./sample_input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(f), "\n")
	s := Schematic{}
	var numbers []Glyph
	for _, line := range lines {
		glyphs := s.AppendRow(line)
		for _, g := range glyphs {
			if g.Kind == "number" {
				numbers = append(numbers, g)
			}
		}
	}

	// examine neighbors of all numbers to see if any contain symbols
	fmt.Printf("numbers: %+v\n", numbers)
	fmt.Printf("neighbors: %+v\n", numbers[0].Neighbors())

	var ans int
	fmt.Println(ans)
}
