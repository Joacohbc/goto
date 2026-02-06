package core

import (
	"goto/src/gpath"
	"goto/src/utils"
	"path/filepath"
)

// ResolvePath resolves the target path based on arguments and flags.
func ResolvePath(args []string, onlyDirectory bool, useTemporal bool) (string, error) {
	path := filepath.Join(args...)

	if onlyDirectory {
		// If only directory flag is passed, check if is a directory
		if err := gpath.ValidPathVar(&path); err != nil {
			return "", err
		}
	} else {
		// Load the config file
		gpathsList, err := utils.LoadGPaths(useTemporal)
		if err != nil {
			return "", err
		}

		// Check if is a index or an abbreviation
		var isIndexOrAbbv bool
		path, isIndexOrAbbv = gpath.GetPathFromIndexOrAbbreviation(gpathsList, path)

		// If it is not, check if is a directory
		if !isIndexOrAbbv {
			if err := gpath.ValidPathVar(&path); err != nil {
				return "", err
			}
		}
	}
	return path, nil
}
