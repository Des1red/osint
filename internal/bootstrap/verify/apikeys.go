package verify

import (
	"fmt"
	"os"
	"strings"

	"osint/internal/models"
)

func verifyIPKeys() {
	if strings.TrimSpace(
		models.ScopeInput.IpAddress,
	) == "" {
		return
	}

	verifyIPWho()
	verifyAbuseIPDB()
	verifyVirusTotal()
	verifyIPQS()
}

func verifyFullNameKeys() {
	if strings.TrimSpace(
		models.ScopeInput.FullName,
	) == "" {
		return
	}
}

func verifyIPWho() {
	if os.Getenv("IPWHO_API_KEY") != "" {
		return
	}

	fmt.Println(
		"optional: IPWHO_API_KEY is not configured",
	)

	fmt.Println(
		"IPWho enrichment will be skipped",
	)

	fmt.Println(
		"get a free key at: https://www.ipwho.org/free-plan",
	)
}

func verifyAbuseIPDB() {
	if os.Getenv("ABUSEIPDB_API_KEY") != "" {
		return
	}

	fmt.Println(
		"optional: ABUSEIPDB_API_KEY is not configured",
	)

	fmt.Println(
		"AbuseIPDB reputation checks will be skipped",
	)

	fmt.Println(
		"get a free key at: https://www.abuseipdb.com/register",
	)
}

func verifyVirusTotal() {
	if os.Getenv("VIRUSTOTAL_API_KEY") != "" {
		return
	}

	fmt.Println(
		"optional: VIRUSTOTAL_API_KEY is not configured",
	)

	fmt.Println(
		"VirusTotal reputation checks will be skipped",
	)

	fmt.Println(
		"get a free key at: https://www.virustotal.com/gui/join-us",
	)
}

func verifyIPQS() {
	if os.Getenv("IPQS_API_KEY") != "" {
		return
	}

	fmt.Println(
		"optional: IPQS_API_KEY is not configured",
	)

	fmt.Println(
		"IPQualityScore reputation checks will be skipped",
	)

	fmt.Println(
		"get a free key at: https://www.ipqualityscore.com/create-account",
	)
}

func VerifyKeys() {
	verifyIPKeys()

	verifyFullNameKeys()
}
