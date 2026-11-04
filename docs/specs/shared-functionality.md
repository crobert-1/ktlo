# Shared functionality

The shared functionality is a set of common functions and utilities that can be used by the various KTLO scripts and tools.
The shared functionality is implemented in a separate package or module to keep the code organized and maintainable.

## Overview

### GraphQL client

Used to connect to provider GraphQL APIs. For the initial credential-free, read-only version, confirm that each
required query can access the requested data anonymously. Public repository visibility does not by itself establish
that every provider API exposes every required field without authentication. If a required query needs credentials,
record that as a design blocker rather than prompting for a token or silently returning incomplete results.

### GraphQL query builder

Used to build GraphQL queries. Should have a set of functions that can be used to build queries for the various
KTLO checks, and should have common queries in shared variables, with the ability to substitute variables into the
queries. Variables that may need substitution include repository name, branch name, and test name.

### GitHub/GitLab API client

Used to connect to GitHub/GitLab APIs for read-only operations in the initial version. Keep provider-specific
authentication and schema differences behind this boundary; do not expose write operations in the initial API.

### Flaky test checker

- Query repo's API to determine success/failure of recent test runs for a given test, and determine if the test is flaky based on the results.
- Check failure logs for known flaky test patterns, and determine if the test is flaky based on the logs.
- Store the results of the flaky test checks in a local cache to avoid re-checking the same tests multiple times.
- Have a set list of log messages that are known to be flaky, and check the failure logs for those messages.
