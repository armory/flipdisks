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

type Snake struct {
	board       virtualboard.VirtualBoard
	tailOffset  int
	snakeLength int
}

func New() {
	board := virtualboard.VirtualBoard{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}

	snake := Snake{
		board:       board,
		tailOffset:  2,
		snakeLength: 3,
	}

	snake.addSnake()
	snake.addEgg()
}

func (s *Snake) addSnake() {
	boardHeight := len(s.board[0])

	headX := s.snakeLength + s.tailOffset - 1
	headY := boardHeight / 2
	s.board[headY][headX] = snakeHeadSpace

	bodyX := headX
	for bodyRemaining := s.snakeLength - 1; bodyRemaining > 0; bodyRemaining-- {
		bodyX--
		s.board[headY][bodyX] = snakeBodySpace
	}
}

func (s *Snake) addEgg() bool {
	boardLength := len(s.board)
	boardHeight := len(s.board[0])

	// let's just try to place it somewhere near an empty area
	// we'll try x first, if that doesn't work, move down a row and try x again
	// we've exhausted all possible positions, game over

	eggX := rand.Intn(boardLength)
	eggY := rand.Intn(boardHeight)

	var xTries, yTries int
	for {
		if s.board[eggX][eggY] == emptySpace {
			s.board[eggX][eggY] = 1
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
