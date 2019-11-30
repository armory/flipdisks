package snake

import (
	"container/ring"
	"fmt"
	"testing"

	"flipdisks/pkg/virtualboard"
	"github.com/stretchr/testify/assert"
)

//func Test_addEgg(t *testing.T) {
//	type fields struct {
//		board virtualboard.VirtualBoard
//	}
//	tests := []struct {
//		name             string
//		fields           fields
//		expectedGameOver bool
//		expectedBoard    virtualboard.VirtualBoard
//	}{
//		{
//			name: "it can find a spot for an egg on an empty board",
//			fields: fields{
//				board: virtualboard.VirtualBoard{
//					{0, 0, 0, 0},
//					{0, 0, 0, 0},
//					{0, 0, 0, 0},
//					{0, 0, 0, 0},
//				},
//			},
//			expectedGameOver: false,
//			expectedBoard: virtualboard.VirtualBoard{
//				{0, 0, 0, 0},
//				{0, 0, 0, 1},
//				{0, 0, 0, 0},
//				{0, 0, 0, 0},
//			},
//		},
//
//		{
//			name: "it still works when the board is 50% full horizontally",
//			fields: fields{
//				board: virtualboard.VirtualBoard{
//					{0, 0, 0, 0},
//					{0, 0, 0, 0},
//					{1, 1, 1, 1},
//					{1, 1, 1, 1},
//				},
//			},
//			expectedGameOver: false,
//			expectedBoard: virtualboard.VirtualBoard{
//				{0, 0, 0, 1},
//				{0, 0, 0, 0},
//				{1, 1, 1, 1},
//				{1, 1, 1, 1},
//			},
//		},
//		{
//			name: "it still works when the board is 50% full vertically",
//			fields: fields{
//				board: virtualboard.VirtualBoard{
//					{0, 0, 1, 1},
//					{0, 0, 1, 1},
//					{0, 0, 1, 1},
//					{0, 0, 1, 1},
//				},
//			},
//			expectedGameOver: false,
//			expectedBoard: virtualboard.VirtualBoard{
//				{1, 0, 1, 1},
//				{0, 0, 1, 1},
//				{0, 0, 1, 1},
//				{0, 0, 1, 1},
//			},
//		},
//		{
//			name: "it still works when the board is has only 1 spot left",
//			fields: fields{
//				board: virtualboard.VirtualBoard{
//					{1, 1, 1, 1},
//					{1, 1, 1, 1},
//					{1, 1, 0, 1},
//					{1, 1, 1, 1},
//				},
//			},
//			expectedGameOver: false,
//			expectedBoard: virtualboard.VirtualBoard{
//				{1, 1, 1, 1},
//				{1, 1, 1, 1},
//				{1, 1, 1, 1},
//				{1, 1, 1, 1},
//			},
//		},
//		{
//			name: "if a full board is found, game is over",
//			fields: fields{
//				board: virtualboard.VirtualBoard{
//					{1, 1, 1, 1},
//					{1, 1, 1, 1},
//					{1, 1, 1, 1},
//					{1, 1, 1, 1},
//				},
//			},
//			expectedGameOver: true,
//			expectedBoard: virtualboard.VirtualBoard{
//				{1, 1, 1, 1},
//				{1, 1, 1, 1},
//				{1, 1, 1, 1},
//				{1, 1, 1, 1},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Snake{
//				board: tt.fields.board,
//			}
//
//			assert.Equal(t, tt.expectedGameOver, s.addEgg(), "returning game status was not expected")
//			assert.Equal(t, tt.expectedBoard, tt.fields.board)
//		})
//	}
//}
//
//func TestSnake_startGame(t *testing.T) {
//	type fields struct {
//		board       virtualboard.VirtualBoard
//		tailOffset  int
//		snakeLength int
//	}
//	tests := []struct {
//		name          string
//		fields        fields
//		expectedBoard virtualboard.VirtualBoard
//	}{
//		{
//			name: "adds a snake on a 11x11",
//			fields: fields{
//				tailOffset:  2,
//				snakeLength: 3,
//				board: virtualboard.VirtualBoard{
//					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//				},
//			},
//			expectedBoard: virtualboard.VirtualBoard{
//				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//				{0, 0, 2, 2, 1, 0, 0, 0, 3, 0, 0},
//				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Snake{
//				board:       tt.fields.board,
//				tailOffset:  tt.fields.tailOffset,
//				snakeLength: tt.fields.snakeLength,
//			}
//
//			s.startGame()
//
//			if !reflect.DeepEqual(tt.expectedBoard, tt.fields.board) {
//				fmt.Println("expected")
//				fmt.Println(tt.expectedBoard)
//				fmt.Println("got")
//				fmt.Println(tt.fields.board)
//				t.Fail()
//			}
//		})
//	}
//}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "play",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			New()
		})
	}
}

func TestSnake_setupGame(t *testing.T) {
	type fields struct {
		boardHeight       int
		boardWidth        int
		startOffset       int
		snakeLength       int
		snakeHead         *ring.Ring
		snakeTail         *ring.Ring
		nextTickDirection direction
		deathBoundaries   deathBoundary
	}
	tests := []struct {
		name         string
		fields       fields
		expectations func(t *testing.T, s *Snake)
	}{
		{
			name: "setup the a 11x11 game correctly",
			fields: fields{
				boardHeight: 11,
				boardWidth:  11,
				startOffset: 2,
				snakeLength: 3,
			},
			expectations: func(t *testing.T, s *Snake) {
				// snake heading east
				assert.Equal(t, mapPoint{4, 5}, s.head.Value)
				assert.Equal(t, mapPoint{2, 5}, s.tail.Value)
				assert.Equal(t, east, s.nextTickDirection)

				// egg in the right spot
				assert.Equal(t, mapPoint{8, 5}, s.eggLoc)

				// deathBoundaries and snake
				sTemp := &Snake{boardHeight: 11, boardWidth: 11}
				sTemp.addOutsideBoundaries() // tested somewhere else
				sTemp.addBoundary(2, 5)      // snake
				sTemp.addBoundary(3, 5)      // snake
				sTemp.addBoundary(4, 5)      // snake
				assert.Equal(t, sTemp.deathBoundaries, s.deathBoundaries)

				// make sure GameBoard is what we expect
				expectedGameBoard := &virtualboard.VirtualBoard{
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 1, 1, 1, 0, 0, 0, 1, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				}
				assert.Equal(t, expectedGameBoard, s.GameBoard, fmt.Sprintf("expected\n%s\ngot\n%s", expectedGameBoard, s.GameBoard))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Snake{
				boardHeight:       tt.fields.boardHeight,
				boardWidth:        tt.fields.boardWidth,
				startOffset:       tt.fields.startOffset,
				snakeLength:       tt.fields.snakeLength,
				head:              tt.fields.snakeHead,
				tail:              tt.fields.snakeTail,
				nextTickDirection: tt.fields.nextTickDirection,
				deathBoundaries:   tt.fields.deathBoundaries,
			}
			s.setupGame()
			tt.expectations(t, s)
		})
	}
}

func TestSnake_addBoundary(t *testing.T) {
	type fields struct {
		deathBoundaries deathBoundary
	}
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name         string
		fields       fields
		args         []args
		expectations func(t *testing.T, s *Snake)
	}{
		{
			name:   "adds a boundary when empty",
			fields: fields{},
			args:   []args{{5, 10}},
			expectations: func(t *testing.T, s *Snake) {
				assert.Equal(t, deathBoundary{
					5: {10: wallExists{}},
				}, s.deathBoundaries)
			},
		},
		{
			name:   "adds multiple",
			fields: fields{},
			args: []args{
				{5, 10},
				{4, 10},
				{3, 10},
			},
			expectations: func(t *testing.T, s *Snake) {
				assert.Equal(t, deathBoundary{
					5: {10: wallExists{}},
					4: {10: wallExists{}},
					3: {10: wallExists{}},
				}, s.deathBoundaries)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Snake{
				deathBoundaries: tt.fields.deathBoundaries,
			}

			for _, arg := range tt.args {
				s.addBoundary(arg.x, arg.y)
			}

			tt.expectations(t, s)
		})
	}
}

func TestSnake_addOutsideBoundaries(t *testing.T) {
	type fields struct {
		boardHeight int
		boardWidth  int
	}
	tests := []struct {
		name         string
		fields       fields
		expectations func(t *testing.T, s *Snake)
	}{
		{
			name: "adds a 2x2 grid",
			fields: fields{
				boardHeight: 2,
				boardWidth:  2,
			},
			expectations: func(t *testing.T, s *Snake) {
				assert.Equal(t, deathBoundary{
					-1: {-1: wallExists{}, 0: wallExists{}, 1: wallExists{}, 2: wallExists{}}, // top
					0: {
						-1: wallExists{}, // left
						2:  wallExists{}, // right
					},
					1: {
						-1: wallExists{}, // left
						2:  wallExists{}, // right
					},
					2: {-1: wallExists{}, 0: wallExists{}, 1: wallExists{}, 2: wallExists{}}, // bottom
				}, s.deathBoundaries)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Snake{
				boardHeight: tt.fields.boardHeight,
				boardWidth:  tt.fields.boardWidth,
			}

			s.addOutsideBoundaries()

			tt.expectations(t, s)
		})
	}
}
