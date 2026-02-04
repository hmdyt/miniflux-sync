package api

import (
	"github.com/pkg/errors"
	"github.com/revett/miniflux-sync/diff"
	miniflux "miniflux.app/v2/client"
)

// GenerateDiffState generates a diff.State struct from a list of feeds.
func GenerateDiffState(
	feeds []*miniflux.Feed, categories []*miniflux.Category,
) (*diff.State, error) {
	state := diff.State{
		FeedURLsByCategoryTitle: map[string][]string{},
		FeedsByCategoryTitle:    map[string][]diff.Feed{},
	}

	// Initialise empty slices for each category.
	for _, category := range categories {
		state.FeedURLsByCategoryTitle[category.Title] = []string{}
		state.FeedsByCategoryTitle[category.Title] = []diff.Feed{}
	}

	// Populate state with values, and create category set.
	for _, feed := range feeds {
		if feed.Category == nil {
			return nil, errors.New("feed has no category")
		}
		categoryTitle := feed.Category.Title

		state.FeedURLsByCategoryTitle[categoryTitle] = append(
			state.FeedURLsByCategoryTitle[categoryTitle], feed.FeedURL,
		)

		options := extractFeedOptions(feed)
		state.FeedsByCategoryTitle[categoryTitle] = append(
			state.FeedsByCategoryTitle[categoryTitle], diff.Feed{
				URL:     feed.FeedURL,
				Options: options,
			},
		)
	}

	return &state, nil
}

// extractFeedOptions extracts the configurable options from a Miniflux feed.
func extractFeedOptions(feed *miniflux.Feed) diff.FeedOptions {
	return diff.FeedOptions{
		Crawler:                     boolPtr(feed.Crawler),
		Username:                    stringPtrIfNotEmpty(feed.Username),
		Password:                    stringPtrIfNotEmpty(feed.Password),
		UserAgent:                   stringPtrIfNotEmpty(feed.UserAgent),
		Cookie:                      stringPtrIfNotEmpty(feed.Cookie),
		Disabled:                    boolPtr(feed.Disabled),
		IgnoreHTTPCache:             boolPtr(feed.IgnoreHTTPCache),
		FetchViaProxy:               boolPtr(feed.FetchViaProxy),
		AllowSelfSignedCertificates: boolPtr(feed.AllowSelfSignedCertificates),
		DisableHTTP2:                boolPtr(feed.DisableHTTP2),
		ScraperRules:                stringPtrIfNotEmpty(feed.ScraperRules),
		RewriteRules:                stringPtrIfNotEmpty(feed.RewriteRules),
		BlocklistRules:              stringPtrIfNotEmpty(feed.BlocklistRules),
		KeeplistRules:               stringPtrIfNotEmpty(feed.KeeplistRules),
		HideGlobally:                boolPtr(feed.HideGlobally),
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func stringPtrIfNotEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
