// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui

import (
	"net/http"
	"time"

	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response"
	"github.com/miniflux/miniflux/ui/static"
)

// Stylesheet renders the CSS.
func (c *Controller) Stylesheet(w http.ResponseWriter, r *http.Request) {
	stylesheet := request.Param(r, "name", "white")
	body := static.Stylesheets["common"]
	etag := static.StylesheetsChecksums["common"]

	if theme, found := static.Stylesheets[stylesheet]; found {
		body += theme
		etag += static.StylesheetsChecksums[stylesheet]
	}

	response.Cache(w, r, "text/css; charset=utf-8", etag, []byte(body), 48*time.Hour)
}
