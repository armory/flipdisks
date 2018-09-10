package image

import (
	"image"
	"io/ioutil"
	"os"
	"testing"

	"github.com/armory/flipdisks/controller/pkg/virtualboard"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	f, err := os.Open("test_fixtures/armory.jpg")
	if err != nil {
		t.Error(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		t.Error(err)
	}
	bounds := img.Bounds()

	v := convertImgToVirtualBoard(img, bounds, false, 140)

	expected, err := ioutil.ReadFile("test_fixtures/armory_virtualboard.txt")
	if err != nil {
		t.Error(err)
	}

	if (*v).String() != string(expected) {
		t.Error("images are not equal")
		t.Errorf("Expected\n%s", expected)
		t.Errorf("Got\n%s", *v)
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

		gifUrls []string
	}{
		".gif": {
			url:     "please display this: http://www.blah.com/doge.gif",
			gifUrls: []string{"http://www.blah.com/doge.gif"},
		},
		"https .gif": {
			url:     "https://www.blah.com/cats.gif",
			gifUrls: []string{"https://www.blah.com/cats.gif"},
		},
		"gif in url path should not match, only if it's a file extension": {
			url:     "https://www.blah.com/cats/gif",
			gifUrls: nil,
		},
		"should be able to handle anchors": {
			url:     "https://www.blah.com/cats.gif#blah",
			gifUrls: []string{"https://www.blah.com/cats.gif#blah"},
		},
		"should be able to handle query params": {
			url:     "https://www.blah.com/cats.gif?one=1",
			gifUrls: []string{"https://www.blah.com/cats.gif?one=1"},
		},
		".png": {
			url:     "http://www.blah.com/ballon.png",
			gifUrls: nil,
		},
		".jpg": {
			url:     "https://www.blah.com/doge.jpg",
			gifUrls: nil,
		},
		".jpeg": {
			url:     "https://www.blah.com/doge.jpeg",
			gifUrls: nil,
		},
		".txt": {
			url:     "https://www.blah.com/names.txt",
			gifUrls: nil,
		},
		"no extension": {
			url:     "http://www.blah.com",
			gifUrls: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gifUrls := GetGifUrl(test.url)
			diffs := deep.Equal(gifUrls, test.gifUrls)

			for _, diff := range diffs {
				t.Errorf(`Test "%s" failed with: %s`, name, diff)
			}
		})
	}
}

func TestConvertGifFromURLToVirtualBoard(t *testing.T) {
	gifBytes, err := ioutil.ReadFile("test_fixtures/fast_parrot.gif")
	if err != nil {
		t.Error(err)
	}

	txtBytes, err := ioutil.ReadFile("test_fixtures/fast_parrot_virtualboard.txt")
	if err != nil {
		t.Error(err)
	}

	gotGif, err := convertGifToVirtualBoard(gifBytes, 50, 50, false, 90)
	if err != nil {
		t.Error(err)
	}

	gotGifFramesTxt := ""
	for _, frame := range gotGif.Flipboards {
		blah := virtualboard.VirtualBoard(*frame)
		gotGifFramesTxt += blah.String() + "\n"
	}

	// Uncomment this line to write gif to txt
	// ioutil.WriteFile("test_fixtures/fast_parrot_virtualboard.txt", []byte(gotGifFramesTxt), os.ModePerm)

	if gotGifFramesTxt != string(txtBytes) {
		t.Error("gif are not equal")
		t.Errorf("Expected\n%s", txtBytes)
		t.Errorf("Got\n%s", gotGifFramesTxt)
	}
}
