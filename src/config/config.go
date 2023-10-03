// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package config stores the global configuration and constants for the UI.
package config

import "embed"

// Zarf UI Global Configuration Variables
var (
	UIAssets  embed.FS
	UIVersion = "unset"
)
