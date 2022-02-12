package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const versionMessage string = "1.6" //Version

var (
	help       bool
	version    bool
	list       bool
	pathQuotes bool
	addPath    string
	delPath    string
	modifyPath string
)

func ArgIsDir(arg string) (string, error) {

	fileInfo, err := os.Stat(arg)

	if err == nil {
		//If it's a directory
		if fileInfo.IsDir() {
			return arg, nil //If a file and is a directory, print and exit

		} else {
			return arg, fmt.Errorf("the path is not a directory") //If the file is a directory, print error and exit
		}

	} else {
		return "", nil //If not a file return nil(because is not a dir), continue
	}
}

func ArgIsShortOrNumber(arg string) (string, error) {

	//Load the config file in memory
	var directories []Directory
	err := loadConfigFile(&directories)

	if err != nil {
		return "", err //In case of error, print the error and exit
	}

	//Check if path is number
	if pathNumber, err := strconv.Atoi(arg); err == nil {

		for i, dir := range directories {

			if pathNumber == i {
				return dir.Path, nil //In case of correct pathNumber, print and exit
			}
		}

		//In case of error, print the error and exit
		return "", fmt.Errorf("the number is invalid(should be: 0-" + strconv.Itoa(len(directories)-1) + "), check config file")

	} else { //If it isn't a number
		for _, dir := range directories {
			if arg == dir.Short {
				return dir.Path, nil //In case of correct abbreviation, print and exit
			}
		}
	}

	return "", nil //In case of args is not a number or a valid abbreviation, continue
}

func helpMessage() string {
	helpMessage := `Three ways to use it, with abbreviations(config file), numbers(index of config file) and paths:

-Abbreviations= "goto <abbreviation>"
-Number="goto <number-of-the-index>"
-Path="goto <path>"

Path of config file: 
`
	return helpMessage + ConfigFile
}

func init() {
	flag.BoolVar(&help, "h", false, "Print help message")
	flag.BoolVar(&help, "help", false, "Print help message")

	flag.BoolVar(&version, "v", false, "Print version")
	flag.BoolVar(&version, "version", false, "Print version")

	flag.BoolVar(&list, "l", false, "Print all path with abbreviations")
	flag.BoolVar(&list, "list", false, "Print all path with abbreviations")

	flag.BoolVar(&pathQuotes, "q", false, "Print the path with quotes: -quotes=[Path/Short/Dir]")
	flag.BoolVar(&pathQuotes, "quotes", false, "Print the path with quotes: -quotes=[Path/Short/Dir]")

	flag.StringVar(&addPath, "a", "", "Add a new path use: -add=[New Path],[New Short]")
	flag.StringVar(&addPath, "add", "", "Add a new path use: -add=[New Path],[New Short]")

	flag.StringVar(&delPath, "d", "", "Delete a path use: --del=[Path to Del]")
	flag.StringVar(&delPath, "del", "", "Delete a path use: --del=[Path to Del]")

	flag.StringVar(&modifyPath, "m", "", "Modify a path: -modif=[Path],[New Short]")
	flag.StringVar(&modifyPath, "modify", "", "Modify a path: -modif=[Path],[New Short]")

	flag.Parse()
}

func main() {

	//Create the config file
	if err := createConfigFile(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	//If the help argument is passed, print help message
	if help {
		fmt.Println(helpMessage())
		return
	}

	//If the version argument is passed, print version message
	if version {
		fmt.Printf("Version of goto: %v", versionMessage)
		return
	}

	//If the list argument is passed, print the list of the config file
	if list {
		var directoriesToList []Directory
		if err := loadConfigFile(&directoriesToList); err != nil {
			fmt.Println("Error:", err)
			return
		}

		for i, dir := range directoriesToList {
			fmt.Printf("%v- Path: \"%v\", Short: \"%v\" \n", i, dir.Path, dir.Short)
		}
		return
	}

	//If the quotes argument is passed, print the dir with quotes
	if pathQuotes {

		//If exists like a Directory
		dir, err := ArgIsDir(flag.Arg(0))
		if err != nil {
			fmt.Println("Error:", err)
			return
		} else if len(dir) != 0 {
			fmt.Println("\"" + dir + "\"")
			return
		}

		//Check if "arg" is an abbreviation or a number index
		path, err := ArgIsShortOrNumber(flag.Arg(0))
		if err != nil {
			fmt.Println("Error:", err)
			return

		} else if len(path) != 0 {
			fmt.Println("\"" + path + "\"")
			return
		}
	}

	//If the add argument is passed, use func add
	if len(addPath) != 0 {

		args := strings.Split(addPath, ",")

		if len(args) != 2 {
			fmt.Println("Error: bad format of --add")
			fmt.Println(helpMessage())
			return
		}

		if len(args[0]) == 0 || len(args[1]) == 0 {
			fmt.Println("Error: path and abbreviation can't be blank spaces")
			return
		}

		dir := Directory{Path: args[0], Short: args[1]}

		if err := addNewPaths(dir); err != nil {
			fmt.Println("Error:", err)
			fmt.Println("The changes were not applied")
			return
		}

		fmt.Println("The changes were applied successfully")
		return
	}

	//If the del argument is passed, use func del
	if len(delPath) != 0 {

		if len(delPath) == 0 {
			fmt.Println("Error: path can't be blank spaces")
			return
		}

		if err := delPaths(delPath); err != nil {
			fmt.Println("Error:", err)
			fmt.Println("The changes were not applied")
			return
		}

		fmt.Println("The changes were applied successfully")
		return
	}

	//If the modify argument is passed, use func modify
	if len(modifyPath) != 0 {

		args := strings.Split(modifyPath, ",")

		if len(args) != 2 {
			fmt.Println("Error: bad format of --modify")
			fmt.Println(helpMessage())
			return
		}

		if len(args[0]) == 0 || len(args[1]) == 0 {
			fmt.Println("Error: path and abbreviation can't be blank spaces")
			return
		}

		if err := modPaths(args[0], args[1]); err != nil {
			fmt.Println("Error:", err)
			fmt.Println("The changes were not applied")
			return
		}

		fmt.Println("The changes were applied successfully")
		return
	}

	//Where the first argument will be stored
	var arg string = flag.Arg(0)

	//Check if "arg" is an abbreviation or a number index
	path, err := ArgIsShortOrNumber(arg)
	if err != nil {
		fmt.Println("Error:", err)
		return

	} else if len(path) != 0 {
		fmt.Println(path)
		return
	}

	//If exists like a Directory
	dir, err := ArgIsDir(arg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	} else if len(dir) != 0 {
		fmt.Println(dir)
		return
	}

	//If the code is here, it means that the arg is invalid
	fmt.Println("Error: invalid argument/s")
}
