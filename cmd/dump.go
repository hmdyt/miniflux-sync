package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/revett/miniflux-sync/api"
	"github.com/revett/miniflux-sync/config"
	"github.com/revett/miniflux-sync/diff"
	"github.com/revett/miniflux-sync/log"
	"gopkg.in/yaml.v2"
	miniflux "miniflux.app/v2/client"
)

func dump(ctx context.Context, flags *config.DumpFlags, client *miniflux.Client) error {
	log.Info(ctx, "exporting data from miniflux")

	feeds, categories, err := api.FetchData(ctx, client)
	if err != nil {
		return errors.Wrap(err, "fetching data")
	}

	remoteState, err := api.GenerateDiffState(feeds, categories)
	if err != nil {
		return errors.Wrap(err, "generating remote state")
	}

	// e.g. "miniflux-sync-remote-20240718105851_opml.xml"
	filename := fmt.Sprintf("./miniflux-sync-remote-%s.yml", time.Now().Format("20060102150405"))
	if flags.Path != "" {
		log.Info(ctx, `using export path from "--path"`, log.Metadata{
			"path": flags.Path,
		})
		filename = flags.Path
	}

	log.Info(ctx, "writing export data to file")

	output := buildDumpOutput(remoteState)
	dat, err := yaml.Marshal(output)
	if err != nil {
		return errors.Wrap(err, "marshalling remote state to yaml")
	}

	if err := os.WriteFile(filename, dat, 0o600); err != nil { //nolint:mnd
		return errors.Wrap(err, "writing export data to file")
	}

	log.Info(ctx, filename)
	return nil
}

// dumpFeedEntry represents a feed entry in the dump output.
type dumpFeedEntry struct {
	URL                         string  `yaml:"url"`
	Crawler                     *bool   `yaml:"crawler,omitempty"`
	Username                    *string `yaml:"username,omitempty"`
	Password                    *string `yaml:"password,omitempty"`
	UserAgent                   *string `yaml:"user_agent,omitempty"`
	Cookie                      *string `yaml:"cookie,omitempty"`
	Disabled                    *bool   `yaml:"disabled,omitempty"`
	IgnoreHTTPCache             *bool   `yaml:"ignore_http_cache,omitempty"`
	FetchViaProxy               *bool   `yaml:"fetch_via_proxy,omitempty"`
	AllowSelfSignedCertificates *bool   `yaml:"allow_self_signed_certificates,omitempty"`
	DisableHTTP2                *bool   `yaml:"disable_http2,omitempty"`
	ScraperRules                *string `yaml:"scraper_rules,omitempty"`
	RewriteRules                *string `yaml:"rewrite_rules,omitempty"`
	BlocklistRules              *string `yaml:"blocklist_rules,omitempty"`
	KeeplistRules               *string `yaml:"keeplist_rules,omitempty"`
	HideGlobally                *bool   `yaml:"hide_globally,omitempty"`
}

// buildDumpOutput builds the output for the dump command.
// If a feed has non-default options, it outputs the object format.
// Otherwise, it outputs just the URL string.
func buildDumpOutput(state *diff.State) map[string][]interface{} {
	output := make(map[string][]interface{})

	for category, feeds := range state.FeedsByCategoryTitle {
		for _, feed := range feeds {
			if hasNonDefaultOptions(feed.Options) {
				entry := dumpFeedEntry{
					URL:                         feed.URL,
					Crawler:                     nonDefaultBool(feed.Options.Crawler, false),
					Username:                    feed.Options.Username,
					Password:                    feed.Options.Password,
					UserAgent:                   feed.Options.UserAgent,
					Cookie:                      feed.Options.Cookie,
					Disabled:                    nonDefaultBool(feed.Options.Disabled, false),
					IgnoreHTTPCache:             nonDefaultBool(feed.Options.IgnoreHTTPCache, false),
					FetchViaProxy:               nonDefaultBool(feed.Options.FetchViaProxy, false),
					AllowSelfSignedCertificates: nonDefaultBool(feed.Options.AllowSelfSignedCertificates, false),
					DisableHTTP2:                nonDefaultBool(feed.Options.DisableHTTP2, false),
					ScraperRules:                feed.Options.ScraperRules,
					RewriteRules:                feed.Options.RewriteRules,
					BlocklistRules:              feed.Options.BlocklistRules,
					KeeplistRules:               feed.Options.KeeplistRules,
					HideGlobally:                nonDefaultBool(feed.Options.HideGlobally, false),
				}
				output[category] = append(output[category], entry)
			} else {
				output[category] = append(output[category], feed.URL)
			}
		}
	}

	return output
}

// hasNonDefaultOptions checks if a feed has any non-default options set.
func hasNonDefaultOptions(opts diff.FeedOptions) bool {
	// Check bool options for non-default (true) values
	if opts.Crawler != nil && *opts.Crawler {
		return true
	}
	if opts.Disabled != nil && *opts.Disabled {
		return true
	}
	if opts.IgnoreHTTPCache != nil && *opts.IgnoreHTTPCache {
		return true
	}
	if opts.FetchViaProxy != nil && *opts.FetchViaProxy {
		return true
	}
	if opts.AllowSelfSignedCertificates != nil && *opts.AllowSelfSignedCertificates {
		return true
	}
	if opts.DisableHTTP2 != nil && *opts.DisableHTTP2 {
		return true
	}
	if opts.HideGlobally != nil && *opts.HideGlobally {
		return true
	}
	// Check string options for non-empty values
	if opts.Username != nil && *opts.Username != "" {
		return true
	}
	if opts.Password != nil && *opts.Password != "" {
		return true
	}
	if opts.UserAgent != nil && *opts.UserAgent != "" {
		return true
	}
	if opts.Cookie != nil && *opts.Cookie != "" {
		return true
	}
	if opts.ScraperRules != nil && *opts.ScraperRules != "" {
		return true
	}
	if opts.RewriteRules != nil && *opts.RewriteRules != "" {
		return true
	}
	if opts.BlocklistRules != nil && *opts.BlocklistRules != "" {
		return true
	}
	if opts.KeeplistRules != nil && *opts.KeeplistRules != "" {
		return true
	}
	return false
}

// nonDefaultBool returns the pointer only if the value differs from the default.
func nonDefaultBool(ptr *bool, defaultVal bool) *bool {
	if ptr == nil || *ptr == defaultVal {
		return nil
	}
	return ptr
}
