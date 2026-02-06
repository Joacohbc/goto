package core

import (
	"fmt"
	"goto/src/utils"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// InitializeConfig sets up the configuration directory, binary, and shell alias.
func InitializeConfig(msgChan chan<- Message) error {
	notifier := NewNotifier(msgChan)
	goos := runtime.GOOS
	if goos == "windows" {
		return fmt.Errorf("initialization is not supported on Windows")
	}

	configDir := utils.GetConfigDir()
	if err := ensureConfigDir(configDir, notifier); err != nil {
		return err
	}

	exePath, err := installBinary(configDir, notifier)
	if err != nil {
		return err
	}

	aliasFile, err := generateAliasFile(configDir, exePath, notifier)
	if err != nil {
		return err
	}

	return configureShell(aliasFile, notifier)
}

func ensureConfigDir(configDir string, notifier *Notifier) error {
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config dir: %w", err)
		}
		notifier.Success("Config directory created.")
	} else {
		notifier.Info("Config directory already exists.")
	}
	return nil
}

func installBinary(configDir string, notifier *Notifier) (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", err
	}

	targetExe := filepath.Join(configDir, "goto.bin")

	if exePath != targetExe {
		input, err := os.ReadFile(exePath)
		if err != nil {
			return "", err
		}

		err = os.WriteFile(targetExe, input, 0755)
		if err != nil {
			return "", err
		}
		notifier.Success("Copied binary to %s", targetExe)
		return targetExe, nil
	}
	return exePath, nil
}

func generateAliasFile(configDir, exePath string, notifier *Notifier) (string, error) {
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
	if err != nil {
		return "", err
	}
	notifier.Success("Generated %s", aliasFile)
	return aliasFile, nil
}

func configureShell(aliasFile string, notifier *Notifier) error {
	shell := os.Getenv("SHELL")
	var shellRC string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

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
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	f, err := os.OpenFile(shellRC, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := os.ReadFile(shellRC)
	if err != nil {
		return err
	}

	sourceCmd := fmt.Sprintf("source %s", aliasFile)
	if strings.Contains(string(content), sourceCmd) {
		notifier.Info("Alias already sourced in %s", shellRC)
		return nil
	} else {
		if _, err := f.WriteString(fmt.Sprintf("\n#Aliases to use goto:\n%s\n", sourceCmd)); err != nil {
			return err
		}
		notifier.Success("Added source command to %s\nPlease restart your terminal or run 'source %s' to activate goto.", shellRC, shellRC)
		return nil
	}
}
