# KTLO

## Goal

The goal is to use golang scripts to automate the KTLO process for the OpenTelemetry Collector and its components,
to reduce the manual workload currently required to fulfill KTLO requirements.

## Implementation

### Technologies

- golang - Used to implement scripts and tools for the KTLO process.
- GraphQL - Used to query the GitHub.com GraphQL API for metadata information. GitLab support is deferred.

### Functionalities

#### Central GoLang script

- Refer to [golang-script.md](./golang-script.md) for the script's spec.
- The executable uses explicit subcommands: `vulnerability-check` (`v`), `flaky-test-analysis` (`f`),
  `generate-report` (`g`), and `dependency-pr-status update` (`d u`). Each command has its own help output.
- `config.yaml` stores shared repository definitions and reusable settings. The command selects the operation; check
  enablement is no longer controlled by config booleans.
- Per-run flags can override the flaky-test lookback, rebuild stored analysis, parse linked workflow failures during PR
  status, and override the stale PR threshold. See [golang-script.md](./golang-script.md) for the flag contract.

#### PR status checks

Run with `ktlo dependency-pr-status update`. Refer to [pr-status.md](./pr-status.md) for its behavior and output.

#### Vulnerability checks

Run with `ktlo vulnerability-check <ID>`. Refer to [vulnerability-check.md](./vulnerability-check.md).

#### Flaky test checks

Run analysis with `ktlo flaky-test-analysis` and render Markdown with `ktlo generate-report`. Refer to
[failing-test-check.md](./failing-test-check.md) and [flaky-test-analysis.md](./flaky-test-analysis.md).


### Repository structure

- One top-level Go executable named `ktlo` that can handle all the KTLO functionalities, with subcommands selecting
  which operation to run and YAML configuration providing shared settings and repository definitions.
- Separate packages or modules for each functionality to keep the code organized and maintainable.
- Use a shared internal package for shared functionality
