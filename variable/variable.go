package variable

import "regexp"

var regexpName *regexp.Regexp = regexp.MustCompile(`@[\w-]+`)

type Variable struct {
	Name  string
	Value string
}

func ExtractVariables(value string) []string {
	return regexpName.FindAllString(value, -1)
}
