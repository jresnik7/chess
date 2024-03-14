package tictactoe

// This struct stores the coordinates of the move and the score of the move, for the purpose of the minimax algorithm
// coordinates are defined in board.go, with an x and a y coordinate, both as integers
type evaluation struct {
	coordinates coord
	score       int
}

func minimax(board [3][3]int, aiMark int, humanMark int, currMark int, depth int) evaluation {
	// Finds the best move, and the fastest win/slowest loss
	if win := WinDeterminer(board, aiMark, humanMark); win != 0 {
		if win == 1 { // AI wins
			return evaluation{coord{x: -1, y: -1}, 15 - depth}
		} else if win == -1 { // AI loses (only when error is made before it can play)
			return evaluation{coord{x: -1, y: -1}, depth - 15}
		}
	} else if CheckIfTie(board) {
		return evaluation{coord{x: -1, y: -1}, 0}
	}
	var threats []coord
	if currMark == aiMark {
		threats = detectImmediateThreats(board, humanMark)
	}
	// List of all possible moves, and their evaluations
	possibleMoves := IsEmpty(board)
	testPlays := make([]evaluation, 0, len(possibleMoves))
	for _, move := range possibleMoves {
		var currPlay evaluation
		currPlay.coordinates.x, currPlay.coordinates.y = move.x, move.y
		board[move.x][move.y] = currMark
		if WinDeterminer(board, aiMark, humanMark) != 0 {
			if currMark == aiMark {
				currPlay.score = 10 - depth
			} else {

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

// Determines the value based on the AI winning vs the human winning
func WinDeterminer(board [3][3]int, aiMark int, humanMark int) int {
	if CheckIfWinner(board, aiMark) {
		return 1
	} else if CheckIfWinner(board, humanMark) {
		return -1
	}
	return 0
}

func detectImmediateThreats(board [3][3]int, opponentMark int) []coord {
	var threats []coord
	moves := IsEmpty(board)
	for _, move := range moves {
		board[move.x][move.y] = opponentMark
		if CheckIfWinner(board, opponentMark) {
			threats = append(threats, move)
		}
		board[move.x][move.y] = 0
	}
	return threats
}

func contains(slice []coord, element coord) bool {
	for _, v := range slice {
		if v.x == element.x && v.y == element.y {
			return true
		}
	}
	return false
}
