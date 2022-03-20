/*
Copyright © 2022 Joaquin Genova <joaquingenovag8@gmail.com>

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
	"goto/src/config"
	"os"

	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Do a backup of the goto-paths file",

	Run: func(cmd *cobra.Command, args []string) {

		//Initial the variables to use config
		config.GotoPathsFile = GotoPathsFile

		//Read the config file
		var gpaths []config.GotoPath
		cobra.CheckErr(config.LoadConfigFile(&gpaths))

		//Get output flag
		output, err := cmd.Flags().GetString("output")
		cobra.CheckErr(err)

		//Check if exists
		if _, err := os.Stat(output); err == nil {
			cobra.CheckErr(fmt.Errorf("the file \"%s\" already exists", output))
		}

		//Create the GotoPath file in the output
		config.GotoPathsFile = output
		cobra.CheckErr(config.CreateJsonFile(gpaths))

		fmt.Println("Backup complete")
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	//Flags
	backupCmd.Flags().StringP("output", "o", GotoPathsFileBackup, "The backup destination path")
}