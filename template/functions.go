// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package template // import "miniflux.app/template"

import (
	"fmt"
	"math"
	"html/template"
	"net/mail"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"miniflux.app/config"
	"miniflux.app/filter"
	"miniflux.app/http/route"
	"miniflux.app/locale"
	"miniflux.app/model"
	"miniflux.app/timezone"
	"miniflux.app/url"
)

type funcMap struct {
	cfg    *config.Config
	router *mux.Router
}

// Map returns a map of template functions that are compiled during template parsing.
func (f *funcMap) Map() template.FuncMap {
	return template.FuncMap{
		"dict":     dict,
		"hasKey":   hasKey,
		"truncate": truncate,
		"isEmail":  isEmail,
		"baseURL": func() string {
			return f.cfg.BaseURL()
		},
		"rootURL": func() string {
			return f.cfg.RootURL()
		},
		"hasOAuth2Provider": func(provider string) bool {
			return f.cfg.OAuth2Provider() == provider
		},
		"route": func(name string, args ...interface{}) string {
			return route.Path(f.router, name, args...)
		},
		"noescape": func(str string) template.HTML {
			return template.HTML(str)
		},
		"proxyFilter": func(data string) string {
			return filter.ImageProxyFilter(f.router, f.cfg, data)
		},
		"proxyURL": func(link string) string {
			proxyImages := f.cfg.ProxyImages()

			if proxyImages == "all" || (proxyImages != "none" && !url.IsHTTPS(link)) {
				return filter.Proxify(f.router, link)
			}

			return link
		},
		"domain": func(websiteURL string) string {
			return url.Domain(websiteURL)
		},
		"hasPrefix": func(str, prefix string) bool {
			return strings.HasPrefix(str, prefix)
		},
		"contains": func(str, substr string) bool {
			return strings.Contains(str, substr)
		},
		"isodate": func(ts time.Time) string {
			return ts.Format("2006-01-02 15:04:05")
		},
		"theme_color": func(theme string) string {
			return model.ThemeColor(theme)
		},

		// These functions are overrided at runtime after the parsing.
		"elapsed": func(timezone string, t time.Time) string {
			return ""
		},
		"t": func(key interface{}, args ...interface{}) string {
			return ""
		},
		"plural": func(key string, n int, args ...interface{}) string {
			return ""
		},
	}
}

func newFuncMap(cfg *config.Config, router *mux.Router) *funcMap {
	return &funcMap{cfg, router}
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("dict expects an even number of arguments")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func hasKey(dict map[string]string, key string) bool {
	if value, found := dict[key]; found {
		return value != ""
	}
	return false
}

func truncate(str string, max int) string {
	runes := 0
	for i := range str {
		runes++
		if runes > max {
			return str[:i] + "…"
		}
	}
	return str
}

func isEmail(str string) bool {
	_, err := mail.ParseAddress(str)
	if err != nil {
		return false
	}
	return true
}

func elapsedTime(language *locale.Language, tz string, t time.Time) string {
	if t.IsZero() {
		return language.Get("time_elapsed.not_yet")
	}

	now := timezone.Now(tz)
	t = timezone.Convert(tz, t)
	if now.Before(t) {
		return language.Get("time_elapsed.not_yet")
	}

	diff := now.Sub(t)
	// Duration in seconds
	s := diff.Seconds()
	// Duration in days
	d := int(s / 86400)
	switch {
	case s < 60:
		return language.Get("time_elapsed.now")
	case s < 3600:
		minutes := int(diff.Minutes())
		return language.Plural("time_elapsed.minutes", minutes, minutes)
	case s < 86400:
		hours := int(diff.Hours())
		return language.Plural("time_elapsed.hours", hours, hours)
	case d == 1:
		return language.Get("time_elapsed.yesterday")
	case d < 7:
		return language.Plural("time_elapsed.days", d, d)
	case d < 31:
		weeks := int(math.Ceil(float64(d) / 7))
		return language.Plural("time_elapsed.weeks", weeks, weeks)
	case d < 365:
		months := int(math.Ceil(float64(d) / 30))
		return language.Plural("time_elapsed.months", months, months)
	default:
		years := int(math.Ceil(float64(d) / 365))
		return language.Plural("time_elapsed.years", years, years)
	}
}
