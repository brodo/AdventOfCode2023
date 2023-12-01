package main

import (
	"fmt"
	"os"
)

var wordList = []string{
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

func main() {
	content, err := os.ReadFile("./01/input.txt")
	if err != nil {
		fmt.Println("Err", err)
	}
	if len(content) == 0 {
		fmt.Println("Empty file")
		return
	}

	digits := []int{}
	sum := 0
	idx := 0
	for {
		r := content[idx]
		switch {
		case r >= '0' && r <= '9':
			digits = append(digits, int(r-'0'))
			idx++
		case r == '\n':
			first := digits[0]
			last := digits[len(digits)-1]
			lineValue := first*10 + last
			sum += lineValue
			digits = []int{}
			idx++
		default:
			skipCount := 1
			for i, word := range wordList {
				if idx+len(word) < len(content) {
					candidate := string(content[idx : idx+len(word)])
					if candidate == word {
						digits = append(digits, i+1)
						skipCount = len(word) - 1
						break
					}
				}
			}
			idx += skipCount

		}
		if idx == len(content) {
			break
		}
	}
	fmt.Println("Sum:", sum)

}
