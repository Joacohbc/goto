package cmd

import (
	"fmt"
	"goto/src/core"
	"goto/src/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "goto",
	Short: "Goto is a \"Path Manager\" that allows you to add a specific path with an identifier and after get it with that identifier (exit with status 2)",
	Long: `
Goto is a "Path Manager" that allows you to add a specific path with an identifier. This path can be used as an abbreviation or an 
index number. These paths are automatically saved in a json file, the goto-paths file. You can add, update, delete and list
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

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Only autocomplete abbreviations when completing the single positional argument of goto
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// If toComplete starts with a path prefix (e.g. "./", "/", "~", "../") or contains "/", suggest directories
		if strings.HasPrefix(toComplete, ".") || strings.HasPrefix(toComplete, "/") || strings.HasPrefix(toComplete, "~") || strings.Contains(toComplete, "/") {
			return nil, cobra.ShellCompDirectiveFilterDirs
		}

		// Otherwise, load saved abbreviations and return them as completions
		gpaths, err := utils.LoadGPaths(utils.TemporalFlagPassed(cmd))
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string
		for _, gp := range gpaths {
			if gp.Abbreviation != "" {
				if toComplete == "" || strings.HasPrefix(gp.Abbreviation, toComplete) {
					completions = append(completions, gp.Abbreviation+"\t"+gp.Path)
				}
			}
		}

		// Also suggest the directories of the current directory alongside the abbreviations
		directive := cobra.ShellCompDirectiveNoFileComp
		if entries, err := os.ReadDir("."); err == nil {
			for _, entry := range entries {
				// Skip hidden directories, they are reachable through the path prefix branch (e.g. "./.config")
				if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") && strings.HasPrefix(entry.Name(), toComplete) {
					completions = append(completions, entry.Name()+string(os.PathSeparator))
					// Avoid adding a space after a directory so the completion can continue into its subdirectories
					directive = cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
				}
			}
		}

		// Use ShellCompDirectiveNoFileComp to show our abbreviations and directories without regular files
		return completions, directive
	},

	Run: runRoot,
}

func runRoot(cmd *cobra.Command, args []string) {

	path, err := core.ResolvePath(args, cmd.Flags().Changed("only-directory"), utils.TemporalFlagPassed(cmd))
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
}

// StartExecution adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func StartExecution() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().BoolP("quotes", "q", false, "Return the path between quotes")
	RootCmd.Flags().BoolP("spaces", "s", false, "Return the path with substituted spaces")
	RootCmd.Flags().BoolP("only-directory", "d", false, "Only check if the argument passed is a directory")
	RootCmd.PersistentFlags().BoolP("temporal", "t", false, "Do the action in the temporal gpath file")
}

// CompleteAbbreviationFlag is a reusable flag completion function that autocompletes saved abbreviations with their paths as descriptions.
func CompleteAbbreviationFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	gpaths, err := utils.LoadGPaths(utils.TemporalFlagPassed(cmd))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var completions []string
	for _, gp := range gpaths {
		if gp.Abbreviation != "" {
			if toComplete == "" || strings.HasPrefix(gp.Abbreviation, toComplete) {
				completions = append(completions, gp.Abbreviation+"\t"+gp.Path)
			}
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteDirectoryFlag is a reusable flag completion function that dynamically suggests directories only.
func CompleteDirectoryFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveFilterDirs
}
