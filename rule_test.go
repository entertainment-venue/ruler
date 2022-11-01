package ruler

import (
	"testing"
)

type TestingCase struct {
	rule   Rule
	actual interface{}
	result bool
}

var TestingArr []TestingCase

func init() {
	case1 := TestingCase{
		rule: Rule{
			comparator: Eq,
			path:       "1",
			value:      "xxx",
		},
		actual: "xxx",
		result: true,
	}
	case2 := TestingCase{
		rule: Rule{
			comparator: Eq,
			path:       "2",
			value:      1,
		},
		actual: 2,
		result: false,
	}
	case3 := TestingCase{
		rule: Rule{
			comparator: Neq,
			path:       "3",
			value:      "xxxx",
		},
		actual: "xxx",
		result: true,
	}
	case4 := TestingCase{
		rule: Rule{
			comparator: Neq,
			path:       "4",
			value:      2,
		},
		actual: 2,
		result: false,
	}
	case5 := TestingCase{
		rule: Rule{
			comparator: Gt,
			path:       "5",
			value:      2,
		},
		actual: 3,
		result: true,
	}
	case6 := TestingCase{
		rule: Rule{
			comparator: Gt,
			path:       "6",
			value:      2,
		},
		actual: 2,
		result: false,
	}
	case7 := TestingCase{
		rule: Rule{
			comparator: Gte,
			path:       "7",
			value:      3,
		},
		actual: 3,
		result: true,
	}
	case8 := TestingCase{
		rule: Rule{
			comparator: Gte,
			path:       "8",
			value:      2,
		},
		actual: 1,
		result: false,
	}
	case9 := TestingCase{
		rule: Rule{
			comparator: Lt,
			path:       "9",
			value:      2,
		},
		actual: 1,
		result: true,
	}
	case10 := TestingCase{
		rule: Rule{
			comparator: Lt,
			path:       "10",
			value:      4,
		},
		actual: 4,
		result: false,
	}
	case11 := TestingCase{
		rule: Rule{
			comparator: Lte,
			path:       "11",
			value:      4,
		},
		actual: 4,
		result: true,
	}
	case12 := TestingCase{
		rule: Rule{
			comparator: Lte,
			path:       "12",
			value:      4,
		},
		actual: 5,
		result: false,
	}
	case13 := TestingCase{
		rule: Rule{
			comparator: Exist,
			path:       "13",
			value:      nil,
		},
		actual: nil,
		result: false,
	}
	case14 := TestingCase{
		rule: Rule{
			comparator: Exist,
			path:       "14",
			value:      "xxx",
		},
		actual: "xxx",
		result: true,
	}
	case15 := TestingCase{
		rule: Rule{
			comparator: Nexist,
			path:       "15",
			value:      nil,
		},
		actual: nil,
		result: true,
	}
	case16 := TestingCase{
		rule: Rule{
			comparator: Nexist,
			path:       "16",
			value:      "xxx",
		},
		actual: "xxx",
		result: false,
	}
	case17 := TestingCase{
		rule: Rule{
			comparator: Regex,
			path:       "17",
			value:      `^1[34578]\d{9}$`,
		},
		actual: "18612045500",
		result: true,
	}
	case18 := TestingCase{
		rule: Rule{
			comparator: Regex,
			path:       "18",
			value:      "/^400[0-9]{7}/",
		},
		actual: "13823292292",
		result: false,
	}
	case19 := TestingCase{
		rule: Rule{
			comparator: Nregex,
			path:       "19",
			value:      "/^400[0-9]{7}/",
		},
		actual: "5001234567",
		result: true,
	}
	case20 := TestingCase{
		rule: Rule{
			comparator: Nregex,
			path:       "20",
			value:      "/^400[0-9]{7}/",
		},
		actual: "5001234567",
		result: true,
	}
	TestingArr = append(TestingArr, case1, case2, case3, case4, case5, case6, case7, case8, case9, case10, case11, case12, case13, case14, case15, case16, case17, case18, case19, case20)
}

func Test_Rule(t *testing.T) {
	for _, v := range TestingArr {
		if v.result != v.rule.Compare(v.actual) {
			t.Errorf("error case %s", v.rule.path)
		}
	}
}
