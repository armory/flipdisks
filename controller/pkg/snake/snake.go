package snake

import (
	"math/rand"

	"flipdisks/pkg/virtualboard"
)

const (
	emptySpace     = 0
	snakeHeadSpace = 1
	snakeBodySpace = 2
	eggSpace       = 3
)

const snakeLength = 3

type Snake struct {

}

func New() {
	board := virtualboard.VirtualBoard{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}

	addSnake(&board)
	addEgg(&board)
}

func addSnake(boardPointer *virtualboard.VirtualBoard) {
	board := *boardPointer
	boardHeight := len(board[0])

	tailOffset := 2
	headX := snakeLength + tailOffset - 1
	headY := boardHeight / 2
	board[headY][headX] = snakeHeadSpace

	bodyX := headX
	for bodyRemaining := snakeLength - 1; bodyRemaining > 0; bodyRemaining-- {
		bodyX--
		board[headY][bodyX] = snakeBodySpace
	}
}


func addEgg(boardPointer *virtualboard.VirtualBoard) bool {
	board := *boardPointer
	boardLength := len(board)
	boardHeight := len(board[0])

	// let's just try to place it somewhere near an empty area
	// we'll try x first, if that doesn't work, move down a row and try x again
	// we've exhausted all possible positions, game over

	eggX := rand.Intn(boardLength)
	eggY := rand.Intn(boardHeight)

	var xTries, yTries int
	for {
		if board[eggX][eggY] == emptySpace {
			board[eggX][eggY] = 1
			return false
		} else {
			eggX = (eggX + 1) % boardLength
			xTries++

			if xTries >= boardLength {
				eggY = (eggY + 1) % boardLength
				yTries++
				xTries = rand.Intn(boardLength)

				if yTries >= boardHeight {
					return true
				}
			}
		}
	}
}
