package snake

import (
	"container/ring"
	"fmt"
	"github.com/go-test/deep"
	"strings"
	"testing"

	"flipdisks/pkg/virtualboard"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
		mockSnaker   func(s *Snake, ctrl *gomock.Controller) snaker
		expectations func(t *testing.T, s *Snake)
	}{
		{
			name: "setup the a 11x11 game correctly",
			fields: fields{
				boardHeight: 11,
				boardWidth:  14,
				startOffset: 2,
				snakeLength: 3,
			},
			expectations: func(t *testing.T, s *Snake) {
				// snake is facing East
				assert.Equal(t, mapPoint{4, 5}, s.head.Value)
				assert.Equal(t, mapPoint{2, 5}, s.tail.Value)

				// egg in the right spot
				assert.Equal(t, mapPoint{11, 5}, s.eggLoc)

				// deathBoundaries and snake
				sTemp := &Snake{boardHeight: 11, boardWidth: 14, deathBoundaries: deathBoundary{}}
				sTemp.addOutsideBoundaries()    // tested somewhere else
				sTemp.deathBoundaries.Add(2, 5) // snake
				sTemp.deathBoundaries.Add(3, 5) // snake
				sTemp.deathBoundaries.Add(4, 5) // snake
				assert.Equal(t, sTemp.deathBoundaries, s.deathBoundaries)

				// make sure GameBoard is what we expect
				expectedGameBoard := (&virtualboard.VirtualBoard{
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 3, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				}).Transpose()
				assertGameBoard(t, expectedGameBoard, s.GameBoard)
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
				boardWidth:  2,
				boardHeight: 2,
			},
			expectations: func(t *testing.T, s *Snake) {
				assert.Equal(t, deathBoundary{
					-1: {-1: wallExists{}, 0: wallExists{}, 1: wallExists{}, 2: wallExists{}}, // top
					0:  {-1: wallExists{} /*       gameBoardInHere       */, 2: wallExists{}},
					1:  {-1: wallExists{} /*       gameBoardInHere       */, 2: wallExists{}},
					2:  {-1: wallExists{}, 0: wallExists{}, 1: wallExists{}, 2: wallExists{}}, // bottom
				}, s.deathBoundaries)
			},
		},
		{
			name: "adds a 4x2 grid",
			fields: fields{
				boardWidth:  2,
				boardHeight: 4,
			},
			//
			expectations: func(t *testing.T, s *Snake) {
				b := deathBoundary{}
				// left side
				b.Add(-1, -1)
				b.Add(-1, 0)
				b.Add(-1, 1)
				b.Add(-1, 2)
				b.Add(-1, 3)
				b.Add(-1, 4)

				// right side
				b.Add(2, -1)
				b.Add(2, 0)
				b.Add(2, 1)
				b.Add(2, 2)
				b.Add(2, 3)
				b.Add(2, 4)

				// top
				b.Add(-1, -1)
				b.Add(0, -1)
				b.Add(1, -1)
				b.Add(2, -1)

				// bottom
				b.Add(-1, 2)
				b.Add(0, 2)
				b.Add(1, 2)
				b.Add(2, 2)

				assert.Equal(t, b, s.deathBoundaries)
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
					4: {0: wallExists{}, 1: wallExists{} /*            */, 3: wallExists{}, 4: wallExists{}},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Snake{
				boardHeight:     tt.fields.boardHeight,
				boardWidth:      tt.fields.boardWidth,
				startOffset:     tt.fields.startOffset,
				snakeLength:     tt.fields.snakeLength,
				eggLoc:          tt.fields.eggLoc,
				deathBoundaries: tt.fields.deathBoundaries,
				GameBoard:       virtualboard.New(tt.fields.boardWidth, tt.fields.boardHeight),
			}
			s.snaker = s
			s.addOutsideBoundaries()

			got := s.addEgg()

			if got != tt.want {
				t.Errorf("addEgg() = %v, want %v", got, tt.want)
			}
			tt.expect(t, s)
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

				expectedGameBoard := (virtualboard.VirtualBoard{
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
				}).Transpose()
				assertGameBoard(t, expectedGameBoard, s.GameBoard)
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

				expectedGameBoard := (virtualboard.VirtualBoard{
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 1, 1, 0, 0, 0, 3, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				}).Transpose()
				assertGameBoard(t, expectedGameBoard, s.GameBoard)
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

				expectedGameBoard := (virtualboard.VirtualBoard{
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 1, 1, 0, 0, 0, 3, 0, 0},
					{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				}).Transpose()
				assertGameBoard(t, expectedGameBoard, s.GameBoard)
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

				expectedGameBoard := (virtualboard.VirtualBoard{
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				}).Transpose()
				assertGameBoard(t, expectedGameBoard, s.GameBoard)
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

				expectedGameBoard := (virtualboard.VirtualBoard{
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				}).Transpose()
				assertGameBoard(t, expectedGameBoard, s.GameBoard)
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

				expectedGameBoard := (virtualboard.VirtualBoard{
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 1, 1, 1, 1, 0, 0, 3, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				}).Transpose()
				assertGameBoard(t, expectedGameBoard, s.GameBoard)
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
		mockSnaker    func(s *Snake, ctrl *gomock.Controller) snaker
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
					mockSnaker: func(s *Snake, ctrl *gomock.Controller) snaker {
						sMock := NewMocksnaker(ctrl)
						sMock.EXPECT().eggNextLoc().Return(9, 3)
						sMock.EXPECT().nextHeadLoc(s.nextTickDirection).Return(s.nextHeadLoc(s.nextTickDirection))
						return sMock
					},
					expectedBoard: (virtualboard.VirtualBoard{
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 1, 1, 3, 0},
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
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 1, 3, 0},
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
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						{0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0},
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
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

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

			t.Logf("start:\n%s", s.GameBoard)
			for i, tick := range tt.ticks {
				s.snaker = s
				if tick.mockSnaker != nil {
					s.snaker = tick.mockSnaker(s, ctrl)
				}
				gotIsGameOver, gotGameWin = s.Tick(tick.nextDirection)
				t.Logf("tick # %d:\n%s", i, s.GameBoard)

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

func TestSnake_checkGameStatus(t *testing.T) {
	type fields struct {
		boardWidth      int
		boardHeight     int
		snakeLength     int
		deathBoundaries deathBoundary
	}
	tests := []struct {
		name           string
		fields         fields
		mockSnaker     func(ctrl *gomock.Controller, s *Snake) snaker
		wantIsGameOver bool
		wantGameWin    bool
	}{
		{
			name: "didn't hit a boundary, game continues on",
			fields: fields{
				boardHeight:     11,
				boardWidth:      11,
				snakeLength:     5,
				deathBoundaries: deathBoundary{},
			},
			mockSnaker: func(ctrl *gomock.Controller, s *Snake) snaker {
				sMock := NewMocksnaker(ctrl)
				sMock.EXPECT().nextHeadLoc(s.nextTickDirection).Return(mapPoint{5, 5})
				return sMock
			},
			wantIsGameOver: false,
			wantGameWin:    false,
		},
		{
			name: "hit a boundary, snake is dead",
			fields: fields{
				boardHeight: 11,
				boardWidth:  11,
				snakeLength: 5,
				deathBoundaries: deathBoundary{
					4: {3: wallExists{}},
				},
			},
			mockSnaker: func(ctrl *gomock.Controller, s *Snake) snaker {
				sMock := NewMocksnaker(ctrl)
				sMock.EXPECT().nextHeadLoc(s.nextTickDirection).Return(mapPoint{4, 3})
				return sMock
			},
			wantIsGameOver: true,
			wantGameWin:    false,
		},
		{
			name: "we won the game!",
			fields: fields{
				boardHeight:     11,
				boardWidth:      10,
				snakeLength:     11 * 10,
				deathBoundaries: deathBoundary{},
			},
			mockSnaker: func(ctrl *gomock.Controller, s *Snake) snaker {
				sMock := NewMocksnaker(ctrl)
				sMock.EXPECT().nextHeadLoc(s.nextTickDirection).Return(mapPoint{4, 3})
				return sMock
			},
			wantIsGameOver: true,
			wantGameWin:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := &Snake{
				boardWidth:      tt.fields.boardWidth,
				boardHeight:     tt.fields.boardHeight,
				snakeLength:     tt.fields.snakeLength,
				deathBoundaries: tt.fields.deathBoundaries,
			}
			s.snaker = s
			if tt.mockSnaker != nil {
				s.snaker = tt.mockSnaker(ctrl, s)
			}

			gotIsGameOver, gotGameWin := s.checkGameStatus()
			if gotIsGameOver != tt.wantIsGameOver {
				t.Errorf("checkGameStatus() gotIsGameOver = %v, want %v", gotIsGameOver, tt.wantIsGameOver)
			}
			if gotGameWin != tt.wantGameWin {
				t.Errorf("checkGameStatus() gotGameWin = %v, want %v", gotGameWin, tt.wantGameWin)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		boardHeight int
		boardWidth  int
		startOffset int
		snakeLength int
	}
	tests := []struct {
		name   string
		args   args
		expect func(t *testing.T, s *Snake)
	}{
		{
			name: "start a new game!",
			args: args{
				boardHeight: 11,
				boardWidth:  20,
				startOffset: 2,
				snakeLength: 4,
			},
			expect: func(t *testing.T, s *Snake) {
				expectedGameBoard := (virtualboard.VirtualBoard{
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				}).Transpose()
				assertGameBoard(t, expectedGameBoard, s.GameBoard)

				assert.Equal(t, East, s.nextTickDirection)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := New(tt.args.boardHeight, tt.args.boardWidth, tt.args.startOffset, tt.args.snakeLength)

			tt.expect(t, got)
		})
	}
}

func TestSnake_copy(t *testing.T) {
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
		disableGameBoard  bool
		GameBoard         *virtualboard.VirtualBoard
		snaker            snaker
	}
	tests := []struct {
		name     string
		fields   fields
		original func() *Snake
		want     func(t *testing.T, orig *Snake, copy *Snake)
	}{
		{
			name: "should copy the struct correctly",
			original: func() *Snake {
				s, _ := New(10, 20, 2, 5)
				s.Tick(East)
				s.Tick(South)
				s.Tick(South)
				s.Tick(East)
				s.DisableGameBoard()
				return s
			},
			want: func(t *testing.T, orig *Snake, copy *Snake) {
				// debugging
				// debugging
				orig.EnableGameBoard()
				copy.EnableGameBoard()
				fmt.Println(orig.GameBoard)
				fmt.Println(copy.GameBoard)
				orig.DisableGameBoard()
				copy.DisableGameBoard()

				assert.Equal(t, 10, copy.boardHeight)
				assert.Equal(t, 20, copy.boardWidth)
				assert.Equal(t, 5, copy.snakeLength) // hasn't eaten the egg

				// calling DisableGameBoard removes it
				assert.Equal(t, true, copy.disableGameBoard)
				assert.Nil(t, copy.GameBoard)

				// do we have the snake at the right points? (from tail to head)
				snakeSegments := []mapPoint{
					{7, 5},
					{8, 5},
					{8, 6},
					{8, 7},
					{9, 7},
				}
				walker := copy.tail
				for i, seg := range snakeSegments {
					assert.Equal(t, seg, walker.Value.(mapPoint), "segment %d is not correct", i)
					walker = walker.Next()
				}
				assert.Equal(t, walker.Prev().Value, copy.head.Value) // the for loop goes 1 too far

				// make sure they're different rings by adding a random value to the copy and checking their Next()/Prev() respectively
				copy.head.Next().Value = mapPoint{99, 99}
				assert.NotEqual(t, copy.head.Next(), orig.head.Next(), "head is pointing to the same ring")

				copy.tail.Prev().Value = mapPoint{99, 99}
				assert.NotEqual(t, copy.tail.Prev(), orig.tail.Prev(), "tail is pointing to the same ring")

				// egg
				assert.Equal(t, mapPoint{15, 5}, copy.eggLoc)

				// death boundaries
				assert.Empty(t, deep.Equal(copy.deathBoundaries, orig.deathBoundaries), "deathBoundaries are not the same")
				copy.deathBoundaries.Add(1, 1)
				deathBoundaryErrors := deep.Equal(copy.deathBoundaries, orig.deathBoundaries)
				assert.NotEmptyf(t, deathBoundaryErrors, "seems like deathBoundaries are the same!\n%s", strings.Join(deathBoundaryErrors, "\n"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := tt.original()

			tt.want(t, orig, orig.copy())
		})
	}
}

func TestSnake_DisableGameBoard(t *testing.T) {
	type fields struct {
		boardHeight       int
		boardWidth        int
		startOffset       int
		snakeLength       int
		head              *ring.Ring
		tail              *ring.Ring
		eggLoc            mapPoint
		nextTickDirection direction
		deathBoundaries   deathBoundary
		disableGameBoard  bool
		GameBoard         *virtualboard.VirtualBoard
		snaker            snaker
	}
	tests := []struct {
		name   string
		fields fields
		expect func(t *testing.T, s *Snake)
	}{
		{
			name: "can disable GameBoard",
			fields: fields{
				GameBoard: virtualboard.New(2, 3),
			},
			expect: func(t *testing.T, s *Snake) {
				assert.True(t, s.disableGameBoard)
				assert.Empty(t, s.GameBoard)
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
				head:              tt.fields.head,
				tail:              tt.fields.tail,
				eggLoc:            tt.fields.eggLoc,
				nextTickDirection: tt.fields.nextTickDirection,
				deathBoundaries:   tt.fields.deathBoundaries,
				disableGameBoard:  tt.fields.disableGameBoard,
				GameBoard:         tt.fields.GameBoard,
				snaker:            tt.fields.snaker,
			}

			s.DisableGameBoard()

			tt.expect(t, s)
		})
	}
}

func TestSnake_EnableGameBoard(t *testing.T) {
	type fields struct {
		boardHeight       int
		boardWidth        int
		startOffset       int
		snakeLength       int
		head              *ring.Ring
		tail              *ring.Ring
		eggLoc            mapPoint
		nextTickDirection direction
		deathBoundaries   deathBoundary
		disableGameBoard  bool
		GameBoard         *virtualboard.VirtualBoard
		snaker            snaker
	}
	tests := []struct {
		name   string
		s      func() *Snake
		expect func(t *testing.T, s *Snake)
	}{
		{
			name: "can enable a disabled gameboard",
			s: func() *Snake {
				s, _ := New(10, 10, 1, 3)
				return s
			},
			expect: func(t *testing.T, s *Snake) {
				assert.False(t, s.disableGameBoard)
				assert.NotEmpty(t, s.GameBoard)

				sTemp, _ := New(10, 10, 1, 3)
				assertGameBoard(t, sTemp.GameBoard, s.GameBoard)
			},
		},
		{
			name: "can be ran multiple times",
			s: func() *Snake {
				s, _ := New(10, 10, 1, 3)
				return s
			},
			expect: func(t *testing.T, s *Snake) {
				s.EnableGameBoard()

				assert.False(t, s.disableGameBoard)
				assert.NotEmpty(t, s.GameBoard)

				sTemp, _ := New(10, 10, 1, 3)
				assertGameBoard(t, sTemp.GameBoard, s.GameBoard)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.s()
			s.EnableGameBoard()
			tt.expect(t, s)
		})
	}
}
