package tests

import (
	"goto/src/utils"
	"testing"

	"github.com/spf13/cobra"
)

func TestTemporalFlagPassed(t *testing.T) {
	// Case 1: Flag not set
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("temporal", "t", false, "")
	if utils.TemporalFlagPassed(cmd) {
		t.Error("Expected TemporalFlagPassed to be false when flag is not set")
	}

	// Case 2: Flag set
	cmd = &cobra.Command{}
	cmd.Flags().BoolP("temporal", "t", false, "")
	_ = cmd.Flags().Set("temporal", "true")
	if !utils.TemporalFlagPassed(cmd) {
		t.Error("Expected TemporalFlagPassed to be true when flag is set")
	}
}

func TestFlagPassed(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("myflag", "m", false, "")

	if utils.FlagPassed(cmd, "myflag") {
		t.Error("Expected FlagPassed to be false when flag is not changed")
	}

	_ = cmd.Flags().Set("myflag", "true")
	if !utils.FlagPassed(cmd, "myflag") {
		t.Error("Expected FlagPassed to be true when flag is changed")
	}
}

func TestOtherFlagsPassed(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("path", "p", "", "")
	cmd.Flags().StringP("abbv", "a", "", "")
	cmd.Flags().IntP("indx", "i", -1, "")

	if utils.PathFlagPassed(cmd) {
		t.Error("Expected PathFlagPassed to be false")
	}
	if utils.AbbreviationFlagPassed(cmd) {
		t.Error("Expected AbbreviationFlagPassed to be false")
	}
	if utils.IndexFlagPassed(cmd) {
		t.Error("Expected IndexFlagPassed to be false")
	}

	_ = cmd.Flags().Set("path", "/some/path")
	_ = cmd.Flags().Set("abbv", "someabbv")
	_ = cmd.Flags().Set("indx", "1")

	if !utils.PathFlagPassed(cmd) {
		t.Error("Expected PathFlagPassed to be true")
	}
	if !utils.AbbreviationFlagPassed(cmd) {
		t.Error("Expected AbbreviationFlagPassed to be true")
	}
	if !utils.IndexFlagPassed(cmd) {
		t.Error("Expected IndexFlagPassed to be true")
	}
}
