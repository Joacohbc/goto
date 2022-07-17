package cmd

import (
	"fmt"
	"goto/src/config"
	"goto/src/utils"
	"os"

	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Do a backup of the goto-paths file",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, _ []string) {

		//Read the config file
		gpaths := utils.LoadGPaths(cmd)

		//Get output flag
		output, err := cmd.Flags().GetString("output")
		cobra.CheckErr(err)

		//Check if exists
		if _, err := os.Stat(output); err == nil {
			cobra.CheckErr(fmt.Errorf("the file \"%s\" already exists", output))
		}

		cobra.CheckErr(config.SaveGPathsFile(gpaths, output))
		fmt.Printf("Backup complete from %s\n", utils.GetFilePath(cmd))
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	//Flags
	backupCmd.Flags().StringP("output", "o", utils.GetDefaultBackupFilePath(), "The backup destination path (must be a file path)")
}
