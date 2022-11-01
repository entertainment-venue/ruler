package ruler

type RuleSetType string

// RuleSet 规则组
type RuleSet struct {
	rules   *[]*Rule
	setType RuleSetType
}
