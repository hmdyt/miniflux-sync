package parse

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/revett/miniflux-sync/diff"
	"github.com/revett/miniflux-sync/log"
	"gopkg.in/yaml.v2"
)

// feedEntry represents a single entry in the YAML that can be either a string or an object.
type feedEntry struct {
	URL     string
	Options diff.FeedOptions
}

// UnmarshalYAML implements custom unmarshaling for mixed format support.
func (f *feedEntry) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try as string first (simple URL format)
	var urlString string
	if err := unmarshal(&urlString); err == nil {
		f.URL = urlString
		return nil
	}

	// Try as object format
	var raw struct {
		URL                         string  `yaml:"url"`
		Crawler                     *bool   `yaml:"crawler"`
		Username                    *string `yaml:"username"`
		Password                    *string `yaml:"password"`
		UserAgent                   *string `yaml:"user_agent"`
		Cookie                      *string `yaml:"cookie"`
		Disabled                    *bool   `yaml:"disabled"`
		IgnoreHTTPCache             *bool   `yaml:"ignore_http_cache"`
		FetchViaProxy               *bool   `yaml:"fetch_via_proxy"`
		AllowSelfSignedCertificates *bool   `yaml:"allow_self_signed_certificates"`
		DisableHTTP2                *bool   `yaml:"disable_http2"`
		ScraperRules                *string `yaml:"scraper_rules"`
		RewriteRules                *string `yaml:"rewrite_rules"`
		BlocklistRules              *string `yaml:"blocklist_rules"`
		KeeplistRules               *string `yaml:"keeplist_rules"`
		HideGlobally                *bool   `yaml:"hide_globally"`
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	if raw.URL == "" {
		return errors.New("feed entry must have a url field")
	}

	f.URL = raw.URL
	f.Options = diff.FeedOptions{
		Crawler:                     raw.Crawler,
		Username:                    raw.Username,
		Password:                    raw.Password,
		UserAgent:                   raw.UserAgent,
		Cookie:                      raw.Cookie,
		Disabled:                    raw.Disabled,
		IgnoreHTTPCache:             raw.IgnoreHTTPCache,
		FetchViaProxy:               raw.FetchViaProxy,
		AllowSelfSignedCertificates: raw.AllowSelfSignedCertificates,
		DisableHTTP2:                raw.DisableHTTP2,
		ScraperRules:                raw.ScraperRules,
		RewriteRules:                raw.RewriteRules,
		BlocklistRules:              raw.BlocklistRules,
		KeeplistRules:               raw.KeeplistRules,
		HideGlobally:                raw.HideGlobally,
	}

	return nil
}

// Parse reads a YAML file to a diff.State struct.
func Parse(ctx context.Context, path string) (*diff.State, error) {
	log.Info(ctx, "reading data from yaml file")
	log.Info(ctx, path)

	data, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		return nil, errors.Wrap(err, "reading data from file")
	}

	var rawData map[string][]feedEntry
	if err := yaml.Unmarshal(data, &rawData); err != nil {
		return nil, errors.Wrap(err, "unmarshalling data")
	}

	state := diff.State{
		FeedURLsByCategoryTitle: map[string][]string{},
		FeedsByCategoryTitle:    map[string][]diff.Feed{},
	}

	for category, entries := range rawData {
		for _, entry := range entries {
			state.FeedURLsByCategoryTitle[category] = append(
				state.FeedURLsByCategoryTitle[category], entry.URL)

			state.FeedsByCategoryTitle[category] = append(
				state.FeedsByCategoryTitle[category], diff.Feed{
					URL:     entry.URL,
					Options: entry.Options,
				})
		}
	}

	if err := validateDuplicateFeedURLs(&state); err != nil {
		return nil, errors.Wrap(err, "validating duplicate feed urls")
	}

	return &state, nil
}

func validateDuplicateFeedURLs(state *diff.State) error {
	feedURLSet := make(map[string]struct{})

	for _, urls := range state.FeedURLsByCategoryTitle {
		for _, url := range urls {
			if _, exists := feedURLSet[url]; exists {
				return errors.Errorf(`duplicate url found across categories: "%s"`, url)
			}

			feedURLSet[url] = struct{}{}
		}
	}

	return nil
}
