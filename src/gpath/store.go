package gpath

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bytedance/sonic"
)

// Create default GotoPath entries
func createDefaultGotoPathsFile(gotoPathsFile string) ([]GotoPath, error) {
	var gpaths []GotoPath

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Default paths
	gpaths = append(gpaths, GotoPath{
		Path:         home,
		Abbreviation: "h",
	})

	gpaths = append(gpaths, GotoPath{
		Path:         filepath.Dir(gotoPathsFile),
		Abbreviation: "config",
	})

	return gpaths, nil
}

// Create the config file if not already exists
func CreateGotoPathsFile(gotoPathsFile string) error {
	if _, err := os.Stat(gotoPathsFile); err == nil {
		return nil
	}

	gpaths, err := createDefaultGotoPathsFile(gotoPathsFile)
	if err != nil {
		return err
	}

	return SaveGPathsFile(gpaths, gotoPathsFile)
}

// Validate the array (using CheckRepeatedItems) and create a paths file from directory array
func SaveGPathsFile(gpaths []GotoPath, gotoPathsFile string) error {

	if err := CheckRepeatedItems(gpaths); err != nil {
		return err
	}

	// Open the file for writing (create or truncate)
	file, err := os.OpenFile(gotoPathsFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	// Use Buffered Writer for efficiency as suggested in the blog
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Encode directly to the stream
	enc := sonic.ConfigDefault.NewEncoder(writer)
	enc.SetIndent("", "\t")
	if err := enc.Encode(gpaths); err != nil {
		return err
	}

	return nil
}

// Load config file into an array
func LoadGPathsFile(gpaths *[]GotoPath, gotoPathsFile string) error {

	// Open the File
	file, err := os.Open(gotoPathsFile)
	if err != nil {
		return fmt.Errorf("error reading config file")
	}
	defer file.Close()

	// Use Buffered Reader for efficiency
	reader := bufio.NewReader(file)

	// Load the Paths using sonic's stream decoder
	dec := sonic.ConfigFastest.NewDecoder(reader)
	if err := dec.Decode(gpaths); err != nil {
		return fmt.Errorf("error parsing config file")
	}

	//If all is OK, check dir and return
	return CheckRepeatedItems(*gpaths)
}
