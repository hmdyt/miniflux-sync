package diff

import "slices"

// State represents either the current remote state of Miniflux, or the intended local state of
// Miniflux.
type State struct {
	FeedURLsByCategoryTitle map[string][]string
	FeedsByCategoryTitle    map[string][]Feed
}

// CategoryExists checks if a category exists in the state.
func (s State) CategoryExists(categoryTitle string) bool {
	_, exists := s.FeedURLsByCategoryTitle[categoryTitle]
	return exists
}

// CategoryTitles returns a list of all category titles in the state.
func (s State) CategoryTitles() []string {
	categorySet := map[string]struct{}{}

	for categoryTitle := range s.FeedURLsByCategoryTitle {
		categorySet[categoryTitle] = struct{}{}
	}

	categoryTitles := make([]string, 0, len(categorySet))
	for categoryTitle := range categorySet {
		categoryTitles = append(categoryTitles, categoryTitle)
	}

	return categoryTitles
}

// FeedExists checks if a feed URL exists in a specific category.
func (s State) FeedExists(feedURL string, categoryTitle string) bool {
	feedURLs, exists := s.FeedURLsByCategoryTitle[categoryTitle]
	if !exists {
		return false
	}

	return slices.Contains(feedURLs, feedURL)
}

// FeedURLs returns a list of all feed URLs in the state.
func (s State) FeedURLs() []string {
	feedURLs := []string{}

	for _, urls := range s.FeedURLsByCategoryTitle {
		feedURLs = append(feedURLs, urls...)
	}

	return feedURLs
}

// GetFeedOptions returns the options for a specific feed URL, or empty options if not found.
func (s State) GetFeedOptions(feedURL string) FeedOptions {
	for _, feeds := range s.FeedsByCategoryTitle {
		for _, feed := range feeds {
			if feed.URL == feedURL {
				return feed.Options
			}
		}
	}
	return FeedOptions{}
}

// GetFeedsByCategory returns the feeds by category.
// If FeedsByCategoryTitle is not set, it falls back to FeedURLsByCategoryTitle.
func (s State) GetFeedsByCategory() map[string][]Feed {
	if len(s.FeedsByCategoryTitle) > 0 {
		return s.FeedsByCategoryTitle
	}

	// Fallback to FeedURLsByCategoryTitle for backward compatibility.
	result := make(map[string][]Feed)
	for categoryTitle, feedURLs := range s.FeedURLsByCategoryTitle {
		feeds := make([]Feed, 0, len(feedURLs))
		for _, url := range feedURLs {
			feeds = append(feeds, Feed{URL: url})
		}
		result[categoryTitle] = feeds
	}
	return result
}
