package tictactoe

// This struct stores the coordinates of the move and the score of the move, for the purpose of the minimax algorithm
// coordinates are defined in board.go, with an x and a y coordinate, both as integers
type evaluation struct {
	coordinates coord
	score       int
}

// Minimax searches all possible moves and searches for the best moves. The best move is rated as the either: the move/moves
// that win the fastest, the move/moves that lose the slowest, or the move that blocks an immediate win for the opponent.
func minimax(board [3][3]int, aiMark int, humanMark int, currMark int, depth int) evaluation {
	// Finds the best move, and the fastest win/slowest loss
	if win := winDeterminer(board, aiMark, humanMark); win != 0 {
		if win == 1 { // AI wins
			return evaluation{coord{x: -1, y: -1}, 15 - depth}
		} else if win == -1 { // AI loses (only when error is made before it can play)
			return evaluation{coord{x: -1, y: -1}, depth - 15}
		}
	} else if checkIfTie(board) {
		return evaluation{coord{x: -1, y: -1}, 0}
	}
	var threats []coord
	if currMark == aiMark {
		threats = detectImmediateThreats(board, humanMark)
	}
	// List of all possible moves, and their evaluations
	possibleMoves := isEmpty(board)
	testPlays := make([]evaluation, 0, len(possibleMoves))
	for _, move := range possibleMoves {
		var currPlay evaluation
		currPlay.coordinates.x, currPlay.coordinates.y = move.x, move.y
		board[move.x][move.y] = currMark
		// if this isn't a tie, otherwise recursively call the function again
		if winDeterminer(board, aiMark, humanMark) != 0 {
			// create a positive score if it is winning for the AI
			if currMark == aiMark {
				currPlay.score = 10 - depth
			} else { // create a negative score if it is losing for the AI
				currPlay.score = depth - 10
			}
		} else if currMark == aiMark {
			result := minimax(board, aiMark, humanMark, humanMark, depth+1)
			currPlay.score = result.score
		} else {
			result := minimax(board, aiMark, humanMark, aiMark, depth+1)
			currPlay.score = result.score
		}
		board[move.x][move.y] = 0
		testPlays = append(testPlays, currPlay)
	}
	// Finds the best move
	var bestTestPlay int
	if currMark == aiMark {
		bestScore := -10
		for i := 0; i < len(testPlays); i++ {
			if testPlays[i].score > bestScore || testPlays[i].score == bestScore && contains(threats, testPlays[i].coordinates) {
				bestScore = testPlays[i].score
				bestTestPlay = i
			}
		}
	} else {
		bestScore := 10
		for i := 0; i < len(testPlays); i++ {
			if testPlays[i].score < bestScore || testPlays[i].score == bestScore && contains(threats, testPlays[i].coordinates) {
				bestScore = testPlays[i].score
				bestTestPlay = i
			}
		}
	}
	return testPlays[bestTestPlay]
}

// winDeterminer returns the value based on the AI winning vs the human winning. It returns 1 for the AI winning,
// and -1 for the human winning. It returns 0 otherwise.
func winDeterminer(board [3][3]int, aiMark int, humanMark int) int {
	if checkIfWinner(board, aiMark) {
		return 1
	} else if checkIfWinner(board, humanMark) {
		return -1
	}
	return 0
}

// detectImmediateThreats finds the next possible winning moves for the opponent. This is important in situations where there are
// multiple winning moves for the opponent, to block at least one of them, and count on a mistake.
func detectImmediateThreats(board [3][3]int, opponentMark int) []coord {
	var threats []coord
	moves := isEmpty(board)
	for _, move := range moves {
		board[move.x][move.y] = opponentMark
		if checkIfWinner(board, opponentMark) {
			threats = append(threats, move)
		}
		board[move.x][move.y] = 0
	}
	return threats
}

// contains simply checks if the slice of coordinates contains the coordinates we are checking against,
// returning true if it is in the slice, and false if it is not in the slice.
func contains(slice []coord, element coord) bool {
	for _, v := range slice {
		if v.x == element.x && v.y == element.y {
			return true
		}
	}
	return false
}
