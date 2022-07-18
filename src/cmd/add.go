package cmd

import (
	"goto/src/gpath"
	"goto/src/utils"

	"github.com/spf13/cobra"
)

// addCmd represents the addGPath command
var addCmd = &cobra.Command{
	Use:     "add-path",
	Aliases: []string{"add", "create-path", "create"},
	Short:   "Add a new path to goto-paths file",
	Long:    `To use the add-path command you need to pass two args: a path and an abbreviation to create a new goto-path`,
	Example: `
# This command add the current directory to the gpaths file with the abbreviation "currentDir"	
goto add-path --current currentDir 

# The same that:
goto add-path ./ currentDir

# To specify the path and abbreviation use:
goto add-path ~/Documents docs

Note: Use -d to search only directorys and don't check if the argument passed is a abbreviation or index
`,
	Args: cobra.RangeArgs(1, 2),
	PreRun: func(cmd *cobra.Command, args []string) {

		// If only 1 args is passed (an abbreviation) the CurrentDirFlag must be passed
		if len(args) == 1 && !utils.CurrentDirFlagPassed(cmd) {
			cobra.CheckErr("must be pass 2 arguments, a path and an abbreviation. Or pass the current directory flag and an abbreviation ")
		}

		// If 2 args are passed the CurrentDirFlag must not be passed
		if len(args) == 2 && utils.CurrentDirFlagPassed(cmd) {
			cobra.CheckErr("if the current directory flag is passed, only pass an abbreviation")
		}
	},

	Run: func(cmd *cobra.Command, args []string) {

		gpaths := utils.LoadGPaths(cmd)

		// Vars to create the gpath and to report a error
		var path, abbv string
		var err error

		// If the CurrentDirFlag is passed, only valid the abbreviation
		if utils.CurrentDirFlagPassed(cmd) {
			path = utils.GetCurrentDirectory()

			// If the , the abbreviation is the argument number 0
			abbv, err = gpath.ValidAbbreviation(args[0])
			cobra.CheckErr(err)

		} else {
			// If CurrentDirFlag is passed, the path is the argument number 0 and the abbreviation the number 1
			path, err = gpath.ValidPath(args[0])
			cobra.CheckErr(err)

			abbv, err = gpath.ValidAbbreviation(args[1])
			cobra.CheckErr(err)
		}

		//Add the new directory to the array
		gpaths = append(gpaths, gpath.GotoPath{
			Path:         path,
			Abbreviation: abbv,
		})

		// And added to the file
		utils.UpdateGPaths(cmd, gpaths)
	},
}

func init() {
	//Add this command to RootCmd
	rootCmd.AddCommand(addCmd)

	//Flags
	addCmd.Flags().BoolP(utils.FlagCurrentDir, "c", false, "The Path to add will be the current directory (\"path\" flag value will be overwrite)")
}
