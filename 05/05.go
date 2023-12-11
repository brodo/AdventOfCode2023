package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
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
)

type SeedList []uint64

type Mapping struct {
	destRangeStart   uint64
	sourceRangeStart uint64
	rangeLength      uint64
}

func (m *Mapping) destRangeEnd() uint64 {
	return m.destRangeStart + m.rangeLength - 1
}

func (m *Mapping) sourceRangeEnd() uint64 {
	return m.sourceRangeStart + m.rangeLength - 1
}

func (m Mapping) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "(Src: %d-%d, Dst: %d-%d)", m.sourceRangeStart, m.sourceRangeEnd(), m.destRangeStart, m.destRangeEnd())
	return builder.String()
}

type MappingList []Mapping

type Almanac struct {
	seedList SeedList
	maps     map[string]MappingList
	keys     []string
}

type NumRange struct {
	start  uint64
	length uint64
}

func (r NumRange) end() uint64 {
	return r.start + r.length - 1
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

func (a *Almanac) CalculateSolutions() [][]uint64 {
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

type Border struct {
	num     uint64
	isInput bool
}

func (a *Almanac) FindMinimum() uint64 {
	numRanges := make([]NumRange, 0)
	for i := 0; i < len(a.seedList); i += 2 {
		numRanges = append(numRanges, NumRange{
			start:  a.seedList[i],
			length: a.seedList[i+1],
		})
	}

	results := make(MappingList, len(numRanges))
	// first step is to fill the results with the initial ranges
	for i, r := range numRanges {
		results[i] = Mapping{
			destRangeStart:   r.start,
			sourceRangeStart: r.start,
			rangeLength:      r.length,
		}
	}

	// iterate over all stages
	for _, key := range a.keys {
		currentStage := a.maps[key]
		borders := make([]Border, 0)
		for _, r := range results {
			borders = append(borders, Border{
				num:     r.sourceRangeEnd(),
				isInput: true,
			}, Border{
				num:     r.sourceRangeStart,
				isInput: true,
			})
		}
		for _, r := range currentStage {
			borders = append(borders, Border{
				num:     r.sourceRangeStart,
				isInput: false,
			}, Border{
				num:     r.sourceRangeEnd(),
				isInput: false,
			})
		}
		slices.SortFunc(borders, func(a, b Border) int {
			return cmp.Compare(a.num, b.num)
		})

		firstSrcPos := -1
		lastSrcPos := -1

		for i := 0; i < len(borders); i++ {
			if borders[i].isInput && firstSrcPos == -1 {
				firstSrcPos = i
				break
			}
		}

		for i := len(borders) - 1; i > 0; i-- {
			if borders[i].isInput && lastSrcPos == -1 {
				lastSrcPos = i
				break
			}
		}

		newResults := make(MappingList, 0)

		for i := firstSrcPos; i < lastSrcPos; i += 2 {
			border := borders[i]
			srs := border.num

			drs := currentStage.Map(results.Map(border.num))
			length := (borders[i+1].num - border.num) + 1
			if length == 0 {
				continue
			}

			newResults = append(newResults, Mapping{
				destRangeStart:   drs,
				sourceRangeStart: srs,
				rangeLength:      length,
			})
		}

		fmt.Printf("results after %s: %v\n", key, newResults)
		results = newResults

	}
	fmt.Printf("results: %v\n", results)

	return results[0].destRangeStart

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
	m := uint64(math.MaxUint64)
	for i, solution := range solutions {
		fmt.Printf("Solution for %d: %v\n", almanac.seedList[i], solution)
		if solution[6] < m {
			m = solution[6]
		}
	}
	println(m)

}

func part2(file *os.File) {
	almanac := parseAlmanac(file)
	fmt.Println(almanac.FindMinimum())
}

func main() {
	file, err := os.Open("05/example.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	part2(file)

}
