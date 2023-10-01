// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package components provides api functions for managing Zarf components.
package components

import (
	"encoding/json"
	"net/http"

	"github.com/defenseunicorns/zarf-ui/src/api/cluster"
	"github.com/defenseunicorns/zarf-ui/src/api/common"
	"github.com/defenseunicorns/zarf/src/config"
	"github.com/defenseunicorns/zarf/src/config/lang"
	"github.com/defenseunicorns/zarf/src/pkg/k8s"
	"github.com/defenseunicorns/zarf/src/pkg/message"
	"github.com/defenseunicorns/zarf/src/types"
	"github.com/go-chi/chi"
)

// ListDeployedComponents writes a list of packages that have been deployed to the connected cluster.
func ListDeployedComponents(w http.ResponseWriter, r *http.Request) {
	pkgName := chi.URLParam(r, "pkg")

	var deployedPackage = types.DeployedPackage{}

	k, err := k8s.New(message.Debugf, cluster.Labels)
	if err != nil {
		message.ErrorWebf(err, w, lang.ErrLoadPackageSecret, pkgName)
	}

	// Get the secret that describes the deployed init package
	secret, err := k.GetSecret("zarf", config.ZarfPackagePrefix+pkgName)
	if err != nil {
		message.ErrorWebf(err, w, lang.ErrLoadPackageSecret, pkgName)
	}

	err = json.Unmarshal(secret.Data["data"], &deployedPackage)
	if err != nil {
		message.ErrorWebf(err, w, lang.ErrLoadPackageSecret, pkgName)
	}
	common.WriteJSONResponse(w, deployedPackage.DeployedComponents, http.StatusOK)
}
