package cmd

import (
	"encoding/json"
	"fmt"
	"goto/src/config"
	"goto/src/gpath"
	"goto/src/utils"
	"os"

	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Do a restore of the goto-paths file",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, _ []string) {

		//Parse all flags
		input, err := cmd.Flags().GetString("input")
		cobra.CheckErr(err)

		//If exists a config backup
		info, err := os.Stat(input)
		if os.IsNotExist(err) {
			cobra.CheckErr(fmt.Errorf("dont have a backup of config file"))
		} else if err != nil {
			cobra.CheckErr(err)
		}

		//If is not a file
		if info.IsDir() {
			cobra.CheckErr(fmt.Errorf("the input can't be a directory"))
		}

		//Read the config backup
		backup, err := os.ReadFile(input)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("cant read the backup of config file: %v", err))
		}

		//Do the unmarshaling of the config backup
		var gpaths []gpath.GotoPath
		if err := json.Unmarshal(backup, &gpaths); err != nil {
			cobra.CheckErr(fmt.Errorf("cant parse the backup of config file"))
		}

		//And re-write the config file with the backup
		cobra.CheckErr(config.SaveGPathsFile(gpaths, utils.GetFilePath(cmd)))

		fmt.Printf("Restore complete in %s\n", utils.GetFilePath(cmd))
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	//Flags
	restoreCmd.Flags().StringP("input", "i", utils.GetDefaultBackupFilePath(), "The ubication of the backup file")
}
