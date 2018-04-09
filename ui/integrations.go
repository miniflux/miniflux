// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui

import (
	"crypto/md5"
	"fmt"

	"github.com/miniflux/miniflux/http/handler"
	"github.com/miniflux/miniflux/ui/form"
)

// ShowIntegrations renders the page with all external integrations.
func (c *Controller) ShowIntegrations(ctx *handler.Context, request *handler.Request, response *handler.Response) {
	user := ctx.LoggedUser()
	integration, err := c.store.Integration(user.ID)
	if err != nil {
		response.HTML().ServerError(err)
		return
	}

	args, err := c.getCommonTemplateArgs(ctx)
	if err != nil {
		response.HTML().ServerError(err)
		return
	}

	response.HTML().Render("integrations", args.Merge(tplParams{
		"menu": "settings",
		"form": form.IntegrationForm{
			PinboardEnabled:      integration.PinboardEnabled,
			PinboardToken:        integration.PinboardToken,
			PinboardTags:         integration.PinboardTags,
			PinboardMarkAsUnread: integration.PinboardMarkAsUnread,
			InstapaperEnabled:    integration.InstapaperEnabled,
			InstapaperUsername:   integration.InstapaperUsername,
			InstapaperPassword:   integration.InstapaperPassword,
			FeverEnabled:         integration.FeverEnabled,
			FeverUsername:        integration.FeverUsername,
			FeverPassword:        integration.FeverPassword,
			WallabagEnabled:      integration.WallabagEnabled,
			WallabagURL:          integration.WallabagURL,
			WallabagClientID:     integration.WallabagClientID,
			WallabagClientSecret: integration.WallabagClientSecret,
			WallabagUsername:     integration.WallabagUsername,
			WallabagPassword:     integration.WallabagPassword,
			NunuxKeeperEnabled:   integration.NunuxKeeperEnabled,
			NunuxKeeperURL:       integration.NunuxKeeperURL,
			NunuxKeeperAPIKey:    integration.NunuxKeeperAPIKey,
		},
	}))
}

// UpdateIntegration updates integration settings.
func (c *Controller) UpdateIntegration(ctx *handler.Context, request *handler.Request, response *handler.Response) {
	user := ctx.LoggedUser()
	integration, err := c.store.Integration(user.ID)
	if err != nil {
		response.HTML().ServerError(err)
		return
	}

	integrationForm := form.NewIntegrationForm(request.Request())
	integrationForm.Merge(integration)

	if integration.FeverUsername != "" && c.store.HasDuplicateFeverUsername(user.ID, integration.FeverUsername) {
		ctx.SetFlashErrorMessage(ctx.Translate("There is already someone else with the same Fever username!"))
		response.Redirect(ctx.Route("integrations"))
		return
	}

	if integration.FeverEnabled {
		integration.FeverToken = fmt.Sprintf("%x", md5.Sum([]byte(integration.FeverUsername+":"+integration.FeverPassword)))
	} else {
		integration.FeverToken = ""
	}

	err = c.store.UpdateIntegration(integration)
	if err != nil {
		response.HTML().ServerError(err)
		return
	}

	response.Redirect(ctx.Route("integrations"))
}
