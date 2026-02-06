package cmd

import (
	"fmt"
	"goto/src/core"
	"goto/src/utils"

	"github.com/spf13/cobra"
)

// ValidCmd represents the addGPath command
var ValidCmd = &cobra.Command{
	Use:     "valid-paths",
	Aliases: []string{"valid", "check-paths", "check"},
	Args:    cobra.ExactArgs(0),
	Short:   "Validate all paths from goto-paths file",

	Run: runValid,
}

func runValid(cmd *cobra.Command, _ []string) {

	cobra.CheckErr(core.ValidatePaths(utils.TemporalFlagPassed(cmd)))

	fmt.Println("All paths are valid <3")
}

func init() {
	//Add this command to RootCmd
	RootCmd.AddCommand(ValidCmd)
}
