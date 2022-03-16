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
	"strconv"

	"github.com/spf13/cobra"
)

var modesToUpdate []string = []string{
	"path-path", // 0
	"path-abbv", // 1
	"path-indx", // 2
	"abbv-path", // 3
	"abbv-abbv", // 4
	"abbv-indx", // 5
	"indx-path", // 6
	"indx-abbv", // 7
	"indx-indx", // 8
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update-path",
	Aliases: []string{"upd", "modify-path", "mod"},
	Short:   "Modify a path from goto-path file",
	Long: `
To use the update-path command you need have 9 modes to update: 
- A "Path" and a new "Path"(path-path)
- A "Path" and a new "Abbreviation" (path-abbv)
- A "Path" and a new "Indx" (path-indx)
- A "Abbreviation" and a new "Path" (abbv-path)
- A "Abbreviation" and a new "Abbreviation" (abbv-path)
- A "Abbreviation" and a new "Indx" (abbv-indx)
- A "Index" and a new "Path" (indx-path)
- A "Index" and a new "Abbreviation" (indx-abbv)
- A "Index" and a new "Index" (indx-indx)

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

		//"Parse" all Flags
		currentPath, err := cmd.Flags().GetBool("current")
		cobra.CheckErr(err)

		newCurrentPath, err := cmd.Flags().GetBool("new-current")
		cobra.CheckErr(err)

		pathToUpd, err := cmd.Flags().GetString("path")
		cobra.CheckErr(err)

		abbvToUpd, err := cmd.Flags().GetString("abbv")
		cobra.CheckErr(err)

		indxToUpd, err := cmd.Flags().GetInt("indx")
		cobra.CheckErr(err)

		mode, err := cmd.Flags().GetString("mode")
		cobra.CheckErr(err)

		new, err := cmd.Flags().GetString("new")
		cobra.CheckErr(err)

		//If CurrentPath is passed, overwrite the path to current directory
		if currentPath {
			currentDir, err := os.Getwd()
			cobra.CheckErr(err)

			absoluteDir, err := filepath.Abs(currentDir)
			cobra.CheckErr(err)

			pathToUpd = absoluteDir
		}

		//If newCurrentPath is passed, overwrite the "new" to current directory
		if newCurrentPath {
			currentDir, err := os.Getwd()
			cobra.CheckErr(err)

			absoluteDir, err := filepath.Abs(currentDir)
			cobra.CheckErr(err)

			new = absoluteDir
		}

		//Initial the variables to use config package
		config.GotoPathsFile = GotoPathsFile

		//Load the goto-paths file to array
		var gpaths []config.GotoPath
		cobra.CheckErr(config.LoadConfigFile(&gpaths))

		// Change the GPath Index 1 for GPath in Index 2 and vice-versa
		changeIndex := func(inx1, inx2 int) {
			gpath1 := gpaths[inx1]
			gpath2 := gpaths[inx2]
			gpaths[inx1] = gpath2
			gpaths[inx2] = gpath1
		}

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

		//path-indx
		case modesToUpdate[2]:
			//Valid the Path and the new Abbreviation
			cobra.CheckErr(config.ValidPath(&pathToUpd))
			cobra.CheckErr(config.IsValidIndex(gpaths, new))

			n, _ := strconv.Atoi(new)

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Path == pathToUpd {
					changeIndex(i, n)
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Path \"%v\" doesn't exist in the goto-paths file", pathToUpd))
				}
			}

		//abbv-path
		case modesToUpdate[3]:
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
		case modesToUpdate[4]:
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

		//abbv-indx
		case modesToUpdate[5]:
			//Valid the Path and the new Abbreviation
			cobra.CheckErr(config.ValidAbbreviation(&abbvToUpd))
			cobra.CheckErr(config.IsValidIndex(gpaths, new))

			n, _ := strconv.Atoi(new)

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Abbreviation == abbvToUpd {
					changeIndex(i, n)
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Path \"%v\" doesn't exist in the goto-paths file", pathToUpd))
				}
			}

		//indx-path
		case modesToUpdate[6]:
			cobra.CheckErr(config.IsValidIndex(gpaths, strconv.Itoa(indxToUpd)))
			cobra.CheckErr(config.ValidPath(&new))

			for i := range gpaths {
				if i == indxToUpd {
					gpaths[indxToUpd].Path = new
					break
				}
			}

		//indx-abbv
		case modesToUpdate[7]:
			cobra.CheckErr(config.IsValidIndex(gpaths, strconv.Itoa(indxToUpd)))
			cobra.CheckErr(config.ValidAbbreviation(&new))

			for i := range gpaths {
				if i == indxToUpd {
					gpaths[indxToUpd].Abbreviation = new
					break
				}
			}

		//indx-indx
		case modesToUpdate[8]:
			cobra.CheckErr(config.IsValidIndex(gpaths, strconv.Itoa(indxToUpd)))
			cobra.CheckErr(config.IsValidIndex(gpaths, new))

			n, _ := strconv.Atoi(new)

			for i := range gpaths {
				if i == indxToUpd {
					changeIndex(i, n)
					break
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
	updateCmd.Flags().IntP("indx", "i", -1, "The Index of the Path")

	updateCmd.Flags().StringP("new", "n", "", "The Path or Abbreviation new")
	updateCmd.Flags().BoolP("new-current", "C", false, "The new Path will be the current directory (\"new\" flag will be overwrite)")

	updateCmd.Flags().StringP("mode", "m", "", "Indicate that update in the format")
	updateCmd.Flags().BoolP("current", "c", false, "The Path to update will be the current directory (\"path\" parameter will be overwrite)")
}
