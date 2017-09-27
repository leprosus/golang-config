package config

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
	"encoding/json"
	"sync"
	"time"
	"crypto/md5"
	"io"
	"github.com/pkg/errors"
	"sync/atomic"
)

type config struct {
	sync.Once

	filePath string
	json     atomic.Value
	md5      string
	logger   logger
	refresh  []func()
}

type logger struct {
	debug func(message string)
	info  func(message string)
	warn  func(message string)
	error func(message string)
	fatal func(message string)
}

var cfg config = config{
	logger: logger{
		debug: func(message string) {},
		info:  func(message string) {},
		warn:  func(message string) {},
		error: func(message string) {},
		fatal: func(message string) {}},
	refresh: []func(){}}

// Init config: set file path to config file and period config refresh
func Init(filePath string) {
	cfg.logger.info("Configuration is initialized")

	cfg.filePath, _ = filepath.Abs(filePath)
	refreshJson()

	cfg.Do(func() {
		go watchFile()
	})
}

func refreshJson() {
	_, err := os.Stat(cfg.filePath)
	if err != nil {
		cfg.logger.fatal(fmt.Sprintf("Can't load config file by %s: %s", cfg.filePath, err.Error()))

		return
	}

	fileMd5, err := getFileMd5(cfg.filePath)
	if err != nil {
		cfg.logger.error(err.Error())

		return
	}

	if cfg.md5 != fileMd5 {
		jsonStr, err := ioutil.ReadFile(cfg.filePath)
		if err != nil {
			cfg.logger.fatal(fmt.Sprintf("Can't load config file by %s: %s", cfg.filePath, err.Error()))

			return
		}

		if !isJson(jsonStr) {
			cfg.logger.fatal(fmt.Sprintf("File %s isn't valid json", cfg.filePath))

			return
		}

		if cfg.json.Load() == nil {
			cfg.logger.info("Configuration is loaded")
		} else {
			cfg.logger.info("Configuration is reloaded")
		}

		cfg.json.Store(string(jsonStr))

		cfg.md5 = fileMd5

		for _, callback := range cfg.refresh {
			callback()
		}
	}
}

func getFileMd5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", errors.Errorf("Can't read file %s because %s", filePath, err.Error())
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", errors.Errorf("Can't copy file %s to hash function because %s", filePath, err.Error())
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))

	return hash, nil
}

func watchFile() {
	for {
		refreshJson()

		time.Sleep(30 * time.Second)
	}
}

func isJson(jsonStr []byte) bool {
	var data map[string]interface{}
	return json.Unmarshal(jsonStr, &data) == nil

}

func getResult(path string) (gjson.Result, bool) {
	cfg.logger.debug(fmt.Sprintf("Try to get value by %s", path))

	jsonStr := cfg.json.Load().(string)
	result := gjson.Get(jsonStr, path)

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
	jsonStr := cfg.json.Load().(string)
	return gjson.Get(jsonStr, path).Exists()
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

// Returns array value by json-path
func Array(path string) []string {
	slice := []string{}

	result, ok := getResult(path)

	if ok {
		for _, el := range result.Array() {
			slice = append(slice, el.String())
		}
	}

	return slice
}
