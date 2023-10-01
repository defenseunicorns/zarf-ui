// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package cluster contains Zarf-specific cluster management functions.
package cluster

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/defenseunicorns/zarf-ui/src/api/common"
	"github.com/defenseunicorns/zarf/src/config"
	"github.com/defenseunicorns/zarf/src/pkg/k8s"
	"github.com/defenseunicorns/zarf/src/pkg/message"
	"github.com/defenseunicorns/zarf/src/types"
	"k8s.io/client-go/tools/clientcmd"
)

// Labels are the Kubernetes labels applied to resources by Zarf
var Labels = k8s.Labels{
	config.ZarfManagedByLabel: "zarf",
}

// Summary returns a summary of cluster status.
func Summary(w http.ResponseWriter, _ *http.Request) {
	message.Debug("cluster.Summary()")

	var state *types.ZarfState
	var reachable bool
	var distro string
	var hasZarf bool
	var k8sRevision string

	k, err := k8s.NewWithWait(message.Debugf, Labels, 5*time.Second)
	rawConfig, _ := clientcmd.NewDefaultClientConfigLoadingRules().GetStartingConfig()

	reachable = err == nil
	if reachable {
		distro, _ = k.DetectDistro()
		state, _ = loadZarfState(k)
		hasZarf = state != nil
		k8sRevision, _ = k.GetServerVersion()
	}

	data := types.ClusterSummary{
		Reachable:   reachable,
		HasZarf:     hasZarf,
		Distro:      distro,
		ZarfState:   state,
		K8sRevision: k8sRevision,
		RawConfig:   rawConfig,
	}

	common.WriteJSONResponse(w, data, http.StatusOK)
}

func loadZarfState(k *k8s.K8s) (state *types.ZarfState, err error) {
	secret, err := k.GetSecret("zarf", "zarf-state")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(secret.Data["state"], &state)
	if err != nil {
		return nil, err
	}

	return state, nil
}
