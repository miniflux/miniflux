// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package model // import "miniflux.app/model"

import (
	"fmt"
	"math"
	"time"

	"miniflux.app/config"
	"miniflux.app/http/client"
)

// List of supported schedulers.
const (
	SchedulerRoundRobin     = "round_robin"
	SchedulerEntryFrequency = "entry_frequency"
	// Default settings for the feed query builder
	DefaultFeedSorting          = "parsing_error_count"
	DefaultFeedSortingDirection = "desc"
)

// Feed represents a feed in the application.
type Feed struct {
	ID                 int64     `json:"id"`
	UserID             int64     `json:"user_id"`
	FeedURL            string    `json:"feed_url"`
	SiteURL            string    `json:"site_url"`
	Title              string    `json:"title"`
	CheckedAt          time.Time `json:"checked_at"`
	NextCheckAt        time.Time `json:"next_check_at"`
	EtagHeader         string    `json:"etag_header"`
	LastModifiedHeader string    `json:"last_modified_header"`
	ParsingErrorMsg    string    `json:"parsing_error_message"`
	ParsingErrorCount  int       `json:"parsing_error_count"`
	ScraperRules       string    `json:"scraper_rules"`
	RewriteRules       string    `json:"rewrite_rules"`
	Crawler            bool      `json:"crawler"`
	BlocklistRules     string    `json:"blocklist_rules"`
	KeeplistRules      string    `json:"keeplist_rules"`
	UserAgent          string    `json:"user_agent"`
	Username           string    `json:"username"`
	Password           string    `json:"password"`
	Disabled           bool      `json:"disabled"`
	IgnoreHTTPCache    bool      `json:"ignore_http_cache"`
	FetchViaProxy      bool      `json:"fetch_via_proxy"`
	PollingInterval    int       `json:"polling_interval_minutes"`
	Category           *Category `json:"category,omitempty"`
	Entries            Entries   `json:"entries,omitempty"`
	Icon               *FeedIcon `json:"icon"`
	UnreadCount        int       `json:"-"`
	ReadCount          int       `json:"-"`
}

func (f *Feed) String() string {
	return fmt.Sprintf("ID=%d, UserID=%d, FeedURL=%s, SiteURL=%s, Title=%s, Category={%s}",
		f.ID,
		f.UserID,
		f.FeedURL,
		f.SiteURL,
		f.Title,
		f.Category,
	)
}

// WithClientResponse updates feed attributes from an HTTP request.
func (f *Feed) WithClientResponse(response *client.Response) {
	f.EtagHeader = response.ETag
	f.LastModifiedHeader = response.LastModified
	f.FeedURL = response.EffectiveURL
}

// WithCategoryID initializes the category attribute of the feed.
func (f *Feed) WithCategoryID(categoryID int64) {
	f.Category = &Category{ID: categoryID}
}

// WithError adds a new error message and increment the error counter.
func (f *Feed) WithError(message string) {
	f.ParsingErrorCount++
	f.ParsingErrorMsg = message
}

// ResetErrorCounter removes all previous errors.
func (f *Feed) ResetErrorCounter() {
	f.ParsingErrorCount = 0
	f.ParsingErrorMsg = ""
}

// CheckedNow set attribute values when the feed is refreshed.
func (f *Feed) CheckedNow() {
	f.CheckedAt = time.Now()

	if f.SiteURL == "" {
		f.SiteURL = f.FeedURL
	}
}

// ScheduleNextCheck set "next_check_at" of a feed based on the scheduler selected from the configuration.
func (f *Feed) ScheduleNextCheck(weeklyCount int, pollingInterval int) {
	const compensationSeconds = 30
	var intervalMinutes int
	if pollingInterval > 0 {
		intervalMinutes = pollingInterval
	} else {
		switch config.Opts.PollingScheduler() {
		case SchedulerEntryFrequency:
			if weeklyCount == 0 {
				intervalMinutes = config.Opts.SchedulerEntryFrequencyMaxInterval()
			} else {
				intervalMinutes = int(math.Round(float64(7*24*60) / float64(weeklyCount)))
			}
			intervalMinutes = int(math.Min(float64(intervalMinutes), float64(config.Opts.SchedulerEntryFrequencyMaxInterval())))
			intervalMinutes = int(math.Max(float64(intervalMinutes), float64(config.Opts.SchedulerEntryFrequencyMinInterval())))
		default:
			intervalMinutes = config.Opts.PollingFrequency()
		}
	}
	// The compensationSeconds compensates the time different between job starts and NextCheckAt is set.
	// For example, the scheduler is round robin and polling interval is 5 minutes.
	// Without the compensationSeconds, the following may happend:
	// (1) The first job starts at 0s; (2) The job sql query runs at 5s, selects a feed, subject to NextCheckAt < now;
	// (3) The NextCheckAt is set at 10s, the value is 5m10s;
	// (4) The next job starts at 5m; (5) The job sql query runs at 5m5s, but it won't get the feed because
	// NextCheckAt (5m10s) is later than now (5m5s). Instead the feed will be fetched by the job run at 10m.
	// The compensationSeconds tries to solve the problem, letting the feed gets fetched by the job run at 5m.
	// The compensation should be smaller than the polling interval. Since the minimum resolution of the polling
	// interval is 1 minute, we use a value smaller than 1 minute.
	// As the compensation is applied to all feeds, it won't change the order when feeds will be fetched.
	f.NextCheckAt = time.Now().Add(time.Second * time.Duration(intervalMinutes*60-compensationSeconds))
}

// FeedCreationRequest represents the request to create a feed.
type FeedCreationRequest struct {
	FeedURL         string `json:"feed_url"`
	CategoryID      int64  `json:"category_id"`
	UserAgent       string `json:"user_agent"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Crawler         bool   `json:"crawler"`
	Disabled        bool   `json:"disabled"`
	IgnoreHTTPCache bool   `json:"ignore_http_cache"`
	FetchViaProxy   bool   `json:"fetch_via_proxy"`
	ScraperRules    string `json:"scraper_rules"`
	RewriteRules    string `json:"rewrite_rules"`
	BlocklistRules  string `json:"blocklist_rules"`
	KeeplistRules   string `json:"keeplist_rules"`
	PollingInterval int    `json:"polling_interval_minutes"`
}

// FeedModificationRequest represents the request to update a feed.
type FeedModificationRequest struct {
	FeedURL         *string `json:"feed_url"`
	SiteURL         *string `json:"site_url"`
	Title           *string `json:"title"`
	ScraperRules    *string `json:"scraper_rules"`
	RewriteRules    *string `json:"rewrite_rules"`
	BlocklistRules  *string `json:"blocklist_rules"`
	KeeplistRules   *string `json:"keeplist_rules"`
	Crawler         *bool   `json:"crawler"`
	UserAgent       *string `json:"user_agent"`
	Username        *string `json:"username"`
	Password        *string `json:"password"`
	CategoryID      *int64  `json:"category_id"`
	Disabled        *bool   `json:"disabled"`
	IgnoreHTTPCache *bool   `json:"ignore_http_cache"`
	FetchViaProxy   *bool   `json:"fetch_via_proxy"`
	PollingInterval *int    `json:"polling_interval_minutes"`
}

// Patch updates a feed with modified values.
func (f *FeedModificationRequest) Patch(feed *Feed) {
	if f.FeedURL != nil && *f.FeedURL != "" {
		feed.FeedURL = *f.FeedURL
	}

	if f.SiteURL != nil && *f.SiteURL != "" {
		feed.SiteURL = *f.SiteURL
	}

	if f.Title != nil && *f.Title != "" {
		feed.Title = *f.Title
	}

	if f.ScraperRules != nil {
		feed.ScraperRules = *f.ScraperRules
	}

	if f.RewriteRules != nil {
		feed.RewriteRules = *f.RewriteRules
	}

	if f.KeeplistRules != nil {
		feed.KeeplistRules = *f.KeeplistRules
	}

	if f.BlocklistRules != nil {
		feed.BlocklistRules = *f.BlocklistRules
	}

	if f.Crawler != nil {
		feed.Crawler = *f.Crawler
	}

	if f.UserAgent != nil {
		feed.UserAgent = *f.UserAgent
	}

	if f.Username != nil {
		feed.Username = *f.Username
	}

	if f.Password != nil {
		feed.Password = *f.Password
	}

	if f.CategoryID != nil && *f.CategoryID > 0 {
		feed.Category.ID = *f.CategoryID
	}

	if f.Disabled != nil {
		feed.Disabled = *f.Disabled
	}

	if f.IgnoreHTTPCache != nil {
		feed.IgnoreHTTPCache = *f.IgnoreHTTPCache
	}

	if f.FetchViaProxy != nil {
		feed.FetchViaProxy = *f.FetchViaProxy
	}

	if f.PollingInterval != nil {
		feed.PollingInterval = *f.PollingInterval
	}
}

// Feeds is a list of feed
type Feeds []*Feed
