package exceptions

import "fmt"

type Pair struct {
	First  bool
	Second *Exception
}

type Exception struct {
	Reason    ExceptionReason `json:"reason"`
	Domain    string          `json:"domain"`
	Operation string          `json:"operation"`
	Message   string          `json:"message"`
	Retryable bool            `json:"retryable"`

	httpStatusCode int
	isInternal     bool
	details        any
	origin         error
	publicFallback *PublicFallback
}

func New(
	reason ExceptionReason,
	domain string,
	operation string,
	message string,
	httpStatusCode int,
	isInternal ...bool,
) *Exception {
	internal := len(isInternal) > 0 && isInternal[0]

	return &Exception{
		Reason:         reason,
		Domain:         domain,
		Operation:      operation,
		Message:        message,
		httpStatusCode: httpStatusCode,
		isInternal:     internal,
	}
}

func (e *Exception) String() string {
	if e == nil {
		return "exception: <nil>"
	}

	message := fmt.Sprintf(
		"exception:\n  reason: %s\n  domain: %s\n  operation: %s\n  message: %s\n  retryable: %t\n  httpStatusCode: %d\n  isInternal: %t",
		e.Reason,
		e.Domain,
		e.Operation,
		e.Message,
		e.Retryable,
		e.httpStatusCode,
		e.isInternal,
	)
	if e.origin != nil {
		message += fmt.Sprintf("\n  origin: %v", e.origin)
	}
	if e.details != nil {
		message += fmt.Sprintf("\n  details: %#v", e.details)
	}
	if e.publicFallback != nil {
		message += fmt.Sprintf("\n  publicFallback: %#v", e.publicFallback)
	}

	return message
}

func (e *Exception) Clone(httpStatusCode ...int) *Exception {
	if e == nil {
		return nil
	}

	clone := *e
	if e.publicFallback != nil {
		publicFallback := *e.publicFallback
		clone.publicFallback = &publicFallback
	}
	if len(httpStatusCode) > 0 {
		clone.httpStatusCode = httpStatusCode[0]
	}

	return &clone
}

func (e *Exception) WithDetails(details any) *Exception {
	e.details = details
	return e
}

func (e *Exception) WithOrigin(origin error) *Exception {
	e.origin = origin
	return e
}

func (e *Exception) WithPublicFallback(fallback PublicFallback) *Exception {
	e.publicFallback = &fallback
	return e
}

func (e *Exception) ToPublic() *Exception {
	if e == nil {
		return InternalServerError()
	}

	if e.isInternal {
		if e.publicFallback == nil {
			return InternalServerError()
		}

		return &Exception{
			Reason:         e.publicFallback.Reason,
			Domain:         e.publicFallback.Domain,
			Operation:      e.publicFallback.Operation,
			Message:        e.publicFallback.Message,
			httpStatusCode: e.publicFallback.HTTPStatusCode,
			Retryable:      e.publicFallback.Retryable,
		}
	}

	return &Exception{
		Reason:         e.Reason,
		Domain:         e.Domain,
		Operation:      e.Operation,
		Message:        e.Message,
		httpStatusCode: e.httpStatusCode,
		Retryable:      e.Retryable,
	}
}

func (e *Exception) HTTPStatusCode() int {
	return e.httpStatusCode
}

func (e *Exception) IsInternal() bool {
	return e.isInternal
}

func (e *Exception) Origin() error {
	return e.origin
}

func (e *Exception) Error() string {
	if e == nil {
		return "exception: <nil>"
	}

	if e.origin != nil {
		return e.origin.Error()
	}

	return e.Message
}

func Cover(existing *Exception, fallbacks []Pair) *Exception {
	if existing != nil {
		return existing
	}

	for _, pair := range fallbacks {
		if pair.First && pair.Second != nil {
			return pair.Second
		}
	}

	return nil
}
