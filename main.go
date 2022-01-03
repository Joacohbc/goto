package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ArgIsDir(arg string) bool {

	fileInfo, err := os.Stat(arg)

	if err == nil {
		//If it's a directory
		if fileInfo.IsDir() {
			fmt.Println(arg)
			return true //If a file and is a directory, print and exit

		} else {
			fmt.Println("Error: the path is not a directory")
			return true //If the file is a directory, print error and exit
		}

	} else {
		return false //If not a file, continue
	}
}

func ArgIsShortOrNumber(arg string) bool {

	//Load the config file in memory
	var directories []directory
	err := loadConfigFile(&directories)

	if err != nil {
		fmt.Println("Error: ", err)
		return false //In case of error, print the error and exit
	}

	//Check if path is number
	if pathNumber, err := strconv.Atoi(arg); err == nil {

		for i, dir := range directories {

			if pathNumber == i {
				fmt.Println(dir.Path)
				return true //In case of correct pathNumber, print and exit
			}
		}

		fmt.Println("Error: the number is invalid(should be: 0-" + strconv.Itoa(len(directories)-1) + "), check config file")
		return true //In case of error, print the error and exit

	} else { //If it isn't a number
		for _, dir := range directories {

			if arg == dir.Short {
				fmt.Println(dir.Path)
				return true //In case of correct abbreviation, print and exit
			}
		}
	}

	return false //In case of args is not a number or a valid abbreviation, continue
}

const versionMessage string = "1.2" //Version

func helpMessage() string {
	helpMessage := `Three ways to use it, with abbreviations(config file), numbers(index of config file) and paths:

-Abbreviations= "goto <abbreviation>"
-Number="goto <number-of-the-index>"
-Path="goto <path>"

Path of config file: 
`
	return helpMessage + configPath()[1]
}

func main() {

	//Create the config file
	if err := createConfigFile(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	//Check if goto have argument
	help := flag.Bool("help", false, "Help message")

	version := flag.Bool("v", false, "Print version")

	list := flag.Bool("l", false, "Print all path with abbreviations")

	pathQuotes := flag.Bool("q", false, "Print the path with quotes")

	addPath := flag.String("add", "", "Add a new path use: --add=\"[New Path],[New Short]\"")

	delPath := flag.String("del", "", "Delete a path use: --del=\"[Path to Del]\"")

	modifyPath := flag.String("modify", "", "Modify a path: --modif=\"[Path],[New Short]\"")

	//Parse the flags
	flag.Parse()

	//If the help argument is passed, print help message
	if *help {
		fmt.Println(helpMessage())
		return
	}

	//If the version argument is passed, print version message
	if *version {
		fmt.Printf("Version of goto: %v", versionMessage)
		return
	}

	if *list {
		var directoriesToList []directory
		if err := loadConfigFile(&directoriesToList); err != nil {
			fmt.Println("Error:", err)
			return
		}

		for i, dir := range directoriesToList {
			fmt.Printf("%v- Path: \"%v\", Short: \"%v\" \n", i, dir.Path, dir.Short)
		}
		return
	}

	if *pathQuotes {
		fmt.Print("\"")
		ArgIsShortOrNumber(flag.Arg(0))
		fmt.Printf("\"")
		return
	}

	if len(*addPath) != 0 {

		args := strings.Split(*addPath, ",")

		if len(args) != 2 {
			fmt.Println("Error: bad format of --add")
			helpMessage()
			return
		}

		if len(args[0]) == 0 || len(args[1]) == 0 {
			fmt.Println("Error: path and abbreviation can't be blank spaces")
			return
		}

		dir := directory{Path: args[0], Short: args[1]}

		if err := addNewPaths(dir); err != nil {
			fmt.Println("Error:", err)
			fmt.Println("The changes were not applied")
			return
		}

		fmt.Println("The changes were applied successfully")
		return
	}

	if len(*delPath) != 0 {

		if len(*delPath) == 0 {
			fmt.Println("Error: path can't be blank spaces")
			return
		}

		if err := delPaths(*delPath); err != nil {
			fmt.Println("Error:", err)
			fmt.Println("The changes were not applied")
			return
		}

		fmt.Println("The changes were applied successfully")
		return
	}

	if len(*modifyPath) != 0 {

		args := strings.Split(*modifyPath, ",")

		if len(args) != 2 {
			fmt.Println("Error: bad format of --modify")
			helpMessage()
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

	//If exists like a Directory
	if ArgIsDir(arg) {
		return
	}

	//Check if "arg" is an abbreviation or a number index
	if ArgIsShortOrNumber(arg) {
		return
	}

	//If the code is here, it means that the arg is invalid
	fmt.Println("Error: Invalid argument/s")

}
