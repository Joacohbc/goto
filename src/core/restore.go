package core

import (
	"bufio"
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

	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cant open the backup of config file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var gpaths []gpath.GotoPath
	if err := sonic.ConfigFastest.NewDecoder(reader).Decode(&gpaths); err != nil {
		return fmt.Errorf("cant parse the backup of config file: %v", err)
	}

	return utils.UpdateGPaths(useTemporal, gpaths)
}
