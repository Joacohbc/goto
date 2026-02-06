package core

import (
	"goto/src/gpath"
	"goto/src/utils"
)

// ValidatePaths checks if all paths in the config file are valid.
func ValidatePaths(useTemporal bool) error {
	gpaths, err := utils.LoadGPaths(useTemporal)
	if err != nil {
		return err
	}

	for _, g := range gpaths {
		if err := g.Valid(); err != nil {
			return err
		}
	}

	return gpath.CheckRepeatedItems(gpaths)
}
