package ruler

import (
	"encoding/json"
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

func NewRule(path, t string, value interface{}) *Rule {
	return &Rule{
		comparator: getComparator(t),
		path:       path,
		value:      value,
		pluck:      getPluck(t),
	}
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
type Comparator func(actual, expected interface{}) bool

// Eq actual和expect相等
func Eq(actual, expect interface{}) bool {
	return actual == expect
}

// Neq actual和expect不相等
func Neq(actual, expect interface{}) bool {
	return !Eq(actual, expect)
}

// Gt actual大于expect
func Gt(actual, expect interface{}) bool {
	switch actual.(type) {
	case string:
		actualStr, ok := actual.(string)
		if !ok {
			return false
		}
		expectStr, ok := expect.(string)
		if !ok {
			return false
		}
		return actualStr > expectStr
	case float64:
		actualNum, ok := actual.(float64)
		if !ok {
			return false
		}
		expectNum, ok := expect.(float64)
		if !ok {
			return false
		}
		return actualNum > expectNum
	default:
		return false
	}
}

// Gte actual大于等于expect
func Gte(actual, expect interface{}) bool {
	switch actual.(type) {
	case string:
		actualStr, ok := actual.(string)
		if !ok {
			return false
		}
		expectStr, ok := expect.(string)
		if !ok {
			return false
		}
		return actualStr >= expectStr
	case float64:
		actualNum, ok := actual.(float64)
		if !ok {
			return false
		}
		expectNum, ok := expect.(float64)
		if !ok {
			return false
		}
		return actualNum >= expectNum
	default:
		return false
	}
}

// Lt actual小于expect
func Lt(actual, expect interface{}) bool {
	switch actual.(type) {
	case string:
		actualStr, ok := actual.(string)
		if !ok {
			return false
		}
		expectStr, ok := expect.(string)
		if !ok {
			return false
		}
		return actualStr < expectStr
	case float64:
		actualNum, ok := actual.(float64)
		if !ok {
			return false
		}
		expectNum, ok := expect.(float64)
		if !ok {
			return false
		}
		return actualNum < expectNum
	default:
		return false
	}
}

// Lte actual小于等于expect
func Lte(actual, expect interface{}) bool {
	switch actual.(type) {
	case string:
		actualStr, ok := actual.(string)
		if !ok {
			return false
		}
		expectStr, ok := expect.(string)
		if !ok {
			return false
		}
		return actualStr <= expectStr
	case float64:
		actualNum, ok := actual.(float64)
		if !ok {
			return false
		}
		expectNum, ok := expect.(float64)
		if !ok {
			return false
		}
		return actualNum <= expectNum
	default:
		return false
	}
}

// Exist actual不是nil
func Exist(actual, expect interface{}) bool {
	return actual != nil
}

// Nexist actual是nil
func Nexist(actual, expect interface{}) bool {
	return actual == nil
}

// Regex actual符合expect的正则匹配
func Regex(actual, expect interface{}) bool {
	return regular(actual, expect)
}

// Nregex actual不符合expect的正则匹配
func Nregex(actual, expect interface{}) bool {
	return !regular(actual, expect)
}

// Contains actual集合包含expect
func Contains(actual, expect interface{}) bool {
	switch bt := expect.(type) {
	case string:
		switch at := actual.(type) {
		case []interface{}:
			for _, v := range at {
				if elem, ok := v.(string); ok && elem == bt {
					return true
				}
			}
			return false
		case []string:
			for _, v := range at {
				if v == bt {
					return true
				}
			}
			return false
		case string:
			return strings.Contains(actual.(string), expect.(string))
		default:
			return false
		}
	case float64:
		switch at := actual.(type) {
		case []interface{}:
			for _, v := range at {
				if elem, ok := v.(float64); ok && elem == bt {
					return true
				}
			}
			return false
		case []float64:
			for _, v := range at {
				if v == bt {
					return true
				}
			}
		default:
			return false
		}
	default:
		return false
	}

	return false
}

// Ncontains actual集合不包含expect
func Ncontains(actual, expect interface{}) bool {
	switch bt := expect.(type) {
	case string:
		switch at := expect.(type) {
		case []interface{}:
			for _, v := range at {
				if elem, ok := v.(string); ok && elem == bt {
					return false
				}
			}
			return true
		case []string:
			for _, v := range at {
				if v == bt {
					return false
				}
			}
			return true
		case string:
			return !strings.Contains(actual.(string), expect.(string))
		default:
			return false
		}
	case float64:
		switch at := actual.(type) {
		case []interface{}:
			for _, v := range at {
				if elem, ok := v.(float64); ok && elem == bt {
					return false
				}
			}
			return true
		case []float64:
			for _, v := range at {
				if v == bt {
					return false
				}
			}
			return true
		default:
			return false
		}
	default:
		return false
	}
}

// OneOf actual在expect集合中
func OneOf(actual, expect interface{}) bool {
	switch expect.(type) {
	case []interface{}:
		for _, v := range expect.([]interface{}) {
			if v == actual {
				return true
			}
		}
		return false
	case map[interface{}]interface{}:
		m, ok := expect.(map[interface{}]struct{})
		if !ok {
			return false
		}
		_, found := m[actual]
		if found {
			return true
		}
		return false
	default:
		return false
	}
}

// NoneOf actual不在expect集合中
func NoneOf(actual, expect interface{}) bool {
	switch expect.(type) {
	case []interface{}:
		for _, v := range expect.([]interface{}) {
			if v == actual {
				return false
			}
		}
		return true
	case map[interface{}]interface{}:
		m, ok := expect.(map[interface{}]struct{})
		if !ok {
			return false
		}
		_, found := m[actual]
		if found {
			return false
		}
		return true
	default:
		return false
	}
}

// Startwith actual以expect开头
func Startwith(actual, expect interface{}) bool {
	var actualStr, expectStr string
	var ok bool
	if actualStr, ok = actual.(string); !ok {
		return false
	}
	if expectStr, ok = expect.(string); !ok {
		return false
	}
	return strings.HasPrefix(actualStr, expectStr)
}

// Nstartwith actual不以expect开头
func Nstartwith(actual, expect interface{}) bool {
	var actualStr, expectStr string
	var ok bool
	if actualStr, ok = actual.(string); !ok {
		return false
	}
	if expectStr, ok = expect.(string); !ok {
		return false
	}
	return !strings.HasPrefix(actualStr, expectStr)
}

// Endwith actual以expect结尾
func Endwith(actual, expect interface{}) bool {
	var actualStr, expectStr string
	var ok bool
	if actualStr, ok = actual.(string); !ok {
		return false
	}
	if expectStr, ok = expect.(string); !ok {
		return false
	}
	return strings.HasSuffix(actualStr, expectStr)
}

// Nendwith actual不以expect结尾
func Nendwith(actual, expect interface{}) bool {
	var actualStr, expectStr string
	var ok bool
	if actualStr, ok = actual.(string); !ok {
		return false
	}
	if expectStr, ok = expect.(string); !ok {
		return false
	}
	return !strings.HasSuffix(actualStr, expectStr)
}

// regular 正则匹配校验
func regular(actual, expect interface{}) bool {
	switch actual.(type) {
	case string:
		actualStr, ok := actual.(string)
		if !ok {
			return false
		}
		expectStr, ok := expect.(string)
		if !ok {
			return false
		}

		r, err := regexp.Compile(expectStr)
		if err != nil {
			return false
		}

		return r.MatchString(actualStr)
	default:
		return false
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
type Pluck func(props map[string]interface{}, path string, decoder func(bts []byte, m *map[string]interface{}) error) interface{}

// PathPluck 根据路径解析map值
func PathPluck(props map[string]interface{}, path string, decoder func(bts []byte, m *map[string]interface{}) error) interface{} {
	parts := strings.Split(path, ".")
	for i := 0; i < len(parts)-1; i++ {
		prev := make(map[string]interface{})
		switch props[parts[i]].(type) {
		case string:
			if err := decoder([]byte(props[parts[i]].(string)), &prev); err != nil {
				return nil
			}
		case map[string]interface{}:
			prev = props[parts[i]].(map[string]interface{})
		default:
			return nil
		}
		props = prev
	}
	return props[parts[len(parts)-1]]
}

// JsonPathPluck 根据路径基于jsonpath解析map值
func JsonPathPluck(props map[string]interface{}, path string, decoder func(bts []byte, m *map[string]interface{}) error) interface{} {
	jpPath, err := jp.ParseString(path)
	if err != nil {
		return nil
	}
	bts, err := json.Marshal(props)
	if err != nil {
		return nil
	}
	obj, err := oj.ParseString(string(bts))
	if err != nil {
		return nil
	}
	return jpPath.First(obj)
}
