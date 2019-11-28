package image

import (
	"image"
	"io/ioutil"
	"os"
	"testing"

	"flipdisks/pkg/virtualboard"
	"github.com/go-test/deep"
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

func TestGetPlainImageUrl(t *testing.T) {
	tests := map[string]struct {
		url string

		plainUrls []string
	}{
		".png": {
			url:       "http://www.blah.com/doge.png",
			plainUrls: []string{"http://www.blah.com/doge.png"},
		},
		"https .png": {
			url:       "https://www.blah.com/cats.png",
			plainUrls: []string{"https://www.blah.com/cats.png"},
		},
		"png in url path should not match, only if it's a file extension": {
			url:       "https://www.blah.com/cats/png",
			plainUrls: nil,
		},
		"should be able to handle anchors for .png": {
			url:       "https://www.blah.com/cats.png#blah",
			plainUrls: []string{"https://www.blah.com/cats.png#blah"},
		},
		"should be able to handle query params for .png": {
			url:       "https://www.blah.com/cats.png?one=1",
			plainUrls: []string{"https://www.blah.com/cats.png?one=1"},
		},
		".jpg": {
			url:       "https://www.blah.com/doge.jpg",
			plainUrls: []string{"https://www.blah.com/doge.jpg"},
		},
		"https .jpg": {
			url:       "https://www.blah.com/cats.jpg",
			plainUrls: []string{"https://www.blah.com/cats.jpg"},
		},
		"jpg in url path should not match, only if it's a file extension": {
			url:       "https://www.blah.com/cats/jpg",
			plainUrls: nil,
		},
		"should be able to handle anchors for .jpg": {
			url:       "https://www.blah.com/cats.jpg#blah",
			plainUrls: []string{"https://www.blah.com/cats.jpg#blah"},
		},
		"should be able to handle query params for .jpg": {
			url:       "https://www.blah.com/cats.jpg?one=1",
			plainUrls: []string{"https://www.blah.com/cats.jpg?one=1"},
		},
		".jpeg": {
			url:       "https://www.blah.com/doge.jpeg",
			plainUrls: []string{"https://www.blah.com/doge.jpeg"},
		},
		"https .jpeg": {
			url:       "https://www.blah.com/cats.jpeg",
			plainUrls: []string{"https://www.blah.com/cats.jpeg"},
		},
		"jpeg in url path should not match, only if it's a file extension": {
			url:       "https://www.blah.com/cats/jpeg",
			plainUrls: nil,
		},
		"should be able to handle anchors with .jpeg": {
			url:       "https://www.blah.com/cats.jpeg#blah",
			plainUrls: []string{"https://www.blah.com/cats.jpeg#blah"},
		},
		"should be able to handle query params with .jpeg": {
			url:       "https://www.blah.com/cats.jpeg?one=1",
			plainUrls: []string{"https://www.blah.com/cats.jpeg?one=1"},
		},
		".txt": {
			url:       "https://www.blah.com/names.txt",
			plainUrls: nil,
		},
		"no extension": {
			url:       "http://www.blah.com",
			plainUrls: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			plainUrls := GetPlainImageUrl(test.url)
			diffs := deep.Equal(plainUrls, test.plainUrls)

			for _, diff := range diffs {
				t.Errorf(`Test "%s" failed with: %s`, name, diff)
			}
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
		"slack: simple url": {
			url:     "<https://www.blah.com/cats.gif>",
			gifUrls: []string{"https://www.blah.com/cats.gif"},
		},
		"slack: should be able to handle query params": {
			url:     "<https://www.blah.com/cats.gif?one=1>",
			gifUrls: []string{"https://www.blah.com/cats.gif?one=1"},
		},
		"slack: should be able to handle query params with dashes": {
			url:     "<https://www.blah.com/cats.gif?one=1&hello=you-animal>",
			gifUrls: []string{"https://www.blah.com/cats.gif?one=1&hello=you-animal"},
		},
		"slack: simple url with anchor": {
			url:     "<https://www.blah.com/cats.gif#mooo>",
			gifUrls: []string{"https://www.blah.com/cats.gif#mooo"},
		},
		"slack: simple url with anchor dash and number": {
			url:     "<https://www.blah.com/cats.gif#mooo-12323-abc>",
			gifUrls: []string{"https://www.blah.com/cats.gif#mooo-12323-abc"},
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
