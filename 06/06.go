package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func part1(file *os.File) {
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	timeLine := strings.Fields(scanner.Text())[1:]
	times := make([]int, len(timeLine))
	for i := 0; i < len(timeLine); i++ {
		times[i], _ = strconv.Atoi(timeLine[i])
	}
	scanner.Scan()
	distanceLine := strings.Fields(scanner.Text())[1:]
	distances := make([]int, len(distanceLine))
	for i := 0; i < len(distanceLine); i++ {
		distances[i], _ = strconv.Atoi(distanceLine[i])
	}

	fmt.Printf("Times    : %v\n", times)
	fmt.Printf("Distances: %v\n", distances)

	product := 1
	for i := 0; i < len(times); i++ {
		minDist := distances[i]
		numberOfWays := 0
		for timePressed := 0; timePressed <= times[i]; timePressed++ {
			dist := (times[i] - timePressed) * timePressed
			if dist > minDist {
				numberOfWays++
			}
		}
		fmt.Printf("Number of ways: %d\n", numberOfWays)
		product = product * numberOfWays
	}
	fmt.Printf("Product: %d\n", product)
}

func main() {
	file, err := os.Open("06/input.txt")
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
