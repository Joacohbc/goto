package core

import (
	"fmt"
	"goto/src/gpath"
	"goto/src/utils"
	"strconv"
)

const msgPathNotExist = "the Path \"%v\" doesn't exist in the goto-paths file"
const msgAbbvNotExist = "the Abbreviation \"%v\" doesn't exist in the goto-paths file"

// UpdatePath updates a path based on the mode and new value.
func UpdatePath(mode string, pathArg, abbvArg string, indexArg int, newValue string, useTemporal bool) error {
	gpaths, err := utils.LoadGPaths(useTemporal)
	if err != nil {
		return err
	}

	changeIndex := func(inx1, inx2 int) {
		gpaths[inx1], gpaths[inx2] = gpaths[inx2], gpaths[inx1]
	}

	modes := [][]string{
		{"path-path", "pp"}, // 0
		{"path-abbv", "pa"}, // 1
		{"path-indx", "pi"}, // 2
		{"abbv-path", "ap"}, // 3
		{"abbv-abbv", "aa"}, // 4
		{"abbv-indx", "ai"}, // 5
		{"indx-path", "ip"}, // 6
		{"indx-abbv", "ia"}, // 7
		{"indx-indx", "ii"}, // 8
	}

	switch mode {

	//path-path
	case modes[0][0], modes[0][1]:
		path, err := gpath.ValidPath(pathArg)
		if err != nil { return err }

		if err := gpath.ValidPathVar(&newValue); err != nil { return err }

		for i := range gpaths {
			if gpaths[i].Path == path {
				gpaths[i].Path = newValue
				break
			}
			if i == len(gpaths)-1 {
				return fmt.Errorf(msgPathNotExist, path)
			}
		}

	//path-abbv
	case modes[1][0], modes[1][1]:
		path, err := gpath.ValidPath(pathArg)
		if err != nil { return err }

		if err := gpath.ValidAbbreviationVar(&newValue); err != nil { return err }

		for i := range gpaths {
			if gpaths[i].Path == path {
				gpaths[i].Abbreviation = newValue
				break
			}
			if i == len(gpaths)-1 {
				return fmt.Errorf(msgPathNotExist, path)
			}
		}

	//path-indx
	case modes[2][0], modes[2][1]:
		path, err := gpath.ValidPath(pathArg)
		if err != nil { return err }

		if err := gpath.IsValidIndex(len(gpaths), newValue); err != nil { return err }
		n, _ := strconv.Atoi(newValue)

		for i := range gpaths {
			if gpaths[i].Path == path {
				changeIndex(i, n)
				break
			}
			if i == len(gpaths)-1 {
				return fmt.Errorf(msgPathNotExist, path)
			}
		}

	//abbv-path
	case modes[3][0], modes[3][1]:
		abbv, err := gpath.ValidAbbreviation(abbvArg)
		if err != nil { return err }

		if err := gpath.ValidPathVar(&newValue); err != nil { return err }

		for i := range gpaths {
			if gpaths[i].Abbreviation == abbv {
				gpaths[i].Path = newValue
				break
			}
			if i == len(gpaths)-1 {
				return fmt.Errorf(msgAbbvNotExist, abbv)
			}
		}

	//abbv-abbv
	case modes[4][0], modes[4][1]:
		abbv, err := gpath.ValidAbbreviation(abbvArg)
		if err != nil { return err }

		if err := gpath.ValidAbbreviationVar(&newValue); err != nil { return err }

		for i := range gpaths {
			if gpaths[i].Abbreviation == abbv {
				gpaths[i].Abbreviation = newValue
				break
			}
			if i == len(gpaths)-1 {
				return fmt.Errorf(msgAbbvNotExist, abbv)
			}
		}

	//abbv-indx
	case modes[5][0], modes[5][1]:
		abbv, err := gpath.ValidAbbreviation(abbvArg)
		if err != nil { return err }

		if err := gpath.IsValidIndex(len(gpaths), newValue); err != nil { return err }
		n, _ := strconv.Atoi(newValue)

		for i := range gpaths {
			if gpaths[i].Abbreviation == abbv {
				changeIndex(i, n)
				break
			}
			if i == len(gpaths)-1 {
				return fmt.Errorf(msgAbbvNotExist, abbv)
			}
		}

	//indx-path
	case modes[6][0], modes[6][1]:
		indx := indexArg
		if err := gpath.IsValidIndex(len(gpaths), strconv.Itoa(indx)); err != nil { return err }

		if err := gpath.ValidPathVar(&newValue); err != nil { return err }

		gpaths[indx].Path = newValue

	//indx-abbv
	case modes[7][0], modes[7][1]:
		indx := indexArg
		if err := gpath.IsValidIndex(len(gpaths), strconv.Itoa(indx)); err != nil { return err }

		if err := gpath.ValidAbbreviationVar(&newValue); err != nil { return err }

		gpaths[indx].Abbreviation = newValue

	//indx-indx
	case modes[8][0], modes[8][1]:
		indx := indexArg
		if err := gpath.IsValidIndex(len(gpaths), strconv.Itoa(indx)); err != nil { return err }

		if err := gpath.IsValidIndex(len(gpaths), newValue); err != nil { return err }
		n, _ := strconv.Atoi(newValue)

		changeIndex(indx, n)

	default:
		return fmt.Errorf("invalid values of modes to update, use goto --modes")
	}

	return utils.UpdateGPaths(useTemporal, gpaths)
}
