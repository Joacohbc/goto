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
	Args:    cobra.ExactArgs(0),
	Short:   "List goto-paths in the goto-paths file",
	Example: `
# To list all goto-paths
goto list

# To list a specific goto-path you can use the Path or the Abbreviation 
goto list --path ~/Documents
goto list --abbv docs
`,

	Run: func(cmd *cobra.Command, _ []string) {

		//Load the goto-paths file to array
		gpaths := utils.LoadGPaths(cmd)

		// If the any path flag is passed
		if utils.CurrentDirFlagPassed(cmd) || utils.PathFlagPassed(cmd) {

			path := utils.GetPath(cmd)

			for i, gpath := range gpaths {
				if gpath.Path == path {
					fmt.Printf("%v - %s\n", i, gpath.String())
					return
				}
			}
			cobra.CheckErr(fmt.Sprintf("the path \"%s\" doesn't exist in the gpaths-file", path))

		} else if utils.AbbvFlagPassed(cmd) { // If the abbreviation flag is passed
			abbv := utils.GetAbbreviation(cmd)
			for i, gpath := range gpaths {

				if gpath.Abbreviation == abbv {
					fmt.Printf("%v - %s\n", i, gpath.String())
					return
				}
			}
			cobra.CheckErr(fmt.Sprintf("doesn't exist a path with that abbreviation \"%s\"", abbv))

		} else if utils.FlagPassed(cmd, "reverse") { //If the flag "reverse" is passed
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
	listCmd.Flags().StringP(utils.FlagPath, "p", "", "The Path to delete")
	listCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "The Abbreviation of the Path")
	listCmd.Flags().BoolP(utils.FlagCurrentDir, "c", false, "The Path to update will be the current directory (\"path\" parameter will be overwrite)")
	listCmd.Flags().BoolP("reverse", "R", false, "List the goto-paths in reverse")
}
