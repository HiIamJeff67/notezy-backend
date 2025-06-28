package test

/* ============================== Test Case for Unit Test ============================== */

type UnitTestCase[ArgType any, ReturnType any] struct {
	Args    ArgType
	Returns ReturnType
}

/* ============================== Test Case for Testing E2E ============================== */

type CommonCookiesType struct {
	AccessToken  string
	RefreshToken string
}

type CommonRequestType struct {
	Header struct {
		UserAgent string
	}
	Body    any
	Cookies CommonCookiesType
}

type CommonResponseType struct {
	HTTPStatusCode int
	Result         *struct {
		Message string
		Data    any
	}
	Exception any
	Cookies   CommonCookiesType
}

type E2ETestCase[RequestType any, ResponseType any] struct {
	Request  RequestType
	Response ResponseType
}
