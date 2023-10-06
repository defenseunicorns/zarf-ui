// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package api provides the UI API server.
package api

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/defenseunicorns/zarf-ui/src/api/auth"
	"github.com/defenseunicorns/zarf-ui/src/api/cluster"
	"github.com/defenseunicorns/zarf-ui/src/api/components"
	"github.com/defenseunicorns/zarf-ui/src/api/packages"
	"github.com/defenseunicorns/zarf-ui/src/config"
	"github.com/defenseunicorns/zarf/src/pkg/layout"
	"github.com/defenseunicorns/zarf/src/pkg/message"
	"github.com/defenseunicorns/zarf/src/pkg/utils"
	"github.com/defenseunicorns/zarf/src/pkg/utils/exec"
	"github.com/defenseunicorns/zarf/src/pkg/utils/helpers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// LaunchAPIServer launches UI API server.
func LaunchAPIServer() {
	message.Debug("api.LaunchAPIServer()")

	// Track the developer port if it's set
	devPort := os.Getenv("API_DEV_PORT")

	// If the env variable API_PORT is set, use that for the listening port
	port := os.Getenv("API_PORT")
	// Otherwise, use a random available port
	if port == "" {
		// If we can't find an available port, just use the default
		if portRaw, err := helpers.GetAvailablePort(); err != nil {
			port = "8080"
		} else {
			port = fmt.Sprintf("%d", portRaw)
		}
	}

	ip := "127.0.0.1"

	// Track the external IP if it's set
	externalIP := os.Getenv("API_EXTERNAL_IP")
	if externalIP != "" {
		ip = externalIP
	}

	// If the env variable API_TOKEN is set, use that for the API secret
	token := os.Getenv("API_TOKEN")
	// Otherwise, generate a random secret
	if token == "" {
		token = utils.RandomString(96)
	}

	// Init the Chi router
	router := chi.NewRouter()

	// Push logs into the message buffer for log persistence
	genericMsg := message.Generic{}
	logFormatter := middleware.DefaultLogFormatter{
		Logger: log.New(&genericMsg, "API CALL | ", log.LstdFlags),
	}

	router.Use(middleware.RequestLogger(&logFormatter))
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/api", func(r chi.Router) {
		// Require a valid token for API calls
		r.Use(auth.RequireSecret(token))
		r.Use(middleware.NoCache)

		r.Head("/", auth.Connect)

		r.Route("/cluster", func(r chi.Router) {
			r.Get("/", cluster.Summary)
		})

		r.Route("/packages", func(r chi.Router) {
			r.Get("/read/{path}", packages.Read)
			r.Get("/list", packages.ListDeployedPackages)
			r.Put("/deploy", packages.DeployPackage)
			r.Get("/deploy-stream", packages.StreamDeployPackage)
			r.Delete("/remove/{name}", packages.RemovePackage)
			r.Put("/{pkg}/connect/{name}", packages.ConnectTunnel)
			r.Delete("/{pkg}/disconnect/{name}", packages.DisconnectTunnel)
			r.Get("/{pkg}/connections", packages.ListPackageConnections)
			r.Get("/{pkg}/components/deployed", components.ListDeployedComponents)
			r.Get("/connections", packages.ListConnections)
			r.Get("/sbom/{path}", packages.ExtractSBOM)
			r.Delete("/sbom", packages.DeleteSBOM)
			r.Route("/find", func(r chi.Router) {
				r.Route("/stream", func(r chi.Router) {
					r.Get("/", packages.FindPackageStream)
					r.Get("/init", packages.FindInitStream)
					r.Get("/home", packages.FindInHomeStream)
				})
			})
		})
	})

	// If no dev port specified, use the server port for the URL and try to open it
	if devPort == "" {
		url := fmt.Sprintf("http://%s:%s/auth?token=%s", ip, port, token)
		message.Infof("Zarf UI connection: %s", url)
		message.Debug(exec.LaunchURL(url))
	} else {
		// Otherwise, use the dev port for the URL and don't try to open
		message.Infof("Zarf UI connection: http://%s:%s/auth?token=%s", ip, devPort, token)
	}

	// Setup the static SBOM server
	sbomSub := os.DirFS(layout.SBOMDir)
	sbomFs := http.FileServer(http.FS(sbomSub))

	// Serve the SBOM viewer files
	router.Get("/sbom-viewer/*", func(w http.ResponseWriter, r *http.Request) {
		message.Debug("api.LaunchAPIServer() - /sbom-viewer/*")

		// Extract the file name from the URL
		file := strings.TrimPrefix(r.URL.Path, "/sbom-viewer/")

		// Ensure SBOM file exists in the config.ZarfSBOMDir
		if test, err := sbomSub.Open(file); err != nil {
			// If the file doesn't exist, redirect to the homepage
			r.URL.Path = "/"
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			// If the file exists, close the file and serve it
			test.Close()
		}
		r.URL.Path = file
		sbomFs.ServeHTTP(w, r)
	})

	// Load the static UI files
	if sub, err := fs.Sub(config.UIAssets, "build/ui"); err != nil {
		message.WarnErr(err, "Unable to load the embedded ui assets")
	} else {
		// Setup a file server for the static UI files
		fs := http.FileServer(http.FS(sub))

		// Catch all routes
		router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			message.Debug("api.LaunchAPIServer() - /*")
			// If the request is not a real file, serve the index.html instead
			if test, err := sub.Open(strings.TrimPrefix(r.URL.Path, "/")); err != nil {
				r.URL.Path = "/"
			} else {
				test.Close()
			}
			fs.ServeHTTP(w, r)
		})
	}

	http.ListenAndServe(":"+port, router)
}
