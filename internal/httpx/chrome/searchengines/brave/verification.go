package brave

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

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
