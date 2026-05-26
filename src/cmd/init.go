package cmd

import (
	"fmt"
	"goto/src/core"
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
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	handler := core.NewMessageHandler(func(msg core.Message) {
		fmt.Println(msg.Content)
	})

	err := core.InitializeConfig(handler.Channel())
	handler.CloseAndWait()

	cobra.CheckErr(err)

	// After successful initialization, generate autocompletion files
	configDir := utils.GetConfigDir()

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		cobra.CheckErr(fmt.Errorf("failed to create config directory for completions: %w", err))
	}

	bashCompFile := filepath.Join(configDir, "completion.bash")
	zshCompFile := filepath.Join(configDir, "completion.zsh")
	fishCompFile := filepath.Join(configDir, "completion.fish")

	// Generate Bash completion
	if f, err := os.Create(bashCompFile); err == nil {
		_ = RootCmd.GenBashCompletion(f)
		f.Close()
	}

	// Generate Zsh completion
	if f, err := os.Create(zshCompFile); err == nil {
		_ = RootCmd.GenZshCompletion(f)
		f.Close()
	}

	// Generate Fish completion
	if f, err := os.Create(fishCompFile); err == nil {
		_ = RootCmd.GenFishCompletion(f, true)
		f.Close()
	}

	// Also attempt to copy fish completion to ~/.config/fish/completions/goto.fish if the folder exists
	if homeDir, err := os.UserHomeDir(); err == nil {
		fishUserCompDir := filepath.Join(homeDir, ".config", "fish", "completions")
		shell := os.Getenv("SHELL")

		// If Fish is the active shell, ensure the completions directory exists
		if strings.Contains(shell, "fish") {
			_ = os.MkdirAll(fishUserCompDir, 0755)
		}

		if _, err := os.Stat(fishUserCompDir); err == nil {
			fishUserCompFile := filepath.Join(fishUserCompDir, "goto.fish")
			if f, err := os.Create(fishUserCompFile); err == nil {
				_ = RootCmd.GenFishCompletion(f, true)
				f.Close()
			}
		}
	}
}

func init() {
	RootCmd.AddCommand(InitCmd)
}
