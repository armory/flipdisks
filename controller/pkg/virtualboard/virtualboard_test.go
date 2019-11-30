package virtualboard

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to convert source-coded matrices to data-encoded matrices
// When drawing arrays in using {{}, {}} format, it's actually [y][x]
// When accessing it everywhere else, we're doing [x][y]
func transposeVirtualBoard(sourceCodedMatrix VirtualBoard) VirtualBoard {
	dataEncodedMatrix := New(len(sourceCodedMatrix[0]), len(sourceCodedMatrix))

	for x := 0; x < len(*dataEncodedMatrix); x++ {
		for y := 0; y < len((*dataEncodedMatrix)[0]); y++ {
			(*dataEncodedMatrix)[x][y] = sourceCodedMatrix[y][x]
		}
	}
	return *dataEncodedMatrix
}

func TestVirtualBoard_String(t *testing.T) {
	tests := []struct {
		name  string
		board VirtualBoard
		want  string
	}{
		{
			name: "draw black on first line",
			board: transposeVirtualBoard(VirtualBoard{
				{1, 1, 1, 1, 1},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			}),
			want: strings.TrimPrefix(`
⚫️⚫️⚫️⚫️⚫️
⚪️⚪️⚪️⚪️⚪️
⚪️⚪️⚪️⚪️⚪️
`, "\n"),
		},
		{
			name: "draw black on left side",
			board: transposeVirtualBoard(VirtualBoard{
				{1, 0, 0, 0, 0},
				{1, 0, 0, 0, 0},
				{1, 0, 0, 0, 0},
			}),
			want: strings.TrimPrefix(`
⚫️⚪️⚪️⚪️⚪️
⚫️⚪️⚪️⚪️⚪️
⚫️⚪️⚪️⚪️⚪️
`, "\n"),
		},
		{
			name: "numbers >= 1 are black",
			board: transposeVirtualBoard(VirtualBoard{
				{1, 0, 1, 0, 2, 3},
				{0, 0, 0, 0, 5, 4},
				{1, 0, 1, 10, 6, 0},
			}),
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
			assert.Equalf(t, []byte(tt.want), []byte(got), "expected\n%s\ngot\n%s", tt.want, got)
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
			assert.Equalf(t, tt.want, got, "expected\n%s\ngot\n%s", tt.want, got)
		})
	}
}
