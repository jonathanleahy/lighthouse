package message

import "github.com/pismo/backoffice-core-bff/internal/app/infrastructure/apierror/message"

var (
	//Internal Error Messages
	ErrApiConnection = message.ErrorMessage{
		ErrorCode:   ErrorCodePrefix + "0001",
		UserMessage: "%s",
	}
	ErrApiBodyParsing = message.ErrorMessage{
		ErrorCode:   ErrorCodePrefix + "0002",
		UserMessage: "Error parsing utils body with message: %s",
	}
	ErrQuery = message.ErrorMessage{
		ErrorCode:   ErrorCodePrefix + "0014",
		UserMessage: "Query not allowed",
	}
	ErrSNSBodyParsingErr = message.ErrorMessage{
		ErrorCode:   ErrorCodePrefix + "0015",
		UserMessage: "Error parsing utils body with message: %s",
	}
	ErrSNSConnectionErr = message.ErrorMessage{
		ErrorCode:   ErrorCodePrefix + "0016",
		UserMessage: "Connection failed with message: %s",
	}

	//Other Error Messages
	ErrForbidden = message.ErrorMessage{
		ErrorCode:   WarningCodePrefix + "0003",
		UserMessage: "Access Denied",
	}

	ErrInternalServerError = message.ErrorMessage{
		ErrorCode:   ErrorCodePrefix + "0500",
		UserMessage: "Internal Server Error",
	}
)

const (
	Debug = iota
	Warning
	Error
	ErrorCodePrefix   = "ECPBFF"
	WarningCodePrefix = "WCPBFF"

	//Error Messages
	ErrorRequestTimeout = "{message: Request timeout}"

	//Success
	RequestSuccess = "Success"
	RequestTimeout = "Request timeout"

	// api erros
	ApiExecInfoMessage  = "API utils successful."
	ApiExecErrorMessage = "API utils error."

	// callers
	CallerApiError           = "api_error"
	CallerApiConnectionError = "api_connection_error"
	CallerBodyParsingError   = "api_request_body_parsing_error"

	// SNS erros
	SNSExecInfoMessage  = "SNS successful."
	SNSExecErrorMessage = "SNS error."
)

