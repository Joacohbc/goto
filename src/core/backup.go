package core

import (
	"fmt"
	"goto/src/gpath"
	"goto/src/utils"
	"os"
)

// BackupGPaths backs up the current goto paths to outputPath.
func BackupGPaths(outputPath string, useTemporal bool) error {
	gpaths, err := utils.LoadGPaths(useTemporal)
	if err != nil {
		return err
	}

	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("the file \"%s\" already exists", outputPath)
	}

	return gpath.SaveGPathsFile(gpaths, outputPath)
}
