package searchengines

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

const braveVerificationTimeout = 2 * time.Minute

// Only one Brave tab may actively perform
// verification at a time.
//
// Other Brave workers keep their own tabs open
// and wait here. Once they acquire the gate they
// inspect their existing page again before doing
// anything.
var braveVerificationGate = make(
	chan struct{},
	1,
)

func brave(
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
			braveSearchTab,
			humanSearchDelay,
			DefaultSearchConcurrency,
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
			parseBraveResults(
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
				"Brave searches failed: %w",
				firstError,
			)
	}

	return result,
		nil
}

func braveSearchTab(
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
	// Every Brave query gets its own tab.
	//
	// The tab remains alive for the complete
	// query lifecycle, including verification,
	// stable result rendering, and HTML capture.
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

	//
	// tabCancel() runs only after the response has
	// been fully captured.
	//
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

func resolveBraveVerification(
	tabContext context.Context,
) error {
	err :=
		acquireBraveVerification(
			tabContext,
		)

	if err != nil {

		return err
	}

	defer releaseBraveVerification()

	//
	// This tab has remained open while waiting
	// for the verification gate.
	//
	// Another worker may have completed a shared
	// Brave verification while we were waiting,
	// so inspect THIS SAME TAB before clicking.
	//
	state,
		err :=
		bravePageState(
			tabContext,
		)

	if err != nil {

		return err
	}

	switch state {

	case "results":

		return waitForBraveStableResults(
			tabContext,
		)

	case "verifying":

		//
		// Brave is already transitioning from a
		// successful verification into results.
		// Keep the tab alive and wait it out.
		//
		return waitForBraveStableResults(
			tabContext,
		)

	case "verify":

		//
		// This specific tab still requires
		// verification. Verify it in place.
		//

	default:

		//
		// The page may be in a tiny transition
		// window. Wait until it becomes meaningful
		// again, then inspect it once more.
		//
		err =
			waitForBraveState(
				tabContext,
			)

		if err != nil {

			return err
		}

		state,
			err =
			bravePageState(
				tabContext,
			)

		if err != nil {

			return err
		}

		if state == "results" ||
			state == "verifying" {

			return waitForBraveStableResults(
				tabContext,
			)
		}

		if state != "verify" {

			return fmt.Errorf(
				"Brave reached an unknown verification state",
			)
		}
	}

	err =
		verify(
			tabContext,
		)

	if err != nil {

		return err
	}

	//
	// IMPORTANT:
	//
	// Do not return when the Verify button merely
	// disappears.
	//
	// Brave briefly shows a "Verified" state and
	// may already have old/result DOM underneath.
	// Keep this exact tab open until the normal
	// result page has remained clean and stable.
	//
	return waitForBraveStableResults(
		tabContext,
	)
}

func acquireBraveVerification(
	ctx context.Context,
) error {
	select {

	case braveVerificationGate <- struct{}{}:

		return nil

	case <-ctx.Done():

		return ctx.Err()
	}
}

func releaseBraveVerification() {
	select {

	case <-braveVerificationGate:

	default:
	}
}

func waitForBraveState(
	tabContext context.Context,
) error {
	return chromedp.Run(
		tabContext,

		chromedp.Poll(
			`
			(() => {
				const visible = element => {
					if (!element) {
						return false;
					}

					const style =
						window.getComputedStyle(
							element
						);

					return (
						style.display !== 'none' &&
						style.visibility !== 'hidden' &&
						element.getClientRects().length > 0
					);
				};

				const exactVerify =
					document.evaluate(
						'/html/body/div/div[1]/div/div[3]/div[2]/div/div/button',
						document,
						null,
						XPathResult.FIRST_ORDERED_NODE_TYPE,
						null
					).singleNodeValue;

				if (visible(exactVerify)) {
					return true;
				}

				const textVerify =
					Array.from(
						document.querySelectorAll(
							'button'
						)
					).find(
						button =>
							visible(button) &&
							(
								button.innerText ||
								button.textContent ||
								''
							)
								.trim()
								.toLowerCase() ===
								'verify'
					);

				if (textVerify) {
					return true;
				}

				const verified =
					Array.from(
						document.querySelectorAll(
							'div, span, p, button'
						)
					).find(
						element =>
							visible(element) &&
							(
								element.innerText ||
								element.textContent ||
								''
							)
								.trim()
								.toLowerCase() ===
								'verified'
					);

				if (verified) {
					return true;
				}

				return (
					document.querySelector(
						'div.snippet[data-type="web"]'
					) !== null ||
					document.querySelector(
						'[data-type="web"][data-pos]'
					) !== null
				);
			})()
			`,
			nil,

			chromedp.WithPollingInterval(
				250*
					time.Millisecond,
			),
		),
	)
}

func bravePageState(
	tabContext context.Context,
) (
	string,
	error,
) {
	var state string

	err :=
		chromedp.Run(
			tabContext,

			chromedp.Evaluate(
				`
				(() => {
					const visible = element => {
						if (!element) {
							return false;
						}

						const style =
							window.getComputedStyle(
								element
							);

						return (
							style.display !== 'none' &&
							style.visibility !== 'hidden' &&
							element.getClientRects().length > 0
						);
					};

					const exactVerify =
						document.evaluate(
							'/html/body/div/div[1]/div/div[3]/div[2]/div/div/button',
							document,
							null,
							XPathResult.FIRST_ORDERED_NODE_TYPE,
							null
						).singleNodeValue;

					if (visible(exactVerify)) {
						return 'verify';
					}

					const textVerify =
						Array.from(
							document.querySelectorAll(
								'button'
							)
						).find(
							button =>
								visible(button) &&
								(
									button.innerText ||
									button.textContent ||
									''
								)
									.trim()
									.toLowerCase() ===
									'verify'
						);

					if (textVerify) {
						return 'verify';
					}

					const verified =
						Array.from(
							document.querySelectorAll(
								'div, span, p, button'
							)
						).find(
							element =>
								visible(element) &&
								(
									element.innerText ||
									element.textContent ||
									''
								)
									.trim()
									.toLowerCase() ===
									'verified'
						);

					if (verified) {
						return 'verifying';
					}

					const hasResults =
						document.querySelector(
							'div.snippet[data-type="web"]'
						) !== null ||
						document.querySelector(
							'[data-type="web"][data-pos]'
						) !== null;

					if (hasResults) {
						return 'results';
					}

					return '';
				})()
				`,
				&state,
			),
		)

	if err != nil {

		return "",
			err
	}

	return strings.TrimSpace(
			state,
		),
		nil
}

func waitForBraveStableResults(
	tabContext context.Context,
) error {
	return chromedp.Run(
		tabContext,

		chromedp.Poll(
			`
			(() => {
				const visible = element => {
					if (!element) {
						return false;
					}

					const style =
						window.getComputedStyle(
							element
						);

					return (
						style.display !== 'none' &&
						style.visibility !== 'hidden' &&
						element.getClientRects().length > 0
					);
				};

				const exactVerify =
					document.evaluate(
						'/html/body/div/div[1]/div/div[3]/div[2]/div/div/button',
						document,
						null,
						XPathResult.FIRST_ORDERED_NODE_TYPE,
						null
					).singleNodeValue;

				const textVerify =
					Array.from(
						document.querySelectorAll(
							'button'
						)
					).find(
						button =>
							visible(button) &&
							(
								button.innerText ||
								button.textContent ||
								''
							)
								.trim()
								.toLowerCase() ===
								'verify'
					);

				const verified =
					Array.from(
						document.querySelectorAll(
							'div, span, p, button'
						)
					).find(
						element =>
							visible(element) &&
							(
								element.innerText ||
								element.textContent ||
								''
							)
								.trim()
								.toLowerCase() ===
								'verified'
					);

				const hasResults =
					document.querySelector(
						'div.snippet[data-type="web"]'
					) !== null ||
					document.querySelector(
						'[data-type="web"][data-pos]'
					) !== null;

				const cleanResults =
					hasResults &&
					!visible(exactVerify) &&
					!textVerify &&
					!verified;

				if (!cleanResults) {
					window.__osintBraveStableSince = 0;

					return false;
				}

				const now =
					Date.now();

				if (!window.__osintBraveStableSince) {
					window.__osintBraveStableSince =
						now;

					return false;
				}

				return (
					now -
					window.__osintBraveStableSince
				) >= 1500;
			})()
			`,
			nil,

			chromedp.WithPollingInterval(
				100*
					time.Millisecond,
			),
		),
	)
}

func verify(
	tabContext context.Context,
) error {
	const verifyXPath = `/html/body/div/div[1]/div/div[3]/div[2]/div/div/button`

	var clicked bool

	err :=
		chromedp.Run(
			tabContext,

			chromedp.WaitVisible(
				verifyXPath,
				chromedp.BySearch,
			),

			chromedp.Evaluate(
				`
				(() => {
					let button =
						document.evaluate(
							'/html/body/div/div[1]/div/div[3]/div[2]/div/div/button',
							document,
							null,
							XPathResult.FIRST_ORDERED_NODE_TYPE,
							null
						).singleNodeValue;

					if (!button) {
						button =
							Array.from(
								document.querySelectorAll(
									'button'
								)
							).find(
								item =>
									(
										item.innerText ||
										item.textContent ||
										''
									)
										.trim()
										.toLowerCase() ===
										'verify'
							);
					}

					if (!button) {
						return false;
					}

					button.scrollIntoView({
						block: 'center',
						inline: 'center'
					});

					button.click();

					return true;
				})()
				`,
				&clicked,
			),
		)

	if err != nil {

		return err
	}

	if !clicked {

		return fmt.Errorf(
			"Brave verification button could not be clicked",
		)
	}

	return nil
}
