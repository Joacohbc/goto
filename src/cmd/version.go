package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VersionGoto = "2.3.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of goto",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Goto version is: " + VersionGoto)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
