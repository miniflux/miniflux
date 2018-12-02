// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui // import "miniflux.app/ui"

import (
	"net/http"

	"miniflux.app/http/request"
	"miniflux.app/http/response/json"
	"miniflux.app/model"
	"miniflux.app/reader/sanitizer"
	"miniflux.app/reader/scraper"
)

func (h *handler) fetchContent(w http.ResponseWriter, r *http.Request) {
	entryID := request.RouteInt64Param(r, "entryID")
	builder := h.store.NewEntryQueryBuilder(request.UserID(r))
	builder.WithEntryID(entryID)
	builder.WithoutStatus(model.EntryStatusRemoved)

	entry, err := builder.GetEntry()
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	if entry == nil {
		json.NotFound(w, r)
		return
	}

	content, err := scraper.Fetch(entry.URL, entry.Feed.ScraperRules, entry.Feed.UserAgent)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}
	
	content = rewrite.Rewriter(entry.URL, content, entry.Feed.RewriteRules)

	entry.Content = sanitizer.Sanitize(entry.URL, content)
	h.store.UpdateEntryContent(entry)

	json.OK(w, r, map[string]string{"content": entry.Content})
}
