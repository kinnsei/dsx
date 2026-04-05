package ref

// Currency holds ISO 4217 currency reference data.
type Currency struct {
	Code     string // ISO 4217: "EUR"
	Name     string // English name: "Euro"
	Symbol   string // Display symbol: "€"
	Decimals int    // Minor unit digits: 2
}

var currencies = []Currency{
	// Major
	{Code: "EUR", Name: "Euro", Symbol: "€", Decimals: 2},
	{Code: "USD", Name: "US Dollar", Symbol: "$", Decimals: 2},
	{Code: "GBP", Name: "British Pound", Symbol: "£", Decimals: 2},
	{Code: "CHF", Name: "Swiss Franc", Symbol: "CHF", Decimals: 2},
	{Code: "JPY", Name: "Japanese Yen", Symbol: "¥", Decimals: 0},
	{Code: "CNY", Name: "Chinese Yuan", Symbol: "¥", Decimals: 2},

	// Middle East
	{Code: "AED", Name: "UAE Dirham", Symbol: "د.إ", Decimals: 2},
	{Code: "SAR", Name: "Saudi Riyal", Symbol: "﷼", Decimals: 2},
	{Code: "QAR", Name: "Qatari Riyal", Symbol: "﷼", Decimals: 2},
	{Code: "BHD", Name: "Bahraini Dinar", Symbol: "BD", Decimals: 3},
	{Code: "KWD", Name: "Kuwaiti Dinar", Symbol: "KD", Decimals: 3},
	{Code: "OMR", Name: "Omani Rial", Symbol: "﷼", Decimals: 3},

	// Asia-Pacific
	{Code: "AUD", Name: "Australian Dollar", Symbol: "A$", Decimals: 2},
	{Code: "NZD", Name: "New Zealand Dollar", Symbol: "NZ$", Decimals: 2},
	{Code: "SGD", Name: "Singapore Dollar", Symbol: "S$", Decimals: 2},
	{Code: "HKD", Name: "Hong Kong Dollar", Symbol: "HK$", Decimals: 2},
	{Code: "KRW", Name: "South Korean Won", Symbol: "₩", Decimals: 0},
	{Code: "INR", Name: "Indian Rupee", Symbol: "₹", Decimals: 2},

	// Americas
	{Code: "CAD", Name: "Canadian Dollar", Symbol: "C$", Decimals: 2},
	{Code: "MXN", Name: "Mexican Peso", Symbol: "MX$", Decimals: 2},
	{Code: "BRL", Name: "Brazilian Real", Symbol: "R$", Decimals: 2},
	{Code: "ARS", Name: "Argentine Peso", Symbol: "AR$", Decimals: 2},
	{Code: "COP", Name: "Colombian Peso", Symbol: "CO$", Decimals: 2},

	// Nordics
	{Code: "SEK", Name: "Swedish Krona", Symbol: "kr", Decimals: 2},
	{Code: "NOK", Name: "Norwegian Krone", Symbol: "kr", Decimals: 2},
	{Code: "DKK", Name: "Danish Krone", Symbol: "kr", Decimals: 2},
	{Code: "ISK", Name: "Icelandic Krona", Symbol: "kr", Decimals: 0},

	// EU non-Euro
	{Code: "PLN", Name: "Polish Zloty", Symbol: "zł", Decimals: 2},
	{Code: "CZK", Name: "Czech Koruna", Symbol: "Kč", Decimals: 2},
	{Code: "HUF", Name: "Hungarian Forint", Symbol: "Ft", Decimals: 2},
	{Code: "RON", Name: "Romanian Leu", Symbol: "lei", Decimals: 2},
	{Code: "BGN", Name: "Bulgarian Lev", Symbol: "лв", Decimals: 2},

	// Africa
	{Code: "ZAR", Name: "South African Rand", Symbol: "R", Decimals: 2},
	{Code: "NGN", Name: "Nigerian Naira", Symbol: "₦", Decimals: 2},
	{Code: "EGP", Name: "Egyptian Pound", Symbol: "E£", Decimals: 2},
	{Code: "MAD", Name: "Moroccan Dirham", Symbol: "MAD", Decimals: 2},
	{Code: "KES", Name: "Kenyan Shilling", Symbol: "KSh", Decimals: 2},

	// Turkey
	{Code: "TRY", Name: "Turkish Lira", Symbol: "₺", Decimals: 2},
}
