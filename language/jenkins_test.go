package language

import (
	"reflect"
	"testing"
)

func TestSearchJenkinsCredentials(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		locale string
		input  string
		want   []string
	}{
		"Credentials separated by colon can be parsed OK": {
			locale: "en",
			input:  "My Jenkins credentials are user:password",
			want:   []string{"user", "password"},
		},
		"Missing password in colon-separated credentials returns nil": {
			locale: "en",
			input:  "My Jenkins credentials are user:    .",
			want:   nil,
		},
		"Missing user in colon-separated credentials returns nil": {
			locale: "en",
			input:  "My Jenkins credentials are    :password.",
			want:   nil,
		},
		"Single word sentence with no credentials returns nil": {
			locale: "en",
			input:  "My",
			want:   nil,
		},
		"No Jenkins context in sentence returns nil": {
			locale: "en",
			input:  "user:password",
			want:   nil,
		},
		"Multi-word sentence can be parsed OK": {
			locale: "en",
			input:  "Here are my Jenkins credentials: user password!!!!!!!",
			want:   []string{"user", "password"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := SearchJenkinsCredentials(tc.locale, tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestSearchJenkinsURL(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input string
		want  string
	}{
		"Sentence with no Jenkins keyword returns empty string": {
			input: "the URL is http://something.com",
			want:  "",
		},
		"No URL in sentence returns empty string": {
			input: "The Jenkins URL is: ",
			want:  "",
		},
		"URL with no host returns empty string": {
			input: "The Jenkins URL is http://",
			want:  "",
		},
		"Invalid URL returns empty string": {
			input: "jenkins url is https wait where is the rest of the URL?",
			want:  "",
		},
		"Sentence with single URL returns URL": {
			input: "The JENKINS URL is https://10.0.0.2/jenkins/api ....ok?",
			want:  "https://10.0.0.2/jenkins/api",
		},
		"Sentence with multiple URLs returns leftmost match": {
			input: "The Jenkins URL is http://10.0.0.1 or https://10.0.0.2",
			want:  "http://10.0.0.1",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := SearchJenkinsURL(tc.input)
			if tc.want != got {
				t.Fatalf("want %s, got %s", tc.want, got)
			}
		})
	}
}
