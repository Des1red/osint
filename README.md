# OSINT-Master

OSINT-Master is a Go-based passive OSINT toolkit for collecting, normalizing, and presenting publicly available intelligence about:

- full names
- IP addresses
- usernames
- domains

The project separates collection, enrichment, attribution, knowledge storage, and presentation so that discovered information is not automatically treated as verified identity data.

> **Scope:** OSINT-Master is designed for passive information gathering and audit-oriented research. It does not perform service exploitation, brute forcing, vulnerability exploitation, or active port scanning.

---

## Features

### Full Name Intelligence

The FullName engine performs bounded identity discovery from a supplied first-and-last name.

Current behavior includes:

- deterministic full-name candidate generation
- provider-specific lookup strategies
- GitHub discovery
- GitLab discovery
- Facebook discovery and canonical-profile handling
- display-name matching after discovery
- canonical deduplication
- shared evidence-aware enrichment
- public WebSearch enrichment
- related-person handling where supported by evidence

The FullName engine keeps **discovery** and **identity matching** separate. A platform result is not treated as the requested person until the engine's matching stage accepts it.

---

### Username Intelligence

The Username engine performs exact username lookup, enrichment, and bounded nearby-variant discovery.

Supported platform collectors:

- GitHub
- GitLab
- Reddit
- Facebook
- Instagram
- Twitter / X
- TikTok

Current behavior includes:

- exact username lookup
- provider-specific account validation
- canonical-profile deduplication
- shared enrichment from exact/root profiles
- bounded username variant generation
- concurrent variant lookup
- Exact / Close / Broad / None username-string matching
- related-account discovery
- ambiguity-preserving attribution

Generated username variants are discovery candidates only. Similarity does not imply real-world identity.

---

### IP Intelligence

The IP engine performs passive IPv4 and IPv6 intelligence collection.

Current components include:

- IP validation
- RDAP registration data
- geolocation from multiple providers
- geolocation consensus and confidence handling
- coordinate disagreement analysis
- cached geolocation results
- routing and ASN intelligence
- announced-prefix information
- ASN neighbours
- reverse DNS
- provider / ISP information
- network classification
- anycast intelligence
- routing/history components
- reputation intelligence

Current reputation integrations include:

- AbuseIPDB
- VirusTotal
- IPQualityScore

The engine preserves provider results rather than collapsing unrelated reputation concepts into one custom risk score.

---

### Domain / DNS Intelligence

The Domain engine performs passive DNS intelligence and bounded subdomain analysis.

Root DNS collection includes:

- A
- AAAA
- CNAME
- MX
- NS
- TXT

Subdomain functionality includes:

- passive Certificate Transparency discovery through `crt.sh`
- normalization and deduplication
- bounded candidate resolution
- concurrent A / AAAA / CNAME resolution
- filtering of stale certificate names without current DNS data
- potential subdomain takeover-risk analysis

Takeover analysis is deliberately conservative. A third-party CNAME is not automatically considered vulnerable. The current implementation requires a recognized third-party target and an upstream DNS `NXDOMAIN` result before reporting a **potential** risk, and manual validation is still required.

---

## Shared Enrichment

`internal/enrich` is the shared evidence-processing layer used by identity-oriented engines.

Its central rule is:

```text
discovered fact != target-owned fact
```

The enrichment pipeline separates:

```text
collection
    ↓
evidence / provenance
    ↓
ownership organization
    ↓
deduplication
    ↓
final structured result
```

Current enrichment stages include:

1. local extraction
2. WebSearch discovery
3. first-hop page inspection
4. already-known relative/person enrichment
5. pre-organization deduplication
6. evidence-aware ownership organization
7. post-organization deduplication

Facts can be assigned to:

- the root target
- a known related person
- an unattributed bucket when ownership is ambiguous

The search path itself is never treated as proof of ownership.

### Browser-backed WebSearch

Public WebSearch uses Chromium through `chromedp`.

Current search-engine support includes:

- DuckDuckGo
- Brave

A shared persistent browser process is used for the enrichment operation, with temporary tabs created for individual searches.

---

## Installation

### Requirements

OSINT-Master requires:

- Go
- Git
- Internet access
- Chrome or Chromium for browser-backed WebSearch

Supported browser executables are searched in the following order:

```text
google-chrome
google-chrome-stable
chromium
chromium-browser
```

Verify the basic environment:

```bash
go version
git --version
which google-chrome
which chromium
```

Only one supported Chrome/Chromium executable is required.

### Clone and build

```bash
git clone <repository-url>
cd osint-master

go mod download
go build ./...
```

Run directly from source:

```bash
go run main.go --help
```

### Install binary

OSINT-Master can install its binary to:

```text
/usr/local/bin/osint
```

Run:

```bash
go run main.go --install
```

Writing to `/usr/local/bin` may require elevated permissions depending on the system.

To remove the installed binary:

```bash
osint --uninstall
```

or run the uninstall path from source if the binary is not available.

For the complete runtime and dependency notes, see [`dependencies.md`](dependencies.md).

---

## Usage

```text
osint [flags]
```

Main flags:

| Flag | Purpose |
|---|---|
| `-n, --name` | Search by full name |
| `-i, --ip` | Search by IPv4 or IPv6 address |
| `-u, --username` | Search by username |
| `-d, --domain` | Inspect DNS, enumerate passive subdomains, and identify potential takeover risks |
| `-o, --output` | Save output using a custom base filename |
| `--full` | Show full / uncapped output where supported |
| `--debug` | Enable debug output |
| `--install` | Build and install OSINT-Master |
| `--uninstall` | Remove the installed binary |
| `-h, --help` | Show help |

### Full name

```bash
go run main.go -n "Jane Doe"
```

Installed binary:

```bash
osint -n "Jane Doe"
```

### Username

```bash
osint -u exampleuser
```

### IP address

```bash
osint -i 8.8.8.8
```

IPv6 is also supported:

```bash
osint -i 2001:4860:4860::8888
```

### Domain

```bash
osint -d example.com
```

### Multiple engines

Multiple inputs can be supplied in one execution:

```bash
osint \
  -u exampleuser \
  -i 8.8.8.8 \
  -d example.com
```

Engines are isolated from one another. A fatal validation/setup failure in one requested engine does not prevent the remaining requested engines from running.

---

## Output

OSINT-Master separates collection from presentation.

Completed engine results are stored in the central `knowledge` package and are then rendered through the `information` layer.

Tree formatting is used for nested information:

```text
DNS Lookup
----------
└─ example.com
   ├─ A
   │  └─ 192.0.2.10
   ├─ MX
   │  └─ 10 mail.example.com
   └─ NS
      ├─ ns1.example.com
      └─ ns2.example.com
```

Large collections are bounded in normal terminal output where appropriate. Full output can expose the complete stored collection.

### Saving results

Use `-o` to provide a custom output base name:

```bash
osint -i 8.8.8.8 -o google
```

Engine-specific naming prevents multiple active engines from overwriting each other's files.

Example:

```text
google_ip.txt
google_domain.txt
google_username.txt
```

Without a custom name, generated output follows the project's engine/target naming convention.

---

## API Keys

Some provider-backed features can use optional API credentials.

Current environment variables include:

```text
IPWHO_API_KEY
ABUSEIPDB_API_KEY
VIRUSTOTAL_API_KEY
IPQS_API_KEY
```

Missing optional credentials should cause the affected provider to be skipped rather than stopping unrelated collection stages.

Do not commit API keys to the repository.

---

## Browser Profile and Cache

### Chromium profile

OSINT-Master uses a persistent Chromium profile.

Typical Linux path:

```text
~/.osint-master/chrome-profile
```

The profile can preserve browser state such as:

- cookies
- local storage
- browser preferences
- session state

Bootstrap prepares the profile before browser-backed discovery runs.

### Cache

OSINT-Master uses the operating system's user cache directory.

Typical Linux geolocation cache:

```text
~/.cache/osint/geo.json
```

The geolocation cache persists between runs and is managed by the IP geolocation package.

---

## Architecture

The project follows a clear responsibility boundary:

```text
CLI
 │
 ▼
cmd
 │
 ▼
bootstrap
 │
 ▼
engines.Manager
 │
 ├── FullName
 ├── IP
 ├── Username
 └── Domain
 │
 ▼
knowledge
 │
 ▼
information
 │
 ▼
terminal / output file
```

Identity engines can also use the shared enrichment system:

```text
engine
  │
  ▼
internal/enrich
  │
  ├── extract
  ├── provenance
  ├── WebSearch
  │    ├── discovery
  │    └── first-hop inspection
  ├── relatives
  ├── organize
  └── dedupe
  │
  ▼
structured enrichment result
```

### Engine lifecycle

The engine manager distinguishes between:

```text
ActiveEngines
    engines requested by the user

CompletedEngines
    engines that completed successfully and are safe to print/store
```

This prevents a failed engine from being presented as a successful empty lookup.

---

## Project Structure

High-level layout:

```text
.
├── cmd/
│   ├── flags.go
│   ├── help.go
│   ├── lifecycle.go
│   └── root.go
│
├── docs/
│   ├── domain_dns_engine_documentation.docx
│   ├── enrich_package_documentation.docx
│   ├── fullname_engine_documentation.docx
│   ├── ip_engine_documentation.docx
│   └── username_engine_documentation.docx
│
├── internal/
│   ├── bootstrap/
│   │   └── verify/
│   │
│   ├── engines/
│   │   ├── domain/
│   │   │   ├── subdomain/
│   │   │   └── takeover/
│   │   ├── fullname/
│   │   ├── ip/
│   │   │   ├── geo/
│   │   │   ├── history/
│   │   │   ├── network/
│   │   │   ├── rdap/
│   │   │   └── reputation/
│   │   ├── username/
│   │   └── manager.go
│   │
│   ├── enrich/
│   │   ├── dedupe/
│   │   ├── extract/
│   │   ├── model/
│   │   ├── opencorp/
│   │   ├── organize/
│   │   ├── provenance/
│   │   ├── relatives/
│   │   └── websearch/
│   │
│   ├── htmlx/
│   ├── httpx/
│   │   └── chrome/
│   ├── information/
│   ├── knowledge/
│   ├── logger/
│   ├── matcher/
│   ├── models/
│   └── platforms/
│
├── dependencies.md
├── go.mod
├── go.sum
└── main.go
```

### Main package responsibilities

| Package | Responsibility |
|---|---|
| `cmd` | CLI flags, help, root command, and lifecycle |
| `bootstrap` | Runtime preparation and dependency verification |
| `engines` | Top-level target-specific intelligence collection |
| `platforms` | Provider-specific account discovery |
| `enrich` | Shared extraction, provenance, WebSearch, ownership, relatives, and dedupe |
| `httpx` | Shared HTTP transport and Chromium browser infrastructure |
| `htmlx` | HTML metadata extraction helpers |
| `knowledge` | Central storage boundary for completed engine results |
| `information` | Human-readable terminal/file presentation |
| `logger` | Error and debug logging |
| `matcher` | Shared matching logic |
| `models` | Shared CLI/runtime/scope models |

---

## Engine Design Principles

OSINT-Master is built around several rules:

- **Passive collection first.**
- **Bound discovery.** Candidate generation and web expansion are intentionally limited.
- **Provider-specific interpretation.** A generic HTTP `200` is not automatically treated as a valid account.
- **Discovery is not identity.**
- **Evidence is not ownership.**
- **Ambiguous information stays ambiguous.**
- **Canonicalize before deduplicating.**
- **Keep provider failures isolated.**
- **Do not invent confidence where providers disagree.**
- **Keep collection, knowledge storage, and presentation separate.**
- **Do not treat a potential subdomain takeover indicator as confirmed exploitability.**

---

## Error Handling

OSINT-Master uses graceful degradation.

Examples:

- malformed engine input can fail that engine
- one failed engine does not terminate other requested engines
- one failed platform collector does not erase other platform results
- one failed IP provider does not erase other IP intelligence
- a Certificate Transparency failure does not erase root DNS information
- temporary DNS failures are not treated as takeover evidence

Errors are logged separately from positive or negative identity decisions.

---

## Documentation

Detailed internal documentation is kept in `docs/`.

- `docs/fullname_engine_documentation.docx`  
  FullName candidate generation, provider lookup, matching, canonicalization, and enrichment.

- `docs/username_engine_documentation.docx`  
  Exact username lookup, supported collectors, enrichment, variants, matching semantics, and concurrency.

- `docs/ip_engine_documentation.docx`  
  RDAP, geolocation, network intelligence, routing, reverse DNS, anycast, reputation, and output behavior.

- `docs/domain_dns_engine_documentation.docx`  
  DNS records, Certificate Transparency subdomain discovery, bounded resolution, and potential takeover-risk analysis.

- `docs/enrich_package_documentation.docx`  
  Evidence, provenance, extraction, WebSearch, first-hop inspection, ownership organization, relatives, and deduplication.

- [`dependencies.md`](dependencies.md)  
  Go modules, Chrome/Chromium, browser profile, cache, API keys, filesystem permissions, and development setup.

---

## Development

Download dependencies:

```bash
go mod download
```

Normalize the module files:

```bash
go mod tidy
```

Build all packages:

```bash
go build ./...
```

Run help:

```bash
go run main.go --help
```

Example development run:

```bash
go run main.go -d example.com --debug
```

Before committing changes:

```bash
gofmt -w .
go build ./...
```

---

## Security and Ethical Use

OSINT-Master is intended for legitimate passive research, defensive assessment, auditing, and authorized security work.

Users are responsible for:

- complying with applicable laws and regulations
- respecting privacy and data-protection requirements
- respecting third-party service terms and rate limits
- using collected information responsibly
- obtaining authorization where an assessment requires it

Subdomain takeover functionality performs risk identification only. It does not claim a vulnerable third-party resource or attempt exploitation.

---

## Current Scope

OSINT-Master currently focuses on four primary target types:

```text
Full Name
IP Address
Username
Domain
```

The project intentionally favors bounded, explainable collection and evidence-preserving output over unrestricted crawling or aggressive scanning.