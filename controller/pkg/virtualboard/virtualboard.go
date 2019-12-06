package virtualboard

import (
	"strings"

	"flipdisks/pkg/fontmap"
)

// use the canvas coord system instead of cartesian, like in math!
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// String is used to draw on the screen
//   note: storage is transposed because we're doing [x][y]
func (b *VirtualBoard) String() string {
	board := *b

	//const onDot = emoji.Sprint(":black_circle:")   // lib can be found: "github.com/kyokomi/emoji"
	//const offDot = emoji.Sprint(":white_circle:")
	//const blackDot = "⚫️"
	//const whiteDot = "⚪️"
	const onDot = " x "
	const offDot = "   "
	const eggSpace = " o " // special value for snake (egg)

	var line strings.Builder                                // apparently this is faster than string appending
	xLen := len(board)                                      // +2 for the header and footer
	yLen := len(board[0])                                   // +2 for the left and right side
	dotLen := max(len([]byte(onDot)), len([]byte(offDot)))  // emoji's sometimes can have an extra 3 bytes for "️ Variation Selector-16". See https://emojipedia.org/variation-selector-16/
	newLineChar := yLen + 1                                 // new line at the end
	line.Grow((xLen+2)*(yLen+2)*dotLen + (newLineChar + 2)) // +1 line for each top/bottom/left/right headers

	// add coord system on top
	intString := []string{" 0 ", " 1 ", " 2 ", " 3 ", " 4 ", " 5 ", " 6 ", " 7 ", " 8 ", " 9 "}
	line.WriteString("   ")
	for x := 0; x < xLen; x++ {
		line.WriteString(intString[x%10])
	}
	line.WriteString("\n")

	for y := 0; y < yLen; y++ { // do y first since we're drawing top down
		rowHeader := intString[y%10]
		line.WriteString(rowHeader) // left side
		for x := 0; x < xLen; x++ { // left to right
			switch board[x][y] {
			case 1:
				line.WriteString(onDot)
			case 3:
				line.WriteString(eggSpace)
			default:
				line.WriteString(offDot)
			}
		}
		line.WriteString(rowHeader + "\n") // right side
	}

	// add coord system on bottom
	line.WriteString("   ")
	for x := 0; x < xLen; x++ {
		line.WriteString(intString[x%10])
	}
	line.WriteString("\n")

	return line.String()
}

// Helper function to convert source-coded matrices to data-encoded matrices, really only for tests
// When drawing arrays in using {{}, {}} format, it's actually [y][x]
// When accessing it everywhere else, we're doing [x][y]
func (b VirtualBoard) Transpose() *VirtualBoard {
	dataEncodedMatrix := New(len(b[0]), len(b))

	for x := 0; x < len(*dataEncodedMatrix); x++ {
		for y := 0; y < len((*dataEncodedMatrix)[0]); y++ {
			(*dataEncodedMatrix)[x][y] = (b)[y][x]
		}
	}
	return dataEncodedMatrix
}
