# Golang Json configuration reader

## Create new configuration reader

``` go
cfg := config.Init("./config.json", 60)
```

## Json example

``` json
{
    "string": "text",
    "digit: {
        "one: 1,
        "two: 2
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

``` go
cfg := config.Init("./config.json", 60)

text := cfg.String("string")
one := cfg.Int("digit.one")
two := cfg.Int("digit.two")
flag := cfg.Bool("flag")
emails := cfg.Array("emails")
```

## List all methods

### Initialization

* config.Init(path, refreshPeriod) - initializes configuration loading

### Getting data

* cgf.String("json.path") - returns string value by json path
* cgf.Int("json.path") - returns int64 value by json path
* cgf.Bool("json.path") - returns boolean value by json path
* cgf.Array("json.path") - returns strings array by json path