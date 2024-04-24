package flexhttp

import (
	"bytes"
	"encoding/json"
	http "github.com/bogdanfinn/fhttp"
	"io"
	"reflect"
)

var useLog bool

func SetLog(l bool) {
	useLog = l
}

func NewFlexHttpClient(config *Config) (*FlexHttp, error) {
	cleanHttp, err := NewCleanHttpClient(config)
	if err != nil {
		return nil, err
	}
	client := &FlexHttp{
		CleanHttp: cleanHttp,
	}
	return client, nil
}

func (f *FlexHttp) Do(req *RQ) (*Response, error) {
	var bodyReader io.Reader
	if req.Body != nil {
		switch reflect.TypeOf(req.Body).Kind() {
		case reflect.String:
			bodyReader = bytes.NewBufferString(req.Body.(string))
			break
		case reflect.Slice:
			bodyBytes, ok := req.Body.([]byte)
			if ok {
				bodyReader = bytes.NewBuffer(bodyBytes)
			}
			break
		default:
			bodyBytes, err := json.Marshal(req.Body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewBuffer(bodyBytes)
		}
	}

	header := http.Header{}
	if req.Headers != nil {
		for k, v := range req.Headers {
			header.Set(k, v[0])
		}
	}
	if req.XHeaders != nil {
		if req.Headers == nil {
			for k, v := range f.CleanHttp.DefaultRequest.Header {
				header.Set(k, v[0])
			}
		}
		for k, v := range req.XHeaders {
			header.Set(k, v[0])
		}
	}

	response, err := f.CleanHttp.Do(RequestOption{
		Method: req.Method,
		Body:   bodyReader,
		Url:    req.URL,
		Header: header,
	})

	if err != nil {
		return nil, err
	}

	var resBody Body
	if response.Body != nil {
		resBody, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
	}

	return &Response{
		Body:   resBody,
		Status: response.StatusCode,
	}, nil
}

func (f *FlexHttp) DoJSON(req *RQ, to any) error {
	res, err := f.Do(req)
	if err != nil {
		return err
	}

	return res.Body.Unmarshal(to)
}

func (fb Body) Unmarshal(to any) error {
	return json.Unmarshal(fb, to)
}
