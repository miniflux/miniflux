// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/miniflux/miniflux/http/cookie"
	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response/html"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/model"
)

// AppSession handles application session middleware.
func (m *Middleware) AppSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		session := m.getAppSessionValueFromCookie(r)

		if session == nil {
			logger.Debug("[Middleware:AppSession] Session not found")

			session, err = m.store.CreateSession()
			if err != nil {
				logger.Error("[Middleware:AppSession] %v", err)
				html.ServerError(w, err)
				return
			}

			http.SetCookie(w, cookie.New(cookie.CookieSessionID, session.ID, m.cfg.IsHTTPS, m.cfg.BasePath()))
		} else {
			logger.Debug("[Middleware:AppSession] %s", session)
		}

		if r.Method == "POST" {
			formValue := r.FormValue("csrf")
			headerValue := r.Header.Get("X-Csrf-Token")

			if session.Data.CSRF != formValue && session.Data.CSRF != headerValue {
				logger.Error(`[Middleware:AppSession] Invalid or missing CSRF token: Form="%s", Header="%s"`, formValue, headerValue)
				html.BadRequest(w, errors.New("invalid or missing CSRF"))
				return
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, SessionIDContextKey, session.ID)
		ctx = context.WithValue(ctx, CSRFContextKey, session.Data.CSRF)
		ctx = context.WithValue(ctx, OAuth2StateContextKey, session.Data.OAuth2State)
		ctx = context.WithValue(ctx, FlashMessageContextKey, session.Data.FlashMessage)
		ctx = context.WithValue(ctx, FlashErrorMessageContextKey, session.Data.FlashErrorMessage)
		ctx = context.WithValue(ctx, UserLanguageContextKey, session.Data.Language)
		ctx = context.WithValue(ctx, PocketRequestTokenContextKey, session.Data.PocketRequestToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) getAppSessionValueFromCookie(r *http.Request) *model.Session {
	cookieValue := request.Cookie(r, cookie.CookieSessionID)
	if cookieValue == "" {
		return nil
	}

	session, err := m.store.Session(cookieValue)
	if err != nil {
		logger.Error("[Middleware:AppSession] %v", err)
		return nil
	}

	return session
}
