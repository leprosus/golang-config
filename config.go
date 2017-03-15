package config

import (
	"path/filepath"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"os"
	"time"
	"fmt"
)

type config struct {
	filePath      string
	json          string
	refreshPeriod int
	lastRefresh   int
}

var cfg config = config{}

// Init config: set file path to config file and period config refresh
func Init(filePath string, refreshPeriod int) {
	if len(cfg.filePath) == 0 {
		cfg = config{
			filePath: filePath,
			refreshPeriod: refreshPeriod,
			lastRefresh: time.Now(),
		}

		cfg.loadJson()
	}
}

func (config *config) loadJson() string {
	timeDiff := time.Now() - config.lastRefresh

	if len(config.json) == 0 || timeDiff > config.refreshPeriod {
		fileName, _ := filepath.Abs(config.filePath)

		json, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Printf("Can't load config file by %s\n", fileName)

			os.Exit(1)
		}

		config.json = string(json)
	}

	return config.json
}

func getResult(path string) gjson.Result {
	return gjson.Get(cfg.loadJson(), path)
}

// Return string value by json-path
func String(path string) string {
	return getResult(path).String()
}

// Return boolean value by json-path
func Bool(path string) bool {
	return getResult(path).Bool()
}

// Return integer value by json-path
func Int(path string) int64 {
	return getResult(path).Int()
}

// Return array value by json-path
func Array(path string) (slice []string) {
	result := getResult(path)

	for _, el := range result.Array() {
		slice = append(slice, el.String())
	}

	return slice
}