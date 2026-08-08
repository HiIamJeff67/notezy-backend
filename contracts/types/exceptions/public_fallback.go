package exceptions

type PublicFallback struct {
	Reason         ExceptionReason
	Domain         string
	Operation      string
	Message        string
	HTTPStatusCode int
	Retryable      bool
}
