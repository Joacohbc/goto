package cmd

import (
	"fmt"
	"goto/src/gpath"
	"goto/src/utils"
	"os"

	"github.com/bytedance/sonic"
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
	Run: func(cmd *cobra.Command, _ []string) {

		//Parse all flags  (if is not passed, have already a default value)
		input, err := cmd.Flags().GetString("input")
		cobra.CheckErr(err)

		//If exists a config backup
		info, err := os.Stat(input)
		cobra.CheckErr(err)

		//If is not a file
		if info.IsDir() {
			cobra.CheckErr("the input can't be a directory")
		}

		//Read the config backup
		backup, err := os.ReadFile(input)
		if err != nil {
			cobra.CheckErr(fmt.Sprintf("cant read the backup of config file: %v", err))
		}

		//Do the unmarshaling of the config backup
		var gpaths []gpath.GotoPath
		if err := sonic.ConfigFastest.Unmarshal(backup, &gpaths); err != nil {
			cobra.CheckErr("cant parse the backup of config file")
		}

		//And overwrite the config file with the backup
		cobra.CheckErr(gpath.SaveGPathsFile(gpaths, utils.GetFilePath(cmd)))

		fmt.Printf("Restore complete in %s\n", utils.GetFilePath(cmd))
	},
}

func init() {
	RootCmd.AddCommand(RestoreCmd)

	//Flags
	RestoreCmd.Flags().StringP("input", "i", utils.GetDefaultBackupFilePath(), "The ubication of the backup file")
}
