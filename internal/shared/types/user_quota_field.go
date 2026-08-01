package types

type UserQuotaField string

func (uqf UserQuotaField) String() string {
	return string(uqf)
}
