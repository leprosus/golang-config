package config

import (
    "path/filepath"
    "io/ioutil"
    "github.com/tidwall/gjson"
    "time"
    "fmt"
)

type config struct {
    filePath      string
    json          string
    refreshPeriod int64
    lastRefresh   int64
    logger        logger
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
        fatal: func(message string) {}}}

// Init config: set file path to config file and period config refresh
func Init(filePath string) {
    cfg.logger.info("Configuration is intialized")

    if len(cfg.filePath) == 0 {
        cfg.filePath = filePath
        cfg.lastRefresh = time.Now().Unix()

        cfg.loadJson()
    }
}

func (config *config) loadJson() string {
    timeDiff := time.Now().Unix() - config.lastRefresh

    if len(config.json) == 0 || (config.refreshPeriod > 0 && timeDiff > config.refreshPeriod) {
        if len(config.json) == 0 {
            cfg.logger.info("Configuration is loaded")
        } else {
            cfg.logger.info("Configuration is reloaded")
        }

        fileName, _ := filepath.Abs(config.filePath)

        json, err := ioutil.ReadFile(fileName)
        if err != nil {
            cfg.logger.fatal(fmt.Sprintf("Can't load config file by %s", fileName))
        }

        config.lastRefresh = time.Now().Unix()

        config.json = string(json)
    }

    return config.json
}

func getResult(path string) (gjson.Result, bool) {
    cfg.logger.debug(fmt.Sprintf("Try to get value by %s", path))

    result := gjson.Get(cfg.loadJson(), path)

    if !result.Exists() {
        cfg.logger.warn(fmt.Sprintf("Value by path `%s` isn't exist", path))

        return gjson.Result{}, true
    }

    cfg.logger.debug(fmt.Sprintf("Value by path `%s` is exist and is set `%s`", path, result.String()))

    return result, false
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

// Refreshes configuration reloading
func RefreshPeriod(refreshPeriod int64) {
    cfg.refreshPeriod = refreshPeriod
    cfg.logger.debug(fmt.Sprintf("Set new refresh period %d s", refreshPeriod))
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
