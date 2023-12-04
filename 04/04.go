package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Card struct {
	num    int
	win    []int
	have   []int
	copies int
}

func (c Card) Points() int {
	points := 0
	for _, e := range c.have {
		if slices.Contains(c.win, e) {
			points += 1
		}
	}
	return points
}

var gameRegex = regexp.MustCompile("Card *(?P<num>\\d+): *(?P<win>(\\d+ *)+)\\| *(?P<have>(\\d+ *)+)")

func readCard(line string) Card {
	match := gameRegex.FindStringSubmatch(line)
	card := Card{copies: 1}
	for i, name := range gameRegex.SubexpNames() {
		switch name {
		case "num":
			card.num, _ = strconv.Atoi(match[i])
		case "win":
			win := strings.Fields(match[i])
			card.win = make([]int, 0)
			for _, str := range win {
				n, _ := strconv.Atoi(str)
				card.win = append(card.win, n)
			}
		case "have":
			have := strings.Fields(match[i])
			card.have = make([]int, 0)
			for _, str := range have {
				n, _ := strconv.Atoi(str)
				card.have = append(card.have, n)
			}
		}
	}
	return card
}

func part1(line string) int {
	card := readCard(line)
	fmt.Printf("card: %v\n", card)
	points := 0
	first := true
	for _, e := range card.have {
		if slices.Contains(card.win, e) {
			if first {
				first = false
				points += 1
			} else {
				points *= 2
			}
		}
	}
	return points
}

type Deck struct {
	cards []Card
}

func (d *Deck) String() string {
	var builder strings.Builder
	for _, card := range d.cards {
		fmt.Fprintf(&builder, "%v\n", card)
	}
	return builder.String()
}

func part2(deck *Deck) int {
	for i, c := range deck.cards {
		for j := 1; j <= c.Points(); j++ {
			if j >= len(deck.cards) {
				break
			}
			deck.cards[i+j].copies += c.copies
		}
	}
	fmt.Println(deck.String())
	sum := 0
	for _, card := range deck.cards {
		sum += card.copies
	}

	return sum
}

func main() {
	file, err := os.Open("04/input.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	deck := Deck{}
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		deck.cards = append(deck.cards, readCard(scanner.Text()))
	}
	fmt.Println(part2(&deck))

}
