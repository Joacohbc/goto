package core

import (
	"goto/src/gpath"
	"goto/src/utils"
)

// AddPath adds a new path to the goto-paths file.
// It validates the input arguments before adding.
func AddPath(pathArg, abbvArg string, useTemporal bool) error {
	gpaths, err := utils.LoadGPaths(useTemporal)
	if err != nil {
		return err
	}

	path, err := gpath.ValidPath(pathArg)
	if err != nil {
		return err
	}

	abbv, err := gpath.ValidAbbreviation(abbvArg)
	if err != nil {
		return err
	}

	gpaths = append(gpaths, gpath.GotoPath{
		Path:         path,
		Abbreviation: abbv,
	})

	// Check for duplicates is handled by UpdateGPaths -> SaveGPathsFile -> CheckRepeatedItems
	// But CheckRepeatedItems requires the array.
	// Logic is consistent.

	return utils.UpdateGPaths(useTemporal, gpaths)
}
