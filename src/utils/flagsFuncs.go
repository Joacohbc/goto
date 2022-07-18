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
)

// Check if the flag (key) was passed
func FlagPassed(cmd *cobra.Command, key string) bool {
	return cmd.Flags().Changed(key)
}

// Check if the FlagPath was passed
func PathFlagPassed(cmd *cobra.Command) bool {
	return cmd.Flags().Changed(FlagPath)
}

// Check if the FlagAbbreviation was passed
func AbbvFlagPassed(cmd *cobra.Command) bool {
	return cmd.Flags().Changed(FlagAbbreviation)
}

// Check if the FlagIndex was passed
func IndexFlagPassed(cmd *cobra.Command) bool {
	return cmd.Flags().Changed(FlagIndex)
}

// Check if the FlagPath was passed
func TemporalFlagPassed(cmd *cobra.Command) bool {
	return cmd.Flags().Changed("temporal")
}

// Returns the value of the FlagPath already validated
func GetPath(cmd *cobra.Command) string {
	path, err := cmd.Flags().GetString(FlagPath)
	cobra.CheckErr(err)

	cobra.CheckErr(gpath.ValidPathVar(&path))
	return path
}

// Returns the value of the FlagAbbreviation already validated
func GetAbbreviation(cmd *cobra.Command) string {
	abbv, err := cmd.Flags().GetString(FlagAbbreviation)
	cobra.CheckErr(err)

	cobra.CheckErr(gpath.ValidAbbreviationVar(&abbv))

	return abbv
}

// Returns the value of the FlagIndex flag already validated
func GetIndex(cmd *cobra.Command) int {
	index, err := cmd.Flags().GetInt(FlagIndex)
	cobra.CheckErr(err)

	gpaths := LoadGPaths(cmd)

	cobra.CheckErr(gpath.IsValidIndex(len(gpaths), strconv.Itoa(index)))
	return index
}
