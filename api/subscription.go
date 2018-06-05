// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/miniflux/miniflux/http/response/json"
	"github.com/miniflux/miniflux/reader/subscription"
)

// GetSubscriptions is the API handler to find subscriptions.
func (c *Controller) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	websiteURL, err := decodeURLPayload(r.Body)
	if err != nil {
		json.BadRequest(w, err)
		return
	}

	subscriptions, err := subscription.FindSubscriptions(websiteURL)
	if err != nil {
		json.ServerError(w, errors.New("Unable to discover subscriptions"))
		return
	}

	if subscriptions == nil {
		json.NotFound(w, fmt.Errorf("No subscription found"))
		return
	}

	json.OK(w, subscriptions)
}
