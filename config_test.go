package config

import (
	"encoding/json"
	"testing"
	"time"
)

var cfgBs = []byte(`{"str":"text", 
"int64": -1, "str_int64": "-1", "uint64": 1, "str_uint64": "1", 
"int32": -1, "str_int32": "-1", "uint32": 1, "str_uint32": "1", 
"float64": 0.1, "str_float64": 0.1,
"float32": 0.1, "str_float32": "0.1", 
"bool": true,
"slice": ["one", "two"],
"map": {"one": "val1", "two": "val2"},
"dur": 5}`)

func TestConfig(t *testing.T) {
	var v map[string]interface{}
	err := json.Unmarshal(cfgBs, &v)
	if err != nil {
		t.Error(err)
	}

	var res Result
	res, err = Parse(v)
	if err != nil {
		t.Error(err)
	}

	InitAsStruct(res)

	if String("str") != "text" {
		t.Error("Can't get value by path `str`")
	}

	if StringOrDefault("", "text") != "text" {
		t.Error("Get unexpected default value from StringOrDefault")
	}

	if Int64("int64") != -1 {
		t.Error("Can't get value by path `int64`")
	}

	if Int64OrDefault("", -1) != -1 {
		t.Error("Get unexpected default value from Int64OrDefault")
	}

	if Int64("str_int64") != -1 {
		t.Error("Can't get value by path `str_int64`")
	}

	if UInt64("uint64") != 1 {
		t.Error("Can't get value by path `int64`")
	}

	if UInt64OrDefault("", 1) != 1 {
		t.Error("Get unexpected default value from UInt64OrDefault")
	}

	if UInt64("str_uint64") != 1 {
		t.Error("Can't get value by path `str_int64`")
	}

	if Int32("int32") != -1 {
		t.Error("Can't get value by path `int32`")
	}

	if Int32OrDefault("", -1) != -1 {
		t.Error("Get unexpected default value from Int32OrDefault")
	}

	if Int32("str_int32") != -1 {
		t.Error("Can't get value by path `str_int32`")
	}

	if UInt32("uint32") != 1 {
		t.Error("Can't get value by path `int32`")
	}

	if UInt32OrDefault("", 1) != 1 {
		t.Error("Get unexpected default value from UInt32OrDefault")
	}

	if UInt32("str_uint32") != 1 {
		t.Error("Can't get value by path `str_int32`")
	}

	if Float64("float64") != 0.1 {
		t.Error("Can't get value by path `float64`")
	}

	if Float64OrDefault("", 0.1) != 0.1 {
		t.Error("Get unexpected default value from Float64OrDefault")
	}

	if Float64("str_float64") != 0.1 {
		t.Error("Can't get value by path `str_float64`")
	}

	if Float32("float32") != 0.1 {
		t.Error("Can't get value by path `float32`")
	}

	if Float32OrDefault("", 0.1) != 0.1 {
		t.Error("Get unexpected default value from Float32OrDefault")
	}

	if Float32("str_float32") != 0.1 {
		t.Error("Can't get value by path `str_float32`")
	}

	if Bool("bool") != true {
		t.Error("Can't get value by path `bool`")
	}

	if !BoolOrDefault("", true) {
		t.Error("Get unexpected default value from BoolOrDefault")
	}

	var sl = Slice("slice")
	if len(sl) != 2 || sl[0].(string) != "one" || sl[1].(string) != "two" {
		t.Error("Can't get value by path `slice`")
	}

	var defSl = SliceOrDefault("", []interface{}{"one", "two"})
	if len(defSl) != 2 || defSl[0].(string) != "one" || defSl[1].(string) != "two" {
		t.Error("Get unexpected default value from SliceOrDefault")
	}

	var ls = List("slice")
	if len(ls) != 2 || ls[0] != "one" || ls[1] != "two" {
		t.Error("Can't get value by path `slice`")
	}

	var defLs = ListOrDefault("", []string{"one", "two"})
	if len(defLs) != 2 || defLs[0] != "one" || defLs[1] != "two" {
		t.Error("Get unexpected default value from ListOrDefault")
	}

	var obj = Map("map")

	var (
		val interface{}
		ok  bool
	)
	if val, ok = obj["one"]; !ok || val.(string) != "val1" {
		t.Error("Can't get value by path `map`")
	}
	if val, ok = obj["two"]; !ok || val.(string) != "val2" {
		t.Error("Can't get value by path `map`")
	}

	var defObj = MapOrDefault("", map[string]interface{}{"one": "val1", "two": "val2"})
	if val, ok = defObj["one"]; !ok || val.(string) != "val1" {
		t.Error("Get unexpected default value from MapOrDefault")
	}
	if val, ok = defObj["two"]; !ok || val.(string) != "val2" {
		t.Error("Get unexpected default value from MapOrDefault")
	}

	var str string
	str, err = res.String("map.one")
	if err != nil {
		t.Error(err)
	}
	if str != "val1" {
		t.Error("Can't get value by path `map.one`")
	}

	var dur time.Duration
	dur, err = res.Duration("dur")
	if err != nil {
		t.Error(err)
	}
	if dur != 5*time.Second {
		t.Error("Can't get value by path `dur`")
	}
}
