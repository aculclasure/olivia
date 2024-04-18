package language

import (
	"fmt"
	"strings"
)

type JiraUserCredentials struct {
	UserName string
	ApiKey   string
}

func GetJiraUserCredentials(sentence string) (JiraUserCredentials, error) {
	if !LevenshteinContains(strings.ToLower(sentence), "jira", 1) {
		return JiraUserCredentials{}, fmt.Errorf(`sentence must contain a word close to "jira", got %q`, sentence)
	}
	if !LevenshteinContains(strings.ToLower(sentence), "credentials", 3) {
		return JiraUserCredentials{}, fmt.Errorf(`sentence must contain a word close to "credentials", got: %q`, sentence)
	}
	for _, word := range strings.Fields(sentence) {
		if strings.Count(word, ":") != 1 {
			continue
		}
		fields := strings.Split(word, ":")
		if len(fields[0]) == 0 || len(fields[1]) == 0 {
			continue
		}
		return JiraUserCredentials{
			UserName: fields[0],
			ApiKey:   fields[1],
		}, nil
	}
	return JiraUserCredentials{}, fmt.Errorf("sentence must contain jira credentials in the form <user-name>:<api-key>, got %q", sentence)
}
