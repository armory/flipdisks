package main

import (
	"reflect"
	"testing"

	"github.com/armory/flipdisks/controller/pkg/flipboard"
	"github.com/armory/flipdisks/controller/pkg/fontmap"
	"github.com/armory/flipdisks/controller/pkg/options"
	"github.com/armory/flipdisks/controller/pkg/virtualboard"
	"github.com/kr/pty"
)

func TestCreateVirtualBoard(t *testing.T) {
	tests := []struct {
		testDescription                string
		panelWidth, numberOfPanelsWide int
		message                        string

		expect virtualboard.VirtualBoard
	}{
		{
			testDescription:    "It should print out a simple Message",
			panelWidth:         7,
			numberOfPanelsWide: 2,
			message:            "ab",

			expect: []fontmap.Row{
				{0, 0, 0, 0, 1, 0, 0, 0},
				{0, 1, 1, 0, 1, 1, 0, 0},
				{1, 0, 1, 0, 1, 0, 1, 0},
				{1, 0, 1, 0, 1, 0, 1, 0},
				{0, 1, 1, 0, 1, 1, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			},
		},

		{
			testDescription:    "It should print out a simple Message with new line",
			panelWidth:         7,
			numberOfPanelsWide: 2,
			message:            "a\nb",

			expect: []fontmap.Row{
				{0, 0, 0, 0},
				{0, 1, 1, 0},
				{1, 0, 1, 0},
				{1, 0, 1, 0},
				{0, 1, 1, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{1, 0, 0, 0},
				{1, 1, 0, 0},
				{1, 0, 1, 0},
				{1, 0, 1, 0},
				{1, 1, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
		},

		{
			testDescription:    "It should word break on board overflow",
			panelWidth:         4, // size of a
			numberOfPanelsWide: 4,
			message:            "aa bbb ccc",

			expect: []fontmap.Row{
				// aa
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0, 1, 1, 0},
				{1, 0, 1, 0, 1, 0, 1, 0},
				{1, 0, 1, 0, 1, 0, 1, 0},
				{0, 1, 1, 0, 0, 1, 1, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				// bbb
				{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0},
				{1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0},
				{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0},
				{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0},
				{1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				// ccc
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0},
				{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0},
				{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0},
				{0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},

		{
			testDescription:    "It should still character break if the word is long and there's no spaces",
			panelWidth:         4, // size of "a"
			numberOfPanelsWide: 4, // number of characters that can fit on a line
			message:            "aaaaaa",

			expect: []fontmap.Row{
				// aaa
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0},
				{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0},
				{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0},
				{0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				// aa
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0, 1, 1, 0},
				{1, 0, 1, 0, 1, 0, 1, 0},
				{1, 0, 1, 0, 1, 0, 1, 0},
				{0, 1, 1, 0, 0, 1, 1, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			},
		},

		{
			testDescription:    "It should print out a special character",
			panelWidth:         4,
			numberOfPanelsWide: 2,
			message:            "$",

			expect: []fontmap.Row{
				{0, 1, 1, 1, 0, 0},
				{1, 0, 1, 0, 0, 0},
				{0, 1, 1, 1, 0, 0},
				{0, 0, 1, 0, 1, 0},
				{1, 1, 1, 1, 0, 0},
				{0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0},
			},
		},
		{
			testDescription:    "It should handle spacing in between special characters",
			panelWidth:         4, // size of "?"
			numberOfPanelsWide: 4, // number of characters that can fit on a line
			message:            "? ?",

			expect: []fontmap.Row{
				//? ?
				{1, 1, 0, 0, 0, 0, 1, 1, 0, 0},
				{0, 0, 1, 0, 0, 0, 0, 0, 1, 0},
				{0, 1, 0, 0, 0, 0, 0, 1, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 1, 0, 0, 0, 0, 0, 1, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		{
			testDescription:    "It should only count the dotChar width if the dotChar exists",
			panelWidth:         4, // size of "?"
			numberOfPanelsWide: 4, // number of characters that can fit on a line
			message:            "ū u",

			expect: []fontmap.Row{
				//█ u
				{1, 1, 1, 0, 0, 0, 0, 0, 0},
				{1, 1, 1, 0, 0, 1, 0, 1, 0},
				{1, 1, 1, 0, 0, 1, 0, 1, 0},
				{1, 1, 1, 0, 0, 1, 0, 1, 0},
				{1, 1, 1, 0, 0, 1, 1, 1, 0},
				{1, 1, 1, 0, 0, 0, 0, 0, 0},
				{1, 1, 1, 0, 0, 0, 0, 0, 0},
			},
		},
	}

	for index, testCase := range tests {
		msgAsDots := fontmap.Render(testCase.message)
		got := flipboard.CreateVirtualBoard(testCase.panelWidth, testCase.numberOfPanelsWide, msgAsDots, testCase.message)
		if !reflect.DeepEqual(testCase.expect, got) {
			t.Errorf("Test %d", index)
			t.Errorf("Expected\n%#v:\n%s", testCase.expect, testCase.expect)
			t.Errorf("Got\n%#v:\n%s", got, got)
		}
	}
}

// These tests are only concerned with not crashing the flipboard when displaying a message
// Todo: we should test the actual virtual board. there's a few options:
// 	- check the cache
// 	- check each panel's value
//	- add a return type and check that
func TestDisplayMessageToPanels(t *testing.T) {
	tests := map[string]struct {
		msg options.FlipboardMessageOptions
	}{
		"simple and autofill": {
			msg: options.FlipboardMessageOptions{
				Message: "Simple String",
			},
		},
		"string with linebreak": {
			msg: options.FlipboardMessageOptions{
				Message: "Simple\nString",
			},
		},
		"image url": {
			msg: func() options.FlipboardMessageOptions {
				o := options.GetDefaultOptions()
				o.Message = "https://cl.ly/2r0k2I1P0d2i/Armory_logo_monochrome_shield%20(2).jpg"
				o.DisplayTime = 1
				return o
			}(),
		},
		"simple string inverted and autofill": {
			msg: options.FlipboardMessageOptions{
				Message:  "Simple String",
				Inverted: true,
			},
		},
		"simple string align center center": {
			msg: options.FlipboardMessageOptions{
				Message: "Simple String",
				Align:   "center center",
			},
		},
		"simple string align left center": {
			msg: options.FlipboardMessageOptions{
				Message: "Simple String",
				Align:   "left center",
			},
		},
		"simple string align right center": {
			msg: options.FlipboardMessageOptions{
				Message: "Simple String",
				Align:   "right center",
			},
		},
		"simple string align right bottom": {
			msg: options.FlipboardMessageOptions{
				Message: "Simple String",
				Align:   "right right",
			},
		},
		"simple string fill": {
			msg: options.FlipboardMessageOptions{
				Message: "Simple String",
				Fill:    "false",
			},
		},
		"simple string autofill": {
			msg: options.FlipboardMessageOptions{
				Message: "Simple String",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// setup a psudo tty to use
			_, tty, _ := pty.Open()
			defer func() { tty.Close() }()

			panelInfo := flipboard.PanelInfo{
				Baud:                     9600,
				Port:                     tty.Name(),
				PanelWidth:               28,
				PanelHeight:              7,
				PhysicallyDisplayedWidth: 7,
			}

			panelLayout := [][]flipboard.PanelAddress{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			}

			board, _ := flipboard.NewFlipboard(panelInfo, panelLayout)
			flipboard.DisplayMessageToPanels(board, &test.msg)
		})
	}
}
