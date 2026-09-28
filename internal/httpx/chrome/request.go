package chrome

import (
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/chromedp"
)

func Get(
	requestURL string,
) (
	Response,
	error,
) {
	return Do(
		Request{
			URL: requestURL,
		},
	)
}

func Do(
	request Request,
) (
	Response,
	error,
) {
	request =
		defaults(
			request,
		)

	request.URL =
		strings.TrimSpace(
			request.URL,
		)

	if request.URL == "" {

		return Response{},
			fmt.Errorf(
				"request URL is empty",
			)
	}

	browserContext,
		err :=
		Context()

	if err != nil {

		return Response{},
			err
	}

	//
	// Every request gets its own temporary tab.
	//
	// Closing this context closes only the tab,
	// not the shared Chromium process.
	//
	tabContext,
		tabCancel :=
		chromedp.NewContext(
			browserContext,
		)

	defer tabCancel()

	requestContext,
		requestCancel :=
		context.WithTimeout(
			tabContext,
			request.Timeout,
		)

	defer requestCancel()

	var html string

	var finalURL string

	err =
		chromedp.Run(
			requestContext,

			chromedp.Navigate(
				request.URL,
			),

			chromedp.WaitReady(
				"body",
				chromedp.ByQuery,
			),

			chromedp.OuterHTML(
				"html",
				&html,
				chromedp.ByQuery,
			),

			chromedp.Location(
				&finalURL,
			),
		)

	if err != nil {

		return Response{},
			err
	}

	return Response{
			HTML: html,

			FinalURL: finalURL,
		},
		nil
}

func defaults(
	request Request,
) Request {
	if request.Timeout <= 0 {

		request.Timeout =
			DefaultTimeout
	}

	return request
}
