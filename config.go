package golang_config

import (
	"fmt"
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
}

func init() {
	result, _ := ParseJson([]byte{})
	cfgJson.Store(result)

	cfgLogger.Store(logger{
		debug: func(message string) {},
		info:  func(message string) {},
		warn:  func(message string) {},
		error: func(message string) {},
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
			for range time.NewTicker(time.Second).C {
				err := refreshJson()

				if err != nil {
					cfgLogger.Load().(logger).error(err.Error())
				}
			}
		}()
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
		var bs []byte
		bs, err = ioutil.ReadFile(fileName)
		if err != nil {
			err = fmt.Errorf("can't load config file %s because: %s", fileName, err.Error())

			return
		}

		var result Result
		result, err = ParseJson(bs)
		if err != nil {
			err = fmt.Errorf("file %s isn't valid json", fileName)

			return
		}
		cfgJson.Store(result)

		atomic.StoreInt64(&cfgTimestamp, info.ModTime().Unix())

		if atomic.LoadUint32(&cfgIsLoaded) == 0 {
			cfgLogger.Load().(logger).info("Configuration is loaded")
		} else {
			cfgLogger.Load().(logger).info("Configuration is reloaded")
		}

		atomic.StoreUint32(&cfgIsLoaded, 1)

		mx := &sync.Mutex{}
		for _, callback := range cfgRefresh.Load().([]func()) {
			mx.Lock()
			callback()
			mx.Unlock()
		}
	}

	return
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

// Adds callback on refresh
func Refresh(callback func()) {
	r := cfgRefresh.Load().([]func())
	r = append(r, callback)
	cfgRefresh.Store(r)

	cfgLogger.Load().(logger).debug("Add callback on refresh")
}

// Returns flag is value existed by json-path
func Exist(path string) bool {
	return cfgJson.Load().(Result).IsExist(path)
}

// Returns string value by json-path
func String(path string) (val string) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	val, err := result.String(path)
	if err != nil {
		switch err.(type) {
		case *ValueNotExist:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` isn't exist", path))
		}

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns bool value by json-path
func Bool(path string) (val bool) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	val, err := result.Bool(path)
	if err != nil {
		switch err.(type) {
		case *ValueNotExist:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` isn't exist", path))
		case *ValueUnexpectedType:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` contains unexpected type of value", path))
		}

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns int32 value by json-path
func Int32(path string) (val int32) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	val, err := result.Int32(path)
	if err != nil {
		switch err.(type) {
		case *ValueNotExist:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` isn't exist", path))
		case *ValueUnexpectedType:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` contains unexpected type of value", path))
		}

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns uint32 value by json-path
func UInt32(path string) (val uint32) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	val, err := result.UInt32(path)
	if err != nil {
		switch err.(type) {
		case *ValueNotExist:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` isn't exist", path))
		case *ValueUnexpectedType:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` contains unexpected type of value", path))
		}

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns int64 value by json-path
func Int64(path string) (val int64) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	val, err := result.Int64(path)
	if err != nil {
		switch err.(type) {
		case *ValueNotExist:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` isn't exist", path))
		case *ValueUnexpectedType:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` contains unexpected type of value", path))
		}

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns uint64 value by json-path
func UInt64(path string) (val uint64) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	val, err := result.UInt64(path)
	if err != nil {
		switch err.(type) {
		case *ValueNotExist:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` isn't exist", path))
		case *ValueUnexpectedType:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` contains unexpected type of value", path))
		}

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns float32 value by json-path
func Float32(path string) (val float32) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	val, err := result.Float32(path)
	if err != nil {
		switch err.(type) {
		case *ValueNotExist:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` isn't exist", path))
		case *ValueUnexpectedType:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` contains unexpected type of value", path))
		}

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns float64 value by json-path
func Float64(path string) (val float64) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	val, err := result.Float64(path)
	if err != nil {
		switch err.(type) {
		case *ValueNotExist:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` isn't exist", path))
		case *ValueUnexpectedType:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` contains unexpected type of value", path))
		}

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns array value by json-path
func Array(path string) (val []string) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	val, err := result.Array(path)
	if err != nil {
		switch err.(type) {
		case *ValueNotExist:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` isn't exist", path))
		case *ValueUnexpectedType:
			cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` contains unexpected type of value", path))
		}

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}
