package slackbot

import (
	"testing"

	"github.com/armory/flipdisks/controller/pkg/options"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func TestCleanupSlackEncodedCharacters(t *testing.T) {
	tests := map[string]struct {
		msg string

		Expected string
	}{
		"< symbol": {
			msg:      "5 &lt; 13",
			Expected: "5 < 13",
		},
		"> symbol": {
			msg:      "10 &gt; 1",
			Expected: "10 > 1",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.True(t, cleanupSlackEncodedCharacters(test.msg) == test.Expected, "Failed!")
		})
	}
}

func TestSplitMessageAndOptions(t *testing.T) {
	tests := map[string]struct {
		msg string

		Expected []options.FlipBoardDisplayOptions
	}{
		//"simple message": {
		//	msg: "hello",
		//	Expected: func() []options.FlipBoardDisplayOptions {
		//		o := options.GetDefaultOptions()
		//		o.Message = "hello"
		//		return []options.FlipBoardDisplayOptions{o}
		//	}(),
		//},
		"simple options": {
			msg: "hello\n---\ninverted: true",
			Expected: func() []options.FlipBoardDisplayOptions {
				o := options.GetDefaultOptions()
				o.Message = "hello"
				o.Inverted = true
				return []options.FlipBoardDisplayOptions{o}
			}(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			diff := deep.Equal(splitMessageAndOptions(test.msg), test.Expected)
			if diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestUnmarshleOptions(t *testing.T) {
	tests := map[string]struct {
		raw string

		Expected options.FlipBoardDisplayOptions
	}{
		"simple": {
			raw: "inverted: true",
			Expected: func() options.FlipBoardDisplayOptions {
				o := options.GetDefaultOptions()
				o.Inverted = true
				return o
			}(),
		},
		"align": {
			raw: "align: center center",
			Expected: func() options.FlipBoardDisplayOptions {
				o := options.GetDefaultOptions()

				// these shouldn't be parsed and set in this function
				o.Align = "center center"
				o.XAlign = ""
				o.YAlign = ""
				return o
			}(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			diff := deep.Equal(unmarshleOptions(test.raw), test.Expected)
			if diff != nil {
				t.Error(diff)
			}
		})
	}
}
