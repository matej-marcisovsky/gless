package main

import (
	"regexp"
	"sort"
	"strings"
	"sync"
)

var regexpCommaSpace *regexp.Regexp = regexp.MustCompile(JoinStrings(COMMA, "\\s{0,1}"))

type Scope struct {
	parent    *Scope
	result    string
	rules     []Rule
	scopes    []*Scope
	selector  string
	variables []Variable
}

func (scope *Scope) AddVariable(variable Variable) {
	scope.variables = append(scope.variables, variable)
}

func (scope *Scope) AddRule(rule Rule) {
	scope.rules = append(scope.rules, rule)
}

func (scope *Scope) AddScope(_scope *Scope) {
	scope.scopes = append(scope.scopes, _scope)
}

func (scope *Scope) Process(waitGroup *sync.WaitGroup, parentSelector string, variables []Variable) {
	scope.prepareVariables(variables)
	selector := scope.computeSelector(parentSelector)

	if scope.selector != "" && len(scope.rules) > 0 {
		scope.result = JoinStrings(selector, SPACE, CURLY_BRACKET_OPEN, NEW_LINE)
	}

	for _, rule := range scope.rules {
		rule.Process(scope.variables)
		scope.result = scope.result + JoinStrings(TABULATOR, rule.property, COLON, SPACE, rule.value, SEMICOLON, NEW_LINE)
	}

	for _, _scope := range scope.scopes {
		waitGroup.Add(1)
		go _scope.Process(waitGroup, selector, scope.variables)
	}

	if scope.selector != "" && len(scope.rules) > 0 {
		scope.result = scope.result + JoinStrings(CURLY_BRACKET_CLOSE, NEW_LINE)
	}

	defer waitGroup.Done()
}

func (scope *Scope) String() string {
	result := scope.result

	for _, _scope := range scope.scopes {
		result = result + _scope.String()
	}

	return result
}

func (scope *Scope) computeSelector(parentSelector string) string {
	result := strings.ReplaceAll(scope.selector, AMPERSAND, parentSelector)

	if strings.Contains(result, COMMA) {
		result = regexpCommaSpace.ReplaceAllString(result, JoinStrings(COMMA, NEW_LINE))
	}

	if !strings.HasPrefix(scope.selector, AMPERSAND) && len(parentSelector) > 0 {
		result = JoinStrings(parentSelector, SPACE, result)
	}

	return result
}

func (scope *Scope) prepareVariables(variables []Variable) {
	scopeVariablesNames := make(map[string]struct{}, len(scope.variables))

	for _, variable := range scope.variables {
		scopeVariablesNames[variable.name] = struct{}{}
	}

	for _, variable := range variables {
		if _, ok := scopeVariablesNames[variable.name]; !ok {
			scope.variables = append(scope.variables, variable)
		}
	}

	sort.Slice(scope.variables, func(i, j int) bool {
		return scope.variables[i].name > scope.variables[j].name
	})
}
