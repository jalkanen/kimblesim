package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var games = flag.Int("games", 10000, "How many games to run")
var out = flag.String("out", "results.txt", "Result file")
var verbose = flag.Bool("verbose", false, "Should output all of the game contents?")

const (
	Popomatic = true
)

func main() {
	rand.Seed(time.Now().UnixNano())

	outFile, err := os.Create(*out)

	if err != nil {
		fmt.Errorf("Cannot open output file: %s", err.Error())
		os.Exit(1)
	}

	defer outFile.Close()

	for i := 0; i < *games; i++ {
		winner := play(i)

		fmt.Fprintln(outFile, winner.Name())
	}
}

func play(game int) Color {
	board := NewBoard()
	players := []Player{
		NewRandom(Red),
		NewRandom(Purple),
		NewRandom(White),
		NewRandom(Blue),
	}
	round := 0

	if *verbose {
		fmt.Println("Starting game simulation...")
	} else {
		fmt.Printf(".")
		if game%100 == 99 {
			fmt.Printf("\n")
		}
	}

	for {
		for ply := 1; ply <= Players; ply++ {
			player := players[ply-1]
			loop := true
			var roll int = 1

			for loop {
				roll = Roll(roll)

				player.Move(board, int(roll))

				if *verbose {
					fmt.Printf("ROUND %d, PLAYER %s, ROLL %d\n", round, player.Color().Name(), roll)
					fmt.Println(board)
				}

				// Check victory
				home := board.home[ply-1]
				if home[0] != Unoccupied && home[1] != Unoccupied && home[2] != Unoccupied && home[3] != Unoccupied {
					if *verbose {
						fmt.Printf("WINNER: %s\n", player.Color().Name())
					}
					return player.Color()
				}

				// Is the user allowed to continue?
				if roll == 6 {
					loop = true
				} else {
					loop = false
				}
			}
		}
		round++

	}

}

var unpairs [6][4]int = [6][4]int{
	{2, 3, 4, 5}, // 1
	{1, 3, 4, 6}, // 2
	{1, 2, 5, 6}, // 3
	{1, 2, 5, 6}, // 4
	{1, 3, 4, 6}, // 5
	{2, 3, 4, 5}, // 6
}

func Roll(prev int) int {
	var roll int

	if Popomatic {
		// Utilizes the Popomatic statistical results from http://statistition.com/?p=440
		r := rand.Float32()
		if r < 0.239 {
			roll = 7 - prev
		} else if r > (1.0 - 0.108) {
			roll = prev
		} else {
			roll = unpairs[prev-1][rand.Int31n(4)]
		}
	} else {
		roll = int(rand.Int31n(6)) + 1
	}

	return roll
}
