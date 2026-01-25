package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"os"
	"path/filepath"
	"testing"
)

func TestRestore(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()

	// Create backup file manually
	backupFile := filepath.Join(os.TempDir(), "goto-test-restore.json")
	// Use absolute path for Path field just in case validation runs (though restore usually just unmarshals)
	// But SaveGPathsFile calls CheckRepeatedItems.
	// We need a valid path.
	absPath, _ := filepath.Abs(".")

	// Manually construct JSON to avoid importing gpath struct if we don't have to
	// But we can use gpath.GotoPath if exported.
	// Or raw json.
	// JSON format: [{"Path":"...","Abbreviation":"..."}]
	content := `[{"Path":"` + absPath + `","Abbreviation":"restored"}]`

	if err := os.WriteFile(backupFile, []byte(content), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(backupFile)

	rstCmd := getTempCmd()
	rstCmd.Flags().StringP("input", "i", "", "")
	rstCmd.Flags().Set("input", backupFile)

	captureOutput(func() {
		cmd.RestoreCmd.Run(rstCmd, []string{})
	})

	gpaths := utils.LoadGPaths(c)
	if len(gpaths) == 0 {
		t.Fatal("Restore failed, gpaths empty")
	}
	if gpaths[0].Abbreviation != "restored" {
		t.Errorf("Expected abbreviation 'restored', got '%s'", gpaths[0].Abbreviation)
	}
}
