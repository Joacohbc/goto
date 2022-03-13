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
	"io/ioutil"
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
		info, err := os.Stat(output)
		if os.IsNotExist(err) {
			cobra.CheckErr(fmt.Errorf("dont have a backup of config file"))
		} else if err != nil {
			cobra.CheckErr(err)
		}

		//And if not a file
		if info.IsDir() {
			cobra.CheckErr(fmt.Errorf("the output can't be a directory"))
		}

		//Do the Array to JSON Bytes
		json, err := json.Marshal(gpaths)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("cant parse the config file: %v", err))
		}

		//If the output is not empty
		//And write the config backup
		if err := ioutil.WriteFile(output, json, 0600); err != nil {
			cobra.CheckErr(fmt.Errorf("cant create the backup of config file: %v", err))
		}

		fmt.Println("Backup complete")
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringP("output", "o", GotoPathsFileBackup, "The backup destination path")
}
