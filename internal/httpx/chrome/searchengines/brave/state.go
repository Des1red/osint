package brave

import (
	"context"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

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
