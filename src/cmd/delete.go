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
	"strconv"

	"github.com/spf13/cobra"
)

// deleteCmd represents the addGPath command
var deleteCmd = &cobra.Command{
	Use:     "delete-path",
	Aliases: []string{"del", "delete", "remove-path", "rem", "remove"},
	Short:   "Delete a path from goto-path file",
	Long: `
To use the delete-path command you need to pass two args: a "Path" and an "Abbreviation" to 
create a new goto-path`,

	Example: `
# This command del the current directory(the "Path") to the gpaths file
goto del -p ~/Documents

# To specify the "Path" and "Abbreviation" use:
goto del --path ~/Documentos -abbv docs
`,

	Run: func(cmd *cobra.Command, args []string) {

		passed := func(flag string) bool { return cmd.Flags().Changed(flag) }

		//Where any error is saved
		var err error = nil

		//Load the goto-paths file to array
		var gpaths []config.GotoPath
		{
			//Initial the variables to use config package
			config.GotoPathsFile = GotoPathsFile

			cobra.CheckErr(config.LoadConfigFile(&gpaths))
		}

		//Before load the flags into de vars check valid number of Flags
		//If the any flag or more than 1 is passed return a error
		if cmd.Flags().NFlag() == 0 || cmd.Flags().NFlag() > 1 {
			cobra.CheckErr(fmt.Errorf("you must specify 1 (only one) flag (Or Path or Abbreviation or Index) to `delete a gpath"))
		}

		//Parse all Flags

		var pathToDel string
		//If the path flag is passed, load the pathToDel
		if passed("path") {
			pathToDel, err = cmd.Flags().GetString("path")
			cobra.CheckErr(err)
			config.ValidPath(&pathToDel)
		}

		//If current flag is passed, overwrite pathToDel with the current directory
		if passed("current") {
			//Get the current path
			pathToDel, err = os.Getwd()
			cobra.CheckErr(err)

			//Valid the path
			cobra.CheckErr(config.ValidPath(&pathToDel))
		}

		var abbvToDel string = ""
		//If the abbv flag is passed, load the abbvToDel
		if passed("abbv") {
			abbvToDel, err = cmd.Flags().GetString("abbv")
			cobra.CheckErr(err)
			config.ValidAbbreviation(&abbvToDel)
		}

		var indxToDel int
		//If the indx flag is passed, load the indxToDel
		if passed("indx") {
			indxToDel, err = cmd.Flags().GetInt("indx")
			cobra.CheckErr(err)
			config.IsValidIndex(gpaths, strconv.Itoa(indxToDel))
		}

		//Delete the directory from the array
		for i, gpath := range gpaths {

			//The gpath passes have the same Path or the same Abbreviation, delete it
			if gpath.Path == pathToDel || gpath.Abbreviation == abbvToDel || i == indxToDel {
				gpaths = append(gpaths[:i], gpaths[i+1:]...)
				break
			}

			if i == len(gpaths)-1 {
				cobra.CheckErr(fmt.Errorf("any gpath match with the flags you passed"))
			}
		}

		//After the changes, valid it
		cobra.CheckErr(config.ValidArray(gpaths))

		//If the array is valid, apply the changes
		cobra.CheckErr(config.CreateJsonFile(gpaths))

		fmt.Println("Changes applied successfully")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	//Flags
	deleteCmd.Flags().StringP("path", "p", "", "The Path to delete")
	deleteCmd.Flags().StringP("abbv", "a", "", "The Abbreviation of the Path")
	deleteCmd.Flags().IntP("indx", "i", -1, "The Index of the Path")
	deleteCmd.Flags().BoolP("current", "c", false, "The Path to remove will be the current directory (\"path\" flag will be overwrite)")
}
