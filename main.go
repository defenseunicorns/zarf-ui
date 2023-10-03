// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package main is the entrypoint for the Zarf UI binary.
package main

import (
	"embed"

	"github.com/defenseunicorns/zarf-ui/src/cmd"
	"github.com/defenseunicorns/zarf-ui/src/config"
	zConfig "github.com/defenseunicorns/zarf/src/config"
)

//go:embed all:build/ui/*
var assets embed.FS

func main() {
	zConfig.ActionsCommandZarfPrefix = "zarf"
	config.UIAssets = assets
	cmd.Execute()
}
