package dtos

type GoogleUserInfoDto struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verifiedEmail"`
	Name          string `json:"name"`
	GivenName     string `json:"givenName"`
	FamilyName    string `json:"familyName"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}
