package main

import "testing"

func TestApplyVariableToProperty(t *testing.T) {
	variable := Variable{name: "property", value: "color"}
	var testCases = map[string]string{
		"@{property}":            "color",
		"background-@{property}": "background-color",
	}

	for property, result := range testCases {
		t.Run(property, func(t *testing.T) {
			rule := Rule{property: property, value: ""}
			rule.applyVariableToProperty(variable)

			if rule.property != result {
				t.Fail()
			}
		})
	}
}

func TestApplyVariableToValue(t *testing.T) {
	variable := Variable{name: "value", value: "red"}
	var testCases = map[string]string{
		"@value": "red",
	}

	for value, result := range testCases {
		t.Run(value, func(t *testing.T) {
			rule := Rule{property: "", value: value}
			rule.applyVariableToValue(variable)

			if rule.value != result {
				t.Fail()
			}
		})
	}
}
