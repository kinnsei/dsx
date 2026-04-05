package ref

// Language holds ISO 639-1 language reference data.
type Language struct {
	Code       string // ISO 639-1: "fr"
	Name       string // English name: "French"
	NativeName string // Native name: "Français"
}

var languages = []Language{
	{Code: "ar", Name: "Arabic", NativeName: "العربية"},
	{Code: "bg", Name: "Bulgarian", NativeName: "Български"},
	{Code: "cs", Name: "Czech", NativeName: "Čeština"},
	{Code: "da", Name: "Danish", NativeName: "Dansk"},
	{Code: "de", Name: "German", NativeName: "Deutsch"},
	{Code: "el", Name: "Greek", NativeName: "Ελληνικά"},
	{Code: "en", Name: "English", NativeName: "English"},
	{Code: "es", Name: "Spanish", NativeName: "Español"},
	{Code: "et", Name: "Estonian", NativeName: "Eesti"},
	{Code: "fi", Name: "Finnish", NativeName: "Suomi"},
	{Code: "fr", Name: "French", NativeName: "Français"},
	{Code: "he", Name: "Hebrew", NativeName: "עברית"},
	{Code: "hi", Name: "Hindi", NativeName: "हिन्दी"},
	{Code: "hr", Name: "Croatian", NativeName: "Hrvatski"},
	{Code: "hu", Name: "Hungarian", NativeName: "Magyar"},
	{Code: "id", Name: "Indonesian", NativeName: "Bahasa Indonesia"},
	{Code: "is", Name: "Icelandic", NativeName: "Íslenska"},
	{Code: "it", Name: "Italian", NativeName: "Italiano"},
	{Code: "ja", Name: "Japanese", NativeName: "日本語"},
	{Code: "ko", Name: "Korean", NativeName: "한국어"},
	{Code: "lt", Name: "Lithuanian", NativeName: "Lietuvių"},
	{Code: "lv", Name: "Latvian", NativeName: "Latviešu"},
	{Code: "ms", Name: "Malay", NativeName: "Bahasa Melayu"},
	{Code: "mt", Name: "Maltese", NativeName: "Malti"},
	{Code: "nl", Name: "Dutch", NativeName: "Nederlands"},
	{Code: "no", Name: "Norwegian", NativeName: "Norsk"},
	{Code: "pl", Name: "Polish", NativeName: "Polski"},
	{Code: "pt", Name: "Portuguese", NativeName: "Português"},
	{Code: "ro", Name: "Romanian", NativeName: "Română"},
	{Code: "ru", Name: "Russian", NativeName: "Русский"},
	{Code: "sk", Name: "Slovak", NativeName: "Slovenčina"},
	{Code: "sl", Name: "Slovenian", NativeName: "Slovenščina"},
	{Code: "sv", Name: "Swedish", NativeName: "Svenska"},
	{Code: "sw", Name: "Swahili", NativeName: "Kiswahili"},
	{Code: "th", Name: "Thai", NativeName: "ไทย"},
	{Code: "tr", Name: "Turkish", NativeName: "Türkçe"},
	{Code: "uk", Name: "Ukrainian", NativeName: "Українська"},
	{Code: "vi", Name: "Vietnamese", NativeName: "Tiếng Việt"},
	{Code: "zh", Name: "Chinese", NativeName: "中文"},
}
