package tictactoe

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// This struct represents coordinates for input onto the board.
type coord struct {
	x, y int
}

// This is the main loop for the game, which plays the game itself. This allows the command line argument of -ai, if the player
// wants to play against an unbeatable AI, and learn the best strategy. Without this flag, the game assumes 2 players.
func Play(withAI bool) {
	// Start determines whether the AI or human starts
	var start int = rand.Int() % 2
	var humanMark, aiMark int
	// Declare all necessary variables
	board := [3][3]int{}
	var turn int = 1
	var x, y int
	var winner, tie bool
	var moves int = 0
	// Determines who begins, sets the marks accordingly
	if start == 0 && withAI {
		fmt.Println("The human player will begin!")
		humanMark = -1
		aiMark = 1
	} else if withAI {
		fmt.Println("The AI will begin!")
		aiMark = -1
		humanMark = 1
	}
	for {
		if withAI {
			// The human begins
			if start == 0 {
				// Only print the board again if it is the first move, as the board hasn't been displayed yet
				if moves == 0 {
					printBoard(board)
				}
				x, y = getCheckInput(board, turn)
				board = updateBoard(board, x, y, turn)
				moves += 1
				winner = checkIfWinner(board, 1) || checkIfWinner(board, -1)
				tie = checkIfTie(board)
				if winner || tie {
					break
				}
				turn = updateTurn(turn)
				optimalMove := minimax(board, aiMark, humanMark, aiMark, moves)
				fmt.Println("The AI has made the following move: ")
				board = updateBoard(board, optimalMove.coordinates.x, optimalMove.coordinates.y, turn)
				moves += 1
				winner = checkIfWinner(board, 1) || checkIfWinner(board, -1)
				tie = checkIfTie(board)
				if winner || tie {
					break
				}
				turn = updateTurn(turn)
			} else { // The AI begins
				optimalMove := minimax(board, aiMark, humanMark, aiMark, moves)
				fmt.Println("The AI has made the following move: ")
				board = updateBoard(board, optimalMove.coordinates.x, optimalMove.coordinates.y, turn)
				moves += 1
				winner = checkIfWinner(board, 1) || checkIfWinner(board, -1)
				tie = checkIfTie(board)
				if winner || tie {
					break
				}
				turn = updateTurn(turn)
				x, y = getCheckInput(board, turn)
				board = updateBoard(board, x, y, turn)
				moves += 1
				winner = checkIfWinner(board, 1) || checkIfWinner(board, -1)
				tie = checkIfTie(board)
				if winner || tie {
					break
				}
				turn = updateTurn(turn)
			}
		} else { // Two human players
			if moves == 0 {
				printBoard(board)
			}
			x, y = getCheckInput(board, turn)
			board = updateBoard(board, x, y, turn)
			moves += 1
			winner = checkIfWinner(board, 1) || checkIfWinner(board, -1)
			tie = checkIfTie(board)
			if winner || tie {
				break
			}
			turn = updateTurn(turn)
		}
	}
	postGameOutput(winner, tie, withAI, turn, start)
}

// printBoard prints the board to the command line, replacing -1s with Xs and 1s with Os,
// while leaving blank cells (0s) filled with spaces
func printBoard(board [3][3]int) {
	for i, row := range board {
		if i != 0 {
			fmt.Println(strings.Repeat("-", (len(row))+3*(len(row)-1)))
		}
		var rowString []string
		for _, cell := range row {
			switch cell {
			case -1:
				rowString = append(rowString, "X")
			case 1:
				rowString = append(rowString, "O")
			case 0:
				rowString = append(rowString, " ")
			}
		}
		fmt.Println(strings.Join(rowString, " | "))
	}
}

// getCheckInput reads the input from the user from command line, and then calls validInput to check if it is a valid move.
// In the case of a valid move, it returns the x and y coordinates, and if it is not, it reprompts the user for a valid move.
func getCheckInput(board [3][3]int, turn int) (int, int) {
	var x, y int
	reader := bufio.NewReader(os.Stdin)
	var checked = false
	for !checked {
		fmt.Printf("Player %d, Enter the coordinates of your move: ", turn)
		input, _ := reader.ReadString('\n')
		x, y = getCoordinates(input)
		if validInput(board, x, y) {
			checked = true
		} else {
			printBoard(board)
		}
	}
	return x, y
}

// validInput checks if the inputted x and y coordinates are valid for the board, both within the bounds and in an unmarked cell.
func validInput(board [3][3]int, x int, y int) bool {
	if x < 0 || x > 2 || y < 0 || y > 2 {
		return false
	}
	if board[x][y] != 0 {
		fmt.Println("Cell already marked, please enter an empty cell.")
		return false
	}
	return true
}

// updateBoard adds the new mark to the board, and prints the board before returning it
func updateBoard(board [3][3]int, x int, y int, turn int) [3][3]int {
	if turn == 1 {
		board[x][y] = -1
	} else if turn == 2 {
		board[x][y] = 1
	}
	printBoard(board)
	return board
}

// This checks the win conditions of tic tac toe, namely the rows, columns, and diagonals of the board
func checkIfWinner(board [3][3]int, mark int) bool {
	for i := 0; i < 3; i++ {
		if board[i][0] == board[i][1] && board[i][1] == board[i][2] && board[i][0] == mark {
			return true
		}
		if board[0][i] == board[1][i] && board[1][i] == board[2][i] && board[0][i] == mark {
			return true
		}
		if board[0][0] == board[1][1] && board[1][1] == board[2][2] && board[0][0] == mark {
			return true
		}
		if board[2][0] == board[1][1] && board[1][1] == board[0][2] && board[2][0] == mark {
			return true
		}
	}
	return false
}

// This function checks if there are any empty spaces on the board, and if so, returns that it is not a tie. If there are no more
// valid cells left on the board, then the game is declared a tie.
func checkIfTie(board [3][3]int) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

// getCoordinates takes the coordinates and finds out if they are valid coordinates for use, and then uses them, or returns invalid
// coordinates such that validInput will return false
func getCoordinates(input string) (int, int) {
	coordinates := strings.FieldsFunc(input, func(r rune) bool {
		return !unicode.IsNumber(r)
	})
	if len(coordinates) == 1 && len(coordinates[0]) == 2 {
		x, _ := strconv.Atoi(string(coordinates[0][0]))
		y, _ := strconv.Atoi(string(coordinates[0][1]))
		return x, y
	} else if len(coordinates) == 2 {
		x, _ := strconv.Atoi(coordinates[0])
		y, _ := strconv.Atoi(coordinates[1])
		return x, y
	}
	return -1, -1
}

// isEmpty checks what spots in the board are empty, and returns them as 1 values in an all 0 array
func isEmpty(board [3][3]int) []coord {
	empty := make([]coord, 0)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == 0 {
				empty = append(empty, coord{x: i, y: j})
			}
		}
	}
	return empty
}

// updateTurn updates the turn variable to the next turn.
func updateTurn(turn int) int {
	return turn%2 + 1
}

// postGameOutput outputs either a congratulations on cheating to beat the AI, a message there is no shame losing to a perfect AI, or
// that they did a good job when they tied the AI. In the case of 2 players, it congratulates the winning player or both in the case
// of a tie.
func postGameOutput(winner bool, tie bool, withAI bool, turn int, start int) {
	if withAI {
		if winner && (start == 0 && turn == 1 || start == 1 && turn == 2) {
			fmt.Println("Congratulations! You beat an unbeatable AI! Probably by modifying the board before playing ;)")
		} else if winner && start == 0 && turn == 2 || winner && start == 1 && turn == 1 {
			fmt.Println("No shame in losing to a perfect AI! Get em next time!")
		} else if tie {
			fmt.Println("Good job! You tied the AI!")
		}
	} else {
		if winner {
			fmt.Printf("Congratulations Player %d! You won!\n", turn)
		} else {
			fmt.Println("Neither player won! Well played to you both!")
		}
	}
}
