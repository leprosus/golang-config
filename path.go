package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	Object map[string]interface{}

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

func Parse(val interface{}) (obj Object, err error) {
	switch val.(type) {
	case map[string]interface{}:
		obj = val.(map[string]interface{})
	default:
		err = fmt.Errorf("can't use the structure")
	}

	return
}

func (o Object) Interface(path string) (value interface{}, err error) {
	slices := strings.Split(path, ".")

	var (
		obj = o

		v  interface{}
		ok bool

		last = len(slices) - 1
	)
	for i, slice := range slices {
		v, ok = obj[slice]
		if !ok {
			err = &ValueNotExist{
				message: fmt.Sprintf("path `%s` isn't exist", path),
			}

			return
		}

		if i < last {
			obj = v.(map[string]interface{})
		} else {
			value = v.(interface{})
		}
	}

	return
}

func (o Object) IsExist(path string) (ok bool) {
	_, err := o.Interface(path)

	ok = err == nil

	return
}

func (o Object) String(path string) (str string, err error) {
	var v interface{}
	v, err = o.Interface(path)
	if err != nil {
		return
	}

	str = fmt.Sprintf("%v", v)

	return
}

func (o Object) IsString(path string) (ok bool) {
	var (
		v   interface{}
		err error
	)
	v, err = o.Interface(path)
	if err != nil {
		return
	}

	_, ok = v.(string)

	return
}

func (o Object) Int32(path string) (i32 int32, err error) {
	var (
		v   interface{}
		f64 float64
		i64 int64
		str string
		ok  bool
	)
	v, err = o.Interface(path)
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

	i64, err = strconv.ParseInt(str, 10, 32)
	if err != nil {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	i32 = int32(i64)

	return
}

func (o Object) IsInt32(path string) (ok bool) {
	_, err := o.Int32(path)

	ok = err == nil

	return
}

func (o Object) UInt32(path string) (ui32 uint32, err error) {
	var (
		v    interface{}
		f64  float64
		ui64 uint64
		str  string
		ok   bool
	)
	v, err = o.Interface(path)
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

func (o Object) IsUInt32(path string) (ok bool) {
	_, err := o.UInt32(path)

	ok = err == nil

	return
}

func (o Object) Int64(path string) (i64 int64, err error) {
	var (
		v   interface{}
		f64 float64
		str string
		ok  bool
	)
	v, err = o.Interface(path)
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

func (o Object) IsInt64(path string) (ok bool) {
	_, err := o.Int64(path)

	ok = err == nil

	return
}

func (o Object) UInt64(path string) (ui64 uint64, err error) {
	var (
		v   interface{}
		f64 float64
		str string
		ok  bool
	)
	v, err = o.Interface(path)
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

func (o Object) IsUInt64(path string) (ok bool) {
	_, err := o.UInt64(path)

	ok = err == nil

	return
}

func (o Object) Float32(path string) (f32 float32, err error) {
	var (
		v   interface{}
		f64 float64
		str string
		ok  bool
	)
	v, err = o.Interface(path)
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

func (o Object) IsFloat32(path string) (ok bool) {
	_, err := o.Float32(path)

	ok = err == nil

	return
}

func (o Object) Float64(path string) (f64 float64, err error) {
	var (
		v   interface{}
		str string
		ok  bool
	)
	v, err = o.Interface(path)
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

func (o Object) IsFloat64(path string) (ok bool) {
	_, err := o.Float64(path)

	ok = err == nil

	return
}

func (o Object) Bool(path string) (flag bool, err error) {
	var (
		v   interface{}
		str string
		ok  bool
	)
	v, err = o.Interface(path)
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

func (o Object) IsBool(path string) (ok bool) {
	_, err := o.Bool(path)

	ok = err == nil

	return
}

func (o Object) List(path string) (list []string, err error) {
	var slices []interface{}
	slices, err = o.Slice(path)
	if err != nil {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	var (
		str string
		ok  bool
	)
	for _, slice := range slices {
		str, ok = slice.(string)
		if !ok {
			err = &ValueUnexpectedType{
				message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
			}

			return
		}

		list = append(list, str)
	}

	return
}

func (o Object) IsList(path string) (ok bool) {
	_, err := o.List(path)

	ok = err == nil

	return
}

func (o Object) Slice(path string) (array []interface{}, err error) {
	var (
		v  interface{}
		ok bool
	)
	v, err = o.Interface(path)
	if err != nil {
		return
	}

	array, ok = v.([]interface{})
	if !ok {
		err = &ValueUnexpectedType{
			message: fmt.Sprintf("path `%s` contains unexpected type of value", path),
		}

		return
	}

	return
}

func (o Object) IsSlice(path string) (ok bool) {
	_, err := o.Slice(path)

	ok = err == nil

	return
}

func (o Object) Map(path string) (object map[string]interface{}, err error) {
	object = map[string]interface{}{}

	var (
		v  interface{}
		ok bool
	)
	v, err = o.Interface(path)
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

func (o Object) IsMap(path string) (ok bool) {
	_, err := o.Map(path)

	ok = err == nil

	return
}

func (o Object) Duration(path string) (dur time.Duration, err error) {
	var i64 int64
	i64, err = o.Int64(path)
	if err != nil {
		return
	}

	dur = time.Duration(i64) * time.Second

	return
}

func (o Object) IsDuration(path string) (ok bool) {
	return o.IsInt64(path)
}
