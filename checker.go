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
	ArrayType    Type = "array"
	JSONType     Type = "json"
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
			message: fmt.Sprintf("can't read the scheme to check some configuration bacause %v", err),
		}

		return
	}

	ruleList, err := c.walkJson(v)
	if err != nil {
		return
	}

	c.RuleSet = ruleList

	return
}

func (c *Checker) walkJson(v interface{}) (set RuleSet, err error) {
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
			subRuleSet, err = c.walkJson(val)
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
	case "array":
		t = ArrayType
	case "json":
		t = JSONType
	case "duration":
		t = DurationType
	default:
		err = &UnexpectedType{
			message: fmt.Sprintf("can't parse type `%v`", str),
		}
	}

	return
}

func (c *Checker) Check(bs []byte) (err error) {
	var ok bool

	var result Result
	result, err = ParseJson(bs)
	if err != nil {
		return
	}

	for path, rule := range c.RuleSet {
		ok = result.IsExist(path)
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
			ok = result.IsString(path)
		case BoolType:
			ok = result.IsBool(path)
		case Int32Type:
			ok = result.IsInt32(path)
		case UInt32Type:
			ok = result.IsUInt32(path)
		case Int64Type:
			ok = result.IsInt64(path)
		case UInt64Type:
			ok = result.IsUInt64(path)
		case Float32Type:
			ok = result.IsFloat32(path)
		case Float64Type:
			ok = result.IsFloat64(path)
		case ArrayType:
			ok = result.IsArray(path)
		case JSONType:
			ok = result.IsJSON(path)
		case DurationType:
			ok = result.IsDuration(path)
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
			str, err = result.String(path)
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
			val, err = result.Interface(path)
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
