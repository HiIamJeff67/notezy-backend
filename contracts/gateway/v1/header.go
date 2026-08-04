package gatewaycontract

type Header string

func (h Header) String() string {
	return string(h)
}

const (
	CoreAuthRefreshed  Header = "X-Core-Auth-Refreshed"
	CoreSetAccessToken Header = "X-Core-Set-Access-Token"
	CoreSetCSRFToken   Header = "X-Core-Set-CSRF-Token"
)
