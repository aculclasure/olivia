package modules

import (
	"github.com/olivia-ai/olivia/language"
	"github.com/olivia-ai/olivia/modules/smartjira"
	"github.com/olivia-ai/olivia/user"
	"github.com/olivia-ai/olivia/util"
)

var (
	JiraUserCredentialsSetterTag = "jira user credentials setter"
)

func JiraUserCredentialsSetterReplacer(locale, entry, _, token string) (string, string) {
	creds, err := language.GetJiraUserCredentials(entry)
	if err != nil {
		return JiraUserCredentialsSetterTag, util.GetMessage(locale, "jira credentials")
	}
	client, err := smartjira.NewClient(creds.UserName, creds.ApiKey)
	if err != nil {
		return JiraUserCredentialsSetterTag, util.GetMessage(locale, "jira client")
	}
	jiraUserID, err := client.CurrentUserId()
	if err != nil {
		return JiraUserCredentialsSetterTag, util.GetMessage(locale, "jira api request")
	}
	user.ChangeUserInformation(token, func(i user.Information) user.Information {
		i.JiraApiToken = creds.ApiKey
		i.JiraUserName = creds.UserName
		i.JiraUserID = jiraUserID
		return i
	})
	return JiraUserCredentialsSetterTag, "Great news! I've saved your Jira credentials and confirmed that I can communicate to Jira with them."
}
