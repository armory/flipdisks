package virtualboard

import "flipdisks/pkg/fontmap"

// use the canvas coord system instead of cartesian like in math!
//   reason being is that it'a intuitive from a data perspective
//   board[x][y] is how we'll be defining out data structure.
//
// (0,0) starts on upper left
// as x gets larger, we'll be moving to the right
// as y gets larger, we'll be moving down the screen
//
// data wise, it looks weird though, but just remember [x][y]
//   board[x][y] = [
//     x0: [ y0 y1 y2...]
//     x1: [ y0 y1 y2...]
//     x2: [ y0 y1 y2...]
//   ]
type VirtualBoard []fontmap.Row

func New(width, height int) *VirtualBoard {
	b := make(VirtualBoard, width)
	for x := 0; x < width; x++ {
		b[x] = make(fontmap.Row, height)
	}

	return &b
}

// String is used to draw on the screen
func (b *VirtualBoard) String() string {
	board := *b
	line := ""
	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[x]); y++ {
			if board[x][y] >= 1 {
				line += "⚫️"
			} else {
				line += "⚪️"
			}
		}
		line += "\n"
	}

	return line
}
