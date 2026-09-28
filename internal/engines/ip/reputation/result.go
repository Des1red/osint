package reputation

type ProviderError struct {
	Provider string
	Error    string
}

type AbuseIPDBResult struct {
	Available bool

	AbuseConfidenceScore int

	TotalReports  int
	DistinctUsers int

	LastReportedAt string

	UsageType string
	Domain    string

	IsPublic      bool
	IsWhitelisted *bool

	Source string
}

type VirusTotalStats struct {
	Harmless   int
	Malicious  int
	Suspicious int
	Undetected int
	Timeout    int
}

type VirusTotalVotes struct {
	Harmless  int
	Malicious int
}

type VirusTotalResult struct {
	Available bool

	Reputation int

	Stats VirusTotalStats
	Votes VirusTotalVotes

	Tags []string

	LastAnalysisDate int64

	Source string
}

type IPQSResult struct {
	Available bool

	FraudScore int

	Proxy bool
	VPN   bool
	Tor   bool

	BotStatus      bool
	RecentAbuse    bool
	FrequentAbuser bool

	AbuseVelocity string

	Source string
}

type ReputationResult struct {
	IP string

	AbuseIPDB  AbuseIPDBResult
	VirusTotal VirusTotalResult
	IPQS       IPQSResult

	Errors []ProviderError
}
