package ref

import "strings"

// SearchCountries returns countries matching a case-insensitive query against
// code and name. Results are ordered: exact code match first, then name prefix,
// then name contains.
func (r *Registry) SearchCountries(query string) []*Country {
	if query == "" {
		return nil
	}
	q := strings.ToLower(query)

	var exact, prefix, contains []*Country
	for i := range countries {
		c := &countries[i]
		code := strings.ToLower(c.Code)
		name := strings.ToLower(c.Name)
		switch {
		case code == q:
			exact = append(exact, c)
		case strings.HasPrefix(name, q):
			prefix = append(prefix, c)
		case strings.Contains(name, q):
			contains = append(contains, c)
		}
	}
	result := make([]*Country, 0, len(exact)+len(prefix)+len(contains))
	result = append(result, exact...)
	result = append(result, prefix...)
	result = append(result, contains...)
	return result
}

// SearchLanguages returns languages matching a case-insensitive query
// against code, name, and native name.
func (r *Registry) SearchLanguages(query string) []*Language {
	if query == "" {
		return nil
	}
	q := strings.ToLower(query)

	var exact, prefix, contains []*Language
	for i := range languages {
		l := &languages[i]
		code := strings.ToLower(l.Code)
		name := strings.ToLower(l.Name)
		native := strings.ToLower(l.NativeName)
		switch {
		case code == q:
			exact = append(exact, l)
		case strings.HasPrefix(name, q) || strings.HasPrefix(native, q):
			prefix = append(prefix, l)
		case strings.Contains(name, q) || strings.Contains(native, q):
			contains = append(contains, l)
		}
	}
	result := make([]*Language, 0, len(exact)+len(prefix)+len(contains))
	result = append(result, exact...)
	result = append(result, prefix...)
	result = append(result, contains...)
	return result
}

// SearchTimezones returns timezones matching a case-insensitive query
// against IANA name and region.
func (r *Registry) SearchTimezones(query string) []*Timezone {
	if query == "" {
		return nil
	}
	q := strings.ToLower(query)

	var exact, contains []*Timezone
	for i := range timezones {
		tz := &timezones[i]
		name := strings.ToLower(tz.Name)
		switch {
		case name == q:
			exact = append(exact, tz)
		case strings.Contains(name, q):
			contains = append(contains, tz)
		}
	}
	result := make([]*Timezone, 0, len(exact)+len(contains))
	result = append(result, exact...)
	result = append(result, contains...)
	return result
}

// SearchCurrencies returns currencies matching a case-insensitive query
// against code, name, and symbol.
func (r *Registry) SearchCurrencies(query string) []*Currency {
	if query == "" {
		return nil
	}
	q := strings.ToLower(query)

	var exact, prefix, contains []*Currency
	for i := range currencies {
		c := &currencies[i]
		code := strings.ToLower(c.Code)
		name := strings.ToLower(c.Name)
		symbol := strings.ToLower(c.Symbol)
		switch {
		case code == q || symbol == q:
			exact = append(exact, c)
		case strings.HasPrefix(name, q) || strings.HasPrefix(code, q):
			prefix = append(prefix, c)
		case strings.Contains(name, q):
			contains = append(contains, c)
		}
	}
	result := make([]*Currency, 0, len(exact)+len(prefix)+len(contains))
	result = append(result, exact...)
	result = append(result, prefix...)
	result = append(result, contains...)
	return result
}
