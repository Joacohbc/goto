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
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goto",
	Short: "Goto is a \"Path Manager\" that allows you to add a specific path with an identifier and after get it with that identifier",
	Long: `
Goto is a "Path Manager" that allows you to add a specific path with an identifier, this path can be used as an abbreviation or an 
index number. Those path are automatically save in a json file, the goto-paths files. From this files can add, update, delete and list
paths and abreviations.
`,

	Example: `
# Move to the destination directory
# "h" is the abbreviation of /home/user
goto h

# You also can use "0" (that is the default index of the /home/user)
goto 0

# Or also you can use goto like cd, use a complete/relative path:
goto /home/user/.config/goto

# For a temporal gpaths you have to use temporal flag(-t / --temporal)
goto -t home
`,
	//If don't have args, return a error
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		path, err := checkIndexOrAbbvOrDir(cmd, args[0])
		cobra.CheckErr(err)

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

func init() {
	cobra.OnInitialize(initVars)
	rootCmd.Flags().BoolP("quotes", "q", false, "Return the path between quotes")
	rootCmd.Flags().BoolP("spaces", "s", false, "Return the path with substituted spaces")
	rootCmd.PersistentFlags().BoolP("temporal", "t", false, "Do the action in the temporal gpath file")
}
