package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"

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
	buffer := bytes.Buffer{}
	xml.EscapeText(&buffer, []byte(m))
	fmt.Fprintf(w, responseTemplate, buffer.String())
}

const helpText = `Welcome to Yelp Roulette!

Randomly select a restaurant/business in a specific location that match provided criteria.
Usage:
Provide a seed search term to select a business. Use key phrases to filter businesses by location, price, etc.
Key Phrases:
- "in": Used to provide a location (e.g. Santa Barbara, CA)
- "within": Provide a search radius with optional units. Default is in meters. (e.g. 10mi, 15 km)
- "$", "$$", "$$$", "$$$$": price levels of businesses to include. If none are provided, all price levels are included.
- "open": restrict search to businesses that are currently open.
Examples:
- "dinner in Santa Barbara, CA"
- "breakfast in NYC open $ $$"
- "italian food in Austin, TX within 10 miles"
`

func main() {
	log.SetLevel(log.DebugLevel)
	log.Info("yelp-roulette server started...")

	http.HandleFunc("/help", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("help()")
		fmt.Fprint(w, helpText)
	})

	http.HandleFunc("/sms", func(w http.ResponseWriter, r *http.Request) {
		client := yelp.NewClient(yelp.ClientConfig{
			Context:     r.Context(),
			AccessToken: os.Getenv("YELP_ROULETTE_ACCESS_TOKEN"),
		})

		body := r.FormValue("Body")

		searchRequest, err := yelp.ParseRequest(body)
		if err != nil {
			respond(w, fmt.Sprintf("ERROR: %v\n\n%s", err, helpText))
			return
		}
		log.Debugf("%#v\n", searchRequest)

		if searchRequest.Location == "" {
			respond(w, helpText)
		} else {
			business, err := client.RandomBusiness(searchRequest)
			if err != nil {
				log.Errorf("querying for random business: %v", err)
				respond(w, "Sorry, there was an issue processing your request. Please try again later.")
			} else {
				respond(w, fmt.Sprintf("%s - %s\n%s", business.Name, business.Price, business.URL))
			}
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
