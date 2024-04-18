package smartjira

import (
	"errors"
	"strings"

	jira "github.com/andygrunwald/go-jira"
)

type Client struct {
	UserName          string
	APIKey            string
	BaseURL           string
	jc                *jira.Client
	currentUserGetter func() (*jira.User, *jira.Response, error)
}

func NewClient(userName, apiKey string) (*Client, error) {
	if strings.TrimSpace(userName) == "" {
		return nil, errors.New("username argument must be non-empty")
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("apikey argument must be non-empty")
	}
	tp := jira.BasicAuthTransport{
		Username: userName,
		Password: apiKey,
	}
	const baseURL = "https://inforwiki.atlassian.net"
	jc, err := jira.NewClient(tp.Client(), baseURL)
	if err != nil {
		return nil, err
	}
	c := &Client{
		UserName: userName,
		APIKey:   apiKey,
		BaseURL:  baseURL,
		jc:       jc,
		currentUserGetter: func() (*jira.User, *jira.Response, error) {
			return jc.User.GetSelf()
		},
	}
	return c, nil
}

func (c *Client) CurrentUserId() (string, error) {
	user, _, err := c.currentUserGetter()
	if err != nil {
		return "", err
	}
	return user.AccountID, nil
}
