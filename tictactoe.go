package main

import (
	"fmt"
	"strconv"
)

type Endgame struct {
	X, O, Tied int
}

func (end Endgame) String() string {
	var total = float64(end.X + end.O + end.Tied)
	return fmt.Sprintf(
		"%.1f%% of X winning, %.1f%% of O winning, %.1f%% of a tie\n",
		float64(end.X)/total*100,
		float64(end.O)/total*100,
		float64(end.Tied)/total*100,
	)
}

// Create a tree structure of all possibilities
type Play struct {
	board   Board
	choices []Play
	end     Endgame
}

func (play *Play) descend(player Player) {
	// Determine the available moves
	options := play.board.available()

	// Will the player win if the make any of these moves?
	// TODO multiple winning moves?
	for _, option := range options {
		// Create a copy of the board
		board := play.board.Copy()

		board[option] = player
		if board.wonBy(player) {
			var end Endgame
			if player == X {
				end.X = 1
			} else {
				end.O = 1
			}

			// Choose this as the sole play
			play.choices = []Play{{board: board, end: end}}
			return
		}
	}

	// Otherwise play every possible move
	play.choices = make([]Play, len(options))

	for i, option := range options {
		board := play.board.Copy()
		board[option] = player
		play.choices[i].board = board

		// If these was the last possible move, record it as a tie
		if len(options) == 1 {
			play.choices[i].end.Tied = 1
		}
	}
	return
}

// Save the plays by board
var choices = make(map[Board]Play)

func descend(play Play, turn int) Play {
	var player Player
	if turn%2 == 0 {
		player = X
	} else {
		player = O
	}

	// TODO call this play?
	play.descend(player)

	// TODO check if the board already exists!

	for i, choice := range play.choices {
		// Only descend if the game isn't over
		if choice.end.X == 0 && choice.end.O == 0 && choice.end.Tied == 0 {
			// TODO Or just mutate
			play.choices[i] = descend(choice, turn+1)
		}

		// Aggregate the endgames
		play.end.X += play.choices[i].end.X
		play.end.O += play.choices[i].end.O
		play.end.Tied += play.choices[i].end.Tied
	}

	choices[play.board] = play
	return play
}

// Recursively construct all game boards
func buildEndgames() {
	// Create the initial move then descend
	var play Play
	play = descend(play, 0)

	// Count all the endgames
	// fmt.Println("X Won:", play.end.X)
	// fmt.Println("Y Won:", play.end.Y)
	// fmt.Println("Tied:", play.end.Tied)

	// fmt.Println("Boards:", len(choices))

	return
}

type Board [9]*int

func (b Board) Copy() Board {
	return b
}

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
	buildEndgames()

	var board Board
	var turn int
	for {
		// Are there available squares?
		if len(board.available()) == 0 {
			// The game should be over - if there is no winner it is a draw
			fmt.Println("The game was a draw")
			return
		}

		endgame := choices[board]

		var current Player
		if turn%2 == 0 {
			current = X
			board.turn(current)
		} else {
			current = O

			// Have the computer pick a play
			// Pick the play with the highest chance of O winning
			var highest float64
			var choice int

			for i, chances := range endgame.choices {
				total := chances.end.O + chances.end.X + chances.end.Tied
				// Go for highest win
				win := float64(chances.end.O) / float64(total)
				if win > highest {
					highest = win
					choice = i
				}

			}
			board = endgame.choices[choice].board
		}

		fmt.Println(board)

		// Show the win probability
		fmt.Println(endgame.end)

		if board.wonBy(current) {
			fmt.Printf("Player %c has won\n", *current)
			return
		}
		turn += 1
	}
}
