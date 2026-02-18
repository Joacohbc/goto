package cmd

import (
	"fmt"
	"goto/src/core"
	"goto/src/utils"

	"github.com/spf13/cobra"
)

// RestoreCmd represents the restore command
var RestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Do a restore of the goto-paths file",
	Example: `
# Format: goto restore [ -t ] [ -i /path/file ]

# Do a restore of goto-paths from a backup in the config directory
goto restore

# If you want to specify the input path
goto restore -i /the/path/file.json.backup
	`,
	Args: cobra.ExactArgs(0),
	Run:  runRestore,
}

func runRestore(cmd *cobra.Command, _ []string) {

	//Parse all flags  (if is not passed, have already a default value)
	input, err := cmd.Flags().GetString("input")
	cobra.CheckErr(err)

	cobra.CheckErr(core.RestoreGPaths(input, cmd.Flags().Changed(utils.FlagTemporal)))

	fmt.Printf("Restore complete in %s\n", utils.GetFilePath(cmd.Flags().Changed(utils.FlagTemporal)))
}

func init() {
	RootCmd.AddCommand(RestoreCmd)

	//Flags
	RestoreCmd.Flags().StringP("input", "i", utils.GetDefaultBackupFilePath(), "The ubication of the backup file")
}
