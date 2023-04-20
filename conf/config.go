package conf

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

var configName = "settings.conf"

func readFile() ([]string, error) {
	content, err := os.ReadFile(configName)
	if err != nil {
		return nil, err
	}
	// cut content into bye slices at newline
	bSlices := bytes.Split(content, []byte("\n"))
	// way too dank gamer move to check if last slice is empty and removing it
	if len(bSlices) > 0 && len(bSlices[len(bSlices)-1]) == 0 {
		bSlices = bSlices[:len(bSlices)-1]
	}
	// init a string slice of len byteSlices - blazingly fast!
	stringSlice := make([]string, 0, len(bSlices))
	// loop through bSlices and add to stringslice
	for i, v := range bSlices {
		str := string(v)
		if len(str) <= 0 {
			return nil, fmt.Errorf("Invalid line %d on config file, aborting.\n", i+1)
		}
		if string(str[0]) != "#" {
			stringSlice = append(stringSlice, str)
		}
	}
	return stringSlice, nil
}

func GetSettings() (map[string]string, error) {
	fileLines, err := readFile()
	if err != nil {
		return nil, err
	}
	// make a map sized the number of lines (1 var / line)
	varMap := make(map[string]string, len(fileLines))

	for _, v := range fileLines {
		fmt.Println(v)
		keyVal := strings.Split(v, "=")
		varMap[keyVal[0]] = keyVal[1]
	}

	return varMap, nil
}
