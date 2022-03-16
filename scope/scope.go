package scope

import (
	"gless/chars"
	"gless/rule"
	"gless/utils"
	"gless/variable"
	"regexp"
	"sort"
	"strings"
	"sync"
)

var regexpCommaSpace *regexp.Regexp = regexp.MustCompile(utils.JoinStrings(chars.COMMA, "\\s{0,1}"))

type Scope struct {
	Parent    *Scope
	Result    string
	Rules     []rule.Rule
	Scopes    []*Scope
	Selector  string
	Variables []variable.Variable
}

func (scope *Scope) AddVariable(variable variable.Variable) {
	scope.Variables = append(scope.Variables, variable)
}

func (scope *Scope) AddRule(rule rule.Rule) {
	scope.Rules = append(scope.Rules, rule)
}

func (scope *Scope) AddScope(_scope *Scope) {
	scope.Scopes = append(scope.Scopes, _scope)
}

func (scope *Scope) Process(waitGroup *sync.WaitGroup, parentSelector string, variables []variable.Variable) {
	scope.prepareVariables(variables)
	selector := scope.computeSelector(parentSelector)

	if scope.Selector != "" && len(scope.Rules) > 0 {
		scope.Result = utils.JoinStrings(selector, chars.SPACE, chars.CURLY_BRACKET_OPEN, chars.NEW_LINE)
	}

	for _, rule := range scope.Rules {
		rule.Process(scope.Variables)
		scope.Result = scope.Result + utils.JoinStrings(chars.TABULATOR, rule.Property, chars.COLON, chars.SPACE, rule.Value, chars.SEMICOLON, chars.NEW_LINE)
	}

	for _, _scope := range scope.Scopes {
		waitGroup.Add(1)
		go _scope.Process(waitGroup, selector, scope.Variables)
	}

	if scope.Selector != "" && len(scope.Rules) > 0 {
		scope.Result = scope.Result + utils.JoinStrings(chars.CURLY_BRACKET_CLOSE, chars.NEW_LINE)
	}

	defer waitGroup.Done()
}

func (scope *Scope) String() string {
	result := scope.Result

	for _, _scope := range scope.Scopes {
		result = result + _scope.String()
	}

	return result
}

func (scope *Scope) computeSelector(parentSelector string) string {
	result := strings.ReplaceAll(scope.Selector, chars.AMPERSAND, parentSelector)

	if strings.Contains(result, chars.COMMA) {
		result = regexpCommaSpace.ReplaceAllString(result, utils.JoinStrings(chars.COMMA, chars.NEW_LINE))
	}

	if !strings.HasPrefix(scope.Selector, chars.AMPERSAND) && len(parentSelector) > 0 {
		result = utils.JoinStrings(parentSelector, chars.SPACE, result)
	}

	return result
}

func (scope *Scope) prepareVariables(variables []variable.Variable) {
	scopeVariablesNames := make(map[string]struct{}, len(scope.Variables))

	for _, variable := range scope.Variables {
		scopeVariablesNames[variable.Name] = struct{}{}
	}

	for _, variable := range variables {
		if _, ok := scopeVariablesNames[variable.Name]; !ok {
			scope.Variables = append(scope.Variables, variable)
		}
	}

	sort.Slice(scope.Variables, func(i, j int) bool {
		return scope.Variables[i].Name > scope.Variables[j].Name
	})
}
