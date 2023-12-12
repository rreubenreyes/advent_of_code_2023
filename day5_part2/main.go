package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Range struct {
	DestinationStart, SourceStart, RangeLength int
}

func (r Range) MaybeMapRange(start, end int) (ranges [][]int) {
	offset := r.DestinationStart - r.SourceStart
	appendRange := func(s, e int) {
		if e-s > 0 {
			ranges = append(ranges, []int{s, e})
		}
	}

	var destStart, destEnd int
	if start >= r.SourceStart && end < r.SourceStart+r.RangeLength {
		// "inner join"
		// fmt.Printf("[%d, %d) is fully within [%d, %d)\n", start, end, r.SourceStart, r.SourceStart+r.RangeLength)
		destStart = start + offset
		destEnd = min(end+offset, r.SourceStart+r.RangeLength)

		appendRange(destStart, destEnd)
	} else if end >= r.SourceStart && end < r.SourceStart+r.RangeLength {
		// "left join"
		// fmt.Printf("[%d, %d) is partially within [%d, %d); source range has a higher starting bound\n", start, end, r.SourceStart, r.SourceStart+r.RangeLength)
		destStart = r.DestinationStart
		destEnd = min(end+offset, r.SourceStart+r.RangeLength)

		appendRange(destStart, destEnd)
		appendRange(start, r.SourceStart)
	} else if start >= r.SourceStart && start < r.SourceStart+r.RangeLength {
		// "right join"
		// fmt.Printf("[%d, %d) is partially within [%d, %d); source range has a lower ending bound\n", start, end, r.SourceStart, r.SourceStart+r.RangeLength)
		destStart = start + offset
		destEnd = r.SourceStart + r.RangeLength + offset

		appendRange(destStart, destEnd)
		appendRange(r.SourceStart+r.RangeLength, end)
	}

	// fmt.Printf("MaybeMapRange result: %+v\n", ranges)
	return ranges
}

type ResourceMap struct {
	data   map[int]int
	Ranges []Range
}

func NewResourceMap() ResourceMap {
	return ResourceMap{
		data: make(map[int]int),
	}
}

func (m *ResourceMap) AddRange(destStart, srcStart, rangeLength int) {
	r := Range{
		DestinationStart: destStart,
		SourceStart:      srcStart,
		RangeLength:      rangeLength,
	}

	m.Ranges = append(m.Ranges, r)
}

func (m ResourceMap) LookupRanges(input [][]int) (result [][]int) {
	for _, i := range input {
		var thisResult [][]int
		var noMatchAlreadyAdded bool
		start, end := i[0], i[1]

		for _, r := range m.Ranges {
			rr := r.MaybeMapRange(start, end)
			thisResult = append(thisResult, rr...)
		}

		if len(thisResult) == 0 && !noMatchAlreadyAdded {
			// fmt.Printf("[%d, %d) not in any ranges\n", start, end)
			thisResult = append(thisResult, []int{start, end})
			noMatchAlreadyAdded = true
		}

		result = append(result, thisResult...)
	}

	return
}

type Almanac struct {
	items      map[struct{ SourceCategory, DestinationCategory string }]ResourceMap
	SeedRanges [][]int
}

func ParseCategories(s string) (src, dest string) {
	r := regexp.MustCompile(`([a-z]+)-to-([a-z]+) map:`)
	categories := r.FindStringSubmatch(s)[1:]
	src, dest = categories[0], categories[1]

	return
}

func (a Almanac) ResourceMap(src, dest string) (ResourceMap, bool) {
	key := struct {
		SourceCategory      string
		DestinationCategory string
	}{
		SourceCategory:      src,
		DestinationCategory: dest,
	}

	v, ok := a.items[key]
	if !ok {
		return NewResourceMap(), false
	}

	return v, true
}

func AlmanacFromFile(f string) (a Almanac) {
	b, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(b), "\n")

	a.items = make(map[struct{ SourceCategory, DestinationCategory string }]ResourceMap)
	var curSrc, curDest string
	var curMap ResourceMap
	for i, line := range lines {
		// first line is always seeds
		if i == 0 {
			seedRanges := strings.Fields(strings.Split(line, "seeds: ")[1])
			var curRangeStart int
			for i, s := range seedRanges {
				if i%2 > 0 {
					var rangeLength int
					v, err := strconv.Atoi(s)
					if err == nil {
						rangeLength = v
					}

					a.SeedRanges = append(a.SeedRanges, []int{curRangeStart, curRangeStart + rangeLength})
				} else {
					v, err := strconv.Atoi(s)
					if err == nil {
						curRangeStart = v
					}
				}
			}

			continue
		}

		// empty line
		if len(line) == 0 {
			continue
		}

		// first map declaration: parse the first map
		if strings.Contains(line, "map:") && curSrc == "" && curDest == "" {
			curSrc, curDest = ParseCategories(line)
			curMap = NewResourceMap()

			continue
		}

		// any other map declaration: save the previous map and parse the next one
		if strings.Contains(line, "map:") && curSrc != "" && curDest != "" {
			a.items[struct {
				SourceCategory      string
				DestinationCategory string
			}{
				SourceCategory:      curSrc,
				DestinationCategory: curDest,
			}] = curMap

			curSrc, curDest = ParseCategories(line)
			curMap = NewResourceMap()

			continue
		}

		// last line, commit whatever we have
		if i == len(lines)-1 {
			a.items[struct {
				SourceCategory      string
				DestinationCategory string
			}{
				SourceCategory:      curSrc,
				DestinationCategory: curDest,
			}] = curMap

			continue
		}

		// range declaration
		fields := strings.Fields(line)
		var d, s, r int
		for j, f := range fields {
			v, _ := strconv.Atoi(f)
			switch j {
			case 0:
				d = v
			case 1:
				s = v
			case 2:
				r = v
			}
		}

		curMap.AddRange(d, s, r)
	}

	return
}

func main() {
	almanac := AlmanacFromFile("./input.txt")

	ranges := almanac.SeedRanges
	seedToSoil, _ := almanac.ResourceMap("seed", "soil")
	soilToFertilizer, _ := almanac.ResourceMap("soil", "fertilizer")
	fertilizerToWater, _ := almanac.ResourceMap("fertilizer", "water")
	waterToLight, _ := almanac.ResourceMap("water", "light")
	lightToTemperature, _ := almanac.ResourceMap("light", "temperature")
	temperatureToHumidity, _ := almanac.ResourceMap("temperature", "humidity")
	humidityToLocation, _ := almanac.ResourceMap("humidity", "location")

	ranges = seedToSoil.LookupRanges(ranges)
	ranges = soilToFertilizer.LookupRanges(ranges)
	ranges = fertilizerToWater.LookupRanges(ranges)
	ranges = waterToLight.LookupRanges(ranges)
	ranges = lightToTemperature.LookupRanges(ranges)
	ranges = temperatureToHumidity.LookupRanges(ranges)
	ranges = humidityToLocation.LookupRanges(ranges)

	var ans int
	for i, r := range ranges {
		if i == 0 {
			ans = min(r[0], r[0])
		} else {
			ans = min(ans, r[0])
		}
	}
	fmt.Println(ans)
}
