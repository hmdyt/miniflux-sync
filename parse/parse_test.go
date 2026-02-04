package parse

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/revett/miniflux-sync/log"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		yaml         string
		expectedURLs map[string][]string
		wantCrawler  map[string]bool // feedURL -> crawler value
	}{
		"StringOnly": {
			yaml: `Tech:
  - https://example.com/feed.xml
  - https://example2.com/feed.xml`,
			expectedURLs: map[string][]string{
				"Tech": {
					"https://example.com/feed.xml",
					"https://example2.com/feed.xml",
				},
			},
			wantCrawler: map[string]bool{},
		},
		"ObjectWithOptions": {
			yaml: `Tech:
  - url: https://example.com/feed.xml
    crawler: true`,
			expectedURLs: map[string][]string{
				"Tech": {"https://example.com/feed.xml"},
			},
			wantCrawler: map[string]bool{
				"https://example.com/feed.xml": true,
			},
		},
		"MixedFormat": {
			yaml: `Tech:
  - https://simple.com/feed.xml
  - url: https://full.com/feed.xml
    crawler: true
    username: user
    password: pass`,
			expectedURLs: map[string][]string{
				"Tech": {
					"https://simple.com/feed.xml",
					"https://full.com/feed.xml",
				},
			},
			wantCrawler: map[string]bool{
				"https://full.com/feed.xml": true,
			},
		},
		"MultipleCategories": {
			yaml: `Tech:
  - https://tech.com/feed.xml
News:
  - url: https://news.com/feed.xml
    crawler: true`,
			expectedURLs: map[string][]string{
				"Tech": {"https://tech.com/feed.xml"},
				"News": {"https://news.com/feed.xml"},
			},
			wantCrawler: map[string]bool{
				"https://news.com/feed.xml": true,
			},
		},
		"AllOptions": {
			yaml: `Tech:
  - url: https://full.com/feed.xml
    crawler: true
    username: user
    password: pass
    user_agent: "Custom UA"
    cookie: "session=abc"
    disabled: true
    ignore_http_cache: true
    fetch_via_proxy: true
    allow_self_signed_certificates: true
    disable_http2: true
    scraper_rules: "article"
    rewrite_rules: "add_image_title"
    blocklist_rules: "(?i)ad"
    keeplist_rules: "(?i)important"
    hide_globally: true`,
			expectedURLs: map[string][]string{
				"Tech": {"https://full.com/feed.xml"},
			},
			wantCrawler: map[string]bool{
				"https://full.com/feed.xml": true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Create temp file with YAML content
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "feeds.yml")
			err := os.WriteFile(tmpFile, []byte(tc.yaml), 0o600)
			require.NoError(t, err)

			// Parse the file with logger context
			logger := log.New()
			ctx := logger.WithContext(context.Background())
			state, err := Parse(ctx, tmpFile)
			require.NoError(t, err)

			// Check FeedURLsByCategoryTitle
			require.Equal(t, tc.expectedURLs, state.FeedURLsByCategoryTitle)

			// Check FeedsByCategoryTitle for crawler options
			for feedURL, expectedCrawler := range tc.wantCrawler {
				opts := state.GetFeedOptions(feedURL)
				require.NotNil(t, opts.Crawler, "crawler should be set for %s", feedURL)
				require.Equal(t, expectedCrawler, *opts.Crawler, "crawler mismatch for %s", feedURL)
			}
		})
	}
}

func TestParse_DuplicateURL(t *testing.T) {
	t.Parallel()

	yaml := `Tech:
  - https://example.com/feed.xml
News:
  - https://example.com/feed.xml`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "feeds.yml")
	err := os.WriteFile(tmpFile, []byte(yaml), 0o600)
	require.NoError(t, err)

	logger := log.New()
	ctx := logger.WithContext(context.Background())
	_, err = Parse(ctx, tmpFile)
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate")
}

func TestParse_MissingURL(t *testing.T) {
	t.Parallel()

	yaml := `Tech:
  - crawler: true`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "feeds.yml")
	err := os.WriteFile(tmpFile, []byte(yaml), 0o600)
	require.NoError(t, err)

	logger := log.New()
	ctx := logger.WithContext(context.Background())
	_, err = Parse(ctx, tmpFile)
	require.Error(t, err)
	require.Contains(t, err.Error(), "url")
}
