package core

import (
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

	var targetIndex int
	var deleted gpath.GotoPath

	if indexArg != -1 {
		if err := gpath.IsValidIndex(len(gpaths), strconv.Itoa(indexArg)); err != nil {
			return nil, err
		}
		targetIndex = indexArg
		deleted = gpaths[targetIndex]
	} else {
		idx, gp, err := findPath(gpaths, pathArg, abbvArg)
		if err != nil {
			return nil, err
		}
		targetIndex = idx
		deleted = *gp
	}

	// Remove the element at targetIndex
	gpaths = append(gpaths[:targetIndex], gpaths[targetIndex+1:]...)

	// Validate again before saving
	if err := gpath.CheckRepeatedItems(gpaths); err != nil {
		return nil, err
	}

	if err := utils.UpdateGPaths(useTemporal, gpaths); err != nil {
		return nil, err
	}

	return &deleted, nil
}
