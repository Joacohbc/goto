package cmd

import (
	"goto/src/core"
	"goto/src/utils"

	"github.com/spf13/cobra"
)

// AddCmd represents the addGPath command
var AddCmd = &cobra.Command{
	Use:     "add-path",
	Aliases: []string{"add", "create-path", "create"},
	Short:   "Add a new path to goto-paths file",
	Long:    `To use the add-path command you need to pass two args: a path and an abbreviation to create a new goto-path`,
	Example: `
# Format: goto add-path [ -t ] path abbv

# This command add the current directory to the gpaths file with the abbreviation "currentDir"
goto add-path ./ currentDir

# To specify the path and abbreviation use:
goto add-path ~/Documents docs
`,
	Args: cobra.ExactArgs(2),

	Run: runAdd,
}

func runAdd(cmd *cobra.Command, args []string) {
	cobra.CheckErr(core.AddPath(args[0], args[1], cmd.Flags().Changed(utils.FlagTemporal)))
}

func init() {
	//Add this command to RootCmd
	RootCmd.AddCommand(AddCmd)
}
