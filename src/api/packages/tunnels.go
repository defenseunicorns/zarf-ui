// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package packages provides api functions for managing Zarf packages.
package packages

import (
	"errors"
	"net/http"

	"github.com/defenseunicorns/zarf-ui/src/api/cluster"
	"github.com/defenseunicorns/zarf-ui/src/api/common"
	"github.com/defenseunicorns/zarf-ui/src/types"
	zConfig "github.com/defenseunicorns/zarf/src/config"
	"github.com/defenseunicorns/zarf/src/pkg/k8s"
	"github.com/defenseunicorns/zarf/src/pkg/message"
	"github.com/go-chi/chi/v5"
)

// PackageTunnel is a struct for storing a tunnel and its connection details
type PackageTunnel struct {
	tunnel     *k8s.Tunnel
	Connection types.APIDeployedPackageConnection `json:"connection,omitempty"`
}

// packageTunnels is a map of package names to PackageTunnel objects
type packageTunnels map[string]map[string]PackageTunnel

// tunnels is a map of package names to tunnel objects used for storing connected tunnels
var tunnels = make(packageTunnels)

// ListConnections returns a map of pkgName to a list of connections
func ListConnections(w http.ResponseWriter, _ *http.Request) {
	allConnections := make(types.APIConnections)
	for name, pkgTunnels := range tunnels {
		for _, pkgTunnel := range pkgTunnels {
			if allConnections[name] == nil {
				allConnections[name] = make(types.APIDeployedPackageConnections, 0)
			}
			allConnections[name] = append(allConnections[name], pkgTunnel.Connection)
		}
	}
	common.WriteJSONResponse(w, allConnections, http.StatusOK)
}

// ListPackageConnections lists all tunnel names
func ListPackageConnections(w http.ResponseWriter, r *http.Request) {
	pkgName := chi.URLParam(r, "pkg")
	if tunnels[pkgName] == nil {
		message.ErrorWebf(errors.New("no tunnels for package %s"), w, pkgName)
		return
	}
	pkgTunnels := make(types.APIDeployedPackageConnections, 0, len(tunnels[pkgName]))
	for _, pkgTunnel := range tunnels[pkgName] {
		pkgTunnels = append(pkgTunnels, pkgTunnel.Connection)
	}

	common.WriteJSONResponse(w, pkgTunnels, http.StatusOK)
}

// ConnectTunnel establishes a tunnel for the requested resource
func ConnectTunnel(w http.ResponseWriter, r *http.Request) {
	pkgName := chi.URLParam(r, "pkg")
	connectionName := chi.URLParam(r, "name")

	if tunnels[pkgName] == nil {
		tunnels[pkgName] = make(map[string]PackageTunnel)
	}

	pkgTunnels := tunnels[pkgName]

	if pkgTunnels[connectionName].tunnel != nil {
		common.WriteJSONResponse(w, tunnels[pkgName][connectionName].Connection, http.StatusOK)
		return
	}

	k, err := k8s.New(message.Debugf, cluster.Labels)
	if err != nil {
		message.ErrorWebf(err, w, "Failed to connect to cluster for %s", connectionName)
		return
	}
	matches, err := k.GetServicesByLabel("", zConfig.ZarfConnectLabelName, connectionName)
	if err != nil {
		message.ErrorWebf(err, w, "Unable to lookup the service: %s", err.Error())
		return
	}

	if len(matches.Items) == 0 {
		message.ErrorWebf(err, w, "Unable to find any matching connect services")
		return
	}

	svc := matches.Items[0]
	remotePort := svc.Spec.Ports[0].TargetPort.IntValue()
	// if remotePort == 0, look for Port (which is required)
	if remotePort == 0 {
		remotePort = k.FindPodContainerPort(svc)
	}

	tunnel, err := k.NewTunnel(svc.Namespace, k8s.SvcResource, svc.Name,
		svc.Annotations[zConfig.ZarfConnectAnnotationURL], 0, remotePort)
	if err != nil {
		message.ErrorWebf(err, w, "Failed to create tunnel for %s", connectionName)
		return
	}

	url, err := tunnel.Connect()
	if err != nil {
		message.ErrorWebf(err, w, "Failed to connect to %s", connectionName)
		return
	}

	tunnels[pkgName][connectionName] = PackageTunnel{
		tunnel: tunnel,
		Connection: types.APIDeployedPackageConnection{
			Name: connectionName,
			URL:  url,
		},
	}
	common.WriteJSONResponse(w, tunnels[pkgName][connectionName].Connection, http.StatusCreated)
}

// DisconnectTunnel closes the tunnel for the requested resource
func DisconnectTunnel(w http.ResponseWriter, r *http.Request) {
	pkgName := chi.URLParam(r, "pkg")
	connectionName := chi.URLParam(r, "name")
	pkgTunnel := tunnels[pkgName][connectionName]
	if pkgTunnel.tunnel == nil {
		message.ErrorWebf(errors.New("tunnel not found"), w, "Failed to disconnect from %s", connectionName)
		return
	}

	pkgTunnel.tunnel.Close()
	delete(tunnels[pkgName], connectionName)

	common.WriteJSONResponse(w, true, http.StatusOK)
}
