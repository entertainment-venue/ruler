package ruler

import (
	"encoding/json"
	"errors"
)

var (
	_ Ruler = new(andRuler)
	_ Ruler = new(orRuler)

	TypeMisMatch = errors.New("ruler type mismatch")
)

// Ruler 这个接口用来规定一组RuleSet以哪种规则匹配
type Ruler interface {
	AddRule(rule *Rule) Ruler
	Validate(msg map[string]interface{}) bool
	ValidateWithResult(msg map[string]interface{}) (map[string]*Result, bool)
}

// NewDefaultRuler 创建基础Ruler
func NewDefaultRuler(rules []*Rule, t string, decoder DataDecoder) (Ruler, error) {
	switch t {
	case "AND":
		return &andRuler{
			rules:   rules,
			decoder: decoder,
		}, nil
	case "OR":
		return &orRuler{
			rules:   rules,
			decoder: decoder,
		}, nil
	default:
		return nil, TypeMisMatch
	}
}

type orRuler struct {
	rules   []*Rule
	decoder DataDecoder
}

func (r *orRuler) AddRule(rule *Rule) Ruler {
	r.rules = append(r.rules, rule)
	return r
}

func (r *orRuler) Validate(msg map[string]interface{}) bool {
	for _, v := range r.rules {
		val, err := v.pluck(msg, v.path, r.decoder)
		if err != nil {
			return false
		}
		if v.Compare(val) {
			return true
		}
	}
	return false
}

func (r *orRuler) ValidateWithResult(msg map[string]interface{}) (map[string]*Result, bool) {
	resMap := make(map[string]*Result)
	flag := true
	for _, v := range r.rules {
		tempRes := &Result{
			expect: v.value,
			actual: nil,
			result: false,
			err:    nil,
		}
		temp, err := v.pluck(msg, v.path, r.decoder)
		if err != nil {
			flag = false
			tempRes.actual = temp
			tempRes.err = err
			tempRes.result = false
			resMap[v.path] = tempRes
			continue
		}
		f, err := v.CompareReturnResult(temp)
		if err != nil {
			flag = false
			tempRes.actual = temp
			tempRes.err = err
			tempRes.result = false
			resMap[v.path] = tempRes
			continue
		}
		if !f {
			flag = false
		}
		tempRes.actual = temp
		tempRes.err = nil
		tempRes.result = f
		resMap[v.path] = tempRes
	}
	return resMap, flag
}

type andRuler struct {
	rules   []*Rule
	decoder DataDecoder
}

func (r *andRuler) AddRule(rule *Rule) Ruler {
	r.rules = append(r.rules, rule)
	return r
}

func (r *andRuler) Validate(msg map[string]interface{}) bool {
	for _, v := range r.rules {
		val, err := v.pluck(msg, v.path, r.decoder)
		if err != nil {
			return false
		}
		if !v.Compare(val) {
			return false
		}
	}
	return true
}

func (r *andRuler) ValidateWithResult(msg map[string]interface{}) (map[string]*Result, bool) {
	resMap := make(map[string]*Result)
	flag := true
	for _, v := range r.rules {
		tempRes := &Result{
			expect: v.value,
			actual: nil,
			result: false,
			err:    nil,
		}
		temp, err := v.pluck(msg, v.path, r.decoder)
		if err != nil {
			flag = false
			tempRes.actual = temp
			tempRes.err = err
			tempRes.result = false
			resMap[v.path] = tempRes
			continue
		}
		f, err := v.CompareReturnResult(temp)
		if err != nil {
			flag = false
			tempRes.actual = temp
			tempRes.err = err
			tempRes.result = false
			resMap[v.path] = tempRes
			continue
		}
		if !f {
			flag = false
		}
		tempRes.actual = temp
		tempRes.err = nil
		tempRes.result = f
		resMap[v.path] = tempRes
	}
	return resMap, flag
}

// Result Rule校验结果
type Result struct {
	expect, actual interface{}
	result         bool
	err            error
}

// JsonString Rule校验结果生成Json类型字符串
func (r *Result) JsonString() string {
	bts, _ := json.Marshal(r)
	return string(bts)
}
