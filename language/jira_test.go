package language

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetJiraUserCredentials(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description string
		input       string
		want        JiraUserCredentials
	}{
		{
			description: "sentence with valid credentials returns expected JiraUserCredentials",
			input:       "My JIRA credentials are user-name:api-token",
			want: JiraUserCredentials{
				UserName: "user-name",
				ApiKey:   "api-token",
			},
		},
		{
			description: "sentence with valid credentials following invalid words returns expected JiraUserCredentials",
			input:       "my ::: jira ::: credentials are user-name:api-token",
			want: JiraUserCredentials{
				UserName: "user-name",
				ApiKey:   "api-token",
			},
		},
		{
			description: "sentence with multiple credentials words returns JiraUserCredentials with first match",
			input:       "my jira credentials are user1:token1 and user2:token2",
			want: JiraUserCredentials{
				UserName: "user1",
				ApiKey:   "token1",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := GetJiraUserCredentials(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestGetJiraUserCredentials_ErrorCases(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description string
		input       string
	}{
		{
			description: `empty sentence returns error`,
			input:       "",
		},
		{
			description: `sentence missing the word "jira" returns error`,
			input:       "my credentials are user-name:api-token",
		},
		{
			description: `sentence missing the word "credentials" returns error`,
			input:       "my jira info is user-name:api-token",
		},
		{
			description: `sentence missing the word "jira" and the word "credentials" returns error`,
			input:       "my info is user-name:api-token",
		},
		{
			description: `sentence with no colon-separated credentials returns error`,
			input:       "my jira credentials are user-name|api-token",
		},
		{
			description: `partial sentence with no credentials returns error`,
			input:       "my jira credentials are",
		},
		{
			description: `sentence with missing user-name and api-key returns error`,
			input:       "my jira api credentials are :",
		},
		{
			description: `sentence missing user-name in credentials word returns error`,
			input:       "my jira credentials are :api-key",
		},
		{
			description: `sentence missing api-key in credentials word returns error`,
			input:       "my jira credentials are user-name:",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, err := GetJiraUserCredentials(tc.input)
			if err == nil {
				t.Error("expected an error but did not get one")
			}
		})
	}
}
