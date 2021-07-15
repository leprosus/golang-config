package config

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

var checkBs = []byte(`{
	"str": {
		"required": true,
		"type": "string",
		"regexp": "^text$",
		"handler": "strHandler"
	},
	"int64": {
		"required": true,
		"type": "int64",
		"regexp": "^-1$",
		"handler": "int64Handler"
	},
	"uint64": {
		"required": true,
		"type": "uint64",
		"regexp": "^1$",
		"handler": "uint64Handler"
	},
	"int32": {
		"required": true,
		"type": "int32",
		"regexp": "^-1$",
		"handler": "int32Handler"
	},
	"uint32": {
		"required": true,
		"type": "uint32",
		"regexp": "^1$",
		"handler": "uint32Handler"
	},
	"float64": {
		"required": true,
		"type": "float64",
		"regexp": "^0\\.1$",
		"handler": "float64Handler"
	},
	"float32": {
		"required": true,
		"type": "float64",
		"regexp": "^0\\.1$",
		"handler": "float64Handler"
	},
	"bool": {
		"required": true,
		"type": "bool",
		"regexp": "^true$",
		"handler": "boolHandler"
	},
	"slice": {
		"required": true,
		"type": "slice",
		"handler": "sliceHandler"
	},
	"list": {
		"required": true,
		"type": "list",
		"handler": "listHandler"
	},
	"map": {
		"required": true,
		"type": "map",
		"handler": "mapHandler"
	},
	"dur": {
		"required": true,
		"type": "duration",
		"regexp": "^5$",
		"handler": "durationHandler"
	}
}`)

func TestChecker(t *testing.T) {
	var (
		handlers = map[string]Handler{
			"strHandler":      strHandler,
			"int64Handler":    int64Handler,
			"uint64Handler":   uint64Handler,
			"int32Handler":    int32Handler,
			"uint32Handler":   uint32Handler,
			"float64Handler":  float64Handler,
			"float32Handler":  float32Handler,
			"boolHandler":     boolHandler,
			"sliceHandler":    sliceHandler,
			"listHandler":     listHandler,
			"mapHandler":      mapHandler,
			"durationHandler": durationHandler,
		}

		checker *Checker
		err     error
	)

	checker, err = NewChecker(checkBs, handlers)
	if err != nil {
		t.Error(err)
	}

	var val map[string]interface{}
	err = json.Unmarshal(jsonBs, &val)
	if err != nil {
		t.Error(err)
	}

	var obj Object
	obj, err = Parse(val)
	if err != nil {
		t.Error(err)
	}

	err = checker.Check(obj)
	if err != nil {
		t.Error(err)
	}
}

func strHandler(v interface{}) (err error) {
	var (
		str string
		ok  bool
	)

	str, ok = v.(string)
	if !ok {
		err = fmt.Errorf("can't convert to string the value `%v`", v)

		return
	}

	if str != "text" {
		err = fmt.Errorf("the value `%v` isn't an expected value", v)

		return
	}

	return
}

func int64Handler(v interface{}) (err error) {
	var (
		f64 float64
		ok  bool
	)

	f64, ok = v.(float64)
	if !ok {
		err = fmt.Errorf("can't convert to f64 the value `%v`", v)

		return
	}

	if int64(f64) != -1 {
		err = fmt.Errorf("the value `%v` isn't an expected value", v)

		return
	}

	return
}

func uint64Handler(v interface{}) (err error) {
	var (
		f64 float64
		ok  bool
	)

	f64, ok = v.(float64)
	if !ok {
		err = fmt.Errorf("can't convert to f64 the value `%v`", v)

		return
	}

	if uint64(f64) != 1 {
		err = fmt.Errorf("the value `%v` isn't an expected value", v)

		return
	}

	return
}

func int32Handler(v interface{}) (err error) {
	var (
		f64 float64
		ok  bool
	)

	f64, ok = v.(float64)
	if !ok {
		err = fmt.Errorf("can't convert to f64 the value `%v`", v)

		return
	}

	if int32(f64) != -1 {
		err = fmt.Errorf("the value `%v` isn't an expected value", v)

		return
	}

	return
}

func uint32Handler(v interface{}) (err error) {
	var (
		f64 float64
		ok  bool
	)

	f64, ok = v.(float64)
	if !ok {
		err = fmt.Errorf("can't convert to f64 the value `%v`", v)

		return
	}

	if uint32(f64) != 1 {
		err = fmt.Errorf("the value `%v` isn't an expected value", v)

		return
	}

	return
}

func float64Handler(v interface{}) (err error) {
	var (
		f64 float64
		ok  bool
	)

	f64, ok = v.(float64)
	if !ok {
		err = fmt.Errorf("can't convert to f64 the value `%v`", v)

		return
	}

	if f64 != 0.1 {
		err = fmt.Errorf("the value `%v` isn't an expected value", v)

		return
	}

	return
}

func float32Handler(v interface{}) (err error) {
	var (
		f32 float32
		ok  bool
	)

	f32, ok = v.(float32)
	if !ok {
		err = fmt.Errorf("can't convert to f64 the value `%v`", v)

		return
	}

	if f32 != 0.1 {
		err = fmt.Errorf("the value `%v` isn't an expected value", v)

		return
	}

	return
}

func boolHandler(v interface{}) (err error) {
	var (
		flag bool
		ok   bool
	)

	flag, ok = v.(bool)
	if !ok {
		err = fmt.Errorf("can't convert to f64 the value `%v`", v)

		return
	}

	if !flag {
		err = fmt.Errorf("the value `%v` isn't an expected value", v)

		return
	}

	return
}

func sliceHandler(v interface{}) (err error) {
	var (
		sl []interface{}
		ok bool
	)

	sl, ok = v.([]interface{})
	if !ok {
		err = fmt.Errorf("can't convert to f64 the value `%v`", v)

		return
	}

	if len(sl) != 2 {
		err = fmt.Errorf("the value `%v` isn't an expected value", v)

		return
	}

	return
}

func listHandler(v interface{}) (err error) {
	var (
		sl []interface{}
		ls []string
		ok bool
	)

	sl, ok = v.([]interface{})
	if !ok {
		err = fmt.Errorf("can't convert to f64 the value `%v`", v)

		return
	}

	for _, val := range sl {
		ls = append(ls, val.(string))
	}

	if len(ls) != 2 {
		err = fmt.Errorf("the value `%v` isn't an expected value", v)

		return
	}

	return
}

func mapHandler(v interface{}) (err error) {
	var (
		hp map[string]interface{}
		ok bool
	)

	hp, ok = v.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("can't convert to f64 the value `%v`", v)

		return
	}

	if len(hp) != 2 {
		err = fmt.Errorf("the value `%v` isn't an expected value", v)

		return
	}

	return
}

func durationHandler(v interface{}) (err error) {
	var (
		f64 float64
		ok  bool
	)

	f64, ok = v.(float64)
	if !ok {
		err = fmt.Errorf("can't convert to f64 the value `%v`", v)

		return
	}

	if time.Duration(f64) != 5*time.Nanosecond {
		err = fmt.Errorf("the value `%v` isn't an expected value", v)

		return
	}

	return
}
