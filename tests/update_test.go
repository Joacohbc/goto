package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"testing"
)

func TestUpdateAbbvAbbv(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()

	cmd.AddCmd.Run(c, []string{".", "oldname"})

	updCmd := getTempCmd()
	updCmd.Flags().StringP(utils.FlagPath, "p", "", "")
	updCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
	updCmd.Flags().IntP(utils.FlagIndex, "i", -1, "")
	updCmd.Flags().StringP("new", "n", "", "")
	updCmd.Flags().BoolP("modes", "m", false, "")

	updCmd.Flags().Set(utils.FlagAbbreviation, "oldname")
	updCmd.Flags().Set("new", "newname")

	cmd.UpdateCmd.Run(updCmd, []string{"abbv-abbv"})

	gpaths := utils.LoadGPaths(c)
	found := false
	for _, gp := range gpaths {
		if gp.Abbreviation == "newname" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find abbreviation 'newname'")
	}
}
