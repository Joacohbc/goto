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
	"encoding/json"
	"fmt"
	"goto/src/config"
	"os"

	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Do a restoreCmd of the goto-paths file",

	Run: func(cmd *cobra.Command, args []string) {

		input, err := cmd.Flags().GetString("input")
		cobra.CheckErr(err)

		//If exists a config backup
		info, err := os.Stat(input)
		if os.IsNotExist(err) {
			cobra.CheckErr(fmt.Errorf("dont have a backup of config file"))
		} else if err != nil {
			cobra.CheckErr(err)
		}

		//If is not a file
		if info.IsDir() {
			cobra.CheckErr(fmt.Errorf("the input can't be a directory"))
		}

		//Read the config backup
		backup, err := os.ReadFile(input)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("cant read the backup of config file: %v", err))
		}

		//Do the unmarshaling of the config backup
		var gpaths []config.GotoPath
		if err := json.Unmarshal(backup, &gpaths); err != nil {
			cobra.CheckErr(fmt.Errorf("cant parse the backup of config file"))
		}

		//Initial the variables to use config
		config.GotoPathsFile = GotoPathsFile

		//And re-write the config file with the backup
		cobra.CheckErr(config.CreateJsonFile(gpaths))

		fmt.Println("Restore complete")
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().StringP("input", "i", GotoPathsFileBackup, "The ubication of the backup file")
}
