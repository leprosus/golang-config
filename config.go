package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

var (
	once = &sync.Once{}

	filePath  = &atomic.Value{}
	format    = &atomic.Value{}
	isLoaded  uint32
	timestamp int64
	cfg       = &atomic.Value{}

	logger = &atomic.Value{}

	withRefresh uint32
	refreshers  = &atomic.Value{}
)

type logFn struct {
	debug func(message string)
	info  func(message string)
	warn  func(message string)
	error func(message string)
}

func init() {
	result, _ := Parse(map[string]interface{}{})
	cfg.Store(result)

	logger.Store(logFn{
		debug: func(message string) {},
		info:  func(message string) {},
		warn:  func(message string) {},
		error: func(message string) {},
	})

	refreshers.Store([]func(){})
}

// Init sets file path to a file with configuration (JSON format) and set periodically refresh data from its
func Init(cfgPath string) (err error) {
	atomic.StoreUint32(&withRefresh, 1)

	cfgPath, err = filepath.Abs(cfgPath)
	if err != nil {
		return
	}

	_, err = os.Stat(cfgPath)
	if err != nil {
		err = fmt.Errorf("can't load config file %s because: %s", cfgPath, err.Error())

		return
	}

	filePath.Store(cfgPath)

	logger.Load().(logFn).info("Configuration is initialized")

	err = refreshJson()
	if err != nil {
		return
	}

	once.Do(func() {
		go func() {
			var e error

			for range time.NewTicker(time.Second).C {
				e = refreshJson()

				if e != nil {
					logger.Load().(logFn).error(e.Error())
				}
			}
		}()
	})

	return
}

// InitAsStruct sets interface (Result struct) as configuration
func InitAsStruct(result Result) {
	atomic.StoreUint32(&withRefresh, 0)

	cfg.Store(result)
}

func refreshJson() (err error) {
	if atomic.LoadUint32(&withRefresh) != 1 {
		return
	}

	cfgPath := filePath.Load().(string)

	var info os.FileInfo
	info, err = os.Stat(cfgPath)
	if err != nil {
		err = fmt.Errorf("can't load config file %s because: %s", cfgPath, err.Error())

		return
	}

	if atomic.LoadUint32(&isLoaded) == 0 || atomic.LoadInt64(&timestamp) != info.ModTime().Unix() {
		var bs []byte
		bs, err = ioutil.ReadFile(cfgPath)
		if err != nil {
			err = fmt.Errorf("can't load config file %s because: %s", cfgPath, err.Error())

			return
		}

		var result Result
		err = json.Unmarshal(bs, &result)
		if err != nil {
			err = fmt.Errorf("file %s isn't suported configuration", cfgPath)

			return
		}
		cfg.Store(result)

		atomic.StoreInt64(&timestamp, info.ModTime().Unix())

		if atomic.LoadUint32(&isLoaded) == 0 {
			logger.Load().(logFn).info("Configuration is loaded")
		} else {
			logger.Load().(logFn).info("Configuration is reloaded")
		}

		atomic.StoreUint32(&isLoaded, 1)

		mx := &sync.Mutex{}
		for _, callback := range refreshers.Load().([]func()) {
			mx.Lock()
			callback()
			mx.Unlock()
		}
	}

	return
}

// Debug sets logger for debug
func Debug(callback func(message string)) {
	l := logger.Load().(logFn)
	l.debug = callback
	logger.Store(l)

	logger.Load().(logFn).debug("Set custom debug logger")
}

// Info sets logger for into
func Info(callback func(message string)) {
	l := logger.Load().(logFn)
	l.info = callback
	logger.Store(l)

	logger.Load().(logFn).debug("Set custom info logger")
}

// Warn sets logger for warning
func Warn(callback func(message string)) {
	l := logger.Load().(logFn)
	l.warn = callback
	logger.Store(l)

	logger.Load().(logFn).debug("Set custom warning logger")
}

// Error sets logger for error
func Error(callback func(message string)) {
	l := logger.Load().(logFn)
	l.error = callback
	logger.Store(l)

	logger.Load().(logFn).debug("Set custom error logger")
}

// Refresh adds callback on refresh
func Refresh(callback func()) {
	r := refreshers.Load().([]func())
	r = append(r, callback)
	refreshers.Store(r)

	logger.Load().(logFn).debug("Add callback on refresh")
}

// Exist returns flag is value existed by path
func Exist(path string) bool {
	return cfg.Load().(Result).IsExist(path)
}

// String returns string value by path
func String(path string) (val string) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

	var err error
	val, err = result.String(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// StringOrDefault returns string value by path or default value
func StringOrDefault(path, defVal string) (val string) {
	if Exist(path) {
		return String(path)
	} else {
		return defVal
	}
}

// Bool returns bool value by path
func Bool(path string) (val bool) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

	var err error
	val, err = result.Bool(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// BoolOrDefault returns bool value by path or default value
func BoolOrDefault(path string, defVal bool) (val bool) {
	if Exist(path) {
		return Bool(path)
	} else {
		return defVal
	}
}

// Int32 returns int32 value by path
func Int32(path string) (val int32) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

	var err error
	val, err = result.Int32(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Int32OrDefault returns int32 value by path or default value
func Int32OrDefault(path string, defVal int32) (val int32) {
	if Exist(path) {
		return Int32(path)
	} else {
		return defVal
	}
}

// UInt32 returns uint32 value by path
func UInt32(path string) (val uint32) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

	var err error
	val, err = result.UInt32(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// UInt32OrDefault returns uint32 value by path or default value
func UInt32OrDefault(path string, defVal uint32) (val uint32) {
	if Exist(path) {
		return UInt32(path)
	} else {
		return defVal
	}
}

// Int64 returns int64 value by path
func Int64(path string) (val int64) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

	var err error
	val, err = result.Int64(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Int64OrDefault returns int64 value by path or default value
func Int64OrDefault(path string, defVal int64) (val int64) {
	if Exist(path) {
		return Int64(path)
	} else {
		return defVal
	}
}

// UInt64 returns uint64 value by path
func UInt64(path string) (val uint64) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

	var err error
	val, err = result.UInt64(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// UInt64OrDefault returns uint64 value by path or default value
func UInt64OrDefault(path string, defVal uint64) (val uint64) {
	if Exist(path) {
		return UInt64(path)
	} else {
		return defVal
	}
}

// Float32 returns float32 value by path
func Float32(path string) (val float32) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

	var err error
	val, err = result.Float32(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Float32OrDefault returns float32 value by path or default value
func Float32OrDefault(path string, defVal float32) (val float32) {
	if Exist(path) {
		return Float32(path)
	} else {
		return defVal
	}
}

// Float64 returns float64 value by path
func Float64(path string) (val float64) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

	var err error
	val, err = result.Float64(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// Float64OrDefault returns float64 value by path or default value
func Float64OrDefault(path string, defVal float64) (val float64) {
	if Exist(path) {
		return Float64(path)
	} else {
		return defVal
	}
}

// List returns slice of strings value by path
func List(path string) (val []string) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

	var err error
	val, err = result.List(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// ListOrDefault returns slice of strings value by path or default value
func ListOrDefault(path string, defVal []string) (val []string) {
	if Exist(path) {
		return List(path)
	} else {
		return defVal
	}
}

// Slice returns slice of interfaces value by path
func Slice(path string) (val []interface{}) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

	var err error
	val, err = result.Slice(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// SliceOrDefault returns array of interfaces value by path or default value
func SliceOrDefault(path string, defVal []interface{}) (val []interface{}) {
	if Exist(path) {
		return Slice(path)
	} else {
		return defVal
	}
}

// Map returns map value by path
func Map(path string) (val map[string]interface{}) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

	var err error
	val, err = result.Map(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// MapOrDefault returns map by path or default value
func MapOrDefault(path string, defVal map[string]interface{}) (val map[string]interface{}) {
	if Exist(path) {
		return Map(path)
	} else {
		return defVal
	}
}

// Duration returns duration value by path
func Duration(path string) (val time.Duration) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

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

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// DurationOrDefault returns duration value by path or default value
func DurationOrDefault(path string, defVal time.Duration) (val time.Duration) {
	if Exist(path) {
		return Duration(path)
	} else {
		return defVal
	}
}

// Path returns path value by path
func Path(path string) (val string) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

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

	cfgPath := filePath.Load().(string)
	cfgPath = filepath.Dir(cfgPath)

	val = filepath.Join(cfgPath, val)

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// PathOrDefault returns path value by path or default value
func PathOrDefault(path, defVal string) (val string) {
	if Exist(path) {
		return Path(path)
	} else {
		return defVal
	}
}

// Interface returns interface value by path
func Interface(path string) (val interface{}) {
	logger.Load().(logFn).debug(fmt.Sprintf("Try to get value by %s", path))

	result := cfg.Load().(Result)

	var err error
	val, err = result.Interface(path)
	if err != nil {
		handleErr(path, err)

		return
	}

	logger.Load().(logFn).debug(fmt.Sprintf("Value by path `%s` is exist and is set `%v`", path, val))

	return
}

// InterfaceOrDefault returns interface value by path or default value
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
		logger.Load().(logFn).warn(fmt.Sprintf("Value by path `%s` isn't exist", path))
	case *ValueUnexpectedType:
		logger.Load().(logFn).warn(fmt.Sprintf("Value by path `%s` contains unexpected type of value", path))
	default:
		logger.Load().(logFn).error(fmt.Sprintf("Parsing by path `%s` returns error: %v", path, err))
	}
}
