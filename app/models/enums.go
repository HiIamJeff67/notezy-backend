package models

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

type Enum interface {
	Name() string
	Scan(value any) error
	Value() (driver.Value, error)
}

func _scanError(value any, e Enum) error {
	// A Helper Function to Get the Error
	return fmt.Errorf("failed to scan %T into %s", value, e.Name())
}

/* ============================== UserRole Definition ============================== */
type UserRole string

const (
	UserRole_Admin  UserRole = "Admin"
	UserRole_Noraml UserRole = "Normal"
	UserRole_Guest  UserRole = "Guest"
)

func (r *UserRole) Name() string {
	return reflect.TypeOf(r).Name()
}

// Scan() makes UserRole support automatically convert type from string in database to UserRole in codebase
func (r *UserRole) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*r = UserRole(string(v))
		return nil
	case string:
		*r = UserRole(v)
		return nil
	}
	return _scanError(value, r)
}

// Value() makes UserRole support automatically convert from UserRole in codebase to string in database
func (r UserRole) Value() (driver.Value, error) {
	return string(r), nil
}

var AllUserRoles = []UserRole{
	UserRole_Admin,
	UserRole_Noraml,
	UserRole_Guest,
}
var AllUserRoleStrings = []string{
	string(UserRole_Admin),
	string(UserRole_Noraml),
	string(UserRole_Guest),
}

/* ============================== UserPlan Definition ============================== */
type UserPlan string

const (
	UserPlan_Enterprise UserPlan = "Enterprise"
	UserPlan_Ultimate   UserPlan = "Ultimate"
	UserPlan_Pro        UserPlan = "Pro"
	UserPlan_Free       UserPlan = "Free"
)

func (p *UserPlan) Name() string {
	return reflect.TypeOf(p).Name()
}

func (p *UserPlan) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*p = UserPlan(string(v))
		return nil
	case string:
		*p = UserPlan(v)
		return nil
	}
	return _scanError(value, p)
}

func (p UserPlan) Value() (driver.Value, error) {
	return string(p), nil
}

var AllUserPlans = []UserPlan{
	UserPlan_Enterprise,
	UserPlan_Ultimate,
	UserPlan_Pro,
	UserPlan_Free,
}
var AllUserPlanStrings = []string{
	string(UserPlan_Enterprise),
	string(UserPlan_Ultimate),
	string(UserPlan_Pro),
	string(UserPlan_Free),
}

/* ============================== UserStatus Definition ============================== */
type UserStatus string

const (
	UserStatus_Online       UserStatus = "Online"
	UserStatus_AFK          UserStatus = "AFK"
	UserStatus_DoNotDisturb UserStatus = "DoNotDisturb"
	UserStatus_Offline      UserStatus = "Offline"
)

func (s *UserStatus) Name() string {
	return reflect.TypeOf(s).Name()
}

func (s *UserStatus) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*s = UserStatus(string(v))
		return nil
	case string:
		*s = UserStatus(v)
		return nil
	}
	return _scanError(value, s)
}

func (s UserStatus) Value() (driver.Value, error) {
	return string(s), nil
}

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

/* ============================== UserGener Definition ============================== */
type UserGender string

const (
	UserGender_Male           UserGender = "Male"
	UserGender_Female         UserGender = "Female"
	UserGender_PreferNotToSay UserGender = "PreferNotToSay"
)

func (g *UserGender) Name() string {
	return reflect.TypeOf(g).Name()
}

func (g *UserGender) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*g = UserGender(string(v))
		return nil
	case string:
		*g = UserGender(v)
		return nil
	}
	return _scanError(value, g)
}

func (g UserGender) Value() (driver.Value, error) {
	return string(g), nil
}

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

/* ============================== Country Definition ============================== */
type Country string

const (
	Country_Default               Country = "Default" // null value
	Country_Taiwan                Country = "Taiwan"
	Country_Japan                 Country = "Japan"
	Country_Malaysia              Country = "Malaysia"
	Country_Singapore             Country = "Singapore"
	Country_China                 Country = "China"
	Country_UnitedStatusOfAmerica Country = "UnitedStatusOfAmerica"
	Country_UnitedKingdom         Country = "UnitedKingdom"
	Country_Australia             Country = "Australia"
	Country_Canada                Country = "Canada"
)

func (c *Country) Name() string {
	return reflect.TypeOf(c).Name()
}

func (c *Country) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*c = Country(string(v))
		return nil
	case string:
		*c = Country(v)
		return nil
	}
	return _scanError(value, c)
}

func (c Country) Value() (driver.Value, error) {
	return string(c), nil
}

var AllCountries = []Country{
	Country_Default,
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
	string(Country_Default),
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

/* ============================== CountryCode Definition ============================== */
type CountryCode string

const (
	CountryCode_Default       CountryCode = "Default" // null value
	CountryCode_Taiwan        CountryCode = "+886"
	CountryCode_Japan         CountryCode = "+81"
	CountryCode_Malaysia      CountryCode = "+60"
	CountryCode_Singapore     CountryCode = "+65"
	CountryCode_China         CountryCode = "+86"
	CountryCode_NANP          CountryCode = "+1"
	CountryCode_UnitedKingdom CountryCode = "+44"
	CountryCode_Australia     CountryCode = "+61"
)

func (cc *CountryCode) Name() string {
	return reflect.TypeOf(cc).Name()
}

func (cc *CountryCode) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*cc = CountryCode(string(v))
		return nil
	case string:
		*cc = CountryCode(v)
		return nil
	}
	return _scanError(value, cc)
}

func (cc CountryCode) Value() (driver.Value, error) {
	return string(cc), nil
}

var AllCountryCodes = []CountryCode{
	CountryCode_Default,
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
	string(CountryCode_Default),
	string(CountryCode_Taiwan),
	string(CountryCode_Japan),
	string(CountryCode_Malaysia),
	string(CountryCode_Singapore),
	string(CountryCode_China),
	string(CountryCode_NANP),
	string(CountryCode_UnitedKingdom),
	string(CountryCode_Australia),
}

/* ============================== Theme Definition ============================== */
type Theme string

const (
	Theme_Light  Theme = "Light"
	Theme_Dark   Theme = "Dark"
	Theme_System Theme = "System"
)

func (t *Theme) Name() string {
	return reflect.TypeOf(t).Name()
}

func (t *Theme) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*t = Theme(string(v))
		return nil
	case string:
		*t = Theme(v)
		return nil
	}
	return _scanError(value, t)
}

func (t Theme) Value() (driver.Value, error) {
	return string(t), nil
}

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

/* ============================== Language Definition ============================== */
type Language string

const (
	Language_English            Language = "English"
	Language_TraditionalChinese Language = "TraditionalChinese"
	Language_SimpleChinese      Language = "SimpleChinese"
	Language_Japanese           Language = "Japanese"
)

func (l *Language) Name() string {
	return reflect.TypeOf(l).Name()
}

func (l *Language) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*l = Language(string(v))
		return nil
	case string:
		*l = Language(v)
		return nil
	}
	return _scanError(value, l)
}

func (l Language) Value() (driver.Value, error) {
	return string(l), nil
}

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

/* ============================== BadgeType Definition ============================== */
type BadgeType string

const (
	BadgeType_Diamond BadgeType = "Diamond"
	BadgeType_Golden  BadgeType = "Golden"
	BadgeType_Silver  BadgeType = "Silver"
	BadgeType_Bronze  BadgeType = "Bronze"
	BadgeType_Steel   BadgeType = "Steel"
)

func (bt *BadgeType) Name() string {
	return reflect.TypeOf(bt).Name()
}

func (bt *BadgeType) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*bt = BadgeType(string(v))
		return nil
	case string:
		*bt = BadgeType(v)
		return nil
	}
	return _scanError(value, bt)
}

func (bt BadgeType) Value() (driver.Value, error) {
	return string(bt), nil
}

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

/* ========================= Validator for Validating Enums ========================= */
func IsValidEnumValues[EnumValue interface {
	UserRole |
		UserPlan |
		UserStatus |
		UserGender |
		Country |
		CountryCode |
		Theme |
		Language |
		BadgeType |
		string
}](value EnumValue, validateValues []EnumValue) bool {
	for _, validateValue := range validateValues {
		if value == validateValue {
			return true
		}
	}
	return false
}
