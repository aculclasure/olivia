package smartjira

import (
	"errors"
	"testing"

	"github.com/andygrunwald/go-jira"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewClient_ErrorCases(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		userName string
		apiKey   string
	}{
		"Returns error given empty user name": {
			userName: "",
			apiKey:   "apikey",
		},
		"Returns error given empty api key": {
			userName: "user",
			apiKey:   "",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := NewClient(tc.userName, tc.apiKey)
			if err == nil {
				t.Error("expected an error but did not get one")
			}
		})
	}
}

func TestNewClient_ReturnsValidClientGivenValidArgs(t *testing.T) {
	t.Parallel()
	got, err := NewClient("user", "apikey")
	if err != nil {
		t.Fatal(err)
	}
	want := &Client{
		UserName: "user",
		APIKey:   "apikey",
		BaseURL:  "https://inforwiki.atlassian.net",
	}
	if !cmp.Equal(want, got, cmpopts.IgnoreUnexported(Client{})) {
		t.Error(cmp.Diff(want, got, cmpopts.IgnoreUnexported(Client{})))
	}
}

func TestClient_CurrentUserIdReturnsExpectedUserIdGivenValidUser(t *testing.T) {
	t.Parallel()
	client := &Client{
		currentUserGetter: func() (*jira.User, *jira.Response, error) {
			return &jira.User{
				AccountID: "123456",
			}, nil, nil
		},
	}
	want := "123456"
	got, err := client.CurrentUserId()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestClient_CurrentUserIdReturnsErrorWhenJiraClientReturnsError(t *testing.T) {
	t.Parallel()
	client := &Client{
		currentUserGetter: func() (*jira.User, *jira.Response, error) {
			return nil, nil, errors.New("oh no!")
		},
	}
	_, err := client.CurrentUserId()
	if err == nil {
		t.Error("expected an error but did not get one")
	}
}
