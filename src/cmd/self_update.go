package cmd

import (
	"fmt"
	"goto/src/core"
	"runtime"

	"github.com/spf13/cobra"
)

// ConsoleLogger implementa core.Logger para imprimir en consola standard
type ConsoleLogger struct{}

func (l *ConsoleLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// UpdateBinaryCmd represents the update command for the binary itself
var UpdateBinaryCmd = &cobra.Command{
	Use:   "update-goto",
	Short: "Update goto to the latest version",
	Long:  `Downloads the latest release from GitHub and updates the current binary if a newer version is available.`,
	Run: func(cmd *cobra.Command, args []string) {
		updateBinary()
	},
}

func init() {
	RootCmd.AddCommand(UpdateBinaryCmd)
}

func updateBinary() {
	goos := runtime.GOOS
	if goos == "windows" {
		fmt.Println("Self-update is not supported on Windows.")
		return
	}

	handler := core.NewMessageHandler(func(msg core.Message) {
		// Por ahora imprimimos todos los niveles igual, sin emojis
		fmt.Print(msg.Content)
	})

	err := core.UpdateBinary(handler.Channel(), VersionGoto)
	handler.CloseAndWait()

	cobra.CheckErr(err)
}
