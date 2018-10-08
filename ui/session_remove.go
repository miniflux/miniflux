// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui // import "miniflux.app/ui"

import (
	"net/http"

	"miniflux.app/http/request"
	"miniflux.app/http/response/html"
	"miniflux.app/http/route"
	"miniflux.app/logger"
)

// RemoveSession remove a user session.
func (c *Controller) RemoveSession(w http.ResponseWriter, r *http.Request) {
	sessionID := request.RouteInt64Param(r, "sessionID")
	err := c.store.RemoveUserSessionByID(request.UserID(r), sessionID)
	if err != nil {
		logger.Error("[Controller:RemoveSession] %v", err)
	}

	html.Redirect(w, r, route.Path(c.router, "sessions"))
}
