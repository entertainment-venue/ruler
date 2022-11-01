package ruler

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

// Rule 规则
// comparator 是比较器, 用于比较两个值是否相等
// path 基于Map的取值路径
// value 期待值
// pluck 取值方式
type Rule struct {
	comparator Comparator
	path       string
	value      interface{}
	pluck      Pluck
}

var (
	TypeMismatchErr   = errors.New("type mismatch error")
	TypeNotSupportErr = errors.New("type not support error")
)

// NewRule 根据条件创建Rule
func NewRule(path, t string, value interface{}) *Rule {
	return NewRuleByOptions(path, value, getComparator(t), getPluck(t))
}

// NewRuleByOptions 自定义Rule
func NewRuleByOptions(path string, value interface{}, comparator Comparator, pluck Pluck) *Rule {
	return &Rule{
		comparator: comparator,
		path:       path,
		value:      value,
		pluck:      pluck,
	}
}

// Compare 规则比较
func (r *Rule) Compare(actual interface{}) bool {
	flag, err := r.compare(actual)
	if err != nil {
		return false
	}
	return flag
}

// CompareReturnResult 比较结果, 并返回error
func (r *Rule) CompareReturnResult(actual interface{}) (bool, error) {
	return r.compare(actual)
}

func (r *Rule) compare(actual interface{}) (bool, error) {
	return r.comparator(actual, r.value)
}

func getComparator(t string) Comparator {
	switch strings.ToUpper(t) {
	case "EQ", "JEQ":
		return Eq
	case "NEQ", "JNEQ":
		return Neq
	case "GT":
		return Gt
	case "GTE":
		return Gte
	case "LT":
		return Lt
	case "LTE":
		return Lte
	case "EXISTS":
		return Exist
	case "NEXISTS":
		return Nexist
	case "REGEX":
		return Regex
	case "NREGEX":
		return Nregex
	case "CONTAINS":
		return Contains
	case "NCONTAINS":
		return Ncontains
	case "ONEOF":
		return OneOf
	case "NONEOF":
		return NoneOf
	case "STARTWITH":
		return Startwith
	case "NSTARTWITH":
		return Nstartwith
	case "ENDWITH":
		return Endwith
	case "NENDWITH":
		return Nendwith
	default:
		return Eq
	}
}

// Comparator 是比较器 目前实现了一些基本的比较器 可以根据需求来实现相应的比较器函数传入Rule
type Comparator func(actual, expected interface{}) (bool, error)

// Eq actual和expect相等
func Eq(actual, expect interface{}) (bool, error) {
	return actual == expect, nil
}

// Neq actual和expect不相等
func Neq(actual, expect interface{}) (bool, error) {
	flag, err := Eq(actual, expect)
	return !flag, err
}

// Gt actual大于expect
func Gt(actual, expect interface{}) (bool, error) {
	switch actual.(type) {
	case string:
		actualStr, ok := actual.(string)
		if !ok {
			return false, TypeMismatchErr
		}
		expectStr, ok := expect.(string)
		if !ok {
			return false, TypeMismatchErr
		}
		return actualStr > expectStr, nil
	case float64:
		actualNum, ok := actual.(float64)
		if !ok {
			return false, TypeMismatchErr
		}
		expectNum, ok := expect.(float64)
		if !ok {
			return false, TypeMismatchErr
		}
		return actualNum > expectNum, nil
	default:
		return false, TypeNotSupportErr
	}
}

// Gte actual大于等于expect
func Gte(actual, expect interface{}) (bool, error) {
	switch actual.(type) {
	case string:
		actualStr, ok := actual.(string)
		if !ok {
			return false, TypeMismatchErr
		}
		expectStr, ok := expect.(string)
		if !ok {
			return false, TypeMismatchErr
		}
		return actualStr >= expectStr, nil
	case float64:
		actualNum, ok := actual.(float64)
		if !ok {
			return false, TypeMismatchErr
		}
		expectNum, ok := expect.(float64)
		if !ok {
			return false, TypeMismatchErr
		}
		return actualNum >= expectNum, nil
	default:
		return false, TypeNotSupportErr
	}
}

// Lt actual小于expect
func Lt(actual, expect interface{}) (bool, error) {
	switch actual.(type) {
	case string:
		actualStr, ok := actual.(string)
		if !ok {
			return false, TypeMismatchErr
		}
		expectStr, ok := expect.(string)
		if !ok {
			return false, TypeMismatchErr
		}
		return actualStr < expectStr, nil
	case float64:
		actualNum, ok := actual.(float64)
		if !ok {
			return false, TypeMismatchErr
		}
		expectNum, ok := expect.(float64)
		if !ok {
			return false, TypeMismatchErr
		}
		return actualNum < expectNum, nil
	default:
		return false, TypeNotSupportErr
	}
}

// Lte actual小于等于expect
func Lte(actual, expect interface{}) (bool, error) {
	switch actual.(type) {
	case string:
		actualStr, ok := actual.(string)
		if !ok {
			return false, TypeMismatchErr
		}
		expectStr, ok := expect.(string)
		if !ok {
			return false, TypeMismatchErr
		}
		return actualStr <= expectStr, nil
	case float64:
		actualNum, ok := actual.(float64)
		if !ok {
			return false, TypeMismatchErr
		}
		expectNum, ok := expect.(float64)
		if !ok {
			return false, TypeMismatchErr
		}
		return actualNum <= expectNum, nil
	default:
		return false, TypeNotSupportErr
	}
}

// Exist actual不是nil
func Exist(actual, expect interface{}) (bool, error) {
	return actual != nil, nil
}

// Nexist actual是nil
func Nexist(actual, expect interface{}) (bool, error) {
	return actual == nil, nil
}

// Regex actual符合expect的正则匹配
func Regex(actual, expect interface{}) (bool, error) {
	return regular(actual, expect)
}

// Nregex actual不符合expect的正则匹配
func Nregex(actual, expect interface{}) (bool, error) {
	flag, err := regular(actual, expect)
	return !flag, err
}

// Contains actual集合包含expect
func Contains(actual, expect interface{}) (bool, error) {
	switch bt := expect.(type) {
	case string:
		switch at := actual.(type) {
		case []interface{}:
			var err error
			for _, v := range at {
				if elem, ok := v.(string); ok && elem == bt {
					return true, nil
				} else if !ok {
					err = TypeMismatchErr
				}
			}
			return false, err
		case []string:
			for _, v := range at {
				if v == bt {
					return true, nil
				}
			}
			return false, nil
		case string:
			return strings.Contains(actual.(string), expect.(string)), nil
		default:
			return false, TypeNotSupportErr
		}
	case float64:
		switch at := actual.(type) {
		case []interface{}:
			var err error
			for _, v := range at {
				if elem, ok := v.(float64); ok && elem == bt {
					return true, nil
				} else if !ok {
					err = TypeMismatchErr
				}
			}
			return false, err
		case []float64:
			for _, v := range at {
				if v == bt {
					return true, nil
				}
			}
		default:
			return false, TypeNotSupportErr
		}
	default:
		return false, TypeNotSupportErr
	}
	return false, nil
}

// Ncontains actual集合不包含expect
func Ncontains(actual, expect interface{}) (bool, error) {
	switch bt := expect.(type) {
	case string:
		switch at := expect.(type) {
		case []interface{}:
			var err error
			for _, v := range at {
				if elem, ok := v.(string); ok && elem == bt {
					return false, nil
				} else if !ok {
					err = TypeMismatchErr
				}
			}
			return true, err
		case []string:
			for _, v := range at {
				if v == bt {
					return false, nil
				}
			}
			return true, nil
		case string:
			return !strings.Contains(actual.(string), expect.(string)), nil
		default:
			return false, TypeNotSupportErr
		}
	case float64:
		switch at := actual.(type) {
		case []interface{}:
			var err error
			for _, v := range at {
				if elem, ok := v.(float64); ok && elem == bt {
					return false, nil
				} else if !ok {
					err = TypeMismatchErr
				}
			}
			return true, err
		case []float64:
			for _, v := range at {
				if v == bt {
					return false, nil
				}
			}
			return true, nil
		default:
			return false, TypeNotSupportErr
		}
	default:
		return false, TypeNotSupportErr
	}
}

// OneOf actual在expect集合中
func OneOf(actual, expect interface{}) (bool, error) {
	switch expect.(type) {
	case []interface{}:
		for _, v := range expect.([]interface{}) {
			if v == actual {
				return true, nil
			}
		}
		return false, nil
	case map[interface{}]interface{}:
		m, ok := expect.(map[interface{}]struct{})
		if !ok {
			return false, TypeMismatchErr
		}
		_, found := m[actual]
		if found {
			return true, nil
		}
		return false, nil
	default:
		return false, TypeNotSupportErr
	}
}

// NoneOf actual不在expect集合中
func NoneOf(actual, expect interface{}) (bool, error) {
	switch expect.(type) {
	case []interface{}:
		for _, v := range expect.([]interface{}) {
			if v == actual {
				return false, nil
			}
		}
		return true, nil
	case map[interface{}]interface{}:
		m, ok := expect.(map[interface{}]struct{})
		if !ok {
			return false, TypeMismatchErr
		}
		_, found := m[actual]
		if found {
			return false, nil
		}
		return true, nil
	default:
		return false, TypeNotSupportErr
	}
}

// Startwith actual以expect开头
func Startwith(actual, expect interface{}) (bool, error) {
	var actualStr, expectStr string
	var ok bool
	if actualStr, ok = actual.(string); !ok {
		return false, TypeMismatchErr
	}
	if expectStr, ok = expect.(string); !ok {
		return false, TypeMismatchErr
	}
	return strings.HasPrefix(actualStr, expectStr), nil
}

// Nstartwith actual不以expect开头
func Nstartwith(actual, expect interface{}) (bool, error) {
	var actualStr, expectStr string
	var ok bool
	if actualStr, ok = actual.(string); !ok {
		return false, TypeMismatchErr
	}
	if expectStr, ok = expect.(string); !ok {
		return false, TypeMismatchErr
	}
	return !strings.HasPrefix(actualStr, expectStr), nil
}

// Endwith actual以expect结尾
func Endwith(actual, expect interface{}) (bool, error) {
	var actualStr, expectStr string
	var ok bool
	if actualStr, ok = actual.(string); !ok {
		return false, TypeMismatchErr
	}
	if expectStr, ok = expect.(string); !ok {
		return false, TypeMismatchErr
	}
	return strings.HasSuffix(actualStr, expectStr), nil
}

// Nendwith actual不以expect结尾
func Nendwith(actual, expect interface{}) (bool, error) {
	var actualStr, expectStr string
	var ok bool
	if actualStr, ok = actual.(string); !ok {
		return false, TypeMismatchErr
	}
	if expectStr, ok = expect.(string); !ok {
		return false, TypeMismatchErr
	}
	return !strings.HasSuffix(actualStr, expectStr), nil
}

// regular 正则匹配校验
func regular(actual, expect interface{}) (bool, error) {
	switch actual.(type) {
	case string:
		actualStr, ok := actual.(string)
		if !ok {
			return false, TypeMismatchErr
		}
		expectStr, ok := expect.(string)
		if !ok {
			return false, TypeMismatchErr
		}

		r, err := regexp.Compile(expectStr)
		if err != nil {
			return false, err
		}
		return r.MatchString(actualStr), nil
	default:
		return false, TypeNotSupportErr
	}
}

func getPluck(t string) Pluck {
	switch strings.ToUpper(t) {
	case "JEQ", "JNEQ":
		return JsonPathPluck
	default:
		return PathPluck
	}
}

// Pluck 根据路径嵌套解析map值
type Pluck func(props map[string]interface{}, path string, decoder DataDecoder) (interface{}, error)

// PathPluck 根据路径解析map值
func PathPluck(props map[string]interface{}, path string, decoder DataDecoder) (interface{}, error) {
	parts := strings.Split(path, ".")
	for i := 0; i < len(parts)-1; i++ {
		prev := make(map[string]interface{})
		switch props[parts[i]].(type) {
		case string:
			if err := decoder([]byte(props[parts[i]].(string)), &prev); err != nil {
				return nil, err
			}
		case map[string]interface{}:
			prev = props[parts[i]].(map[string]interface{})
		default:
			return nil, TypeNotSupportErr
		}
		props = prev
	}
	return props[parts[len(parts)-1]], nil
}

// JsonPathPluck 根据路径基于jsonpath解析map值
func JsonPathPluck(props map[string]interface{}, path string, decoder DataDecoder) (interface{}, error) {
	jpPath, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}
	bts, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}
	obj, err := oj.ParseString(string(bts))
	if err != nil {
		return nil, err
	}
	return jpPath.First(obj), nil
}
