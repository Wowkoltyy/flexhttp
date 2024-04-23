package flexhttp

import (
	"errors"
	"github.com/bogdanfinn/tls-client/profiles"
	"log"
	"net/url"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

func GetChromeProfile(version int) profiles.ClientProfile {
	switch version {
	case 120:
		return profiles.Chrome_120
	case 117:
		return profiles.Chrome_117
	case 112:
		return profiles.Chrome_112
	case 111:
		return profiles.Chrome_111
	case 110:
		return profiles.Chrome_110
	case 109:
		return profiles.Chrome_109
	case 108:
		return profiles.Chrome_108
	case 107:
		return profiles.Chrome_107
	case 106:
		return profiles.Chrome_106
	case 105:
		return profiles.Chrome_105
	case 104:
		return profiles.Chrome_104
	case 103:
		return profiles.Chrome_103
	default:
		return profiles.Chrome_120 // default profile
	}
}

// create http client and return *CleanHttp. Take *Config as params.
func NewCleanHttpClient(config *Config) (*CleanHttp, error) {
	if config.Timeout < 30 {
		config.Timeout = 30
	}

	if config.Profile == nil {
		chromeProfile := GetChromeProfile(config.DefaultRequest.ChromeVersion)
		config.Profile = &chromeProfile
	}

	options := []tls_client.HttpClientOption{
		tls_client.WithClientProfile(*config.Profile),
		tls_client.WithInsecureSkipVerify(),
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
		tls_client.WithRandomTLSExtensionOrder(),
	}

	if config.Proxy != "" {
		options = append(options, tls_client.WithProxyUrl(config.Proxy))
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	c := CleanHttp{
		Config:         config,
		Client:         client,
		Log:            config.Log,
		DefaultRequest: config.DefaultRequest,
	}

	return &c, nil
}

func (c *CleanHttp) SetCookies(url *url.URL, cookies []*http.Cookie) {
	c.Client.SetCookies(url, cookies)
}

func (c *CleanHttp) Do(request RequestOption) (*http.Response, error) {
	if request.Url == "" {
		return nil, errors.New("please provide valid url")
	}

	if request.Header == nil || len(request.Header) < 1 {
		request.Header = c.DefaultRequest.Header
	}

	req, err := http.NewRequest(request.Method, request.Url, request.Body)
	if err != nil {
		return nil, err
	}

	for k, v := range request.Header {
		req.Header.Set(k, v[0])
	}

	u, _ := url.Parse(request.Url)
	for _, cook := range c.Client.GetCookieJar().Cookies(u) {
		req.AddCookie(cook)

		if c.Log {
			log.Println("add cookie", cook.Name, cook.Value, cook.Domain)
		}
	}

	if c.Log {
		log.Println(req.URL, req.Header, req.Body)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
