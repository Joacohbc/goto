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
	"path/filepath"

	"github.com/spf13/cobra"
)

var modesToUpdate []string = []string{
	"path-path",
	"path-abbv",
	"abbv-path",
	"abbv-abbv",
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update-path",
	Aliases: []string{"upd", "modify-path", "mod"},
	Short:   "Modify a path from goto-path file",
	Long: `
To use the update-path command you need have 4 modes to update: 
- A "Path" and a new "Path"(path-path)
- A "Path" and a new "Abbreviation" (path-abbv)
- A "Abbreviation" and a new "Path" (abbv-path)
- A "Abbreviation" and a new "Abbreviation" (abbv-path)

To update a path from goto-path file
`,

	Example: `
# If you wanto to update need to select a modes and pass
# two args

# Update the home of the user
goto update-path --mode path-path --path /home/myuser --new /home/mynewuser

# "h" the default abbreviation to home directory
goto update-path --mode abbv-path --abbv h --new /home/mynewuser

# Change the abbreviation of the come
goto update-path --mode path-abbv --path /home/myuser --new home

# Or
goto update-path --mode abbv-abbv --abbv h --new home
`,

	Run: func(cmd *cobra.Command, args []string) {

		var gpaths []config.GotoPath
		currentPath, err := cmd.Flags().GetBool("current")
		cobra.CheckErr(err)

		pathToUpd, err := cmd.Flags().GetString("path")
		cobra.CheckErr(err)

		abbvToUpd, err := cmd.Flags().GetString("abbv")
		cobra.CheckErr(err)

		//Initial the variables to use config package
		config.GotoPathsFile = GotoPathsFile
		config.ConfigDir = ConfigDir

		//Load the goto-paths file to array
		cobra.CheckErr(config.LoadConfigFile(&gpaths))

		//If CurrentPath is passed, the path to add is current directory
		if currentPath {
			currentDir, err := os.Getwd()
			cobra.CheckErr(err)

			absoluteDir, err := filepath.Abs(currentDir)
			cobra.CheckErr(err)

			pathToUpd = absoluteDir
		}

		mode, err := cmd.Flags().GetString("mode")
		cobra.CheckErr(err)

		new, err := cmd.Flags().GetString("new")
		cobra.CheckErr(err)

		switch mode {

		//path-path
		case modesToUpdate[0]:
			//Valid the Path and the new Path
			cobra.CheckErr(config.ValidPath(&pathToUpd))
			cobra.CheckErr(config.ValidPath(&new))

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Path == pathToUpd {
					gpaths[i].Path = new
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Path \"%v\" doesn't exist in the goto-paths file", pathToUpd))
				}
			}

		//path-abbv
		case modesToUpdate[1]:
			//Valid the Path and the new Abbreviation
			cobra.CheckErr(config.ValidPath(&pathToUpd))
			cobra.CheckErr(config.ValidAbbreviation(&new))

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Path == pathToUpd {
					gpaths[i].Abbreviation = new
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Path \"%v\" doesn't exist in the goto-paths file", pathToUpd))
				}
			}

		//abbv-path
		case modesToUpdate[2]:
			//Valid the Abbreviation and the new Path
			cobra.CheckErr(config.ValidAbbreviation(&abbvToUpd))
			cobra.CheckErr(config.ValidPath(&new))

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Abbreviation == abbvToUpd {
					gpaths[i].Path = new
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Abbreviation \"%v\" doesn't exist in the goto-paths file", abbvToUpd))
				}
			}

		//abbv-abbv
		case modesToUpdate[3]:
			//Valid the Abbreviation and the new Path
			cobra.CheckErr(config.ValidAbbreviation(&abbvToUpd))
			cobra.CheckErr(config.ValidAbbreviation(&new))

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Abbreviation == abbvToUpd {
					gpaths[i].Abbreviation = new
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Abbreviation \"%v\" doesn't exist in the goto-paths file", abbvToUpd))
				}
			}

		default:
			cobra.CheckErr(fmt.Errorf("the valid values to flag mode are: %v", modesToUpdate))
		}

		//If the array is valid, apply the changes
		cobra.CheckErr(config.CreateJsonFile(gpaths))

		fmt.Println("Changes applied successfully")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringP("path", "p", "", "The Path to delete")
	updateCmd.Flags().StringP("abbv", "a", "", "The Abbreviation of the Path")
	updateCmd.Flags().StringP("new", "n", "", "The Path or Abbreviation new")

	updateCmd.Flags().StringP("mode", "m", "", "Indicate that update in the format")

	updateCmd.Flags().BoolP("current", "c", false, "The Path to remove will be the current directory")
}
