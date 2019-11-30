package flipboard

import (
	"regexp"

	"flipdisks/pkg/fontmap"
	"flipdisks/pkg/options"
	"flipdisks/pkg/virtualboard"
)

func renderTextToVirtualBoard(msg *options.FlipboardMessageOptions, board *Flipboard) *virtualboard.VirtualBoard {
	msgCharsAsDots := fontmap.Render(msg.Message)
	virtualBoardPointer := CreateVirtualBoard(board.PanelInfo.PhysicallyDisplayedWidth, len(board.PanelAddressesLayout[0]), msgCharsAsDots, msg.Message)
	virtualBoard := *virtualBoardPointer

	// todo, it would be nice to just invert it without through the whole board again
	// handle inverting for words
	if msg.Inverted {
		for _, row := range virtualBoard {
			for charIndex, x := range row {
				if x == 0 {
					row[charIndex] = 1
				} else {
					row[charIndex] = 0
				}
			}
		}
	}

	return virtualBoardPointer
}

func CreateVirtualBoard(panelWidth int, numberOfPanelsWide int, msgCharsAsDots []fontmap.Letter, msg string) *virtualboard.VirtualBoard {
	// we have to convert our long array of dotCharacters to a virtual board
	var longestLine, lineNumber int
	longestLine = 0
	lineNumber = 0
	lineMaxWidth := panelWidth * numberOfPanelsWide
	virtualBoard := virtualboard.VirtualBoard{}

	// join the letters together to form one long string
	for charIndexInMessage := 0; charIndexInMessage < len(msgCharsAsDots); charIndexInMessage++ {
		charAsDots := msgCharsAsDots[charIndexInMessage]

		// handle line breaks
		if msg[charIndexInMessage] == '\n' && charAsDots == nil {
			lineNumber++
			longestLine = 0
			continue
		}

		// try to word break, if it doesn't work, then we'lll need to character break
		if msg[charIndexInMessage] == ' ' {
			unprocessedStringMsg := msg[charIndexInMessage:] // msg will look something like: "   bbb"
			unprocessedDotMessage := msgCharsAsDots[charIndexInMessage:]

			matchPos := regexp.MustCompile(`\S+`).FindStringIndex(unprocessedStringMsg) // matchPos[0] will be the first "b"
			nextDotWord := unprocessedDotMessage[matchPos[0]:matchPos[1]]

			// find the width of dots for the word
			wordDotWidth := 0
			for _, dotChar := range nextDotWord {
				if len(dotChar) > 0 {
					wordDotWidth += len(dotChar[0])
				}
			}

			// since we're breaking on the word, we should discard all the whitespace before the word
			if longestLine+wordDotWidth > lineMaxWidth {
				lineNumber++
				longestLine = 0

				// advance our pointer to the beginning of the next word
				charIndexInMessage += matchPos[0]
				charAsDots = msgCharsAsDots[charIndexInMessage]
			}
		} else if longestLine+len(charAsDots[0]) > lineMaxWidth {
			// if there's no spaces, and the word is super long, let's fallback and do a character break
			lineNumber++
			longestLine = 0
		}

		// write character to the virtual board
		for charRowIndex, charRow := range charAsDots {
			boardCharRowIndex := charRowIndex + (lineNumber * fontmap.TI84.Metadata.MaxHeight)

			// create all missing rows from the virtual board, up to our current boardCharRowIndex
			for len(virtualBoard) <= boardCharRowIndex {
				virtualBoard = append(virtualBoard, fontmap.Row{})
			}

			virtualBoard[boardCharRowIndex] = append(virtualBoard[boardCharRowIndex], charRow...)

			// keep track of the longest char row for the line
			if longestLine < len(virtualBoard[boardCharRowIndex]) {
				longestLine = len(virtualBoard[boardCharRowIndex])
			}

		}
	}
	return &virtualBoard
}
