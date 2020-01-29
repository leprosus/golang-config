package config

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

var (
	once         = &sync.Once{}
	cfgFilePath  = &atomic.Value{}
	cfgJson      = &atomic.Value{}
	cfgIsLoaded  uint32
	cfgTimestamp int64
	cfgLogger    = &atomic.Value{}
	cfgRefresh   = &atomic.Value{}
)

type logger struct {
	debug func(message string)
	info  func(message string)
	warn  func(message string)
	error func(message string)
	fatal func(message string)
}

func init() {
	cfgJson.Store(gjson.ParseBytes([]byte{}))

	cfgLogger.Store(logger{
		debug: func(message string) {},
		info:  func(message string) {},
		warn:  func(message string) {},
		error: func(message string) {},
		fatal: func(message string) {},
	})

	cfgRefresh.Store([]func(){})
}

// Init config: set file path to config file and period config refresh
func Init(filePath string) (err error) {
	cfgFilePath.Store(filePath)

	cfgLogger.Load().(logger).info("Configuration is initialized")

	once.Do(func() {
		err = refreshJson()
		if err != nil {
			return
		}

		go func() {
			ticker := time.NewTicker(time.Second)

			for range ticker.C {
				err := refreshJson()

				if err != nil {
					cfgLogger.Load().(logger).error(err.Error())
				}
			}
		}()

		time.Sleep(time.Second)
	})

	return
}

func refreshJson() (err error) {
	fileName, err := filepath.Abs(cfgFilePath.Load().(string))
	if err != nil {
		return
	}

	info, err := os.Stat(fileName)
	if err != nil {
		err = fmt.Errorf("can't load config file %s because: %s", fileName, err.Error())

		return
	}

	if atomic.LoadUint32(&cfgIsLoaded) == 0 || atomic.LoadInt64(&cfgTimestamp) != info.ModTime().Unix() {
		var jsonStr []byte
		jsonStr, err = ioutil.ReadFile(fileName)
		if err != nil {
			err = fmt.Errorf("can't load config file %s because: %s", fileName, err.Error())

			return
		}

		if !isJson(jsonStr) {
			err = fmt.Errorf("file %s isn't valid json", fileName)

			return
		}

		atomic.StoreInt64(&cfgTimestamp, info.ModTime().Unix())

		if atomic.LoadUint32(&cfgIsLoaded) == 0 {
			cfgLogger.Load().(logger).info("Configuration is loaded")
		} else {
			cfgLogger.Load().(logger).info("Configuration is reloaded")
		}

		atomic.StoreUint32(&cfgIsLoaded, 1)
		cfgJson.Store(gjson.ParseBytes(jsonStr))

		mx := &sync.Mutex{}
		for _, callback := range cfgRefresh.Load().([]func()) {
			mx.Lock()
			callback()
			mx.Unlock()
		}
	}

	return
}

func isJson(jsonStr []byte) bool {
	var data map[string]interface{}

	return json.Unmarshal(jsonStr, &data) == nil
}

func getResult(path string) (*gjson.Result, bool) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(gjson.Result).Get(path)
	if !result.Exists() {
		cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` isn't exist", path))

		return &gjson.Result{}, false
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%s`", path, result.String()))

	return &result, true
}

// Sets logger for debug
func Debug(callback func(message string)) {
	l := cfgLogger.Load().(logger)
	l.debug = callback
	cfgLogger.Store(l)

	cfgLogger.Load().(logger).debug("Set custom debug logger")
}

// Sets logger for into
func Info(callback func(message string)) {
	l := cfgLogger.Load().(logger)
	l.info = callback
	cfgLogger.Store(l)

	cfgLogger.Load().(logger).debug("Set custom info logger")
}

// Sets logger for warning
func Warn(callback func(message string)) {
	l := cfgLogger.Load().(logger)
	l.warn = callback
	cfgLogger.Store(l)

	cfgLogger.Load().(logger).debug("Set custom warning logger")
}

// Sets logger for error
func Error(callback func(message string)) {
	l := cfgLogger.Load().(logger)
	l.error = callback
	cfgLogger.Store(l)

	cfgLogger.Load().(logger).debug("Set custom error logger")
}

// Sets logger for fatal
func Fatal(callback func(message string)) {
	l := cfgLogger.Load().(logger)
	l.fatal = callback
	cfgLogger.Store(l)

	cfgLogger.Load().(logger).debug("Set custom fatal logger")
}

// Adds callback on refresh
func Refresh(callback func()) {
	r := cfgRefresh.Load().([]func())
	r = append(r, callback)
	cfgRefresh.Store(r)

	cfgLogger.Load().(logger).debug("Add callback on refresh")
}

// Returns flag is value existed by json-path
func Exist(path string) bool {
	return cfgJson.Load().(gjson.Result).Exists()
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
