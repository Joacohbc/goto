package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const versionMessage string = "1.7" //Version

var (
	Help         bool
	Version      bool
	List         bool
	PathQuotes   bool
	ConfFilePath bool
	AddPath      string
	DelPath      string
	ModifyPath   string
)

func ArgIsShortOrNumber(arg string) (string, error) {

	//Load the config file in memory
	var directories []Directory
	if err := loadConfigFile(&directories); err != nil {
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
			if arg == dir.Abbreviation {
				return dir.Path, nil //In case of correct abbreviation, print and exit
			}
		}
	}

	return "", nil //In case of args is not a number or a valid abbreviation, continue
}

func init() {

	flag.BoolVar(&ConfFilePath, "path", false, "Print path of the config.json")

	flag.BoolVar(&Help, "h", false, "Print help message")
	flag.BoolVar(&Help, "help", false, "Print help message")

	flag.BoolVar(&Version, "v", false, "Print version")
	flag.BoolVar(&Version, "version", false, "Print version")

	flag.BoolVar(&List, "l", false, "Print all path with abbreviations")
	flag.BoolVar(&List, "list", false, "Print all path with abbreviations")

	flag.BoolVar(&PathQuotes, "q", false, "Print the path with quotes: -q")
	flag.BoolVar(&PathQuotes, "quotes", false, "Print the path with quotes: -quotes")

	flag.StringVar(&AddPath, "a", "", "Add a new path use: -a=[New Path],[New Short]")
	flag.StringVar(&AddPath, "add", "", "Add a new path use: -add=[New Path],[New Short]")

	flag.StringVar(&DelPath, "d", "", "Delete a path use: --d=[Path to Del]")
	flag.StringVar(&DelPath, "del", "", "Delete a path use: --del=[Path to Del]")

	flag.StringVar(&ModifyPath, "m", "", "Modify a path: -m=[Path],[New Short]")
	flag.StringVar(&ModifyPath, "modify", "", "Modify a path: -modif=[Path],[New Short]")

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

	if ConfFilePath {
		fmt.Println(ConfigFile)
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
	if len(AddPath) != 0 {

		dir, err := ToDirectory(AddPath)
		if err != nil {
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
	if len(DelPath) != 0 {

		if len(DelPath) == 0 {
			fmt.Println("Error: path can't be blank spaces")
			return
		}

		if err := delPaths(DelPath); err != nil {
			fmt.Println("Error:", err)
			fmt.Println("The changes were not applied")
			return
		}

		fmt.Println("The changes were applied successfully")
		return
	}

	//If the modify argument is passed, use func modify
	if len(ModifyPath) != 0 {

		dir, err := ToDirectory(ModifyPath)
		if err != nil {
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

	//Check if "arg" is an abbreviation or a number index
	path, err := ArgIsShortOrNumber(arg)
	if err != nil {
		fmt.Println("Error:", err)
		return

	} else if len(path) != 0 {
		fmt.Println(path)
		return
	}

	//If exists like afile
	if fileInfo, err := os.Stat(arg); err == nil {
		//If it's a directory
		if fileInfo.IsDir() {
			fmt.Println(arg)
			return

		} else {
			fmt.Println("Error: the path is not a directory")
			return
		}
	}

	//If the code is here, it means that the arg is invalid
	fmt.Println("Error: invalid argument/s")
}
