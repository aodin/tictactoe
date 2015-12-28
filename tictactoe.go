package main

import (
	"fmt"
	"strconv"
)

type Board [9]*int

func (b Board) String() (out string) {
	for i, square := range b {
		if square == nil {
			out += "."
		} else {
			out += fmt.Sprintf("%c", *square)
		}
		if i%3 == 2 {
			out += "\n"
		}
	}
	return
}

func (b Board) Options() (out string) {
	for i, square := range b {
		if square == nil {
			out += strconv.Itoa(i)
		} else {
			out += "-"
		}
		if i%3 == 2 {
			out += "\n"
		}
	}
	return
}

// TODO automatically construct the wins
var wins = [...][3]int{
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},
	{0, 4, 8},
	{2, 4, 6},
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},
}

func (b Board) available() (squares []int) {
	for i, square := range b {
		if square == nil {
			squares = append(squares, i)
		}
	}
	return
}

func (b Board) wonBy(current Player) bool {
	// Check for win
	for _, win := range wins {
		over := true
		for _, square := range win {
			if b[square] == nil {
				over = false
				break
			}
			over = over && (b[square] == current)
		}
		if over {
			return true
		}
	}
	return false
}

type Player *int

var xxx = 88 // X
var ooo = 79 // O

var X = Player(&xxx)
var O = Player(&ooo)

func (board *Board) turn(current Player) {
	fmt.Printf("Player %c - pick a square!\n", *current)
	fmt.Println(board.Options())

	var picked int
	for {
		var raw string
		if _, err := fmt.Scanln(&raw); err != nil {
			fmt.Println(err)
			return
		}

		square, err := strconv.Atoi(raw)
		if err != nil {
			// Pick again until the square is valid
			fmt.Printf("Could not understand %s - pick again\n", raw)
			continue
		}

		if square < 0 || square > 8 {
			fmt.Println("Square must be between 0 and 8 - pick again")
		} else if board[square] != nil {
			fmt.Println("Square is taken - pick again")
		} else {
			picked = square
			break
		}
	}
	board[picked] = current
}

func main() {
	var board Board

	var turn int
	for {
		// Are there available squares?
		if len(board.available()) == 0 {
			// The game should be over - if there is no winner it is a draw
			fmt.Println("The game was a draw")
			return
		}

		var current Player
		if turn%2 == 0 {
			current = X
		} else {
			current = O
		}

		board.turn(current)
		fmt.Println(board)

		if board.wonBy(current) {
			fmt.Printf("Player %c has won\n", *current)
			return
		}
		turn += 1
	}
}
