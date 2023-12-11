package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ScratchCard struct {
	Winning map[int]struct{}
	Have    map[int]struct{}
}

func (sc ScratchCard) Score() (score int) {
	for k := range sc.Have {
		_, ok := sc.Winning[k]
		if ok && score == 0 {
			score = 1
		} else if ok && score > 0 {
			score *= 2
		}
	}

	return
}

func ScratchCardFromLine(line string) (sc ScratchCard) {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		panic(errors.New("unexpected line format"))
	}

	numbers := parts[1]
	nn := strings.Split(numbers, "|")
	if len(nn) != 2 {
		panic(errors.New("unexpected line format"))
	}

	sc.Winning = make(map[int]struct{})
	for _, n := range strings.Fields(nn[0]) {
		v, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		sc.Winning[v] = struct{}{}
	}

	sc.Have = make(map[int]struct{})
	for _, n := range strings.Fields(nn[1]) {
		v, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		sc.Have[v] = struct{}{}
	}

	return
}

func ReadLines(f string) []string {
	b, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}

	return strings.Split(string(b), "\n")
}

func main() {
	lines := ReadLines("./input.txt")
	var ans int
	for _, line := range lines {
		card := ScratchCardFromLine(line)
		ans += card.Score()
	}

	fmt.Println(ans)
}
