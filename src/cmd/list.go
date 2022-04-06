/*
Copyright © 2022 Joacohbc <joacog48@gmail.com>

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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listGPathCmd represents the listGPath command
var listCmd = &cobra.Command{
	Use:     "list-path",
	Aliases: []string{"list"},
	Args:    cobra.ExactArgs(0),
	Short:   "List goto-paths in the goto-paths file",
	Example: `
#To list all goto-paths
goto list

#To list a specific goto-path you can use the Path or the Abbreviation 
goto list --path ~/Documents
goto list --abbv docs
`,

	Run: func(cmd *cobra.Command, args []string) {

		passed := func(flag string) bool { return cmd.Flags().Changed(flag) }

		//Load the goto-paths file to array
		var gpaths []config.GotoPath
		LoadGPath(cmd, &gpaths)

		//If the flag "path" or "current" are passed
		if passed("path") || passed("current") {

			path := viper.GetString("path")

			//If current flag is passed, overwrite the path to current directory
			if passed("current") {
				path = GetCurrentDirectory()
			}

			//Valid the path
			cobra.CheckErr(config.ValidPathVar(&path))

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

			abbv := viper.GetString("abbv")
			cobra.CheckErr(config.ValidAbbreviationVar(&abbv))

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
