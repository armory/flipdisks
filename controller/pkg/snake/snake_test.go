package snake

import (
	"fmt"
	"reflect"
	"testing"

	"flipdisks/pkg/virtualboard"
	"github.com/stretchr/testify/assert"
)

func Test_addEgg(t *testing.T) {
	type args struct {
		boardPointer *virtualboard.VirtualBoard
	}
	tests := []struct {
		name             string
		args             args
		expectedGameOver bool
		expectedBoard    *virtualboard.VirtualBoard
	}{
		{
			name: "it can find a spot for an egg on an empty board",
			args: args{
				boardPointer: &virtualboard.VirtualBoard{
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
				},
			},
			expectedGameOver: false,
			expectedBoard: &virtualboard.VirtualBoard{
				{0, 0, 0, 0},
				{0, 0, 0, 1},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
		},
		{
			name: "it still works when the board is 50% full horizontally",
			args: args{
				boardPointer: &virtualboard.VirtualBoard{
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{1, 1, 1, 1},
					{1, 1, 1, 1},
				},
			},
			expectedGameOver: false,
			expectedBoard: &virtualboard.VirtualBoard{
				{0, 0, 0, 1},
				{0, 0, 0, 0},
				{1, 1, 1, 1},
				{1, 1, 1, 1},
			},
		},
		{
			name: "it still works when the board is 50% full vertically",
			args: args{
				boardPointer: &virtualboard.VirtualBoard{
					{0, 0, 1, 1},
					{0, 0, 1, 1},
					{0, 0, 1, 1},
					{0, 0, 1, 1},
				},
			},
			expectedGameOver: false,
			expectedBoard: &virtualboard.VirtualBoard{
				{1, 0, 1, 1},
				{0, 0, 1, 1},
				{0, 0, 1, 1},
				{0, 0, 1, 1},
			},
		},
		{
			name: "it still works when the board is has only 1 spot left",
			args: args{
				boardPointer: &virtualboard.VirtualBoard{
					{1, 1, 1, 1},
					{1, 1, 1, 1},
					{1, 1, 0, 1},
					{1, 1, 1, 1},
				},
			},
			expectedGameOver: false,
			expectedBoard: &virtualboard.VirtualBoard{
				{1, 1, 1, 1},
				{1, 1, 1, 1},
				{1, 1, 1, 1},
				{1, 1, 1, 1},
			},
		},
		{
			name: "if a full board is found, game is over",
			args: args{
				boardPointer: &virtualboard.VirtualBoard{
					{1, 1, 1, 1},
					{1, 1, 1, 1},
					{1, 1, 1, 1},
					{1, 1, 1, 1},
				},
			},
			expectedGameOver: true,
			expectedBoard: &virtualboard.VirtualBoard{
				{1, 1, 1, 1},
				{1, 1, 1, 1},
				{1, 1, 1, 1},
				{1, 1, 1, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedGameOver, addEgg(tt.args.boardPointer), "returning game status was not expected")
			assert.Equal(t, *tt.expectedBoard, *tt.args.boardPointer)
		})
	}
}

func Test_addSnake(t *testing.T) {
	type args struct {
		boardPointer *virtualboard.VirtualBoard
	}
	tests := []struct {
		name          string
		args          args
		expectedBoard *virtualboard.VirtualBoard
	}{
		{
			name: "adds a snake on a 11x11",
			args: args{
				boardPointer: &virtualboard.VirtualBoard{
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
			expectedBoard: &virtualboard.VirtualBoard{
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
			addSnake(tt.args.boardPointer)
			if !reflect.DeepEqual(tt.expectedBoard, tt.args.boardPointer) {
				fmt.Println("expected")
				fmt.Println(tt.expectedBoard)
				fmt.Println("got")
				fmt.Println(tt.args.boardPointer)
				t.Fail()
			}
		})
	}
}
