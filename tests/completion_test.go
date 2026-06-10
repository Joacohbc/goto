package tests

import (
	"goto/src/cmd"
	"goto/src/gpath"
	"goto/src/utils"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCompleteAbbreviationFlag(t *testing.T) {
	dummyCmd, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Load gpaths and add test items
	gpathsList, err := utils.LoadGPaths(false)
	if err != nil {
		t.Fatalf("failed to load gpaths: %v", err)
	}

	gpathsList = append(gpathsList, gpath.GotoPath{
		Path:         "/home/devuser/test1",
		Abbreviation: "abc",
	}, gpath.GotoPath{
		Path:         "/home/devuser/test2",
		Abbreviation: "def",
	})
	if err := utils.UpdateGPaths(false, gpathsList); err != nil {
		t.Fatalf("failed to update gpaths: %v", err)
	}

	// Test CompleteAbbreviationFlag without prefix
	completions, directive := cmd.CompleteAbbreviationFlag(dummyCmd, nil, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected directive %v, got %v", cobra.ShellCompDirectiveNoFileComp, directive)
	}

	foundAbc := false
	foundDef := false
	for _, c := range completions {
		if strings.HasPrefix(c, "abc\t") {
			foundAbc = true
		}
		if strings.HasPrefix(c, "def\t") {
			foundDef = true
		}
	}
	if !foundAbc || !foundDef {
		t.Errorf("expected to find abc and def, got: %v", completions)
	}

	// Test CompleteAbbreviationFlag with matching prefix
	completions, directive = cmd.CompleteAbbreviationFlag(dummyCmd, nil, "ab")
	foundAbc = false
	foundDef = false
	for _, c := range completions {
		if strings.HasPrefix(c, "abc\t") {
			foundAbc = true
		}
		if strings.HasPrefix(c, "def\t") {
			foundDef = true
		}
	}
	if !foundAbc {
		t.Error("expected to find abc with prefix 'ab'")
	}
	if foundDef {
		t.Error("expected not to find def with prefix 'ab'")
	}

	// Test CompleteAbbreviationFlag with non-matching prefix
	completions, _ = cmd.CompleteAbbreviationFlag(dummyCmd, nil, "xyz")
	if len(completions) != 0 {
		t.Errorf("expected 0 completions for non-matching prefix, got: %v", completions)
	}
}

func TestRootCmdValidArgsFunction(t *testing.T) {
	dummyCmd, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Load gpaths and add test items
	gpathsList, err := utils.LoadGPaths(false)
	if err != nil {
		t.Fatalf("failed to load gpaths: %v", err)
	}

	gpathsList = append(gpathsList, gpath.GotoPath{
		Path:         "/home/devuser/test1",
		Abbreviation: "abc",
	})
	if err := utils.UpdateGPaths(false, gpathsList); err != nil {
		t.Fatalf("failed to update gpaths: %v", err)
	}

	validArgsFunc := cmd.RootCmd.ValidArgsFunction

	// Test with existing args (should suggest nothing)
	completions, directive := validArgsFunc(dummyCmd, []string{"arg1"}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp || len(completions) != 0 {
		t.Errorf("expected empty completions with existing args, got completions=%v, directive=%v", completions, directive)
	}

	// Test with path prefix "./" (should trigger directory filter)
	completions, directive = validArgsFunc(dummyCmd, nil, "./")
	if directive != cobra.ShellCompDirectiveFilterDirs || len(completions) != 0 {
		t.Errorf("expected ShellCompDirectiveFilterDirs, got directive=%v, completions=%v", directive, completions)
	}

	// Test with path prefix "/" (should trigger directory filter)
	completions, directive = validArgsFunc(dummyCmd, nil, "/var")
	if directive != cobra.ShellCompDirectiveFilterDirs || len(completions) != 0 {
		t.Errorf("expected ShellCompDirectiveFilterDirs, got directive=%v, completions=%v", directive, completions)
	}

	// Test with normal abbreviation prefix
	completions, directive = validArgsFunc(dummyCmd, nil, "ab")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected directive %v, got %v", cobra.ShellCompDirectiveNoFileComp, directive)
	}
	foundAbc := false
	for _, c := range completions {
		if strings.HasPrefix(c, "abc\t") {
			foundAbc = true
		}
	}
	if !foundAbc {
		t.Errorf("expected to find abc, got: %v", completions)
	}

	// Test that directories of the current directory are suggested alongside the abbreviations
	tempDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(tempDir, "abdir"), 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}
	if err := os.Mkdir(filepath.Join(tempDir, ".hidden"), 0755); err != nil {
		t.Fatalf("failed to create hidden test directory: %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	}()

	completions, directive = validArgsFunc(dummyCmd, nil, "ab")
	expectedDirective := cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
	if directive != expectedDirective {
		t.Errorf("expected directive %v, got %v", expectedDirective, directive)
	}
	foundAbc = false
	foundDir := false
	for _, c := range completions {
		if strings.HasPrefix(c, "abc\t") {
			foundAbc = true
		}
		if c == "abdir"+string(os.PathSeparator) {
			foundDir = true
		}
	}
	if !foundAbc || !foundDir {
		t.Errorf("expected to find abc and abdir%c, got: %v", os.PathSeparator, completions)
	}

	// Hidden directories should not be suggested
	completions, _ = validArgsFunc(dummyCmd, nil, "")
	for _, c := range completions {
		if strings.HasPrefix(c, ".hidden") {
			t.Errorf("expected hidden directories to be excluded, got: %v", completions)
		}
	}
}

func TestAddCmdValidArgsFunction(t *testing.T) {
	validArgsFunc := cmd.AddCmd.ValidArgsFunction

	// First argument should complete directories
	completions, directive := validArgsFunc(cmd.AddCmd, nil, "")
	if directive != cobra.ShellCompDirectiveFilterDirs || len(completions) != 0 {
		t.Errorf("expected directory filter for first arg, got directive=%v, completions=%v", directive, completions)
	}

	// Subsequent arguments should not complete files
	completions, directive = validArgsFunc(cmd.AddCmd, []string{"path1"}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp || len(completions) != 0 {
		t.Errorf("expected no file comp for second arg, got directive=%v, completions=%v", directive, completions)
	}
}
