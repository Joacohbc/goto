package cmd

import (
	"fmt"
	"goto/src/core"
	"goto/src/utils"

	"github.com/spf13/cobra"
)

const msgPathDeleted = "The path %s (%s) was deleted\n"

// DeleteCmd represents the addGPath command
var DeleteCmd = &cobra.Command{
	Use:     "delete-path",
	Aliases: []string{"del", "delete", "remove-path", "rem", "remove"},
	Short:   "Delete a path from goto-path file",
	Long:    `To use the delete-path command you need to provide a Path, Abbreviation or Index to identify the goto-path to delete`,

	Example: `
# Format : goto delete-path [ -t ] { -p path | -a abbreviation | -i index } 

# To specify the "Path", "Abbreviation" or Index. use:

# Delete the gpath with the path "/home/user/Documents"
goto delete-path --path ~/Documents

# Delete the gpath with the abbreviation "docs"
goto delete-path --abbv docs

# Delete the gpath in the index "2"
goto delete-path --indx 2
`,
	Args: cobra.ExactArgs(0),
	PreRun: func(cmd *cobra.Command, _ []string) {

		/*
			Valid cases:
			- Specify only one flag, to indicate which gpath will be deleted
			- Specify 2 flags, to indicate which gpath will be deleted and the temporary flag

			Invalid cases:
			- 0 flags
			- +2 flags
			- 1 flag that it is the temporary flag
			- 2 flags and one of these it is not the temporary flag
		*/

		//If the number or flags are 0 or more than 2, return an error
		if cmd.Flags().NFlag() == 0 || cmd.Flags().NFlag() > 2 {
			cobra.CheckErr("you must specify only one flag to delete a gpath (Or Path or Abbreviation or Index)")
		}

		//If only one flags is passed and it is the temporary flags, return an error
		if cmd.Flags().NFlag() == 1 && utils.TemporalFlagPassed(cmd) {
			cobra.CheckErr("you must specify one flag to delete a gpath (Or Path or Abbreviation or Index)")
		}

		//If 2 flags are passed and none of these two is not the temporary flag
		//it means that two flags were passed to identify the gpath
		/*
			2 flags to identify the gpath may cause an error to delete the path.
			For example: -p /home/user -i 2, the index not match with the gpath, so delete one of the paths
		*/
		if cmd.Flags().NFlag() == 2 && !utils.TemporalFlagPassed(cmd) {
			cobra.CheckErr(fmt.Errorf("you must specify only one flag to delete a gpath (Or Path or Abbreviation or Index)"))
		}
	},
	Run: runDelete,
}

func runDelete(cmd *cobra.Command, _ []string) {
	path, _ := cmd.Flags().GetString(utils.FlagPath)
	abbv, _ := cmd.Flags().GetString(utils.FlagAbbreviation)
	indx, _ := cmd.Flags().GetInt(utils.FlagIndex)

	deleted, err := core.DeletePath(path, abbv, indx, utils.TemporalFlagPassed(cmd))
	cobra.CheckErr(err)

	fmt.Printf(msgPathDeleted, deleted.Path, deleted.Abbreviation)
}

func init() {
	RootCmd.AddCommand(DeleteCmd)

	//Flags
	DeleteCmd.Flags().StringP(utils.FlagPath, "p", "", "The Path to delete")
	DeleteCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "The Abbreviation of the Path")
	DeleteCmd.Flags().IntP(utils.FlagIndex, "i", -1, "The Index of the Path")
}
