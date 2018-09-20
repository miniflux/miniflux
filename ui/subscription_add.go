// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui  // import "miniflux.app/ui"

import (
	"net/http"

	"miniflux.app/http/client"
	"miniflux.app/http/response/html"
	"miniflux.app/http/request"
	"miniflux.app/ui/session"
	"miniflux.app/ui/view"
)

// AddSubscription shows the form to add a new feed.
func (c *Controller) AddSubscription(w http.ResponseWriter, r *http.Request) {
	sess := session.New(c.store, request.SessionID(r))
	view := view.New(c.tpl, r, sess)

	user, err := c.store.UserByID(request.UserID(r))
	if err != nil {
		html.ServerError(w, err)
		return
	}

	categories, err := c.store.Categories(user.ID)
	if err != nil {
		html.ServerError(w, err)
		return
	}

	view.Set("categories", categories)
	view.Set("menu", "feeds")
	view.Set("user", user)
	view.Set("countUnread", c.store.CountUnreadEntries(user.ID))
	view.Set("countErrorFeeds", c.store.CountErrorFeeds(user.ID))
	view.Set("defaultUserAgent", client.DefaultUserAgent)

	html.OK(w, r, view.Render("add_subscription"))
}
