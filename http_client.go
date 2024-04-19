package main

import (
	"fmt"
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client
}

var (
	// DefaultHeader is a map containing default HTTP headers commonly used in web requests.
	// These headers mimic the behavior of popular web browsers to ensure compatibility
	// with web servers and services.
	// I take it from browser so you replace with yours or add what you like
	DefaultHeader = map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Encoding": "gzip, deflate, br, zstd",
		"ACCEPT-LANGUAGE": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,ar-MA;q=0.6,ar;q=0.5,es;q=0.4",
	}
)

func NewHttpClient() *HttpClient {
	return &HttpClient{
		client: &http.Client{
			// CheckRedirect is a custom function to handle redirects.
			// I increase the maximum number of redirects to 20 to accommodate
			// for the specific requirements of the target URL, which may involve
			// a larger number of redirects before reaching the final destination.
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 20 {
					return fmt.Errorf("stopped after 20 redirects")
				}
				return nil
			},
			// Some URLs are not trusted so due to an issues in TLS certificate verification
			// so you have to ignore certificate verification.
			// Transport: &http.Transport{
			// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			// },
		},
	}
}

// NewRequest creates a new HTTP request with the specified method, URL, and headers,
// and sends the request using the underlying HTTP client. It retries the request
// up to 10 times in case of network errors or server timeouts.
//
// Parameters:
//   - method: The HTTP method (GET, POST, PUT, DELETE, etc.) for the request.
//   - url: The URL to send the request to.
//   - headers: A map of HTTP headers to include in the request.
//
// Returns:
//   - (*http.Response): A pointer to the HTTP response if the request is successful.
//   - (error): An error if the request fails or exceeds the maximum number of retries.
func (c *HttpClient) NewRequest(method, url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	tries := 10
	for {
		resp, err := c.client.Do(req)
		if err != nil {
			if tries == 0 {
				return nil, err
			}
			fmt.Println("Error:", err.Error(), "trying again...")
			tries--
			time.Sleep(time.Second * 2)
			continue
		}
		return resp, nil
	}
}
