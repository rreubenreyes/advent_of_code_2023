package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ScratchCardTable struct {
	UniqueCards      map[int]ScratchCard
	Copies           map[int]int
	totalCardsCount  int
	winningsMemoized bool
}

func (t *ScratchCardTable) AddCard(card ScratchCard) {
	t.UniqueCards[card.ID] = card
	t.totalCardsCount++
}

func (t *ScratchCardTable) Winnings() int {
	if t.winningsMemoized {
		return t.totalCardsCount
	}
	registerCopies := func(origID, matching int, times int) {
		for i := origID; i < origID+matching && i < len(t.UniqueCards); i++ {
			// register all copies for the original card
			copyID := i + 1
			fmt.Printf("%dx card %d wins card %d; ", times, origID, copyID)
			cp, ok := t.Copies[copyID]
			if !ok {
				t.Copies[copyID] = times
			}
			t.Copies[copyID] = cp + times
			t.totalCardsCount += times
			fmt.Printf("card %d has %d copies; have %d total cards\n", copyID, cp+times, t.totalCardsCount)
		}
	}

	for i := 0; i < len(t.UniqueCards); i++ {
		id := i + 1
		card, ok := t.UniqueCards[id]
		if !ok {
			panic(errors.New("unexpected card index out of range"))
		}

		// register all copies for the current card
		cp, ok := t.Copies[id]
		if !ok {
			cp = 0
		}
		registerCopies(id, card.MatchingCount(), cp+1)
	}

	return t.totalCardsCount
}

type ScratchCard struct {
	Winning map[int]struct{}
	Have    map[int]struct{}
	ID      int
}

func (sc ScratchCard) MatchingCount() (count int) {
	for k := range sc.Have {
		_, ok := sc.Winning[k]
		if ok {
			count++
		}
	}

	return
}

func ScratchCardFromLine(line string) (sc ScratchCard) {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		panic(errors.New("unexpected line format"))
	}

	idPart, numbers := parts[0], parts[1]
	nn := strings.Split(numbers, "|")
	if len(nn) != 2 {
		panic(errors.New("unexpected line format"))
	}

	ii := strings.Split(idPart, "Card")
	if len(ii) != 2 {
		panic(errors.New("unexpected line format"))
	}
	f := strings.Fields(ii[1])
	if len(f) != 1 {
		panic(errors.New("unexpected line format"))
	}

	v, err := strconv.Atoi(f[0])
	if err != nil {
		panic(err)
	}

	sc.ID = v

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

	t := ScratchCardTable{
		UniqueCards: make(map[int]ScratchCard),
		Copies:      make(map[int]int),
	}
	for _, line := range lines {
		card := ScratchCardFromLine(line)
		t.AddCard(card)
	}

	ans := t.Winnings()
	fmt.Println(ans)
}
