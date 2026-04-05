package ref

// Timezone holds IANA timezone reference data.
type Timezone struct {
	Name        string // IANA name: "Europe/Brussels"
	UTCOffset   string // Standard UTC offset: "+01:00"
	Region      string // Region: "Europe"
	CountryCode string // Primary country: "BE"
}

var timezones = []Timezone{
	// Europe
	{Name: "Europe/London", UTCOffset: "+00:00", Region: "Europe", CountryCode: "GB"},
	{Name: "Europe/Dublin", UTCOffset: "+00:00", Region: "Europe", CountryCode: "IE"},
	{Name: "Atlantic/Reykjavik", UTCOffset: "+00:00", Region: "Atlantic", CountryCode: "IS"},
	{Name: "Europe/Lisbon", UTCOffset: "+00:00", Region: "Europe", CountryCode: "PT"},
	{Name: "Europe/Brussels", UTCOffset: "+01:00", Region: "Europe", CountryCode: "BE"},
	{Name: "Europe/Amsterdam", UTCOffset: "+01:00", Region: "Europe", CountryCode: "NL"},
	{Name: "Europe/Luxembourg", UTCOffset: "+01:00", Region: "Europe", CountryCode: "LU"},
	{Name: "Europe/Paris", UTCOffset: "+01:00", Region: "Europe", CountryCode: "FR"},
	{Name: "Europe/Berlin", UTCOffset: "+01:00", Region: "Europe", CountryCode: "DE"},
	{Name: "Europe/Vienna", UTCOffset: "+01:00", Region: "Europe", CountryCode: "AT"},
	{Name: "Europe/Zurich", UTCOffset: "+01:00", Region: "Europe", CountryCode: "CH"},
	{Name: "Europe/Vaduz", UTCOffset: "+01:00", Region: "Europe", CountryCode: "LI"},
	{Name: "Europe/Rome", UTCOffset: "+01:00", Region: "Europe", CountryCode: "IT"},
	{Name: "Europe/Madrid", UTCOffset: "+01:00", Region: "Europe", CountryCode: "ES"},
	{Name: "Europe/Malta", UTCOffset: "+01:00", Region: "Europe", CountryCode: "MT"},
	{Name: "Europe/Copenhagen", UTCOffset: "+01:00", Region: "Europe", CountryCode: "DK"},
	{Name: "Europe/Stockholm", UTCOffset: "+01:00", Region: "Europe", CountryCode: "SE"},
	{Name: "Europe/Oslo", UTCOffset: "+01:00", Region: "Europe", CountryCode: "NO"},
	{Name: "Europe/Warsaw", UTCOffset: "+01:00", Region: "Europe", CountryCode: "PL"},
	{Name: "Europe/Prague", UTCOffset: "+01:00", Region: "Europe", CountryCode: "CZ"},
	{Name: "Europe/Budapest", UTCOffset: "+01:00", Region: "Europe", CountryCode: "HU"},
	{Name: "Europe/Ljubljana", UTCOffset: "+01:00", Region: "Europe", CountryCode: "SI"},
	{Name: "Europe/Zagreb", UTCOffset: "+01:00", Region: "Europe", CountryCode: "HR"},
	{Name: "Europe/Bratislava", UTCOffset: "+01:00", Region: "Europe", CountryCode: "SK"},
	{Name: "Europe/Helsinki", UTCOffset: "+02:00", Region: "Europe", CountryCode: "FI"},
	{Name: "Europe/Tallinn", UTCOffset: "+02:00", Region: "Europe", CountryCode: "EE"},
	{Name: "Europe/Riga", UTCOffset: "+02:00", Region: "Europe", CountryCode: "LV"},
	{Name: "Europe/Vilnius", UTCOffset: "+02:00", Region: "Europe", CountryCode: "LT"},
	{Name: "Europe/Athens", UTCOffset: "+02:00", Region: "Europe", CountryCode: "GR"},
	{Name: "Europe/Bucharest", UTCOffset: "+02:00", Region: "Europe", CountryCode: "RO"},
	{Name: "Europe/Sofia", UTCOffset: "+02:00", Region: "Europe", CountryCode: "BG"},
	{Name: "Europe/Nicosia", UTCOffset: "+02:00", Region: "Europe", CountryCode: "CY"},
	{Name: "Europe/Istanbul", UTCOffset: "+03:00", Region: "Europe", CountryCode: "TR"},

	// Middle East
	{Name: "Asia/Dubai", UTCOffset: "+04:00", Region: "Asia", CountryCode: "AE"},
	{Name: "Asia/Muscat", UTCOffset: "+04:00", Region: "Asia", CountryCode: "OM"},
	{Name: "Asia/Bahrain", UTCOffset: "+03:00", Region: "Asia", CountryCode: "BH"},
	{Name: "Asia/Qatar", UTCOffset: "+03:00", Region: "Asia", CountryCode: "QA"},
	{Name: "Asia/Kuwait", UTCOffset: "+03:00", Region: "Asia", CountryCode: "KW"},
	{Name: "Asia/Riyadh", UTCOffset: "+03:00", Region: "Asia", CountryCode: "SA"},

	// Africa
	{Name: "Africa/Cairo", UTCOffset: "+02:00", Region: "Africa", CountryCode: "EG"},
	{Name: "Africa/Casablanca", UTCOffset: "+01:00", Region: "Africa", CountryCode: "MA"},
	{Name: "Africa/Johannesburg", UTCOffset: "+02:00", Region: "Africa", CountryCode: "ZA"},
	{Name: "Africa/Lagos", UTCOffset: "+01:00", Region: "Africa", CountryCode: "NG"},
	{Name: "Africa/Nairobi", UTCOffset: "+03:00", Region: "Africa", CountryCode: "KE"},

	// Americas
	{Name: "America/New_York", UTCOffset: "-05:00", Region: "America", CountryCode: "US"},
	{Name: "America/Chicago", UTCOffset: "-06:00", Region: "America", CountryCode: "US"},
	{Name: "America/Denver", UTCOffset: "-07:00", Region: "America", CountryCode: "US"},
	{Name: "America/Los_Angeles", UTCOffset: "-08:00", Region: "America", CountryCode: "US"},
	{Name: "America/Toronto", UTCOffset: "-05:00", Region: "America", CountryCode: "CA"},
	{Name: "America/Vancouver", UTCOffset: "-08:00", Region: "America", CountryCode: "CA"},
	{Name: "America/Mexico_City", UTCOffset: "-06:00", Region: "America", CountryCode: "MX"},
	{Name: "America/Sao_Paulo", UTCOffset: "-03:00", Region: "America", CountryCode: "BR"},
	{Name: "America/Argentina/Buenos_Aires", UTCOffset: "-03:00", Region: "America", CountryCode: "AR"},
	{Name: "America/Bogota", UTCOffset: "-05:00", Region: "America", CountryCode: "CO"},

	// Asia-Pacific
	{Name: "Asia/Singapore", UTCOffset: "+08:00", Region: "Asia", CountryCode: "SG"},
	{Name: "Asia/Hong_Kong", UTCOffset: "+08:00", Region: "Asia", CountryCode: "HK"},
	{Name: "Asia/Tokyo", UTCOffset: "+09:00", Region: "Asia", CountryCode: "JP"},
	{Name: "Asia/Seoul", UTCOffset: "+09:00", Region: "Asia", CountryCode: "KR"},
	{Name: "Asia/Shanghai", UTCOffset: "+08:00", Region: "Asia", CountryCode: "CN"},
	{Name: "Asia/Kolkata", UTCOffset: "+05:30", Region: "Asia", CountryCode: "IN"},
	{Name: "Australia/Sydney", UTCOffset: "+10:00", Region: "Australia", CountryCode: "AU"},
	{Name: "Australia/Melbourne", UTCOffset: "+10:00", Region: "Australia", CountryCode: "AU"},
	{Name: "Pacific/Auckland", UTCOffset: "+12:00", Region: "Pacific", CountryCode: "NZ"},
}
