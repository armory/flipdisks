package virtualboard

import "github.com/armory/flipdisks/controller/pkg/fontmap"

type VirtualBoard []fontmap.Row

func (board VirtualBoard) String() string {
	line := ""
	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[x]); y++ {
			if board[x][y] == 1 {
				line += "⚫️"
			} else {
				line += "⚪️"
			}
		}
		line += "\n"
	}

	return line
}
