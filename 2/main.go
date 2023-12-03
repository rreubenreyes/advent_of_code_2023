package main

import (
	"fmt"
	"os"
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

		_, ok := cur.Children[ss]
		if !ok {
			return false
		}
	}

	return true
}

type CalibrationValues struct {
	Digits *Trie
	Lines  []string
}

func (v CalibrationValues) Parse() int {
	// TODO: implement
	return -1
}

func main() {
	// build trie
	trie := &Trie{}
	digits := []string{
		"one", "two", "three",
		"four", "five", "six",
		"seven", "eight", "nine",
	}
	for _, digit := range digits {
		trie.Insert(digit)
	}

	f, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(f), "\n")

	v := CalibrationValues{
		Lines:  lines,
		Digits: trie,
	}

	fmt.Println(v.Parse())
}
