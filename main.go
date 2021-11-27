package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("expected filename to be the first argument to the program")
	}

	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error while opening file: %v", err)
	}

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("error while reading csv from given file: %v", err)
	}

	distances := map[string]map[string]int{}
	for _, record := range records {
		distance, err := strconv.Atoi(record[2])
		if err != nil {
			log.Fatalf("error while reading distance from csv file: %v", err)
		}

		from, to := record[0], record[1]
		if _, ok := distances[from]; !ok {
			distances[from] = map[string]int{}
		}

		if _, ok := distances[from][to]; ok {
			distances[from][to] = min(distances[from][to], distance)
		} else {
			distances[from][to] = distance
		}
	}

	jpnagar := os.Args[2]
	indiranagar := os.Args[3]

	beenHere := map[string]bool{}
	reachable, via, cost := traverseDepth(beenHere, distances, jpnagar, indiranagar)
	fmt.Printf("Reachable?: %t, Via: %s, Cost (time): %d\n", reachable, via, cost)
}

func traverseDepth(beenHere map[string]bool, distances map[string]map[string]int, from string, destination string) (bool, string, int) {
	beenHere[from] = true

	comparisons := map[string]int{}
	for to, cost := range distances[from] {
		if ok := beenHere[to]; ok {
			continue
		}

		if to == destination {
			comparisons[to] = cost
		} else {
			reachable, _, newCost := traverseDepth(beenHere, distances, to, destination)
			if !reachable {
				continue
			}

			comparisons[to] = cost + newCost
		}
	}

	if len(comparisons) == 0 {
		return false, "", 0
	}

	bestpath, bestcost, err := minMap(comparisons)
	if err != nil {
		log.Fatalf("error with minMap: %v", err)
	}

	return true, bestpath, bestcost
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func minMap(m map[string]int) (string, int, error) {
	if len(m) == 0 {
		return "", 0, errors.New("length of map m is zero")
	}

	minimumKey := ""
	minimumSoFar := 9223372036854775807

	for k, v := range m {
		if v < minimumSoFar {
			minimumKey = k
			minimumSoFar = v
		}
	}

	return minimumKey, minimumSoFar, nil
}
