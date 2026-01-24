package cmd

import (
	"fmt"
	"goto/src/gpath"
	"goto/src/utils"

	"github.com/spf13/cobra"
)

// addCmd represents the addGPath command
var validCmd = &cobra.Command{
	Use:     "valid-paths",
	Aliases: []string{"valid", "check-paths", "check"},
	Args:    cobra.ExactArgs(0),
	Short:   "Valid all path from goto-paths file",

	Run: func(cmd *cobra.Command, _ []string) {

		gpaths := utils.LoadGPaths(cmd)

		// Check all paths are valid from goto-paths file
		for _, g := range gpaths {
			if err := g.Valid(); err != nil {
				cobra.CheckErr(err)
			}
		}

		//Check the whole gpath array
		if err := gpath.CheckRepeatedItems(gpaths); err != nil {
			cobra.CheckErr(err)
		}

		fmt.Println("All paths are valid <3")
	},
}

func init() {
	//Add this command to RootCmd
	rootCmd.AddCommand(validCmd)
}
