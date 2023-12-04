package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
)

func isNumber(num rune) bool {
	return num >= '0' && num <= '9'
}

func isSymbol(sym uint8) bool {
	return !isNumber(rune(sym)) && sym != '.'
}

func part1(schematic []string) int {
	sum := 0
	for y, line := range schematic {
		currNumber := 0
		hasAdjacent := false
		inNumber := false
		for x, char := range line {
			if isNumber(char) {
				hasAdjacent = hasAdjacent ||
					(x > 0 && line[x-1] != '.' && !inNumber) || // left
					(x < len(line)-1 && isSymbol(line[x+1])) || // right
					(y > 0 && isSymbol(schematic[y-1][x])) || // top
					(y < len(schematic)-1 && isSymbol(schematic[y+1][x])) || // bottom
					(x > 0 && y > 0 && !inNumber && schematic[y-1][x-1] != '.') || // diagonal left top
					(x > 0 && y < len(schematic)-1 && !inNumber && isSymbol(schematic[y+1][x-1])) || // diagonal left bottom
					(x < len(line)-1 && y > 0 && isSymbol(schematic[y-1][x+1])) || // diagonal right top
					(x < len(line)-1 && y < len(schematic)-1 && isSymbol(schematic[y+1][x+1])) // diagonal right bottom

				inNumber = true
				currNumber = currNumber*10 + int(char-'0')
			}
			if !isNumber(char) || x == len(line)-1 {
				if hasAdjacent {
					sum += currNumber
					fmt.Println(currNumber)
				}
				inNumber = false
				hasAdjacent = false
				currNumber = 0
			}
		}
	}
	return sum
}

type EntityType uint8

const (
	Nothing EntityType = iota
	Gear
	PartNumber
)

type Entity struct {
	value int
	eType EntityType
}

func (e *Entity) String() string {
	switch e.eType {
	case Gear:
		return "*"
	case Nothing:
		return "o"
	default:
		return "#"
	}

}

func (pn *Entity) Len() int {
	return int(math.Log10(float64(pn.value))) + 1
}

func part2(inputLines []string) int {
	schematic := make([][]Entity, len(inputLines))
	gearXs := make([]int, 0)
	gearYs := make([]int, 0)

	for y, line := range inputLines {
		currNumber := 0
		inNumber := false
		schematic[y] = make([]Entity, len(line))
		for x, char := range line {
			if isNumber(char) {
				currNumber = currNumber*10 + int(char-'0')
				inNumber = true
			}
			if !isNumber(char) || x == len(line)-1 {
				if inNumber {
					inNumber = false
					entity := Entity{
						value: currNumber,
						eType: PartNumber,
					}
					for i := x - entity.Len(); i < x; i++ {
						schematic[y][i] = entity
					}
					currNumber = 0
				}
			}
			if char == '*' {
				schematic[y][x] = Entity{
					value: 0,
					eType: Gear,
				}
				gearXs = append(gearXs, x)
				gearYs = append(gearYs, y)
			}
		}
	}

	sum := 0
	for i, x := range gearXs {
		y := gearYs[i]
		entities := make([]Entity, 0)
		if x > 0 { // left
			e := schematic[y][x-1]
			if e.eType == PartNumber {
				if !slices.Contains(entities, e) {
					entities = append(entities, e)
				}
			}
		}
		if x < len(schematic[y])-1 { // right
			e := schematic[y][x+1]
			if e.eType == PartNumber {
				if !slices.Contains(entities, e) {
					entities = append(entities, e)
				}
			}
		}
		if y > 0 { // top
			e := schematic[y-1][x]
			if e.eType == PartNumber {
				if !slices.Contains(entities, e) {
					entities = append(entities, e)
				}
			}
		}
		if y < len(schematic)-1 { // bottom
			e := schematic[y+1][x]
			if e.eType == PartNumber {
				if !slices.Contains(entities, e) {
					entities = append(entities, e)
				}
			}
		}
		if x > 0 && y > 0 { // top left
			e := schematic[y-1][x-1]
			if e.eType == PartNumber {
				if !slices.Contains(entities, e) {
					entities = append(entities, e)
				}
			}
		}
		if x > 0 && y < len(schematic)-1 { // bottom left
			e := schematic[y+1][x-1]
			if e.eType == PartNumber {
				if !slices.Contains(entities, e) {
					entities = append(entities, e)
				}
			}
		}
		if x < len(schematic[y])-1 && y > 0 { // top right
			e := schematic[y-1][x+1]
			if e.eType == PartNumber {
				if !slices.Contains(entities, e) {
					entities = append(entities, e)
				}
			}
		}

		if x < len(schematic[y])-1 && y < len(schematic)-1 { // bottom right
			e := schematic[y+1][x+1]
			if e.eType == PartNumber {
				if !slices.Contains(entities, e) {
					entities = append(entities, e)
				}
			}
		}

		if len(entities) == 2 {
			partNumber := entities[0].value * entities[1].value
			sum += partNumber
		}

	}

	return sum
}

func main() {
	file, err := os.Open("03/input.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	schematic := make([]string, 0)

	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		schematic = append(schematic, scanner.Text())
	}

	fmt.Printf("Sum: %d\n", part2(schematic))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
