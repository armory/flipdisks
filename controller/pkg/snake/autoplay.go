package snake

import (
	"container/heap"
	"fmt"
	"math/rand"
	"time"
)

func (s *Snake) AutoPlayRandomly() (isGameOver bool) {
	directions := []direction{East, North, South, West}

	fmt.Println(s.GameBoard) // print starting game
	for !isGameOver {
		i := rand.Intn(len(directions)) % len(directions)
		isGameOver, _ := s.Tick(directions[i])

		// go into survival mode
		if isGameOver {
			for ohGod := 0; ohGod <= len(directions)-1; ohGod++ {
				isGameOver, _ = s.Tick(directions[i])
			}
		}

		fmt.Println(s.GameBoard)
		time.Sleep(50 * time.Millisecond)
	}

	return true
}

type autoPlay struct {
	q priorityQueue
}

func (s *Snake) AutoPlay() (gameOver bool) {
	for !gameOver {
		ap := autoPlay{q: priorityQueue{}}

		sCopy := s.copy()
		directions, found := ap.getPath(sCopy)

		fmt.Println(s.GameBoard) // print the start of the game

		if found {
			// start walking
			for _, d := range directions {
				gameOver, _ = s.Tick(d)
				fmt.Println(s.GameBoard)
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	return true
}

func (a autoPlay) getPath(snake *Snake) (directions []direction, found bool) {
	node := &autoPlayNode{
		s:                 snake,
		heuristicDistance: MaxInt,
	}
	heap.Push(&a.q, node)

	for {
		// we couldn't find a path
		if a.q.Len() == 0 {
			return nil, false
		}

		curr := heap.Pop(&a.q).(*autoPlayNode)

		// yay! found the egg, return the directions we should follow
		if curr.gotEgg {
			for curr.parent != nil {
				directions = append([]direction{curr.s.nextTickDirection}, directions...)
				curr = curr.parent
			}
			return directions, true
		}

		// start exploring our neighbors
		for _, neighbor := range curr.explore() {
			neighborHeuristic := neighbor.travelCost() + neighbor.getHeuristicDistance()
			if neighbor.gotEgg {
				neighborHeuristic = 0 // omg yay!
			}

			// the neighbor is better than our current node, we should really explore this one
			if neighborHeuristic <= curr.heuristicDistance {
				neighbor.heuristicDistance = neighborHeuristic
				heap.Push(&a.q, neighbor)
			}
		}
	}
}

type autoPlayNode struct {
	s                 *Snake
	parent            *autoPlayNode
	gotEgg            bool
	heuristicDistance int // sort by heuristicDistance

	index int // used internally by sort package
}

func (a *autoPlayNode) explore() []*autoPlayNode {
	var snakeFeelers []*autoPlayNode

	for _, direction := range []direction{East, North, South, West} {
		snake := a.s.copy()
		snake.nextTickDirection = direction
		gotEgg := snake.willGetEgg()
		isGameOver, win := snake.Tick(direction)

		// game is over, we don't need to try this spot
		if isGameOver && !win {
			continue
		}

		nextMove := &autoPlayNode{
			s:      snake,
			parent: a,
			gotEgg: gotEgg,
		}

		snakeFeelers = append(snakeFeelers, nextMove)
	}
	return snakeFeelers
}

// sometime in the future, we might want to change
func (a *autoPlayNode) travelCost() int {
	return 1
}

func (a *autoPlayNode) getHeuristicDistance() int {
	currentPos := a.s.head.Value.(mapPoint)

	// use the manhattan distance because it's fast and good enough
	xDist := intAbs(currentPos.x - a.s.eggLoc.x)
	yDist := intAbs(currentPos.y - a.s.eggLoc.y)

	return xDist + yDist
}

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

func intAbs(n int) int {
	y := n >> 31
	return (n ^ y) - y
}
