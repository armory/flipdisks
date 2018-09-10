package flipboard

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/armory/flipdisks/controller/pkg/image"
	"github.com/armory/flipdisks/controller/pkg/options"
	"github.com/armory/flipdisks/controller/pkg/virtualboard"
	"github.com/kevinawoo/flipdots/panel"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Flipboard struct {
	panels               *[][]panel.Panel
	PanelInfo            PanelInfo
	PanelAddressesLayout [][]PanelAddress
	displayQueue         []*options.FlipboardMessageOptions
	countdownDate        string
	newMessage           chan bool
	msgCurrentlyPlaying  bool
}

type Opts func(*Flipboard) error

func NewFlipboard(info PanelInfo, layout [][]PanelAddress, opts ...Opts) (*Flipboard, error) {
	panels, err := CreatePanels(info, layout)
	if err != nil {
		return &Flipboard{}, errors.New("couldn't create panels: " + err.Error())
	}

	board := Flipboard{
		panels:               panels,
		PanelInfo:            info,
		PanelAddressesLayout: layout,
		newMessage:           make(chan bool),
	}

	for _, opt := range opts {
		err := opt(&board)
		if err != nil {
			log.Error("couldn't set options: " + err.Error())
		}
	}

	return &board, nil
}

func CountdownDate(date string) Opts {
	return func(flipboard *Flipboard) error {
		if date == "" {
			return errors.New("countdown date called without date being set")
		}

		flipboard.countdownDate = date

		fmt.Println("starting countdown clock")
		go func() {
			for {
				if len(flipboard.displayQueue) == 0 && flipboard.msgCurrentlyPlaying == false {
					tick := flipboard.getNextCountdown()
					flipboard.Enqueue(&tick)
				}
				time.Sleep(time.Duration(time.Second * 1))
			}
		}()

		return nil
	}
}

func (b *Flipboard) Enqueue(msg *options.FlipboardMessageOptions) {
	fmt.Printf("Enqueuing Message: %+v\n", msg.Message)
	b.displayQueue = append(b.displayQueue, msg)
	b.newMessage <- true
}

func (b *Flipboard) dequeue() *options.FlipboardMessageOptions {
	var msg *options.FlipboardMessageOptions
	if len(b.displayQueue) > 0 {
		msg, b.displayQueue = b.displayQueue[0], b.displayQueue[1:]
	}

	return msg
}

func Play(board *Flipboard) {
	log.Info("listening")
	for {
		select {
		case <-board.newMessage:
			board.msgCurrentlyPlaying = true
			fmt.Println("playing")
			msg := board.dequeue()
			fmt.Println("dequeed")
			DisplayMessageToPanels(board, msg)

			fmt.Printf("keeping message displayed for: %dms ...\n", msg.DisplayTime)
			time.Sleep(time.Millisecond * time.Duration(msg.DisplayTime))
			fmt.Println("Done! Listening for next message...")
			board.msgCurrentlyPlaying = false
		}
	}
}

func DisplayMessageToPanels(board *Flipboard, msg *options.FlipboardMessageOptions) {
	if msg.Message == "debug all panels" || msg.Message == "debug panels" {
		msg.DisplayTime = 0
		board.DebugPanelAddressByGoingInOrder()
		return
	}
	if strings.Contains(msg.Message, "debug panel") {
		panelAddress, _ := strconv.Atoi(strings.Replace(msg.Message, "debug panel ", "", -1))
		msg.DisplayTime = 0
		board.DebugSinglePanel(panelAddress)
		return
	}

	// we got a virtualBoard yay! Lets just display it!
	if msg.VirtualBoard != nil {
		displayVirtualBoardToPhysicalBoard(msg, msg.VirtualBoard, board)
		return
	}

	maxWidth := uint(board.PanelInfo.PanelHeight * len(board.PanelAddressesLayout[0]))
	maxHeight := uint(board.PanelInfo.PanelWidth * len(board.PanelAddressesLayout))

	fmt.Printf("Enqueuing Message: %+v\n", msg.Message)

	gifUrls := image.GetGifUrl(msg.Message)
	if gifUrls != nil {
		for _, gifUrl := range gifUrls {
			fmt.Println("Got gif! rendering...")

			msg.DisplayTime = 0 // we'll be controlling the frame display time
			frames, _ := image.ConvertGifFromURLToVirtualBoard(gifUrl, maxWidth, maxHeight, msg.Inverted, msg.BWThreshold)

			for frameIndex, frame := range frames.Flipboards {
				frameDuration := frames.Delay[frameIndex]

				// a gif really is 1 "message", so we're not going to enqueue it, because someone else could put in a random message in it
				displayVirtualBoardToPhysicalBoard(msg, frame, board)

				time.Sleep(frameDuration)
				time.Sleep(1 * time.Second)
			}
		}
	} else if image.IsPlainImageUrl(msg.Message) {
		v := image.ConvertImageUrlToVirtualBoard(maxWidth, maxHeight, msg.Message, msg.Inverted, msg.BWThreshold)
		displayVirtualBoardToPhysicalBoard(msg, v, board)
	} else { // plain text
		v := renderTextToVirtualBoard(msg, board)
		displayVirtualBoardToPhysicalBoard(msg, v, board)
	}
}

func displayVirtualBoardToPhysicalBoard(msg *options.FlipboardMessageOptions, vBoardPointer *virtualboard.VirtualBoard, board *Flipboard) {
	virtualBoard := *vBoardPointer

	setPhysicalBoardFill(msg, virtualBoard, board)

	// set alignment options
	msg.XAlign, msg.YAlign = options.GetAlignOptions(msg.Align)

	fmt.Println(virtualBoard)

	// the library flipped height and width by accident, we'll work around it
	panelWidth := board.PanelInfo.PanelHeight
	panelHeight := board.PanelInfo.PanelWidth
	// convert virtual virtualBoard to a physical virtualBoard
	boardWidth := panelWidth * len(board.PanelAddressesLayout[0])
	boardHeight := panelHeight * len(board.PanelAddressesLayout)
	xOffSet, yOffSet := findOffSets(msg, &virtualBoard, boardWidth, boardHeight)
	for y := 0; y < len(virtualBoard); y++ {
		for x := 0; x < len(virtualBoard[y]); x++ {
			// which dot should we set?
			panelXCoord := (x + xOffSet) / panelWidth
			panelYCoord := (y + yOffSet) / panelHeight
			dotXCoord := (x + xOffSet) % panelWidth
			dotYCoord := (y + yOffSet) % panelHeight

			if dotXCoord < 0 || dotYCoord < 0 || panelXCoord < 0 || panelYCoord < 0 {
				continue
			}

			if panelYCoord >= len(board.PanelAddressesLayout) {
				log.Printf("Warning: row %d, exceeds specified HEIGHT %d, dropping the rest of it.", x, panelHeight)
				continue
			}

			if panelXCoord >= len(board.PanelAddressesLayout[panelYCoord]) {
				log.Printf("Warning: cell(%d,%d) exceeds specified WIDTH %d, dropping the rest of it.", y, x, panelWidth)
				continue
			}

			//log.Printf("Setting panel(%d,%d), adddress %d, dot(%d,%d) with %t", panelYCoord, panelXCoord, p.Address, dotYCoord, dotXCoord, dotValue)

			// there's a bug in this library, where x and y are flipped. we need to handle this later
			//p := panels[panelYCoord][panelXCoord]
			p := board.GetPanel(panelYCoord, panelXCoord)
			dotValue := virtualBoard[y][x] == 1
			p.Set(dotYCoord, dotXCoord, dotValue)
		}
	}
	// send our virtual panels to the physical virtualBoard
	board.Send()
}

func setPhysicalBoardFill(msg *options.FlipboardMessageOptions, virtualBoard virtualboard.VirtualBoard, board *Flipboard) {
	fill := msg.Fill == "true"
	// if no fill is provided, let's try to set autofill
	if msg.Fill == "" {
		var sum int

		// Go across the top to add up all the values
		for x := range virtualBoard[0] {
			sum += virtualBoard[0][x]
		}

		// go across the bottom to add up all the values
		for x := range virtualBoard[len(virtualBoard)-1] {
			sum += virtualBoard[len(virtualBoard)-1][x]
		}

		// go on the left and right side to add up all the values
		for _, row := range virtualBoard {
			// sometimes the row will be empty, because of a \n, let's just ignore it
			if len(row) > 0 {
				sum += row[0] // left y going down
			}

			// if for some reason it's just a single row, we'll have already taken care of adding the sum before
			if len(row) > 1 {
				sum += row[len(row)-1] // right y going down
			}
		}

		height := len(virtualBoard)
		width := len(virtualBoard[0])
		fill = float32(sum)/float32(2*(width+height)) >= .5 // magic number
		//fmt.Println("setting autofill to be: ", fill)
	}
	// set the fill value
	board.SetAll(fill)
}

func findOffSets(options *options.FlipboardMessageOptions, vBoardPointer *virtualboard.VirtualBoard, boardWidth, boardHeight int) (int, int) {
	virtualBoard := *vBoardPointer
	var xOffSet int

	switch options.XAlign {
	case "left":
		xOffSet = 0
	case "center":
		xOffSet = (boardWidth - len(virtualBoard[0])) / 2
		fmt.Println(xOffSet)
	case "right":
		xOffSet = boardWidth - len(virtualBoard[0])
	default:
		xOffSet, _ = strconv.Atoi(options.XAlign)
	}

	var yOffSet int

	switch options.YAlign {
	case "top":
		// we don't do anything
	case "center":
		yOffSet = (boardHeight - len(virtualBoard)) / 2
	case "bottom":
		yOffSet = boardHeight - len(virtualBoard)
	default:
		yOffSet, _ = strconv.Atoi(options.YAlign)
	}

	return xOffSet, yOffSet
}

func (b *Flipboard) getNextCountdown() options.FlipboardMessageOptions {
	horizonEventTime, err := time.Parse("2006-01-02", b.countdownDate)
	if err != nil {
		fmt.Println(err)
	}
	t := time.Now()
	elapsed := horizonEventTime.Sub(t)
	days := int(elapsed.Hours() / 24)
	hours := int(elapsed.Hours()) % 24
	mins := int(elapsed.Minutes()) % 60
	secs := int(elapsed.Seconds()) % 60
	msg := options.FlipboardMessageOptions{
		Message:     fmt.Sprintf("HORIZON EVENT\n%d:%02d:%02d:%02d", days, hours, mins, secs),
		DisplayTime: 0, // the CountdownDate option controls the timing
		Align:       "center center",
	}
	return msg
}
