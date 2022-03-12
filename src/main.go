package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const versionMessage string = "1.10" //Version

var (
	Help             bool   // -h -help
	Version          bool   // -v -version
	List             bool   // -list
	ConfFilePath     bool   // config-path
	PathQuotes       bool   // -q
	AddPath          bool   // -add
	DelPath          bool   // -del
	ModifyPath       bool   // -modify
	DoBackup         bool   // -backup
	DoRestore        bool   // -restore
	PathFlag         string // -path
	AbbreviationFlag string // -abbv
)

// Check if the arg is a valid abbreviation, and if a
// number check that is valid index in the config file.
//
// If the arg is not a number or abbreviation, check if
// it is a directory, if not return error
func ArgIsAbbvOrNumberOrDir(arg string) (string, error) {

	//Load the config file in memory
	var directories []Directory
	if err := loadConfigFile(&directories); err != nil {
		return "", err //In case of error, print the error and exit
	}

	//Check if path is number
	if pathNumber, err := strconv.Atoi(arg); err == nil {

		//If the path is over the max index return error
		if pathNumber < 0 || pathNumber > len(directories)-1 {
			//In case of error, print the error and exit
			return "", fmt.Errorf("the number is invalid(should be: 0-" + strconv.Itoa(len(directories)-1) + "), check config file")
		}

		for i, dir := range directories {
			if pathNumber == i {
				return dir.Path, nil //In case of correct pathNumber, print and exit
			}
		}
	}

	//If not a number, check if is an abbreviation
	for _, dir := range directories {
		if arg == dir.Abbreviation {
			return dir.Path, nil //In case of correct abbreviation, print and exit
		}
	}

	//If it is neither a number nor an abbreviation, check if is exists file
	fileInfo, err := os.Stat(arg)
	if err == nil {
		//If exists, check if it's a directory
		if fileInfo.IsDir() {
			return filepath.Clean(arg), nil
		}

		//If not a directory
		return "", fmt.Errorf("the path \"%s\" is not a directory", arg)
	}

	//If the path not exists
	if os.IsNotExist(err) {
		return "", fmt.Errorf("the path \"%s\" is not exists", arg)
	}

	return "", err
}

func init() {

	//Info flags
	flag.BoolVar(&Help, "h", false, "Print help message")
	flag.BoolVar(&Help, "help", false, "Print help message")
	flag.BoolVar(&Version, "v", false, "Print version")
	flag.BoolVar(&Version, "version", false, "Print version")
	flag.BoolVar(&ConfFilePath, "config-path", false, "Print path of the config.json")

	//Path and Abbreviations flags
	flag.StringVar(&PathFlag, "path", "", "The Path to add, delete or modify actions")
	flag.StringVar(&AbbreviationFlag, "abbv", "", "The Abbreviation to add, delete or modify actions")

	//Add, Delete, Modify and List Directory flags
	flag.BoolVar(&AddPath, "add", false, "Add a new Path with a Abbreviation")
	flag.BoolVar(&DelPath, "del", false, "Delete a the path")
	flag.BoolVar(&ModifyPath, "modify", false, "Modify a Abbreviation from the Path")
	flag.BoolVar(&List, "list", false, "Print all path with abbreviations")

	//Backup and Restore
	flag.BoolVar(&DoBackup, "backup", false, "Do backup of the config file")
	flag.BoolVar(&DoRestore, "restore", false, "Do restore of the config file")

	//Other funcs flags
	flag.BoolVar(&PathQuotes, "q", false, "Print the path (abbreviation or index) or directory with quotes")

	flag.Parse()
}

func main() {

	//Create the config file
	if err := createConfigFile(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	//If the help argument is passed, print help message
	if Help {
		flag.Usage()
		return
	}

	//If the version argument is passed, print version message
	if Version {
		fmt.Printf("Version of goto: %v", versionMessage)
		return
	}

	//Print the path of the config file
	if ConfFilePath {
		fmt.Println(ConfigFile)
		return
	}

	//Do a backup of the config file
	if DoBackup {
		if err := doBackup(); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Backup complete")
		return
	}

	//Do a restore of the backup of config file
	if DoRestore {
		if err := doRestore(); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Restore complete")
		return
	}

	//If the list argument is passed, print the list of the config file
	if List {
		var directoriesToList []Directory
		if err := loadConfigFile(&directoriesToList); err != nil {
			fmt.Println("Error:", err)
			return
		}

		for i, dir := range directoriesToList {
			fmt.Printf("%v - %s\n", i, dir.String())
		}
		return
	}

	//If the quotes argument is passed, print the dir with quotes
	if PathQuotes {
		//Check if "arg" is an abbreviation or a number index
		path, err := ArgIsAbbvOrNumberOrDir(flag.Arg(0))
		if err != nil {
			fmt.Println("Error:", err)
			return

		}
		fmt.Println("\"" + path + "\"")
		return
	}

	//If the add argument is passed, use func add
	if AddPath {

		dir := Directory{Path: PathFlag, Abbreviation: AbbreviationFlag}
		if err := dir.ValidDirectory(); err != nil {
			fmt.Println("Error:", err)
			fmt.Println("The changes were not applied")
			return
		}

		if err := addNewPaths(dir); err != nil {
			fmt.Println("Error:", err)
			fmt.Println("The changes were not applied")
			return
		}

		fmt.Println("The changes were applied successfully")
		return
	}

	//If the del argument is passed, use func del
	if DelPath {

		dir := Directory{Path: PathFlag, Abbreviation: "not-necessary"}
		if err := dir.ValidDirectory(); err != nil {
			fmt.Println("Error: ", err)
			fmt.Println("The changes were not applied")
			return
		}

		if err := delPaths(PathFlag); err != nil {
			fmt.Println("Error:", err)
			fmt.Println("The changes were not applied")
			return
		}

		fmt.Println("The changes were applied successfully")
		return
	}

	//If the modify argument is passed, use func modify
	if ModifyPath {

		dir := Directory{Path: PathFlag, Abbreviation: AbbreviationFlag}
		if err := dir.ValidDirectory(); err != nil {
			fmt.Println("Error:", err)
			fmt.Println("The changes were not applied")
			return
		}

		if err := modPaths(dir.Path, dir.Abbreviation); err != nil {
			fmt.Println("Error:", err)
			fmt.Println("The changes were not applied")
			return
		}

		fmt.Println("The changes were applied successfully")
		return
	}

	//Where the first argument will be stored
	var arg string = flag.Arg(0)

	//Check if "arg" is an abbreviation, a number index or directory
	path, err := ArgIsAbbvOrNumberOrDir(arg)
	if err != nil {
		fmt.Println("Error:", err)
		return

	}
	fmt.Println(path)

	//If the code is here, it means that the arg is invalid
	//fmt.Println("invalid argument/s")
}
