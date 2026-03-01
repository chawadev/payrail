package payrailCore

import (
	"fmt"
	"log"
	"os"
)

// ErrorLevel defines severity of errors
type ErrorLevel int

const (
	ErrorLevelInfo ErrorLevel = iota
	ErrorLevelWarning
	ErrorLevelError
	ErrorLevelFatal
)

// PaymentError wraps errors from the payment framework with context
type PaymentError struct {
	Level   ErrorLevel
	Message string
	Err     error  // underlying error, may be nil
	Code    string // error code for programmatic handling
}

// Error implements error interface
func (pe *PaymentError) Error() string {
	if pe.Err != nil {
		return fmt.Sprintf("%s: %v", pe.Message, pe.Err)
	}
	return pe.Message
}

// Common error codes
const (
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeProvider     = "PROVIDER_ERROR"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeNetwork      = "NETWORK_ERROR"
	ErrCodeParsing      = "PARSING_ERROR"
	ErrCodeInternal     = "INTERNAL_ERROR"
)

// NewPaymentError creates a new PaymentError
func NewPaymentError(level ErrorLevel, code, message string, err error) *PaymentError {
	return &PaymentError{
		Level:   level,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string) *PaymentError {
	return NewPaymentError(ErrorLevelError, ErrCodeValidation, message, nil)
}

// NewProviderError creates a provider error
func NewProviderError(message string, err error) *PaymentError {
	return NewPaymentError(ErrorLevelError, ErrCodeProvider, message, err)
}

// NewNetworkError creates a network error
func NewNetworkError(message string, err error) *PaymentError {
	return NewPaymentError(ErrorLevelError, ErrCodeNetwork, message, err)
}

// NewParsingError creates a parsing error
func NewParsingError(message string, err error) *PaymentError {
	return NewPaymentError(ErrorLevelError, ErrCodeParsing, message, err)
}

// ErrorLogger provides structured error logging for the framework
type ErrorLogger struct {
	writer *os.File
}

// NewErrorLogger creates a new error logger
func NewErrorLogger() *ErrorLogger {
	return &ErrorLogger{writer: os.Stderr}
}

// LogError logs a payment error with context
func (el *ErrorLogger) LogError(pe *PaymentError) {
	levelStr := levelToString(pe.Level)
	msg := fmt.Sprintf("[%s] %s (Code: %s): %v\n", levelStr, pe.Message, pe.Code, pe.Err)
	el.writer.WriteString(msg)
}

// LogErrorf logs a formatted error message
func (el *ErrorLogger) LogErrorf(level ErrorLevel, format string, args ...interface{}) {
	levelStr := levelToString(level)
	msg := fmt.Sprintf("[%s] %s\n", levelStr, fmt.Sprintf(format, args...))
	el.writer.WriteString(msg)
}

// Fatal logs an error and exits the program
func (el *ErrorLogger) Fatal(pe *PaymentError) {
	el.LogError(pe)
	log.Fatal(pe.Error())
}

// Fatalf logs a formatted error and exits the program
func (el *ErrorLogger) Fatalf(format string, args ...interface{}) {
	el.LogErrorf(ErrorLevelFatal, format, args...)
	log.Fatalf(format, args...)
}

func levelToString(level ErrorLevel) string {
	switch level {
	case ErrorLevelInfo:
		return "INFO"
	case ErrorLevelWarning:
		return "WARNING"
	case ErrorLevelError:
		return "ERROR"
	case ErrorLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
