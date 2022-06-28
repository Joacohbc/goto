package utils

import (
	"goto/src/gpath"
	"strconv"

	"github.com/spf13/cobra"
)

const (
	FlagPath         string = "path"
	FlagAbbreviation string = "abbv"
	FlagIndex        string = "indx"
	FlagCurretDir    string = "current"
)

// Check if the flag (key) was passed
func FlagPassed(cmd *cobra.Command, key string) bool {
	return cmd.Flags().Changed(key)
}

// Check if the Path flag  was passed
func PathFlagPassed(cmd *cobra.Command) bool {
	return cmd.Flags().Changed(FlagPath)
}

// Check if the Abbreviation flag  was passed
func AbbvFlagPassed(cmd *cobra.Command) bool {
	return cmd.Flags().Changed(FlagAbbreviation)
}

// Check if the Index flag  was passed
func IndexFlagPassed(cmd *cobra.Command) bool {
	return cmd.Flags().Changed(FlagIndex)
}

// Check if the FlagCurretDir flag was passed
func CurrentDirFlagPassed(cmd *cobra.Command) bool {
	return cmd.Flags().Changed(FlagCurretDir)
}

// Returns the FlagPath flag already valided and checking the FlagCurretDir flag
// In case of any error, the use cobra.CheckErr() to print and exit
func GetPath(cmd *cobra.Command) string {
	path, err := cmd.Flags().GetString(FlagPath)
	cobra.CheckErr(err)

	//If current is passed, overwrite the path to current directory
	if FlagPassed(cmd, FlagCurretDir) {
		path = GetCurrentDirectory()
	}

	cobra.CheckErr(gpath.ValidPathVar(&path))
	return path
}

// Returns the FlagAbbreviation flag already valided
func GetAbbreviation(cmd *cobra.Command) string {
	abbv, err := cmd.Flags().GetString(FlagAbbreviation)
	cobra.CheckErr(err)

	cobra.CheckErr(gpath.ValidAbbreviationVar(&abbv))

	return abbv
}

// Returns the FlagIndex flag already valided
func GetIndex(cmd *cobra.Command) int {
	index, err := cmd.Flags().GetInt(FlagIndex)
	cobra.CheckErr(err)

	gpaths := LoadGPaths(cmd)

	cobra.CheckErr(gpath.IsValidIndex(len(gpaths), strconv.Itoa(index)))
	return index
}
