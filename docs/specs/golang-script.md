# GoLang Script

## Overview

The GoLang script is the entry point for running enabled KTLO functionality against configured repositories.

## Input

The script accepts a required `--config` argument pointing to a YAML file. The YAML has top-level camelCase keys for
each functionality and a `repositories` list shared by enabled functionality. The initial configuration is:

```yaml
prStatus:
  enabled: true
  updateFlakyAnalysisReport: true

failingTestCheck:
  enabled: false
  lookback: 7d

repositories:
  - https://github.com/owner/repository
```

`prStatus.enabled` defaults to `true` if omitted. Setting it to `false` skips PR status. At least one repository is
required when any functionality is enabled. `failingTestCheck.enabled` defaults to `false`; `lookback` defaults to
`7d` and accepts positive Go duration values such as `30m` and `12h`, plus `d` for days. Unknown configuration keys
are errors so misspelled settings are not silently ignored. Future functionalities will add their own top-level
camelCase blocks.
`prStatus.updateFlakyAnalysisReport` defaults to `true`. When false, PR Status reads the existing local flaky-analysis
report without updating it; this setting does not disable writing the PR Status `RESULTS.md` output.
The top-level `run` delegates command-line parsing to `parseOptions` and YAML loading to a dedicated `loadConfig` helper.

The CLI parsing API may return an optional flaky-analysis selector:

```go
type cliOptions struct {
    ConfigPath string
    Selector   *testresults.FailureSelector
}

func parseOptions(args []string, stderr io.Writer) (cliOptions, int)
func runLookup(config config.Config, selector testresults.FailureSelector, stdout io.Writer) error
```

When a selector is supplied, lookup reads local reports and does not create a GitHub client or make network requests.
The configured repository URL is parsed without credentials or network access to locate each report. Without a
selector, the executable follows the normal enabled-check flow and requires the GitHub token.

The top-level script also provides a read-only flaky-analysis lookup mode. It accepts one test name, workflow name, or
job name as a selector and returns the matching flaky-analysis information from the local report. The selector is
looked up for the repositories listed in the config. This mode reads the existing report and does not fetch logs or
start a new analysis.

The current implementation supports GitHub.com repositories. GitLab support is deferred. The script requires a token
from `GH_TOKEN`, falling back to `GITHUB_TOKEN` when the first value is unset or blank. Credentials are not accepted in
the YAML file or command-line arguments.

## CLI contract

The executable is named `ktlo`:

```text
ktlo --config config.yaml
ktlo --config config.yaml --test package.TestName/subtest
ktlo --config config.yaml --workflow CI
ktlo --config config.yaml --job unit-tests
```

Without a lookup selector, the script runs enabled functionality for every configured repository. With a lookup
selector, it reads the local report for each configured repository and returns matching information. Select at most
one of `--test`, `--workflow`, or `--job` in a single invocation. Invalid configuration or repository URLs
are usage errors; connection failures, timeouts, unsuccessful HTTP responses, and GraphQL errors are operational
errors. Tokens must never appear in standard output, standard error, or other log messages, including when an API
response echoes one. The script attempts each configured repository even if a previous repository fails, then exits
non-zero if any configured check failed.

Lookup output includes the matching test or workflow/job classification, pass/fail counts and percentages, `main`
history, PR branch and commit history, partial-data state, any `Unknown` reason, and representative evidence links.
The output preserves the distinction between `main` and each PR commit.

For each repository, `runChecks` creates one GitHub client and reuses it across checks. Run the failing-test analyzer
once before PR Status when either `failingTestCheck.enabled` is true or enabled PR Status has
`updateFlakyAnalysisReport` set to true. If both request an update, do not run the analyzer twice. Continue to PR Status
after a partial-analysis error if the report was written, and return a non-zero exit code for collected errors.

PR status results are written to `results/github/<owner>/<repository>/RESULTS.md` in the script's working directory,
grouped into `Passed`, `Failed`, and `Others` Markdown tables with a separate author column, a clickable `title
(#number)` pull request link, and the CI/CD rollup state. After each results file is completely written, the script
reports its path and successful completion on standard output. Review and merge status classification is future work.
Usage and operational errors are written to standard error. Invalid arguments, unsupported configuration keys, or a
check that cannot complete result in a non-zero exit code.

## Build

Run `make ktlo` from the repository root to build the Go executable and its dependencies. The output is
`./bin/ktlo`.

## Functionalities

The script delegates each functionality to its own package or module. Initial functionality must not create, update,
merge, or comment on repository resources.

## Future requirements

- Add progress log messages at appropriate log levels so users can follow repository processing and check progress.
- Make the PR status dependency-bump author list configurable in `config.yaml` instead of keeping it in Go constants.
