// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui

import (
	"net/http"

	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/response"
	"github.com/miniflux/miniflux/http/response/html"
	"github.com/miniflux/miniflux/http/route"
	"github.com/miniflux/miniflux/locale"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/ui/form"
	"github.com/miniflux/miniflux/ui/session"
	"github.com/miniflux/miniflux/ui/view"
)

// UpdateSettings update the settings.
func (c *Controller) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	sess := session.New(c.store, ctx)
	view := view.New(c.tpl, ctx, sess)

	user, err := c.store.UserByID(ctx.UserID())
	if err != nil {
		html.ServerError(w, err)
		return
	}

	timezones, err := c.store.Timezones()
	if err != nil {
		html.ServerError(w, err)
		return
	}

	settingsForm := form.NewSettingsForm(r)

	view.Set("form", settingsForm)
	view.Set("themes", model.Themes())
	view.Set("languages", locale.AvailableLanguages())
	view.Set("timezones", timezones)
	view.Set("menu", "settings")
	view.Set("user", user)
	view.Set("countUnread", c.store.CountUnreadEntries(user.ID))

	if err := settingsForm.Validate(); err != nil {
		view.Set("errorMessage", err.Error())
		html.OK(w, view.Render("settings"))
		return
	}

	if c.store.AnotherUserExists(user.ID, settingsForm.Username) {
		view.Set("errorMessage", "This user already exists.")
		html.OK(w, view.Render("settings"))
		return
	}

	err = c.store.UpdateUser(settingsForm.Merge(user))
	if err != nil {
		logger.Error("[Controller:UpdateSettings] %v", err)
		view.Set("errorMessage", "Unable to update this user.")
		html.OK(w, view.Render("settings"))
		return
	}

	sess.SetLanguage(user.Language)
	sess.NewFlashMessage(c.translator.GetLanguage(ctx.UserLanguage()).Get("Preferences saved!"))
	response.Redirect(w, r, route.Path(c.router, "settings"))
}
