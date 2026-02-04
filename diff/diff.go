package diff

import "sort"

// CalculateDiff calculates the differences between the local and remote state and returns the
// actions to be performed.
func CalculateDiff(local *State, remote *State) ([]Action, error) { //nolint:cyclop
	actions := []Action{}

	// Iterate over remote feeds and check if they exist in the local feeds.
	for categoryTitle, feedURLs := range remote.FeedURLsByCategoryTitle {
		for _, feedURL := range feedURLs {
			if !local.FeedExists(feedURL, categoryTitle) {
				actions = append(actions, Action{
					Type:          DeleteFeed,
					CategoryTitle: categoryTitle,
					FeedURL:       feedURL,
				})
			}
		}
	}

	// Iterate over remote categories and check if they exist in the local categories.
	for categoryTitle := range remote.FeedURLsByCategoryTitle {
		if !local.CategoryExists(categoryTitle) {
			actions = append(actions, Action{
				Type:          DeleteCategory,
				CategoryTitle: categoryTitle,
			})
		}
	}

	// Iterate over local categories and check if they exist in the remote categories.
	for categoryTitle := range local.FeedURLsByCategoryTitle {
		if !remote.CategoryExists(categoryTitle) {
			actions = append(actions, Action{
				Type:          CreateCategory,
				CategoryTitle: categoryTitle,
			})
		}
	}

	// Iterate over local feeds and check if they exist in the remote feeds.
	for categoryTitle, feeds := range local.GetFeedsByCategory() {
		for _, feed := range feeds {
			if !remote.FeedExists(feed.URL, categoryTitle) {
				actions = append(actions, Action{
					Type:          CreateFeed,
					CategoryTitle: categoryTitle,
					FeedURL:       feed.URL,
					FeedOptions:   feed.Options,
				})
			}
		}
	}

	// Check for feed option updates (both feeds exist, but options differ).
	for categoryTitle, feeds := range local.GetFeedsByCategory() {
		for _, localFeed := range feeds {
			if remote.FeedExists(localFeed.URL, categoryTitle) {
				remoteOptions := remote.GetFeedOptions(localFeed.URL)
				if needsUpdate(localFeed.Options, remoteOptions) {
					actions = append(actions, Action{
						Type:          UpdateFeed,
						CategoryTitle: categoryTitle,
						FeedURL:       localFeed.URL,
						FeedOptions:   localFeed.Options,
					})
				}
			}
		}
	}

	sort.Sort(ActionSorter(actions))

	return actions, nil
}

// needsUpdate checks if local options differ from remote and require an update.
func needsUpdate(local, remote FeedOptions) bool {
	if local.IsEmpty() {
		return false // No options specified locally, no update needed.
	}
	return !local.Equal(remote)
}
