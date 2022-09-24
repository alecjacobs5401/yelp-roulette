package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/alecjacobs5401/yelp-roulette/pkg/yelp"
	log "github.com/sirupsen/logrus"
)

const responseTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<Response>
	<Message>
		<Body>
			%s
		</Body>
	</Message>
</Response>
`

func respond(w http.ResponseWriter, m string) {
	log.Debugf("respond(): %q", m)
	fmt.Fprintf(w, responseTemplate, m)
}

const helpText = `Welcome to Yelp Roulette!

Randomly select a restaurant/business in a specific location that match provided criteria.
Usage:
Provide a seed search term to select a business. Use key phrases to filter businesses by location, price, etc.
Key Phrases:
- "in": Used to provide a location (e.g. Santa Barbara, CA)
- "$", "$$", "$$$", "$$$$": price levels of businesses to include. If none are provided, all price levels are included.
- "open": restrict search to businesses that are currently open.
Examples:
- "dinner in Santa Barbara, CA"
- "breakfast in NYC open $ $$"
`

func main() {
	log.SetLevel(log.DebugLevel)
	log.Info("yelp-roulette server started...")

	client := yelp.NewClient(yelp.ClientConfig{
		Context:     context.Background(),
		AccessToken: os.Getenv("YELP_ROULETTE_ACCESS_TOKEN"),
	})

	http.HandleFunc("/sms", func(w http.ResponseWriter, r *http.Request) {
		body := r.FormValue("Body")
		searchRequest, _ := parseBody(body)
		fmt.Printf("%#v\n", searchRequest)

		business, err := client.RandomBusiness(searchRequest)
		if err != nil {
			respond(w, "Sorry, there was an issue processing your request. Please try again later.")
		} else {
			respond(w, fmt.Sprintf("%s - %s\n%s", business.Name, business.Price, business.URL))
		}
	})

	if err := http.ListenAndServe(getPort(), nil); err != nil {
		log.Fatal(err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	return fmt.Sprintf(":%s", port)
}

// takes a message and parses out the text into a SearchRequest
func parseBody(body string) (yelp.SearchRequest, error) {
	tokens := strings.Split(body, " ")

	inIndex := -1
	openIndex := -1
	priceLevels := []string{}
	locationTokens := []string{}
	termTokens := []string{}
	for index, token := range tokens {
		switch strings.ToLower(token) {
		case "in":
			inIndex = index
		case "open":
			openIndex = index
		case "$", "$$", "$$$", "$$$$":
			priceLevels = append(priceLevels, fmt.Sprint(strings.Count(token, "$")))
		default:
			if inIndex >= 0 && index > inIndex {
				locationTokens = append(locationTokens, token)
			} else {
				termTokens = append(termTokens, token)
			}
		}
	}

	return yelp.SearchRequest{
		Term:     strings.Join(termTokens, " "),
		Location: strings.Join(locationTokens, " "),
		OpenNow:  openIndex >= 0,
		Price:    priceLevels,
	}, nil
}
