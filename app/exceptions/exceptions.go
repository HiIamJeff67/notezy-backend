package exceptions

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"

	logs "notezy-backend/app/logs"
)

/* ============================== Exception Field Type Definition ============================== */

type ExceptionCode int
type ExceptionPrefix string

const (
	// the first 3 digits are the class of exceptions
	// the last 5 digits are the individual labels for each exceptions
	ExceptionSubDomainCodeShiftAmount = 100000
	MaxExceptionCode                  = 99999999 // 999 99999
	MinExceptionCode                  = 0        // 000 00000
	// reserve some codes for general use purpose
	// see the below general exceptions ex. NotFound(), FailedToCreate()
	ReservedExceptionCode = 100 // *** **100, the codes >= *** **100 will be use in the general domain
)

// all the domain prefix shown here, defined in their corresponded files
// we have 100 codes available to set
const (
// ExceptionPrefix_Util ExceptionPrefix = "Util"       					 1
// ExceptionPrefix_Cookie ExceptionPrefix = "Cookie"					 2
// ExceptionPrefix_Cache ExceptionPrefix = "Cache"	   					 3
// ExceptionPrefix_Context ExceptionPrefix = "Context"					 4
// ExceptionPrefix_Email ExceptionPrefix = "Email"					     5
// ExceptionPrefix_Test ExceptionPrefix = "Test"						 6
// ExceptionPrefix_Search ExceptionPrefix = "Search"			         7
// ExceptionPrefix_Storage ExceptionPrefix = "Storage"				     8
// ExceptionPrefix_Adapter ExceptionPrefix = "Adapter"					 9

// ExceptionPrefix_Auth ExceptionPrefix = "Auth" 			 		     31
// ExceptionPrefix_User ExceptionPrefix = "User"                         32
// ExceptionPrefix_UserInfo ExceptionPrefix = "UserInfo"                 33
// ExceptionPrefix_UserAccount ExceptionPrefix = "UserAccount"           34
// ExceptionPrefix_UserSetting ExceptionPrefix = "UserSetting"           35
// ExceptionPrefix_UsersToBadges ExceptionPrefix = "UsersToBadges"       36
// ExceptionPrefix_Badge ExceptionPrefix = "Badge"                       37
// ExceptionPrefix_Theme ExceptionPrefix = "Theme"						 38
// ExceptionPrefix_Parser ExceptionPrefix = "Parser"				     39
// ExceptionPrefix_Shelf ExceptionPrefix = "Shelf"   					 40
// ExceptionPrefix_Material ExceptionPrefix = "Material"				 41
)

func IsExceptionCode(exceptionCode int) bool {
	return exceptionCode >= MinExceptionCode && exceptionCode <= MaxExceptionCode
}

/* ============================== Location & StackTrace for Exceptions ============================== */

type StackFrame struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
	Package  string `json:"package"`
}

func GetStackTrace(skip int, maxTraceDepth int) []StackFrame {
	var frames []StackFrame

	for i := skip; i < skip+maxTraceDepth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break // the end of the trace stack
		}

		funcName := "unknown"
		packageName := "unknown"

		if fn := runtime.FuncForPC(pc); fn != nil {
			fullFuncName := fn.Name()
			parts := strings.Split(fullFuncName, ".")
			if len(parts) >= 2 {
				packageName = strings.Join(parts[:len(parts)-1], ".")
				funcName = parts[len(parts)-1]
			} else {
				funcName = fullFuncName
			}
		}

		frames = append(frames, StackFrame{
			File:     filepath.Base(file),
			Line:     line,
			Function: funcName,
			Package:  packageName,
		})
	}

	return frames
}

/* ============================== General Exception Structure Definition ============================== */

type Exception struct {
	Code           ExceptionCode   // custom exception code
	Prefix         ExceptionPrefix // custom exception prefix
	Reason         string          // exception reason(for the convenience of frontend to error handling)
	IsInternal     bool            // to indicate whether this exception can passing to the frontend or not
	Message        string          // custom exception message
	HTTPStatusCode int             // http status code
	Details        any             // additional error details (optional)
	Error          error           // original error (optional)
	LastStackFrame *StackFrame     // the last location where the exception happened
	StackTrace     []StackFrame    // the entire path to where the exception actually take place
}

type ExceptionCompareOption struct {
	WithCode           bool
	WithPrefix         bool
	WithReason         bool
	WithIsInteral      bool
	WithMessage        bool
	WithHTTPStatusCode bool
	WithDetails        bool
	WithError          bool
}

func (e *Exception) GetString() string {
	if e.Error != nil {
		return fmt.Sprintf("[%v]%s: %v", e.Code, e.Reason, e.Error)
	}
	return fmt.Sprintf("[%v]%s: %s", e.Code, e.Reason, e.Message)
}

func (e *Exception) GetGinH() *gin.H {
	return &gin.H{
		"code":    e.Code,
		"reason":  e.Reason,
		"prefix":  e.Prefix,
		"message": e.Message,
		"status":  e.HTTPStatusCode,
		"details": e.Details,
		"error":   e.Error,
	}
}

func (e *Exception) ResponseWithJSON(ctx *gin.Context) {
	ctx.JSON(e.HTTPStatusCode, gin.H{
		"success":   false,
		"data":      nil,
		"exception": e.GetGinH(),
	})
}

func (e *Exception) SafelyResponseWithJSON(ctx *gin.Context) {
	if e.IsInternal {
		e = e.InternalServerWentWrong(e)
	}
	ctx.JSON(e.HTTPStatusCode, gin.H{
		"success":   false,
		"data":      nil,
		"exception": e.GetGinH(),
	})
}

func (e *Exception) WithDetails(details any) *Exception {
	e.Details = details
	return e
}

func (e *Exception) WithError(err error) *Exception {
	e.Error = err
	return e
}

func (e *Exception) Log() *Exception {
	if e.Error != nil {
		logs.Alert(e)
		logs.FError("[%d]%s: %s(%v)", e.Code, e.Reason, e.Message, e.Error.Error())
	} else {
		logs.FError("[%d]%s %s", e.Code, e.Reason, e.Message)
	}
	return e
}

func (e *Exception) Panic() {
	if e.Error != nil {
		panic(fmt.Sprintf("[%d]%s: %s(%v)", e.Code, e.Reason, e.Message, e.Error.Error()))
	} else {
		panic(fmt.Sprintf("[%d]%s (%v)", e.Code, e.Reason, e.Message))
	}
}

func (e *Exception) PanicVerbose() {
	if e.Error != nil {
		panic(fmt.Sprintf("[%d]%s: %v", e.Code, e.Reason, e.Error.Error()))
	} else {
		panic(fmt.Sprintf("[%d]%s", e.Code, e.Reason))
	}
}

func (e *Exception) Trace(skip int, maxTraceDepth int) {
	e.StackTrace = GetStackTrace(skip, maxTraceDepth)
}

func (e *Exception) ToGraphQLError(ctx context.Context) *gqlerror.Error {
	extensions := map[string]interface{}{
		"code":       e.Code,
		"prefix":     e.Prefix,
		"httpStatus": e.HTTPStatusCode,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	if e.Details != nil {
		extensions["details"] = e.Details
	}

	var path ast.Path
	var locations []gqlerror.Location

	if ctx != nil {
		if fieldContext := graphql.GetFieldContext(ctx); fieldContext != nil {
			path = fieldContext.Path()

			if fieldContext.Field.Position != nil {
				locations = []gqlerror.Location{
					{
						Line:   fieldContext.Field.Position.Line,
						Column: fieldContext.Field.Position.Column,
					},
				}
			}
		}

		if requestOperationContext := graphql.GetOperationContext(ctx); requestOperationContext != nil {
			if requestOperationContext.OperationName != "" {
				extensions["operationName"] = requestOperationContext.OperationName
			}
		}
	}

	gqlError := &gqlerror.Error{
		Message:    e.Message,
		Path:       path,
		Locations:  locations,
		Extensions: extensions,
	}

	if e.Error != nil {
		gqlError.Err = e.Error
	}

	return gqlError
}

func CompareExceptions(e1 *Exception, e2 *Exception, opt ExceptionCompareOption) bool {
	if opt.WithCode && e1.Code != e2.Code {
		return false
	}
	if opt.WithPrefix && e1.Prefix != e2.Prefix {
		return false
	}
	if opt.WithReason && e1.Reason != e2.Reason {
		return false
	}
	if opt.WithIsInteral && e1.IsInternal != e2.IsInternal {
		return false
	}
	if opt.WithMessage && e1.Message != e2.Message {
		return false
	}
	if opt.WithHTTPStatusCode && e1.HTTPStatusCode != e2.HTTPStatusCode {
		return false
	}
	if opt.WithDetails && fmt.Sprintf("%v", e1.Details) != fmt.Sprintf("%v", e2.Details) {
		return false
	}
	if opt.WithError && fmt.Sprintf("%v", e1.Error) != fmt.Sprintf("%v", e2.Error) {
		return false
	}
	return true
}

func CompareCommonExceptions(e1 *Exception, e2 *Exception, withMessage bool) bool {
	if e1.Code != e2.Code {
		return false
	}
	if e1.Prefix != e2.Prefix {
		return false
	}
	if e1.Reason != e2.Reason {
		return false
	}
	if e1.IsInternal != e2.IsInternal {
		return false
	}
	if withMessage && e1.Message != e2.Message {
		return false
	}
	return true
}

/* ============================== General Exception Define in the Top Layer ============================== */

func UndefinedError(optionalMessage ...string) *Exception {
	message := "Undefined error happened"
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           99900001,
		Prefix:         "General",
		Reason:         "UndefinedError",
		IsInternal:     true,
		Message:        message,
		HTTPStatusCode: http.StatusNotImplemented,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func NotImplemented(optionalMessage ...string) *Exception {
	message := "Not yet implemented the methods"
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           99900002,
		Prefix:         "General",
		Reason:         "NotImplemented",
		IsInternal:     true,
		Message:        message,
		HTTPStatusCode: http.StatusNotImplemented,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (e *Exception) InternalServerWentWrong(originalException *Exception, optionalMessage ...string) *Exception {
	message := "Something went wrong"
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	exception := &Exception{
		Code:           99900003,
		Prefix:         "General",
		Reason:         "InternalServerWentWrong",
		IsInternal:     true,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
	if originalException == nil {
		return exception
	}

	if originalException.Error != nil {
		exception.Error = originalException.Error
	}
	if originalException.Details != nil {
		exception.Message = originalException.Message
	}

	return exception
}

/* ============================== Database Exception Domain Definition ============================== */

type DatabaseExceptionDomain struct {
	_BaseCode ExceptionCode
	_Prefix   ExceptionPrefix
}

func (d *DatabaseExceptionDomain) NotFound(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("%s not found", string(d._Prefix))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 1,
		Prefix:         d._Prefix,
		Reason:         "NotFound",
		IsInternal:     false,
		Message:        message,
		HTTPStatusCode: http.StatusNotFound,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *DatabaseExceptionDomain) FailedToCreate(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to create the %s", string(d._Prefix))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 2,
		Prefix:         d._Prefix,
		Reason:         "FailedToCreate",
		IsInternal:     false,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *DatabaseExceptionDomain) FailedToUpdate(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to update the %s", string(d._Prefix))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 3,
		Prefix:         d._Prefix,
		Reason:         "FailedToUpdate",
		IsInternal:     false,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *DatabaseExceptionDomain) FailedToDelete(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to delete the %s", string(d._Prefix))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 4,
		Prefix:         d._Prefix,
		Reason:         "FailedToDelete",
		IsInternal:     false,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *DatabaseExceptionDomain) NoChanges() *Exception {
	return &Exception{
		Code:           d._BaseCode + 5,
		Prefix:         d._Prefix,
		Reason:         "NoChanges",
		IsInternal:     false,
		Message:        fmt.Sprintf("No Changes on %s", string(d._Prefix)),
		HTTPStatusCode: http.StatusNotModified,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *DatabaseExceptionDomain) FailedToCommitTransaction(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to commit the transaction in %s", string(d._Prefix))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 6,
		Prefix:         d._Prefix,
		Reason:         "FailedToCommitTransaction",
		IsInternal:     true,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== API Exception Domain Definition ============================== */

type APIExceptionDomain struct {
	_BaseCode ExceptionCode
	_Prefix   ExceptionPrefix
}

func (d *APIExceptionDomain) Timeout(time time.Duration, optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Timeout in %s with %v", string(d._Prefix), time)
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 11,
		Prefix:         d._Prefix,
		Reason:         "Timeout",
		IsInternal:     false,
		Message:        message,
		HTTPStatusCode: http.StatusRequestTimeout,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== GraphQL Exception Domain Definition ============================== */

type GraphQLExceptionDomain struct {
	_BaseCode ExceptionCode
	_Prefix   ExceptionPrefix
}

func (d *GraphQLExceptionDomain) InvalidSourceInBatchFunction() *Exception {
	return &Exception{
		Code:           d._BaseCode + 21,
		Prefix:         d._Prefix,
		Reason:         "InvalidSourceInBatchFunction",
		IsInternal:     true,
		Message:        fmt.Sprintf("Invalid source field detected while working on jobs in the batch function of %s", d._Prefix),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== Type Exception Domain Definition ============================== */

type TypeExceptionDomain struct {
	_BaseCode ExceptionCode
	_Prefix   ExceptionPrefix
}

func (d *TypeExceptionDomain) InvalidInput(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Invalid input object detected in %s", string(d._Prefix))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 31,
		Prefix:         d._Prefix,
		Reason:         "InvalidInput",
		IsInternal:     true,
		Message:        message,
		HTTPStatusCode: http.StatusBadRequest,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TypeExceptionDomain) InvalidDto(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Invalid dto detected in %s", string(d._Prefix))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 32,
		Prefix:         d._Prefix,
		Reason:         "InvalidDto",
		IsInternal:     false,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TypeExceptionDomain) InvalidType(value any) *Exception {
	return &Exception{
		Code:           d._BaseCode + 33,
		Prefix:         d._Prefix,
		Reason:         "InvalidType",
		IsInternal:     true,
		Message:        fmt.Sprintf("Invalid type in %s", string(d._Prefix)),
		HTTPStatusCode: http.StatusInternalServerError,
		Details: map[string]any{
			"actualType": fmt.Sprintf("%T", value),
			"value":      value,
		},
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== File Exception Domain Definition ============================== */

type FileExceptionDomain struct {
	_BaseCode ExceptionCode
	_Prefix   ExceptionPrefix
}

// action should be for examples: "read", "write", "execute", "update", "delete", etc.
func (d *FileExceptionDomain) NoPermission(action string) *Exception {
	return &Exception{
		Code:           d._BaseCode + 41,
		Prefix:         d._Prefix,
		Reason:         "NoPermission",
		IsInternal:     false,
		Message:        fmt.Sprintf("You don't have any permission to %s this file", action),
		HTTPStatusCode: http.StatusBadRequest,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *FileExceptionDomain) FileTooLarge(fileSize int64, maxFileSize int64) *Exception {
	return &Exception{
		Code:           d._BaseCode + 42,
		Prefix:         d._Prefix,
		Reason:         "FileTooLarge",
		IsInternal:     false,
		Message:        fmt.Sprintf("File size of %d bytes is too large which excceed the max size of %d bytes", fileSize, maxFileSize),
		HTTPStatusCode: http.StatusTooManyRequests,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *FileExceptionDomain) TooManyFiles(numberOfFiles int64) *Exception {
	return &Exception{
		Code:           d._BaseCode + 43,
		Prefix:         d._Prefix,
		Reason:         "TooManyFiles",
		IsInternal:     false,
		Message:        fmt.Sprintf("Passing %d of files is not allowed in this operation", numberOfFiles),
		HTTPStatusCode: http.StatusTooManyRequests,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *FileExceptionDomain) CannotGetFileObjects() *Exception {
	return &Exception{
		Code:           d._BaseCode + 44,
		Prefix:         d._Prefix,
		Reason:         "CannotGetFileObjects",
		IsInternal:     true,
		Message:        "Failed to get file objects",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *FileExceptionDomain) CannotOpenFiles() *Exception {
	return &Exception{
		Code:           d._BaseCode + 45,
		Prefix:         d._Prefix,
		Reason:         "CannotOpenFiles",
		IsInternal:     true,
		Message:        "Failed to open the files",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *FileExceptionDomain) CannotPeekFiles() *Exception {
	return &Exception{
		Code:           d._BaseCode + 46,
		Prefix:         d._Prefix,
		Reason:         "CannotPeekFiles",
		IsInternal:     true,
		Message:        "Failed to peek the files",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *FileExceptionDomain) CannotCloseFiles() *Exception {
	return &Exception{
		Code:           d._BaseCode + 47,
		Prefix:         d._Prefix,
		Reason:         "CannotCloseFiles",
		IsInternal:     true,
		Message:        "Failed to close the files",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *FileExceptionDomain) CannotReadFileBytes() *Exception {
	return &Exception{
		Code:           d._BaseCode + 48,
		Prefix:         d._Prefix,
		Reason:         "CannotReadFileBytes",
		IsInternal:     true,
		Message:        "Failed to read the file into bytes",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *FileExceptionDomain) FailedToDetectContentType() *Exception {
	return &Exception{
		Code:           d._BaseCode + 49,
		Prefix:         d._Prefix,
		Reason:         "FailedToDetectContentType",
		IsInternal:     true,
		Message:        "Failed to detect content type, the given file may be invalid",
		HTTPStatusCode: http.StatusUnsupportedMediaType,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== Test Exception Domain Definition ============================== */

type TestExceptionDomain struct {
	_BaseCode ExceptionCode
	_Prefix   ExceptionPrefix
}

func (d *TestExceptionDomain) FailedToMarshalTestdata(testdataPath string) *Exception {
	message := fmt.Sprintf("Failed to marshal testdata from %v", testdataPath)

	return &Exception{
		Code:           d._BaseCode + 91,
		Prefix:         d._Prefix,
		Reason:         "FailedToMarshalTestdata",
		IsInternal:     true,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TestExceptionDomain) FailedToUnmarshalTestdata(testdataPath string) *Exception {
	message := fmt.Sprintf("Failed to unmarshal testdata from %v", testdataPath)

	return &Exception{
		Code:           d._BaseCode + 92,
		Prefix:         d._Prefix,
		Reason:         "FailedToUnmarshalTestdata",
		IsInternal:     true,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TestExceptionDomain) InvalidTestdataJSONForm(testdataPath string) *Exception {
	message := fmt.Sprintf("Invalid testdata json form from %v", testdataPath)

	return &Exception{
		Code:           d._BaseCode + 93,
		Prefix:         d._Prefix,
		Reason:         "InvalidTestdataJSONForm",
		IsInternal:     true,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
