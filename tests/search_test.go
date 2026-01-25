package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"strings"
	"testing"
)

func TestSearchByAbbreviation(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()

	cmd.AddCmd.Run(c, []string{".", "p1"})

	searchCtx := getTempCmd()
	searchCtx.Flags().StringP(utils.FlagPath, "p", "", "")
	searchCtx.Flags().StringP(utils.FlagAbbreviation, "a", "", "")

	searchCtx.Flags().Set(utils.FlagAbbreviation, "p1")

	output := captureOutput(func() {
		cmd.SearchCmd.Run(searchCtx, []string{})
	})

	if !strings.Contains(output, "p1") {
		t.Errorf("Expected 'p1' in search results, got: %s", output)
	}
}
