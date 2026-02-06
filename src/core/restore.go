package core

import (
	"fmt"
	"goto/src/gpath"
	"goto/src/utils"
	"os"

	"github.com/bytedance/sonic"
)

// RestoreGPaths restores goto paths from inputPath.
func RestoreGPaths(inputPath string, useTemporal bool) error {
	info, err := os.Stat(inputPath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("the input can't be a directory")
	}

	backup, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("cant read the backup of config file: %v", err)
	}

	var gpaths []gpath.GotoPath
	if err := sonic.ConfigFastest.Unmarshal(backup, &gpaths); err != nil {
		return fmt.Errorf("cant parse the backup of config file: %v", err)
	}

	return utils.UpdateGPaths(useTemporal, gpaths)
}
