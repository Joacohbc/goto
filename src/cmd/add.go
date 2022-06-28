/*
Copyright © 2022 Joacohbc <joacog48@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
	Args:    cobra.ExactArgs(0),
	Short:   "Add a new path to goto-paths file",
	Long: `
To use the add-path command you need to pass two args: a "Path" and an "Abbreviation" to 
create a new goto-path`,

	Example: `
# This command add the current directory(the "Path") to the gpaths file with
# the abbreviation "currentDir"	
goto add-path --current -abbv currentDir

# To specify the "Path" and "Abbreviation" use:
goto add-path --path ~/Documents -abbv docs
`,

	Run: func(cmd *cobra.Command, args []string) {

		gpaths := utils.LoadGPaths(cmd)

		//Add the new directory to the array
		gpaths = append(gpaths, gpath.GotoPath{
			Path:         utils.GetPath(cmd),
			Abbreviation: utils.GetAbbreviation(cmd),
		})

		// And added to the file
		utils.UpdateGPaths(cmd, gpaths)
	},
}

func init() {
	//Add this command to RootCmd
	rootCmd.AddCommand(addCmd)

	//Flags
	addCmd.Flags().StringP(utils.FlagPath, "p", "", "The Path to add")
	addCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "The Abbreviation of the Path")
	addCmd.Flags().BoolP(utils.FlagCurretDir, "c", false, "The Path to add will be the current directory (\"path\" flag will be overwrite)")
}
