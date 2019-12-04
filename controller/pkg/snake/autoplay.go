package snake

import (
	"container/heap"
	"fmt"
	"math/rand"
	"time"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

func (s *Snake) AutoPlayRandomly() {
	directions := []direction{East, North, South, West}
	var isGameOver, _ bool

	fmt.Println(s.GameBoard)
	for {
		i := rand.Intn(len(directions)) % len(directions)
		isGameOver, _ = s.Tick(directions[i])

		// go into survival mode
		if isGameOver {
			deaths := 0
			for ohGod := 0; ohGod <= len(directions)-1; ohGod++ {
				isGameOver, _ = s.Tick(directions[i])
				if isGameOver {
					deaths++
				}
			}

			// really died
			if deaths == len(directions) {
				return
			}
		}

		fmt.Println(s.GameBoard)
		time.Sleep(50 * time.Millisecond)
	}
}

func (s *Snake) AutoPlay() (gameOver bool) {
	for !gameOver {
		ss := autoPlay{
			q: priorityQueue{},
		}

		sCopy := s.copy()
		directions, found := ss.getPath(sCopy) // the 2nd param is empty because egg is already included everywhere

		fmt.Println(s.GameBoard)

		if found {
			for _, d := range directions {
				s.Tick(d)
				fmt.Println(s.GameBoard)
				time.Sleep(10 * time.Millisecond)
			}
		} else {
			break
		}
	}
	return true
}

type autoPlay struct {
	q priorityQueue
}

func (a autoPlay) getPath(snake *Snake) (directions []direction, found bool) {
	node := &autoPlayNode{
		s:                 snake,
		heuristicDistance: MaxInt,
	}
	heap.Push(&a.q, node)

	for {
		fmt.Println("queue length: ", a.q.Len())
		if a.q.Len() == 0 {
			return nil, false
		}

		curr := heap.Pop(&a.q).(*autoPlayNode)
		//fmt.Println("pop", curr.s.nextTickDirection, curr.heuristicDistance, fmt.Sprintf("%p", curr))

		// found the egg, return the directions!
		if curr.gotEgg {
			for curr.parent != nil {
				directions = append([]direction{curr.s.nextTickDirection}, directions...)
				curr = curr.parent
			}
			return directions, true
		}

		for _, neighbor := range curr.explore() {
			//logrus.Infof("%v->%v %s - egg%v", neighbor.s.head.Value.(mapPoint), neighbor.s.head.Value.(mapPoint), neighbor.s.nextTickDirection, neighbor.s.eggLoc)

			neighborHeuristic := neighbor.travelCost() + neighbor.getHeuristicDistance()
			if neighbor.gotEgg {
				neighborHeuristic = 0
			}

			// neighbor is better than our current, we should explore this one
			if neighborHeuristic <= curr.heuristicDistance {
				neighbor.heuristicDistance = neighborHeuristic
				//fmt.Println("push", neighbor.s.nextTickDirection, neighbor.heuristicDistance, fmt.Sprintf("%p", neighbor))
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

func (a *autoPlayNode) travelCost() int {
	return 1
}

func (a *autoPlayNode) getHeuristicDistance() int {
	currentPos := a.s.head.Value.(mapPoint)

	xDist := abs(currentPos.x - a.s.eggLoc.x)
	yDist := abs(currentPos.y - a.s.eggLoc.y)

	return xDist + yDist
}

func (a *autoPlayNode) copy() *autoPlayNode {
	newSnakeSight := &autoPlayNode{
		s: a.s.copy(),
	}

	return newSnakeSight
}

func abs(n int) int {
	y := n >> 31
	return (n ^ y) - y
}
