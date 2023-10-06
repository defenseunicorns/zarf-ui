// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package auth provides an endpoint for authenticating against the Zarf UI server.
package auth

import (
	"net/http"
)

// RequireSecret ensures the request has a valid token.
func RequireSecret(validToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Header.Get("Authorization") != validToken {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Connect is a head-only request to test the connection.
func Connect(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
