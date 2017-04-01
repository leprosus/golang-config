package config

import (
	"path/filepath"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"os"
	"time"
	"fmt"
	"github.com/pkg/errors"
)

type config struct {
	filePath      string
	json          string
	refreshPeriod int64
	lastRefresh   int64
	debug         bool
}

var cfg config = config{}

// Init config: set file path to config file and period config refresh
func Init(filePath string, refreshPeriod int64) {
	if len(cfg.filePath) == 0 {
		cfg = config{
			filePath: filePath,
			refreshPeriod: refreshPeriod,
			lastRefresh: time.Now().Unix()}

		cfg.debug = false

		cfg.loadJson()
	}
}

func (config *config) loadJson() string {
	timeDiff := time.Now().Unix() - config.lastRefresh

	if len(config.json) == 0 || timeDiff > config.refreshPeriod {
		fileName, _ := filepath.Abs(config.filePath)

		json, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Printf("Can't load config file by %s\n", fileName)

			os.Exit(1)
		}

		config.lastRefresh = time.Now().Unix()

		config.json = string(json)
	}

	return config.json
}

func getResult(path string) (gjson.Result, bool) {
	if cfg.debug {
		fmt.Printf("Try to get value by %s\n", path)
	}

	result := gjson.Get(cfg.loadJson(), path)

	if !result.Exists() {
		if cfg.debug {
			err := errors.New(fmt.Sprintf("Can't get value by `%s`", path))

			fmt.Printf("Catch error %s\n", err.Error())
		}

		return gjson.Result{}, true
	}

	return result, false
}

// Return string value by json-path
func String(path string) (string, bool) {
	result, err := getResult(path)

	return result.String(), err == nil
}

// Return boolean value by json-path
func Bool(path string) (bool, bool) {
	result, err := getResult(path)

	return result.Bool(), err == nil
}

// Return integer value by json-path
func Int(path string) (int64, bool) {
	result, err := getResult(path)

	return result.Int(), err == nil
}

// Return array value by json-path
func Array(path string) ([]string, bool) {
	slice := []string{}

	result, err := getResult(path)

	for _, el := range result.Array() {
		slice = append(slice, el.String())
	}

	return slice, err == nil
}