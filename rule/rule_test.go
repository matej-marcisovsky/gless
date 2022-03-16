package rule

import (
	"gless/variable"
	"testing"
)

func TestApplyVariableToProperty(t *testing.T) {
	variable := variable.Variable{Name: "property", Value: "color"}
	var testCases = map[string]string{
		"@{property}":            "color",
		"background-@{property}": "background-color",
	}

	for property, result := range testCases {
		t.Run(property, func(t *testing.T) {
			rule := Rule{Property: property, Value: ""}
			rule.applyVariableToProperty(variable)

			if rule.Property != result {
				t.Fail()
			}
		})
	}
}

func TestApplyVariableToValue(t *testing.T) {
	variable := variable.Variable{Name: "value", Value: "red"}
	var testCases = map[string]string{
		"@value": "red",
	}

	for value, result := range testCases {
		t.Run(value, func(t *testing.T) {
			rule := Rule{Property: "", Value: value}
			rule.applyVariableToValue(variable)

			if rule.Value != result {
				t.Fail()
			}
		})
	}
}
