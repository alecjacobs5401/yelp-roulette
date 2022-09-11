package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecjacobs5401/yelp-roulette/pkg/yelp"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yelp-roulette [options] [query]",
	Short: "Randomly select a restaurant from the Yelp API",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if accessToken == "" {
			accessToken = os.Getenv("YELP_ROULETTE_ACCESS_TOKEN")
		}

		yelp.NewClient(yelp.ClientConfig{Context: context.Background(), AccessToken: accessToken})
	},
}

var (
	accessToken string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&accessToken, "access-token", "t", "", "Yelp Developer API Access Token")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
