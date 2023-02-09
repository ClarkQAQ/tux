package tux

import (
	"errors"
	"net/textproto"
)

// HTTP headers
const (
	HeaderAccept                        = "Accept"
	HeaderAcceptCharset                 = "Accept-Charset"
	HeaderAcceptEncoding                = "Accept-Encoding"
	HeaderAcceptLanguage                = "Accept-Language"
	HeaderAuthorization                 = "Authorization"
	HeaderCacheControl                  = "Cache-Control"
	HeaderContentLength                 = "Content-Length"
	HeaderContentMD5                    = "Content-MD5"
	HeaderContentType                   = "Content-Type"
	HeaderDoNotTrack                    = "DNT"
	HeaderIfMatch                       = "If-Match"
	HeaderIfModifiedSince               = "If-Modified-Since"
	HeaderIfNoneMatch                   = "If-None-Match"
	HeaderIfRange                       = "If-Range"
	HeaderIfUnmodifiedSince             = "If-Unmodified-Since"
	HeaderMaxForwards                   = "Max-Forwards"
	HeaderProxyAuthorization            = "Proxy-Authorization"
	HeaderPragma                        = "Pragma"
	HeaderRange                         = "Range"
	HeaderReferer                       = "Referer"
	HeaderUserAgent                     = "User-Agent"
	HeaderTE                            = "TE"
	HeaderVia                           = "Via"
	HeaderWarning                       = "Warning"
	HeaderCookie                        = "Cookie"
	HeaderOrigin                        = "Origin"
	HeaderAcceptDatetime                = "Accept-Datetime"
	HeaderXRequestedWith                = "X-Requested-With"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAcceptPatch                   = "Accept-Patch"
	HeaderAcceptRanges                  = "Accept-Ranges"
	HeaderAllow                         = "Allow"
	HeaderContentEncoding               = "Content-Encoding"
	HeaderContentLanguage               = "Content-Language"
	HeaderContentLocation               = "Content-Location"
	HeaderContentDisposition            = "Content-Disposition"
	HeaderContentRange                  = "Content-Range"
	HeaderETag                          = "ETag"
	HeaderExpires                       = "Expires"
	HeaderLastModified                  = "Last-Modified"
	HeaderLink                          = "Link"
	HeaderLocation                      = "Location"
	HeaderP3P                           = "P3P"
	HeaderProxyAuthenticate             = "Proxy-Authenticate"
	HeaderRefresh                       = "Refresh"
	HeaderRetryAfter                    = "Retry-After"
	HeaderServer                        = "Server"
	HeaderSetCookie                     = "Set-Cookie"
	HeaderStrictTransportSecurity       = "Strict-Transport-Security"
	HeaderTransferEncoding              = "Transfer-Encoding"
	HeaderUpgrade                       = "Upgrade"
	HeaderVary                          = "Vary"
	HeaderWWWAuthenticate               = "WWW-Authenticate"

	// Non-Standard
	HeaderXFrameOptions          = "X-Frame-Options"
	HeaderXXSSProtection         = "X-XSS-Protection"
	HeaderContentSecurityPolicy  = "Content-Security-Policy"
	HeaderXContentSecurityPolicy = "X-Content-Security-Policy"
	HeaderXWebKitCSP             = "X-WebKit-CSP"
	HeaderXContentTypeOptions    = "X-Content-Type-Options"
	HeaderXPoweredBy             = "X-Powered-By"
	HeaderXUACompatible          = "X-UA-Compatible"
	HeaderXForwardedProto        = "X-Forwarded-Proto"
	HeaderXHTTPMethodOverride    = "X-HTTP-Method-Override"
	HeaderXForwardedFor          = "X-Forwarded-For"
	HeaderXRealIP                = "X-Real-IP"
	HeaderXCSRFToken             = "X-CSRF-Token"
	HeaderXRatelimitLimit        = "X-Ratelimit-Limit"
	HeaderXRatelimitRemaining    = "X-Ratelimit-Remaining"
	HeaderXRatelimitReset        = "X-Ratelimit-Reset"

	defaultStatusCode = 200
	StatusOK          = 200
)

var (
	ErrWriterAlreadyExported = errors.New("writer already exported")
	ErrResponseAlreadySent   = errors.New("response already sent")
)

// Normalize formats the input header to the formation of "Xxx-Xxx".
func Normalize(header string) string {
	return textproto.CanonicalMIMEHeaderKey(header)
}
