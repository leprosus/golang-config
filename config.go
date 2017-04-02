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

// Returns flag is value existed by json-path
func Exist(path string) (bool) {
	_, ok := getResult(path)

	return ok
}

// Returns string value by json-path
func String(path string) (string) {
	result, _ := getResult(path)

	return result.String()
}

// Returns boolean value by json-path
func Bool(path string) (bool) {
	result, _ := getResult(path)

	return result.Bool()
}

// Returns integer value by json-path
func Int(path string) (int64) {
	result, _ := getResult(path)

	return result.Int()
}

// Returns array value by json-path
func Array(path string) ([]string) {
	slice := []string{}

	result, ok := getResult(path)

	if ok {
		for _, el := range result.Array() {
			slice = append(slice, el.String())
		}
	}

	return slice
}