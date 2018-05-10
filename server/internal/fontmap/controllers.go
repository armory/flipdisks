package fontmap

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// todo return the LIST of fonts names and meta data (without charmap)
func GetFonts(context *gin.Context) {
	context.JSON(200, TI84)
}

func generateSpace(width int, height int) Letter {
	var space Letter
	for j := 0; j < height; j++ {
		var row Row
		for i := 0; i < width; i++ {
			row = append(row, "⚪️")
		}

		space = append(space, row)
	}

	return space
}

func addKerning(letter Letter, amountOfKerning int) Letter {
	var kernedLetter Letter
	for _, row := range letter {
		kernedRow := row
		for j := 0; j < amountOfKerning; j++ {
			kernedRow = append(kernedRow, "⚪️")
		}
		kernedLetter = append(kernedLetter, kernedRow)
	}

	fmt.Println(kernedLetter)

	return kernedLetter
}

// todo for Jimmy
// - kerning (space between letters)
// - a space need to render to bunch of white dots
func Render(context *gin.Context) {
	// this is the requesting payload we expect our clients to send us
	type requestPayloadType struct {
		FontName   string `json:"fontName"`
		Text       string `json:"text"`
		Kerning    int    `json:"kerning"`
		SpaceWidth int    `json:"spaceWidth"`
	}
	payload := &requestPayloadType{}
	context.Bind(payload) // store the data in our type

	//this will determine our max height should be

	letterBoardHeight := TI84.Metadata.AverageHeight

	for tallerCharacter, height := range TI84.Metadata.TallerCharacters {
		if letterBoardHeight < height && strings.Contains(payload.Text, tallerCharacter) {
			letterBoardHeight = height
		}
	}

	// map each character to a rendered font using the font map
	board := [][]Letter{
		[]Letter{}, // line
	}

	line := &board[0]

	for _, char := range strings.Split(payload.Text, "") {
		if char == " " {
			*line = append(*line, generateSpace(payload.SpaceWidth, letterBoardHeight))
		} else if char == "\n" {
			board = append(board, []Letter{})
			line = &board[len(board)-1]
		} else {
			*line = append(*line, addKerning(TI84.Charmap[char], payload.Kerning))
		}
	}

	context.JSON(200, board)
}
