package constants

// 错误码集中维护；各 service/handler 按实体与字段手动拼接 message 透传。
const (
	CodeOK               = 0
	CodeBadRequest       = 40000
	CodeUnauthorized     = 40100
	CodeForbidden        = 40300
	CodeNotFound         = 40400
	CodeConflict         = 40900
	CodeValidationFailed = 42200
	CodeRateLimited      = 42900
	CodeInternalError    = 50000

	CodeUserExists         = 10001
	CodeUserNotFound       = 10002
	CodeInvalidCredential  = 10003
	CodeUserBanned         = 10004
	CodeWishNotFound       = 20001
	CodeWishNotOwner       = 20002
	CodeWishAlreadyClaimed = 20003
	CodeWishStatusInvalid  = 20004
	CodeClaimNotFound      = 30001
	CodeClaimNotOwner      = 30002
	CodeClaimStatusInvalid = 30003

	CodeAcceptanceNotFound         = 31001
	CodeAcceptancePending          = 31002
	CodeAcceptanceAlreadySubmitted = 31003
	CodeAcceptanceStatusInvalid    = 31004
	CodeAcceptanceNotPublisher     = 31005
	CodeAcceptanceRequired         = 31006

	CodeBlessingFailed     = 40001
	CodeCapsuleNotFound    = 50001
	CodeCapsuleLocked      = 50002
	CodeBadgeNotGranted    = 60001
	CodeAuditDenied        = 70001
	CodeUploadFailed       = 80001
	CodeFileTypeNotAllowed = 80002
)
