# OSINT-Master Dependencies

OSINT-Master is a Go-based passive OSINT tool.

Most application dependencies are managed through Go modules, while a small number of system and runtime dependencies must be available on the host machine.

---

## 1. Go

OSINT-Master requires a working Go installation to build from source.

Verify Go is installed:

```bash
go version
```

Project Go dependencies are defined in:

```text
go.mod
go.sum
```

Install or update them with:

```bash
go mod download
```

or:

```bash
go mod tidy
```

---

## 2. Go Module Dependencies

The project uses external Go packages including:

### pflag

Package:

```text
github.com/spf13/pflag
```

Used for command-line flag parsing.

Examples:

```text
-n, --name
-i, --ip
-u, --username
-d, --domain
-o, --output
--full
--debug
--install
--uninstall
```

---

### chromedp

Package:

```text
github.com/chromedp/chromedp
```

Used for Chromium-based browser automation during WebSearch discovery.

`chromedp` is a compile-time Go dependency.

It does **not** provide a Chromium browser itself.

A compatible Chrome or Chromium executable must also exist on the host system.

---

### clihelp

Package:

```text
github.com/Des1red/clihelp
```

Used to format command-line help output.

---

## 3. Chrome / Chromium

WebSearch requires a locally installed Chrome or Chromium browser.

OSINT-Master searches for the following executables in order:

```text
google-chrome
google-chrome-stable
chromium
chromium-browser
```

At least one must be available through the system `PATH`.

Verify manually with:

```bash
which google-chrome
```

or:

```bash
which chromium
```

During startup, bootstrap verifies that a supported browser executable exists.

If none can be found, OSINT-Master cannot start browser-backed WebSearch operations.

---

## 4. Browser Profile

OSINT-Master uses a persistent Chromium profile.

On Linux, the profile is stored under:

```text
~/.osint-master/chrome-profile
```

The directory is created automatically by bootstrap with permissions:

```text
0700
```

The persistent profile allows browser state to survive between OSINT-Master runs, including data such as:

```text
cookies
local storage
browser preferences
session state
```

The browser package itself does not create or discover this directory.

Bootstrap prepares the browser environment before the Chrome package is used.

---

## 5. Cache

OSINT-Master uses the operating system's user cache directory.

On a typical Linux system this resolves to:

```text
~/.cache/osint
```

The exact location is determined through Go's:

```go
os.UserCacheDir()
```

Geolocation cache data is stored at:

```text
<cache directory>/osint/geo.json
```

Bootstrap verifies the cache during application startup.

---

## 6. Network Access

OSINT-Master requires outbound network access for OSINT collection.

Depending on the enabled engine, this may include access to:

```text
search engines
RDAP servers
public websites
social platforms
geolocation providers
network intelligence providers
IP reputation providers
```

Blocking outbound HTTPS or DNS traffic may prevent individual engines from working.

---

## 7. API Keys

Some OSINT providers require API credentials.

Bootstrap verifies configured API keys before normal execution.

The exact keys required depend on the engines enabled in the project.

Current provider-backed features include areas such as:

```text
IP reputation
geolocation
platform lookups
network intelligence
```

Examples of provider integrations currently used by the IP intelligence system include:

```text
AbuseIPDB
VirusTotal
IPQualityScore
```

Provider credentials should not be committed directly into the repository.

Use the configuration or environment mechanism defined by the project.

---

## 8. Filesystem Permissions

Normal OSINT searches do not require root privileges.

However, installation to:

```text
/usr/local/bin/osint
```

requires permission to write to:

```text
/usr/local/bin
```

The install operation builds the project and copies the resulting binary to:

```text
/usr/local/bin/osint
```

Depending on the system, this may require:

```bash
sudo
```

The same applies when removing an installed binary that is owned by root.

---

## 9. Installation Dependencies

To build and install OSINT-Master from source, the machine needs:

```text
Go
Git
Chrome or Chromium
Internet access
```

A typical development setup can be checked with:

```bash
go version
git --version
which google-chrome
which chromium
```

Only one supported Chrome/Chromium executable is required.

---

## 10. Runtime Dependency Flow

The runtime dependency flow is:

```text
OSINT-Master
│
├── Bootstrap
│   ├── verifies API/configuration requirements
│   ├── verifies cache
│   ├── verifies Chrome/Chromium
│   └── prepares the persistent browser profile
│
├── Native HTTP engines
│   └── public APIs and web endpoints
│
└── Chromium WebSearch
    ├── chromedp
    ├── installed Chrome/Chromium executable
    └── persistent browser profile
```

The Chrome package is responsible for browser lifecycle and browser configuration.

System dependency discovery and filesystem preparation belong to bootstrap.

---

## 11. Development Setup

Clone the repository and enter the project:

```bash
git clone <repository>
cd osint-master
```

Download Go dependencies:

```bash
go mod download
```

Verify the project builds:

```bash
go build ./...
```

Run directly:

```bash
go run main.go --help
```

Example:

```bash
go run main.go -n "Example Name"
```

Full output:

```bash
go run main.go -n "Example Name" --full
```

Debug output:

```bash
go run main.go -n "Example Name" --debug
```

---

## 12. Dependency Responsibilities

Dependencies are intentionally separated by responsibility.

```text
Go modules
    Compile-time application dependencies

Bootstrap
    System/runtime dependency verification

models
    Shared paths and resolved runtime configuration

chrome
    Chromium lifecycle and browser behavior

engines
    External OSINT providers and collection

information
    Result presentation only
```

This separation prevents collection or presentation packages from taking responsibility for system setup.