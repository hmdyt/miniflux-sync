package diff

// FeedOptions represents the configurable options for a Miniflux feed.
// All fields are pointers to distinguish between "not set" and "set to zero value".
type FeedOptions struct {
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

// Feed represents a feed with its URL and optional configuration.
type Feed struct {
	URL     string
	Options FeedOptions
}

// IsEmpty returns true if no options are set.
func (o FeedOptions) IsEmpty() bool {
	return o.Crawler == nil &&
		o.Username == nil &&
		o.Password == nil &&
		o.UserAgent == nil &&
		o.Cookie == nil &&
		o.Disabled == nil &&
		o.IgnoreHTTPCache == nil &&
		o.FetchViaProxy == nil &&
		o.AllowSelfSignedCertificates == nil &&
		o.DisableHTTP2 == nil &&
		o.ScraperRules == nil &&
		o.RewriteRules == nil &&
		o.BlocklistRules == nil &&
		o.KeeplistRules == nil &&
		o.HideGlobally == nil
}

// Equal compares two FeedOptions for equality.
// Only compares fields that are set in the receiver (local).
func (o FeedOptions) Equal(other FeedOptions) bool {
	if o.Crawler != nil && !boolPtrEqual(o.Crawler, other.Crawler) {
		return false
	}
	if o.Username != nil && !stringPtrEqual(o.Username, other.Username) {
		return false
	}
	if o.Password != nil && !stringPtrEqual(o.Password, other.Password) {
		return false
	}
	if o.UserAgent != nil && !stringPtrEqual(o.UserAgent, other.UserAgent) {
		return false
	}
	if o.Cookie != nil && !stringPtrEqual(o.Cookie, other.Cookie) {
		return false
	}
	if o.Disabled != nil && !boolPtrEqual(o.Disabled, other.Disabled) {
		return false
	}
	if o.IgnoreHTTPCache != nil && !boolPtrEqual(o.IgnoreHTTPCache, other.IgnoreHTTPCache) {
		return false
	}
	if o.FetchViaProxy != nil && !boolPtrEqual(o.FetchViaProxy, other.FetchViaProxy) {
		return false
	}
	if o.AllowSelfSignedCertificates != nil && !boolPtrEqual(o.AllowSelfSignedCertificates, other.AllowSelfSignedCertificates) {
		return false
	}
	if o.DisableHTTP2 != nil && !boolPtrEqual(o.DisableHTTP2, other.DisableHTTP2) {
		return false
	}
	if o.ScraperRules != nil && !stringPtrEqual(o.ScraperRules, other.ScraperRules) {
		return false
	}
	if o.RewriteRules != nil && !stringPtrEqual(o.RewriteRules, other.RewriteRules) {
		return false
	}
	if o.BlocklistRules != nil && !stringPtrEqual(o.BlocklistRules, other.BlocklistRules) {
		return false
	}
	if o.KeeplistRules != nil && !stringPtrEqual(o.KeeplistRules, other.KeeplistRules) {
		return false
	}
	if o.HideGlobally != nil && !boolPtrEqual(o.HideGlobally, other.HideGlobally) {
		return false
	}
	return true
}

func boolPtrEqual(a, b *bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func stringPtrEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
