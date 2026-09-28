package searchengines

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"

	"osint/internal/httpx/chrome"
)

const duckDuckGoMaxSearchResults = 40

const duckDuckGoMaxLoadRounds = 8

const duckDuckGoLoadTimeout = 2 * time.Second

func duckDuckGo(
	browserContext context.Context,
	queries []string,
) (
	[]Result,
	error,
) {
	responses :=
		runSearchMany(
			browserContext,
			queries,
			duckDuckGoSearchTab,
			duckDuckGoSearchDelay,
			DuckDuckGoSearchConcurrency,
		)

	var result []Result

	successfulQueries :=
		0

	var firstError error

	for _, response := range responses {

		if response.Err != nil {

			if firstError == nil {

				firstError =
					response.Err
			}

			continue
		}

		items,
			err :=
			parseDuckDuckGoResults(
				response.HTML,
			)

		if err != nil {

			if firstError == nil {

				firstError =
					err
			}

			continue
		}

		successfulQueries++

		for index := range items {

			items[index].Query =
				response.Query
		}

		result =
			append(
				result,
				items...,
			)
	}

	if successfulQueries == 0 &&
		firstError != nil {

		return nil,
			fmt.Errorf(
				"DuckDuckGo searches failed: %w",
				firstError,
			)
	}

	return result,
		nil
}

func duckDuckGoSearchTab(
	browserContext context.Context,
	query string,
) pageResponse {
	result :=
		pageResponse{
			Query: query,
		}

	query =
		strings.TrimSpace(
			query,
		)

	if query == "" {

		result.Err =
			fmt.Errorf(
				"search query is empty",
			)

		return result
	}

	//
	// Every query receives its own temporary
	// tab inside the shared Chromium process.
	//
	tabContext,
		tabCancel :=
		chromedp.NewContext(
			browserContext,
		)

	defer tabCancel()

	timeoutContext,
		timeoutCancel :=
		context.WithTimeout(
			tabContext,
			chrome.DefaultTimeout,
		)

	defer timeoutCancel()

	searchURL :=
		"https://duckduckgo.com/?q=" +
			url.QueryEscape(
				query,
			)

	err :=
		chromedp.Run(
			timeoutContext,

			chromedp.Navigate(
				searchURL,
			),

			chromedp.WaitReady(
				"body",
				chromedp.ByQuery,
			),

			chromedp.WaitVisible(
				`article[data-testid="result"],
				 div[data-testid="result"],
				 div.result,
				 div.web-result`,
				chromedp.ByQuery,
			),
		)

	if err != nil {

		result.Err =
			err

		return result
	}

	for round := 0; round <
		duckDuckGoMaxLoadRounds; round++ {

		currentCount :=
			0

		err =
			chromedp.Run(
				timeoutContext,

				chromedp.Evaluate(
					`
					document.querySelectorAll(
						'article[data-testid="result"], ' +
						'div[data-testid="result"], ' +
						'div.result, ' +
						'div.web-result'
					).length
					`,
					&currentCount,
				),
			)

		if err != nil {

			break
		}

		if currentCount >=
			duckDuckGoMaxSearchResults {

			break
		}

		var clickedMore bool

		err =
			chromedp.Run(
				timeoutContext,

				chromedp.Evaluate(
					`
					window.scrollTo(
						0,
						document.body.scrollHeight
					)
					`,
					nil,
				),

				chromedp.Sleep(
					700*time.Millisecond,
				),

				chromedp.Evaluate(
					`
					(() => {
						let more =
							document.querySelector(
								'#more-results'
							);

						if (!more) {
							more =
								document.querySelector(
									'[data-testid="more-results"]'
								);
						}

						if (!more) {
							const controls =
								Array.from(
									document.querySelectorAll(
										'button, a'
									)
								);

							more =
								controls.find(
									control =>
										/more results/i.test(
											(
												control.innerText ||
												control.textContent ||
												''
											).trim()
										)
								);
						}

						if (!more) {
							return false;
						}

						more.scrollIntoView({
							block: 'center'
						});

						more.click();

						return true;
					})()
					`,
					&clickedMore,
				),
			)

		if err != nil {

			break
		}

		if !clickedMore {

			break
		}

		//
		// Wait only as long as DuckDuckGo
		// actually needs to append more results.
		//
		// Previously this always slept for 1.5s,
		// even when results appeared immediately.
		//
		loadContext,
			loadCancel :=
			context.WithTimeout(
				timeoutContext,
				duckDuckGoLoadTimeout,
			)

		countExpression :=
			fmt.Sprintf(
				`
				document.querySelectorAll(
					'article[data-testid="result"], ' +
					'div[data-testid="result"], ' +
					'div.result, ' +
					'div.web-result'
				).length > %d
				`,
				currentCount,
			)

		err =
			chromedp.Run(
				loadContext,

				chromedp.Poll(
					countExpression,
					nil,

					chromedp.WithPollingInterval(
						100*
							time.Millisecond,
					),
				),
			)

		loadCancel()

		if err != nil {

			//
			// Nothing else appeared within the
			// small expansion window.
			//
			// This is not a failed search. It just
			// means there are no more useful
			// results to load.
			//
			break
		}
	}

	var html string

	var finalURL string

	err =
		chromedp.Run(
			timeoutContext,

			chromedp.Evaluate(
				`
				window.scrollTo(
					0,
					0
				)
				`,
				nil,
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

		result.Err =
			err

		return result
	}

	result.HTML =
		html

	result.FinalURL =
		finalURL

	return result
}
