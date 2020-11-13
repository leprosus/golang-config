# Golang thread-safe, lock-free JSON configuration reader

## Import
```go
import "github.com/leprosus/golang-config"
```

## Create new configuration reader

```go
cfg := config.Init("./config.json")

// To print debug information in stdout
cfg.Stdout(true)
```

## Json example

```json
{
    "string": "text",
    "digit": {
        "one": 1,
        "two": 2
    },
    "flag": true,
    "emails": [
        "user1@domain",
        "user2@domain"
    ]
}
```

NB: to get config values to need to use simple json-path requests [jsonpath.com](http://jsonpath.com)

## Getting example

```go
cfg := config.Init("./config.json", 60)

text := cfg.String("string")
one := cfg.Int("digit.one")
two := cfg.Int("digit.two")
flag := cfg.Bool("flag")
emails := cfg.List("emails")
```

## List all methods

### Initialization

* config.Init(path) - initializes configuration loading.
* config.Stdout(mode) - to set out all debug information into stdout.
* config.Debug(func(message string)) - sets custom logger for debug.
* config.Info(func(message string)) - sets custom logger for info.
* config.Warn(func(message string)) - sets custom logger for warn.
* config.Error(func(message string)) - sets custom logger for error.
* config.Refresh(func()) - adds callback on refresh.

### Getting data

* cgf.String("json.path") - returns string value by json path.
* cgf.StringOrDefault("json.path", "default") - returns string value by json path or default value.
* cgf.Int32("json.path") - returns int32 value by json path.
* cgf.Int32OrDefault("json.path", 0) - returns int32 value by json path or default value.
* cgf.UInt32("json.path") - returns uint32 value by json path.
* cgf.UInt32OrDefault("json.path", 0) - returns uint32 value by json path or default value.
* cgf.Int64("json.path") - returns int64 value by json path.
* cgf.Int64OrDefault("json.path", 0) - returns int64 value by json path or default value.
* cgf.UInt64("json.path") - returns uint64 value by json path.
* cgf.UInt64OrDefault("json.path", 0) - returns uint64 value by json path or default value.
* cgf.Float32("json.path") - returns float32 value by json path.
* cgf.Float32OrDefault("json.path", 0) - returns float32 value by json path or default value.
* cgf.Float64("json.path") - returns float64 value by json path.
* cgf.Float64OrDefault("json.path", 0) - returns float64 value by json path or default value.
* cgf.Bool("json.path") - returns bool value by json path.
* cgf.BoolOrDefault("json.path", false) - returns bool value by json path or default value.
* cgf.List("json.path") - returns strings array by json path.
* cgf.ListOrDefault("json.path", []string{""}) - returns strings array by json path or default value.
* cgf.Array("json.path") - returns array of interfaces by json path.
* cgf.ArrayOrDefault("json.path") - returns array of interfaces by json path or default value.
* cgf.JSON("json.path") - returns map[string]interface by json path.
* cgf.JSONOrDefault("json.path", []interface{}{""}) - returns map[string]interface by json path or default value.
* cgf.Duration("json.path") - returns duration in seconds by json path.
* cgf.DurationOrDefault("json.path", time.Second) - returns duration in seconds by json path or default value.
* cgf.Interface("json.path") - returns interface by json path.
* cgf.InterfaceOrDefault("json.path", "") - returns interface by json path or default value.