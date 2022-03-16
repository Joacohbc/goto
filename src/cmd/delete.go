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

// deleteCmd represents the addGPath command
var deleteCmd = &cobra.Command{
	Use:     "delete-path",
	Aliases: []string{"del", "delete", "remove-path", "rem", "remove"},
	Short:   "Delete a path from goto-path file",
	Long: `
To use the delete-path command you need to pass two args: a "Path" and an "Abbreviation" to 
create a new goto-path`,

	Example: `
# This command add the current directory(the "Path") to the config file with
# the abbreviation "currentDir"	
goto add --current -abbv currentDir

# To specify the "Path" and "Abbreviation" use:
goto add --path ~/Documentos -abbv docs
`,

	Run: func(cmd *cobra.Command, args []string) {

		var gpathToDel config.GotoPath
		var gpaths []config.GotoPath
		{
			//Parse Flags//
			currentPath, err := cmd.Flags().GetBool("current")
			cobra.CheckErr(err)

			pathToDel, err := cmd.Flags().GetString("path")
			cobra.CheckErr(err)

			abbvToDel, err := cmd.Flags().GetString("abbv")
			cobra.CheckErr(err)

			//Initial the variables to use config package
			config.GotoPathsFile = GotoPathsFile
			config.ConfigDir = ConfigDir

			//Load the goto-paths file to array
			cobra.CheckErr(config.LoadConfigFile(&gpaths))

			//If CurrentPath is passed, the path to add is current directory
			if currentPath {
				//Get the current path
				currentDir, err := os.Getwd()
				cobra.CheckErr(err)

				//Valid the path
				cobra.CheckErr(config.ValidPath(&currentDir))

				//If all ok, overwrite "pathToDel" variable
				pathToDel = currentDir
			}

			//Create and valid the GPath
			gpathToDel = config.GotoPath{
				Path:         pathToDel,
				Abbreviation: abbvToDel,
			}
			cobra.CheckErr(gpathToDel.Valid())
		}

		//Delete the directory from the array
		for i, gpath := range gpaths {

			//The gpath passes have the same Path or the same Abbreviation, delete it
			if gpath.Path == gpathToDel.Path || gpath.Abbreviation == gpathToDel.Abbreviation {
				gpaths = append(gpaths[:i], gpaths[i+1:]...)
				break
			}

			if i == len(gpaths)-1 {
				cobra.CheckErr(fmt.Errorf("the path \"%v\" or abbreviation \"%s\" doesn't exist in the goto-paths file", gpathToDel.Path, gpathToDel.Abbreviation))
			}
		}
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
	deleteCmd.Flags().BoolP("current", "c", false, "The Path to remove will be the current directory (\"path\" flag will be overwrite)")
}
