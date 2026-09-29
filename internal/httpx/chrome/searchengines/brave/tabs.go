package brave

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

func runBraveSearchMany(
	browserContext context.Context,
	queries []string,
) []pageResponse {
	responses :=
		make(
			[]pageResponse,
			len(queries),
		)

	//
	// Create every Brave tab up front.
	//
	// IMPORTANT:
	//
	// None of these tabs are closed until EVERY
	// Brave query has completed.
	//
	tabContexts :=
		make(
			[]context.Context,
			len(queries),
		)

	tabCancels :=
		make(
			[]context.CancelFunc,
			len(queries),
		)

	for index := range queries {

		tabContexts[index],
			tabCancels[index] =
			chromedp.NewContext(
				browserContext,
			)
	}

	//
	// Close all Brave tabs together only after the
	// complete Brave batch has finished.
	//
	defer func() {

		for _, cancel := range tabCancels {

			if cancel != nil {

				cancel()
			}
		}
	}()

	semaphore :=
		make(
			chan struct{},
			DefaultSearchConcurrency,
		)

	launchDelays :=
		searchLaunchDelays(
			len(queries),
			humanSearchDelay,
		)

	var waitGroup sync.WaitGroup

	for index, query := range queries {

		waitGroup.Add(
			1,
		)

		go func(
			index int,
			query string,
			tabContext context.Context,
		) {
			defer waitGroup.Done()

			delay :=
				launchDelays[index]

			if delay > 0 {

				timer :=
					time.NewTimer(
						delay,
					)

				defer timer.Stop()

				select {

				case <-browserContext.Done():

					responses[index] =
						pageResponse{
							Query: query,

							Err: browserContext.Err(),
						}

					return

				case <-timer.C:
				}
			}

			select {

			case semaphore <- struct{}{}:

			case <-browserContext.Done():

				responses[index] =
					pageResponse{
						Query: query,

						Err: browserContext.Err(),
					}

				return
			}

			defer func() {

				<-semaphore
			}()

			responses[index] =
				braveSearchTab(
					tabContext,
					query,
				)
		}(
			index,
			query,
			tabContexts[index],
		)
	}

	waitGroup.Wait()

	return responses
}

func braveSearchTab(
	tabContext context.Context,
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
	// The Brave tab was created by
	// runBraveSearchMany().
	//
	// Do NOT close it here.
	//
	// It remains alive until the entire Brave
	// search batch has completed.
	//

	timeoutContext,
		timeoutCancel :=
		context.WithTimeout(
			tabContext,
			braveVerificationTimeout,
		)

	defer timeoutCancel()

	searchURL :=
		"https://search.brave.com/search?q=" +
			url.QueryEscape(
				query,
			) +
			"&source=web"

	err :=
		navigateBraveSearch(
			timeoutContext,
			searchURL,
		)

	if err != nil {

		result.Err =
			err

		return result
	}

	state,
		err :=
		bravePageState(
			timeoutContext,
		)

	if err != nil {

		result.Err =
			err

		return result
	}

	switch state {

	case "results":

		//
		// Normal fast path.
		//

	case "verify",
		"verifying":

		err =
			resolveBraveVerification(
				timeoutContext,
			)

		if err != nil {

			result.Err =
				err

			return result
		}

	default:

		result.Err =
			fmt.Errorf(
				"Brave reached an unknown page state",
			)

		return result
	}

	//
	// Do not capture or close this tab merely
	// because result nodes appeared once.
	//
	// Brave can briefly keep old/result DOM under
	// its verification UI. The page must remain
	// in a clean result state for a short period.
	//
	err =
		waitForBraveStableResults(
			timeoutContext,
		)

	if err != nil {

		result.Err =
			err

		return result
	}

	//
	// Give lazy-rendered result content a chance
	// to populate before capturing the DOM.
	//
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
				700*
					time.Millisecond,
			),
		)

	if err != nil {

		result.Err =
			err

		return result
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

func navigateBraveSearch(
	tabContext context.Context,
	searchURL string,
) error {
	err :=
		chromedp.Run(
			tabContext,

			chromedp.Navigate(
				searchURL,
			),

			chromedp.WaitReady(
				"body",
				chromedp.ByQuery,
			),
		)

	if err != nil {

		return err
	}

	return waitForBraveState(
		tabContext,
	)
}
