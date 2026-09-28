# API Configuration

OSINT-Master can use optional third-party API services to enrich collected information.

API keys are **not required to run OSINT-Master**.

When an optional API key is unavailable, the associated provider or enrichment step is skipped while the rest of the engine continues normally.

---

## Supported API Keys

The following environment variables are currently supported:

| Environment Variable | Provider | Used By | Purpose |
|---|---|---|---|
| `IPWHO_API_KEY` | IPWho | IP Engine | Additional IP and network enrichment |
| `ABUSEIPDB_API_KEY` | AbuseIPDB | IP Engine | IP abuse and reputation information |
| `VIRUSTOTAL_API_KEY` | VirusTotal | IP Engine | Reputation and security intelligence |
| `IPQS_API_KEY` | IPQualityScore | IP Engine | IP reputation and risk information |

All currently supported API keys are optional.

---

# IPWho

Environment variable:

```bash
IPWHO_API_KEY
```

IPWho provides additional enrichment for IP investigations.

If the key is not configured, OSINT-Master prints:

```text
optional: IPWHO_API_KEY is not configured
IPWho enrichment will be skipped
get a free key at: https://www.ipwho.org/free-plan
```

Create an API key at:

```text
https://www.ipwho.org/free-plan
```

Example:

```bash
export IPWHO_API_KEY="your_api_key"
```

---

# AbuseIPDB

Environment variable:

```bash
ABUSEIPDB_API_KEY
```

AbuseIPDB is used by the IP reputation subsystem.

If the key is not configured, OSINT-Master prints:

```text
optional: ABUSEIPDB_API_KEY is not configured
AbuseIPDB reputation checks will be skipped
get a free key at: https://www.abuseipdb.com/register
```

Create an account and obtain an API key at:

```text
https://www.abuseipdb.com/register
```

Example:

```bash
export ABUSEIPDB_API_KEY="your_api_key"
```

---

# VirusTotal

Environment variable:

```bash
VIRUSTOTAL_API_KEY
```

VirusTotal is used as an additional IP reputation and security intelligence source.

If the key is not configured, OSINT-Master prints:

```text
optional: VIRUSTOTAL_API_KEY is not configured
VirusTotal reputation checks will be skipped
get a free key at: https://www.virustotal.com/gui/join-us
```

Create an account at:

```text
https://www.virustotal.com/gui/join-us
```

Example:

```bash
export VIRUSTOTAL_API_KEY="your_api_key"
```

---

# IPQualityScore

Environment variable:

```bash
IPQS_API_KEY
```

IPQualityScore is used by the IP reputation subsystem.

If the key is not configured, OSINT-Master prints:

```text
optional: IPQS_API_KEY is not configured
IPQualityScore reputation checks will be skipped
get a free key at: https://www.ipqualityscore.com/create-account
```

Create an account at:

```text
https://www.ipqualityscore.com/create-account
```

Example:

```bash
export IPQS_API_KEY="your_api_key"
```

---

# Temporary Configuration

Environment variables can be configured for the current terminal session using `export`.

Example:

```bash
export IPWHO_API_KEY="your_ipwho_key"
export ABUSEIPDB_API_KEY="your_abuseipdb_key"
export VIRUSTOTAL_API_KEY="your_virustotal_key"
export IPQS_API_KEY="your_ipqs_key"
```

The variables remain available until the terminal session is closed.

You can verify that a variable exists with:

```bash
printenv IPWHO_API_KEY
```

or:

```bash
echo "$IPWHO_API_KEY"
```

Avoid displaying API keys in shared terminals, screenshots, logs, or recordings.

---

# Persistent Configuration

On Linux, API keys can be added to the shell configuration file.

For Bash:

```bash
nano ~/.bashrc
```

Add:

```bash
export IPWHO_API_KEY="your_ipwho_key"
export ABUSEIPDB_API_KEY="your_abuseipdb_key"
export VIRUSTOTAL_API_KEY="your_virustotal_key"
export IPQS_API_KEY="your_ipqs_key"
```

Reload the configuration:

```bash
source ~/.bashrc
```

For Zsh, use:

```bash
nano ~/.zshrc
```

and then:

```bash
source ~/.zshrc
```

---

# Running With One API Key

It is not necessary to configure every provider.

For example:

```bash
export VIRUSTOTAL_API_KEY="your_api_key"

osint -i 8.8.8.8
```

VirusTotal enrichment will be available while providers without configured API keys will be skipped.

---

# Running Without API Keys

OSINT-Master can also run without any API credentials:

```bash
osint -i 8.8.8.8
```

The key verification layer detects which optional credentials are unavailable.

For an IP investigation, it currently checks:

```text
IPWHO_API_KEY
ABUSEIPDB_API_KEY
VIRUSTOTAL_API_KEY
IPQS_API_KEY
```

Missing credentials do not stop the IP engine.

Instead, only the corresponding optional provider is skipped.

Other available IP intelligence sources and collectors continue to run.

---

# Key Verification

API key verification is handled before provider-specific enrichment is performed.

For IP investigations, the current verification flow is:

```text
VerifyKeys
│
└─ verifyIPKeys
   │
   ├─ verifyIPWho
   ├─ verifyAbuseIPDB
   ├─ verifyVirusTotal
   └─ verifyIPQS
```

The checks are performed only when an IP address was supplied.

If no IP scope is active, IP API key warnings are not printed.

The project also contains a Full Name key-verification stage:

```text
VerifyKeys
│
└─ verifyFullNameKeys
```

No Full Name API credentials are currently required by that stage.

---

# Security

API keys are credentials and should be treated as secrets.

Do not:

```text
commit API keys to Git
hard-code API keys into Go source files
include API keys in README examples
store real API keys in documentation
publish keys in screenshots or terminal recordings
```

OSINT-Master reads credentials from environment variables so that secrets remain separate from the source code.

For example, avoid:

```go
const apiKey = "actual-secret-key"
```

Use environment configuration instead:

```go
key :=
	os.Getenv(
		"VIRUSTOTAL_API_KEY",
	)
```

---

# Git and Environment Files

If API keys are stored in a local environment file, ensure that file is excluded from Git.

For example:

```text
.env
.env.local
.env.*
```

A repository may provide a safe example such as:

```text
.env.example
```

containing only variable names:

```bash
IPWHO_API_KEY=
ABUSEIPDB_API_KEY=
VIRUSTOTAL_API_KEY=
IPQS_API_KEY=
```

Never place real credentials in the example file.

---

# API Limits

Third-party providers may impose:

- request limits;
- daily or monthly quotas;
- account-level restrictions;
- endpoint-specific limits;
- free-plan limitations.

OSINT-Master should treat failures from optional providers as partial enrichment failures rather than failures of the entire investigation whenever possible.

The exact quotas and restrictions are controlled by each provider and may change independently of OSINT-Master.

---

# Current Configuration Summary

A completely configured IP environment may contain:

```bash
export IPWHO_API_KEY="..."
export ABUSEIPDB_API_KEY="..."
export VIRUSTOTAL_API_KEY="..."
export IPQS_API_KEY="..."
```

A minimal installation may contain none of them.

Both configurations are valid.

API credentials expand the amount of enrichment available to OSINT-Master, but they are not required for the core application to operate.