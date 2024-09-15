package i18n

import "fmt"

// Instance is default instance
var Instance II18N

// init func
func init() {
	Instance = New()
}

// Translate is a shortcut for Instance.Translate
func Translate(text string, args ...any) string {
	return Instance.Translate(text, args...)
}

// TranslateInLocale is a shortcut for Instance.TranslateInLocale
func TranslateInLocale(text, locale string, args ...any) string {
	return Instance.TranslateInLocale(text, locale, args...)
}

// TranslateAsError is a shortcut for Instance.TranslateAsError
func TranslateAsError(text string, args ...any) error {
	return fmt.Errorf(Instance.Translate(text, args...))
}

// TranslateAsErrorInLocale is a shortcut for Instance.TranslateAsErrorInLocale
func TranslateAsErrorInLocale(text, locale string, args ...any) error {
	return fmt.Errorf(Instance.TranslateInLocale(text, locale, args...))
}
