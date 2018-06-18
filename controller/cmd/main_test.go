package main

import (
	"reflect"
	"testing"

	"github.com/armory/flipdisks/controller/pkg/fontmap"
)

func TestCreateVirtualBoard(t *testing.T) {
	tests := []struct {
		testDescription                string
		panelWidth, numberOfPanelsWide int
		message                        string

		expect VirtualBoard
	}{
		{
			testDescription:    "It should print out a simple message",
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
			testDescription:    "It should print out a simple message with new line",
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
	}

	for index, testCase := range tests {
		msgAsDots := fontmap.Render(testCase.message)
		got := createVirtualBoard(testCase.panelWidth, testCase.numberOfPanelsWide, msgAsDots, testCase.message)
		if !reflect.DeepEqual(testCase.expect, got) {
			t.Errorf("Test %d", index)
			t.Errorf("Expected\n%#v:\n%s", testCase.expect, testCase.expect)
			t.Errorf("Got\n%#v:\n%s", got, got)
		}
	}
}
