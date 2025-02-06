package modules

import (
	"fmt"

	"github.com/olivia-ai/olivia/user"
)

type Session interface {
	Get() ([]byte, error)
	Send([]byte) error
}

type ResponseMessage struct {
	Content     string           `json:"content"`
	Tag         string           `json:"tag"`
	Information user.Information `json:"information"`
}

// Module is a structure for dynamic intents with a Tag, some Patterns and Responses and
// a Replacer function to execute the dynamic changes.
type Module struct {
	Tag       string
	Patterns  []string
	Responses []string
	Replacer  func(string, string, string, string) (string, string)
	Run       func(msg, token string, session Session) error
	Context   string
}

var modules = map[string][]Module{}

// RegisterModule registers a module into the map
func RegisterModule(locale string, module Module) {
	modules[locale] = append(modules[locale], module)
}

// RegisterModules registers an array of modules into the map
func RegisterModules(locale string, _modules []Module) {
	modules[locale] = append(modules[locale], _modules...)
}

// GetModules returns all the registered modules
func GetModules(locale string) []Module {
	return modules[locale]
}

// GetModuleByTag returns a module found by the given tag and locale
func GetModuleByTag(tag, locale string) Module {
	for _, module := range modules[locale] {
		if tag != module.Tag {
			continue
		}

		return module
	}

	return Module{}
}

// ReplaceContent apply the Replacer of the matching module to the response and returns it
func ReplaceContent(locale, tag, entry, response, token string) (string, string) {
	for _, module := range modules[locale] {
		if module.Tag != tag {
			continue
		}

		return module.Replacer(locale, entry, response, token)
	}

	return tag, response
}

func RunModule(message, locale, tag, token string, session Session) error {
	for _, module := range modules[locale] {
		if module.Tag != tag {
			continue
		}
		return module.Run(message, token, session)
	}
	return fmt.Errorf("no module found to handle intent %s in locale %s", tag, locale)
}
