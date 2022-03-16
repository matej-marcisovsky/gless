package expression

import (
	"gless/chars"
	"gless/variable"
	"strings"
	"testing"
)

func TestIsBalanced(t *testing.T) {
	expression := Expression("(3 * 2 / (2 + 1))")

	if !expression.isBalanced() {
		t.Fail()
	}
}

func TestEvaluate(t *testing.T) {
	variables := make(map[string]variable.Variable)
	variables["@a"] = variable.Variable{Name: "@a", Value: "3px"}
	variables["@b"] = variable.Variable{Name: "@b", Value: "@c + 5px"}
	variables["@c"] = variable.Variable{Name: "@c", Value: "8px"}

	expression := Expression("@a * (@b + C * D) + E")
	evaluate, err := expression.Evaluate(variables)

	if err != nil {
		t.Fail()
	}

	if evaluate != "3px 8px 5px + C D * + * E +" {
		t.Log(evaluate)
		t.Fail()
	}
}

func TestToPostfix(t *testing.T) {
	expression := Expression("3px * (8px + 5px + C * D) + E")
	postfix := strings.Join(expression.toPostfix(), chars.SPACE)

	if postfix != "3px 8px 5px + C D * + * E +" {
		t.Log(postfix)
		t.Fail()
	}
}
