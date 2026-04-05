// Package ref provides static reference data for ISO countries, currencies,
// languages, and timezones, along with locale-aware formatting for numbers,
// money, dates, percentages, IBANs, and phone numbers.
package ref

import "strings"

// Registry holds indexed reference data for fast lookups.
type Registry struct {
	countriesByCode     map[string]*Country
	currenciesByCode    map[string]*Currency
	languagesByCode     map[string]*Language
	timezonesByName     map[string]*Timezone
	countriesByPrefix   map[string][]*Country
	countriesByCurrency map[string][]*Country
	countriesByLanguage map[string][]*Country
	timezonesByCountry  map[string][]*Timezone
}

// New creates a Registry with all reference data indexed.
func New() *Registry {
	r := &Registry{
		countriesByCode:     make(map[string]*Country, len(countries)),
		currenciesByCode:    make(map[string]*Currency, len(currencies)),
		languagesByCode:     make(map[string]*Language, len(languages)),
		timezonesByName:     make(map[string]*Timezone, len(timezones)),
		countriesByPrefix:   make(map[string][]*Country),
		countriesByCurrency: make(map[string][]*Country),
		countriesByLanguage: make(map[string][]*Country),
		timezonesByCountry:  make(map[string][]*Timezone),
	}
	for i := range currencies {
		r.currenciesByCode[currencies[i].Code] = &currencies[i]
	}
	for i := range languages {
		r.languagesByCode[languages[i].Code] = &languages[i]
	}
	for i := range timezones {
		tz := &timezones[i]
		r.timezonesByName[tz.Name] = tz
		r.timezonesByCountry[tz.CountryCode] = append(r.timezonesByCountry[tz.CountryCode], tz)
	}
	for i := range countries {
		c := &countries[i]
		r.countriesByCode[c.Code] = c
		r.countriesByPrefix[c.PhonePrefix] = append(r.countriesByPrefix[c.PhonePrefix], c)
		r.countriesByCurrency[c.CurrencyCode] = append(r.countriesByCurrency[c.CurrencyCode], c)
		for _, lang := range c.Languages {
			r.countriesByLanguage[lang] = append(r.countriesByLanguage[lang], c)
		}
	}
	return r
}

// ---------------------------------------------------------------------------
// Countries
// ---------------------------------------------------------------------------

// Countries returns all countries.
func (r *Registry) Countries() []*Country {
	out := make([]*Country, len(countries))
	for i := range countries {
		out[i] = &countries[i]
	}
	return out
}

// Country looks up a country by ISO 3166-1 alpha-2 code.
func (r *Registry) Country(code string) (*Country, bool) {
	c, ok := r.countriesByCode[strings.ToUpper(code)]
	return c, ok
}

// CountriesForCurrency returns all countries using a currency.
func (r *Registry) CountriesForCurrency(currencyCode string) ([]*Country, bool) {
	cs, ok := r.countriesByCurrency[strings.ToUpper(currencyCode)]
	return cs, ok && len(cs) > 0
}

// CountriesForLanguage returns all countries where a language is spoken.
func (r *Registry) CountriesForLanguage(languageCode string) ([]*Country, bool) {
	cs, ok := r.countriesByLanguage[strings.ToLower(languageCode)]
	return cs, ok && len(cs) > 0
}

// ---------------------------------------------------------------------------
// Currencies
// ---------------------------------------------------------------------------

// Currencies returns all currencies.
func (r *Registry) Currencies() []*Currency {
	out := make([]*Currency, len(currencies))
	for i := range currencies {
		out[i] = &currencies[i]
	}
	return out
}

// Currency looks up a currency by ISO 4217 code.
func (r *Registry) Currency(code string) (*Currency, bool) {
	c, ok := r.currenciesByCode[strings.ToUpper(code)]
	return c, ok
}

// CurrencyForCountry returns the primary currency for a country.
func (r *Registry) CurrencyForCountry(countryCode string) (*Currency, bool) {
	c, ok := r.countriesByCode[strings.ToUpper(countryCode)]
	if !ok {
		return nil, false
	}
	return r.Currency(c.CurrencyCode)
}

// ---------------------------------------------------------------------------
// Phone prefixes
// ---------------------------------------------------------------------------

// PhonePrefix returns the international dialing prefix for a country.
func (r *Registry) PhonePrefix(countryCode string) (string, bool) {
	c, ok := r.countriesByCode[strings.ToUpper(countryCode)]
	if !ok {
		return "", false
	}
	return c.PhonePrefix, true
}

// CountryForPhonePrefix returns the primary country for a phone prefix.
// Some prefixes are shared (e.g., +1 for US/CA); this returns the first match.
// Use CountriesForPhonePrefix to get all matches.
func (r *Registry) CountryForPhonePrefix(prefix string) (*Country, bool) {
	cs, ok := r.countriesByPrefix[prefix]
	if !ok || len(cs) == 0 {
		return nil, false
	}
	return cs[0], true
}

// CountriesForPhonePrefix returns all countries sharing a phone prefix.
func (r *Registry) CountriesForPhonePrefix(prefix string) ([]*Country, bool) {
	cs, ok := r.countriesByPrefix[prefix]
	return cs, ok && len(cs) > 0
}

// ---------------------------------------------------------------------------
// Languages
// ---------------------------------------------------------------------------

// Languages returns all languages.
func (r *Registry) Languages() []*Language {
	out := make([]*Language, len(languages))
	for i := range languages {
		out[i] = &languages[i]
	}
	return out
}

// Language looks up a language by ISO 639-1 code.
func (r *Registry) Language(code string) (*Language, bool) {
	l, ok := r.languagesByCode[strings.ToLower(code)]
	return l, ok
}

// LanguagesForCountry returns the languages spoken in a country.
func (r *Registry) LanguagesForCountry(countryCode string) ([]*Language, bool) {
	c, ok := r.countriesByCode[strings.ToUpper(countryCode)]
	if !ok {
		return nil, false
	}
	out := make([]*Language, 0, len(c.Languages))
	for _, code := range c.Languages {
		if l, ok := r.languagesByCode[code]; ok {
			out = append(out, l)
		}
	}
	return out, len(out) > 0
}

// ---------------------------------------------------------------------------
// Timezones
// ---------------------------------------------------------------------------

// Timezones returns all timezones.
func (r *Registry) Timezones() []*Timezone {
	out := make([]*Timezone, len(timezones))
	for i := range timezones {
		out[i] = &timezones[i]
	}
	return out
}

// Timezone looks up a timezone by IANA name.
func (r *Registry) Timezone(name string) (*Timezone, bool) {
	tz, ok := r.timezonesByName[name]
	return tz, ok
}

// TimezonesForCountry returns all timezones for a country.
func (r *Registry) TimezonesForCountry(countryCode string) ([]*Timezone, bool) {
	tzs, ok := r.timezonesByCountry[strings.ToUpper(countryCode)]
	return tzs, ok && len(tzs) > 0
}

// DefaultTimezone returns the primary timezone for a country.
func (r *Registry) DefaultTimezone(countryCode string) (*Timezone, bool) {
	c, ok := r.countriesByCode[strings.ToUpper(countryCode)]
	if !ok {
		return nil, false
	}
	return r.Timezone(c.DefaultTimezone)
}
