package main

import (
	"fmt"
	"os"
	"time"
)

type Directory struct {
	Path  string `json:"Path"`
	Short string `json:"Short"`
}

type DirectoryTemp struct {
	Directory
	CreateAt time.Time `json:"CreateAt"`
}

type Error struct {
	Error error
	At    time.Time
}

func (e Error) PrintError() {
	fmt.Printf("[%s] Error: %s", e.At, e.Error.Error())
	os.Exit(1)
}
