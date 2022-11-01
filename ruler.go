package ruler

// Ruler 这个接口用来规定一组RuleSet以哪种规则匹配
type Ruler interface {
	Validate(msg map[string]interface{}) bool
	ValidateWithRes(msg map[string]interface{}) (map[string]interface{}, bool)
	//Type() RulerType
	AddRule(rule Rule) Ruler
}
