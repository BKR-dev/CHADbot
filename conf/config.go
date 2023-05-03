package conf

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

var configName = "settings.conf"

// config struct
type Config struct {
	TopicId   string `json:"topicId"`
	ApiKey    string `json:"apiKey"`
	ApiUser   string `json:"apiUser"`
	ApiUserId string `json:"apiUserId"`
	Url       string `json:"url"`
}

// reads from file "settings.conf"
func readFile() ([]string, error) {
	file, err := os.Open(configName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
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

// remove spaces from string but return it
func removeSpaces(str string) string {
	// byte sloice size string
	b := make([]byte, 0, len(str))
	// go through string check if rune is space
	for _, c := range str {
		if !unicode.IsSpace(c) {
			// push to byte slice cast to byte
			b = append(b, byte(c))
		}
	}
	// return cast to string from byte sloice
	return string(b)
}

// Returns all entries in the config file
func GetSettings() (*Config, error) {
	fileLines, err := os.ReadFile("settings.conf")
	if err != nil {
		return nil, err
	}
	// make a map sized the number of lines (1 var / line)
	varMap := make(map[string]string, len(fileLines))
	reader := strings.NewReader(string(fileLines))
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		keyVal := strings.Split(scanner.Text(), "=")
		varMap[keyVal[0]] = keyVal[1]
	}
	// go thorugh the lines
	config := Config{
		TopicId:   varMap["topicId"],
		ApiKey:    varMap["apiKey"],
		ApiUser:   varMap["apiUser"],
		ApiUserId: varMap["apiUserId"],
		Url:       varMap["url"],
	}

	return &config, nil
}
