package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecjacobs5401/yelp-roulette/pkg/convert"
	"github.com/alecjacobs5401/yelp-roulette/pkg/yelp"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yelp-roulette [query]",
	Short: "Randomly select a restaurant from the Yelp API",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if accessToken == "" {
			accessToken = os.Getenv("YELP_ROULETTE_ACCESS_TOKEN")
		}

		searchWithin, err := convert.ToMeters(within, withinUnit)
		if err != nil {
			fatalError(err)
		}

		client := yelp.NewClient(yelp.ClientConfig{Context: context.Background(), AccessToken: accessToken})
		business, err := client.RandomBusiness(yelp.SearchRequest{
			Term:          args[0],
			Location:      location,
			OpenNow:       openNow,
			Price:         price,
			MaxSampleSize: maxSampleSize,
			Within:        int(searchWithin),
		})
		if err != nil {
			fatalError(err)
		}
		fmt.Printf("%s - %s\n%s\n", business.Name, business.Price, business.URL)
	},
}

var (
	accessToken   string
	location      string
	openNow       bool
	price         []string
	maxSampleSize int
	within        float64
	withinUnit    string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&accessToken, "access-token", "t", "", "Yelp Developer API Access Token")
	rootCmd.Flags().StringVarP(&location, "location", "l", "Santa Barbara, CA", "Location to base search results off of")
	rootCmd.Flags().BoolVarP(&openNow, "open-now", "", false, "Filters results based on if business is open now")
	rootCmd.Flags().StringArrayVarP(&price, "price", "p", []string{}, "Pricing levels to filter the search result with: 1 = $, 2 = $$, 3 = $$$, 4 = $$$$")
	rootCmd.Flags().IntVarP(&maxSampleSize, "max-sample-size", "m", 50, "Maximum sample size for random business selection")
	rootCmd.Flags().Float64VarP(&within, "within", "w", 20, "Limit search results to given radius")
	rootCmd.Flags().StringVarP(&withinUnit, "within-unit", "U", "miles", "Units for the within option. Allowed units: miles, meters, kilometers")

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fatalError(err)
	}
}

func fatalError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
