package cmd

import (
	"fmt"
	"goto/src/gpath"
	"goto/src/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goto",
	Short: "Goto is a \"Path Manager\" that allows you to add a specific path with an identifier and after get it with that identifier (exit with status 2)",
	Long: `
Goto is a "Path Manager" that allows you to add a specific path with an identifier, this path can be used as an abbreviation or an 
index number. Those path are automatically save in a json file, the goto-paths files. From this files can add, update, delete and list
paths and abbreviations.
`,

	Example: `
# Format: goto [ -t ] { abbreviation | path | index }

# Move to the destination directory
# "h" is the abbreviation of /home/user
goto h

# You also can use "0" (that is the default index of the /home/user)
goto 0

# Or also you can use goto like cd, use a complete/relative path:
goto /home/user/.config/goto

# For a temporal gpaths you have to use temporal flag(-t / --temporal)
goto -t home

# If you have a directory named like a number or like abbreviation you should use -d / --only-directory flag
goto -d 1 # This will move to the directory "1" and don't move to the first path in the gpaths file
goto -d h # This will move to the directory "h" and don't move to the path with the abbreviation "h"
`,
	//If don't have args, return a error
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		path := filepath.Join(args...)

		// If only directory flag is passed, check if is a directory a continue
		if cmd.Flags().Changed("only-directory") {
			/*
				This Flags is use if the directory is a named "123/". Also if the directory
				is a named like a direcotry an abbreviation in the gpath files and you want
				to go to the directory and not to the abbreviation
			*/
			cobra.CheckErr(gpath.ValidPathVar(&path))
		} else {
			// If it is not passed
			var isIndexOrAbbv bool

			// Check if is a index or an abbreviation
			path, isIndexOrAbbv = utils.IsIndexOrAbbv(cmd, path)

			// If it is not, check if is a directory
			if !isIndexOrAbbv {
				cobra.CheckErr(gpath.ValidPathVar(&path))
			}

		}

		//If quote flag is passed
		if cmd.Flags().Changed("quotes") {
			fmt.Println("\"" + path + "\"")
			return
		}

		//If spaces flag is passed
		if cmd.Flags().Changed("spaces") {
			fmt.Println(strings.ReplaceAll(path, " ", "\\ "))
			return
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
	rootCmd.Flags().BoolP("quotes", "q", false, "Return the path between quotes")
	rootCmd.Flags().BoolP("spaces", "s", false, "Return the path with substituted spaces")
	rootCmd.Flags().BoolP("only-directory", "d", false, "Only check if the argument passed is a directory")
	rootCmd.PersistentFlags().BoolP("temporal", "t", false, "Do the action in the temporal gpath file")
}
