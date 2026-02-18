package cmd

import (
	"fmt"
	"goto/src/core"
	"goto/src/utils"

	"github.com/spf13/cobra"
)

// ListCmd represents the listGPath command
var ListCmd = &cobra.Command{
	Use:     "list-path",
	Aliases: []string{"list"},
	Short:   "List goto-paths in the goto-paths file",
	Example: `
# Format: goto list [ -t ] [ -R ]

# List all gpaths
goto list

# List all gpaths form temporal file
goto list -t
`,
	Run: runList,
}

func runList(cmd *cobra.Command, _ []string) {

	//Load the goto-paths file to array
	gpaths, err := core.ListPaths(cmd.Flags().Changed(utils.FlagTemporal))
	cobra.CheckErr(err)

	if cmd.Flags().Changed("reverse") { // If the reverse flag is passed
		for i := range gpaths {
			fmt.Printf("%v - %s\n", len(gpaths)-i-1, gpaths[len(gpaths)-i-1].String())
		}
		return
	}

	//If any flag is passed
	for i, gpath := range gpaths {
		fmt.Printf("%v - %s\n", i, gpath.String())
	}
}

func init() {
	//Add this command to RootCommand
	RootCmd.AddCommand(ListCmd)

	//Flags
	ListCmd.Flags().BoolP("reverse", "R", false, "List the goto-paths in reverse")
}
