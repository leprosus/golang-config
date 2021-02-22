package config

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
	filePath, err = filepath.Abs(filePath)
	if err != nil {
		return
	}

	_, err = os.Stat(filePath)
	if err != nil {
		err = fmt.Errorf("can't load config file %s because: %s", filePath, err.Error())

		return
	}

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
	cfgPath := cfgFilePath.Load().(string)

	var info os.FileInfo
	info, err = os.Stat(cfgPath)
	if err != nil {
		err = fmt.Errorf("can't load config file %s because: %s", cfgPath, err.Error())

		return
	}

	if atomic.LoadUint32(&cfgIsLoaded) == 0 || atomic.LoadInt64(&cfgTimestamp) != info.ModTime().Unix() {
		var bs []byte
		bs, err = ioutil.ReadFile(cfgPath)
		if err != nil {
			err = fmt.Errorf("can't load config file %s because: %s", cfgPath, err.Error())

			return
		}

		var result Result
		result, err = ParseJson(bs)
		if err != nil {
			err = fmt.Errorf("file %s isn't valid json", cfgPath)

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

	var err error
	val, err = result.String(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns string value by json-path or default value
func StringOrDefault(path, defVal string) (val string) {
	if Exist(path) {
		return String(path)
	} else {
		return defVal
	}
}

// Returns bool value by json-path
func Bool(path string) (val bool) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	var err error
	val, err = result.Bool(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns bool value by json-path or default value
func BoolOrDefault(path string, defVal bool) (val bool) {
	if Exist(path) {
		return Bool(path)
	} else {
		return defVal
	}
}

// Returns int32 value by json-path
func Int32(path string) (val int32) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	var err error
	val, err = result.Int32(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns int32 value by json-path or default value
func Int32OrDefault(path string, defVal int32) (val int32) {
	if Exist(path) {
		return Int32(path)
	} else {
		return defVal
	}
}

// Returns uint32 value by json-path
func UInt32(path string) (val uint32) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	var err error
	val, err = result.UInt32(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns uint32 value by json-path or default value
func UInt32OrDefault(path string, defVal uint32) (val uint32) {
	if Exist(path) {
		return UInt32(path)
	} else {
		return defVal
	}
}

// Returns int64 value by json-path
func Int64(path string) (val int64) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	var err error
	val, err = result.Int64(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns int64 value by json-path or default value
func Int64OrDefault(path string, defVal int64) (val int64) {
	if Exist(path) {
		return Int64(path)
	} else {
		return defVal
	}
}

// Returns uint64 value by json-path
func UInt64(path string) (val uint64) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	var err error
	val, err = result.UInt64(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns uint64 value by json-path or default value
func UInt64OrDefault(path string, defVal uint64) (val uint64) {
	if Exist(path) {
		return UInt64(path)
	} else {
		return defVal
	}
}

// Returns float32 value by json-path
func Float32(path string) (val float32) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	var err error
	val, err = result.Float32(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns float32 value by json-path or default value
func Float32OrDefault(path string, defVal float32) (val float32) {
	if Exist(path) {
		return Float32(path)
	} else {
		return defVal
	}
}

// Returns float64 value by json-path
func Float64(path string) (val float64) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	var err error
	val, err = result.Float64(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns float64 value by json-path or default value
func Float64OrDefault(path string, defVal float64) (val float64) {
	if Exist(path) {
		return Float64(path)
	} else {
		return defVal
	}
}

// Returns array of strings value by json-path
func List(path string) (val []string) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	var err error
	val, err = result.List(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns array of strings value by json-path or default value
func ListOrDefault(path string, defVal []string) (val []string) {
	if Exist(path) {
		return List(path)
	} else {
		return defVal
	}
}

// Returns array of interfaces value by json-path
func Array(path string) (val []interface{}) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	var err error
	val, err = result.Array(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns array of interfaces value by json-path or default value
func ArrayOrDefault(path string, defVal []interface{}) (val []interface{}) {
	if Exist(path) {
		return Array(path)
	} else {
		return defVal
	}
}

// Returns json value by json-path
func JSON(path string) (val map[string]interface{}) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	var err error
	val, err = result.JSON(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns json value by json-path or default value
func JSONOrDefault(path string, defVal map[string]interface{}) (val map[string]interface{}) {
	if Exist(path) {
		return JSON(path)
	} else {
		return defVal
	}
}

// Returns duration value by json-path
func Duration(path string) (val time.Duration) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	var (
		i64 int64
		err error
	)
	i64, err = result.Int64(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	val = time.Duration(i64) * time.Second

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns duration value by json-path or default value
func DurationOrDefault(path string, defVal time.Duration) (val time.Duration) {
	if Exist(path) {
		return Duration(path)
	} else {
		return defVal
	}
}

// Returns path value by json-path
func Path(path string) (val string) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	var err error
	val, err = result.String(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	if len(val) > 0 && val[0:1] == "/" {
		return
	}

	if len(val) > 2 && val[1:3] == ":\\" {
		return
	}

	cfgPath := cfgFilePath.Load().(string)
	cfgPath = filepath.Dir(cfgPath)

	val = filepath.Join(cfgPath, val)

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns path value by json-path or default value
func PathOrDefault(path, defVal string) (val string) {
	if Exist(path) {
		return Path(path)
	} else {
		return defVal
	}
}

// Returns interface value by json-path
func Interface(path string) (val interface{}) {
	cfgLogger.Load().(logger).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfgJson.Load().(Result)

	var err error
	val, err = result.Interface(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	cfgLogger.Load().(logger).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Returns interface value by json-path or default value
func InterfaceOrDefault(path string, defVal interface{}) (val interface{}) {
	if Exist(path) {
		return Interface(path)
	} else {
		return defVal
	}
}

func handleErr(path string, err error) {
	switch err.(type) {
	case *ValueNotExist:
		cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` isn't exist", path))
	case *ValueUnexpectedType:
		cfgLogger.Load().(logger).warn(fmt.Sprintf("Value by path `%s` contains unexpected type of value", path))
	default:
		cfgLogger.Load().(logger).error(fmt.Sprintf("Parsing by path `%s` returns error: %v", path, err))
	}
}
