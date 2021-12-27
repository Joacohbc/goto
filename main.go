package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func ArgIsDir(arg string) bool {

	fileInfo, err := os.Stat(arg)

	if err == nil {
		//If it's a directory
		if fileInfo.IsDir() {
			fmt.Println(arg)
			return true //If a file and is a directory, print and exit

		} else {
			fmt.Println("Error: The path is not a directory")
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
		return true //In case of error, print th error and exit
	}

	//Check if path is number
	if pathNumber, err := strconv.Atoi(arg); err == nil {

		for i, dir := range directories {

			if pathNumber == i {
				fmt.Println(dir.Path)
				return true //In case of correct pathNumber, print and exit
			}
		}

		fmt.Println("Error: The number is invalid(should be: 0-" + strconv.Itoa(len(directories)-1) + "), check config file")
		return true //In case of error, print the error and exit

	} else { //If it isn't a number
		for _, dir := range directories {

			if arg == dir.Short {
				fmt.Println(dir.Path)
				return true //In case of correct abbreviation, print and exit
			}
		}
	}
	return false //In caes of args is not a number or a valid abbreviation, continue
}

const versionMessage string = "1.0" //Version

func helpMessage() string {
	helpMessage := `Goto is a command to move between folders, it has 3 way to use it, with abbreviations(config file), numbers(index of config file) and paths:

-Abbreviations= "goto <abbreviation>"
-Number="goto <number-of-the-index>"
-Path="goto <path>"

Path of config file: 
`
	return helpMessage + configPath()[1]
}

func main() {

	//Create the config file
	createConfigFile()

	//Check if goto have argument
	path := flag.String("path", "", "Path to go")

	help := flag.Bool("help", false, "Help message")

	version := flag.Bool("version", false, "Print version")

	//Parse the flags
	flag.Parse()

	//If the help argument is passed, print help message
	if *help {
		fmt.Print(helpMessage())
		return
	}

	//If the version argument is passed, print version message
	if *version {
		fmt.Printf("Version of goto: %v", versionMessage)
		return
	}

	//Where the first argument will be stored
	var arg string = *path

	//If exists like a Directory
	if ArgIsDir(arg) {
		return
	}

	//Check if "arg" is an abbreviation or a number index
	if ArgIsShortOrNumber(arg) {
		return
	}

	//If the code is here, it means that the arg is invalid
	fmt.Print("Error: Invalid argument/s")
}
