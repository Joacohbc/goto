package cmd

import (
	"fmt"
	"goto/src/core"
	"goto/src/utils"

	"github.com/spf13/cobra"
)

// BackupCmd represents the backup command
var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Do a backup of the goto-paths file",
	Example: `
# Format: goto backup [ -o path ]

# Made a backup of goto-paths in the config directory
goto backup

# If you want to specify the output path
goto backup -o /the/path/file.json.backup
`,
	Args: cobra.ExactArgs(0),
	Run:  runBackup,
}

func runBackup(cmd *cobra.Command, _ []string) {

	//Get output flag (if is not passed, have already a default value)
	output, err := cmd.Flags().GetString("output")
	cobra.CheckErr(err)

	cobra.CheckErr(core.BackupGPaths(output, utils.TemporalFlagPassed(cmd)))
	fmt.Printf("Backup complete from %s\n", utils.GetFilePath(utils.TemporalFlagPassed(cmd)))
}

func init() {
	RootCmd.AddCommand(BackupCmd)

	//Flags
	BackupCmd.Flags().StringP("output", "o", utils.GetDefaultBackupFilePath(), "The backup destination path (must be a file path)")
}
