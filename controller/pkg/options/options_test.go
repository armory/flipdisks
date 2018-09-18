package options

import (
	"testing"

	"github.com/go-test/deep"
)

func TestUnmarshleOptions(t *testing.T) {
	tests := map[string]struct {
		raw string

		Expected FlipboardMessageOptions
	}{
		"simple": {
			raw: "inverted: true",
			Expected: func() FlipboardMessageOptions {
				o := GetDefaultOptions()
				o.Inverted = true
				return o
			}(),
		},
		"align": {
			raw: "align: center center",
			Expected: func() FlipboardMessageOptions {
				o := GetDefaultOptions()

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
