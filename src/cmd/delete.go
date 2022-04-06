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
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		passed := func(flag string) bool { return cmd.Flags().Changed(flag) }

		//Load the goto-paths file to array
		var gpaths []config.GotoPath
		LoadGPath(cmd, &gpaths)

		//
		//Parse all Flags
		//

		var pathToDel string = viper.GetString("path")
		if passed("path") {
			cobra.CheckErr(config.ValidPathVar(&pathToDel))
		}

		//If current flag is passed, overwrite pathToDel with the current directory
		if passed("current") {
			pathToDel = GetCurrentDirectory()
		}

		var abbvToDel string = viper.GetString("abbv")
		if passed("abbv") {
			cobra.CheckErr(config.ValidAbbreviationVar(&abbvToDel))
		}

		var indxToDel int = viper.GetInt("indx")
		if passed("indx") {
			cobra.CheckErr(config.IsValidIndex(gpaths, strconv.Itoa(indxToDel)))
		}

		//Delete the directory from the array
		for i, gpath := range gpaths {

			//The gpath passes have the same Path or the same Abbreviation, delete it
			if gpath.Path == pathToDel || gpath.Abbreviation == abbvToDel || i == indxToDel {
				gpaths = append(gpaths[:i], gpaths[i+1:]...)
				fmt.Printf("The path %s (Index %v - Abbreviation %s) was deleted\n", gpath.Path, i, gpath.Abbreviation)
				break
			}

			if i == len(gpaths)-1 {
				cobra.CheckErr(fmt.Errorf("any gpath match with the flags you passed"))
			}
		}

		//After the changes, valid it
		cobra.CheckErr(config.ValidArray(gpaths))

		CreateGPath(cmd, gpaths)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	//Flags
	deleteCmd.Flags().StringP("path", "p", "", "The Path to delete") // A intial value that can't match with any gpath
	viper.BindPFlag("path", deleteCmd.Flags().Lookup("path"))

	deleteCmd.Flags().StringP("abbv", "a", "", "The Abbreviation of the Path") // A intial value that can't match with any gpath
	viper.BindPFlag("abbv", deleteCmd.Flags().Lookup("abbv"))

	deleteCmd.Flags().IntP("indx", "i", -1, "The Index of the Path") // A intial value that can't match with any gpath
	viper.BindPFlag("indx", deleteCmd.Flags().Lookup("indx"))

	deleteCmd.Flags().BoolP("current", "c", false, "The Path to add will be the current directory (\"path\" flag will be overwrite)")
}
