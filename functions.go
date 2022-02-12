package main

import "fmt"

func addNewPaths(newDir Directory) error {

	var directories []Directory
	loadConfigFile(&directories)

	for _, dir := range directories {
		if newDir.Path == dir.Path {
			return fmt.Errorf("the path: \"%v\" already exists", newDir.Path)
		}
	}

	directories = append(directories, newDir)

	if err := validConfiguredPaths(directories); err != nil {
		return fmt.Errorf("invalid new config file,\n%v", err)
	}

	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}

func delPaths(pathToDel string) error {

	var directories []Directory
	loadConfigFile(&directories)

	var find bool = false
	for i, dir := range directories {
		if dir.Path == pathToDel {
			directories = append(directories[:i], directories[i+1:]...)
			find = true
			break
		}
	}

	if !find {
		return fmt.Errorf("path \"%v\" doesn't exist", pathToDel)
	}

	if err := validConfiguredPaths(directories); err != nil {
		return fmt.Errorf("invalid new config file,\n%v", err)
	}

	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}

func modPaths(pathToModif string, newShort string) error {

	var directories []Directory
	err := loadConfigFile(&directories)
	if err != nil {
		return err
	}

	var exist bool = false
	for i, dir := range directories {

		if dir.Path == pathToModif {
			directories[i].Short = newShort
			exist = true
		}
	}

	if !exist {
		return fmt.Errorf("the path that you are trying to modify is not exists")
	}

	if err := validConfiguredPaths(directories); err != nil {
		return fmt.Errorf("invalid new config file,\n%v", err)
	}

	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}
