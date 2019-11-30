package snake

import (
	"container/ring"
	"fmt"
	"strings"
	"testing"
	"time"

	"flipdisks/pkg/virtualboard"
	"github.com/stretchr/testify/assert"
)

//func TestNew(t *testing.T) {
//	tests := []struct {
//		name string
//	}{
//		{
//			name: "play",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			New()
//		})
//	}
//}

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
				// snake heading East
				assert.Equal(t, mapPoint{4, 5}, s.head.Value)
				assert.Equal(t, mapPoint{2, 5}, s.tail.Value)
				assert.Equal(t, East, s.nextTickDirection)

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
				expectedGameBoard := (&virtualboard.VirtualBoard{
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 1, 1, 1, 0, 0, 0, 3, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				}).Transpose()
				assert.Equalf(t, expectedGameBoard, s.GameBoard, "expected\n%s\ngot\n%s", expectedGameBoard, s.GameBoard)
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
			name: "adds a boundary when empty",
			fields: fields{
				deathBoundaries: deathBoundary{},
			},
			args: []args{{5, 10}},
			expectations: func(t *testing.T, s *Snake) {
				assert.Equal(t, deathBoundary{
					5: {10: wallExists{}},
				}, s.deathBoundaries)
			},
		},
		{
			name: "adds multiple",
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

func TestDeathBoundary_Remove(t *testing.T) {
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
			name: "removes a boundary point, but there's still y left",
			fields: fields{
				deathBoundaries: deathBoundary{
					5: {10: wallExists{}, 11: wallExists{}},
				},
			},
			args: []args{{5, 10}},
			expectations: func(t *testing.T, s *Snake) {
				assert.Equal(t, deathBoundary{
					5: {11: wallExists{}},
				}, s.deathBoundaries)
			},
		},
		{
			name: "removes a boundary x point, but no more x left",
			fields: fields{
				deathBoundaries: deathBoundary{
					5:  {10: wallExists{}},
					99: {1: wallExists{}},
				},
			},
			args: []args{{5, 10}},
			expectations: func(t *testing.T, s *Snake) {
				assert.Equal(t, deathBoundary{
					5:  {},
					99: {1: wallExists{}},
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
				s.deathBoundaries.Remove(arg.x, arg.y)
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
				assert.Equal(t, mapPoint{1, 2}, s.eggLoc)
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
				s.GameBoard=virtualboard.New(s.boardWidth, s.boardHeight)
				s.addOutsideBoundaries()

				got := s.addEgg()
				if got != tt.want {
					t.Errorf("addEgg() = %v, want %v", got, tt.want)
				}

				tt.expect(t, s)

				done <- struct{}{}
			}()

			select {
			case <-time.After(500 * time.Millisecond):
				t.Fatal("timed out trying to add an egg! fix it!")
			case <-done:
			}
		})
	}
}

func TestSnake_moveSnake(t *testing.T) {
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
	}
	type args struct {
		gotLonger bool
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func(s *Snake)

		expect func(t *testing.T, s *Snake)
	}{
		{
			name: "moves to the east and checks death",
			fields: fields{
				boardHeight:       11,
				boardWidth:        11,
				startOffset:       2,
				snakeLength:       3,
				nextTickDirection: East,
			},
			args: args{false},
			expect: func(t *testing.T, s *Snake) {
				assert.Equal(t, mapPoint{5, 5}, s.head.Value, "head is not in the correct pos")
				_, exists := s.deathBoundaries[5][5]
				assert.True(t, exists, "head is not setup as death")

				assert.Equal(t, mapPoint{3, 5}, s.tail.Value, "tail is not in the correct pos")
				_, exists = s.deathBoundaries[2][5]
				assert.False(t, exists, "old tail is still set as death")

				expectedGameBoard := strings.TrimPrefix(`
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚫️⚫️⚫️⚪️⚪️⚫️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
`, "\n")
				assert.Equalf(t, expectedGameBoard, s.GameBoard.String(), "expected\n%s\ngot\n%s", expectedGameBoard, s.GameBoard)
			},
		},
		{
			name: "moves to the north",
			fields: fields{
				boardHeight:       11,
				boardWidth:        11,
				startOffset:       2,
				snakeLength:       3,
				nextTickDirection: North,
			},
			args: args{false},
			expect: func(t *testing.T, s *Snake) {
				assert.Equal(t, mapPoint{4, 4}, s.head.Value, "head is not in the correct pos")
				assert.Equal(t, mapPoint{3, 5}, s.tail.Value, "tail is not in the correct pos")

				expectedGameBoard := strings.TrimPrefix(`
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚫️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚫️⚫️⚪️⚪️⚪️⚫️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
`, "\n")
				assert.Equalf(t, []byte(expectedGameBoard), []byte(s.GameBoard.String()), "expected\n%s\ngot\n%s", expectedGameBoard, s.GameBoard)
			},
		},
		{
			name: "moves to the south",
			fields: fields{
				boardHeight:       11,
				boardWidth:        11,
				startOffset:       2,
				snakeLength:       3,
				nextTickDirection: South,
			},
			args: args{false},
			expect: func(t *testing.T, s *Snake) {
				assert.Equal(t, mapPoint{4, 6}, s.head.Value, "head is not in the correct pos")
				assert.Equal(t, mapPoint{3, 5}, s.tail.Value, "tail is not in the correct pos")

				expectedGameBoard := strings.TrimPrefix(`
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚫️⚫️⚪️⚪️⚪️⚫️⚪️⚪️
⚪️⚪️⚪️⚪️⚫️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
`, "\n")
				assert.Equalf(t, []byte(expectedGameBoard), []byte(s.GameBoard.String()), "expected\n%s\ngot\n%s", expectedGameBoard, s.GameBoard)
			},
		},
		{
			name: "moves to the west",
			fields: fields{
				boardHeight:       11,
				boardWidth:        11,
				startOffset:       2,
				snakeLength:       3,
				nextTickDirection: West,
			},
			beforeTest: func(s *Snake) {
				s.nextTickDirection = South
				s.moveSnake(false)
				s.nextTickDirection = South
				s.moveSnake(false)
				s.nextTickDirection = South
				s.moveSnake(false)
			},
			args: args{false},
			expect: func(t *testing.T, s *Snake) {
				assert.Equal(t, mapPoint{3, 8}, s.head.Value, "head is not in the correct pos")
				assert.Equal(t, mapPoint{4, 7}, s.tail.Value, "tail is not in the correct pos")

				expectedGameBoard := strings.TrimPrefix(`
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚫️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚫️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚫️⚫️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
`, "\n")
				assert.Equalf(t, []byte(expectedGameBoard), []byte(s.GameBoard.String()), "expected\n%s\ngot\n%s", expectedGameBoard, s.GameBoard)
			},
		},
		{
			name: "moves to the west",
			fields: fields{
				boardHeight:       11,
				boardWidth:        11,
				startOffset:       2,
				snakeLength:       3,
				nextTickDirection: West,
			},
			beforeTest: func(s *Snake) {
				s.nextTickDirection = South
				s.moveSnake(false)
				s.nextTickDirection = South
				s.moveSnake(false)
				s.nextTickDirection = South
				s.moveSnake(false)
			},
			args: args{false},
			expect: func(t *testing.T, s *Snake) {
				assert.Equal(t, mapPoint{3, 8}, s.head.Value, "head is not in the correct pos")
				assert.Equal(t, mapPoint{4, 7}, s.tail.Value, "tail is not in the correct pos")

				expectedGameBoard := strings.TrimPrefix(`
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚫️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚫️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚫️⚫️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
`, "\n")
				assert.Equalf(t, []byte(expectedGameBoard), []byte(s.GameBoard.String()), "expected\n%s\ngot\n%s", expectedGameBoard, s.GameBoard)
			},
		},
		{
			name: "moves to the east, got longer",
			fields: fields{
				boardHeight:       11,
				boardWidth:        11,
				startOffset:       2,
				snakeLength:       3,
				nextTickDirection: East,
			},
			args: args{true},
			expect: func(t *testing.T, s *Snake) {
				assert.Equal(t, mapPoint{5, 5}, s.head.Value)
				assert.Equal(t, mapPoint{2, 5}, s.tail.Value)

				expectedGameBoard := strings.TrimPrefix(`
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚫️⚫️⚫️⚫️⚪️⚪️⚫️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️⚪️
`, "\n")
				assert.Equalf(t, expectedGameBoard, s.GameBoard.String(), "expected\n%s\ngot\n%s", expectedGameBoard, s.GameBoard)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Snake{
				boardHeight: tt.fields.boardHeight,
				boardWidth:  tt.fields.boardWidth,
				startOffset: tt.fields.startOffset,
				snakeLength: tt.fields.snakeLength,
			}
			s.setupGame()
			if tt.beforeTest != nil {
				tt.beforeTest(s)
			}
			//fmt.Println(s.GameBoard)

			s.nextTickDirection = tt.fields.nextTickDirection
			s.moveSnake(tt.args.gotLonger)
			tt.expect(t, s)
		})
	}
}

func TestSnake_Tick(t *testing.T) {
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
	}
	type ticks struct {
		nextDirection direction
		assert        func(t *testing.T, s *Snake) bool
		expectedBoard *virtualboard.VirtualBoard
	}
	tests := []struct {
		name           string
		fields         fields
		ticks          []ticks
		wantIsGameOver bool
		wantGameWin    bool
	}{
		{
			name: "simple game - do nothing",
			fields: fields{
				boardWidth:  11,
				boardHeight: 11,
				snakeLength: 3,
				startOffset: 2,
			},
			ticks: []ticks{
				{nextDirection: East},
				{nextDirection: East},
				{nextDirection: East},
				{nextDirection: East},
				{nextDirection: East},
				{nextDirection: East},
				{nextDirection: East},
			},
			wantIsGameOver: true,
			wantGameWin:    false,
		},
		{
			name: "simple game - moving around",
			fields: fields{
				boardWidth:  11,
				boardHeight: 11,
				snakeLength: 3,
				startOffset: 2,
			},
			ticks: []ticks{
				{
					nextDirection: East,
					expectedBoard: (virtualboard.VirtualBoard{
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 1, 1, 1, 0, 0, 3, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					}).Transpose(),
				},
				{
					nextDirection: North,
					expectedBoard: (virtualboard.VirtualBoard{
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 1, 1, 0, 0, 3, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					}).Transpose(),
				},
				{
					nextDirection: North,
					expectedBoard: (virtualboard.VirtualBoard{
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 1, 0, 0, 3, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					}).Transpose(),
				},
				{
					nextDirection: East,
					expectedBoard: (virtualboard.VirtualBoard{
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					}).Transpose(),
				},
				{
					nextDirection: East,
					expectedBoard: (virtualboard.VirtualBoard{
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					}).Transpose(),
				},
				{
					nextDirection: East,
					expectedBoard: (virtualboard.VirtualBoard{
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					}).Transpose(),
				},
				{
					nextDirection: South,
					expectedBoard: (virtualboard.VirtualBoard{
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					}).Transpose(),
				},
				{
					nextDirection: South,
					expectedBoard: (virtualboard.VirtualBoard{
						{0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					}).Transpose(),
				},
				{
					nextDirection: West,
					expectedBoard: (virtualboard.VirtualBoard{
						{0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					}).Transpose(),
				},
				{
					nextDirection: West,
					expectedBoard: (virtualboard.VirtualBoard{
						{0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
						{0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					}).Transpose(),
				},
			},
			wantIsGameOver: false,
			wantGameWin:    false,
		},
	}
	for _, tt := range tests {
		if tt.name != "simple game - moving around" {
			t.Log("WARNING!!! WE'RE SKIPPING THIS TEST! This was prob commited by accident", "snake_test.go:909-11/30/19-04:49")
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			s := &Snake{
				boardHeight:       tt.fields.boardHeight,
				boardWidth:        tt.fields.boardWidth,
				startOffset:       tt.fields.startOffset,
				snakeLength:       tt.fields.snakeLength,
				head:              tt.fields.head,
				tail:              tt.fields.tail,
				nextTickDirection: tt.fields.nextTickDirection,
				eggLoc:            tt.fields.eggLoc,
				deathBoundaries:   tt.fields.deathBoundaries,
				GameBoard:         tt.fields.GameBoard,
			}
			s.setupGame()

			var gotIsGameOver, gotGameWin bool

			fmt.Printf("start:\n%s", s.GameBoard)
			for i, tick := range tt.ticks {
				gotIsGameOver, gotGameWin = s.Tick(tick.nextDirection)
				fmt.Printf("tick # %d:\n%s", i, s.GameBoard)

				if tick.expectedBoard != nil {
					assert.Equal(t, tick.expectedBoard, s.GameBoard, "expected\n%s\ngot\n%s", tick.expectedBoard, s.GameBoard.String())
				}

				if tick.assert != nil {
					tick.assert(t, s)
				}
			}

			assert.Equal(t, tt.wantIsGameOver, gotIsGameOver, "gameOver?")
			assert.Equal(t, tt.wantGameWin, gotGameWin, "won game?")
		})
	}
}

func Test_deathBoundary_IsBoundary(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		b    deathBoundary
		args args
		want bool
	}{
		{
			name: "can find a boundary easily",
			b: deathBoundary{
				1: {1: wallExists{}},
				2: {2: wallExists{}},
			},
			args: args{1, 1},
			want: true,
		},
		{
			name: "if it x doesn't exist, it's not a boundary",
			b: deathBoundary{
				1: {1: wallExists{}},
				2: {2: wallExists{}},
			},
			args: args{99, 99},
			want: false,
		},
		{
			name: "if it y doesn't exist, it's not a boundary",
			b: deathBoundary{
				1: {1: wallExists{}},
				2: {2: wallExists{}},
			},
			args: args{1, 99},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsBoundary(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("IsBoundary() = %v, want %v", got, tt.want)
			}
		})
	}
}
