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
var listCmd = &cobra.Command{
	Use:     "list-path",
	Aliases: []string{"list"},

	Short: "List all gpaths in the gpaths file",

	Run: func(cmd *cobra.Command, args []string) {

		passed := func(flag string) bool { return cmd.Flags().Changed(flag) }

		//Load the goto-paths file to array
		var gpaths []config.GotoPath
		{
			//Initial the variables to use config
			config.GotoPathsFile = GotoPathsFile
			cobra.CheckErr(config.LoadConfigFile(&gpaths))
		}

		//Where any error is saved
		var err error = nil

		//Parse all Flags

		//If the flag "path" or "current" are passed
		if passed("path") || passed("current") {

			var path string = ""
			//If the path flag is passed, load the path
			if passed("path") {
				path, err = cmd.Flags().GetString("path")
				cobra.CheckErr(err)
			}

			//If current flag is passed, overwrite the path to current directory
			if passed("current") {
				//Get the current path, and overwrite "path" variable
				path, err = os.Getwd()
				cobra.CheckErr(err)
			}

			//Valid the path
			cobra.CheckErr(config.ValidPath(&path))

			for i, gpath := range gpaths {
				if gpath.Path == path {
					fmt.Printf("%v - %s\n", i, gpath.String())
					return
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the path \"%s\" doesn't exist in the gpaths-file", path))
				}
			}
			return
		}

		//If the flag "abbv" is passed
		if passed("abbv") {

			var abbv string
			abbv, err = cmd.Flags().GetString("abbv")
			cobra.CheckErr(err)

			cobra.CheckErr(config.ValidAbbreviation(&abbv))

			for i, gpath := range gpaths {
				if gpath.Abbreviation == abbv {
					fmt.Printf("%v - %s\n", i, gpath.String())
					return
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("doesn't exist a path with that abbreviation \"%s\"", abbv))
				}
			}
			return
		}

		//If the flag "reverse" is passed
		if passed("reverse") {
			for i := range gpaths {
				fmt.Printf("%v - %s\n", len(gpaths)-i-1, gpaths[len(gpaths)-i-1].String())
			}
			return
		}

		//If any flag is passed
		for i, gpath := range gpaths {
			fmt.Printf("%v - %s\n", i, gpath.String())
		}
	},
}

func init() {
	//Add this command to RootCommand
	rootCmd.AddCommand(listCmd)

	//Flags
	listCmd.Flags().StringP("path", "p", "", "The Path to delete")
	listCmd.Flags().StringP("abbv", "a", "", "The Abbreviation of the Path")
	listCmd.Flags().BoolP("current", "c", false, "The Path to update will be the current directory (\"path\" parameter will be overwrite)")
	listCmd.Flags().BoolP("reverse", "R", false, "List the goto-paths in reverse")
}
