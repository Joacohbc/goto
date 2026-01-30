package gpath

import (
	"fmt"
)

// GotoPath Type
type GotoPath struct {
	Path         string `json:"path"`
	Abbreviation string `json:"abbreviation"`
}

// Return gpath in String format
func (d *GotoPath) String() string {
	return fmt.Sprintf("\"%s\" - %s", d.Path, d.Abbreviation)
}

// This function valid a directory with ValidPathVar() and ValidAbbreviationVar()
func (d GotoPath) Valid() error {

	if _, err := ValidPath(d.Path); err != nil {
		return err
	}

	if _, err := ValidAbbreviation(d.Abbreviation); err != nil {
		return err
	}

	return nil
}
