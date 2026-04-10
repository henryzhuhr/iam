// Package constant provides error codes and constants for the IAM application.
package constant

import "errors"

// Business error codes (module 2-digit + business 3-digit).
const (
	CodeOK              = 0
	CodeAuthFailed      = 10001 // wrong username or password
	CodeAuthLocked      = 10002 // account locked
	CodeAuthMFAFail     = 10003 // MFA verification failed
	CodeAuthCodeErr     = 10004 // verification code error
	CodeAuthCodeExp     = 10005 // verification code expired

	CodeTokenExpired  = 11001 // token expired
	CodeTokenInvalid  = 11002 // token invalid
	CodeTokenRevoked  = 11003 // token revoked
	CodeTokenRefresh  = 11004 // refresh token invalid

	CodeUserNotFound     = 20001 // user not found
	CodeUserEmailExists  = 20002 // email already exists
	CodeUserDisabled     = 20003 // user disabled
	CodeUserPasswordPolicy = 20004 // password does not meet policy

	CodeRoleNotFound    = 30001 // role not found
	CodeRoleCodeExists  = 30002 // role code duplicate
	CodeRoleSoDConflict = 30003 // SoD conflict
	CodeRolePermDenied  = 30004 // insufficient permission

	CodeTenantNotFound   = 40001 // tenant not found
	CodeTenantExpired    = 40002 // tenant expired
	CodeTenantDisabled   = 40003 // tenant disabled
	CodeTenantQuotaExceed = 40004 // quota exceeded

	CodeAppNotFound      = 50001 // application not found
	CodeAppDisabled      = 50002 // application disabled
	CodeAppNotAuthorized = 50003 // user not authorized for this app

	CodeClientNotFound = 60001 // client not found
	CodeClientAKSKInvalid = 60002 // AK/SK invalid
	CodeClientDisabled = 60003 // client disabled

	CodeInternalError = 99001 // internal server error
	CodeDBError       = 99002 // database error
	CodeRedisError    = 99003 // Redis error
	CodeKafkaError    = 99004 // Kafka error
)

// Sentinel errors for internal use.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrTenantExpired  = errors.New("tenant expired")
	ErrTenantDisabled = errors.New("tenant disabled")
)
