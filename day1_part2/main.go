package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Trie struct {
	Children map[string]*Trie
	Value    string
}

func (t *Trie) Insert(s string) {
	cur := t
	for _, r := range s {
		ss := string(r)

		v, ok := cur.Children[ss]
		if ok {
			cur = v
		} else {
			cur.Children[ss] = &Trie{
				Value:    s,
				Children: make(map[string]*Trie),
			}
			cur = cur.Children[ss]
		}
	}
}

func (t Trie) Has(s string) bool {
	cur := &t
	for _, r := range s {
		ss := string(r)

		v, ok := cur.Children[ss]
		if !ok {
			return false
		}

		cur = v
	}

	return true
}

type LineParser struct {
	Trie   *Trie
	Digits map[string]string
}

func (p LineParser) Parse(line string) int {
	var first, last string

	for left, right := 0, 1; right <= len(line); {
		window := line[left:right]
		// fmt.Println(right, len(line))
		// fmt.Printf("spelled digit check %s %d %d %d ", window, left, right, len(line))
		// current window not in trie, window is too small; reset window
		if !p.Trie.Has(window) && right-left <= 1 {
			// fmt.Println("current window not in trie, window is too small; reset window")
			right++
			left++
			continue
		}

		// current window not in trie, window can be closed; close window
		if !p.Trie.Has(window) && right-left > 1 {
			// fmt.Println("current window not in trie, window can be closed; close window")
			left++
			continue
		}

		v, ok := p.Digits[window]
		// current window in trie, is not a digit; open window
		if !ok {
			// fmt.Println("current window in trie, is not a digit; open window")
			right++
			continue
		}

		// current window in trie, is a digit; reset window, assign
		// fmt.Println("current window in trie, is a digit; reset window, assign")
		if first == "" {
			first = v
			right++
			left = right - 1
			continue
		}
		if first != "" {
			last = v
			right++
			left = right - 1
			continue
		}

		panic(errors.New("something has gone terribly wrong"))
	}

	if last == "" {
		last = first
	}

	i, err := strconv.Atoi(first + last)
	if err != nil {
		panic(err)
	}

	return i
}

var digits map[string]string = map[string]string{
	"1":     "1",
	"2":     "2",
	"3":     "3",
	"4":     "4",
	"5":     "5",
	"6":     "6",
	"7":     "7",
	"8":     "8",
	"9":     "9",
	"one":   "1",
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
}

func main() {
	// build trie
	trie := &Trie{Children: make(map[string]*Trie)}
	for k := range digits {
		trie.Insert(k)
	}

	f, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(f), "\n")
	var ans int
	for _, line := range lines {
		p := LineParser{
			Trie:   trie,
			Digits: digits,
		}

		v := p.Parse(line)
		fmt.Printf("%s %d\n", line, v)
		ans += v
	}

	fmt.Println(ans)
}
