/*
Copyright © 2022 Joaquin Genova <joaquingenovag8@gmail.com>

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
	"fmt"
	"goto/src/config"
	"os"

	"github.com/spf13/cobra"
)

// listGPathCmd represents the listGPath command
var listGPathCmd = &cobra.Command{
	Use:   "list",
	Short: "List all goto-paths in the goto-paths file",

	Run: func(cmd *cobra.Command, args []string) {

		reverse, err := cmd.Flags().GetBool("reverse")
		cobra.CheckErr(err)

		//Initial the variables to use config
		config.GotoPathsFile = GotoPathsFile

		var gpaths []config.GotoPath
		cobra.CheckErr(config.LoadConfigFile(&gpaths))

		if reverse {
			listReverse(gpaths)
			os.Exit(0)
		}

		list(gpaths)
		os.Exit(0)
	},
}

func list(gpaths []config.GotoPath) {
	for i, gpath := range gpaths {
		fmt.Printf("%v - %s\n", i, gpath.String())
	}
}

func listReverse(gpaths []config.GotoPath) {
	for i := range gpaths {
		fmt.Printf("%v - %s\n", len(gpaths)-i-1, gpaths[len(gpaths)-i-1].String())
	}
}

func init() {
	//Add this command to RootCommand
	rootCmd.AddCommand(listGPathCmd)

	//Flags
	listGPathCmd.Flags().BoolP("reverse", "R", false, "List the goto-paths in reverse")
}
