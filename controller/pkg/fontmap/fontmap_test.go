package fontmap

import (
	"reflect"
	"testing"
)

func TestGenerateSpace(t *testing.T) {
	tests := []struct {
		width, height, fill int
		expected            Letter
	}{
		{1, 1, 0,
			Letter{
				Row{0},
			},
		},
		{1, 1, 1,
			Letter{
				Row{1},
			},
		},

		// example of what a actual whitespace looks like
		{2, 6, 0,
			Letter{
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
			},
		},

		// example of what an unknown character looks like
		{4, 6, 1,
			Letter{
				Row{1, 1, 1, 1},
				Row{1, 1, 1, 1},
				Row{1, 1, 1, 1},
				Row{1, 1, 1, 1},
				Row{1, 1, 1, 1},
				Row{1, 1, 1, 1},
			},
		},
	}

	for _, testCase := range tests {
		got := GenerateSpace(testCase.width, testCase.height, testCase.fill)
		if !reflect.DeepEqual(testCase.expected, got) {
			t.Errorf("Expected %#v, but got %#v", testCase.expected, got)
			t.Errorf("Expected : \n%s", testCase.expected)
			t.Errorf("Got: \n%s", got)
		}
	}
}
