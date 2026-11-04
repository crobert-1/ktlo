# PR Status

## Overview

The goal of this check is to review open PRs in the target repository and inform the user of their status,
especially which PRs are ready to merge, and which PRs require attention.

## Functionality

The PR Status check will perform the following actions:

### Dependabot and Renovate PRs

#### Implementation

The script will check for open Dependabot and Renovate PRs in the target repository. If any are found, it will:

- Check if they are ready to merge:
  - Check if their CI checks have passed.
  - Check if they have been approved by a reviewer.
  - If the PRs have passed CI and have been approved, the script will inform the user that the PRs are ready to merge.

- If the PR's CI is blocked from running:
  - The script will inform the user that the PRs need their CI to be approved to run.

- If the PR has a merge conflict, the script will inform the user the PR needs to be re-generated to resolve the merge conflict.

- If the PR has failed CI, use the shared functionality to check the failure logs for known flaky test patterns,
  and determine if the PR's CI failures are due to flaky tests. If the failures are due to flaky tests, the script will
  inform the user that the PR's CI failures are due to flaky tests.

#### Output

- Present the results by PR status (determined by the script):
  - Ready to merge
  - CI blocked
  - Merge conflict
  - CI failed due to flaky tests
  - CI failed due to non-flaky tests

- PRs should be listed with following information (and format):
  - PR author: PR Title (PR URL that is clickable)