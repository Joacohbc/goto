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
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ModesToUpdate []string = []string{
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

	ModesToUpdateShort []string = []string{
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
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update-path",
	Aliases: []string{"upd", "update", "modify-path", "mod"},
	Args:    cobra.MaximumNArgs(1),
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

	PreRun: func(cmd *cobra.Command, args []string) {
		//If modes is passed, show all modes and exit
		if cmd.Flags().Changed("modes") {
			for i := range ModesToUpdate {
				fmt.Println("Long form:", ModesToUpdate[i], "|", "Short form:", ModesToUpdateShort[i])
			}
			os.Exit(0)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {

		//Load the goto-paths file to array
		var gpaths []config.GotoPath
		LoadGPath(cmd, &gpaths)

		// Change the GPath Index 1 for GPath in Index 2 and vice-versa
		changeIndex := func(inx1, inx2 int) {
			gpath1 := gpaths[inx1]
			gpath2 := gpaths[inx2]
			gpaths[inx1] = gpath2
			gpaths[inx2] = gpath1
		}

		current := func(key, currentKey string) string {
			if cmd.Flags().Changed(currentKey) {
				return GetCurrentDirectory()
			}

			var path string = viper.GetString(key)
			cobra.CheckErr(config.ValidPathVar(&path))
			return path
		}

		//Arg 0 indicate the Mode of the update
		switch args[0] {

		//path-path
		case ModesToUpdate[0], ModesToUpdateShort[0]:
			path := current("path", "current")
			new := current("new", "new-current")

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
		case ModesToUpdate[1], ModesToUpdateShort[1]:
			//Valid the Path and the new Abbreviation
			path := current("path", "current")

			new := viper.GetString("new")
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
		case ModesToUpdate[2], ModesToUpdateShort[2]:
			path := current("path", "current")

			new := viper.GetInt("new")
			cobra.CheckErr(config.IsValidIndex(gpaths, strconv.Itoa(new)))

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Path == path {
					changeIndex(i, new)
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Path \"%v\" doesn't exist in the goto-paths file", path))
				}
			}

		//abbv-path
		case ModesToUpdate[3], ModesToUpdateShort[3]:
			abbv := viper.GetString("abbv")
			cobra.CheckErr(config.ValidAbbreviationVar(&abbv))

			new := current("new", "new-current")

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
		case ModesToUpdate[4], ModesToUpdateShort[4]:
			abbv := viper.GetString("abbv")
			cobra.CheckErr(config.ValidAbbreviationVar(&abbv))

			new := viper.GetString("new")
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
		case ModesToUpdate[5], ModesToUpdateShort[5]:
			abbv := viper.GetString("abbv")
			cobra.CheckErr(config.ValidAbbreviationVar(&abbv))

			new := viper.GetInt("new")
			cobra.CheckErr(config.IsValidIndex(gpaths, strconv.Itoa(new)))

			//And search in the array
			for i := range gpaths {
				if gpaths[i].Abbreviation == abbv {
					changeIndex(i, new)
					break
				}

				if i == len(gpaths)-1 {
					cobra.CheckErr(fmt.Errorf("the Abbreviation \"%v\" doesn't exist in the goto-paths file", abbv))
				}
			}

		//indx-path
		case ModesToUpdate[6], ModesToUpdateShort[6]:
			index := viper.GetInt("indx")
			cobra.CheckErr(config.IsValidIndex(gpaths, strconv.Itoa(index)))

			new := current("new", "new-current")

			for i := range gpaths {
				if i == index {
					gpaths[index].Path = new
					break
				}
			}

		//indx-abbv
		case ModesToUpdate[7], ModesToUpdateShort[7]:
			index := viper.GetInt("indx")
			cobra.CheckErr(config.IsValidIndex(gpaths, strconv.Itoa(index)))

			new := viper.GetString("new")
			cobra.CheckErr(config.ValidAbbreviationVar(&new))

			for i := range gpaths {
				if i == index {
					gpaths[index].Abbreviation = new
					break
				}
			}

		//indx-indx
		case ModesToUpdate[8], ModesToUpdateShort[8]:

			index := viper.GetInt("indx")
			cobra.CheckErr(config.IsValidIndex(gpaths, strconv.Itoa(index)))

			new := viper.GetInt("new")
			cobra.CheckErr(config.IsValidIndex(gpaths, strconv.Itoa(new)))

			for i := range gpaths {
				if i == index {
					changeIndex(i, new)
					break
				}
			}

		default:
			cobra.CheckErr(fmt.Errorf("invalid values of modes to update, use goto --modes"))
		}

		CreateGPath(cmd, gpaths)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	//Flags//

	//Flags "To Update"
	updateCmd.Flags().StringP("path", "p", "", "The Path to delete")
	viper.BindPFlag("path", updateCmd.Flags().Lookup("path"))

	updateCmd.Flags().BoolP("current", "c", false, "The Path to update will be the current directory (\"path\" flag will be overwrite)")
	viper.BindPFlag("current", updateCmd.Flags().Lookup("current"))

	updateCmd.Flags().StringP("abbv", "a", "", "The Abbreviation of the Path")
	viper.BindPFlag("abbv", updateCmd.Flags().Lookup("abbv"))

	updateCmd.Flags().IntP("indx", "i", -1, "The Index of the Path")
	viper.BindPFlag("indx", updateCmd.Flags().Lookup("indx"))

	//Flags "Update To"
	updateCmd.Flags().StringP("new", "n", "", "The Path or Abbreviation new")
	viper.BindPFlag("new", updateCmd.Flags().Lookup("new"))

	updateCmd.Flags().BoolP("new-current", "C", false, "The new Path will be the current directory (\"new\" flag will be overwrite)")
	viper.BindPFlag("new-current", updateCmd.Flags().Lookup("new-current"))

	//Flag info
	updateCmd.Flags().BoolP("modes", "m", false, "Print all modes formats")
}
