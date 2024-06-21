package httpcore

import (
	"errors"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

func Req(method string, URL string, header map[string]string, body interface{}, from string) (resp *resty.Response, err error) {
	// intialise header
	client := resty.New()
	builder := client.R()
	// todo iterate header
	for headerName, headerValue := range header {
		builder.SetHeader(headerName, headerValue)
	}
	// query parameter
	// form data
	// body
	switch body := body.(type) {
	case string:
		builder.SetBody(body)
	case map[string]interface{}:
		// reqBodyData := url.Values{}
		// for reqName, reqValue := range body {
		// 	reqBodyData.Set(reqName, reqValue.(string))
		// }
		// encodedData := reqBodyData.Encode()
		// builder.SetBody(encodedData)
		builder.SetBody(body)
	case []byte:
		builder.SetBody(body)
	case map[string]string:
		reqBodyData := url.Values{}
		for reqName, reqValue := range body {
			reqBodyData.Set(reqName, reqValue)
		}
		encodedData := reqBodyData.Encode()
		builder.SetBody(strings.NewReader(encodedData))
	default:
		if body == nil {
			break
		}
		builder.SetBody(body)
	}
	// todo validate url
	// resp, err := builder.Execute(resty.MethodPost, URL)
	switch method {
	case "GET":
		resp, err = builder.Execute(resty.MethodGet, URL)
		return
	case "POST":
		resp, err = builder.Execute(resty.MethodPost, URL)
		return
	case "PUT":
		resp, err = builder.Execute(resty.MethodPut, URL)
		return
	case "DELETE":
		resp, err = builder.Execute(resty.MethodDelete, URL)
		return
	default:
		return nil, errors.New("provided http method doesn't exists")
	}
}
