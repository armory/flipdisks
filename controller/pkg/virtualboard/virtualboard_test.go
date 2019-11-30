package virtualboard

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVirtualBoard_String(t *testing.T) {
	tests := []struct {
		name  string
		board VirtualBoard
		want  string // CAUTION: ignore the first newline, it's just easier to see in src
	}{
		{
			name: "draw black on first line",
			board: VirtualBoard{
				{1, 1, 1, 1, 1},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			want: strings.TrimPrefix(`
⚫️⚫️⚫️⚫️⚫️
⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️
`, "\n"),
		},
		{
			name: "draw black on left side",
			board: VirtualBoard{
				{1, 0, 0, 0, 0},
				{1, 0, 0, 0, 0},
				{1, 0, 0, 0, 0},
			},
			want: strings.TrimPrefix(`
⚫️⚪️⚪️⚪️⚪️
⚫️⚪️⚪️⚪️⚪️
⚫️⚪️⚪️⚪️⚪️
`, "\n"),
		},
		{
			name: "numbers >= 1 are black",
			board: VirtualBoard{
				{1, 0, 1, 0, 2, 3},
				{0, 0, 0, 0, 5, 4},
				{1, 0, 1, 10, 6, 0},
			},
			want: strings.TrimPrefix(`
⚫️⚪️⚫️⚪️⚫️⚫️
⚪️⚪️⚪️⚪️⚫️⚫️
⚫️⚪️⚫️⚫️⚫️⚪️
`, "\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.board.String()
			assert.Equal(t, []byte(tt.want), []byte(got), fmt.Sprintf("expected\n%s\ngot\n%s", tt.want, got))
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		height int
		width  int
	}
	tests := []struct {
		name string
		args args
		want *VirtualBoard
	}{
		{
			name: "create a new 2x4",
			args: args{
				width:  2,
				height: 4,
			},
			want: &VirtualBoard{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.width, tt.args.height)
			assert.Equal(t, tt.want, got, fmt.Sprintf("expected\n%s\ngot\n%s", tt.want, got))
		})
	}
}
