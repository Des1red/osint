package network

func applyIPAPI(
	result *NetworkResult,
	data ipAPIResponse,
) {
	if result.ISP == "" {
		result.ISP = data.ISP
	}

	if result.Organization == "" {
		result.Organization = data.Org
	}

	if data.Reverse != "" {
		result.ReverseDNS =
			appendUnique(
				result.ReverseDNS,
				data.Reverse,
			)
	}

	result.Mobile = FlagResult{
		Available: true,
		Value:     data.Mobile,
		Source:    "ip-api.com",
	}

	result.Proxy = FlagResult{
		Available: true,
		Value:     data.Proxy,
		Source:    "ip-api.com",
	}

	result.Hosting = FlagResult{
		Available: true,
		Value:     data.Hosting,
		Source:    "ip-api.com",
	}

	result.Sources = appendUnique(
		result.Sources,
		"ip-api.com",
	)
}

func applyIPWho(
	result *NetworkResult,
	data ipWhoResult,
) {
	if result.ISP == "" {
		result.ISP = data.ISP
	}

	if result.Organization == "" {
		result.Organization =
			data.Organization
	}

	if result.Domain == "" {
		result.Domain = data.Domain
	}

	if result.ConnectionType == "" {
		result.ConnectionType =
			data.ConnectionType
	}

	result.VPN = FlagResult{
		Available: true,
		Value:     data.VPN,
		Source:    "IPWho",
	}

	result.Tor = FlagResult{
		Available: true,
		Value:     data.Tor,
		Source:    "IPWho",
	}

	result.Threat = data.Threat

	result.Sources = appendUnique(
		result.Sources,
		"IPWho",
	)
}

func appendUnique(
	values []string,
	value string,
) []string {
	if value == "" {
		return values
	}

	for _, existing := range values {
		if existing == value {
			return values
		}
	}

	return append(
		values,
		value,
	)
}
