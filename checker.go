package config

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Checker struct {
	HandlerBinds map[string]Handler
	RuleSet      RuleSet
}

type Handler func(v interface{}) (err error)

type Rule struct {
	IsRequired bool
	Type       Type
	RegExp     *regexp.Regexp
	Handler    Handler
}

type RuleSet map[string]Rule

type Type string

type UnexpectedScheme struct {
	message string
}

func (e *UnexpectedScheme) Error() string {
	return e.message
}

type UnexpectedRule struct {
	message string
}

func (e *UnexpectedRule) Error() string {
	return e.message
}

type UnexpectedType struct {
	message string
}

func (e *UnexpectedType) Error() string {
	return e.message
}

type UnexpectedHandler struct {
	message string
}

func (e *UnexpectedHandler) Error() string {
	return e.message
}

type UnexpectedValue struct {
	message string
}

func (e *UnexpectedValue) Error() string {
	return e.message
}

const (
	StringType   Type = "string"
	BoolType     Type = "bool"
	Int32Type    Type = "int32"
	UInt32Type   Type = "uint32"
	Int64Type    Type = "int64"
	UInt64Type   Type = "uint64"
	Float32Type  Type = "float32"
	Float64Type  Type = "float64"
	ListType     Type = "list"
	SliceType    Type = "slice"
	MapType      Type = "map"
	DurationType Type = "duration"
)

func NewChecker(jsonRules []byte, handlerBinds map[string]Handler) (c *Checker, err error) {
	c = &Checker{
		HandlerBinds: handlerBinds,
	}

	var v interface{}
	err = json.Unmarshal(jsonRules, &v)
	if err != nil {
		err = &UnexpectedScheme{
			message: fmt.Sprintf("can't read the scheme to check a configuration bacause %v", err),
		}

		return
	}

	c.RuleSet, err = c.walkRules(v)
	if err != nil {
		return
	}

	return
}

func (c *Checker) walkRules(v interface{}) (set RuleSet, err error) {
	set = RuleSet{}

	heap, ok := v.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("can't read json rules")

		return
	}

	var (
		subHeap       map[string]interface{}
		subName       string
		rule, subRule Rule
		subRuleSet    RuleSet
	)
	for name, val := range heap {
		subHeap, ok = val.(map[string]interface{})
		if !ok {
			err = &UnexpectedRule{
				message: fmt.Sprintf("can't parse scheme `%v`", val),
			}

			return
		}

		_, ok = subHeap["required"]
		if ok {
			rule, err = c.parseRule(subHeap)
			if err != nil {
				return
			}

			set[name] = rule
		}

		delete(subHeap, "required")
		delete(subHeap, "type")
		delete(subHeap, "regexp")
		delete(subHeap, "handler")

		if len(subHeap) > 0 {
			subRuleSet, err = c.walkRules(val)
			if err != nil {
				return
			}
			for subName, subRule = range subRuleSet {
				set[name+"."+subName] = subRule
			}
		}
	}

	return
}

func (c *Checker) parseRule(v map[string]interface{}) (rule Rule, err error) {
	var (
		ok  bool
		val interface{}
		str string
	)

	val, ok = v["required"]
	if !ok {
		err = &UnexpectedRule{
			message: fmt.Sprintf("can't parse rule `%v` because of `required` is absent", v),
		}

		return
	}

	rule.IsRequired, ok = val.(bool)
	if !ok {
		err = &UnexpectedRule{
			message: fmt.Sprintf("can't parse rule `%v` because of wrong `required`", v),
		}

		return
	}

	val, ok = v["type"]
	if !ok {
		err = &UnexpectedRule{
			message: fmt.Sprintf("can't parse rule `%v` because of `type` is absent", v),
		}

		return
	}

	str, ok = val.(string)
	if !ok {
		err = &UnexpectedRule{
			message: fmt.Sprintf("can't parse rule `%v` because of wrong `type`", v),
		}

		return
	}

	rule.Type, err = convType(str)
	if err != nil {
		return
	}

	val, ok = v["regexp"]
	if ok {
		str, ok = val.(string)
		if !ok {
			err = &UnexpectedRule{
				message: fmt.Sprintf("can't parse rule `%v` because of wrong `regexp`", v),
			}

			return
		}

		rule.RegExp, err = regexp.Compile(str)
		if err != nil {
			err = &UnexpectedRule{
				message: fmt.Sprintf("can't parse rule `%v` because of `regexp` can't be compiled", v),
			}

			return
		}
	}

	val, ok = v["handler"]
	if ok {
		str, ok = val.(string)
		if !ok {
			err = &UnexpectedRule{
				message: fmt.Sprintf("can't parse rule `%v` because of wrong `handler`", v),
			}

			return
		}

		rule.Handler, ok = c.HandlerBinds[str]
		if !ok {
			err = &UnexpectedHandler{
				message: fmt.Sprintf("can't parse handler `%v` because of `handler` isn't preset in the checker", v),
			}
		}
	}

	return
}

func convType(str string) (t Type, err error) {
	switch strings.ToLower(str) {
	case "string":
		t = StringType
	case "bool":
		t = BoolType
	case "int32":
		t = Int32Type
	case "uint32":
		t = UInt32Type
	case "int64":
		t = Int64Type
	case "uint64":
		t = UInt64Type
	case "float32":
		t = Float32Type
	case "float64":
		t = Float64Type
	case "list":
		t = ListType
	case "slice":
		t = SliceType
	case "map":
		t = MapType
	case "duration":
		t = DurationType
	default:
		err = &UnexpectedType{
			message: fmt.Sprintf("can't parse type `%v`", str),
		}
	}

	return
}

func (c *Checker) Check(obj Object) (err error) {
	var ok bool

	for path, rule := range c.RuleSet {
		ok = obj.IsExist(path)
		if !ok {
			if rule.IsRequired {
				err = &UnexpectedValue{
					message: fmt.Sprintf("path `%v` isn't set but `required`", path),
				}

				return
			} else {
				continue
			}
		}

		switch rule.Type {
		case StringType:
			ok = obj.IsString(path)
		case BoolType:
			ok = obj.IsBool(path)
		case Int32Type:
			ok = obj.IsInt32(path)
		case UInt32Type:
			ok = obj.IsUInt32(path)
		case Int64Type:
			ok = obj.IsInt64(path)
		case UInt64Type:
			ok = obj.IsUInt64(path)
		case Float32Type:
			ok = obj.IsFloat32(path)
		case Float64Type:
			ok = obj.IsFloat64(path)
		case ListType:
			ok = obj.IsList(path)
		case SliceType:
			ok = obj.IsSlice(path)
		case MapType:
			ok = obj.IsMap(path)
		case DurationType:
			ok = obj.IsDuration(path)
		default:
			ok = false
		}

		if !ok {
			err = &UnexpectedValue{
				message: fmt.Sprintf("path `%v` has wrong `type`", path),
			}

			return
		}

		if rule.RegExp != nil {
			var str string
			str, err = obj.String(path)
			if err != nil {
				err = &UnexpectedValue{
					message: fmt.Sprintf("path `%v` has wrong `type`", path),
				}

				return
			}

			ok = rule.RegExp.MatchString(str)
			if !ok {
				err = &UnexpectedValue{
					message: fmt.Sprintf("path `%v` has wrong value by `regexp`", path),
				}

				return
			}
		}

		if rule.Handler != nil {
			var val interface{}
			val, err = obj.Interface(path)
			if err != nil {
				err = &UnexpectedValue{
					message: fmt.Sprintf("path `%v` has wrong `type`", path),
				}

				return
			}

			err = rule.Handler(val)
			if err != nil {
				err = &UnexpectedValue{
					message: fmt.Sprintf("path `%v` has wrong value by `handler`: %v", path, err),
				}

				return
			}
		}
	}

	return
}
