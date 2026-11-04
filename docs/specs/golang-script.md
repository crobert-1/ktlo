# GoLang Script

## Overview

The GoLang script is used as the entry point for automating the KTLO process.

## Input arguments

- Target repository: A full URL for a GitHub or GitLab repository, e.g. `https://github.com/signalfx/splunk-otel-collector`
- The initial version is read-only and does not accept or require credentials. The target repository and requested
  information must be accessible anonymously; private repositories are out of scope for this version.
- Requested checks: A list of checks to perform, e.g. `pr-status`, `open-vuln`, `flaky-tests`.
  The script will perform the requested checks on the target repository and return the results. Default is to perform all checks if no specific checks are requested.

## Functionalities

The script itself will be as thin of a wrapper as possible, delegating the actual work to subcommands for each
functionality. The subcommands will be implemented in separate packages or modules to keep the code organized
and maintainable. Initial subcommands must not create, update, merge, or comment on repository resources.
