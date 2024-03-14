package main

import (
	"chess/tictactoe"
	"flag"
)

func main() {
	// Flag to check if the game is against an AI or two people
	withAI := flag.Bool("ai", false, "Determines whether you play an AI.")
	flag.Parse()
	// Start the game
	tictactoe.Play(*withAI)
}
