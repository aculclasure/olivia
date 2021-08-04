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
