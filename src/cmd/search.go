package cmd

import (
	"fmt"
	"goto/src/utils"

	"github.com/spf13/cobra"
)

// SearchCmd represents the searchGPath command
var SearchCmd = &cobra.Command{
	Use:     "search-path",
	Aliases: []string{"search", "find-path", "find"},
	Short:   "search goto-paths in the goto-paths file",
	Example: `
# Format: goto search [ -t ] { -p path | -a abbreviation }
# To search a specific goto-path you can use the Path or the Abbreviation 
goto search --path ~/Documents
goto search --abbv docs
`,
	PreRun: func(cmd *cobra.Command, args []string) {

		//If the number or flags are 0 or more than 2, return an error
		if cmd.Flags().NFlag() == 0 || cmd.Flags().NFlag() > 2 {
			cobra.CheckErr("you must specify only one flag to find a gpath (Or Path or Abbreviation)")
		}

		//If only one flags is passed and it is the temporary flags, return an error
		if cmd.Flags().NFlag() == 1 && utils.TemporalFlagPassed(cmd) {
			cobra.CheckErr("you must specify one flag to find a gpath (Or Path or Abbreviation)")
		}
	},
	Run: func(cmd *cobra.Command, _ []string) {

		//Load the goto-paths file to array
		gpaths := utils.LoadGPaths(cmd)

		// If the any path flag is passed
		if utils.PathFlagPassed(cmd) {

			path := utils.GetPath(cmd)

			for i, gpath := range gpaths {
				if gpath.Path == path {
					fmt.Printf("%v - %s\n", i, gpath.String())
					return
				}
			}
			cobra.CheckErr(fmt.Errorf("the path \"%s\" doesn't exist in the gpaths-file", path))

		} else {
			abbv := utils.GetAbbreviation(cmd)

			for i, gpath := range gpaths {
				if gpath.Abbreviation == abbv {
					fmt.Printf("%v - %s\n", i, gpath.String())
					return
				}
			}
			cobra.CheckErr(fmt.Errorf("doesn't exist a path with that abbreviation \"%s\"", abbv))
		}
	},
}

func init() {
	//Add this command to RootCommand
	RootCmd.AddCommand(SearchCmd)

	//Flags
	SearchCmd.Flags().StringP(utils.FlagPath, "p", "", "The Path to delete")
	SearchCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "The Abbreviation of the Path")
}
