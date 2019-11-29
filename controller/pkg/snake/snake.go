package snake

import (
	"math/rand"

	"flipdisks/pkg/virtualboard"
)

//func New() {
//	 virtualboard.VirtualBoard{
//		fontmap.Row{0, 0, 0, 0},
//		fontmap.Row{0, 0, 0, 0},
//		fontmap.Row{0, 0, 0, 0},
//		fontmap.Row{0, 0, 0, 0},
//	}
//}

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
		if board[eggX][eggY] == 0 {
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
