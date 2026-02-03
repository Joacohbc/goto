package core

import (
	"fmt"
	"goto/src/gpath"
	"goto/src/utils"
	"strconv"
)

// DeletePath deletes a path identified by path, abbreviation or index.
// Returns the deleted path info or error.
func DeletePath(pathArg, abbvArg string, indexArg int, useTemporal bool) (*gpath.GotoPath, error) {
	gpaths, err := utils.LoadGPaths(useTemporal)
	if err != nil {
		return nil, err
	}

	var deleted gpath.GotoPath
	found := false

	if pathArg != "" {
		path, err := gpath.ValidPath(pathArg)
		if err != nil {
			return nil, err
		}

		for i, gp := range gpaths {
			if gp.Path == path {
				deleted = gp
				gpaths = append(gpaths[:i], gpaths[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("any gpath match with the path \"%s\"", path)
		}

	} else if abbvArg != "" {
		abbv, err := gpath.ValidAbbreviation(abbvArg)
		if err != nil {
			return nil, err
		}

		for i, gp := range gpaths {
			if gp.Abbreviation == abbv {
				deleted = gp
				gpaths = append(gpaths[:i], gpaths[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("any gpath match with the abbreviation \"%s\"", abbv)
		}

	} else if indexArg != -1 {
		if err := gpath.IsValidIndex(len(gpaths), strconv.Itoa(indexArg)); err != nil {
			return nil, err
		}

		deleted = gpaths[indexArg]
		gpaths = append(gpaths[:indexArg], gpaths[indexArg+1:]...)
		found = true

	} else {
		return nil, fmt.Errorf("no identifier provided")
	}

	// Validate again before saving
	if err := gpath.CheckRepeatedItems(gpaths); err != nil {
		return nil, err
	}

	if err := utils.UpdateGPaths(useTemporal, gpaths); err != nil {
		return nil, err
	}

	return &deleted, nil
}
