package cmd

import (
	"fmt"
	"goto/src/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize goto configuration and shell alias",
	Long:  `Initialize the .config directory, goto-path.json, and generate alias.sh. It also adds the alias to your shell configuration.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		configDir := utils.GetConfigDir()
		ensureConfigDir(configDir)

		exePath := installBinary(configDir)
		aliasFile := generateAliasFile(configDir, exePath)

		configureShell(aliasFile)
		fmt.Println("Initialization complete. You can delete this file now (was already copied to the config directory). Enjoy using goto!")
	},
}

func init() {
	RootCmd.AddCommand(InitCmd)
}

func ensureConfigDir(configDir string) {
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			cobra.CheckErr(fmt.Errorf("failed to create config dir: %w", err))
		}
		fmt.Println("Config directory created.")
	} else {
		fmt.Println("Config directory already exists.")
	}
}

func installBinary(configDir string) string {
	exePath, err := os.Executable()
	cobra.CheckErr(err)

	exePath, err = filepath.EvalSymlinks(exePath)
	cobra.CheckErr(err)

	targetExe := filepath.Join(configDir, "goto.bin")

	if exePath != targetExe {
		input, err := os.ReadFile(exePath)
		cobra.CheckErr(err)

		err = os.WriteFile(targetExe, input, 0755)
		cobra.CheckErr(err)
		fmt.Printf("Copied binary to %s\n", targetExe)
		return targetExe
	}
	return exePath
}

func generateAliasFile(configDir, exePath string) string {
	aliasContent := fmt.Sprintf(`#!/bin/bash
GOTO_FILE="%s"
#GOTO FUNC
goto() {
    OUTPUT=$("$GOTO_FILE" $@)
    
    #If the return "2", the program return a gpath successfully
    if [ $? -eq 2 ]; then
        cd "$OUTPUT"   
        echo "Go to:" $OUTPUT
    elif [ $? -eq 1 ]; then # If error exit with status 1
        echo "$OUTPUT"
        return 1
    else
        echo "$OUTPUT"
    fi
}

#cd is change by goto function
alias cd="goto"
alias cdt="goto -t"
`, exePath)

	aliasFile := filepath.Join(configDir, "alias.sh")
	err := os.WriteFile(aliasFile, []byte(aliasContent), 0644)
	cobra.CheckErr(err)
	fmt.Printf("Generated %s\n", aliasFile)
	return aliasFile
}

func configureShell(aliasFile string) {
	shell := os.Getenv("SHELL")
	var shellRC string
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	if strings.Contains(shell, "zsh") {
		shellRC = filepath.Join(homeDir, ".zshrc")
	} else if strings.Contains(shell, "bash") {
		shellRC = filepath.Join(homeDir, ".bashrc")
	} else if strings.Contains(shell, "fish") {
		shellRC = filepath.Join(homeDir, ".config", "fish", "config.fish")
	} else if strings.Contains(shell, "dash") {
		shellRC = filepath.Join(homeDir, ".profile")
	} else if strings.Contains(shell, "tcsh") {
		shellRC = filepath.Join(homeDir, ".tcshrc")
	} else if strings.Contains(shell, "csh") {
		shellRC = filepath.Join(homeDir, ".cshrc")
	} else if strings.Contains(shell, "ksh") {
		shellRC = filepath.Join(homeDir, ".kshrc")
	} else if strings.HasSuffix(shell, "/sh") || shell == "sh" {
		shellRC = filepath.Join(homeDir, ".profile")
	} else {
		cobra.CheckErr(fmt.Errorf("unsupported shell: %s", shell))
	}

	f, err := os.OpenFile(shellRC, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	cobra.CheckErr(err)
	defer f.Close()

	content, err := os.ReadFile(shellRC)
	cobra.CheckErr(err)

	sourceCmd := fmt.Sprintf("source %s", aliasFile)
	if strings.Contains(string(content), sourceCmd) {
		fmt.Printf("Alias already sourced in %s\n", shellRC)
	} else {
		if _, err := f.WriteString(fmt.Sprintf("\n#Aliases to use goto:\n%s\n", sourceCmd)); err != nil {
			cobra.CheckErr(err)
		}
		fmt.Printf("Added source command to %s\n", shellRC)
		fmt.Println("Please restart your terminal or run 'source " + shellRC + "' to activate goto.")
	}
}
