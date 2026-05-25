package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionGoto = "unknown"

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of goto",
	Run:   runVersion,
}

func runVersion(_ *cobra.Command, _ []string) {
	fmt.Println("Goto version is: " + VersionGoto)
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
