package main

import (
	"fmt"
	"os"
	"strings"
)

func ReadLines(f string) []string {
	b, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}

	return strings.Split(string(b), "\n")
}

func main() {
	lines := ReadLines("./sample_input.txt")
	for _, line := range lines {
		// do something
	}

	var ans int
	fmt.Println(ans)
}
