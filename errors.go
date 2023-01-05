package ggpo

import "fmt"

type Error struct {
	Code ErrorCode
	Name string
}

func (e Error) Error() string {
	return fmt.Sprintf("ggpo: %s:%d", e.Name, e.Code)
}

type ErrorCode int

const (
	Ok                           ErrorCode = 0
	ErrorCodeSuccess             ErrorCode = 0
	ErrorCodeGeneralFailure      ErrorCode = -1
	ErrorCodeInvalidSession      ErrorCode = 1
	ErrorCodeInvalidPlayerHandle ErrorCode = 2
	ErrorCodePlayerOutOfRange    ErrorCode = 3
	ErrorCodePredictionThreshod  ErrorCode = 4
	ErrorCodeUnsupported         ErrorCode = 5
	ErrorCodeNotSynchronized     ErrorCode = 6
	ErrorCodeInRollback          ErrorCode = 7
	ErrorCodeInputDropped        ErrorCode = 8
	ErrorCodePlayerDisconnected  ErrorCode = 9
	ErrorCodeTooManySpectators   ErrorCode = 10
	ErrorCodeInvalidRequest      ErrorCode = 11
)

func Success(result ErrorCode) bool {
	return result == ErrorCodeSuccess
}
