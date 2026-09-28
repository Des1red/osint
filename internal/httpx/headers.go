package httpx

import "net/http"

func browserHeaders() http.Header {
	headers :=
		make(
			http.Header,
		)

	headers.Set(
		"User-Agent",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36",
	)

	headers.Set(
		"Accept",
		"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
	)

	headers.Set(
		"Accept-Language",
		"en-US,en;q=0.9",
	)

	headers.Set(
		"Upgrade-Insecure-Requests",
		"1",
	)

	headers.Set(
		"Sec-Fetch-Dest",
		"document",
	)

	headers.Set(
		"Sec-Fetch-Mode",
		"navigate",
	)

	headers.Set(
		"Sec-Fetch-Site",
		"none",
	)

	headers.Set(
		"Sec-Fetch-User",
		"?1",
	)

	headers.Set(
		"Sec-CH-UA",
		`"Chromium";v="140", "Not=A?Brand";v="24"`,
	)

	headers.Set(
		"Sec-CH-UA-Mobile",
		"?0",
	)

	headers.Set(
		"Sec-CH-UA-Platform",
		`"Linux"`,
	)

	return headers
}

func mergeHeaders(
	defaults http.Header,
	custom http.Header,
) http.Header {
	final :=
		make(
			http.Header,
		)

	for key, values := range defaults {

		for _, value := range values {

			final.Add(
				key,
				value,
			)
		}
	}

	for key, values := range custom {

		final.Del(
			key,
		)

		for _, value := range values {

			final.Add(
				key,
				value,
			)
		}
	}

	return final
}
