# Repository instructions

These instructions apply to future work in this repository. Follow the user's request when it explicitly narrows or overrides them.

## Working in this repository

- For requests that change functionality, update the relevant specification as well as the Go code unless the user asks for spec-only work, says not to implement, or otherwise narrows the task.
- Once the user has reviewed and agreed on a package API, implement against that API without reopening the API review unless implementation reveals a concrete blocker or the user asks to revisit it.
- Generate or implement one Go file at a time, then stop for the user's review before moving to another file. Apply the user's code review comments to that file before continuing; treat those comments as durable guidance for later generated code.
- Before presenting each generated or implemented file for review, check it against the agreed package API, applicable specifications, and all prior review feedback in this file. Do not carry unresolved review comments into the next file.
- Do not run tests locally. GitHub CI runs them. Use `gofmt` on Go files that you change.
- Do not post PR comments or messages on the user's behalf. Give review feedback in the CLI/conversation. Do not check CI status when reviewing a PR.
- Keep errors and logs specific about what failed, include identifying context, and use a severity appropriate to the issue.
- Keep comments only when they explain a non-obvious constraint or behavior that readers need to understand. Remove comments that merely narrate the code.

## Go style

- Put each file-level constant declaration on its own line. Keep constants of the same type alphabetically ordered, and place enum constants beside their enum type rather than in a large file-level block.
- As much as practical, put behavior that validates, compares, or operates on a value on that object's type as a pointer-receiver method, rather than a standalone function that takes the object as an argument. Define the method beside the type; use names such as `Equals` for comparisons.
- At the beginning of a method, declare variables one per line and alphabetically order declarations of the same type.
- Loop index variables do not need comments explaining their position or iteration range. For other string or slice indexing at an index other than 0 or the last valid index, add a short nearby comment explaining what value or portion that index is expected to select. For parsing, identify the expected token or byte offset.
- Whenever using `strings.Fields`, add a very short nearby comment describing the expected fields.
- Centralize plain string values in constants, using a shared package when the value is shared. Keep string literals in Go code to the minimum required by the language or to initialize a variable, log message, or error message; required syntax such as struct tags is an exception.
- Give enum types their own validity method (for example, `IsValid`) rather than repeating lists of accepted values in compound conditionals.
- Normalize case whenever converting an input string into a string-backed enum, including parser and JSON decoding paths.
- When adding or changing tests, use Go's table-driven test pattern for cases that exercise the same behavior. Do not run those tests locally.

## Current design and layout decisions

- The flaky test component is called an **Analyzer**, rather than a Runner.
- Keep `internal/common` dependency-free so packages can share common types without importing higher-level functionality.
- In test results, keep report data and persistence in `report.go`; keep assessment, classification, outcome history, and rate calculations together in `assessment.go` because those APIs overlap. A `Report` contains `Workflow` values; workflows contain jobs, jobs contain steps, and tests belong under `step.tests` whenever the source step is known.
- Keep `Test.FailureMessage` populated only for failures.
- Parallel analyzer execution is a future requirement; it does not need to shape the current implementation.
- When adding a new style or layout preference, append it under the running feedback list below and apply it to future work.

## Running list of user feedback on style and layout

- Review and agree on package APIs before implementing them; then implement one Go file at a time for review.
- Prefer clear names that describe their role: use `Analyzer` instead of `Runner`, and use `PRNumber` instead of an ambiguous `Number` field.
- Put tests beneath their known source step so the report hierarchy communicates that workflows contain jobs, jobs contain steps, and steps contain tests.
- Keep report persistence in `report.go` and assessment/classification rules with history/rate calculations in `assessment.go` because those APIs overlap.
- Keep enum constants next to their corresponding type, and use a type method to validate enum strings.
- Store the failure message on each test observation, and populate it only when that observation failed.
- Put shared types in a dependency-free common package. Treat parallel analyzer execution as a future TODO rather than a current implementation constraint.
- Keep specifications in sync with implementation requests unless the user explicitly asks for spec-only work or says not to implement.
- After package APIs have been reviewed, proceed to code generation and file-by-file review; do not repeat the package API review as a prerequisite.
- Treat user comments on generated code as required review feedback: incorporate them before generating the next file, and add durable style or layout preferences to this list.
- Review each generated file against the agreed API and applicable repository instructions before asking the user to review it.
- Prefer pointer-receiver methods on the object being validated, compared, or operated on over standalone functions that take that object as an argument; keep those methods beside their type when practical.
- Explain interior string and slice indexes with a short comment stating the expected element or parsed offset; this applies to string and slice parsing as well as ordinary indexing. Loop index variables are exempt when the index only identifies the current iteration position.
- Add a very short comment at each `strings.Fields` call describing its expected field sequence.
- Normalize input case when converting strings into enums, including parser and JSON inputs.
- Normalize input case before converting strings to enums, including values decoded from JSON.
