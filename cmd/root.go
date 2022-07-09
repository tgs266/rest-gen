/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
	restgen "github.com/tgs266/rest-gen/rest-gen"
	"github.com/tgs266/rest-gen/rest-gen/config"
)

var rootCmd = &cobra.Command{
	Use: "rgen",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		inputDir, _ := cmd.Flags().GetString("input-dir")
		outputDir, _ := cmd.Flags().GetString("output-dir")

		noPath := configPath == ""
		noInputDir := inputDir == ""
		noOutputDir := outputDir == ""

		// no config passed and not both input and output
		if noPath && (noInputDir || noOutputDir) {
			panic("must pass either a config or both input and output directories")
		}
		var cfg config.Config
		if !noPath {
			cfg = config.Read(configPath)
		} else {
			cfg = config.Config{
				Definitions: config.Defintions{
					InputDir:  inputDir,
					OutputDir: outputDir,
				},
			}
		}

		restgen.Generate(cfg)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringP("config", "c", "", "rest-gen yaml config")
	rootCmd.Flags().StringP("input-dir", "i", "", "input directory")
	rootCmd.Flags().StringP("output-dir", "o", "", "output directory")
}
