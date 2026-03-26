package types

type RefreshableResponseFieldName string

const (
	RefreshableResponseFieldName_NewAccessToken RefreshableResponseFieldName = "newAccessToken"
	RefreshableResponseFieldName_NewCSRFToken   RefreshableResponseFieldName = "newCSRFToken"
)

func (rrfn RefreshableResponseFieldName) String() string {
	return string(rrfn)
}
