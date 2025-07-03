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

type CommonRequestType[BodyType any, CookiesType any] struct {
	Header struct {
		UserAgent *string
	}
	Body    BodyType
	Cookies *CookiesType
}

type CommonResponseType[DataType any, CookiesType any] struct {
	HTTPStatusCode int
	Body           struct {
		Success   bool      `json:"success"`
		Data      *DataType `json:"data"`
		Exception any       `json:"exception"`
	}
	Cookies *CookiesType
}

type E2ETestCase[RequestType any, ResponseType any] struct {
	Request  RequestType
	Response ResponseType
}
