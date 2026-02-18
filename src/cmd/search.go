package cmd

import (
	"fmt"
	"goto/src/core"
	"goto/src/utils"

	"github.com/spf13/cobra"
)

// SearchCmd represents the searchGPath command
var SearchCmd = &cobra.Command{
	Use:     "search-path",
	Aliases: []string{"search", "find-path", "find"},
	Short:   "Search goto-paths in the goto-paths file",
	Example: `
# Format: goto search [ -t ] { -p path | -a abbreviation }
# To search a specific goto-path you can use the Path or the Abbreviation 
goto search --path ~/Documents
goto search --abbv docs
`,
	PreRun: preRunSearch,
	Run:    runSearch,
}

func preRunSearch(cmd *cobra.Command, args []string) {

	//If the number or flags are 0 or more than 2, return an error
	if cmd.Flags().NFlag() == 0 || cmd.Flags().NFlag() > 2 {
		cobra.CheckErr("you must specify only one flag to find a gpath (Or Path or Abbreviation)")
	}

	//If only one flags is passed and it is the temporary flags, return an error
	if cmd.Flags().NFlag() == 1 && cmd.Flags().Changed(utils.FlagTemporal) {
		cobra.CheckErr("you must specify one flag to find a gpath (Or Path or Abbreviation)")
	}
}

func runSearch(cmd *cobra.Command, _ []string) {
	path, _ := cmd.Flags().GetString(utils.FlagPath)
	abbv, _ := cmd.Flags().GetString(utils.FlagAbbreviation)

	idx, gpath, err := core.SearchPath(path, abbv, cmd.Flags().Changed(utils.FlagTemporal))
	cobra.CheckErr(err)

	fmt.Printf("%v - %s\n", idx, gpath.String())
}

func init() {
	//Add this command to RootCommand
	RootCmd.AddCommand(SearchCmd)

	//Flags
	SearchCmd.Flags().StringP(utils.FlagPath, "p", "", "The Path to delete")
	SearchCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "The Abbreviation of the Path")
}
