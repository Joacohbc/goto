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
	"goto/src/gpath"
	"goto/src/utils"

	"github.com/spf13/cobra"
)

// deleteCmd represents the addGPath command
var deleteCmd = &cobra.Command{
	Use:     "delete-path",
	Aliases: []string{"del", "delete", "remove-path", "rem", "remove"},
	Short:   "Delete a path from goto-path file",
	Args:    cobra.ExactArgs(0),
	Long: `
To use the delete-path command you need to Path, Abbreviation or Index from a goto-path`,

	Example: `
# This command delete-path the current directory(the "Path") to the gpaths file
goto delete-path --current

# To specify the "Path", "Abbreviation" or Index. use:

# Delete the gpath with the path "/home/user/Documents"
goto delete-path --path ~/Documents

# Delete the gpath with the abbreviation "docs"
goto delete-path --abbv docs

# Delete the gpath in the index "2"
goto delete-path --indx 2
`,
	PreRun: func(cmd *cobra.Command, args []string) {

		passed := func(flag string) bool { return cmd.Flags().Changed(flag) }

		//If any flags is passed or more than 2, return an error
		if cmd.Flags().NFlag() == 0 || cmd.Flags().NFlag() > 2 {
			cobra.CheckErr(fmt.Errorf("you must specify one flag to delete a gpath (Or Path or Abbreviation or Index)"))
		}

		//If one flags is passed and it is the temporary flags, return an error
		if cmd.Flags().NFlag() == 1 && passed("temporal") {
			cobra.CheckErr(fmt.Errorf("you must specify one flag to delete a gpath (Or Path or Abbreviation or Index)"))
		}

		//If 2 flags are passed and the temporary flag is not
		//passed it means that two flags were passed to identify the gpath
		/*
			2 flags to identify the gpath may cause an error to delete the path.
			For example, -p /home/user -i 2
			the index not match with the gpath, so delete one of the paths
		*/
		if cmd.Flags().NFlag() == 2 && !passed("temporal") {
			cobra.CheckErr(fmt.Errorf("you must specify one flag to delete a gpath (Or Path or Abbreviation or Index)"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		//Load the goto-paths file to array
		gpaths := utils.LoadGPaths(cmd)

		if utils.FlagPassed(cmd, utils.FlagPath) {
			path := utils.GetPath(cmd)

			//Delete the directory from the array
			for i, gpath := range gpaths {

				//The gpath passes have the same Path or the same Abbreviation, delete it
				if gpath.Path == path {
					gpaths = append(gpaths[:i], gpaths[i+1:]...)
					fmt.Printf("The path %s (%s) was deleted\n", gpath.Path, gpath.String())
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("any gpath match with the path \"%s\"", path))
				}
			}
		}

		if utils.FlagPassed(cmd, utils.FlagAbbreviation) {
			abbv := utils.GetAbbreviation(cmd)

			//Delete the directory from the array
			for i, gpath := range gpaths {

				//The gpath passes have the same Path or the same Abbreviation, delete it
				if gpath.Abbreviation == abbv {
					gpaths = append(gpaths[:i], gpaths[i+1:]...)
					fmt.Printf("The path %s (%s) was deleted\n", gpath.Path, gpath.String())
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("any gpath match with the abbreviation \"%s\"", abbv))
				}
			}
		}

		if utils.FlagPassed(cmd, utils.FlagIndex) {
			indx := utils.GetIndex(cmd)

			//Delete the directory from the array
			for i, gpath := range gpaths {

				//The gpath passes have the same Path or the same Abbreviation, delete it
				if i == indx {
					gpaths = append(gpaths[:i], gpaths[i+1:]...)
					fmt.Printf("The path %s (%s) was deleted\n", gpath.Path, gpath.String())
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("any gpath match with the index %d", i))
				}
			}
		}

		//After the changes, valid it
		cobra.CheckErr(gpath.ValidArray(gpaths))
		utils.UpdateGPaths(cmd, gpaths)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	//Flags
	deleteCmd.Flags().StringP(utils.FlagPath, "p", "", "The Path to delete")
	deleteCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "The Abbreviation of the Path")
	deleteCmd.Flags().IntP(utils.FlagIndex, "i", -1, "The Index of the Path")
}
