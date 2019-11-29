package snake

import (
	"fmt"
	"reflect"
	"testing"

	"flipdisks/pkg/virtualboard"
	"github.com/stretchr/testify/assert"
)

func Test_addEgg(t *testing.T) {
	type fields struct {
		board virtualboard.VirtualBoard
	}
	tests := []struct {
		name             string
		fields           fields
		expectedGameOver bool
		expectedBoard    virtualboard.VirtualBoard
	}{
		{
			name: "it can find a spot for an egg on an empty board",
			fields: fields{
				board: virtualboard.VirtualBoard{
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
				},
			},
			expectedGameOver: false,
			expectedBoard: virtualboard.VirtualBoard{
				{0, 0, 0, 0},
				{0, 0, 0, 1},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
		},

		{
			name: "it still works when the board is 50% full horizontally",
			fields: fields{
				board: virtualboard.VirtualBoard{
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{1, 1, 1, 1},
					{1, 1, 1, 1},
				},
			},
			expectedGameOver: false,
			expectedBoard: virtualboard.VirtualBoard{
				{0, 0, 0, 1},
				{0, 0, 0, 0},
				{1, 1, 1, 1},
				{1, 1, 1, 1},
			},
		},
		{
			name: "it still works when the board is 50% full vertically",
			fields: fields{
				board: virtualboard.VirtualBoard{
					{0, 0, 1, 1},
					{0, 0, 1, 1},
					{0, 0, 1, 1},
					{0, 0, 1, 1},
				},
			},
			expectedGameOver: false,
			expectedBoard: virtualboard.VirtualBoard{
				{1, 0, 1, 1},
				{0, 0, 1, 1},
				{0, 0, 1, 1},
				{0, 0, 1, 1},
			},
		},
		{
			name: "it still works when the board is has only 1 spot left",
			fields: fields{
				board: virtualboard.VirtualBoard{
					{1, 1, 1, 1},
					{1, 1, 1, 1},
					{1, 1, 0, 1},
					{1, 1, 1, 1},
				},
			},
			expectedGameOver: false,
			expectedBoard: virtualboard.VirtualBoard{
				{1, 1, 1, 1},
				{1, 1, 1, 1},
				{1, 1, 1, 1},
				{1, 1, 1, 1},
			},
		},
		{
			name: "if a full board is found, game is over",
			fields: fields{
				board: virtualboard.VirtualBoard{
					{1, 1, 1, 1},
					{1, 1, 1, 1},
					{1, 1, 1, 1},
					{1, 1, 1, 1},
				},
			},
			expectedGameOver: true,
			expectedBoard: virtualboard.VirtualBoard{
				{1, 1, 1, 1},
				{1, 1, 1, 1},
				{1, 1, 1, 1},
				{1, 1, 1, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Snake{
				board: tt.fields.board,
			}

			assert.Equal(t, tt.expectedGameOver, s.addEgg(), "returning game status was not expected")
			assert.Equal(t, tt.expectedBoard, tt.fields.board)
		})
	}
}

func TestSnake_addSnake(t *testing.T) {
	type fields struct {
		board       virtualboard.VirtualBoard
		tailOffset  int
		snakeLength int
	}
	tests := []struct {
		name          string
		fields        fields
		expectedBoard virtualboard.VirtualBoard
	}{
		{
			name: "adds a snake on a 11x11",
			fields: fields{
				tailOffset:  2,
				snakeLength: 3,
				board: virtualboard.VirtualBoard{
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
			},
			expectedBoard: virtualboard.VirtualBoard{
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Snake{
				board:       tt.fields.board,
				tailOffset:  tt.fields.tailOffset,
				snakeLength: tt.fields.snakeLength,
			}

			s.addSnake()

			if !reflect.DeepEqual(tt.expectedBoard, tt.fields.board) {
				fmt.Println("expected")
				fmt.Println(tt.expectedBoard)
				fmt.Println("got")
				fmt.Println(tt.fields.board)
				t.Fail()
			}
		})
	}
}
