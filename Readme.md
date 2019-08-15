# Golang thread-safe, lock-free JSON configuration reader

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

NB: to get config values to need to use json-path requests [jsonpath.com](http://jsonpath.com)

## Getting example

```go
cfg := config.Init("./config.json", 60)

text := cfg.String("string")
one := cfg.Int("digit.one")
two := cfg.Int("digit.two")
flag := cfg.Bool("flag")
emails := cfg.Array("emails")
```

## List all methods

### Initialization

* config.Init(path) - initializes configuration loading
* config.Stdout(mode) - to set out all debug information into stdout
* config.Debug(func(message string)) - sets custom logger for debug
* config.Info(func(message string)) - sets custom logger for info
* config.Warn(func(message string)) - sets custom logger for warn
* config.Error(func(message string)) - sets custom logger for error
* config.Fatal(func(message string)) - sets custom logger for fatal
* config.Refresh(func()) - adds callback on refresh

### Getting data

* cgf.String("json.path") - returns string value by json path
* cgf.Int("json.path") - returns int64 value by json path
* cgf.Bool("json.path") - returns boolean value by json path
* cgf.Array("json.path") - returns strings array by json path