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

// Return the path of
func CheckIndexOrAbbvOrDir(arg string) (string, error) {
	//Initial the variables to use config package
	config.GotoPathsFile = GotoPathsFile
	config.ConfigDir = ConfigDir

	//Load the config file in memory
	var gpaths []config.GotoPath
	if err := config.LoadConfigFile(&gpaths); err != nil {
		return "", err
	}

	//Check if path is number
	if err := config.IsValidIndex(gpaths, arg); err == nil {

		//I already kwow that "arg" is a number
		pathNumber, _ := strconv.Atoi(arg)

		for i, gpath := range gpaths {
			if pathNumber == i {
				return gpath.Path, nil
			}
		}
	}

	//If not a number, check if is an abbreviation
	for _, gpath := range gpaths {
		if arg == gpath.Abbreviation {
			return gpath.Path, nil
		}
	}

	//If it is neither a number nor an abbreviation, check if is exists file
	fileInfo, err := os.Stat(arg)
	if err == nil {
		//If exists, check if it's a directory
		if fileInfo.IsDir() {
			////Valid the path
			//if err := config.ValidPath(&arg); err != nil {
			//	return "", err
			//}
			return arg, nil
		}

		//If not a directory
		return "", fmt.Errorf("the path \"%s\" is not a directory", arg)
	}

	//If the path not exists
	if os.IsNotExist(err) {
		return "", fmt.Errorf("the path \"%s\" is not a directory", arg)
	}

	return "", fmt.Errorf("invalid argument/s")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goto",
	Short: "The ultimate way path manager in the command line",
	Long: `
Goto is a command that you can use it like a Path Manger you are allow to 
add specific path with a identifier to move faster, this path can be used like 
abbreviation or a index number.

If you use Goto with cd (for example, with aliases) you have the ultimate way 
to move between folders in the command line.

Fast and Easy to use and install`,

	//If don't have args, return a error
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		if path, err := CheckIndexOrAbbvOrDir(args[0]); err == nil {
			fmt.Println(path)

			//Return 2 because is easier for the alias.sh
			//only need if [[ "$?" == "2"]]
			os.Exit(2)
		} else {
			cobra.CheckErr(err)
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
