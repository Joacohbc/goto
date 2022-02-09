package main

import (
	"fmt"
	"os"
	"time"
)

//
// Directory Type
//
type Directory struct {
	Path  string `json:"Path"`
	Short string `json:"Short"`
}

//
// Temporal Directory Type
//
type DirectoryTemp struct {
	Directory
	CreateAt time.Time `json:"CreateAt"`
}

//Print a Error message
func PrintError(e error, previousMessage ...string) {
	fmt.Println("Error: ", e.Error())

	for _, s := range previousMessage {
		fmt.Println(s)
	}

	os.Exit(1)
}

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
