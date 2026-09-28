package httpx

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultTimeout = 10 * time.Second

	DefaultMaxRedirects = 10

	DefaultMaxBodyBytes = 2 * 1024 * 1024
)

var transport = func() *http.Transport {
	base :=
		http.DefaultTransport.(*http.Transport).
			Clone()

	base.MaxIdleConns =
		100

	base.MaxIdleConnsPerHost =
		20

	base.IdleConnTimeout =
		90 * time.Second

	base.ForceAttemptHTTP2 =
		true

	return base
}()

type Request struct {
	Method string
	URL    string

	Headers http.Header

	Body []byte

	Browser bool

	Timeout time.Duration

	MaxRedirects int
	MaxBodyBytes int64
}

type Response struct {
	StatusCode int

	Body []byte

	FinalURL string

	Headers http.Header

	Truncated bool
}

func Get(
	requestURL string,
) (Response, error) {
	return Do(
		Request{
			Method: http.MethodGet,

			URL: requestURL,

			Browser: true,
		},
	)
}

func Head(
	requestURL string,
) (Response, error) {
	return Do(
		Request{
			Method: http.MethodHead,

			URL: requestURL,

			Browser: true,
		},
	)
}

func Do(
	options Request,
) (Response, error) {
	options =
		defaults(
			options,
		)

	if strings.TrimSpace(
		options.URL,
	) == "" {

		return Response{},
			fmt.Errorf(
				"request URL is empty",
			)
	}

	request, err :=
		http.NewRequest(
			options.Method,
			options.URL,
			bytes.NewReader(
				options.Body,
			),
		)

	if err != nil {
		return Response{},
			err
	}

	if options.Browser {

		request.Header =
			mergeHeaders(
				browserHeaders(),
				options.Headers,
			)

	} else {

		request.Header =
			options.Headers.Clone()
	}

	client :=
		&http.Client{
			Transport: transport,

			Timeout: options.Timeout,

			CheckRedirect: func(
				req *http.Request,
				via []*http.Request,
			) error {
				if len(via) >=
					options.MaxRedirects {

					return fmt.Errorf(
						"too many redirects",
					)
				}

				return nil
			},
		}

	response, err :=
		client.Do(
			request,
		)

	if err != nil {
		return Response{},
			err
	}

	defer response.Body.Close()

	result :=
		Response{
			StatusCode: response.StatusCode,

			FinalURL: response.Request.URL.String(),

			Headers: response.Header.Clone(),
		}

	if options.Method ==
		http.MethodHead {

		return result,
			nil
	}

	body, err :=
		io.ReadAll(
			io.LimitReader(
				response.Body,
				options.MaxBodyBytes+1,
			),
		)

	if err != nil {
		return Response{},
			err
	}

	if int64(
		len(body),
	) > options.MaxBodyBytes {

		result.Truncated =
			true

		body =
			body[:options.MaxBodyBytes]
	}

	result.Body =
		body

	return result,
		nil
}

func defaults(
	options Request,
) Request {
	if strings.TrimSpace(
		options.Method,
	) == "" {

		options.Method =
			http.MethodGet
	}

	if options.Headers == nil {

		options.Headers =
			make(
				http.Header,
			)
	}

	if options.Timeout <= 0 {

		options.Timeout =
			DefaultTimeout
	}

	if options.MaxRedirects <= 0 {

		options.MaxRedirects =
			DefaultMaxRedirects
	}

	if options.MaxBodyBytes <= 0 {

		options.MaxBodyBytes =
			DefaultMaxBodyBytes
	}

	return options
}
