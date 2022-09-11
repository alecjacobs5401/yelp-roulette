package yelp

import (
	"context"

	"github.com/shurcooL/graphql"
)

// ClientConfig provides configuration for a new Yelp Client
type ClientConfig struct {
	Context     context.Context
	AccessToken string
}

// Client provides a wrapped to the Yelp GraphQL API
type Client struct {
	context context.Context
	client  *graphql.Client
}

// SearchRequest represents the parameters the Yelp Search request
type SearchRequest struct {
	Term          string
	Location      string
	OpenNow       bool
	Price         []string
	MaxSampleSize int
}

type searchQuery struct {
	Search struct {
		Total    graphql.Int
		Business []YelpBusiness
	} `graphql:"search(term: $term, location: $location, limit: $limit, offset: $offset, price: $price, open_now: $openNow)"`
}

// YelpBusiness represents a Yelp Business returned from the Search query
type YelpBusiness struct {
	ID    graphql.ID
	Name  graphql.String
	URL   graphql.String
	Price graphql.String
}
