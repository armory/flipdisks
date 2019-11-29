package virtualboard

import (
	"fmt"
	"testing"
)

func TestVirtualBoard_String(t *testing.T) {
	tests := []struct {
		name  string
		board VirtualBoard
		want  string
	}{
		{
			name: "simple",
			board: VirtualBoard{
				{1, 0, 1, 0, 1, 1},
				{0, 0, 0, 0, 1, 1},
				{1, 0, 1, 1, 1, 0},
			},
			want: `⚫️⚪️⚫️⚪️⚫️⚫️
⚪️⚪️⚪️⚪️⚫️⚫️
⚫️⚪️⚫️⚫️⚫️⚪️
`,
		},
		{
			name: "numbers >= 1 are black",
			board: VirtualBoard{
				{1, 0, 1, 0, 2, 3},
				{0, 0, 0, 0, 5, 4},
				{1, 0, 1, 10, 6, 0},
			},
			want: `⚫️⚪️⚫️⚪️⚫️⚫️
⚪️⚪️⚪️⚪️⚫️⚫️
⚫️⚪️⚫️⚫️⚫️⚪️
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.board.String()
			if got != tt.want {
				fmt.Println("expected")
				fmt.Println(tt.want)
				fmt.Println("got")
				fmt.Println(got)
				t.Fail()
			}
		})
	}
}
