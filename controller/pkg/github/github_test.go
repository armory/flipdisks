package github

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	t.Run("Providing a token", func(t *testing.T) {
		g, _ := New(Token("hello"))
		assert.True(t, g.context != nil, "Using the token, it should have set a context")
		assert.True(t, g.client != nil, "Using the token, it should have set a client")
	})

	t.Run("Not providing a token should create an anonymous client", func(t *testing.T) {
		g, _ := New()
		assert.True(t, g.context != nil, "Using the token, it should have set a context")
		assert.True(t, g.client != nil, "Using the token, it should have set a client")
	})
}

func TestClient_GetEmojis(t *testing.T) {
	t.Run("We should be able to return a list of cached emojis", func(t *testing.T) {
		g := Client{
			Emojis: EmojiLookup{"yay": "https://blah.com/yay.jpg"},
		}

		foundEmojis, _ := g.GetEmojis()
		assert.True(t, reflect.DeepEqual(foundEmojis, g.Emojis), "We didn't get the cached list of emojis")
	})

	t.Run("We should be able to go fetch a new list of emojis", func(t *testing.T) {
		// not sure how to test this one
	})
}
