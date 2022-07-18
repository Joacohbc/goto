package cmd

import (
	"fmt"
	"goto/src/utils"

	"github.com/spf13/cobra"
)

// listGPathCmd represents the listGPath command
var listCmd = &cobra.Command{
	Use:     "list-path",
	Aliases: []string{"list"},
	Short:   "List goto-paths in the goto-paths file",
	Run: func(cmd *cobra.Command, _ []string) {

		//Load the goto-paths file to array
		gpaths := utils.LoadGPaths(cmd)

		if utils.FlagPassed(cmd, "reverse") { // If the reverse flag is passed
			for i := range gpaths {
				fmt.Printf("%v - %s\n", len(gpaths)-i-1, gpaths[len(gpaths)-i-1].String())
			}
			return
		}

		//If any flag is passed
		for i, gpath := range gpaths {
			fmt.Printf("%v - %s\n", i, gpath.String())
		}
	},
}

func init() {
	//Add this command to RootCommand
	rootCmd.AddCommand(listCmd)

	//Flags
	listCmd.Flags().BoolP("reverse", "R", false, "List the goto-paths in reverse")
}
