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

	Run: func(cmd *cobra.Command, args []string) {

		gpaths := utils.LoadGPaths(cmd)

		for _, g := range gpaths {

			if err := g.Valid(); err != nil {
				fmt.Println("Error in the file:", err)
				return
			}
		}

		if err := gpath.ValidArray(gpaths); err != nil {
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
