package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Glyph struct {
	Schematic   *Schematic
	ID          int
	Value       string
	Kind        string
	Coordinates Point
}

func (g Glyph) Neighbors() (neighbors []Glyph) {
	// TODO: make this function aware of Glyph.ID; don't return duplicates
	getNeighbors := func(x, y int) {
		for i := x - 1; i <= x+len(g.Value); i++ {
			v, ok := g.Schematic.GlyphAt(i, y)
			if ok {
				neighbors = append(neighbors, v)
			}
		}
	}

	// look up
	getNeighbors(g.Coordinates.X, g.Coordinates.Y-1)
	// look down
	getNeighbors(g.Coordinates.X, g.Coordinates.Y+1)
	// look left
	v, ok := g.Schematic.GlyphAt(g.Coordinates.X-1, g.Coordinates.Y)
	if ok {
		neighbors = append(neighbors, v)
	}
	// look right
	v, ok = g.Schematic.GlyphAt(g.Coordinates.X+len(g.Value), g.Coordinates.Y)
	if ok {
		neighbors = append(neighbors, v)
	}

	return
}

func NewGlyph(s *Schematic, id int, value string, coords Point) Glyph {
	var kind string
	if _, err := strconv.Atoi(value); err == nil {
		kind = "number"
	} else if value == "." {
		kind = "blank"
	} else if value == "*" {
		kind = "gear"
	} else {
		kind = "symbol"
	}

	return Glyph{
		Schematic:   s,
		ID:          id,
		Value:       value,
		Kind:        kind,
		Coordinates: coords,
	}
}

type Point struct {
	X, Y int
}

type Schematic struct {
	rows   [][]Glyph
	nextID int
}

func (s *Schematic) NextID() int {
	id := s.nextID
	s.nextID += 1

	return id
}

func (s Schematic) GlyphAt(x, y int) (Glyph, bool) {
	if y < 0 || y >= len(s.rows) {
		// row out of bounds
		return Glyph{}, false
	}
	if x < 0 || x >= len(s.rows[y]) {
		// column out of bounds
		return Glyph{}, false
	}

	return s.rows[y][x], true
}

func (sch *Schematic) AppendRow(row string) (glyphs []Glyph) {
	buf := []Point{}
	rowNum := len(sch.rows)
	var cur string
	id := sch.NextID()
	for colNum, rune := range row {
		s := string(rune)
		_, err := strconv.Atoi(s)
		if err != nil {
			// if not digit, assign all previous coordinates to accumulated value,
			for _, coord := range buf {
				glyphs = append(glyphs, NewGlyph(sch, id, cur, coord))
			}
			id = sch.NextID()
			buf = []Point{}

			// then immediately assign the current coordinate and value
			glyphs = append(glyphs, NewGlyph(sch, id, s, Point{X: colNum, Y: rowNum}))
			id = sch.NextID()
			cur = ""
		} else {
			// if digit, don't immediately assign the current coordinate
			buf = append(buf, Point{X: colNum, Y: rowNum})

			// store the current digit to be assigned later
			cur += s
			if colNum == len(row)-1 {
				// if the last column has a digit, always assign it
				for _, coord := range buf {
					glyphs = append(glyphs, NewGlyph(sch, id, cur, coord))
				}
				id = sch.NextID()
			}
		}
	}

	sch.rows = append(sch.rows, glyphs)

	return
}

func main() {
	f, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(f), "\n")
	s := Schematic{}
	var gears []Glyph
	for _, line := range lines {
		glyphs := s.AppendRow(line)
		for i := 0; i < len(glyphs); {
			g := glyphs[i]
			if g.Kind == "gear" {
				gears = append(gears, g)
			}

			i += len(g.Value)
		}
	}

	// examine neighbors of all numbers to see if any contain symbols
	var ans int
	for _, g := range gears {
		neighbors := g.Neighbors()
		for _, n := range neighbors {
			if n.Kind == "number" {
				v, err := strconv.Atoi(g.Value)
				if err != nil {
					panic(err)
				}
				ans += v
				break
			}
		}
	}

	fmt.Println(ans)
}
