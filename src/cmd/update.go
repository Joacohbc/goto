package cmd

import (
	"fmt"
	"goto/src/core"
	"goto/src/utils"

	"github.com/spf13/cobra"
)

// UpdateCmd represents the update command
var UpdateCmd = &cobra.Command{
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

# The same that:
goto update ap --abbv h --new /home/mynewuser

# Change the abbreviation of the come
goto update-path path-abbv --path /home/myuser --new home

# Or if you want to update the abbreviation of the home
goto update abbv-abbv --abbv h --new home
`,
	Args:   cobra.RangeArgs(0, 1),
	PreRun: preRunUpdate,
	Run:    runUpdate,
}

func preRunUpdate(cmd *cobra.Command, args []string) {

	// If no arguments are passed and neither the modes flag is passed, return a error.
	if len(args) == 0 && !utils.FlagPassed(cmd, "modes") {
		cobra.CheckErr("must be specify a mode to update")
	}

	// If no value for new flags is passed, return a error
	if !utils.FlagPassed(cmd, "new") {
		cobra.CheckErr("must be specify the new filed to update (path/abbreviation/index)")
	}

}

func runUpdate(cmd *cobra.Command, args []string) {

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
	newVal, err := cmd.Flags().GetString("new")
	cobra.CheckErr(err)

	path, _ := cmd.Flags().GetString(utils.FlagPath)
	abbv, _ := cmd.Flags().GetString(utils.FlagAbbreviation)
	indx, _ := cmd.Flags().GetInt(utils.FlagIndex)

	cobra.CheckErr(core.UpdatePath(args[0], path, abbv, indx, newVal, utils.TemporalFlagPassed(cmd)))
}

func init() {
	RootCmd.AddCommand(UpdateCmd)

	//Flags//

	//Flags "To Update"
	UpdateCmd.Flags().StringP(utils.FlagPath, "p", "", "The Path to delete")
	UpdateCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "The Abbreviation of the Path")
	UpdateCmd.Flags().IntP(utils.FlagIndex, "i", -1, "The Index of the Path")

	//Flags "Update To"
	UpdateCmd.Flags().StringP("new", "n", "", "The Path or Abbreviation new")

	//Flag info
	UpdateCmd.Flags().BoolP("modes", "m", false, "Print all modes formats")
}
