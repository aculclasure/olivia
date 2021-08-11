package language

import (
	"net/url"
	"regexp"
	"strings"
)

var jenkinsCredsPatterns = map[string][]*regexp.Regexp{
	"en": {
		regexp.MustCompile(`^(?i).*Jenkins credentials are:?\s*(?P<name>\w+)(?::|\s+)(?P<password>\w+)[[:punct:]]*$`),
		regexp.MustCompile(`^(?i).*my Jenkins credentials:?\s*(?P<name>\w+)(?::|\s+)(?P<password>\w+)[[:punct:]]*$`),
	},
}

// SearchJenkinsCredentials searches for the 2 Jenkins credentials (user name
// and password/API key) in the given sentence and returns them as a slice of
// strings.
func SearchJenkinsCredentials(locale, sentence string) []string {
	for _, re := range jenkinsCredsPatterns[locale] {
		match := re.FindStringSubmatch(sentence)
		if match == nil {
			continue
		}

		result := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		wantSubexpNames := []string{"name", "password"}
		var values []string
		for _, subexp := range wantSubexpNames {
			v, ok := result[subexp]
			if !ok {
				return nil
			}
			values = append(values, v)
		}
		return values
	}

	return nil
}

// SearchJenkinsURL searches for the Jenkins API URL in the given sentence
// and returns it.
func SearchJenkinsURL(sentence string) string {
	if !strings.Contains(strings.ToLower(sentence), "jenkins") {
		return ""
	}

	var match string
	fields := strings.Fields(sentence)
	for _, fld := range fields {
		if strings.HasPrefix(strings.ToLower(fld), "http") {
			match = fld
			break
		}
	}

	u, err := url.Parse(match)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}

	return match
}
