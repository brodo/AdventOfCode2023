package main

import (
	"fmt"
	"os"
)

type Color uint8

const (
	Red = iota
	Green
	Blue
)

type Cubes struct {
	count int
	color Color
}

type CubeSet []Cubes

type Game struct {
	id   int
	sets []CubeSet
}

var colorMap = map[string]Color{
	"red":   Red,
	"green": Green,
	"blue":  Blue,
}

func parseLine(line string) Game {
	//Game 9: 5 blue, 1 green, 4 red; 2 green, 6 red, 12 blue; 2 green, 7 blue, 1 red; 12 blue, 2 green, 1 red
	//Game 10: 1 red, 16 blue, 18 green; 14 green, 13 blue; 4 green, 7 blue; 5 red, 16 blue, 11 green; 14 green, 2 red, 5 blue; 10 blue, 3 red, 6 green

	game := Game{}
	game.sets = []CubeSet{}
	idx := 5

	for line[idx] >= '0' && line[idx] <= '9' {
		game.id = game.id*10 + int(line[idx]-'0')
		idx += 1
	}

	idx += 2 // skip colon and space

	currNum := 0
	currCubeSet := CubeSet{}
	for idx < len(line) {
		for line[idx] >= '0' && line[idx] <= '9' {
			currNum = currNum*10 + int(line[idx]-'0')
			idx += 1
		}
		idx++ // skip space
		for k := range colorMap {
			if idx+len(k) > len(line) {
				continue
			}
			if line[idx:idx+len(k)] == k {
				currCubeSet = append(currCubeSet, Cubes{count: currNum, color: colorMap[k]})
				currNum = 0
				idx += len(k)
				break
			}
		}

		if len(line) <= idx || line[idx] == ';' {
			game.sets = append(game.sets, currCubeSet)
			currCubeSet = CubeSet{}
			idx += 2
		} else if line[idx] == ',' {
			idx += 2 // skip comma and space
		}
	}

	return game
}

func checkIfPossible(game Game) bool {
	maxRed := 12
	maxGreen := 13
	maxBlue := 14
	for _, cubeSet := range game.sets {
		for _, c := range cubeSet {
			switch c.color {
			case Red:
				if c.count > maxRed {
					return false
				}
			case Green:
				if c.count > maxGreen {
					return false
				}
			case Blue:
				if c.count > maxBlue {
					return false
				}
			}
		}

	}
	return true
}

func calculatePower(game Game) int {
	minRed := 0
	minGreen := 0
	minBlue := 0
	for _, cubeSet := range game.sets {
		for _, c := range cubeSet {
			switch c.color {
			case Red:
				minRed = max(minRed, c.count)
			case Green:
				minGreen = max(minGreen, c.count)
			case Blue:
				minBlue = max(minBlue, c.count)
			}
		}

	}
	return minGreen * minRed * minBlue
}

func main() {
	content, err := os.ReadFile("./02/input.txt")
	if err != nil {
		fmt.Println("Err", err)
	}

	lastNewLine := 0
	possibleSum := 0
	powerSum := 0
	for i, r := range content {
		if r == '\n' || i == len(content)-1 {
			game := parseLine(string(content[lastNewLine:i]))
			if checkIfPossible(game) {
				possibleSum += game.id
				fmt.Printf("Game: %d is possible\n", game.id)
			}
			power := calculatePower(game)
			fmt.Printf("Game: %d has power %d\n", game.id, power)
			powerSum += power
			lastNewLine = i + 1
		}
	}
	fmt.Printf("Number of possible games: %d\n", possibleSum)
	fmt.Printf("Power sum: %d\n", powerSum)

}
