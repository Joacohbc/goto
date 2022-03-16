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

	Short: "List all goto-paths in the goto-paths file",

	Run: func(cmd *cobra.Command, args []string) {

		//Parse "all" Flags
		currentPath, err := cmd.Flags().GetBool("current")
		cobra.CheckErr(err)

		path, err := cmd.Flags().GetString("path")
		cobra.CheckErr(err)

		abbv, err := cmd.Flags().GetString("abbv")
		cobra.CheckErr(err)

		reverse, err := cmd.Flags().GetBool("reverse")
		cobra.CheckErr(err)

		//Initial the variables to use config
		config.GotoPathsFile = GotoPathsFile

		var gpaths []config.GotoPath
		cobra.CheckErr(config.LoadConfigFile(&gpaths))

		//If CurrentPath is passed, the path to add is current directory
		if currentPath {
			//Get the current path
			currentDir, err := os.Getwd()
			cobra.CheckErr(err)

			//Valid the path
			cobra.CheckErr(config.ValidPath(&currentDir))

			//If all ok, overwrite "pathToAdd" variable
			path = currentDir
		}

		//If the flag "path" is passed
		if path != "" {
			//Valid the path
			cobra.CheckErr(config.ValidPath(&path))

			for i, gpath := range gpaths {
				if gpath.Path == path {
					fmt.Printf("%v - %s\n", i, gpath.String())
					return
				}

				if i == len(gpaths)-1 {
					fmt.Printf("The path \"%s\" doesn't exist in the gpaths-file\n", path)
				}
			}
			return
		}

		//If the flag "abbv" is passed
		if abbv != "" {
			cobra.CheckErr(config.ValidAbbreviation(&abbv))

			for i, gpath := range gpaths {
				if gpath.Abbreviation == abbv {
					fmt.Printf("%v - %s\n", i, gpath.String())
					return
				}

				if i == len(gpaths)-1 {
					fmt.Printf("The path \"%s\" doesn't exist in the gpaths-file\n", path)
				}
			}
			return
		}

		//If the flag "reverse" is passed
		if reverse {
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
