package exceptions

const (
	_ExceptionBaseCode_Cache ExceptionCode = 700000
	ExceptionBaseCode_Cache ExceptionCode = _ExceptionBaseCode_Cache + ReservedExceptionCode
	ExceptionPrefix_Cache ExceptionPrefix = "Cache"
)

type CacheExceptionDomain struct {
	APIExceptionDomain
}

var Cache = &CacheExceptionDomain{
	APIExceptionDomain{ BaseCode: _ExceptionBaseCode_Cache, Prefix: ExceptionPrefix_Cache },
}
