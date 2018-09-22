// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package locale // import "miniflux.app/locale"

// Translator manage supported locales.
type Translator struct {
	locales catalog
}

// AddLanguage loads a new language into the system.
func (t *Translator) AddLanguage(language, data string) (err error) {
	t.locales[language], err = parseCatalogMessages(data)
	return err
}

// GetLanguage returns the given language handler.
func (t *Translator) GetLanguage(language string) *Language {
	translations, found := t.locales[language]
	if !found {
		return &Language{language: language}
	}

	return &Language{language: language, translations: translations}
}

// NewTranslator creates a new Translator.
func NewTranslator() *Translator {
	return &Translator{locales: make(catalog)}
}
