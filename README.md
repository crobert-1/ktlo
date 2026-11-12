# ktlo

`ktlo` runs enabled checks against repositories listed in a YAML configuration file. PR status results are written as
Markdown under `results/github/<owner>/<repository>/RESULTS.md`.

## Prerequisite

Set `GH_TOKEN` to a GitHub token that can read the target repository and its pull requests. `GITHUB_TOKEN` is also
accepted as a fallback.

```sh
export GH_TOKEN="your-github-token"
```

## Build and run

```sh
make ktlo
cp config.example.yaml config.yaml
./bin/ktlo --config config.yaml
```

## Run modes

The configuration selects which checks run. Each enabled check runs for every configured repository, and both checks
can run in the same invocation.

- **PR status** (`prStatus.enabled`, enabled by default) checks open dependency-update pull requests and their CI status.
  It writes a Markdown report to `results/github/<owner>/<repository>/RESULTS.md`.
- **Failing test analysis** (`failingTestCheck.enabled`, disabled by default) collects workflow runs within the configured
  lookback, parses Go test logs, and records test outcomes and flaky-test evidence. It writes a JSON report to
  `results/github/<owner>/<repository>/flaky-test-results.json`.

Both checks use the configured GitHub repositories and require a token from `GH_TOKEN` or its `GITHUB_TOKEN` fallback.

## Configuration

The YAML file passed with `--config` selects which checks to run, their settings, and the repositories to check. The
same options are shown in [config.example.yaml](config.example.yaml):

```yaml
prStatus:
  enabled: true

failingTestCheck:
  enabled: false
  lookback: 7d
  maxArtifactAgeHours: 48

repositories:
  - https://github.com/owner/repository
```

- `repositories` (default: none): GitHub repository URLs to check. At least one is required when a check is enabled.
- `prStatus`
  - `enabled` (default: `true`): Runs the PR status check.
- `failingTestCheck`
  - `enabled` (default: `false`): Runs the flaky-test analysis check.
  - `lookback` (default: `7d`): How far back to include workflow attempts. Accepts Go duration units such as `30m` and `12h`, plus `d` for days. Must be positive.
  - `maxArtifactAgeHours` (default: `48`): Intended freshness limit, in hours, for consuming a flaky-analysis report. `0` disables the limit; the current CLI does not yet enforce this option.

Unknown YAML keys are rejected. `failingTestCheck` is opt-in; `prStatus` runs by default unless explicitly disabled.
