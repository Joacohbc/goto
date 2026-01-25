package tests

import (
	"bytes"
	"goto/src/utils"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

// Helper to reset the temporal file before each test
func resetTempFile(t *testing.T) {
	cmd := getTempCmd()

	path := utils.GetFilePath(cmd)

	// Write minimal valid content (cannot be empty per business logic)
	// We use TempDir as default entry
	tmp := os.TempDir()
	content := `[{"Path":"` + tmp + `","Abbreviation":"default_test_entry"}]`
	err := os.WriteFile(path, []byte(content), 0666)
	if err != nil {
		t.Fatal(err)
	}
}

// Helper to get a command configured with temporal flag
func getTempCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("temporal", "t", false, "")
	_ = cmd.Flags().Set("temporal", "true")
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
