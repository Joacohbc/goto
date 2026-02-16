package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VersionGoto = "2.4.9"

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
