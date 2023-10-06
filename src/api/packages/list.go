// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package packages provides api functions for managing Zarf packages.
package packages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/defenseunicorns/zarf-ui/src/api/cluster"
	"github.com/defenseunicorns/zarf-ui/src/api/common"
	"github.com/defenseunicorns/zarf/src/config"
	"github.com/defenseunicorns/zarf/src/pkg/k8s"
	"github.com/defenseunicorns/zarf/src/pkg/message"
	zTypes "github.com/defenseunicorns/zarf/src/types"
)

// ListDeployedPackages writes a list of packages that have been deployed to the connected cluster.
func ListDeployedPackages(w http.ResponseWriter, _ *http.Request) {

	k, err := k8s.New(message.Debugf, cluster.Labels)
	if err != nil {
		message.ErrorWebf(err, w, "Could not connect to cluster")
		return
	}

	var deployedPackages = []zTypes.DeployedPackage{}
	var errorList []error
	// Get the secrets that describe the deployed packages
	secrets, err := k.GetSecretsWithLabel("zarf", "package-deploy-info")
	if err != nil {
		message.ErrorWebf(err, w, "Unable to get a list of the deployed Zarf packages")
		return
	}

	// Process the k8s secret into our internal structs
	for _, secret := range secrets.Items {
		if strings.HasPrefix(secret.Name, config.ZarfPackagePrefix) {
			var deployedPackage zTypes.DeployedPackage
			err := json.Unmarshal(secret.Data["data"], &deployedPackage)
			// add the error to the error list
			if err != nil {
				errorList = append(errorList, fmt.Errorf("unable to unmarshal the secret %s/%s", secret.Namespace, secret.Name))
			} else {
				deployedPackages = append(deployedPackages, deployedPackage)
			}
		}
	}

	// TODO #1312: Handle errors where some deployedPackages were able to be parsed
	if len(errorList) > 0 && len(deployedPackages) == 0 {
		message.ErrorWebf(err, w, "Unable to get a list of the deployed Zarf packages")
		return
	}

	common.WriteJSONResponse(w, deployedPackages, http.StatusOK)
}
