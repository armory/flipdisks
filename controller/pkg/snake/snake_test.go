package snake

import (
	"container/ring"
	"fmt"
	"testing"
	"time"

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
				sTemp := &Snake{boardHeight: 11, boardWidth: 11, deathBoundaries: deathBoundary{}}
				sTemp.addOutsideBoundaries()    // tested somewhere else
				sTemp.deathBoundaries.Add(2, 5) // snake
				sTemp.deathBoundaries.Add(3, 5) // snake
				sTemp.deathBoundaries.Add(4, 5) // snake
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

func TestDeathBoundary_Add(t *testing.T) {
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
			fields: fields{
				deathBoundaries: deathBoundary{},
			},
			args:   []args{{5, 10}},
			expectations: func(t *testing.T, s *Snake) {
				assert.Equal(t, deathBoundary{
					5: {10: wallExists{}},
				}, s.deathBoundaries)
			},
		},
		{
			name:   "adds multiple",
			fields: fields{
				deathBoundaries: deathBoundary{},
			},
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
				s.deathBoundaries.Add(arg.x, arg.y)
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
				boardHeight:     tt.fields.boardHeight,
				boardWidth:      tt.fields.boardWidth,
				deathBoundaries: deathBoundary{},
			}

			s.addOutsideBoundaries()

			tt.expectations(t, s)
		})
	}
}

func TestSnake_addEgg(t *testing.T) {
	// this is a tad bit bigger than 4K which is 3840x2160
	// that's pretty big...
	largestSupportedBoardWidthHeight := 3840

	type fields struct {
		boardHeight     int
		boardWidth      int
		startOffset     int
		snakeLength     int
		eggLoc          mapPoint
		deathBoundaries deathBoundary
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
		expect func(t *testing.T, s *Snake)
	}{
		{
			name: "randomly places an egg",
			fields: fields{
				boardWidth:  5,
				boardHeight: 5,
			},
			want: true,
			expect: func(t *testing.T, s *Snake) {
				assert.Equal(t, mapPoint{1, 3}, s.eggLoc)
			},
		},
		{
			name: "randomly places an egg in a really full board",
			fields: fields{
				boardWidth:  5,
				boardHeight: 5,
				deathBoundaries: deathBoundary{
					0: {0: wallExists{}, 1: wallExists{}, 2: wallExists{}, 3: wallExists{}, 4: wallExists{}},
					1: {0: wallExists{}, 1: wallExists{}, 2: wallExists{}, 3: wallExists{}, 4: wallExists{}},
					2: {0: wallExists{}, 1: wallExists{}, 2: wallExists{}, 3: wallExists{}, 4: wallExists{}},
					3: {0: wallExists{}, 1: wallExists{}, 2: wallExists{}, 3: wallExists{}, 4: wallExists{}},
					4: {0: wallExists{}, 1: wallExists{}, 3: wallExists{}, 4: wallExists{}},
				},
			},
			want: true,
			expect: func(t *testing.T, s *Snake) {
				assert.Equal(t, mapPoint{4, 2}, s.eggLoc)
			},
		},
		{
			name: "can't place an egg because the board is full",
			fields: fields{
				boardWidth:  5,
				boardHeight: 5,
				deathBoundaries: deathBoundary{
					0: {0: wallExists{}, 1: wallExists{}, 2: wallExists{}, 3: wallExists{}, 4: wallExists{}},
					1: {0: wallExists{}, 1: wallExists{}, 2: wallExists{}, 3: wallExists{}, 4: wallExists{}},
					2: {0: wallExists{}, 1: wallExists{}, 2: wallExists{}, 3: wallExists{}, 4: wallExists{}},
					3: {0: wallExists{}, 1: wallExists{}, 2: wallExists{}, 3: wallExists{}, 4: wallExists{}},
					4: {0: wallExists{}, 1: wallExists{}, 2: wallExists{}, 3: wallExists{}, 4: wallExists{}},
				},
			},
			want:   false,
			expect: func(t *testing.T, s *Snake) {},
		},
		{
			name: "it can find that 1 spot in a very large board (essentially a load test)",
			fields: fields{
				boardWidth:  largestSupportedBoardWidthHeight,
				boardHeight: largestSupportedBoardWidthHeight,
				deathBoundaries: func() deathBoundary {
					b := deathBoundary{}
					for x := 0; x < largestSupportedBoardWidthHeight; x++ {
						for y := 0; y < largestSupportedBoardWidthHeight; y++ {
							if !(x == 393 && y == 488) {
								b.Add(x, y)
							}
						}
					}
					return b
				}(),
			},
			want: true,
			expect: func(t *testing.T, s *Snake) {
				assert.Equal(t, mapPoint{393, 488}, s.eggLoc)
			},
		},
		{
			name: "can't place an egg because the board is full in a really large board",
			fields: fields{
				boardWidth:  largestSupportedBoardWidthHeight,
				boardHeight: largestSupportedBoardWidthHeight,
				deathBoundaries: func() deathBoundary {
					b := deathBoundary{}
					for x := 0; x < largestSupportedBoardWidthHeight; x++ {
						for y := 0; y < largestSupportedBoardWidthHeight; y++ {
							b.Add(x, y)
						}
					}
					return b
				}(),
			},
			want:   false,
			expect: func(t *testing.T, s *Snake) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan struct{})

			go func() {
				s := &Snake{
					boardHeight:     tt.fields.boardHeight,
					boardWidth:      tt.fields.boardWidth,
					startOffset:     tt.fields.startOffset,
					snakeLength:     tt.fields.snakeLength,
					eggLoc:          tt.fields.eggLoc,
					deathBoundaries: tt.fields.deathBoundaries,
				}
				s.addOutsideBoundaries()

				got := s.addEgg()
				if got != tt.want {
					t.Errorf("addEgg() = %v, want %v", got, tt.want)
				}

				tt.expect(t, s)

				done <- struct{}{}
			}()

			select {
			case <-time.After(250 * time.Millisecond):
				t.Fatal("timed out trying to add an egg! fix it!")
			case <-done:
			}
		})
	}
}
