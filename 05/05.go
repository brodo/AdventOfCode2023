package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	Seed        = "seed"
	Soil        = "soil"
	Fertilizer  = "fertilizer"
	Water       = "water"
	Light       = "light"
	Temperature = "temperature"
	Humidity    = "humidity"
	Location    = "location"
)

const (
	StartState = iota
	MapStart
	Map
	Gap
	StartGap
)

type SeedList []uint64

type Mapping struct {
	destRangeStart   uint64
	sourceRangeStart uint64
	rangeLength      uint64
}

type MappingList []Mapping

type Almanac struct {
	seedList SeedList
	maps     map[string]MappingList
	keys     []string
}

func (a Almanac) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "seeds: ")
	for _, s := range a.seedList {
		fmt.Fprintf(&builder, "%d ", s)
	}
	fmt.Fprintf(&builder, "\n\n")
	for _, key := range a.keys {
		fmt.Fprintf(&builder, "%s map:\n", key)
		for _, mapping := range a.maps[key] {
			fmt.Fprintf(&builder, "%d %d %d\n", mapping.destRangeStart, mapping.sourceRangeStart, mapping.rangeLength)
		}
		fmt.Fprintf(&builder, "\n")
	}

	return builder.String()
}

func (m *MappingList) Map(number uint64) uint64 {
	for _, mapping := range *m {
		if mapping.sourceRangeStart <= number && mapping.sourceRangeStart+mapping.rangeLength >= number {
			return mapping.destRangeStart + (number - mapping.sourceRangeStart)
		}
	}
	return number
}

func (a Almanac) CalculateSolutions() [][]uint64 {
	solutions := make([][]uint64, len(a.seedList))
	for i, s := range a.seedList {
		solutions[i] = make([]uint64, 7)
		currentRes := s
		for j, key := range a.keys {
			mappingList := a.maps[key]
			currentRes = mappingList.Map(currentRes)
			solutions[i][j] = currentRes
		}
	}
	return solutions
}

func toIntSlice(strs []string) []uint64 {
	res := make([]uint64, len(strs))
	for i, s := range strs {
		res[i], _ = strconv.ParseUint(s, 10, 64)
	}
	return res
}

func parseAlmanac(file *os.File) Almanac {
	scanner := bufio.NewScanner(file)
	state := StartState
	almanac := Almanac{
		seedList: make(SeedList, 0),
		maps:     make(map[string]MappingList),
		keys:     make([]string, 0),
	}
	currentSubject := ""
	mappings := make([]Mapping, 0)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		switch state {
		case StartState:
			strFields := strings.Fields(scanner.Text())[1:]
			almanac.seedList = toIntSlice(strFields)
			scanner.Scan()
			state = MapStart
			continue
		case MapStart:
			currentSubject = strings.Fields(scanner.Text())[0]
			state = Map
			continue
		case Map:
			text := scanner.Text()
			if text == "" {
				almanac.maps[currentSubject] = mappings
				almanac.keys = append(almanac.keys, currentSubject)
				mappings = make([]Mapping, 0)
				state = MapStart
				continue
			}
			ints := toIntSlice(strings.Fields(text))
			mappings = append(mappings, Mapping{
				destRangeStart:   ints[0],
				sourceRangeStart: ints[1],
				rangeLength:      ints[2],
			})
		}
	}
	almanac.maps[currentSubject] = mappings
	almanac.keys = append(almanac.keys, currentSubject)
	mappings = make([]Mapping, 0)
	return almanac
}

func part1(file *os.File) {
	almanac := parseAlmanac(file)
	solutions := almanac.CalculateSolutions()
	min := uint64(math.MaxUint64)
	for i, solution := range solutions {
		fmt.Printf("Solution for %d: %v\n", almanac.seedList[i], solution)
		if solution[6] < min {
			min = solution[6]
		}
	}
	println(min)

}

func main() {
	file, err := os.Open("05/input.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	part1(file)

}
