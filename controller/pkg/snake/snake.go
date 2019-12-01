package snake

import (
	"container/ring"
	"fmt"
	"math"
	"math/rand"

	"flipdisks/pkg/virtualboard"
	"github.com/beefsack/go-astar"
)

const (
	emptySpace     = 0
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
	nextHeadLoc(d direction) mapPoint
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
	disableGameBoard bool
	GameBoard        *virtualboard.VirtualBoard

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

type Option func(s *Snake) error

func New(boardHeight, boardWidth, startOffset, snakeLength int, opts ...Option) (*Snake, error) {
	s := &Snake{
		boardHeight: boardHeight,
		boardWidth:  boardWidth,
		startOffset: startOffset,
		snakeLength: snakeLength,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return &Snake{}, err
		}
	}

	s.setupGame()

	_, _ = s.Tick(East)

	return s, nil
}

func DisableGameBoard() Option {
	return func(s *Snake) error {
		s.disableGameBoard = true
		return nil
	}
}

func (s *Snake) setupGame() {
	s.snaker = s
	snakeBody := ring.New(s.boardWidth * s.boardHeight)

	if !s.disableGameBoard {
		s.GameBoard = virtualboard.New(s.boardWidth, s.boardHeight)
	}
	s.deathBoundaries = deathBoundary{}
	s.addOutsideBoundaries()

	bodyX := s.snakeLength + s.startOffset - 1 // subtract 1 because arrays start at 0
	bodyY := s.boardHeight / 2

	s.head = snakeBody

	for bodyRemaining := s.snakeLength; bodyRemaining > 0; bodyRemaining-- {
		point := mapPoint{bodyX, bodyY}

		snakeBody.Value = point
		s.deathBoundaries.Add(bodyX, bodyY)
		(*s.GameBoard)[point.x][point.y] = snakeBodySpace

		s.tail = snakeBody

		// advance pointer to where we would place the next body segment
		snakeBody = snakeBody.Prev()
		bodyX--
	}

	// add an egg in the same place as where the head is, but on the East side
	// adding 1 because it looks good
	s.eggLoc = mapPoint{s.boardWidth - s.head.Value.(mapPoint).x + 1, bodyY}
	(*s.GameBoard)[s.eggLoc.x][s.eggLoc.y] = eggSpace

	s.nextTickDirection = East
}

func (s *Snake) DisableGameBoard() {
	s.disableGameBoard = true
	s.GameBoard = nil
}

func (s *Snake) EnableGameBoard() {
	s.disableGameBoard = false
	s.GameBoard = virtualboard.New(s.boardWidth, s.boardHeight)

	(*s.GameBoard)[s.eggLoc.x][s.eggLoc.y] = eggSpace

	snakeWalk := s.head
	for snakeWalk.Value != s.tail.Value {
		bodyLoc := snakeWalk.Value.(mapPoint)
		(*s.GameBoard)[bodyLoc.x][bodyLoc.y] = snakeBodySpace
		snakeWalk.Prev()
	}
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
	if getLonger {
		s.snakeLength++
	} else {
		oldTail := s.tail.Value.(mapPoint)
		s.deathBoundaries.Remove(oldTail.x, oldTail.y)
		if !s.disableGameBoard {
			(*s.GameBoard)[oldTail.x][oldTail.y] = emptySpace
		}

		s.tail = s.tail.Next()
	}

	nextHead := s.nextHeadLoc(s.nextTickDirection)

	s.head = s.head.Next()
	s.head.Value = nextHead
	s.deathBoundaries.Add(nextHead.x, nextHead.y)
	if !s.disableGameBoard {
		(*s.GameBoard)[nextHead.x][nextHead.y] = snakeBodySpace
	}
}

func (s *Snake) nextHeadLoc(d direction) mapPoint {
	currentHead := s.head.Value.(mapPoint)

	var nextHead mapPoint
	switch d {
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
	nextHead := s.snaker.nextHeadLoc(s.nextTickDirection)
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
	nextHead := s.nextHeadLoc(s.nextTickDirection)

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
					if !s.disableGameBoard {
						(*s.GameBoard)[x][y] = eggSpace
					}
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
