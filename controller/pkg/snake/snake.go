package snake

import (
	"container/ring"
	"math/rand"

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
	North direction = iota
	South
	East
	West
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

	_, _ = snake.Tick(East)

	return snake
}

func (s *Snake) setupGame() {
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
		gameBoard[int(point.x)][int(point.y)] = 1

		s.tail = snakeBody

		// advance pointer to where we would place the next body segment
		snakeBody = snakeBody.Prev()
		bodyX--
	}

	// add an egg in the same place as where the head is, but on the East side
	// adding 1 because it looks good
	s.eggLoc = mapPoint{s.boardWidth - s.head.Value.(mapPoint).x + 1, bodyY}
	gameBoard[int(s.eggLoc.x)][int(s.eggLoc.y)] = 1

	s.nextTickDirection = East
}

func (s *Snake) addOutsideBoundaries() {
	if s.deathBoundaries == nil {
		s.deathBoundaries = deathBoundary{}
	}

	// draw top
	for i := -1; i < s.boardWidth+1; i++ {
		s.deathBoundaries.Add(-1, i)
	}

	// draw bottom
	for i := -1; i < s.boardWidth+1; i++ {
		s.deathBoundaries.Add(s.boardHeight, i)
	}

	// draw left (would start from -1 to height +1, but that's double drawing corners)
	for i := 0; i < s.boardHeight; i++ {
		s.deathBoundaries.Add(i, -1)
	}

	// draw right (would start from -1 to height +1, but that's double drawing corners)
	for i := 0; i < s.boardHeight; i++ {
		s.deathBoundaries.Add(i, s.boardHeight)
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

func (s *Snake) moveSnake(gotLonger bool) {
	gameBoard := *s.GameBoard

	if !gotLonger {
		oldTail := s.tail.Value.(mapPoint)
		gameBoard[oldTail.x][oldTail.y] = 0
		s.deathBoundaries.Remove(oldTail.x, oldTail.y)

		s.tail = s.tail.Next()
	} else {
		// do nothing to the tail, lazy eval
	}

	oldHead := s.head.Value.(mapPoint)

	var newHead mapPoint
	switch s.nextTickDirection {
	case North:
		newHead = mapPoint{oldHead.x, oldHead.y - 1}
	case South:
		newHead = mapPoint{oldHead.x, oldHead.y + 1}
	case East:
		newHead = mapPoint{oldHead.x + 1, oldHead.y}
	case West:
		newHead = mapPoint{oldHead.x - 1, oldHead.y}
	}

	s.head = s.head.Next()
	s.head.Value = newHead
	s.deathBoundaries.Add(newHead.x, newHead.y)
	gameBoard[newHead.x][newHead.y] = 1
}

func (s *Snake) checkGameStatus() (isGameOver, gameWin bool) {
	_, dead := s.deathBoundaries[xPos(s.head.Value.(mapPoint).x)][yPos(s.head.Value.(mapPoint).y)]
	if dead {
		return true, false
	}

	if s.snakeLength == s.boardWidth*s.boardWidth {
		return true, true
	}

	return false, false
}

func (s *Snake) willGetEgg() bool {
	if s.head.Value.(mapPoint).x == s.eggLoc.x && s.head.Value.(mapPoint).y == s.eggLoc.y {
		return true
	}
	return false
}

func (s *Snake) addEgg() bool {
	full := make(chan struct{})
	added := make(chan struct{})

	// let's make sure the board isn't full
	go func() {
		howFull := 0
		for _, deathY := range s.deathBoundaries {
			howFull = howFull + len(deathY)
		}

		// +2 because deathBoundaries are created outside the board (top/bottom left/right)
		if howFull == (s.boardHeight+2)*(s.boardWidth+2) {
			full <- struct{}{}
		}
	}()

	// try placing an egg randomly, if we can't then lets just start iterating from that location
	go func() {
		eggX := rand.Intn(s.boardWidth)
		eggY := rand.Intn(s.boardHeight)

		for xTries := 0; xTries < s.boardWidth; xTries++ {
			for yTries := 0; yTries < s.boardHeight; yTries++ {
				_, hitBoundary := s.deathBoundaries[xPos(eggX)][yPos(eggY)]
				if !hitBoundary {
					s.eggLoc = mapPoint{x: eggX, y: eggY}
					added <- struct{}{}
				}

				eggY = (eggY + 1) % s.boardHeight
			}
			eggX = (eggX + 1) % s.boardWidth
		}
	}()

	// it's a race! which goroutine will finish first?
	// yeah, goroutines aren't parallel, but for some reason it's faster 🤷
	// something must not be 100% synchronous, so we're kind of relying on that hack
	select {
	case <-full:
		return false
	case <-added:
		return true
	}
}
