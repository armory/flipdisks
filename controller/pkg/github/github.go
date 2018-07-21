package github

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Client struct {
	Emojis  EmojiLookup
	context context.Context
	client  *github.Client
}

type EmojiLookup map[string]string

func New(options ...func(*Client) error) (*Client, error) {
	g := Client{}

	for _, option := range options {
		err := option(&g)
		if err != nil {
			return nil, err
		}
	}

	// I guess we want an unauthenticated github call
	if g.client == nil {
		g.client = github.NewClient(nil)
		g.context = context.Background()
	}

	return &g, nil
}

func Token(token string) func(*Client) error {
	return func(g *Client) error {
		g.context = context.Background()
		tokenSource := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)

		tc := oauth2.NewClient(g.context, tokenSource)

		g.client = github.NewClient(tc)

		return nil
	}
}

// GetEmojis will return github list of emojis, it'll cache the list that we
// got from github
func (g *Client) GetEmojis() (EmojiLookup, error) {
	if g.Emojis == nil {
		githubEmojiLookup, _, err := g.client.ListEmojis(g.context)
		g.Emojis = githubEmojiLookup
		if err != nil {
			return nil, err
		}
	}

	return g.Emojis, nil
}
