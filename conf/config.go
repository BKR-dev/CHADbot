package conf

import (
	"fmt"
	"os"
)

var configName = "settings.conf"

func readFile() ([]string, error) {
	content, err := os.ReadFile(configName)
	if err != nil {
		return nil, err
	}
	fmt.Println(content)
	return nil, nil
}

func GetSettings() {
	readFile()
}
