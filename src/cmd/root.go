// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package cmd contains the CLI commands for Zarf.
package cmd

import (
	"os"

	"github.com/defenseunicorns/zarf-ui/src/api"
	zCmd "github.com/defenseunicorns/zarf/src/cmd"
	"github.com/defenseunicorns/zarf/src/cmd/common"
	"github.com/defenseunicorns/zarf/src/config"
	"github.com/defenseunicorns/zarf/src/config/lang"
	"github.com/defenseunicorns/zarf/src/pkg/message"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "zarf-ui",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip for vendor-only commands
		if common.CheckVendorOnlyFromPath(cmd) {
			return
		}

		// Don't log the help command
		if cmd.Parent() == nil {
			config.SkipLogFile = true
		}

		common.SetupCLI()
	},
	Short: lang.RootCmdShort,
	Long:  lang.RootCmdLong,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api.LaunchAPIServer()
	},
}

var zarfCmd = &cobra.Command{
	Use:                "zarf",
	Short:              "Aliased internal Zarf commands.",
	Long:               "Aliased internal Zarf commands that can be referenced in Zarf 'actions'",
	DisableFlagParsing: true,
	Hidden:             true,
	Run: func(cmd *cobra.Command, args []string) {
		os.Args = os.Args[1:]
		zCmd.Execute()
	},
}

// Execute is the entrypoint for the CLI.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(zarfCmd)

	// Skip for vendor-only commands
	if common.CheckVendorOnlyFromArgs() {
		return
	}

	v := common.InitViper()

	v.SetDefault(common.VLogLevel, "info")
	v.SetDefault(common.VZarfCache, config.ZarfDefaultCachePath)

	rootCmd.PersistentFlags().StringVarP(&common.LogLevelCLI, "log-level", "l", v.GetString(common.VLogLevel), lang.RootCmdFlagLogLevel)
	rootCmd.PersistentFlags().StringVarP(&config.CLIArch, "architecture", "a", v.GetString(common.VArchitecture), lang.RootCmdFlagArch)
	rootCmd.PersistentFlags().BoolVar(&config.SkipLogFile, "no-log-file", v.GetBool(common.VNoLogFile), lang.RootCmdFlagSkipLogFile)
	rootCmd.PersistentFlags().BoolVar(&message.NoProgress, "no-progress", v.GetBool(common.VNoProgress), lang.RootCmdFlagNoProgress)
	rootCmd.PersistentFlags().BoolVar(&config.NoColor, "no-color", v.GetBool(common.VNoColor), lang.RootCmdFlagNoColor)
	rootCmd.PersistentFlags().StringVar(&config.CommonOptions.CachePath, "zarf-cache", v.GetString(common.VZarfCache), lang.RootCmdFlagCachePath)
	rootCmd.PersistentFlags().StringVar(&config.CommonOptions.TempDirectory, "tmpdir", v.GetString(common.VTmpDir), lang.RootCmdFlagTempDir)
	rootCmd.PersistentFlags().BoolVar(&config.CommonOptions.Insecure, "insecure", v.GetBool(common.VInsecure), lang.RootCmdFlagInsecure)
}
