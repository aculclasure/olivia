package modules

import (
	"context"
	"fmt"
	"strings"

	"github.com/bndr/gojenkins"
	"github.com/olivia-ai/olivia/language"
	"github.com/olivia-ai/olivia/user"
	"github.com/olivia-ai/olivia/util"
)

var (
	// JenkinsSetterTag is the intent tag for its module.
	JenkinsCredsSetterTag = "jenkins credentials setter"

	// JenkinsURLSetterTag is the intent tag for its module.
	JenkinsURLSetterTag = "jenkins url setter"

	// JenkinsJobNamesGetterTag is the intent tag for its module.
	JenkinsJobNamesGetterTag = "jenkins job names getter"
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

func JenkinsJobNamesGetterReplacer(locale, entry, response, token string) (string, string) {
	info := user.GetUserInformation(token)
	if info.JenkinsURL == "" {
		return JenkinsJobNamesGetterTag, util.GetMessage(locale, "no jenkins url")
	}
	if info.JenkinsUser == "" || info.JenkinsPassword == "" {
		return JenkinsJobNamesGetterTag, util.GetMessage(locale, "no jenkins credentials")
	}

	ctx := context.Background()
	jenkins := gojenkins.CreateJenkins(nil, info.JenkinsURL, info.JenkinsUser, info.JenkinsPassword)
	if _, err := jenkins.Init(ctx); err != nil {
		return JenkinsJobNamesGetterTag, util.GetMessage(locale, "no jenkins connection")
	}

	jobs, err := jenkins.GetAllJobNames(ctx)
	if err != nil {
		return JenkinsJobNamesGetterTag, util.GetMessage(locale, "no jenkins connection")
	}
	var names []string
	for _, j := range jobs {
		names = append(names, j.Name)
	}

	return JenkinsJobNamesGetterTag, fmt.Sprintf(response, strings.Join(names, ","))
}
