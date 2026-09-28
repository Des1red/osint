package chrome

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"osint/internal/models"

	"github.com/chromedp/chromedp"
)

var browserMu sync.RWMutex

var browserContext context.Context

var browserCancel context.CancelFunc

func Open() error {
	browserMu.Lock()

	defer browserMu.Unlock()

	//
	// Already open.
	//
	if browserContext != nil {

		select {

		case <-browserContext.Done():

			//
			// The browser died or was cancelled.
			// Clear the old session and reopen it.
			//
			if browserCancel != nil {

				browserCancel()
			}

			browserContext =
				nil

			browserCancel =
				nil

		default:

			return nil
		}
	}

	contextValue,
		cancel,
		err :=
		newBrowser()

	if err != nil {
		return err
	}

	browserContext =
		contextValue

	browserCancel =
		cancel

	return nil
}

func Context() (
	context.Context,
	error,
) {
	browserMu.RLock()

	current :=
		browserContext

	browserMu.RUnlock()

	if current == nil {

		return nil,
			fmt.Errorf(
				"browser session is not open",
			)
	}

	select {

	case <-current.Done():

		return nil,
			fmt.Errorf(
				"browser session is closed",
			)

	default:
	}

	return current,
		nil
}

func Close() {
	browserMu.Lock()

	cancel :=
		browserCancel

	browserContext =
		nil

	browserCancel =
		nil

	browserMu.Unlock()

	if cancel != nil {

		cancel()
	}
}

func newBrowser() (
	context.Context,
	context.CancelFunc,
	error,
) {
	executable :=
		strings.TrimSpace(
			models.Browser.Executable,
		)

	profileDirectory :=
		strings.TrimSpace(
			models.Browser.ProfileDirectory,
		)

	if executable == "" {

		return nil,
			nil,
			fmt.Errorf(
				"browser executable was not initialized by bootstrap",
			)
	}

	if profileDirectory == "" {

		return nil,
			nil,
			fmt.Errorf(
				"browser profile was not initialized by bootstrap",
			)
	}

	options :=
		append(
			chromedp.DefaultExecAllocatorOptions[:],

			chromedp.ExecPath(
				executable,
			),

			//
			// Persistent browser profile.
			//
			// The directory is prepared by
			// bootstrap before the browser layer
			// is used.
			//
			chromedp.UserDataDir(
				profileDirectory,
			),

			//
			// Visible Chromium window.
			//
			chromedp.Flag(
				"headless",
				false,
			),

			chromedp.Flag(
				"disable-gpu",
				true,
			),

			chromedp.Flag(
				"disable-dev-shm-usage",
				true,
			),

			chromedp.Flag(
				"no-first-run",
				true,
			),

			chromedp.Flag(
				"no-default-browser-check",
				true,
			),

			chromedp.Flag(
				"disable-blink-features",
				"AutomationControlled",
			),
		)

	allocatorContext,
		allocatorCancel :=
		chromedp.NewExecAllocator(
			context.Background(),
			options...,
		)

	contextValue,
		contextCancel :=
		chromedp.NewContext(
			allocatorContext,
		)

	//
	// Force Chromium to start once here.
	//
	err :=
		chromedp.Run(
			contextValue,

			chromedp.Navigate(
				"about:blank",
			),
		)

	if err != nil {

		contextCancel()

		allocatorCancel()

		return nil,
			nil,
			err
	}

	cancel :=
		func() {

			contextCancel()

			allocatorCancel()
		}

	return contextValue,
		cancel,
		nil
}
