package config

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type config struct {
	filePath      string
	json          string
	fileTimestamp int64
	logger        logger
	refresh       []func()
	mx            sync.Mutex
}

type logger struct {
	debug func(message string)
	info  func(message string)
	warn  func(message string)
	error func(message string)
	fatal func(message string)
}

var (
	cfg config = config{
		logger: logger{
			debug: func(message string) { fmt.Fprintln(os.Stdout, message) },
			info:  func(message string) { fmt.Fprintln(os.Stdout, message) },
			warn:  func(message string) { fmt.Fprintln(os.Stdout, message) },
			error: func(message string) { fmt.Fprintln(os.Stderr, message) },
			fatal: func(message string) { fmt.Fprintln(os.Stderr, message) }},
		refresh: []func(){}}
	once = sync.Once{}
)

// Init config: set file path to config file and period config refresh
func Init(filePath string) (err error) {
	cfg.logger.info("Configuration is initialized")

	cfg.filePath = filePath

	once.Do(func() {
		err = refreshJson()

		go func() {
			ticker := time.NewTicker(time.Second)

			for range ticker.C {
				err = refreshJson()

				cfg.logger.error(err.Error())
			}
		}()
	})

	return
}

func getJson() string {
	return cfg.json
}

func refreshJson() (err error) {
	fileName, _ := filepath.Abs(cfg.filePath)

	info, err := os.Stat(cfg.filePath)
	if err != nil {
		err = fmt.Errorf("can't load config file %s because: %s", fileName, err.Error())

		return
	}

	if len(cfg.json) == 0 || cfg.fileTimestamp != info.ModTime().Unix() {
		jsonStr, err := ioutil.ReadFile(fileName)
		if err != nil {
			err = fmt.Errorf("can't load config file %s because: %s", fileName, err.Error())

			return
		}

		if !isJson(jsonStr) {
			err = fmt.Errorf("file %s isn't valid json", fileName)

			return
		}

		cfg.fileTimestamp = info.ModTime().Unix()

		if len(cfg.json) == 0 {
			cfg.logger.info("Configuration is loaded")
		} else {
			cfg.logger.info("Configuration is reloaded")
		}

		cfg.json = string(jsonStr)

		for _, callback := range cfg.refresh {
			go callback()
		}
	}

	return
}

func isJson(jsonStr []byte) bool {
	var data map[string]interface{}
	return json.Unmarshal(jsonStr, &data) == nil

}

func getResult(path string) (gjson.Result, bool) {
	cfg.logger.debug(fmt.Sprintf("Try to get value by %s", path))

	result := gjson.Get(getJson(), path)

	if !result.Exists() {
		cfg.logger.warn(fmt.Sprintf("Value by path `%s` isn't exist", path))

		return gjson.Result{}, false
	}

	cfg.logger.debug(fmt.Sprintf("Value by path `%s` is exist and is set `%s`", path, result.String()))

	return result, true
}

// Sets logger for debig
func Debug(callback func(message string)) {
	cfg.logger.debug = callback
	cfg.logger.debug("Set custom debug logger")
}

// Sets logger for into
func Info(callback func(message string)) {
	cfg.logger.info = callback
	cfg.logger.debug("Set custom info logger")
}

// Sets logger for warning
func Warn(callback func(message string)) {
	cfg.logger.warn = callback
	cfg.logger.debug("Set custom warning logger")
}

// Sets logger for error
func Error(callback func(message string)) {
	cfg.logger.error = callback
	cfg.logger.debug("Set custom error logger")
}

// Sets logger for fatal
func Fatal(callback func(message string)) {
	cfg.logger.fatal = callback
	cfg.logger.debug("Set custom fatal logger")
}

// Adds callback on refresh
func Refresh(callback func()) {
	cfg.refresh = append(cfg.refresh, callback)
	cfg.logger.debug("Add callback on refresh")
}

// Returns flag is value existed by json-path
func Exist(path string) bool {
	return gjson.Get(getJson(), path).Exists()
}

// Returns string value by json-path
func String(path string) string {
	result, _ := getResult(path)

	return result.String()
}

// Returns boolean value by json-path
func Bool(path string) bool {
	result, _ := getResult(path)

	return result.Bool()
}

// Returns integer value by json-path
func Int(path string) int64 {
	result, _ := getResult(path)

	return result.Int()
}

// Returns float value by json-path
func Float(path string) float64 {
	result, _ := getResult(path)

	return result.Float()
}

// Returns array value by json-path
func Array(path string) (slice []string) {
	result, ok := getResult(path)

	if ok {
		for _, el := range result.Array() {
			slice = append(slice, el.String())
		}
	}

	return
}
