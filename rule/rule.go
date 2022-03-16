package rule

import (
	"gless/chars"
	"gless/utils"
	"gless/variable"
	"strings"
)

type Rule struct {
	Property string
	Value    string
}

func (rule *Rule) Process(variables []variable.Variable) {
	propertyContainsVariables := strings.Contains(rule.Property, utils.JoinStrings(chars.ASPERAND, chars.CURLY_BRACKET_OPEN))
	valueContainsVariables := strings.Contains(rule.Value, chars.ASPERAND)

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

func (rule *Rule) applyVariableToProperty(variable variable.Variable) {
	rule.Property = strings.ReplaceAll(rule.Property, utils.JoinStrings(chars.ASPERAND, chars.CURLY_BRACKET_OPEN, variable.Name, chars.CURLY_BRACKET_CLOSE), variable.Value)
}

func (rule *Rule) applyVariableToValue(variable variable.Variable) {
	rule.Value = strings.ReplaceAll(rule.Value, utils.JoinStrings(chars.ASPERAND, variable.Name), variable.Value)
}
