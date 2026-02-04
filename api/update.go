package api

import (
	"context"

	"github.com/pkg/errors"
	"github.com/revett/miniflux-sync/diff"
	"github.com/revett/miniflux-sync/log"
	miniflux "miniflux.app/v2/client"
)

// Update performs a series of actions on the Miniflux instance.
func Update( //nolint:cyclop,funlen
	ctx context.Context,
	client *miniflux.Client,
	actions []diff.Action,
	feeds []*miniflux.Feed,
	categories []*miniflux.Category,
) error {
	log.Info(ctx, "performing actions")

	for _, action := range actions {
		switch action.Type {
		case diff.CreateCategory:
			log.Info(ctx, "creating category", log.Metadata{
				"title": action.CategoryTitle,
			})

			category, err := client.CreateCategory(action.CategoryTitle)
			if err != nil {
				return errors.Wrap(err, "creating category")
			}

			categories = append(categories, category)

		case diff.CreateFeed:
			log.Info(ctx, "creating feed", log.Metadata{
				"category": action.CategoryTitle,
				"url":      action.FeedURL,
			})

			categoryID, err := findCategoryIDByTitle(action.CategoryTitle, categories)
			if err != nil {
				return errors.Wrap(err, "finding category id")
			}

			req := miniflux.FeedCreationRequest{
				FeedURL:    action.FeedURL,
				CategoryID: categoryID,
			}

			applyOptionsToCreationRequest(&req, action.FeedOptions)

			feedID, err := client.CreateFeed(&req)
			if err != nil {
				return errors.Wrap(err, "creating feed")
			}

			feed, err := client.Feed(feedID)
			if err != nil {
				return errors.Wrap(err, "fetching feed")
			}

			feeds = append(feeds, feed)

		case diff.UpdateFeed:
			log.Info(ctx, "updating feed", log.Metadata{
				"category": action.CategoryTitle,
				"url":      action.FeedURL,
			})

			feedID, err := findFeedIDByURL(action.FeedURL, feeds)
			if err != nil {
				return errors.Wrap(err, "finding feed id")
			}

			modReq := miniflux.FeedModificationRequest{}
			applyOptionsToModificationRequest(&modReq, action.FeedOptions)

			_, err = client.UpdateFeed(feedID, &modReq)
			if err != nil {
				return errors.Wrap(err, "updating feed")
			}

		case diff.DeleteCategory:
			log.Info(ctx, "deleting category", log.Metadata{
				"title": action.CategoryTitle,
			})

			categoryID, err := findCategoryIDByTitle(action.CategoryTitle, categories)
			if err != nil {
				return errors.Wrap(err, "finding category id")
			}

			if err := client.DeleteCategory(categoryID); err != nil {
				return errors.Wrap(err, "deleting category")
			}

			categories = removeCategoryByID(categoryID, categories)

		case diff.DeleteFeed:
			log.Info(ctx, "deleting feed", log.Metadata{
				"category": action.CategoryTitle,
				"url":      action.FeedURL,
			})

			feedID, err := findFeedIDByURL(action.FeedURL, feeds)
			if err != nil {
				return errors.Wrap(err, "finding feed id")
			}

			if err := client.DeleteFeed(feedID); err != nil {
				return errors.Wrap(err, "deleting feed")
			}

			feeds = removeFeedByID(feedID, feeds)

		default:
			return errors.Errorf(`unknown action type: "%s"`, action.Type)
		}
	}

	return nil
}

func findCategoryIDByTitle(title string, categories []*miniflux.Category) (int64, error) {
	for _, category := range categories {
		if category.Title == title {
			return category.ID, nil
		}
	}

	return 0, errors.Errorf(`category not found: "%s"`, title)
}

func findFeedIDByURL(url string, feeds []*miniflux.Feed) (int64, error) {
	for _, feed := range feeds {
		if feed.FeedURL == url {
			return feed.ID, nil
		}
	}

	return 0, errors.Errorf(`feed not found: "%s"`, url)
}

func removeCategoryByID(id int64, categories []*miniflux.Category) []*miniflux.Category {
	for i, category := range categories {
		if category.ID == id {
			return append(categories[:i], categories[i+1:]...)
		}
	}

	return categories
}

func removeFeedByID(id int64, feeds []*miniflux.Feed) []*miniflux.Feed {
	for i, feed := range feeds {
		if feed.ID == id {
			return append(feeds[:i], feeds[i+1:]...)
		}
	}

	return feeds
}

// applyOptionsToCreationRequest applies feed options to a FeedCreationRequest.
func applyOptionsToCreationRequest(req *miniflux.FeedCreationRequest, opts diff.FeedOptions) {
	if opts.Crawler != nil {
		req.Crawler = *opts.Crawler
	}
	if opts.Username != nil {
		req.Username = *opts.Username
	}
	if opts.Password != nil {
		req.Password = *opts.Password
	}
	if opts.UserAgent != nil {
		req.UserAgent = *opts.UserAgent
	}
	if opts.Cookie != nil {
		req.Cookie = *opts.Cookie
	}
	if opts.Disabled != nil {
		req.Disabled = *opts.Disabled
	}
	if opts.IgnoreHTTPCache != nil {
		req.IgnoreHTTPCache = *opts.IgnoreHTTPCache
	}
	if opts.FetchViaProxy != nil {
		req.FetchViaProxy = *opts.FetchViaProxy
	}
	if opts.AllowSelfSignedCertificates != nil {
		req.AllowSelfSignedCertificates = *opts.AllowSelfSignedCertificates
	}
	if opts.DisableHTTP2 != nil {
		req.DisableHTTP2 = *opts.DisableHTTP2
	}
	if opts.ScraperRules != nil {
		req.ScraperRules = *opts.ScraperRules
	}
	if opts.RewriteRules != nil {
		req.RewriteRules = *opts.RewriteRules
	}
	if opts.BlocklistRules != nil {
		req.BlocklistRules = *opts.BlocklistRules
	}
	if opts.KeeplistRules != nil {
		req.KeeplistRules = *opts.KeeplistRules
	}
	if opts.HideGlobally != nil {
		req.HideGlobally = *opts.HideGlobally
	}
}

// applyOptionsToModificationRequest applies feed options to a FeedModificationRequest.
func applyOptionsToModificationRequest(req *miniflux.FeedModificationRequest, opts diff.FeedOptions) {
	req.Crawler = opts.Crawler
	req.Username = opts.Username
	req.Password = opts.Password
	req.UserAgent = opts.UserAgent
	req.Cookie = opts.Cookie
	req.Disabled = opts.Disabled
	req.IgnoreHTTPCache = opts.IgnoreHTTPCache
	req.FetchViaProxy = opts.FetchViaProxy
	req.AllowSelfSignedCertificates = opts.AllowSelfSignedCertificates
	req.DisableHTTP2 = opts.DisableHTTP2
	req.ScraperRules = opts.ScraperRules
	req.RewriteRules = opts.RewriteRules
	req.BlocklistRules = opts.BlocklistRules
	req.KeeplistRules = opts.KeeplistRules
	req.HideGlobally = opts.HideGlobally
}
