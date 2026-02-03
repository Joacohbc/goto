package core

import (
	"goto/src/gpath"
	"goto/src/utils"
)

// ListPaths returns the list of goto paths.
func ListPaths(useTemporal bool) ([]gpath.GotoPath, error) {
	return utils.LoadGPaths(useTemporal)
}
