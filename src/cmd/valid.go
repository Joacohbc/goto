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
	Aliases: []string{"valid", "check-paths"},
	Args:    cobra.ExactArgs(0),
	Short:   "Valid all path from goto-paths file",

	Run: func(cmd *cobra.Command, _ []string) {

		gpaths := utils.LoadGPaths(cmd)

		for _, g := range gpaths {

			if err := g.Valid(); err != nil {
				fmt.Println("Error in the file:", err)
				return
			}
		}

		if err := gpath.DontRepeatInArray(gpaths); err != nil {
			fmt.Println("Error in the file:", err)
			return
		}

		fmt.Println("All paths are valid <3")
	},
}

func init() {
	//Add this command to RootCmd
	rootCmd.AddCommand(validCmd)
}
