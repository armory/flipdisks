package virtualboard

import (
	"strings"

	"flipdisks/pkg/fontmap"
	"github.com/kyokomi/emoji"
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

	//blackDot := emoji.Sprint(":black_circle:")
	//whiteDot := emoji.Sprint(":white_circle:")
	blackDot := " x "
	whiteDot := "   "
	//blackDot := emoji.Sprint(":snake:")
	//whiteDot := emoji.Sprint(":white_large_square:")

	var line strings.Builder

	xLen := len(board)
	yLen := len(board[0])
	dotLen := max(len([]byte(blackDot)), len([]byte(whiteDot))) // there's an extra 3 bytes for "️ Variation Selector-16" see https://emojipedia.org/variation-selector-16/
	newLines := yLen + 1                                        // new line at the end
	line.Grow(xLen*yLen*dotLen + newLines)

	// add coord system
	//iToEmoji := []string{":zero:", ":one:", ":two:", ":three:", ":four:", ":five:", ":six:", ":seven:", ":eight:", ":nine:"}
	iToEmoji := []string{" 0 ", " 1 ", " 2 ", " 3 ", " 4 ", " 5 ", " 6 ", " 7 ", " 8 ", " 9 "}
	line.WriteString("   ")
	for x := 0; x < xLen; x++ {
		line.WriteString(emoji.Sprint(iToEmoji[x%10]))
	}
	line.WriteString("\n")

	for y := 0; y < yLen; y++ { // do y first since we're drawing top down
		line.WriteString(emoji.Sprint(iToEmoji[y%10]))
		for x := 0; x < xLen; x++ { // then left right
			//line.WriteString(strconv.Itoa(board[x][y]))
			if board[x][y] == 1 {
				line.WriteString(blackDot)
			} else if board[x][y] == 3 {
				line.WriteString(" o ")
			} else {
				line.WriteString(whiteDot)
			}
		}
		line.WriteString(emoji.Sprint(iToEmoji[y%10])+"\n")
	}

	line.WriteString("   ")
	for x := 0; x < xLen; x++ {
		line.WriteString(emoji.Sprint(iToEmoji[x%10]))
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
