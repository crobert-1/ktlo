# <Feature or component name>

## Purpose

Describe the problem this feature solves and who uses it.

## Scope

### In scope

- <Observable behavior this specification covers>

### Out of scope

- <Related behavior this specification does not cover>

## Behavior

Describe the expected behavior in user terms. Include representative inputs, outputs, and edge cases.

## Command-line contract

If this feature has a command-line interface, specify:

- Command and subcommand names
- Flags, arguments, defaults, and validation
- Whether credentials are required, supported credential sources, and minimum permissions
- Data or operations unavailable without authentication
- Output format and exit codes

Otherwise, write `Not applicable`.

## Go package/API contract

List the package's responsibility and proposed exported API. Include function signatures, input/output types,
error behavior, and any provider-specific boundaries. Mark undecided items as open questions instead of guessing.

```go
// Proposed API signatures go here; these are design contracts, not implementation.
```

## Data and external services

Specify external APIs, required fields/scopes, pagination, rate limits, caching, and data retention where relevant.

## Failure handling

Describe expected handling for invalid input, authentication/authorization failures, network/API errors, partial
results, and retryable conditions.

## Security and privacy

Describe how credentials and sensitive data are supplied, used, logged, and stored. Request only the permissions
needed for the specified behavior.

## Acceptance criteria

- [ ] <Observable condition that demonstrates the feature works>
- [ ] <Observable condition for an important failure or edge case>

## Dependencies and relationships

List related specs and explain the contract between them. Avoid duplicating requirements owned by another spec.

## Open questions

- <Decision that must be resolved before implementation>
