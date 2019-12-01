package snake

import (
	"container/ring"
	"math/rand"

	"flipdisks/pkg/virtualboard"
)

const (
	emptySpace     = 0
	//snakeHeadSpace = 1
	snakeBodySpace = 1
	eggSpace       = 3
)

type direction int

const (
	North direction = iota
	South
	East
	West
)

type snaker interface {
	nextHeadLoc() mapPoint
	eggNextLoc() (x, y int)
}

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

	// allow our internal functions to be mocked out
	snaker
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

	_, _ = snake.Tick(East)

	return snake
}

func (s *Snake) setupGame() {
	s.snaker = s
	s.GameBoard = virtualboard.New(s.boardWidth, s.boardHeight)
	snakeBody := ring.New(s.boardWidth * s.boardHeight)

	gameBoard := *s.GameBoard

	s.deathBoundaries = deathBoundary{}
	s.addOutsideBoundaries()

	bodyX := s.snakeLength + s.startOffset - 1 // subtract 1 because arrays start at 0
	bodyY := s.boardWidth / 2

	s.head = snakeBody

	for bodyRemaining := s.snakeLength; bodyRemaining > 0; bodyRemaining-- {
		point := mapPoint{bodyX, bodyY}

		snakeBody.Value = point
		s.deathBoundaries.Add(bodyX, bodyY)
		gameBoard[int(point.x)][int(point.y)] = snakeBodySpace

		s.tail = snakeBody

		// advance pointer to where we would place the next body segment
		snakeBody = snakeBody.Prev()
		bodyX--
	}

	// add an egg in the same place as where the head is, but on the East side
	// adding 1 because it looks good
	s.eggLoc = mapPoint{s.boardWidth - s.head.Value.(mapPoint).x + 1, bodyY}
	gameBoard[int(s.eggLoc.x)][int(s.eggLoc.y)] = eggSpace

	s.nextTickDirection = East
}

func (s *Snake) addOutsideBoundaries() {
	if s.deathBoundaries == nil {
		s.deathBoundaries = deathBoundary{}
	}

	// draw top/bottom
	for i := -1; i <= s.boardWidth; i++ {
		s.deathBoundaries.Add(i, -1)           // top side
		s.deathBoundaries.Add(i, s.boardWidth) // bottom side
	}

	// draw sides (would start from -1 to height +1, but that's double drawing corners)
	for i := 0; i <= s.boardHeight; i++ {
		s.deathBoundaries.Add(-1, i)           // left side
		s.deathBoundaries.Add(s.boardWidth, i) // right side
	}
}

func (b *deathBoundary) Add(x, y int) {
	boundary := *b
	_, ok := boundary[xPos(x)]
	if !ok {
		boundary[xPos(x)] = map[yPos]wallExists{}
	}

	boundary[xPos(x)][yPos(y)] = wallExists{}
}

func (b *deathBoundary) Remove(x, y int) {
	boundary := *b
	_, found := boundary[xPos(x)]
	if !found {
		return // easy peasy
	}

	delete(boundary[xPos(x)], yPos(y))
}

func (b *deathBoundary) IsBoundary(x, y int) bool {
	_, exists := (*b)[xPos(x)]
	if !exists {
		return false
	}

	_, dead := (*b)[xPos(x)][yPos(y)]
	return dead
}

func (s *Snake) Tick(nextDirection direction) (isGameOver, gameWin bool) {
	s.nextTickDirection = nextDirection

	isGameOver, gameWin = s.checkGameStatus()
	if isGameOver {
		return isGameOver, gameWin
	}

	canGetEgg := s.willGetEgg()

	if !canGetEgg {
		s.moveSnake(false)
	} else {
		s.moveSnake(true)
		ableToAddEgg := s.addEgg()
		if !ableToAddEgg {
			return true, true
		}
	}

	return false, false
}

func (s *Snake) moveSnake(getLonger bool) {
	gameBoard := *s.GameBoard

	if getLonger {
		s.snakeLength++
	} else {
		oldTail := s.tail.Value.(mapPoint)
		gameBoard[oldTail.x][oldTail.y] = emptySpace
		s.deathBoundaries.Remove(oldTail.x, oldTail.y)

		s.tail = s.tail.Next()
	}

	nextHead := s.nextHeadLoc()

	s.head = s.head.Next()
	s.head.Value = nextHead
	s.deathBoundaries.Add(nextHead.x, nextHead.y)
	gameBoard[nextHead.x][nextHead.y] = snakeBodySpace
}

func (s *Snake) nextHeadLoc() mapPoint {
	currentHead := s.head.Value.(mapPoint)

	var nextHead mapPoint
	switch s.nextTickDirection {
	case North:
		nextHead = mapPoint{currentHead.x, currentHead.y - 1}
	case South:
		nextHead = mapPoint{currentHead.x, currentHead.y + 1}
	case East:
		nextHead = mapPoint{currentHead.x + 1, currentHead.y}
	case West:
		nextHead = mapPoint{currentHead.x - 1, currentHead.y}
	}
	return nextHead
}

func (s *Snake) checkGameStatus() (isGameOver, gameWin bool) {
	nextHead := s.snaker.nextHeadLoc()
	dead := s.deathBoundaries.IsBoundary(nextHead.x, nextHead.y)
	if dead {
		return true, false
	}

	if s.snakeLength == s.boardWidth*s.boardHeight {
		return true, true
	}

	return false, false
}

func (s *Snake) willGetEgg() bool {
	nextHead := s.nextHeadLoc()

	if nextHead.x == s.eggLoc.x && nextHead.y == s.eggLoc.y {
		return true
	}
	return false
}

func (s *Snake) eggNextLoc() (x, y int) {
	x = rand.Intn(s.boardWidth)
	y = rand.Intn(s.boardHeight)
	return x, y
}

func (s *Snake) addEgg() bool {
	itsFull, added := make(chan struct{}), make(chan struct{})

	// let's make sure the board isn't full
	go func() {
		howFull := 0
		for _, deathY := range s.deathBoundaries {
			howFull = howFull + len(deathY)
		}

		// +2 because deathBoundaries are created outside the board (top/bottom left/right)
		if howFull >= (s.boardHeight+2)*(s.boardWidth+2) {
			itsFull <- struct{}{}
			return
		}

		// just for sanity sake, if that condition didn't work, lets try this one
		if s.snakeLength >= s.boardHeight*s.boardWidth {
			itsFull <- struct{}{}
			return
		}
	}()

	// try placing an egg randomly, if we can't then lets just start iterating from that location
	go func() {
		x, y := s.snaker.eggNextLoc()

		for xTries := 0; xTries < s.boardWidth; xTries++ {
			for yTries := 0; yTries < s.boardHeight; yTries++ {
				boundaryExists := s.deathBoundaries.IsBoundary(x, y)
				if !boundaryExists {
					s.eggLoc = mapPoint{x: x, y: y}
					(*s.GameBoard)[x][y] = eggSpace
					added <- struct{}{}
					return // exit early
				}

				y = (y + 1) % s.boardHeight
			}
			x = (x + 1) % s.boardWidth
		}
		// we tried everywhere, it's full... or there's a bug...
		itsFull <- struct{}{}
	}()

	// it's a race! which goroutine will finish first?
	// yeah, goroutines aren't parallel, but for some reason it's faster 🤷
	// something must not be 100% synchronous, so we're kind of relying on that hack
	select {
	case <-itsFull:
		return false
	case <-added:
		return true
	}
}
