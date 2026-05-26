package core

import (
	"goto/src/gpath"
	"goto/src/utils"
	"path/filepath"
	"strings"
)

// ResolvePath resolves the target path based on arguments and flags.
func ResolvePath(args []string, onlyDirectory bool, useTemporal bool) (string, error) {
	rawPath := ""
	if len(args) > 0 {
		rawPath = args[0]
	}

	path := filepath.Join(args...)

	// If onlyDirectory flag is passed or the raw argument looks like a path, skip abbreviation resolution
	isExplicitPath := onlyDirectory || 
		strings.HasPrefix(rawPath, ".") || 
		strings.HasPrefix(rawPath, "/") || 
		strings.HasPrefix(rawPath, "~") || 
		strings.Contains(rawPath, "/")

	if isExplicitPath {
		// If only directory flag is passed or is an explicit path, check if is a directory
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
