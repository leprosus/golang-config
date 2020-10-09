package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	Result map[string]interface{}

	ValueNotExist struct {
		message string
	}
	ValueUnexpectedType struct {
		message string
	}
)

func (e *ValueNotExist) Error() string {
	return e.message
}

func (e *ValueUnexpectedType) Error() string {
	return e.message
}

func ParseJson(bs []byte) (result Result, err error) {
	var v interface{}
	err = json.Unmarshal(bs, &v)
	if err != nil {
		return
	}

	switch v.(type) {
	case map[string]interface{}:
		result = v.(map[string]interface{})
	default:
		err = fmt.Errorf("get not json")
	}

	return
}

func (r Result) Interface(path string) (value interface{}, err error) {
	slices := strings.Split(path, ".")

	var (
		result = r

		v  interface{}
		ok bool

		last = len(slices) - 1
	)
	for i, slice := range slices {
		v, ok = result[slice]
		if !ok {
			err = &ValueNotExist{
				message: fmt.Sprintf("path `%s` isn't exist", path),
			}

			return
		}

		if i < last {
			result = v.(map[string]interface{})
		} else {
			value = v.(interface{})
		}
	}

	return
}

func (r Result) IsExist(path string) (ok bool) {
	_, err := r.Interface(path)

	ok = err == nil

	return
}

func (r Result) String(path string) (str string, err error) {
	var v interface{}
	v, err = r.Interface(path)
	if err != nil {
		return
	}

	str = fmt.Sprintf("%v", v)

	return
}

func (r Result) IsString(path string) (ok bool) {
	var (
		v   interface{}
		err error
	)
	v, err = r.Interface(path)
	if err != nil {
		return
	}

	_, ok = v.(string)

	return
}

func (r Result) Int32(path string) (i32 int32, err error) {
	var (
		v   interface{}
		f64 float64
		i64 uint64
		str string
		ok  bool
	)
	v, err = r.Interface(path)
	if err != nil {
		return
	}

	f64, ok = v.(float64)
	if ok {
		i32 = int32(f64)

		return
	}

	str, ok = v.(string)
	if !ok {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	i64, err = strconv.ParseUint(str, 10, 32)
	if err != nil {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	i32 = int32(i64)

	return
}

func (r Result) IsInt32(path string) (ok bool) {
	_, err := r.Int32(path)

	ok = err == nil

	return
}

func (r Result) UInt32(path string) (ui32 uint32, err error) {
	var (
		v    interface{}
		f64  float64
		ui64 uint64
		str  string
		ok   bool
	)
	v, err = r.Interface(path)
	if err != nil {
		return
	}

	f64, ok = v.(float64)
	if ok {
		ui32 = uint32(f64)

		return
	}

	str, ok = v.(string)
	if !ok {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	ui64, err = strconv.ParseUint(str, 10, 32)
	if err != nil {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	ui32 = uint32(ui64)

	return
}

func (r Result) IsUInt32(path string) (ok bool) {
	_, err := r.UInt32(path)

	ok = err == nil

	return
}

func (r Result) Int64(path string) (i64 int64, err error) {
	var (
		v   interface{}
		f64 float64
		str string
		ok  bool
	)
	v, err = r.Interface(path)
	if err != nil {
		return
	}

	f64, ok = v.(float64)
	if ok {
		i64 = int64(f64)

		return
	}

	str, ok = v.(string)
	if !ok {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	i64, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	return
}

func (r Result) IsInt64(path string) (ok bool) {
	_, err := r.Int64(path)

	ok = err == nil

	return
}

func (r Result) UInt64(path string) (ui64 uint64, err error) {
	var (
		v   interface{}
		f64 float64
		str string
		ok  bool
	)
	v, err = r.Interface(path)
	if err != nil {
		return
	}

	f64, ok = v.(float64)
	if ok {
		ui64 = uint64(f64)

		return
	}

	str, ok = v.(string)
	if !ok {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	ui64, err = strconv.ParseUint(str, 10, 64)
	if err != nil {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	return
}

func (r Result) IsUInt64(path string) (ok bool) {
	_, err := r.UInt64(path)

	ok = err == nil

	return
}

func (r Result) Float32(path string) (f32 float32, err error) {
	var (
		v   interface{}
		f64 float64
		str string
		ok  bool
	)
	v, err = r.Interface(path)
	if err != nil {
		return
	}

	f64, ok = v.(float64)
	if ok {
		f32 = float32(f64)

		return
	}

	str, ok = v.(string)
	if !ok {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	f64, err = strconv.ParseFloat(str, 32)
	if err != nil {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	f32 = float32(f64)

	return
}

func (r Result) IsFloat32(path string) (ok bool) {
	_, err := r.Float32(path)

	ok = err == nil

	return
}

func (r Result) Float64(path string) (f64 float64, err error) {
	var (
		v   interface{}
		str string
		ok  bool
	)
	v, err = r.Interface(path)
	if err != nil {
		return
	}

	f64, ok = v.(float64)
	if ok {
		return
	}

	str, ok = v.(string)
	if !ok {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	f64, err = strconv.ParseFloat(str, 64)
	if err != nil {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	return
}

func (r Result) IsFloat64(path string) (ok bool) {
	_, err := r.Float64(path)

	ok = err == nil

	return
}

func (r Result) Bool(path string) (flag bool, err error) {
	var (
		v   interface{}
		str string
		ok  bool
	)
	v, err = r.Interface(path)
	if err != nil {
		return
	}

	flag, ok = v.(bool)
	if ok {
		return
	}

	str, ok = v.(string)
	if !ok {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	flag, err = strconv.ParseBool(str)
	if err != nil {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	return
}

func (r Result) IsBool(path string) (ok bool) {
	_, err := r.Bool(path)

	ok = err == nil

	return
}

func (r Result) Array(path string) (array []string, err error) {
	var (
		v   interface{}
		str string
		ok  bool
	)
	v, err = r.Interface(path)
	if err != nil {
		return
	}

	slices, ok := v.([]interface{})
	if !ok {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	for _, slice := range slices {
		str, ok = slice.(string)
		if !ok {
			err = &ValueUnexpectedType{
				message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
			}

			return
		}

		array = append(array, str)
	}

	return
}

func (r Result) IsArray(path string) (ok bool) {
	_, err := r.Array(path)

	ok = err == nil

	return
}

func (r Result) JSON(path string) (object map[string]interface{}, err error) {
	object = map[string]interface{}{}

	var (
		v  interface{}
		ok bool
	)
	v, err = r.Interface(path)
	if err != nil {
		return
	}

	object, ok = v.(map[string]interface{})
	if !ok {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	return
}

func (r Result) IsJSON(path string) (ok bool) {
	_, err := r.JSON(path)

	ok = err == nil

	return
}

func (r Result) Duration(path string) (dur time.Duration, err error) {
	var i64 int64
	i64, err = r.Int64(path)
	if err != nil {
		return
	}

	dur = time.Duration(i64) * time.Second

	return
}

func (r Result) IsDuration(path string) (ok bool) {
	return r.IsInt64(path)
}
