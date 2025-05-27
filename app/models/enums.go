package models

type UserRole string
const (
	UserRole_Admin UserRole = "Admin"
	UserRole_Noraml UserRole = "Normal"
	UserRole_Guest UserRole = "Guest"
)
var AllUserRoles = []UserRole{ 
	UserRole_Admin,
	UserRole_Noraml,
	UserRole_Noraml,
}
var AllUserRoleStrings = []string{
	string(UserRole_Admin),
	string(UserRole_Noraml),
	string(UserRole_Guest),
}

type UserPlan string
const (
	UserPlan_Enterprise UserPlan = "Enterprise"
	UserPlan_Ultimate UserPlan = "Ultimate"
	UserPlan_Pro UserPlan = "Pro"
	UserPlan_Free UserPlan = "Free"
)
var AllUserPlans = []UserPlan{
	UserPlan_Enterprise, 
	UserPlan_Ultimate, 
	UserPlan_Pro, 
	UserPlan_Free, 
}
var AllUserPlanStrings = []string {
	string(UserPlan_Enterprise), 
	string(UserPlan_Ultimate), 
	string(UserPlan_Pro), 
	string(UserPlan_Free), 
}

type UserStatus string
const (
	UserStatus_Online UserStatus = "Online"
	UserStatus_AFK UserStatus = "AFK"
	UserStatus_DoNotDisturb UserStatus = "DoNotDisturb"
	UserStatus_Offline UserStatus = "Offline"
)
var AllUserStatuses = []UserStatus{
	UserStatus_Online, 
	UserStatus_AFK, 
	UserStatus_DoNotDisturb, 
	UserStatus_Offline, 
}
var AllUserStatusStrings = []string{
	string(UserStatus_Online), 
	string(UserStatus_AFK), 
	string(UserStatus_DoNotDisturb), 
	string(UserStatus_Offline), 
}

type UserGender string
const (
	UserGender_Male UserGender = "Male"
	UserGender_Female UserGender = "Female"
	UserGender_PreferNotToSay UserGender = "PreferNotToSay"
)
var AllUserGenders = []UserGender{
	UserGender_Male, 
	UserGender_Female, 
	UserGender_PreferNotToSay, 
}
var AllUserGenderStrings = []string{
	string(UserGender_Male), 
	string(UserGender_Female), 
	string(UserGender_PreferNotToSay), 
}

type Country string
const (
	Country_Taiwan Country = "Taiwan"
	Country_Japan Country = "Japan"
	Country_Malaysia Country = "Malaysia"
	Country_Singapore Country = "Singapore"
	Country_China Country = "China"
	Country_UnitedStatusOfAmerica Country = "UnitedStatusOfAmerica"
	Country_UnitedKingdom Country = "UnitedKingdom"
	Country_Australia Country = "Australia"
	Country_Canada Country = "Canada"
)
var AllCountries = []Country{
	Country_Taiwan, 
	Country_Japan, 
	Country_Malaysia, 
	Country_Singapore, 
	Country_China, 
	Country_UnitedStatusOfAmerica, 
	Country_UnitedKingdom, 
	Country_Australia, 
	Country_Canada, 
}
var AllCountryStrings = []string{
	string(Country_Taiwan), 
	string(Country_Japan), 
	string(Country_Malaysia), 
	string(Country_Singapore), 
	string(Country_China), 
	string(Country_UnitedStatusOfAmerica), 
	string(Country_UnitedKingdom), 
	string(Country_Australia), 
	string(Country_Canada), 
}

type CountryCode string
const (
	CountryCode_Taiwan CountryCode = "+886"
	CountryCode_Japan CountryCode = "+81"
	CountryCode_Malaysia CountryCode = "+60"
	CountryCode_Singapore CountryCode = "+65"
	CountryCode_China CountryCode = "+86"
	CountryCode_NANP CountryCode = "+1"
	CountryCode_UnitedKingdom CountryCode = "+44"
	CountryCode_Australia CountryCode = "+61"
)
var AllCountryCodes = []CountryCode{
	CountryCode_Taiwan, 
	CountryCode_Japan, 
	CountryCode_Malaysia, 
	CountryCode_Singapore, 
	CountryCode_China, 
	CountryCode_NANP, // NANP stands for North American Numbering Plan, it's used in United States of America and Canada
	CountryCode_UnitedKingdom, 
	CountryCode_Australia, 
}
var AllCountryCodeStrings = []string{
	string(CountryCode_Taiwan), 
	string(CountryCode_Japan), 
	string(CountryCode_Malaysia), 
	string(CountryCode_Singapore), 
	string(CountryCode_China), 
	string(CountryCode_NANP), 
	string(CountryCode_UnitedKingdom), 
	string(CountryCode_Australia), 
}

type Theme string
const (
	Theme_Light Theme = "Light"
	Theme_Dark Theme = "Dark"
	Theme_System Theme = "System"
)
var AllThemes = []Theme{
	Theme_Light,
	Theme_Dark,
	Theme_System,
}
var AllThemeStrings = []string{
	string(Theme_Light), 
	string(Theme_Dark),
	string(Theme_System),
}

type Language string
const (
	Language_English Language = "English"
	Language_TraditionalChinese Language = "TraditionalChinese"
	Language_SimpleChinese Language = "SimpleChinese"
	Language_Japanese Language = "Japanese"
)
var AllLanguages = []Language{
	Language_English, 
	Language_TraditionalChinese, 
	Language_SimpleChinese, 
	Language_Japanese,
}
var AllLanguageStrings = []string{
	string(Language_English), 
	string(Language_TraditionalChinese), 
	string(Language_SimpleChinese), 
	string(Language_Japanese),
}

type BadgeType string
const (
	BadgeType_Diamond BadgeType = "Diamond"
	BadgeType_Golden BadgeType = "Golden"
	BadgeType_Silver BadgeType = "Silver"
	BadgeType_Bronze BadgeType = "Bronze"
	BadgeType_Steel BadgeType = "Steel"
)
var AllBadgeTypes = []BadgeType{
	BadgeType_Diamond, 
	BadgeType_Golden, 
	BadgeType_Silver, 
	BadgeType_Bronze, 
	BadgeType_Steel, 
}
var AllBadgeTypeStrings = []string{
	string(BadgeType_Diamond), 
	string(BadgeType_Golden), 
	string(BadgeType_Silver), 
	string(BadgeType_Bronze), 
	string(BadgeType_Steel), 
}

/* ============================== Validator for Validating Enums ============================== */
func IsValidEnumValues[EnumValue interface {
    UserRole       | 
    UserPlan       | 
    UserStatus     | 
    UserGender     | 
    Country        | 
    CountryCode    | 
    Theme          | 
    Language       | 
    BadgeType	   |
	string
}](value EnumValue, validateValues []EnumValue) bool {
    for _, validateValue := range validateValues {
		if value == validateValue { return true; }
	}
	return false;
}
/* ============================== Validator for Validating Enums ============================== */