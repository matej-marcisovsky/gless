package main

import (
	"bufio"
	"fmt"
	"gless/chars"
	"gless/rule"
	"gless/scope"
	"gless/utils"
	"gless/variable"
	"os"
	"strings"
	"sync"
)

func main() {
	args := os.Args[1:]
	filename := args[0]

	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)

	if err != nil {
		fmt.Printf("Error opening file '%s': %v\n", filename, err)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	root := scope.Scope{}
	currentScope := &root

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 || strings.HasPrefix(line, utils.JoinStrings(chars.SLASH, chars.SLASH)) || strings.HasPrefix(line, utils.JoinStrings(chars.SLASH, chars.ASTERISK)) || strings.HasPrefix(line, chars.ASTERISK) {
			continue
		}

		if strings.HasPrefix(line, chars.ASPERAND) && strings.HasSuffix(line, chars.SEMICOLON) {
			parts := strings.SplitN(line, chars.COLON, 2)
			currentScope.AddVariable(variable.Variable{Name: strings.Trim(parts[0], chars.ASPERAND), Value: strings.Trim(parts[1], utils.JoinStrings(chars.SPACE, chars.SEMICOLON))})

			continue
		}

		if strings.HasSuffix(line, chars.SEMICOLON) {
			parts := strings.SplitN(line, chars.COLON, 2)
			currentScope.AddRule(rule.Rule{Property: parts[0], Value: strings.Trim(parts[1], utils.JoinStrings(chars.SPACE, chars.SEMICOLON))})

			continue
		}

		if strings.HasSuffix(line, chars.CURLY_BRACKET_OPEN) {
			scope := scope.Scope{Parent: currentScope, Selector: strings.Trim(line, utils.JoinStrings(chars.SPACE, chars.CURLY_BRACKET_OPEN))}

			currentScope.AddScope(&scope)
			currentScope = &scope

			continue
		}

		if line == chars.CURLY_BRACKET_CLOSE {
			currentScope = currentScope.Parent
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filename, err)
		return
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go root.Process(&waitGroup, "", make([]variable.Variable, 0))

	waitGroup.Wait()

	result := root.String()
	fmt.Print(result)
}
