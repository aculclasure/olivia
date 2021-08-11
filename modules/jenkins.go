package modules

import (
	"github.com/olivia-ai/olivia/language"
	"github.com/olivia-ai/olivia/user"
	"github.com/olivia-ai/olivia/util"
)

var (
	// JenkinsSetterTag is the intent tag for its module.
	JenkinsCredsSetterTag = "jenkins credentials setter"

	// JenkinsURLSetterTag is the intent tag for its module.
	JenkinsURLSetterTag = "jenkins url setter"
)

// JenkinsSetterReplacer extracts a Jenkins username and password from the
// given entry and stores them in the user information.
// See modules/modules.go#Module.Replacer() for more details.
func JenkinsCredsSetterReplacer(locale, entry, response, token string) (string, string) {
	creds := language.SearchJenkinsCredentials(locale, entry)
	if len(creds) != 2 {
		return JenkinsCredsSetterTag, util.GetMessage(locale, "no jenkins credentials")
	}

	user.ChangeUserInformation(token, func(info user.Information) user.Information {
		info.JenkinsUser = creds[0]
		info.JenkinsPassword = creds[1]
		return info
	})

	return JenkinsCredsSetterTag, response
}

// JenkinsURLSetterReplacer extracts a Jenkins API URL from the given entry and
// stores them in the user information.
// See modules/modules.go#Module.Replacer() for more details.
func JenkinsURLSetterReplacer(locale, entry, response, token string) (string, string) {
	apiURL := language.SearchJenkinsURL(entry)
	if apiURL == "" {
		return JenkinsURLSetterTag, util.GetMessage(locale, "no jenkins url")
	}

	user.ChangeUserInformation(token, func(info user.Information) user.Information {
		info.JenkinsURL = apiURL
		return info
	})

	return JenkinsURLSetterTag, response
}
