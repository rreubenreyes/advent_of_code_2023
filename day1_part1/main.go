package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CalibrationValueFromLine(s string) int {
	var first, last string
	for _, ss := range s {
		_, err := strconv.Atoi(string(ss))
		if err == nil && first == "" {
			first = string(ss)
		}
		if err == nil && first != "" {
			last = string(ss)
		}
	}

	if last == "" {
		last = first
	}

	v, err := strconv.Atoi(first + last)
	if err != nil {
		panic(err)
	}

	return v
}

func main() {
	f, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(f), "\n")
	var ans int
	for _, line := range lines {
		v := CalibrationValueFromLine(line)
		ans += v
	}

	fmt.Println(ans)
}
