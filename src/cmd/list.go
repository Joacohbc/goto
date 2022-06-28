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

	"github.com/spf13/cobra"
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

		//Load the goto-paths file to array
		gpaths := loadGPath(cmd)

		// If the any path flag is passed
		if passed(cmd, FlagCurretDir) || passed(cmd, FlagPath) {

			path := getPath(cmd)

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

		// If the abbreviation flag is passed
		if passed(cmd, FlagAbbreviation) {
			abbv := getAbbreviation(cmd)
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
		if passed(cmd, "reverse") {
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
	listCmd.Flags().StringP(FlagPath, "p", "", "The Path to delete")
	listCmd.Flags().StringP(FlagAbbreviation, "a", "", "The Abbreviation of the Path")
	listCmd.Flags().BoolP(FlagCurretDir, "c", false, "The Path to update will be the current directory (\"path\" parameter will be overwrite)")
	listCmd.Flags().BoolP("reverse", "R", false, "List the goto-paths in reverse")
}
