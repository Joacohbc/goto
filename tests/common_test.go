package tests

import (
	"bytes"
	"goto/src/utils"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// Helper to reset the temporal file before each test
func resetConfigFile(t *testing.T) (*cobra.Command, func()) {
	cmd := getCmd()
	path := utils.GetFilePath(cmd)

	return cmd, func() {
		if err := os.RemoveAll(filepath.Dir(path)); err != nil {
			t.Fatalf("Failed to reset config file: %v", err)
		}
	}
}

// Helper to get a command configured with temporal flag
func getCmd() *cobra.Command {
	cmd := &cobra.Command{}
	// TODO: This method should set the temporal flag to true/false by passing argument
	// cmd.Flags().BoolP("temporal", "t", false, "")
	// _ = cmd.Flags().Set("temporal", "true")
	return cmd
}

// Helper to capture stdout
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// RunExpectedExit runs the logical part of a test in a subprocess and asserts that it exits with status 1.
// testName: The name of the test function to run (e.g. "TestAddPathRepeated")
// envKey: The environment variable key used to trigger the subprocess logic (e.g. "TEST_SUBPROCESS")
// The application uses os.Exit(1) (via cobra.CheckErr) when it encounters an error, such as a duplicate path.
// Standard panic/recover mechanisms cannot catch os.Exit, so we must run this test case in a subprocess.
func RunExpectedExit(t *testing.T, testName string, envKey string) {
	cmd := exec.Command(os.Args[0], "-test.run="+testName)
	cmd.Env = append(os.Environ(), envKey+"=1")
	err := cmd.Run()

	// Verify that the subprocess failed as expected (exit status 1)
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return // The test passed
	}
	t.Fatalf("process ran successfully (err: %v), but expected exit status 1", err)
}
