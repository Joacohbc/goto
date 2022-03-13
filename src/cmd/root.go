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
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	ConfigDir           string
	ConfigFile          string
	GotoPathsFile       string
	GotoPathsFileBackup string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goto",
	Short: "The ultimate way to move between folders in the command line",
	Long: `
Goto is a command that can be used like cd, and also allows you to 
add specific path to move faster, this path can be used like abbreviation or 
a index number`,

	//If don't have args, return a error
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		//Initial the variables to use config package
		config.GotoPathsFile = GotoPathsFile
		config.ConfigDir = ConfigDir

		printAndExit := func(args ...interface{}) {
			fmt.Println(args...)

			//Return 3 because is easier for the alias.sh
			//only need if [[ "$?" == "3"]]
			os.Exit(3)
		}

		//Load the config file in memory
		var gpaths []config.GotoPath
		cobra.CheckErr(config.LoadConfigFile(&gpaths))

		//Check if path is number
		if pathNumber, err := strconv.Atoi(args[0]); err == nil {

			//If the path is over the max index return error
			if pathNumber < 0 || pathNumber > len(gpaths)-1 {
				cobra.CheckErr(fmt.Errorf("the number is invalid(should be: 0-" + strconv.Itoa(len(gpaths)-1) + "), check config file"))
			}

			for i, gpath := range gpaths {
				if pathNumber == i {
					printAndExit(gpath.Path)
				}
			}
		}

		//If not a number, check if is an abbreviation
		for _, gpath := range gpaths {
			if args[0] == gpath.Abbreviation {
				printAndExit(gpath.Path)
			}
		}

		//If it is neither a number nor an abbreviation, check if is exists file
		fileInfo, err := os.Stat(args[0])
		if err == nil {
			//If exists, check if it's a directory
			if fileInfo.IsDir() {
				printAndExit(filepath.Clean(args[0]))
			}

			//If not a directory
			cobra.CheckErr(fmt.Errorf("the path \"%s\" is not a directory", args[0]))
		}

		//If the path not exists
		if os.IsNotExist(err) {
			cobra.CheckErr(fmt.Errorf("the path \"%s\" is not a directory", args[0]))
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initOfConfigVars() {
	//Get the directory
	configPath, err := os.UserConfigDir()
	cobra.CheckErr(err)

	ConfigDir = filepath.Join(configPath, "/goto/")
	GotoPathsFile = filepath.Join(ConfigDir, "config.json")
	GotoPathsFile = filepath.Join(ConfigDir, "goto-paths.json")
	GotoPathsFileBackup = filepath.Clean(GotoPathsFile + ".backup")

	//Create the json file

	//Initial the variables to use config package
	config.GotoPathsFile = GotoPathsFile
	config.ConfigDir = ConfigDir

	cobra.CheckErr(config.CreateConfigFile())
}

func init() {
	cobra.OnInitialize(initOfConfigVars)
}
