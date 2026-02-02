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
func resetConfigFile(t *testing.T, temporalDir bool) (*cobra.Command, func()) {
	cmd := &cobra.Command{}

	cmd.Flags().BoolP("temporal", "t", false, "")
	if temporalDir {
		_ = cmd.Flags().Set("temporal", "true")
	}

	path := utils.GetFilePath(cmd)
	dir := filepath.Dir(path)
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("Failed to reset config file: %v", err)
	}

	utils.SetupConfigFile()

	return cmd, func() {
		// If path is in a temp dir like /tmp, do not remove the directory, just the file
		if dir == os.TempDir() || filepath.Clean(dir) == "/tmp" {
			_ = os.Remove(path)
		} else {
			if err := os.RemoveAll(dir); err != nil {
				t.Fatalf("Failed to reset config file: %v", err)
			}
		}
	}
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
	cmd.Env = append(os.Environ(), envKey+"=1", utils.TESTING_ENV_VAR+"="+utils.TESTING_ENV_VAR_VALUE)
	err := cmd.Run()

	// Verify that the subprocess failed as expected (exit status 1)
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return // The test passed
	}
	t.Fatalf("process ran successfully (err: %v), but expected exit status 1", err)
}
