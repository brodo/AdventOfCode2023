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

	digits := []int{}
	sum := 0
	idx := 0
	for {
		r := content[idx]
		skipCount := 1
		switch {
		case r >= '0' && r <= '9':
			digits = append(digits, int(r-'0'))
		case r == '\n' || idx == len(content)-1: // end of line or end of file
			sum += digits[0]*10 + digits[len(digits)-1]
			digits = []int{}
		default:
			for i, word := range wordList {
				if idx+len(word) >= len(content) {
					continue
				}
				candidate := string(content[idx : idx+len(word)])
				if candidate == word {
					digits = append(digits, i+1)
					skipCount = len(word) - 1
					break
				}
			}
		}
		idx += skipCount
		if idx == len(content) {
			break
		}
	}
	fmt.Println("Sum:", sum)
}
