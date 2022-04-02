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
	"goto/src/config"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	ConfigDir           string
	ConfigFile          string
	GotoPathsFile       string
	TempGotoPathsFile   string
	GotoPathsFileBackup string
)

// Return the path of
func CheckIndexOrAbbvOrDir(arg string) (string, error) {

	//Load the config file in memory
	var gpaths []config.GotoPath
	if err := config.LoadConfigFile(&gpaths, GotoPathsFile); err != nil {
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

	//Valid the path
	if err := config.ValidPath(&arg); err != nil {
		return "", err
	}

	//If the Path is valid, return it
	return arg, nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goto",
	Short: "Goto is a \"Path Manager\" that allows you to add a specific path with an identifier and after get it with that identifier",
	Long: `
Goto is a "Path Manager" that allows you to add a specific path with an identifier, this path can be used as an abbreviation or an 
index number. Those path are automatically save in a json file, the goto-paths files. From this files can add, update, delete and list
paths and abreviations.
`,
	//If don't have args, return a error
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		if path, err := CheckIndexOrAbbvOrDir(args[0]); err == nil {

			//If quote flag is passed
			if cmd.Flags().Changed("quotes") {
				fmt.Println("\"" + path + "\"")
				os.Exit(0)
			}

			//If spaces flag is passed
			if cmd.Flags().Changed("spaces") {
				fmt.Println(strings.ReplaceAll(path, " ", "\\ "))
				os.Exit(0)
			}

			//If quote flag is not passed
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
	TempGotoPathsFile = filepath.Join(os.TempDir(), "goto-paths-temp.json")

	cobra.CheckErr(config.CreateGotoPathFile(GotoPathsFile))
}

func init() {
	cobra.OnInitialize(initOfConfigVars)
	rootCmd.Flags().BoolP("quotes", "q", false, "Return the path between quotes")
	rootCmd.Flags().BoolP("spaces", "s", false, "Return the path with substituted spaces")
}
