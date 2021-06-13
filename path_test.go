package config

import (
	"encoding/json"
	"testing"
	"time"
)

var jsonBs = []byte(`{"str":"text", 
"int64": -1, "str_int64": "-1", "uint64": 1, "str_uint64": "1", 
"int32": -1, "str_int32": "-1", "uint32": 1, "str_uint32": "1", 
"float64": 0.1, "str_float64": 0.1,
"float32": 0.1, "str_float32": "0.1", 
"bool": true,
"slice": ["one", "two"],
"map": {"one": "val1", "two": "val2"},
"dur": 5}`)

func TestJsonParser(t *testing.T) {
	var val map[string]interface{}
	err := json.Unmarshal(jsonBs, &val)
	if err != nil {
		t.Error(err)
	}

	var res Result
	res, err = Parse(val)
	if err != nil {
		t.Error(err)
	}

	testParser(res, t)
}

func testParser(res Result, t *testing.T) {
	var err error

	var str string
	str, err = res.String("str")
	if err != nil {
		t.Error(err)
	}
	if str != "text" {
		t.Error("Method String returns unexpected result")
	}

	if !res.IsString("str") {
		t.Error("Method IsString returns unexpected result")
	}

	var i64 int64
	i64, err = res.Int64("int64")
	if err != nil {
		t.Error(err)
	}
	if i64 != -1 {
		t.Error("Method Int64 returns unexpected result")
	}

	if !res.IsInt64("int64") {
		t.Error("Method IsInt64 returns unexpected result")
	}

	i64, err = res.Int64("str_int64")
	if err != nil {
		t.Error(err)
	}
	if i64 != -1 {
		t.Error("Method Int64 returns unexpected result")
	}

	if !res.IsInt64("str_int64") {
		t.Error("Method IsInt64 returns unexpected result")
	}

	var ui64 uint64
	ui64, err = res.UInt64("uint64")
	if err != nil {
		t.Error(err)
	}
	if ui64 != 1 {
		t.Error("Method UInt64 returns unexpected result")
	}

	if !res.IsUInt64("uint64") {
		t.Error("Method IsUInt64 returns unexpected result")
	}

	ui64, err = res.UInt64("str_uint64")
	if err != nil {
		t.Error(err)
	}
	if ui64 != 1 {
		t.Error("Method UInt64 returns unexpected result")
	}

	if !res.IsUInt64("str_uint64") {
		t.Error("Method IsUInt64 returns unexpected result")
	}

	var i32 int32
	i32, err = res.Int32("int32")
	if err != nil {
		t.Error(err)
	}
	if i32 != -1 {
		t.Error("Method Int32 returns unexpected result")
	}

	if !res.IsInt32("int32") {
		t.Error("Method IsInt32 returns unexpected result")
	}

	i32, err = res.Int32("str_int32")
	if err != nil {
		t.Error(err)
	}
	if i32 != -1 {
		t.Error("Method Int32 returns unexpected result")
	}

	if !res.IsInt32("str_int32") {
		t.Error("Method IsInt32 returns unexpected result")
	}

	var ui32 uint32
	ui32, err = res.UInt32("uint32")
	if err != nil {
		t.Error(err)
	}
	if ui32 != 1 {
		t.Error("Method UInt32 returns unexpected result")
	}

	if !res.IsUInt32("uint32") {
		t.Error("Method IsUInt32 returns unexpected result")
	}

	ui32, err = res.UInt32("str_uint32")
	if err != nil {
		t.Error(err)
	}
	if ui32 != 1 {
		t.Error("Method UInt32 returns unexpected result")
	}

	if !res.IsUInt32("str_uint32") {
		t.Error("Method IsUInt32 returns unexpected result")
	}

	var f64 float64
	f64, err = res.Float64("float64")
	if err != nil {
		t.Error(err)
	}
	if f64 != 0.1 {
		t.Error("Method Float64 returns unexpected result")
	}

	if !res.IsFloat64("float64") {
		t.Error("Method IsFloat64 returns unexpected result")
	}

	f64, err = res.Float64("str_float64")
	if err != nil {
		t.Error(err)
	}
	if f64 != 0.1 {
		t.Error("Method Float64 returns unexpected result")
	}

	if !res.IsFloat64("str_float64") {
		t.Error("Method IsFloat64 returns unexpected result")
	}

	var f32 float32
	f32, err = res.Float32("float32")
	if err != nil {
		t.Error(err)
	}
	if f32 != 0.1 {
		t.Error("Method Float32 returns unexpected result")
	}

	if !res.IsFloat32("float32") {
		t.Error("Method IsFloat32 returns unexpected result")
	}

	f32, err = res.Float32("str_float32")
	if err != nil {
		t.Error(err)
	}
	if f32 != 0.1 {
		t.Error("Method Float32 returns unexpected result")
	}

	if !res.IsFloat32("str_float32") {
		t.Error("Method IsFloat32 returns unexpected result")
	}

	var b bool
	b, err = res.Bool("bool")
	if err != nil {
		t.Error(err)
	}
	if !b {
		t.Error("Method Bool returns unexpected result")
	}

	if !res.IsBool("bool") {
		t.Error("Method IsBool returns unexpected result")
	}

	var sl []interface{}
	sl, err = res.Slice("slice")
	if err != nil {
		t.Error(err)
	}
	if len(sl) != 2 {
		t.Error("Method Slice returns unexpected result")
	}

	if !res.IsSlice("slice") {
		t.Error("Method IsSlice returns unexpected result")
	}

	var ls []string
	ls, err = res.List("slice")
	if err != nil {
		t.Error(err)
	}
	if len(ls) != 2 {
		t.Error("Method List returns unexpected result")
	}

	if !res.IsList("slice") {
		t.Error("Method IsList returns unexpected result")
	}

	var obj map[string]interface{}
	obj, err = res.Map("map")
	if err != nil {
		t.Error(err)
	}
	if obj["one"] != "val1" || obj["two"] != "val2" {
		t.Error("Method Map returns unexpected result")
	}

	if !res.IsMap("map") {
		t.Error("Method IsMap returns unexpected result")
	}

	_, err = res.String("map.one")
	if err != nil {
		t.Error(err)
	}

	var dur time.Duration
	dur, err = res.Duration("dur")
	if err != nil {
		t.Error(err)
	}
	if dur != 5*time.Second {
		t.Error("Method Duration returns unexpected result")
	}

	if !res.IsDuration("dur") {
		t.Error("Method IsDuration returns unexpected result")
	}
}
