// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package packages provides api functions for managing Zarf packages.
package packages

import (
	"io"
	"net/http"
	"net/url"

	"github.com/defenseunicorns/zarf-ui/src/api/common"
	"github.com/defenseunicorns/zarf-ui/src/types"
	"github.com/defenseunicorns/zarf/src/pkg/layout"
	"github.com/defenseunicorns/zarf/src/pkg/message"
	"github.com/go-chi/chi/v5"
	goyaml "github.com/goccy/go-yaml"
	"github.com/mholt/archiver/v3"
)

// Read reads a package from the local filesystem and writes the Zarf.yaml json to the response.
func Read(w http.ResponseWriter, r *http.Request) {
	message.Debug("packages.Read()")

	path := chi.URLParam(r, "path")

	if pkg, err := ReadPackage(path); err != nil {
		message.ErrorWebf(err, w, "Unable to read the package at: `%s`", path)
	} else {
		common.WriteJSONResponse(w, pkg, http.StatusOK)
	}
}

// ReadPackage reads a packages yaml from the local filesystem and returns an APIZarfPackage.
func ReadPackage(path string) (pkg types.APIZarfPackage, err error) {
	var file []byte

	pkg.Path, err = url.QueryUnescape(path)
	if err != nil {
		return pkg, err
	}

	// Check for zarf.yaml in the package and read into file
	err = archiver.Walk(pkg.Path, func(f archiver.File) error {
		if f.Name() == layout.ZarfYAML {
			file, err = io.ReadAll(f)
			if err != nil {
				return err
			}
			return archiver.ErrStopWalk
		}

		return nil
	})
	if err != nil {
		return pkg, err
	}

	err = goyaml.Unmarshal(file, &pkg.ZarfPackage)
	return pkg, err
}
