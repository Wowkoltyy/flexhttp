package flexhttp

import (
	"github.com/Wowkoltyy/harparser"
	"github.com/bogdanfinn/tls-client/profiles"
	"io"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

type CleanHttp struct {
	Config         *Config
	Client         tls_client.HttpClient
	DefaultRequest *harparser.RequestInfo
	Log            bool
}

type Body []byte

type Response struct {
	Body   Body
	Status int
}

type RQ struct {
	Method   string
	Body     any
	URL      string
	SkipBody bool
	Headers  http.Header
	XHeaders http.Header

	RequestInfo    harparser.RequestInfo
	XHeaderManager func(http.Header) http.Header
}

type FlexHttp struct {
	CleanHttp *CleanHttp
}

type Config struct {
	Proxy          string
	Timeout        int
	Log            bool
	Profile        *profiles.ClientProfile
	DefaultRequest *harparser.RequestInfo

	ReadTimeout, WriteTimeout, MaxIdleConnDuration time.Duration
}

type RequestOption struct {
	Method string
	Body   io.Reader
	Url    string
	Header http.Header
}

type UserAgentInfo struct {
	BrowserName    string
	BrowserVersion string
	OSName         string
	OSVersion      string
	UaVersion      string
}
