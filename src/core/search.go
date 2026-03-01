package core

import (
	"fmt"
	"goto/src/gpath"
	"goto/src/utils"
)

// SearchPath searches for a path by Path or Abbreviation.
// Returns the index, the path, and error if not found.
func SearchPath(pathArg, abbvArg string, useTemporal bool) (int, *gpath.GotoPath, error) {
	gpaths, err := utils.LoadGPaths(useTemporal)
	if err != nil {
		return -1, nil, err
	}

	return findPath(gpaths, pathArg, abbvArg)
}

// findPath is an internal helper to find the index and path in a slice of GotoPath.
func findPath(gpaths []gpath.GotoPath, pathArg, abbvArg string) (int, *gpath.GotoPath, error) {
	if pathArg != "" {
		path, err := gpath.ValidPath(pathArg)
		if err != nil {
			return -1, nil, err
		}

		for i := range gpaths {
			if gpaths[i].Path == path {
				return i, &gpaths[i], nil
			}
		}
		return -1, nil, fmt.Errorf("the path \"%s\" doesn't exist in the gpaths-file", path)
	}

	if abbvArg != "" {
		abbv, err := gpath.ValidAbbreviation(abbvArg)
		if err != nil {
			return -1, nil, err
		}

		for i := range gpaths {
			if gpaths[i].Abbreviation == abbv {
				return i, &gpaths[i], nil
			}
		}
		return -1, nil, fmt.Errorf("doesn't exist a path with that abbreviation \"%s\"", abbv)
	}

	return -1, nil, fmt.Errorf("no identifier provided")
}
