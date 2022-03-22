package expression

import (
	"errors"
	"fmt"
	"gless/chars"
	"gless/dimension"
	"gless/unit"
	"gless/utils"
	"gless/variable"
	"strings"
)

// http://mathcenter.oxford.emory.edu/site/cs171/shuntingYardAlgorithm/

var operators string = utils.JoinStrings(chars.ASTERISK, chars.MINUS, chars.PLUS, chars.SLASH)

type Expression string

// TODO strict units
func (expression Expression) Evaluate(variables map[string]variable.Variable) (string, error) {
	if !expression.isBalanced() {
		return "", errors.Unwrap(fmt.Errorf("expression '%s' is not balanced", string(expression)))
	}

	err := expression.processVariables(variables)

	if err != nil {
		return "", err
	}

	postfixNotation := expression.toPostfix()
	stack := make([]*dimension.Dimension, 0)
	finalUnit := unit.SINGULAR
	for _, token := range postfixNotation {
		currentUnit, isUnit := unit.GetUnit(token)

		if isUnit && finalUnit == unit.SINGULAR {
			finalUnit = currentUnit
		}

		if expression.isOperator(token) {
			last := stack[len(stack)-1]
			penultimate := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			dimension := dimension.Dimension{Unit: finalUnit}

			switch token {
			case chars.ASTERISK:
				dimension.Value = penultimate.Value * last.Value
			case chars.MINUS:
				dimension.Value = penultimate.Value - last.Value
			case chars.PLUS:
				dimension.Value = penultimate.Value + last.Value
			case chars.SLASH:
				dimension.Value = penultimate.Value / last.Value
			default:
				return "", errors.Unwrap(fmt.Errorf("unknown operator '%s'", token))
			}

			stack = append(stack, &dimension)
		} else {
			dimension, err := dimension.NewDimension(token)
			if err != nil {
				return "", err
			}

			stack = append(stack, dimension)
		}
	}

	if len(stack) != 1 {
		return "", errors.Unwrap(fmt.Errorf("invalid expression '%s'", string(expression)))
	}

	return stack[0].String(), nil
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

func (expression Expression) isOperator(token string) bool {
	return strings.ContainsAny(token, operators)
}

func (expression *Expression) processVariables(variables map[string]variable.Variable) error {
	foundVariables := variable.ExtractVariables(string(*expression))

	for len(foundVariables) > 0 {
		for _, variableName := range foundVariables {
			variable, variableExists := variables[variableName]

			if !variableExists {
				return errors.Unwrap(fmt.Errorf("variable '%s' not find", variableName))
			}

			*expression = Expression(strings.ReplaceAll(string(*expression), variable.Name, variable.Value))
		}

		foundVariables = variable.ExtractVariables(string(*expression))
	}

	return nil
}

func (expression Expression) toPostfix() []string {
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

		if len(token) == 1 && expression.isOperator(token) {
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
