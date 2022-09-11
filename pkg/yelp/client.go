package yelp

import (
	"context"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type ClientConfig struct {
	Context     context.Context
	AccessToken string
}

type Client struct {
	client *graphql.Client
}

const endpoint = "https://api.yelp.com/v3/graphql"

// NewClient returns a Yelp GraphQL Client Wrapper
func NewClient(config ClientConfig) *Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.AccessToken},
	)
	httpClient := oauth2.NewClient(config.Context, src)
	client := graphql.NewClient(endpoint, httpClient)

	return &Client{client: client}
}
