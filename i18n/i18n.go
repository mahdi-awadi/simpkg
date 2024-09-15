package i18n

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-per/simpkg/parse"
	"github.com/go-per/simpkg/str"
)

// Errors
var (
	ErrInvalidTranslation       = errors.New("invalid translations")
	ErrTranslationAlreadyExists = errors.New("translation already exists")
	ErrNoTranslationFound       = errors.New("no translations found")
)

// II18N interface
type II18N interface {
	SetDefaultLocale(lang string)
	GetDefaultLocale() string
	SetFilesExtension(ext string)
	SetSupportedLocales(l []string)
	AddSupportedLocale(l string)
	RemoveSupportedLocale(l string)
	GetSupportedLocales() []string
	IsSupportedLocale(locale string) bool
	AddTranslation(t *Translation) error
	Translate(text string, args ...any) string
	TranslateInLocale(text, locale string, args ...any) string
	Load(s string) error
}

// I18N struct
type I18N struct {
	locale           string
	filesExt         string
	translations     map[string]*Translation
	supportedLocales []string
	locker           sync.RWMutex
}

// New I18N instance
func New() *I18N {
	return &I18N{
		locale:           "en",
		translations:     make(map[string]*Translation, 0),
		supportedLocales: []string{"en"},
		filesExt:         "json",
		locker:           sync.RWMutex{},
	}
}

// SetDefaultLocale set default locale
func (i *I18N) SetDefaultLocale(lang string) {
	i.locale = lang
	i.AddSupportedLocale(lang)
}

// GetDefaultLocale get default locale
func (i *I18N) GetDefaultLocale() string {
	return i.locale
}

// SetFilesExtension set files extension
func (i *I18N) SetFilesExtension(ext string) {
	i.filesExt = strings.TrimPrefix(ext, ".")
}

// SetSupportedLocales set supported locales
func (i *I18N) SetSupportedLocales(l []string) {
	i.supportedLocales = l
}

// AddSupportedLocale add new supported locale
func (i *I18N) AddSupportedLocale(l string) {
	i.supportedLocales = str.UniqueSlice(append(i.supportedLocales, l))
}

// RemoveSupportedLocale remove supported locale
func (i *I18N) RemoveSupportedLocale(l string) {
	for index, locale := range i.supportedLocales {
		if locale == l {
			i.supportedLocales = append(i.supportedLocales[:index], i.supportedLocales[index+1:]...)
			break
		}
	}
}

// GetSupportedLocales get supported locales
func (i *I18N) GetSupportedLocales() []string {
	return i.supportedLocales
}

// IsSupportedLocale check if locale is supported
func (i *I18N) IsSupportedLocale(locale string) bool {
	for _, l := range i.supportedLocales {
		if l == locale {
			return true
		}
	}

	return false
}

// AddTranslation add new translations
func (i *I18N) AddTranslation(t *Translation) error {
	if t == nil {
		return ErrInvalidTranslation
	}

	exists := i.getTranslation(t.locale)
	if exists != nil {
		return ErrTranslationAlreadyExists
	}

	i.locker.RLock()
	i.translations[t.locale] = t
	i.locker.RUnlock()

	return nil
}

// Translate translate key
func (i *I18N) Translate(text string, args ...any) string {
	return i.TranslateInLocale(text, i.locale, args...)
}

// TranslateInLocale translate key in locale
func (i *I18N) TranslateInLocale(text, locale string, args ...any) string {
	t := i.getTranslation(locale)
	if t == nil {
		return text
	}

	return t.Translate(text, args...)
}

// Load load translations
func (i *I18N) Load(s string) error {
	entries, err := os.ReadDir(s)
	if err != nil {
		return err
	}
	if entries == nil || len(entries) == 0 {
		return ErrNoTranslationFound
	}

	for _, entry := range entries {
		if !i.IsSupportedLocale(entry.Name()) {
			continue
		}

		// find translations files
		files, err := filepath.Glob(filepath.Join(s, entry.Name(), "*."+i.filesExt))
		if err != nil {
			return err
		}

		for _, file := range files {
			content, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			var translations map[string]string
			err = parse.Decode(content, &translations)
			if err != nil {
				return err
			}

			if translations != nil || len(translations) > 0 {
				translation := i.getTranslation(entry.Name())
				if translation == nil {
					translation = NewTranslation(entry.Name())
					_ = i.AddTranslation(translation)
				}

				// set translations
				translation.Set(translations)
			}
		}
	}

	return nil
}

// getTranslation get translations by locale
func (i *I18N) getTranslation(locale string) *Translation {
	t, ok := i.translations[locale]
	if !ok {
		return nil
	}

	return t
}
