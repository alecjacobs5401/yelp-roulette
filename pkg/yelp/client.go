package yelp

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

const endpoint = "https://api.yelp.com/v3/graphql"

// NewClient returns a Yelp GraphQL Client Wrapper
func NewClient(config ClientConfig) *Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.AccessToken},
	)
	httpClient := oauth2.NewClient(config.Context, src)
	client := graphql.NewClient(endpoint, httpClient)

	return &Client{context: config.Context, client: client}
}

// init sets initial values for variables used in the function.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Search returns all businesses matching the SearchRequest terms up to a provided maximum sample size
func (c *Client) Search(params SearchRequest) ([]YelpBusiness, error) {
	variables := map[string]interface{}{
		"term":     graphql.String(params.Term),
		"location": graphql.String(params.Location),
		"limit":    graphql.Int(50),
		"offset":   graphql.Int(0),
		"openNow":  graphql.Boolean(params.OpenNow),
		"price":    graphql.String(strings.Join(params.Price, ", ")),
	}
	query := searchQuery{}
	if err := c.client.Query(c.context, &query, variables); err != nil {
		return []YelpBusiness{}, fmt.Errorf("executing query: %w", err)
	}

	return query.Search.Business, nil
}

func (c *Client) RandomBusiness(params SearchRequest) (YelpBusiness, error) {
	businesses, err := c.Search(params)
	if err != nil {
		return YelpBusiness{}, err
	}

	return businesses[rand.Intn(len(businesses))], nil
}
