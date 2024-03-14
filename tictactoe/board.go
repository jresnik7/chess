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

type coord struct {
	x int
	y int
}

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
			if start == 0 {
				// Only print the board again if it is the first move, as the board hasn't been displayed yet
				if moves == 0 {
					printBoard(board)
				}
				x, y = GetCheckInput(board, turn)
				board = UpdateBoard(board, x, y, turn)
				moves += 1
				winner = CheckIfWinner(board, 1) || CheckIfWinner(board, -1)
				tie = CheckIfTie(board)
				if winner || tie {
					break
				}
				turn = UpdateTurn(turn)
				optimalMove := minimax(board, aiMark, humanMark, aiMark, moves)
				fmt.Println("The AI has made the following move: ")
				board = UpdateBoard(board, optimalMove.coordinates.x, optimalMove.coordinates.y, turn)
				moves += 1
				winner = CheckIfWinner(board, 1) || CheckIfWinner(board, -1)
				tie = CheckIfTie(board)
				if winner || tie {
					break
				}
				turn = UpdateTurn(turn)
			} else {
				optimalMove := minimax(board, aiMark, humanMark, aiMark, moves)
				fmt.Println("The AI has made the following move: ")
				board = UpdateBoard(board, optimalMove.coordinates.x, optimalMove.coordinates.y, turn)
				moves += 1
				winner = CheckIfWinner(board, 1) || CheckIfWinner(board, -1)
				tie = CheckIfTie(board)
				if winner || tie {
					break
				}
				turn = UpdateTurn(turn)
				x, y = GetCheckInput(board, turn)
				board = UpdateBoard(board, x, y, turn)
				moves += 1
				winner = CheckIfWinner(board, 1) || CheckIfWinner(board, -1)
				tie = CheckIfTie(board)
				if winner || tie {
					break
				}
				turn = UpdateTurn(turn)
			}
		} else {
			printBoard(board)
			x, y = GetCheckInput(board, turn)
			board = UpdateBoard(board, x, y, turn)
			winner = CheckIfWinner(board, 1) || CheckIfWinner(board, -1)
			tie = CheckIfTie(board)
			if winner || tie {
				break
			}
			turn = UpdateTurn(turn)
		}
	}
	PostGameOutput(winner, tie, withAI, turn, start)
}

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

func GetCheckInput(board [3][3]int, turn int) (int, int) {
	var x, y int
	reader := bufio.NewReader(os.Stdin)
	var checked = false
	for !checked {
		fmt.Printf("Player %d, Enter the coordinates of your move: ", turn)
		input, _ := reader.ReadString('\n')
		x, y = GetCoordinates(input)
		if ValidInput(board, x, y) {
			checked = true
		} else {
			printBoard(board)
		}
	}
	return x, y
}

func ValidInput(board [3][3]int, x int, y int) bool {
	if x < 0 || x > 2 || y < 0 || y > 2 {
		return false
	}
	if board[x][y] != 0 {
		fmt.Println("Cell already marked, please enter an empty cell.")
		return false
	}
	return true
}

func UpdateBoard(board [3][3]int, x int, y int, turn int) [3][3]int {
	if turn == 1 {
		board[x][y] = -1
	} else if turn == 2 {
		board[x][y] = 1
	}
	printBoard(board)
	return board
}

func CheckIfWinner(board [3][3]int, mark int) bool {
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
func CheckIfTie(board [3][3]int) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

// GetCoordinates takes the coordinates and finds out if they are valid coordinates for use, and then uses them, or returns invalid
// coordinates such that ValidInput will return false
func GetCoordinates(input string) (int, int) {
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

// IsEmpty checks what spots in the board are empty, and returns them as 1 values in an all 0 array
func IsEmpty(board [3][3]int) []coord {
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
