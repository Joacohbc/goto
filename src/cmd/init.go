package cmd

import (
	"fmt"
	"goto/src/core"

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
	msg, err := core.InitializeConfig()
	cobra.CheckErr(err)
	if msg != "" {
		fmt.Println(msg)
	}
}

func init() {
	RootCmd.AddCommand(InitCmd)
}
