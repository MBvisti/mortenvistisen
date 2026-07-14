package validation

import "testing"

func TestValidationBuilderRules(t *testing.T) {
	b := NewBuilder()
	b.Required("title", "")
	b.MaxLen("title", "", 60)
	b.RecommendedLenBetween("title", 50, 60)
	b.URL("imageLink", "")
	b.MinInt("readTime", 0, 0)

	rules := b.Rules()

	assertRuleParam(t, rules, "title", "max", "max", 60)
	assertRuleParam(t, rules, "title", "recommended_length", "min", 50)
	assertRuleParam(t, rules, "title", "recommended_length", "max", 60)
	assertRuleParam(t, rules, "readTime", "min", "min", int64(0))
	assertRule(t, rules, "title", "required")
	assertRule(t, rules, "imageLink", "url")
}

func assertRule(t *testing.T, rules Rules, field string, code string) Rule {
	t.Helper()

	for _, rule := range rules[field] {
		if rule.Code == code {
			return rule
		}
	}

	t.Fatalf("Rules() missing rule %q for field %q", code, field)
	return Rule{}
}

func assertRuleParam(
	t *testing.T,
	rules Rules,
	field string,
	code string,
	param string,
	want any,
) {
	t.Helper()

	rule := assertRule(t, rules, field, code)
	if got := rule.Params[param]; got != want {
		t.Fatalf(
			"Rules()[%q] rule %q param %q = %#v, want %#v",
			field,
			code,
			param,
			got,
			want,
		)
	}
}
