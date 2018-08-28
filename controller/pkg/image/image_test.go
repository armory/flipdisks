package image

import (
	"image"
	"os"
	"reflect"
	"testing"

	"github.com/armory/flipdisks/controller/pkg/fontmap"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	armory := []fontmap.Row{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
		{0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 0},
		{0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	f, err := os.Open("armory.jpg")
	if err != nil {
		t.Error(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		t.Error(err)
	}
	bounds := img.Bounds()

	v := convert(img, bounds, false, 140)
	if !reflect.DeepEqual(armory, v) {
		t.Error("images are not equal")
		t.Errorf("Got %v", v)
		t.Errorf("Expected %v", armory)
	}
}

func TestIsPlainImageUrl(t *testing.T) {
	tests := map[string]struct {
		url string

		expected bool
	}{
		".png": {
			url:      "http....",
			expected: true,
		},
		".jpg": {
			url:      "https://blah.jpg",
			expected: true,
		},
		".jpeg": {
			url:      "https://blah.jpeg",
			expected: true,
		},
		".txt": {
			url:      "http://blah.txt",
			expected: true,
		},
		"no extension": {
			url:      "http://google.com",
			expected: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.True(t, IsPlainImageUrl(test.url) == test.expected, "failed message")

		})
	}
}

func TestIsGifUrl(t *testing.T) {
	tests := map[string]struct {
		url string

		expected bool
	}{
		".gif": {
			url:      "http....",
			expected: true,
		},
		".png": {
			url:      "http....",
			expected: false,
		},
		".jpg": {
			url:      "https://blah.jpg",
			expected: false,
		},
		".jpeg": {
			url:      "https://blah.jpeg",
			expected: false,
		},
		".txt": {
			url:      "http....",
			expected: false,
		},
		"no extension": {
			url:      "http....",
			expected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.True(t, IsGifUrl(test.url) == test.expected, "failed message")

		})
	}
}
