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
- The executable takes `--config <yaml-file>`. The configuration has a `repositories` list and one camelCase block per
  functionality, such as `prStatus: { enabled: true }`.
- `prStatus.enabled` defaults to true and may be disabled explicitly. Enabled functionality runs for each configured
  repository.
- The top-level executable also provides a read-only lookup mode accepting one test name, workflow name, or job name.
  It reads the local flaky-analysis report for configured repositories without fetching logs or starting analysis.

#### PR status checks

Implemented first. Refer to [pr-status.md](./pr-status.md) for the initial check's behavior and output.


#### Open vulnerability checks

Future functionality; not part of the initial configuration contract.


#### Flaky test checks

The `failingTestCheck` configuration block is available as an opt-in scaffold. Refer to
[failing-test-check.md](./failing-test-check.md); analysis and reporting behavior are future work.


### Repository structure

- One top-level Go executable named `ktlo` that can handle all the KTLO functionalities, with YAML configuration
  selecting which functionalities to run.
- Separate packages or modules for each functionality to keep the code organized and maintainable.
- Use a shared internal package for shared functionality
