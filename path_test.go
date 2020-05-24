package config

import (
	"testing"
)

var bs = []byte(`{"str":"text", 
"int64": -1, "str_int64": "-1", "uint64": 1, "str_uint64": "1", 
"int32": -1, "str_int32": "-1", "uint32": 1, "str_uint32": "1", 
"float32": 0.1, "str_float32": "0.1", 
"float64": 0.1, "str_float64": 0.1,
"bool": true,
"list": ["one", "two"],
"obj": {"one": "val1", "two": "val2"}}`)

func TestParser(t *testing.T) {
	res, err := ParseJson(bs)
	if err != nil {
		t.Error(err)
	}

	_, err = res.String("str")
	if err != nil {
		t.Error(err)
	}

	_, err = res.Int64("int64")
	if err != nil {
		t.Error(err)
	}

	_, err = res.Int64("str_int64")
	if err != nil {
		t.Error(err)
	}

	_, err = res.UInt64("uint64")
	if err != nil {
		t.Error(err)
	}

	_, err = res.UInt64("str_uint64")
	if err != nil {
		t.Error(err)
	}

	_, err = res.Int64("int32")
	if err != nil {
		t.Error(err)
	}

	_, err = res.Int64("str_int32")
	if err != nil {
		t.Error(err)
	}

	_, err = res.UInt64("uint32")
	if err != nil {
		t.Error(err)
	}

	_, err = res.UInt64("str_uint32")
	if err != nil {
		t.Error(err)
	}

	_, err = res.Float32("float32")
	if err != nil {
		t.Error(err)
	}

	_, err = res.Float32("str_float32")
	if err != nil {
		t.Error(err)
	}

	_, err = res.Float64("float64")
	if err != nil {
		t.Error(err)
	}

	_, err = res.Float64("str_float64")
	if err != nil {
		t.Error(err)
	}

	_, err = res.Bool("bool")
	if err != nil {
		t.Error(err)
	}

	_, err = res.Array("list")
	if err != nil {
		t.Error(err)
	}

	_, err = res.JSON("obj")
	if err != nil {
		t.Error(err)
	}

	_, err = res.String("obj.one")
	if err != nil {
		t.Error(err)
	}
}
