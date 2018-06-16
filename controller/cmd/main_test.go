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
			panelWidth:         4,         // size of a
			numberOfPanelsWide: 2 + 1 + 2, // character breaking creates "aa aa\na" but we want "aa\naaa"
			message:            "aa aaa",

			// "aa\naaa"
			expect: []fontmap.Row{
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0, 1, 1, 0},
				{1, 0, 1, 0, 1, 0, 1, 0},
				{1, 0, 1, 0, 1, 0, 1, 0},
				{0, 1, 1, 0, 0, 1, 1, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0},
				{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0},
				{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0},
				{0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
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
