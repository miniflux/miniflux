// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package middleware

import (
	"context"
	"net/http"

	"github.com/miniflux/miniflux/http/cookie"
	"github.com/miniflux/miniflux/http/route"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/storage"

	"github.com/gorilla/mux"
)

// UserSessionMiddleware represents a user session middleware.
type UserSessionMiddleware struct {
	store  *storage.Storage
	router *mux.Router
}

// Handler execute the middleware.
func (s *UserSessionMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := s.getSessionFromCookie(r)

		if session == nil {
			logger.Debug("[Middleware:UserSession] Session not found")
			if s.isPublicRoute(r) {
				next.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, route.Path(s.router, "login"), http.StatusFound)
			}
		} else {
			logger.Debug("[Middleware:UserSession] %s", session)
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDContextKey, session.UserID)
			ctx = context.WithValue(ctx, IsAuthenticatedContextKey, true)
			ctx = context.WithValue(ctx, UserSessionTokenContextKey, session.Token)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func (s *UserSessionMiddleware) isPublicRoute(r *http.Request) bool {
	route := mux.CurrentRoute(r)
	switch route.GetName() {
	case "login",
		"checkLogin",
		"stylesheet",
		"javascript",
		"oauth2Redirect",
		"oauth2Callback",
		"appIcon",
		"favicon",
		"webManifest":
		return true
	default:
		return false
	}
}

func (s *UserSessionMiddleware) getSessionFromCookie(r *http.Request) *model.UserSession {
	sessionCookie, err := r.Cookie(cookie.CookieUserSessionID)
	if err == http.ErrNoCookie {
		return nil
	}

	session, err := s.store.UserSessionByToken(sessionCookie.Value)
	if err != nil {
		logger.Error("[Middleware:UserSession] %v", err)
		return nil
	}

	return session
}

// NewUserSessionMiddleware returns a new UserSessionMiddleware.
func NewUserSessionMiddleware(s *storage.Storage, r *mux.Router) *UserSessionMiddleware {
	return &UserSessionMiddleware{store: s, router: r}
}
