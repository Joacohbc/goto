package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//
// Directory Type
//

//const DirectoryTypeSeparator string = ","

type Directory struct {
	Path         string `json:"path"`
	Abbreviation string `json:"abbreviation"`
}

//Return directory in String format
func (d *Directory) String() string {
	return fmt.Sprintf("Path: \"%s\" - Abbreviation: \"%s\"", d.Path, d.Abbreviation)
}

// This function valid:
// - The Path and the Abbreviation don't be empty
// - Clean the Path
// - That Abbreviation is not a number letter
// - That Path exist and is a Directory
func (d *Directory) ValidDirectory() error {

	d.Path = strings.TrimSpace(d.Path)
	d.Path = filepath.Clean(d.Path)
	d.Abbreviation = strings.TrimSpace(d.Abbreviation)

	if len(d.Path) < 1 || len(d.Abbreviation) < 1 {
		return fmt.Errorf("path and abbreviation can't be empty or be blank space")
	}

	if _, err := strconv.Atoi(d.Abbreviation); err == nil {
		return fmt.Errorf("the Abbreviation can't be a number'")
	}

	info, err := os.Stat(d.Path)

	if os.IsNotExist(err) {
		return fmt.Errorf("file \"%s\" don't exist", d.Path)
	}

	if err != nil {
		return fmt.Errorf("error to get info of \"%s\": %s", d.Path, err.Error())
	}

	if !info.IsDir() {
		return fmt.Errorf("the file \"%s\" is not a directory", d.Path)
	}

	return nil
}

//Parse a string to Directory (using the DirectoryTypeSeparator to split)
//func ToDirectory(s string) (Directory, error) {
//
//	s = strings.TrimSpace(s)
//	//Use Trim to avoid blank spaces Abbreviations
//	if len(s) < 3 {
//		return Directory{}, fmt.Errorf("need at least 3 characters for create a path")
//	}
//
//	args := strings.Split(s, DirectoryTypeSeparator)
//
//	if len(args) != 2 {
//		return Directory{}, fmt.Errorf("need 2 args to make a create a new path")
//	}
//
//	dir := Directory{Path: args[0], Abbreviation: args[1]}
//	if err := dir.ValidDirectory(); err != nil {
//		return Directory{}, err
//	}
//
//	return dir, nil
//}

//Check that the any directory has the same Path or same Abbreviation that other
func ValidArray(dirs []Directory) error {

	if len(dirs) == 0 {
		return fmt.Errorf("the config file is empty")
	}

	for i, dir := range dirs {

		//Check that 2 Path don't have the same abbreviation, where the indexs are diferents
		//(With diferent index beacause obviously the same index have the same abbreviation and the same path)
		for indexRepeated, dirRepeated := range dirs {

			//If have the same path and is not the same index
			if (dir.Path == dirRepeated.Path) && (i != indexRepeated) {
				return fmt.Errorf("the path: \"%v\" already exists", dirRepeated.Path)
			}

			//If have the same path and is not the same index
			if (dir.Abbreviation == dirRepeated.Abbreviation) && (i != indexRepeated) {
				return fmt.Errorf("the path: \"%v\"(index %v) have the same abbreviation that \"%v\"(index %v)", dir.Path, i, dirRepeated.Path, indexRepeated)
			}
		}
	}

	return nil
}

/*

//
// Temporal Directory Type
//
type DirectoryTemp struct {
	Directory
	CreateAt time.Time `json:"CreateAt"`
}

*/

/*
//
// Error Type
//
type Error struct {
	Error error
	At    time.Time
}

//Use interface{} to accept "string" type and "error" type
//Create a new object of Error type
func NewError(e interface{}) Error {

	//If e == nill return Error with nil
	//because fmt.Errorf(fmt.Sprint(nil)) is not the same that nil
	if e == nil {
		return Error{
			Error: nil,
			At:    time.Now(),
		}
	}

	return Error{
		Error: fmt.Errorf(fmt.Sprint(e)),
		At:    time.Now(),
	}
}

func (e Error) IsNil() bool {
	return e.Error == nil
}

func (e Error) IsNotNil() bool {
	return e.Error != nil
}

func (e Error) Print() {
	fmt.Printf("Error: %s \n", e.Error.Error())
	os.Exit(1)
}

*/
