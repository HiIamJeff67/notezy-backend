package types

type UserQuotaField string

const ()

func (uqf UserQuotaField) String() string {
	return string(uqf)
}
