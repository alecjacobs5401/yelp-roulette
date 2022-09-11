# yelp-roulette

Simple app to randomly select a restaurant using the Yelp API

## Usage

```
Randomly select a restaurant from the Yelp API

Usage:
  yelp-roulette [query] [flags]

Flags:
  -t, --access-token string   Yelp Developer API Access Token
  -h, --help                  help for yelp-roulette
  -l, --location string       Location to base search results off of (default "Santa Barbara, CA")
      --open-now              Filters results based on if business is open now
  -p, --price stringArray     Pricing levels to filter the search result with: 1 = $, 2 = $$, 3 = $$$, 4 = $$$$
```

The Access Token can also be configured by the environment variable `YELP_ROULETTE_ACCESS_TOKEN`

### Filter by Location

The Yelp API allows for fuzzy like searching of Businesses via geographic location terms (e.g. `Yosemite, CA`, `NYC`, `123 Main Street, Fake City, Fake State`).

The location can be filtered by providing the `--location` flag

```
yelp-roulette dinner --location "Albequerque, NM"
```

### Filter by Price Options

The Yelp API allows for limiting results based on price buckets.
Multiple price options can be provided and will be concatenated together. For example,

```
yelp-roulette dinner --price 1 --price 2
```

will limit results to those that match either the `$` or `$$` pricing level.

### Filter by Open Businesses

By default, all businesses that match the search terms will be in
sample pool for results. To restrict the sample pool to only those
that are currently open, provide the `--open-now` flag.

```
yelp-roulette dinner --open-now
```

## Development

Local development requires go1.19+. Follow these steps for local development:

1. Pull the latest `main` source
2. Create a new Feature Branch
3. Make changes to the source code
4. Run `make build` to build the `yelp-roulette` CLI
5. Test your changes with `bin/yelp-roulette`
