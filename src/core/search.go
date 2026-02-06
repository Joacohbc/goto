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

	if pathArg != "" {
		path, err := gpath.ValidPath(pathArg)
		if err != nil {
			return -1, nil, err
		}
		// Note: Original code uses GetPath which uses ValidPathVar.
		// ValidPath returns error if path does not exist.
		// For search, we might want to search even if path doesn't exist on disk?
		// But existing logic validates it. So we stick to it.

		for i := range gpaths {
			if gpaths[i].Path == path {
				return i, &gpaths[i], nil
			}
		}
		return -1, nil, fmt.Errorf("the path \"%s\" doesn't exist in the gpaths-file", path)

	} else {
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
}
