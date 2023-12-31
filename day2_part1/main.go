package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Pull struct {
	Blue, Red, Green int
}

type Game struct {
	Pulls []Pull
	ID    int
}

func GameFromLine(line string) (g Game) {
	// parse ID
	gameExpr := regexp.MustCompile(`Game (?P<ID>\d+):`)
	gameSubexprs := gameExpr.FindStringSubmatch(line)
	if len(gameSubexprs) != 2 {
		panic("unexpected input")
	}

	i, err := strconv.Atoi(gameSubexprs[1])
	if err != nil {
		panic("unexpected input")
	}

	g.ID = i

	// parse block counts
	pullExpr := regexp.MustCompile(`(?P<Count>[0-9]+) (?P<Color>\w+)`)
	pullsStr := strings.Split(line, ":")[1]
	pulls := strings.Split(pullsStr, ";")
	g.Pulls = make([]Pull, len(pulls))
	for j, p := range pulls {
		p = strings.Trim(p, " ")
		allSubexprs := pullExpr.FindAllStringSubmatch(p, -1)

		for _, s := range allSubexprs {
			count, color := s[1], s[2]
			jj, err := strconv.Atoi(count)
			if err != nil {
				panic("unexpected input")
			}

			switch color {
			case "blue":
				g.Pulls[j].Blue = jj
			case "green":
				g.Pulls[j].Green = jj
			case "red":
				g.Pulls[j].Red = jj
			}
		}
	}

	return
}

func main() {
	f, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(f), "\n")
	var ans int
	for _, line := range lines {
		game := GameFromLine(line)

		isPossible := true
		for _, pull := range game.Pulls {
			if pull.Blue > 14 || pull.Green > 13 || pull.Red > 12 {
				isPossible = false
				break
			}
		}
		if isPossible {
			ans += game.ID
		}
	}
	fmt.Println(ans)
}
