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

	for index, testCase := range tests {
		got := GenerateSpace(testCase.width, testCase.height, testCase.fill)
		if !reflect.DeepEqual(testCase.expected, got) {
			t.Errorf("Test %d: Expected %#v, but got %#v", index, testCase.expected, got)
			t.Errorf("Expected : \n%s", testCase.expected)
			t.Errorf("Got: \n%s", got)
		}
	}
}

func TestAddKerning(t *testing.T) {
	tests := []struct {
		// conditions
		letter          Letter
		amountOfKerning int

		// check against our result
		expected Letter
	}{
		{
			letter:          Letter{},
			amountOfKerning: 0,
			expected:        Letter{},
		},

		{
			letter:          Letter{},
			amountOfKerning: 1,
			expected:        Letter{},
		},


		{
			letter: Letter{
				Row{},
			},
			amountOfKerning: 0,
			expected: Letter{
				Row{},
			},
		},

		{
			letter: Letter{
				Row{},
			},
			amountOfKerning: 1,
			expected: Letter{
				Row{0},
			},
		},

		{
			letter: Letter{
				Row{},
			},
			amountOfKerning: 2,
			expected: Letter{
				Row{0, 0},
			},
		},

		{
			letter: Letter{
				Row{},
			},
			amountOfKerning: 3,
			expected: Letter{
				Row{0, 0, 0},
			},
		},

		{
			letter: Letter{
				Row{1},
			},
			amountOfKerning: 0,
			expected: Letter{
				Row{1},
			},
		},

		{
			letter: Letter{
				Row{1},
			},
			amountOfKerning: 1,
			expected: Letter{
				Row{1, 0},
			},
		},

		{
			letter: Letter{
				Row{1},
			},
			amountOfKerning: -1,
			expected: Letter{
				Row{1},
			},
		},


		{
			letter: Letter{
				{1, 0, 1},
				{0, 1, 0},
				{1, 0, 1},
			},
			amountOfKerning: 1,
			expected: Letter{
				{1, 0, 1, 0},
				{0, 1, 0, 0},
				{1, 0, 1, 0},
			},
		},
	}

	for index, testCase := range tests {
		got := AddKerning(testCase.letter, testCase.amountOfKerning)
		if !reflect.DeepEqual(testCase.expected, got) {
			t.Errorf("Test %d: Expected %#v, but got %#v", index, testCase.expected, got)
			t.Errorf("Expected : \n%s", testCase.expected)
			t.Errorf("Got: \n%s", got)
		}
	}
}

func TestRender(t *testing.T) {
	tests := []struct {
		message string
		expect  []Letter
	}{
		{message: "", expect: nil},
		{message: "a", expect: []Letter{{
			Row{0, 0, 0, 0},
			Row{0, 1, 1, 0},
			Row{1, 0, 1, 0},
			Row{1, 0, 1, 0},
			Row{0, 1, 1, 0},
			Row{0, 0, 0, 0},
			Row{0, 0, 0, 0},
		}}},

		{message: "a b", expect: []Letter{
			{
				Row{0, 0, 0, 0},
				Row{0, 1, 1, 0},
				Row{1, 0, 1, 0},
				Row{1, 0, 1, 0},
				Row{0, 1, 1, 0},
				Row{0, 0, 0, 0},
				Row{0, 0, 0, 0},
			},
			{
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
			},
			{
				Row{1, 0, 0, 0},
				Row{1, 1, 0, 0},
				Row{1, 0, 1, 0},
				Row{1, 0, 1, 0},
				Row{1, 1, 0, 0},
				Row{0, 0, 0, 0},
				Row{0, 0, 0, 0},
			},
		}},
		{message: "aa bb", expect: []Letter{
			{
				Row{0, 0, 0, 0},
				Row{0, 1, 1, 0},
				Row{1, 0, 1, 0},
				Row{1, 0, 1, 0},
				Row{0, 1, 1, 0},
				Row{0, 0, 0, 0},
				Row{0, 0, 0, 0},
			},
			{
				Row{0, 0, 0, 0},
				Row{0, 1, 1, 0},
				Row{1, 0, 1, 0},
				Row{1, 0, 1, 0},
				Row{0, 1, 1, 0},
				Row{0, 0, 0, 0},
				Row{0, 0, 0, 0},
			},
			{
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
				Row{0, 0},
			},
			{
				Row{1, 0, 0, 0},
				Row{1, 1, 0, 0},
				Row{1, 0, 1, 0},
				Row{1, 0, 1, 0},
				Row{1, 1, 0, 0},
				Row{0, 0, 0, 0},
				Row{0, 0, 0, 0},
			},
			{
				Row{1, 0, 0, 0},
				Row{1, 1, 0, 0},
				Row{1, 0, 1, 0},
				Row{1, 0, 1, 0},
				Row{1, 1, 0, 0},
				Row{0, 0, 0, 0},
				Row{0, 0, 0, 0},
			},
		}},
		{message: "a\nb", expect: []Letter{
			{
				Row{0, 0, 0, 0},
				Row{0, 1, 1, 0},
				Row{1, 0, 1, 0},
				Row{1, 0, 1, 0},
				Row{0, 1, 1, 0},
				Row{0, 0, 0, 0},
				Row{0, 0, 0, 0},
			},
			nil,
			{
				Row{1, 0, 0, 0},
				Row{1, 1, 0, 0},
				Row{1, 0, 1, 0},
				Row{1, 0, 1, 0},
				Row{1, 1, 0, 0},
				Row{0, 0, 0, 0},
				Row{0, 0, 0, 0},
			},
		}},
		{message: "Á", expect: []Letter{
			{
				Row{1, 1, 1},
				Row{1, 1, 1},
				Row{1, 1, 1},
				Row{1, 1, 1},
				Row{1, 1, 1},
				Row{1, 1, 1},
				Row{1, 1, 1},
			},
		}},
		{message: "🔥", expect: []Letter{
			{
				Row{1, 1, 1},
				Row{1, 1, 1},
				Row{1, 1, 1},
				Row{1, 1, 1},
				Row{1, 1, 1},
				Row{1, 1, 1},
				Row{1, 1, 1},
			},
		}},

	}
	for index, testCase := range tests {
		got := Render(testCase.message)
		if !reflect.DeepEqual(testCase.expect, got) {
			t.Errorf("Test %d", index)
			t.Errorf("Expected\n%#v:\n%s", testCase.expect, testCase.expect)
			t.Errorf("Got\n%#v:\n%s", got, got)
		}
	}
}
