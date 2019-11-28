package slackbot

import (
	"testing"

	"flipdisks/pkg/options"
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

		Expected []options.FlipboardMessageOptions
	}{
		"simple message": {
			msg: "hello",
			Expected: func() []options.FlipboardMessageOptions {
				o := options.GetDefaultOptions()
				o.Message = "hello"
				return []options.FlipboardMessageOptions{o}
			}(),
		},
		"simple options": {
			msg: "hello\n---\ninverted: true",
			Expected: func() []options.FlipboardMessageOptions {
				o := options.GetDefaultOptions()
				o.Message = "hello"
				o.Inverted = true
				return []options.FlipboardMessageOptions{o}
			}(),
		},
		"playlist options": {
			msg: `---
- message: hello
  inverted: true
- message: world
  bwThreshold: 33
`,
			Expected: func() []options.FlipboardMessageOptions {
				h := options.GetDefaultOptions()
				h.Message = "hello"
				h.Inverted = true

				w := options.GetDefaultOptions()
				w.Message = "world"
				w.BWThreshold = 33
				return []options.FlipboardMessageOptions{h, w}
			}(),
		},
		"playlist with newline in front": {
			msg: `
---
- message: hello
  inverted: true
- message: world
  bwThreshold: 33
`,
			Expected: func() []options.FlipboardMessageOptions {
				h := options.GetDefaultOptions()
				h.Message = "hello"
				h.Inverted = true

				w := options.GetDefaultOptions()
				w.Message = "world"
				w.BWThreshold = 33
				return []options.FlipboardMessageOptions{h, w}
			}(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			diff := deep.Equal(options.SplitMessageAndOptions(test.msg), test.Expected)
			if diff != nil {
				t.Error(diff)
			}
		})
	}
}
