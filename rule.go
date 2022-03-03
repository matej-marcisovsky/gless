package main

import (
	"strings"
)

type Rule struct {
	property string
	value    string
}

func (rule *Rule) Process(variables []Variable) {
	propertyContainsVariables := strings.Contains(rule.property, JoinStrings(ASPERAND, CURLY_BRACKET_OPEN))
	valueContainsVariables := strings.Contains(rule.value, ASPERAND)

	if propertyContainsVariables || valueContainsVariables {
		for _, variable := range variables {
			if propertyContainsVariables {
				rule.applyVariableToProperty(variable)
			}

			if valueContainsVariables {
				rule.applyVariableToValue(variable)
			}
		}
	}
}

func (rule *Rule) applyVariableToProperty(variable Variable) {
	rule.property = strings.ReplaceAll(rule.property, JoinStrings(ASPERAND, CURLY_BRACKET_OPEN, variable.name, CURLY_BRACKET_CLOSE), variable.value)
}

func (rule *Rule) applyVariableToValue(variable Variable) {
	rule.value = strings.ReplaceAll(rule.value, JoinStrings(ASPERAND, variable.name), variable.value)
}
