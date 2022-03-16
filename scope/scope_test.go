package scope

import "testing"

func TestComputeSelector(t *testing.T) {
	const PARENT_SELECTOR = ".link"
	var testCases = map[string]string{
		"& + &":  ".link + .link",
		"& &":    ".link .link",
		"&&":     ".link.link",
		"&,&ish": ".link,\n.linkish",
	}

	for selector, result := range testCases {
		t.Run(selector, func(t *testing.T) {
			scope := Scope{Selector: selector}

			if scope.computeSelector(PARENT_SELECTOR) != result {
				t.Fail()
			}
		})
	}
}
