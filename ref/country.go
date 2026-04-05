package ref

// DateOrder represents the order of date components in locale formatting.
type DateOrder int

const (
	DateOrderDMY DateOrder = iota // 22/03/1985 (most of world)
	DateOrderMDY                  // 03/22/1985 (US)
	DateOrderYMD                  // 1985/03/22 (Japan, Korea, China)
)

// Country holds ISO 3166-1 alpha-2 country reference data.
type Country struct {
	Code            string    // ISO 3166-1 alpha-2: "BE"
	Name            string    // English name: "Belgium"
	PhonePrefix     string    // International dialing: "+32"
	CurrencyCode    string    // Primary ISO 4217 currency: "EUR"
	Languages       []string  // ISO 639-1 language codes: ["nl", "fr", "de"]
	DefaultTimezone string    // Primary IANA timezone: "Europe/Brussels"
	DecimalSep      string    // Number decimal separator: ","
	ThousandsSep    string    // Number thousands separator: "."
	SymbolAfter     bool      // Currency symbol placed after amount
	DateOrder       DateOrder // Date component order
	DateSep         string    // Date separator: "/"
}

var countries = []Country{
	// EU / EEA — Eurozone
	{Code: "AT", Name: "Austria", PhonePrefix: "+43", CurrencyCode: "EUR", Languages: []string{"de"}, DefaultTimezone: "Europe/Vienna", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "BE", Name: "Belgium", PhonePrefix: "+32", CurrencyCode: "EUR", Languages: []string{"nl", "fr", "de"}, DefaultTimezone: "Europe/Brussels", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "CY", Name: "Cyprus", PhonePrefix: "+357", CurrencyCode: "EUR", Languages: []string{"el", "tr"}, DefaultTimezone: "Europe/Nicosia", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "DE", Name: "Germany", PhonePrefix: "+49", CurrencyCode: "EUR", Languages: []string{"de"}, DefaultTimezone: "Europe/Berlin", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "EE", Name: "Estonia", PhonePrefix: "+372", CurrencyCode: "EUR", Languages: []string{"et"}, DefaultTimezone: "Europe/Tallinn", DecimalSep: ",", ThousandsSep: " ", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "ES", Name: "Spain", PhonePrefix: "+34", CurrencyCode: "EUR", Languages: []string{"es"}, DefaultTimezone: "Europe/Madrid", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "FI", Name: "Finland", PhonePrefix: "+358", CurrencyCode: "EUR", Languages: []string{"fi", "sv"}, DefaultTimezone: "Europe/Helsinki", DecimalSep: ",", ThousandsSep: " ", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "FR", Name: "France", PhonePrefix: "+33", CurrencyCode: "EUR", Languages: []string{"fr"}, DefaultTimezone: "Europe/Paris", DecimalSep: ",", ThousandsSep: " ", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "GR", Name: "Greece", PhonePrefix: "+30", CurrencyCode: "EUR", Languages: []string{"el"}, DefaultTimezone: "Europe/Athens", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "HR", Name: "Croatia", PhonePrefix: "+385", CurrencyCode: "EUR", Languages: []string{"hr"}, DefaultTimezone: "Europe/Zagreb", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "IE", Name: "Ireland", PhonePrefix: "+353", CurrencyCode: "EUR", Languages: []string{"en"}, DefaultTimezone: "Europe/Dublin", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "IT", Name: "Italy", PhonePrefix: "+39", CurrencyCode: "EUR", Languages: []string{"it"}, DefaultTimezone: "Europe/Rome", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "LT", Name: "Lithuania", PhonePrefix: "+370", CurrencyCode: "EUR", Languages: []string{"lt"}, DefaultTimezone: "Europe/Vilnius", DecimalSep: ",", ThousandsSep: " ", SymbolAfter: true, DateOrder: DateOrderYMD, DateSep: "-"},
	{Code: "LU", Name: "Luxembourg", PhonePrefix: "+352", CurrencyCode: "EUR", Languages: []string{"fr", "de"}, DefaultTimezone: "Europe/Luxembourg", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "LV", Name: "Latvia", PhonePrefix: "+371", CurrencyCode: "EUR", Languages: []string{"lv"}, DefaultTimezone: "Europe/Riga", DecimalSep: ",", ThousandsSep: " ", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "MT", Name: "Malta", PhonePrefix: "+356", CurrencyCode: "EUR", Languages: []string{"mt", "en"}, DefaultTimezone: "Europe/Malta", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "NL", Name: "Netherlands", PhonePrefix: "+31", CurrencyCode: "EUR", Languages: []string{"nl"}, DefaultTimezone: "Europe/Amsterdam", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "-"},
	{Code: "PT", Name: "Portugal", PhonePrefix: "+351", CurrencyCode: "EUR", Languages: []string{"pt"}, DefaultTimezone: "Europe/Lisbon", DecimalSep: ",", ThousandsSep: " ", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "SI", Name: "Slovenia", PhonePrefix: "+386", CurrencyCode: "EUR", Languages: []string{"sl"}, DefaultTimezone: "Europe/Ljubljana", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "SK", Name: "Slovakia", PhonePrefix: "+421", CurrencyCode: "EUR", Languages: []string{"sk"}, DefaultTimezone: "Europe/Bratislava", DecimalSep: ",", ThousandsSep: " ", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},

	// EU — Non-Eurozone
	{Code: "BG", Name: "Bulgaria", PhonePrefix: "+359", CurrencyCode: "BGN", Languages: []string{"bg"}, DefaultTimezone: "Europe/Sofia", DecimalSep: ",", ThousandsSep: " ", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "CZ", Name: "Czech Republic", PhonePrefix: "+420", CurrencyCode: "CZK", Languages: []string{"cs"}, DefaultTimezone: "Europe/Prague", DecimalSep: ",", ThousandsSep: " ", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "DK", Name: "Denmark", PhonePrefix: "+45", CurrencyCode: "DKK", Languages: []string{"da"}, DefaultTimezone: "Europe/Copenhagen", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "HU", Name: "Hungary", PhonePrefix: "+36", CurrencyCode: "HUF", Languages: []string{"hu"}, DefaultTimezone: "Europe/Budapest", DecimalSep: ",", ThousandsSep: " ", SymbolAfter: true, DateOrder: DateOrderYMD, DateSep: "."},
	{Code: "PL", Name: "Poland", PhonePrefix: "+48", CurrencyCode: "PLN", Languages: []string{"pl"}, DefaultTimezone: "Europe/Warsaw", DecimalSep: ",", ThousandsSep: " ", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "RO", Name: "Romania", PhonePrefix: "+40", CurrencyCode: "RON", Languages: []string{"ro"}, DefaultTimezone: "Europe/Bucharest", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "SE", Name: "Sweden", PhonePrefix: "+46", CurrencyCode: "SEK", Languages: []string{"sv"}, DefaultTimezone: "Europe/Stockholm", DecimalSep: ",", ThousandsSep: " ", SymbolAfter: true, DateOrder: DateOrderYMD, DateSep: "-"},

	// EEA / EFTA
	{Code: "CH", Name: "Switzerland", PhonePrefix: "+41", CurrencyCode: "CHF", Languages: []string{"de", "fr", "it"}, DefaultTimezone: "Europe/Zurich", DecimalSep: ".", ThousandsSep: "'", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "IS", Name: "Iceland", PhonePrefix: "+354", CurrencyCode: "ISK", Languages: []string{"is"}, DefaultTimezone: "Atlantic/Reykjavik", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "LI", Name: "Liechtenstein", PhonePrefix: "+423", CurrencyCode: "CHF", Languages: []string{"de"}, DefaultTimezone: "Europe/Vaduz", DecimalSep: ".", ThousandsSep: "'", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
	{Code: "NO", Name: "Norway", PhonePrefix: "+47", CurrencyCode: "NOK", Languages: []string{"no"}, DefaultTimezone: "Europe/Oslo", DecimalSep: ",", ThousandsSep: " ", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},

	// UK
	{Code: "GB", Name: "United Kingdom", PhonePrefix: "+44", CurrencyCode: "GBP", Languages: []string{"en"}, DefaultTimezone: "Europe/London", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},

	// Americas
	{Code: "AR", Name: "Argentina", PhonePrefix: "+54", CurrencyCode: "ARS", Languages: []string{"es"}, DefaultTimezone: "America/Argentina/Buenos_Aires", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "BR", Name: "Brazil", PhonePrefix: "+55", CurrencyCode: "BRL", Languages: []string{"pt"}, DefaultTimezone: "America/Sao_Paulo", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "CA", Name: "Canada", PhonePrefix: "+1", CurrencyCode: "CAD", Languages: []string{"en", "fr"}, DefaultTimezone: "America/Toronto", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderYMD, DateSep: "-"},
	{Code: "CO", Name: "Colombia", PhonePrefix: "+57", CurrencyCode: "COP", Languages: []string{"es"}, DefaultTimezone: "America/Bogota", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "MX", Name: "Mexico", PhonePrefix: "+52", CurrencyCode: "MXN", Languages: []string{"es"}, DefaultTimezone: "America/Mexico_City", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "US", Name: "United States", PhonePrefix: "+1", CurrencyCode: "USD", Languages: []string{"en", "es"}, DefaultTimezone: "America/New_York", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderMDY, DateSep: "/"},

	// Middle East
	{Code: "AE", Name: "United Arab Emirates", PhonePrefix: "+971", CurrencyCode: "AED", Languages: []string{"ar", "en"}, DefaultTimezone: "Asia/Dubai", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "BH", Name: "Bahrain", PhonePrefix: "+973", CurrencyCode: "BHD", Languages: []string{"ar"}, DefaultTimezone: "Asia/Bahrain", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "KW", Name: "Kuwait", PhonePrefix: "+965", CurrencyCode: "KWD", Languages: []string{"ar"}, DefaultTimezone: "Asia/Kuwait", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "OM", Name: "Oman", PhonePrefix: "+968", CurrencyCode: "OMR", Languages: []string{"ar"}, DefaultTimezone: "Asia/Muscat", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "QA", Name: "Qatar", PhonePrefix: "+974", CurrencyCode: "QAR", Languages: []string{"ar"}, DefaultTimezone: "Asia/Qatar", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "SA", Name: "Saudi Arabia", PhonePrefix: "+966", CurrencyCode: "SAR", Languages: []string{"ar"}, DefaultTimezone: "Asia/Riyadh", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},

	// Asia-Pacific
	{Code: "AU", Name: "Australia", PhonePrefix: "+61", CurrencyCode: "AUD", Languages: []string{"en"}, DefaultTimezone: "Australia/Sydney", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "CN", Name: "China", PhonePrefix: "+86", CurrencyCode: "CNY", Languages: []string{"zh"}, DefaultTimezone: "Asia/Shanghai", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderYMD, DateSep: "/"},
	{Code: "HK", Name: "Hong Kong", PhonePrefix: "+852", CurrencyCode: "HKD", Languages: []string{"zh", "en"}, DefaultTimezone: "Asia/Hong_Kong", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "IN", Name: "India", PhonePrefix: "+91", CurrencyCode: "INR", Languages: []string{"hi", "en"}, DefaultTimezone: "Asia/Kolkata", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "JP", Name: "Japan", PhonePrefix: "+81", CurrencyCode: "JPY", Languages: []string{"ja"}, DefaultTimezone: "Asia/Tokyo", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderYMD, DateSep: "/"},
	{Code: "KR", Name: "South Korea", PhonePrefix: "+82", CurrencyCode: "KRW", Languages: []string{"ko"}, DefaultTimezone: "Asia/Seoul", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderYMD, DateSep: "."},
	{Code: "NZ", Name: "New Zealand", PhonePrefix: "+64", CurrencyCode: "NZD", Languages: []string{"en"}, DefaultTimezone: "Pacific/Auckland", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "SG", Name: "Singapore", PhonePrefix: "+65", CurrencyCode: "SGD", Languages: []string{"en", "zh", "ms"}, DefaultTimezone: "Asia/Singapore", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},

	// Africa
	{Code: "EG", Name: "Egypt", PhonePrefix: "+20", CurrencyCode: "EGP", Languages: []string{"ar"}, DefaultTimezone: "Africa/Cairo", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "KE", Name: "Kenya", PhonePrefix: "+254", CurrencyCode: "KES", Languages: []string{"sw", "en"}, DefaultTimezone: "Africa/Nairobi", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "MA", Name: "Morocco", PhonePrefix: "+212", CurrencyCode: "MAD", Languages: []string{"ar", "fr"}, DefaultTimezone: "Africa/Casablanca", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "NG", Name: "Nigeria", PhonePrefix: "+234", CurrencyCode: "NGN", Languages: []string{"en"}, DefaultTimezone: "Africa/Lagos", DecimalSep: ".", ThousandsSep: ",", SymbolAfter: false, DateOrder: DateOrderDMY, DateSep: "/"},
	{Code: "ZA", Name: "South Africa", PhonePrefix: "+27", CurrencyCode: "ZAR", Languages: []string{"en"}, DefaultTimezone: "Africa/Johannesburg", DecimalSep: ".", ThousandsSep: " ", SymbolAfter: false, DateOrder: DateOrderYMD, DateSep: "/"},

	// Turkey
	{Code: "TR", Name: "Turkey", PhonePrefix: "+90", CurrencyCode: "TRY", Languages: []string{"tr"}, DefaultTimezone: "Europe/Istanbul", DecimalSep: ",", ThousandsSep: ".", SymbolAfter: true, DateOrder: DateOrderDMY, DateSep: "."},
}
