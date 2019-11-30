package snake

import (
	"container/ring"

	"flipdisks/pkg/virtualboard"
)

const (
	emptySpace     = 0
	snakeHeadSpace = 1
	snakeBodySpace = 2
	eggSpace       = 3
)

type direction int

const (
	north direction = iota
	south
	east
	west
)

type Snake struct {
	boardHeight int
	boardWidth  int

	startOffset int
	snakeLength int

	// [ nil,... tail, (body,body,body,...) head, nil,nil, ...]   <--- circular buffer
	//            ^tailPointer               ^headPointer
	head *ring.Ring
	tail *ring.Ring

	nextTickDirection direction

	eggLoc mapPoint

	// if (x,y) exists, there's a boundary here, (snake or wall)
	deathBoundaries deathBoundary

	// exposed for you to view the board anytime
	GameBoard *virtualboard.VirtualBoard
}

type xPos int
type yPos int
type deathBoundary map[xPos]map[yPos]wallExists
type wallExists struct{}

type mapPoint struct {
	x, y int
}

func New() *Snake {
	snake := &Snake{
		boardHeight: 11,
		boardWidth:  11,
		startOffset: 2,
		snakeLength: 3,
	}

	snake.setupGame()

	_, _ = snake.Tick()

	return snake
}

func (s *Snake) setupGame() {
	s.GameBoard = virtualboard.New(s.boardWidth, s.boardHeight)
	snakeBody := ring.New(s.boardWidth * s.boardHeight)

	gameBoard := *s.GameBoard

	s.addOutsideBoundaries()

	bodyX := s.snakeLength + s.startOffset - 1 // subtract 1 because arrays start at 0
	bodyY := s.boardWidth / 2

	s.head = snakeBody

	for bodyRemaining := s.snakeLength; bodyRemaining > 0; bodyRemaining-- {
		point := mapPoint{bodyX, bodyY}

		snakeBody.Value = point
		s.addBoundary(bodyX, bodyY)
		gameBoard[int(point.y)][int(point.x)] = 1

		s.tail = snakeBody

		// advance pointer to where we would place the next body segment
		snakeBody = snakeBody.Prev()
		bodyX--
	}

	// add an egg in the same place as where the head is, but on the east side
	// adding 1 because it looks good
	s.eggLoc = mapPoint{s.boardWidth - s.head.Value.(mapPoint).x + 1, bodyY}
	gameBoard[int(s.eggLoc.y)][int(s.eggLoc.x)] = 1

	s.nextTickDirection = east
}

func (s *Snake) addOutsideBoundaries() {
	// draw top
	for i := -1; i < s.boardWidth+1; i++ {
		s.addBoundary(-1, i)
	}

	// draw bottom
	for i := -1; i < s.boardWidth+1; i++ {
		s.addBoundary(s.boardHeight, i)
	}

	// draw left (would start from -1 to height +1, but that's double drawing corners)
	for i := 0; i < s.boardHeight; i++ {
		s.addBoundary(i, -1)
	}

	// draw right (would start from -1 to height +1, but that's double drawing corners)
	for i := 0; i < s.boardHeight; i++ {
		s.addBoundary(i, s.boardHeight)
	}
}

func (s *Snake) addBoundary(x, y int) {
	if s.deathBoundaries == nil {
		s.deathBoundaries = deathBoundary{}
	}

	_, ok := s.deathBoundaries[xPos(x)]
	if !ok {
		s.deathBoundaries[xPos(x)] = map[yPos]wallExists{}
	}

	s.deathBoundaries[xPos(x)][yPos(y)] = wallExists{}
}

//func (s *Snake) startGame() {
//	boardHeight := len(s.board)
//	boardLength := len(s.board[0])
//
//	headX := s.snakeLength + s.tailOffset - 1
//	headY := boardHeight / 2
//	s.board[headY][headX] = snakeHeadSpace
//
//	bodyX := headX
//	for bodyRemaining := s.snakeLength - 1; bodyRemaining > 0; bodyRemaining-- {
//		bodyX--
//		s.board[headY][bodyX] = snakeBodySpace
//	}
//
//	s.board[headY][boardLength-headX+1] = eggSpace
//}
//
//func (s *Snake) addEgg() bool {
//	boardLength := len(s.board)
//	boardHeight := len(s.board[0])
//
//	// let's just try to place it somewhere near an empty area
//	// we'll try x first, if that doesn't work, move down a row and try x again
//	// we've exhausted all possible positions, game over
//
//	eggX := rand.Intn(boardLength)
//	eggY := rand.Intn(boardHeight)
//
//	var xTries, yTries int
//	for {
//		if s.board[eggX][eggY] == emptySpace {
//			s.board[eggX][eggY] = eggSpace
//			return false
//		} else {
//			eggX = (eggX + 1) % boardLength
//			xTries++
//
//			if xTries >= boardLength {
//				eggY = (eggY + 1) % boardLength
//				yTries++
//				xTries = rand.Intn(boardLength)
//
//				if yTries >= boardHeight {
//					return true
//				}
//			}
//		}
//	}
//}

func (s *Snake) Tick() (isGameOver, gameWin bool) {
	s.moveSnake(false)

	isGameOver, gameWin = s.checkBoundaries()
	if isGameOver {
		return isGameOver, gameWin
	}

	gotEgg := s.gotEgg()

	if gotEgg == true {
		s.moveSnake(true)
		ableToAddEgg := s.addEgg()
		if !ableToAddEgg {
			return true, true
		}
	} else {
		s.moveSnake(false)
	}

	return false, false
}

func (s *Snake) moveSnake(gotLonger bool) {
	if !gotLonger {
		s.tail = s.tail.Next()
	}

	switch s.nextTickDirection {
	case north:
	case south:
	case east:
	case west:
	}
}

func (s *Snake) checkBoundaries() (isGameOver, gameWin bool) {
	return false, false
}

func (s *Snake) gotEgg() bool {
	return false
}

func (s *Snake) addEgg() bool {
	// remove old position of egg
	// add in new egg randomly
	return true
}
