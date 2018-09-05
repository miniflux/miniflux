// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package api // import "miniflux.app/api"

import (
	"errors"
	"net/http"

	"miniflux.app/http/request"
	"miniflux.app/http/response/json"
)

// FeedIcon returns a feed icon.
func (c *Controller) FeedIcon(w http.ResponseWriter, r *http.Request) {
	feedID, err := request.IntParam(r, "feedID")
	if err != nil {
		json.BadRequest(w, err)
		return
	}

	if !c.store.HasIcon(feedID) {
		json.NotFound(w, errors.New("This feed doesn't have any icon"))
		return
	}

	icon, err := c.store.IconByFeedID(request.UserID(r), feedID)
	if err != nil {
		json.ServerError(w, err)
		return
	}

	if icon == nil {
		json.NotFound(w, errors.New("This feed doesn't have any icon"))
		return
	}

	json.OK(w, r, &feedIcon{
		ID:       icon.ID,
		MimeType: icon.MimeType,
		Data:     icon.DataURL(),
	})
}
