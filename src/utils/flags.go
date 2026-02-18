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
	FlagTemporal     string = "temporal"
)

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

	gpaths, err := LoadGPaths(cmd.Flags().Changed(FlagTemporal))
	cobra.CheckErr(err)

	cobra.CheckErr(gpath.IsValidIndex(len(gpaths), strconv.Itoa(index)))
	return index
}
