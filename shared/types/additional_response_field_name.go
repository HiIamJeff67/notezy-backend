package types

type AdditionalResponseFieldDomainName string

const (
	AdditionalResponseFieldDomainName_RefreshableTokens AdditionalResponseFieldDomainName = "refreshableTokens"
	AdditionalResponseFieldDomainName_Embedded          AdditionalResponseFieldDomainName = "embedded"
)

func (arfdn AdditionalResponseFieldDomainName) String() string {
	return string(arfdn)
}

type RefreshableResponseFieldName string

const (
	RefreshableResponseFieldName_NewAccessToken RefreshableResponseFieldName = "newAccessToken"
	RefreshableResponseFieldName_NewCSRFToken   RefreshableResponseFieldName = "newCSRFToken"
)

func (rrfn RefreshableResponseFieldName) String() string {
	return string(rrfn)
}

type EmbedAuthorizedResponseFieldName string

const (
	EmbeddedAuthorizedResponseFieldName_PublicId EmbedAuthorizedResponseFieldName = "embeddedPublicId"
)

func (earfn EmbedAuthorizedResponseFieldName) String() string {
	return string(earfn)
}
