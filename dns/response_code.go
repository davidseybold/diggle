package dns

type ResponseCode uint8

const (
	ResponseCodeNoError          ResponseCode = 0
	ReponseCodeFormError         ResponseCode = 1
	ResponseCodeServerFailure    ResponseCode = 2
	ResponseCodeNXDomain         ResponseCode = 3
	ResponseCodeNotImplemented   ResponseCode = 4
	ReponseCodeRefused           ResponseCode = 5
	ResponseCodeYXDomain         ResponseCode = 6
	ResponseCodeYXRRSet          ResponseCode = 7
	ResponseCodeNXRRSet          ResponseCode = 8
	ResponseCodeNotAuthoritative ResponseCode = 9
	ResponseCodeNotAuthorized    ResponseCode = 9
	ResponseCodeNotZone          ResponseCode = 10
	ResponseCodeDSOTypeNI        ResponseCode = 11
	ResponseCodeBadOptVersion    ResponseCode = 16
	ResponseCodeBadSignature     ResponseCode = 16
	ResponseCodeBadKey           ResponseCode = 17
	ResponseCodeBadTime          ResponseCode = 18
	ResponseCodeCodeBadMode      ResponseCode = 19
	ResponseCodeBadName          ResponseCode = 20
	ResponseCodeBadAlgorithm     ResponseCode = 21
	ResponseCodeBadTruncation    ResponseCode = 22
	ResponseCodeBadCookie        ResponseCode = 23
)
