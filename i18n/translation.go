package i18n

import (
	"regexp"
	"strings"
	"sync"

	"github.com/go-per/simpkg/format"
)

// compile regex for arguments
var re = regexp.MustCompile(`{(.*?)}+`)

// Translation struct
type Translation struct {
	Translations map[string]string
	locale       string
	locker       sync.RWMutex
	argPrefix    string
	argSuffix    string
}

// NewTranslation instance
func NewTranslation(locale string) *Translation {
	return &Translation{
		Translations: make(map[string]string, 0),
		locale:       locale,
		locker:       sync.RWMutex{},
		argPrefix:    "{",
		argSuffix:    "}",
	}
}

// Add new translations item
func (t *Translation) Add(key, value string) {
	t.locker.Lock()
	t.Translations[key] = value
	t.locker.Unlock()
}

// Set translations items
func (t *Translation) Set(m map[string]string) {
	if m == nil || len(m) == 0 {
		return
	}

	t.locker.Lock()
	for key, value := range m {
		t.Translations[key] = value
	}
	t.locker.Unlock()
}

// Get translations item string
func (t *Translation) Get(key string) (string, bool) {
	v, ok := t.Translations[key]
	if !ok {
		return key, false
	}

	return v, true
}

// Translate text
func (t *Translation) Translate(text string, args ...any) string {
	translate, ok := t.Get(text)
	if !ok {
		return text
	}

	// format string
	if args != nil && len(args) > 0 && args[0] != nil {
		translated := translate
		matches := re.FindAllStringSubmatch(translated, -1)
		if matches != nil && len(matches) > 0 {
			for index, match := range matches {
				if index < len(args) {
					replace := args[index]
					if replace != nil && match != nil && len(match) > 0 {
						translated = strings.ReplaceAll(translated, match[0], format.String(replace))
					}
				}
			}
		}

		return translated
	}

	return translate
}
