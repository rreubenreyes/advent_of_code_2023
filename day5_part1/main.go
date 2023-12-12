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

func (m ResourceMap) Lookup(src int) int {
	for _, r := range m.Ranges {
		// fmt.Printf("is key %d in range [%d, %d)?: ", src, r.SourceStart, r.SourceStart+r.RangeLength)
		if src >= r.SourceStart && src < r.SourceStart+r.RangeLength {
			// fmt.Print("yes; ")
			offset := r.DestinationStart - r.SourceStart
			// fmt.Printf("dest=%d\n", src+offset)

			return src + offset
		}
		// fmt.Println("no")
	}

	// fmt.Printf("key %d not in any range; dest=%d\n", src, src)
	return src
}

type Almanac struct {
	items map[struct{ SourceCategory, DestinationCategory string }]ResourceMap
	Seeds []int
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
			seeds := strings.Fields(line)
			for _, s := range seeds {
				v, err := strconv.Atoi(s)
				if err == nil {
					a.Seeds = append(a.Seeds, v)
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

	// I feel like part 2 is going to be a ridiculously long transitive lookup
	seedToSoil, _ := almanac.ResourceMap("seed", "soil")
	soilToFertilizer, _ := almanac.ResourceMap("soil", "fertilizer")
	fertilizerToWater, _ := almanac.ResourceMap("fertilizer", "water")
	waterToLight, _ := almanac.ResourceMap("water", "light")
	lightToTemperature, _ := almanac.ResourceMap("light", "temperature")
	temperatureToHumidity, _ := almanac.ResourceMap("temperature", "humidity")
	humidityToLocation, _ := almanac.ResourceMap("humidity", "location")

	var ans int
	for i, seed := range almanac.Seeds {
		soil := seedToSoil.Lookup(seed)
		fertilizer := soilToFertilizer.Lookup(soil)
		water := fertilizerToWater.Lookup(fertilizer)
		light := waterToLight.Lookup(water)
		temperature := lightToTemperature.Lookup(light)
		humidity := temperatureToHumidity.Lookup(temperature)
		location := humidityToLocation.Lookup(humidity)

		if i == 0 {
			ans = min(location, location)
		} else {
			ans = min(ans, location)
		}
	}

	fmt.Println(ans)
}
