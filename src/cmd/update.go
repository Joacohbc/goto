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
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update-path",
	Aliases: []string{"upd", "update", "modify-path", "mod"},
	Args:    cobra.ExactArgs(1),
	Short:   "Update a path from goto-path file",
	Long: `
To use the update-path command you have 9 modes to update, each mode needs two args, 
the first to identify the goto-path and the second specific to what is to be updated. 

Modes:
- A "Path" and a new "Path" (path-path)
- A "Path" and a new "Abbreviation" (path-abbv)
- A "Path" and a new "Indx" (path-indx)
- A "Abbreviation" and a new "Path" (abbv-path)
- A "Abbreviation" and a new "Abbreviation" (abbv-path)
- A "Abbreviation" and a new "Indx" (abbv-indx)
- A "Index" and a new "Path" (indx-path)
- A "Index" and a new "Abbreviation" (indx-abbv)
- A "Index" and a new "Index" (indx-indx)
`,

	Example: `
# Update the home of the user
goto update-path path-path --path /home/myuser --new /home/mynewuser

# "h" the default abbreviation to home directory
goto update-path abbv-path --abbv h --new /home/mynewuser

# Change the abbreviation of the come
goto update-path path-abbv --path /home/myuser --new home

# Or
goto update-path abbv-abbv --abbv h --new home
`,

	Run: func(cmd *cobra.Command, args []string) {

		modesToUpdate := []string{
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

		modesToUpdateShort := []string{
			"pp", // 0
			"pa", // 1
			"pi", // 2
			"ap", // 3
			"aa", // 4
			"ai", // 5
			"ip", // 6
			"ia", // 7
			"ii", // 8
		}

		//If modes is passed, show all modes
		if passed(cmd, "modes") {
			for i := range modesToUpdate {
				fmt.Println("Long form:", modesToUpdate[i], "|", "Short form:", modesToUpdateShort[i])
			}
			return
		}

		//Parse the new falg
		new, err := cmd.Flags().GetString("new")
		cobra.CheckErr(err)

		//If new-current is passed, overwrite the "new" to current directory
		if passed(cmd, "new-current") {
			new = getCurrentDirectory()
		}

		//Load the goto-paths file to array
		var gpaths []config.GotoPath
		loadGPath(cmd, &gpaths)

		// Change the GPath Index 1 for GPath in Index 2 and vice-versa
		changeIndex := func(inx1, inx2 int) {
			gpaths[inx1], gpaths[inx2] = gpaths[inx2], gpaths[inx1]
		}

		//Arg 0 indicate the Mode of the update
		switch args[0] {

		//path-path
		case modesToUpdate[0], modesToUpdateShort[0]:
			//Valid the Path and the new Path
			path := getPath(cmd)
			cobra.CheckErr(config.ValidPathVar(&new))

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Path == path {
					gpaths[i].Path = new
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Path \"%v\" doesn't exist in the goto-paths file", path))
				}
			}

		//path-abbv
		case modesToUpdate[1], modesToUpdateShort[1]:
			//Valid the Path and the new Abbreviation
			path := getPath(cmd)
			cobra.CheckErr(config.ValidAbbreviationVar(&new))

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Path == path {
					gpaths[i].Abbreviation = new
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Path \"%v\" doesn't exist in the goto-paths file", path))
				}
			}

		//path-indx
		case modesToUpdate[2], modesToUpdateShort[2]:
			//Valid the Path and the new Abbreviation
			path := getPath(cmd)
			cobra.CheckErr(config.IsValidIndex(gpaths, new))

			n, _ := strconv.Atoi(new)

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Path == path {
					changeIndex(i, n)
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Path \"%v\" doesn't exist in the goto-paths file", path))
				}
			}

		//abbv-path
		case modesToUpdate[3], modesToUpdateShort[3]:
			//Valid the Abbreviation and the new Path
			abbv := getAbbreviation(cmd)
			cobra.CheckErr(config.ValidPathVar(&new))

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Abbreviation == abbv {
					gpaths[i].Path = new
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Abbreviation \"%v\" doesn't exist in the goto-paths file", abbv))
				}
			}

		//abbv-abbv
		case modesToUpdate[4], modesToUpdateShort[4]:
			//Valid the Abbreviation and the new Path
			abbv := getAbbreviation(cmd)
			cobra.CheckErr(config.ValidAbbreviationVar(&new))

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Abbreviation == abbv {
					gpaths[i].Abbreviation = new
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Abbreviation \"%v\" doesn't exist in the goto-paths file", abbv))
				}
			}

		//abbv-indx
		case modesToUpdate[5], modesToUpdateShort[5]:
			//Valid the Path and the new Abbreviation
			abbv := getAbbreviation(cmd)
			cobra.CheckErr(config.IsValidIndex(gpaths, new))

			n, _ := strconv.Atoi(new)

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Abbreviation == abbv {
					changeIndex(i, n)
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Abbreviation \"%v\" doesn't exist in the goto-paths file", abbv))
				}
			}

		//indx-path
		case modesToUpdate[6], modesToUpdateShort[6]:
			indx := getIndex(cmd)
			cobra.CheckErr(config.ValidPathVar(&new))

			for i := range gpaths {
				if i == indx {
					gpaths[indx].Path = new
					break
				}
			}

		//indx-abbv
		case modesToUpdate[7], modesToUpdateShort[7]:
			indx := getIndex(cmd)
			cobra.CheckErr(config.ValidAbbreviationVar(&new))

			for i := range gpaths {
				if i == indx {
					gpaths[indx].Abbreviation = new
					break
				}
			}

		//indx-indx
		case modesToUpdate[8], modesToUpdateShort[8]:
			indx := getIndex(cmd)
			cobra.CheckErr(config.IsValidIndex(gpaths, new))

			n, _ := strconv.Atoi(new)

			for i := range gpaths {
				if i == indx {
					changeIndex(i, n)
					break
				}
			}

		default:
			cobra.CheckErr(fmt.Errorf("invalid values of modes to update, use goto --modes"))
		}

		createGPath(cmd, gpaths)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	//Flags//

	//Flags "To Update"
	updateCmd.Flags().StringP(FlagPath, "p", "", "The Path to delete")
	updateCmd.Flags().BoolP(FlagCurretDir, "c", false, "The Path to update will be the current directory (\"path\" flag will be overwrite)")
	updateCmd.Flags().StringP(FlagAbbreviation, "a", "", "The Abbreviation of the Path")
	updateCmd.Flags().IntP(FlagIndex, "i", -1, "The Index of the Path")

	//Flags "Update To"
	updateCmd.Flags().StringP("new", "n", "", "The Path or Abbreviation new")
	updateCmd.Flags().BoolP("new-current", "C", false, "The new Path will be the current directory (\"new\" flag will be overwrite)")

	//Flag info
	updateCmd.Flags().BoolP("modes", "m", false, "Print all modes formats")
}
