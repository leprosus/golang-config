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
	refreshPeriod int64
	lastRefresh   int64
	stdout        bool
	logger        func(message string)
}

var cfg config = config{}

// Init config: set file path to config file and period config refresh
func Init(filePath string) {
	if len(cfg.filePath) == 0 {
		cfg = config{
			filePath: filePath,
			lastRefresh: time.Now().Unix()}

		cfg.stdout = false

		cfg.logger = func(message string) {
			if cfg.stdout {
				fmt.Println(message)
			}
		}

		cfg.loadJson()
	}
}

func (config *config) loadJson() string {
	timeDiff := time.Now().Unix() - config.lastRefresh

	if len(config.json) == 0 || (config.refreshPeriod > 0 && timeDiff > config.refreshPeriod) {
		fileName, _ := filepath.Abs(config.filePath)

		json, err := ioutil.ReadFile(fileName)
		if err != nil {
			message := fmt.Sprintf("Can't load config file by %s\n", fileName)
			cfg.logger(message)

			os.Exit(1)
		}

		config.lastRefresh = time.Now().Unix()

		config.json = string(json)
	}

	return config.json
}

func getResult(path string) (gjson.Result, bool) {
	cfg.logger(fmt.Sprintf("Try to get value by %s\n", path))

	result := gjson.Get(cfg.loadJson(), path)

	if !result.Exists() {
		message := fmt.Sprintf("Can't get value by `%s`", path)
		cfg.logger(message)

		return gjson.Result{}, true
	}

	return result, false
}

func Stdout(mode bool) {
	cfg.stdout = mode
}

func Logger(callback func(message string)) {
	cfg.logger = func(message string) {
		callback(message)

		if cfg.stdout {
			fmt.Println(message)
		}
	}
}

func RefreshPeriod(refreshPeriod int64) {
	cfg.refreshPeriod = refreshPeriod
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