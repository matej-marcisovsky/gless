package expression

import (
	"errors"
	"fmt"
	"gless/chars"
	"gless/unit"
	"gless/utils"
	"gless/variable"
	"strings"
)

// http://mathcenter.oxford.emory.edu/site/cs171/shuntingYardAlgorithm/

type Expression string

func (expression Expression) Evaluate(variables map[string]variable.Variable) (string, error) {
	if !expression.isBalanced() {
		return "", errors.Unwrap(fmt.Errorf("Expression '%s' is not balanced", string(expression)))
	}

	err := expression.processVariables(variables)

	if err != nil {
		return "", err
	}

	postfixNotation := expression.toPostfix()
	finalUnit := unit.NUMBER
	for _, notation := range postfixNotation {
		_unit, isUnit := unit.GetUnit(notation)

		if isUnit {
			finalUnit = _unit
			break
		}
	}

	fmt.Println(finalUnit)

	return strings.Join(postfixNotation, chars.SPACE), nil
}

func (expression Expression) isBalanced() bool {
	stack := 0

	for _, char := range string(expression) {
		stringChar := string(char)

		if stringChar == chars.PARENTHESIS_CLOSE && stack == 0 {
			return false
		}

		if stringChar == chars.PARENTHESIS_OPEN {
			stack++
		}

		if stringChar == chars.PARENTHESIS_CLOSE {
			stack--
		}
	}

	return stack == 0
}

func (expression Expression) isHigherPrecedence(a string, b string) bool {
	higherPrecedenceOperators := utils.JoinStrings(chars.ASTERISK, chars.SLASH)

	return strings.ContainsAny(a, higherPrecedenceOperators) && !strings.ContainsAny(b, higherPrecedenceOperators)
}

func (expression Expression) isSamePrecedence(a string, b string) bool {
	higherPrecedenceOperators := utils.JoinStrings(chars.ASTERISK, chars.SLASH)

	return (strings.ContainsAny(a, higherPrecedenceOperators) && strings.ContainsAny(b, higherPrecedenceOperators)) ||
		(!strings.ContainsAny(a, higherPrecedenceOperators) && !strings.ContainsAny(b, higherPrecedenceOperators))
}

func (expression *Expression) processVariables(variables map[string]variable.Variable) error {
	foundVariables := variable.ExtractVariables(string(*expression))

	for len(foundVariables) > 0 {
		for _, variableName := range foundVariables {
			variable, variableExists := variables[variableName]

			if !variableExists {
				return errors.Unwrap(fmt.Errorf("Variable '%s' not find", variableName))
			}

			*expression = Expression(strings.ReplaceAll(string(*expression), variable.Name, variable.Value))
		}

		foundVariables = variable.ExtractVariables(string(*expression))
	}

	return nil
}

func (expression Expression) toPostfix() []string {
	operators := utils.JoinStrings(chars.ASTERISK, chars.MINUS, chars.PLUS, chars.SLASH)
	operatorStack := make([]string, 0)
	queue := make([]string, 0)

	expression = Expression(strings.ReplaceAll(string(expression), chars.PARENTHESIS_OPEN, utils.JoinStrings(chars.PARENTHESIS_OPEN, chars.SPACE)))
	expression = Expression(strings.ReplaceAll(string(expression), chars.PARENTHESIS_CLOSE, utils.JoinStrings(chars.SPACE, chars.PARENTHESIS_CLOSE)))

	for _, token := range strings.Split(string(expression), chars.SPACE) {
		if token == chars.PARENTHESIS_OPEN {
			operatorStack = append(operatorStack, token)
			continue
		}

		if token == chars.PARENTHESIS_CLOSE {
			for i := len(operatorStack) - 1; i >= 0; i-- {
				operator := operatorStack[i]

				if operator == chars.PARENTHESIS_OPEN {
					operatorStack = operatorStack[0:i]
					break
				}

				queue = append(queue, operator)
			}
			continue
		}

		if len(token) == 1 && strings.ContainsAny(token, operators) {
			if len(operatorStack) == 0 {
				operatorStack = append(operatorStack, token)
				continue
			}

			lastOperator := operatorStack[len(operatorStack)-1]

			if lastOperator == chars.PARENTHESIS_OPEN {
				operatorStack = append(operatorStack, token)
				continue
			}

			if expression.isHigherPrecedence(token, lastOperator) {
				operatorStack = append(operatorStack, token)
				continue
			}

			if !expression.isHigherPrecedence(token, lastOperator) || expression.isSamePrecedence(token, lastOperator) {
				for i := len(operatorStack) - 1; i >= 0; i-- {
					operator := operatorStack[i]

					if operator == chars.PARENTHESIS_OPEN || operator == chars.PARENTHESIS_CLOSE {
						break
					}

					queue = append(queue, operator)
					operatorStack = operatorStack[0:i]

					if expression.isHigherPrecedence(token, operator) {
						break
					}
				}

				operatorStack = append(operatorStack, token)
				continue
			}
		} else {
			queue = append(queue, token)
		}
	}

	queue = append(queue, operatorStack...)

	return queue
}
