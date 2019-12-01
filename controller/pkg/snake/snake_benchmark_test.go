package snake

import (
	"testing"

	"flipdisks/pkg/virtualboard"
	"github.com/stretchr/testify/assert"
)

func BenchmarkSnake_addEgg(b *testing.B) {
	// this is a tad bit bigger than 4K which is 3840x2160
	// that's pretty big...
	width4k := 3840
	height4k := 2160

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
			name: "find the only empty spot in the board",
			fields: fields{
				boardWidth:  width4k,
				boardHeight: height4k,
				deathBoundaries: func() deathBoundary {
					b := deathBoundary{}
					for x := 0; x < width4k; x++ {
						for y := 0; y < height4k; y++ {
							if !(x == 100 && y == 200) {
								b.Add(x, y)
							}
						}
					}
					return b
				}(),
			},
			want: true,
			expect: func(t *testing.T, s *Snake) {
				assert.Equal(t, mapPoint{100, 200}, s.eggLoc)
			},
		},
		{
			name: "board is full, can't place egg",
			fields: fields{
				boardWidth:  width4k,
				boardHeight: height4k,
				deathBoundaries: func() deathBoundary {
					b := deathBoundary{}
					for x := 0; x < width4k; x++ {
						for y := 0; y < height4k; y++ {
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
		s := &Snake{
			boardWidth:      tt.fields.boardWidth,
			boardHeight:     tt.fields.boardHeight,
			deathBoundaries: tt.fields.deathBoundaries,
			GameBoard:       virtualboard.New(tt.fields.boardWidth, tt.fields.boardHeight),
		}
		s.snaker = s
		s.addOutsideBoundaries()

		b.Run(tt.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				b.ResetTimer()
				s.addEgg()
			}
		})
	}
}

