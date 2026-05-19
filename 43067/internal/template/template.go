package template

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"snippetbox/internal/logger"
)

var varPattern = regexp.MustCompile(`\{\{(\w+)(?::([^}]+))?\}\}`)

type Variable struct {
	Name         string
	DefaultValue string
	Prompt       string
}

func ExtractVariables(content string) []Variable {
	matches := varPattern.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)
	var variables []Variable

	for _, match := range matches {
		name := match[1]
		if !seen[name] {
			seen[name] = true
			var defaultValue string
			if len(match) > 2 {
				defaultValue = match[2]
			}
			variables = append(variables, Variable{
				Name:         name,
				DefaultValue: defaultValue,
				Prompt:       fmt.Sprintf("Enter value for %s", name),
			})
		}
	}

	logger.Debug("Extracted %d variables from template", len(variables))
	return variables
}

func PromptForVariables(variables []Variable) (map[string]string, error) {
	values := make(map[string]string)
	scanner := bufio.NewScanner(os.Stdin)

	for _, v := range variables {
		if v.DefaultValue != "" {
			fmt.Printf("%s [%s]: ", v.Prompt, v.DefaultValue)
		} else {
			fmt.Printf("%s: ", v.Prompt)
		}

		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		if input == "" && v.DefaultValue != "" {
			values[v.Name] = v.DefaultValue
		} else {
			values[v.Name] = input
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	return values, nil
}

func ReplaceVariables(content string, values map[string]string) string {
	result := varPattern.ReplaceAllStringFunc(content, func(match string) string {
		groups := varPattern.FindStringSubmatch(match)
		if len(groups) < 2 {
			return match
		}
		name := groups[1]
		if val, ok := values[name]; ok {
			return val
		}
		return match
	})

	logger.Debug("Template variables replaced")
	return result
}

func ProcessTemplate(content string) (string, error) {
	variables := ExtractVariables(content)
	if len(variables) == 0 {
		return content, nil
	}

	fmt.Println("\nTemplate variables detected:")
	values, err := PromptForVariables(variables)
	if err != nil {
		return "", err
	}

	return ReplaceVariables(content, values), nil
}
