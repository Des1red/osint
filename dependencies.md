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

## 4. Xvfb

OSINT-Master uses Xvfb to provide Chromium with a virtual X11 display.

Chromium continues to run in normal headed mode:

```text
headless=false
```

but its window is rendered inside the virtual display rather than appearing on the user's desktop.

This allows OSINT-Master to preserve normal browser behavior without displaying Chromium windows during execution.

Bootstrap requires the following executable to exist:

```text
Xvfb
```

It must be available through the system `PATH`.

Verify manually with:

```bash
which Xvfb
```

### Fedora

Install Xvfb with:

```bash
sudo dnf install xorg-x11-server-Xvfb
```

### Debian / Ubuntu

Install Xvfb with:

```bash
sudo apt install xvfb
```

OSINT-Master is currently intended for Linux systems.

During startup, bootstrap verifies that Xvfb is available before preparing the virtual browser environment.

If Xvfb is not installed, startup fails before browser-backed engines are started.

---

## 5. Virtual Display Environment

OSINT-Master starts and manages Xvfb directly.

Users do not need to launch the tool through:

```bash
xvfb-run
```

Instead, the runtime creates its own Xvfb process.

Xvfb is started using:

```text
-displayfd
```

which allows the X server itself to select an available display number.

The resulting display is assigned to the process environment through:

```text
DISPLAY=:<number>
```

Chromium inherits this environment and renders inside the virtual X11 display.

The configured virtual screen is:

```text
1920x1080x24
```

TCP access to the X server is disabled.

The runtime flow is approximately:

```text
OSINT-Master
    │
    ├── start Xvfb
    │
    ├── obtain available display number
    │
    ├── set DISPLAY
    │
    └── start Chromium
            │
            └── render into virtual X11 display
```

No Chromium window is displayed on the user's normal desktop.

---

## 6. Xvfb Runtime State

OSINT-Master stores information about its Xvfb process so that stale virtual-display processes can be safely recovered after an abnormal previous run.

The Xvfb runtime state is stored under:

```text
~/.osint-master/xvfb-state.json
```

The state records Linux process identity information including:

```text
PID
executable path
Linux process start time
display number
```

The executable path and process start time are used together with the PID to protect against PID reuse.

A stale PID alone is not considered sufficient evidence that a process belongs to OSINT-Master.

During bootstrap state preparation:

```text
xvfb-state.json exists?
    │
    ├── no
    │   └── continue
    │
    └── yes
        │
        ├── verify PID
        ├── verify executable
        └── verify Linux process start time
                │
                ├── exact match
                │   └── terminate stale OSINT Xvfb process
                │
                └── mismatch
                    └── do not terminate process
```

The stale state file is removed after recovery.

OSINT-Master does not blindly remove global X11 lock files or terminate unrelated Xvfb processes.

---

## 7. Xvfb Process Cleanup

Xvfb is protected through multiple cleanup mechanisms.

### Normal shutdown

When OSINT-Master finishes normally:

```text
virtualenv.Close()
```

performs the following cleanup:

```text
terminate Xvfb
restore the previous DISPLAY value
remove the Xvfb runtime state file
```

### Parent process termination

Because OSINT-Master currently targets Linux, the Xvfb child process is configured with a Linux parent-death signal.

If the OSINT-Master process dies unexpectedly, the operating system sends:

```text
SIGKILL
```

to its Xvfb child.

This reduces the chance of an orphan Xvfb process remaining after abnormal termination.

### Previous-run recovery

If a previous run somehow leaves a valid OSINT-owned Xvfb process behind, bootstrap verifies its persisted Linux process identity before terminating it.

Together these mechanisms provide:

```text
normal cleanup
+
parent-death cleanup
+
next-run stale-process recovery
```

---

## 8. Browser Profile

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

Before each run, stale Chromium profile lock files created by previous interrupted sessions are cleared from the OSINT-Master-owned Chromium profile.

---

## 9. Cache

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

## 10. Network Access

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

## 11. API Keys

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
IPWho
AbuseIPDB
VirusTotal
IPQualityScore
```

Provider credentials should not be committed directly into the repository.

Use environment variables or the configuration mechanism defined by the project.

---

## 12. Filesystem Permissions

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

Runtime state stored under:

```text
~/.osint-master
```

belongs to the current user and does not require root access.

---

## 13. Installation Dependencies

To build and run OSINT-Master from source, the machine needs:

```text
Go
Git
Chrome or Chromium
Xvfb
Internet access
```

A typical development setup can be checked with:

```bash
go version
git --version
which google-chrome
which chromium
which Xvfb
```

Only one supported Chrome/Chromium executable is required.

Xvfb is required for the virtual browser environment.

---

## 14. Runtime Dependency Flow

The runtime dependency flow is:

```text
OSINT-Master
│
├── Bootstrap
│   ├── verifies API/configuration requirements
│   ├── verifies cache
│   ├── verifies Chrome/Chromium
│   ├── verifies Xvfb
│   ├── prepares Chromium runtime state
│   └── prepares virtual-environment state
│
├── Virtual Environment
│   ├── starts Xvfb
│   ├── selects an available X11 display
│   ├── sets DISPLAY
│   ├── persists Linux process identity
│   └── manages Xvfb cleanup
│
├── Native HTTP engines
│   └── public APIs and web endpoints
│
└── Chromium WebSearch
    ├── chromedp
    ├── installed Chrome/Chromium executable
    ├── persistent browser profile
    └── Xvfb virtual display
```

The Chrome package is responsible for Chromium lifecycle and browser behavior.

The virtual environment package is responsible for the Xvfb lifecycle and virtual display.

System dependency discovery and runtime-state preparation belong to bootstrap.

---

## 15. Development Setup

Clone the repository and enter the project:

```bash
git clone <repository>
cd osint-master
```

Install the required Linux runtime dependencies.

For Fedora:

```bash
sudo dnf install xorg-x11-server-Xvfb
```

Install Chrome or Chromium if one is not already available.

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

No separate `xvfb-run` wrapper is required.

OSINT-Master starts and manages its own virtual display.

---

## 16. Dependency Responsibilities

Dependencies are intentionally separated by responsibility.

```text
Go modules
    Compile-time application dependencies

Bootstrap
    System/runtime dependency verification
    Runtime-state preparation

models
    Shared paths and resolved runtime configuration

virtualenv
    Virtual-environment lifecycle
    Xvfb startup and cleanup

xvfb
    Linux Xvfb process management
    DISPLAY configuration
    Persistent Linux process identity

chrome
    Chromium lifecycle and browser behavior

engines
    External OSINT providers and collection

information
    Result presentation only
```

This separation prevents collection or presentation packages from taking responsibility for system setup and keeps browser behavior separate from virtual-display management.