// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package types contains all the types used by Zarf.
package types

import (
	zTypes "github.com/defenseunicorns/zarf/src/types"
	"k8s.io/client-go/tools/clientcmd/api"
)

// RestAPI is the struct that is used to marshal/unmarshal the top-level API objects.
type RestAPI struct {
	ZarfPackage                   zTypes.ZarfPackage            `json:"zarfPackage"`
	ZarfState                     zTypes.ZarfState              `json:"zarfState"`
	ZarfCommonOptions             zTypes.ZarfCommonOptions      `json:"zarfCommonOptions"`
	ZarfCreateOptions             zTypes.ZarfCreateOptions      `json:"zarfCreateOptions"`
	ZarfPackageOptions            zTypes.ZarfPackageOptions     `json:"zarfPackageOptions"`
	ZarfInitOptions               zTypes.ZarfInitOptions        `json:"zarfInitOptions"`
	ConnectStrings                zTypes.ConnectStrings         `json:"connectStrings"`
	ClusterSummary                APIClusterSummary             `json:"clusterSummary"`
	DeployedPackage               zTypes.DeployedPackage        `json:"deployedPackage"`
	APIZarfPackage                APIZarfPackage                `json:"apiZarfPackage"`
	APIZarfDeployPayload          APIZarfDeployPayload          `json:"apiZarfDeployPayload"`
	APIZarfPackageConnection      APIDeployedPackageConnection  `json:"apiZarfPackageConnection"`
	APIDeployedPackageConnections APIDeployedPackageConnections `json:"apiZarfPackageConnections"`
	APIConnections                APIConnections                `json:"apiConnections"`
	APIPackageSBOM                APIPackageSBOM                `json:"apiPackageSBOM"`
}

// APIClusterSummary contains the summary of a cluster for the API.
type APIClusterSummary struct {
	Reachable   bool              `json:"reachable"`
	HasZarf     bool              `json:"hasZarf"`
	Distro      string            `json:"distro"`
	ZarfState   *zTypes.ZarfState `json:"zarfState"`
	K8sRevision string            `json:"k8sRevision"`
	RawConfig   *api.Config       `json:"rawConfig"`
}

// APIZarfPackage represents a ZarfPackage and its path for the API.
type APIZarfPackage struct {
	Path        string             `json:"path"`
	ZarfPackage zTypes.ZarfPackage `json:"zarfPackage"`
}

// APIZarfDeployPayload represents the needed data to deploy a ZarfPackage/ZarfInit
type APIZarfDeployPayload struct {
	PackageOpts zTypes.ZarfPackageOptions `json:"packageOpts"`
	InitOpts    *zTypes.ZarfInitOptions   `json:"initOpts,omitempty"`
}

// APIPackageSBOM represents the SBOM viewer files for a package
type APIPackageSBOM struct {
	Path  string   `json:"path"`
	SBOMS []string `json:"sboms"`
}

// APIConnections represents all of the existing connections
type APIConnections map[string]APIDeployedPackageConnections

// APIDeployedPackageConnections represents all of the connections for a deployed package
type APIDeployedPackageConnections []APIDeployedPackageConnection

// APIDeployedPackageConnection represents a single connection from a deployed package
type APIDeployedPackageConnection struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}
