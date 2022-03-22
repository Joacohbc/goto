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
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the addGPath command
var addCmd = &cobra.Command{
	Use:     "add-path",
	Aliases: []string{"add", "create-path", "create"},
	Short:   "Add a new path to goto-paths file",
	Long: `
To use the add-path command you need to pass two args: a "Path" and an "Abbreviation" to 
create a new goto-path`,

	Example: `
# This command add the current directory(the "Path") to the gpaths file with
# the abbreviation "currentDir"	
goto add-path --current -abbv currentDir

# To specify the "Path" and "Abbreviation" use:
goto add-path --path ~/Documents -abbv docs
`,

	Run: func(cmd *cobra.Command, args []string) {

		//Where load all gpaths
		var gpaths []config.GotoPath

		//The gpath to add
		var gpathToAdd config.GotoPath
		{
			//Create the GPath to add
			pathToAdd, err := cmd.Flags().GetString("path")
			cobra.CheckErr(err)

			abbvToAdd, err := cmd.Flags().GetString("abbv")
			cobra.CheckErr(err)

			//Initialze the variables to use config package
			config.GotoPathsFile = GotoPathsFile

			//Load the goto-paths file to array
			cobra.CheckErr(config.LoadConfigFile(&gpaths))

			//If CurrentPath is passed, the path to add is current directory
			if cmd.Flags().Changed("current") {
				//Get the current path
				currentDir, err := os.Getwd()
				cobra.CheckErr(err)

				//Valid the path
				cobra.CheckErr(config.ValidPath(&currentDir))

				//If all ok, overwrite "pathToAdd" variable
				pathToAdd = currentDir
			}

			//Create and valid the GPath
			gpathToAdd = config.GotoPath{
				Path:         pathToAdd,
				Abbreviation: abbvToAdd,
			}
			cobra.CheckErr(gpathToAdd.Valid())
		}

		//Add the new directory to the array and valid it
		gpaths = append(gpaths, gpathToAdd)

		//If the array is valid, apply the changes
		cobra.CheckErr(config.CreateJsonFile(gpaths))

		fmt.Println("Changes applied successfully")
	},
}

func init() {
	//Add this command to RootCmd
	rootCmd.AddCommand(addCmd)

	//Flags
	addCmd.Flags().StringP("path", "p", "", "The Path to add")
	addCmd.Flags().StringP("abbv", "a", "", "The Abbreviation of the Path")
	addCmd.Flags().BoolP("current", "c", false, "The Path to add will be the current directory (\"path\" flag will be overwrite)")
}
