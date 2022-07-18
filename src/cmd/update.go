package cmd

import (
	"fmt"
	"goto/src/gpath"
	"goto/src/utils"
	"strconv"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update-path",
	Aliases: []string{"upd", "update", "modify-path", "mod"},
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
# Format: goto update-path [ -t ] mode { -p path | -a abbreviation | -i index } { -p path | -a abbreviation | -i index } 

# Update the home of the user
goto update-path path-path --path /home/myuser --new /home/mynewuser

# "h" the default abbreviation to home directory
goto update-path abbv-path --abbv h --new /home/mynewuser

# Change the abbreviation of the come
goto update-path path-abbv --path /home/myuser --new home

# Or
goto update-path abbv-abbv --abbv h --new home
`,
	Args: cobra.RangeArgs(0, 1),
	PreRun: func(cmd *cobra.Command, args []string) {

		// If no arguments are passed and neither the modes flag is passed, return a error.
		if len(args) == 0 && !utils.FlagPassed(cmd, "modes") {
			cobra.CheckErr("must be specify a mode to update")
		}

		// If no value for new flags is passed, return a error
		if !utils.FlagPassed(cmd, "new") && !utils.FlagPassed(cmd, "new-current") {
			cobra.CheckErr("must be specify the new filed to update (path/abbreviation/index)")
		}

	},
	Run: func(cmd *cobra.Command, args []string) {

		modes := [][]string{
			{"path-path", "pp"}, // 0
			{"path-abbv", "pa"}, // 1
			{"path-indx", "pi"}, // 2
			{"abbv-path", "ap"}, // 3
			{"abbv-abbv", "aa"}, // 4
			{"abbv-indx", "ai"}, // 5
			{"indx-path", "ip"}, // 6
			{"indx-abbv", "ia"}, // 7
			{"indx-indx", "ii"}, // 8
		}

		//If modes is passed, show all modes
		if utils.FlagPassed(cmd, "modes") {
			for i := range modes {
				fmt.Println("Long form:", modes[i][0], "|", "Short form:", modes[i][1])
			}
			return
		}

		//Parse the new flag
		new, err := cmd.Flags().GetString("new")
		cobra.CheckErr(err)

		//If new-current is passed, overwrite the "new" to current directory
		if utils.FlagPassed(cmd, "new-current") {
			new = utils.GetCurrentDirectory()
		}

		//Load the goto-paths file to array
		gpaths := utils.LoadGPaths(cmd)

		// Change the GPath Index 1 for GPath in Index 2 and vice-versa
		changeIndex := func(inx1, inx2 int) {
			gpaths[inx1], gpaths[inx2] = gpaths[inx2], gpaths[inx1]
		}

		//Arg 0 indicate the Mode of the update
		switch args[0] {

		//path-path
		case modes[0][0], modes[0][1]:
			//Valid the Path and the new Path
			path := utils.GetPath(cmd)
			cobra.CheckErr(gpath.ValidPathVar(&new))

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
		case modes[1][0], modes[1][1]:
			//Valid the Path and the new Abbreviation
			path := utils.GetPath(cmd)
			cobra.CheckErr(gpath.ValidAbbreviationVar(&new))

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
		case modes[2][0], modes[2][1]:
			//Valid the Path and the new Abbreviation
			path := utils.GetPath(cmd)
			cobra.CheckErr(gpath.IsValidIndex(len(gpaths), new))

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
		case modes[3][0], modes[3][1]:
			//Valid the Abbreviation and the new Path
			abbv := utils.GetAbbreviation(cmd)
			cobra.CheckErr(gpath.ValidPathVar(&new))

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
		case modes[4][0], modes[4][1]:
			//Valid the Abbreviation and the new Path
			abbv := utils.GetAbbreviation(cmd)
			cobra.CheckErr(gpath.ValidAbbreviationVar(&new))

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
		case modes[5][0], modes[5][1]:
			//Valid the Path and the new Abbreviation
			abbv := utils.GetAbbreviation(cmd)
			cobra.CheckErr(gpath.IsValidIndex(len(gpaths), new))

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
		case modes[6][0], modes[6][1]:
			indx := utils.GetIndex(cmd)
			cobra.CheckErr(gpath.ValidPathVar(&new))

			for i := range gpaths {
				if i == indx {
					gpaths[indx].Path = new
					break
				}
			}

		//indx-abbv
		case modes[7][0], modes[7][1]:
			indx := utils.GetIndex(cmd)
			cobra.CheckErr(gpath.ValidAbbreviationVar(&new))

			for i := range gpaths {
				if i == indx {
					gpaths[indx].Abbreviation = new
					break
				}
			}

		//indx-indx
		case modes[8][0], modes[8][1]:
			indx := utils.GetIndex(cmd)
			cobra.CheckErr(gpath.IsValidIndex(len(gpaths), new))

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

		utils.UpdateGPaths(cmd, gpaths)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	//Flags//

	//Flags "To Update"
	updateCmd.Flags().StringP(utils.FlagPath, "p", "", "The Path to delete")
	updateCmd.Flags().BoolP(utils.FlagCurrentDir, "c", false, "The Path to update will be the current directory (\"path\" flag value will be overwrite)")
	updateCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "The Abbreviation of the Path")
	updateCmd.Flags().IntP(utils.FlagIndex, "i", -1, "The Index of the Path")

	//Flags "Update To"
	updateCmd.Flags().StringP("new", "n", "", "The Path or Abbreviation new")
	updateCmd.Flags().BoolP("new-current", "C", false, "The new Path will be the current directory (\"new\" flag value will be overwrite)")

	//Flag info
	updateCmd.Flags().BoolP("modes", "m", false, "Print all modes formats")
}
