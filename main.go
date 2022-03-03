package main

import (
	"bufio"
	"fmt"
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
	root := Scope{}
	currentScope := &root

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 || strings.HasPrefix(line, JoinStrings(SLASH, SLASH)) || strings.HasPrefix(line, JoinStrings(SLASH, ASTERISK)) || strings.HasPrefix(line, ASTERISK) {
			continue
		}

		if strings.HasPrefix(line, ASPERAND) && strings.HasSuffix(line, SEMICOLON) {
			parts := strings.SplitN(line, COLON, 2)
			currentScope.AddVariable(Variable{name: strings.Trim(parts[0], ASPERAND), value: strings.Trim(parts[1], JoinStrings(SPACE, SEMICOLON))})

			continue
		}

		if strings.HasSuffix(line, SEMICOLON) {
			parts := strings.SplitN(line, COLON, 2)
			currentScope.AddRule(Rule{property: parts[0], value: strings.Trim(parts[1], JoinStrings(SPACE, SEMICOLON))})

			continue
		}

		if strings.HasSuffix(line, CURLY_BRACKET_OPEN) {
			scope := Scope{parent: currentScope, selector: strings.Trim(line, JoinStrings(SPACE, CURLY_BRACKET_OPEN))}

			currentScope.AddScope(&scope)
			currentScope = &scope

			continue
		}

		if line == CURLY_BRACKET_CLOSE {
			currentScope = currentScope.parent
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filename, err)
		return
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go root.Process(&waitGroup, "", make([]Variable, 0))

	waitGroup.Wait()

	result := root.String()
	fmt.Print(result)
	fmt.Print("\n")
}
