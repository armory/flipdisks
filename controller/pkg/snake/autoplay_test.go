package snake

import (
	"container/ring"
	"flipdisks/pkg/virtualboard"
	"testing"
)

func Test_abs(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "positive returns positive",
			args: args{10},
			want: 10,
		},
		{
			name: "negative returns positive",
			args: args{-10},
			want: 10,
		},
		{
			name: "0 returns 0",
			args: args{-10},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intAbs(tt.args.n); got != tt.want {
				t.Errorf("intAbs() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestSnake_AutoPlay(t *testing.T) {
	tests := []struct {
		name   string

		snakeMock snaker
	}{
		{
			name:   "play",
			//snakeMock : func (ctrl *gomock.Controller) snaker {
			//
			//},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//ctrl := gomock.NewController(t)

			//s, _ := New(30, 30, 2, 3)
			s, _ := New(28+28, 7*10, 2, 3)

			//s.snaker = tt.snakeMock
			s.AutoPlay()
		})
	}
}

func TestSnake_AutoPlayRandomly(t *testing.T) {
	type fields struct {
		boardHeight       int
		boardWidth        int
		startOffset       int
		snakeLength       int
		head              *ring.Ring
		tail              *ring.Ring
		nextTickDirection direction
		eggLoc            mapPoint
		deathBoundaries   deathBoundary
		GameBoard         *virtualboard.VirtualBoard
		snaker            snaker
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "play",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for {
				s, _ := New(11, 11, 2, 4)
				s.AutoPlayRandomly()
			}
		})
	}
}

//function reconstruct_path(cameFrom, current)
//    total_path := {current}
//    while current in cameFrom.Keys:
//        current := cameFrom[current]
//        total_path.prepend(current)
//    return total_path
//
//// A* finds a path from start to goal.
//// h is the heuristic function. h(n) estimates the cost to reach goal from node n.
//function A_Star(start, goal, h)
//    // The set of discovered nodes that may need to be (re-)expanded.
//    // Initially, only the start node is known.
//    openSet := {start}
//
//    // For node n, cameFrom[n] is the node immediately preceding it on the cheapest path from start to n currently known.
//    cameFrom := an empty map
//
//    // For node n, gScore[n] is the cost of the cheapest path from start to n currently known.
//    gScore := map with default value of Infinity
//    gScore[start] := 0
//
//    // For node n, fScore[n] := gScore[n] + h(n).
//    fScore := map with default value of Infinity
//    fScore[start] := h(start)
//
//    while openSet is not empty
//        current := the node in openSet having the lowest fScore[] value
//        if current = goal
//            return reconstruct_path(cameFrom, current)
//
//        openSet.Remove(current)
//        for each neighbor of current
//            // d(current,neighbor) is the weight of the edge from current to neighbor
//            // tentative_gScore is the distance from start to the neighbor through current
//            tentative_gScore := gScore[current] + d(current, neighbor)
//            if tentative_gScore < gScore[neighbor]
//                // This path to neighbor is better than any previous one. Record it!
//                cameFrom[neighbor] := current
//                gScore[neighbor] := tentative_gScore
//                fScore[neighbor] := gScore[neighbor] + h(neighbor)
//                if neighbor not in openSet
//                    openSet.add(neighbor)
//
//    // Open set is empty but goal was never reached
//    return failure

