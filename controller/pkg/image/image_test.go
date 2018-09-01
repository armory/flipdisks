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

	v := convertImgToVirtualBoard(img, bounds, false, 140)
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
			url:      "http://www.blah.com/doge.png",
			expected: true,
		},
		"https .png": {
			url:      "https://www.blah.com/cats.png",
			expected: true,
		},
		"png in url path should not match, only if it's a file extension": {
			url:      "https://www.blah.com/cats/png",
			expected: false,
		},
		"should be able to handle anchors for .png": {
			url:      "https://www.blah.com/cats.png#blah",
			expected: true,
		},
		"should be able to handle query params for .png": {
			url:      "https://www.blah.com/cats.png?one=1",
			expected: true,
		},
		".jpg": {
			url:      "https://www.blah.com/doge.jpg",
			expected: true,
		},
		"https .jpg": {
			url:      "https://www.blah.com/cats.jpg",
			expected: true,
		},
		"jpg in url path should not match, only if it's a file extension": {
			url:      "https://www.blah.com/cats/jpg",
			expected: false,
		},
		"should be able to handle anchors for .jpg": {
			url:      "https://www.blah.com/cats.jpg#blah",
			expected: true,
		},
		"should be able to handle query params for .jpg": {
			url:      "https://www.blah.com/cats.jpg?one=1",
			expected: true,
		},
		".jpeg": {
			url:      "https://www.blah.com/doge.jpeg",
			expected: true,
		},
		"https .jpeg": {
			url:      "https://www.blah.com/cats.jpeg",
			expected: true,
		},
		"jpeg in url path should not match, only if it's a file extension": {
			url:      "https://www.blah.com/cats/jpeg",
			expected: false,
		},
		"should be able to handle anchors with .jpeg": {
			url:      "https://www.blah.com/cats.jpeg#blah",
			expected: true,
		},
		"should be able to handle query params with .jpeg": {
			url:      "https://www.blah.com/cats.jpeg?one=1",
			expected: true,
		},
		".txt": {
			url:      "https://www.blah.com/names.txt",
			expected: false,
		},
		"no extension": {
			url:      "http://www.blah.com",
			expected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.True(t, IsPlainImageUrl(test.url) == test.expected, "Failed")

		})
	}
}

func TestIsGifUrl(t *testing.T) {
	tests := map[string]struct {
		url string

		expected bool
	}{
		".gif": {
			url:      "http://www.blah.com/doge.gif",
			expected: true,
		},
		"https .gif": {
			url:      "https://www.blah.com/cats.gif",
			expected: true,
		},
		"gif in url path should not match, only if it's a file extension": {
			url:      "https://www.blah.com/cats/gif",
			expected: false,
		},
		"should be able to handle anchors": {
			url:      "https://www.blah.com/cats.gif#blah",
			expected: true,
		},
		"should be able to handle query params": {
			url:      "https://www.blah.com/cats.gif?one=1",
			expected: true,
		},
		".png": {
			url:      "http://www.blah.com/ballon.png",
			expected: false,
		},
		".jpg": {
			url:      "https://www.blah.com/doge.jpg",
			expected: false,
		},
		".jpeg": {
			url:      "https://www.blah.com/doge.jpeg",
			expected: false,
		},
		".txt": {
			url:      "https://www.blah.com/names.txt",
			expected: false,
		},
		"no extension": {
			url:      "http://www.blah.com",
			expected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.True(t, IsGifUrl(test.url) == test.expected, "Failed")

		})
	}
}
