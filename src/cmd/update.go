package cmd

import (
	"fmt"
	"goto/src/core"
	"goto/src/utils"
	"strings"

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

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		modes := []string{
			"path-path\tUpdate a Path with a new Path",
			"pp\tUpdate a Path with a new Path (short)",
			"path-abbv\tUpdate a Path with a new Abbreviation",
			"pa\tUpdate a Path with a new Abbreviation (short)",
			"path-indx\tUpdate a Path with a new Index",
			"pi\tUpdate a Path with a new Index (short)",
			"abbv-path\tUpdate an Abbreviation with a new Path",
			"ap\tUpdate an Abbreviation with a new Path (short)",
			"abbv-abbv\tUpdate an Abbreviation with a new Abbreviation",
			"aa\tUpdate an Abbreviation with a new Abbreviation (short)",
			"abbv-indx\tUpdate an Abbreviation with a new Index",
			"ai\tUpdate an Abbreviation with a new Index (short)",
			"indx-path\tUpdate an Index with a new Path",
			"ip\tUpdate an Index with a new Path (short)",
			"indx-abbv\tUpdate an Index with a new Abbreviation",
			"ia\tUpdate an Index with a new Abbreviation (short)",
			"indx-indx\tUpdate an Index with a new Index",
			"ii\tUpdate an Index with a new Index (short)",
		}
		return modes, cobra.ShellCompDirectiveNoFileComp
	},

	Run:    runUpdate,
}

func preRunUpdate(cmd *cobra.Command, args []string) {

	// Either positional arg or --modes flag must be passed
	modesFlag, _ := cmd.Flags().GetString("modes")
	if len(args) == 0 && modesFlag == "" {
		cobra.CheckErr("must specify a mode to update")
	}

	// If no value for new flags is passed, and we are not listing/showing modes, return a error
	if modesFlag != "show" && modesFlag != "list" && !utils.FlagPassed(cmd, "new") {
		cobra.CheckErr("must specify the new field to update (path/abbreviation/index)")
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

	// Resolve the mode: either from the `--modes`/`-m` flag or from the first positional argument `args[0]`
	var modeVal string
	modesFlag, _ := cmd.Flags().GetString("modes")
	if modesFlag != "" {
		modeVal = modesFlag
	} else if len(args) > 0 {
		modeVal = args[0]
	}

	// If modes is passed as "show" or "list", show all modes
	if modeVal == "show" || modeVal == "list" {
		for i := range modes {
			fmt.Println("Long form:", modes[i][0], "|", "Short form:", modes[i][1])
		}
		return
	}

	if modeVal == "" {
		cobra.CheckErr("must specify a mode to update (positional argument or flag --modes / -m)")
	}

	//Parse the new flag
	newVal, err := cmd.Flags().GetString("new")
	cobra.CheckErr(err)

	path, _ := cmd.Flags().GetString(utils.FlagPath)
	abbv, _ := cmd.Flags().GetString(utils.FlagAbbreviation)
	indx, _ := cmd.Flags().GetInt(utils.FlagIndex)

	cobra.CheckErr(core.UpdatePath(modeVal, path, abbv, indx, newVal, utils.TemporalFlagPassed(cmd)))
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
	UpdateCmd.Flags().StringP("modes", "m", "", "Specify the update mode (or use 'show' to list all modes)")

	_ = UpdateCmd.RegisterFlagCompletionFunc(utils.FlagAbbreviation, CompleteAbbreviationFlag)
	_ = UpdateCmd.RegisterFlagCompletionFunc(utils.FlagPath, CompleteDirectoryFlag)

	_ = UpdateCmd.RegisterFlagCompletionFunc("modes", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		modes := []string{
			"path-path\tUpdate a Path with a new Path",
			"pp\tUpdate a Path with a new Path (short)",
			"path-abbv\tUpdate a Path with a new Abbreviation",
			"pa\tUpdate a Path with a new Abbreviation (short)",
			"path-indx\tUpdate a Path with a new Index",
			"pi\tUpdate a Path with a new Index (short)",
			"abbv-path\tUpdate an Abbreviation with a new Path",
			"ap\tUpdate an Abbreviation with a new Path (short)",
			"abbv-abbv\tUpdate an Abbreviation with a new Abbreviation",
			"aa\tUpdate an Abbreviation with a new Abbreviation (short)",
			"abbv-indx\tUpdate an Abbreviation with a new Index",
			"ai\tUpdate an Abbreviation with a new Index (short)",
			"indx-path\tUpdate an Index with a new Path",
			"ip\tUpdate an Index with a new Path (short)",
			"indx-abbv\tUpdate an Index with a new Abbreviation",
			"ia\tUpdate an Index with a new Abbreviation (short)",
			"indx-indx\tUpdate an Index with a new Index",
			"ii\tUpdate an Index with a new Index (short)",
			"show\tShow all modes to update",
		}
		var completions []string
		for _, m := range modes {
			if toComplete == "" || strings.HasPrefix(m, toComplete) {
				completions = append(completions, m)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})

	// This completion function is used for the "new" flag, it checks the mode being updated and if it's a path update, 
	// it will only show directories for completion
	_ = UpdateCmd.RegisterFlagCompletionFunc("new", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var mode string
		modesFlag, _ := cmd.Flags().GetString("modes")
		if modesFlag != "" {
			mode = modesFlag
		} else if len(args) > 0 {
			mode = args[0]
		}

		if mode == "path-path" || mode == "pp" || mode == "abbv-path" || mode == "ap" || mode == "indx-path" || mode == "ip" {
			return nil, cobra.ShellCompDirectiveFilterDirs
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	})
}
