# Dependency Bump PR Status

## Failed

<table>
<thead>
<tr><th>Author</th><th>Pull Request</th><th>CI/CD</th></tr>
</thead>
<tbody>
<tr><td>srv-gh-o11y-gdi</td><td><a href="https://github.com/signalfx/splunk-otel-collector/pull/8458">Update splunk-otel-js to latest (#8458)</a></td><td><a href="#pr-8458">FAILURE</a> — 376 passed, 6 failed, 15 skipped, 19 cancelled</td></tr>
<tr><td>dependabot</td><td><a href="https://github.com/signalfx/splunk-otel-collector/pull/8442">Bump github.com/gobwas/glob from 0.2.4-0.20181002190808-e7a84e9525fe to 1.1.0 in /internal/signalfx-agent (#8442)</a></td><td><a href="#pr-8442">FAILURE</a> — 23 passed, 18 failed, 40 skipped, 28 cancelled</td></tr>
<tr><td>dependabot</td><td><a href="https://github.com/signalfx/splunk-otel-collector/pull/8440">Bump github.com/prometheus/prometheus from 0.315.1 to 0.316.0 (#8440)</a></td><td><a href="#pr-8440">FAILURE</a> — 17 passed, 8 failed, 23 skipped, 23 cancelled</td></tr>
<tr><td>dependabot</td><td><a href="https://github.com/signalfx/splunk-otel-collector/pull/8397">Bump github.com/hashicorp/vault-plugin-auth-gcp from 0.23.2-0.20260604163449-108858b5ffea to 0.24.2 (#8397)</a></td><td><a href="#pr-8397">FAILURE</a> — 11 passed, 12 failed, 25 skipped, 23 cancelled</td></tr>
<tr><td>dependabot</td><td><a href="https://github.com/signalfx/splunk-otel-collector/pull/8076">Bump ruby from 3.3.8-bookworm to 4.0.0-bookworm in /deployments/puppet (#8076)</a></td><td><a href="#pr-8076">FAILURE</a> — 87 passed, 2 failed, 15 skipped</td></tr>
<tr><td>dependabot</td><td><a href="https://github.com/signalfx/splunk-otel-collector/pull/7879">Bump actionshub/chef-install from 5.0.2 to 6.0.0 (#7879)</a></td><td><a href="#pr-7879">FAILURE</a> — 102 passed, 8 failed, 13 skipped</td></tr>
</tbody>
</table>

## Stale

<table>
<thead>
<tr><th>Author</th><th>Pull Request</th></tr>
</thead>
<tbody>
<tr><td>srv-gh-o11y-gdi</td><td><a href="https://github.com/signalfx/splunk-otel-collector/pull/8420">Update splunk-otel-dotnet to latest (#8420)</a></td></tr>
<tr><td>dependabot</td><td><a href="https://github.com/signalfx/splunk-otel-collector/pull/8415">Bump ansible/ansible-lint from 26.10.1 to 26.11.0 (#8415)</a></td></tr>
</tbody>
</table>

<a id="pr-8458"></a>
### [Update splunk-otel-js to latest (#8458)](https://github.com/signalfx/splunk-otel-collector/pull/8458)

Partial — the report is incomplete; findings use only available observations.

<table>
<thead>
<tr><th>Failure</th><th>Classification</th><th>Latest failure</th><th>Main history</th><th>PR commit history</th><th>Related CI history</th></tr>
</thead>
<tbody>
<tr><td>build-and-test / cross-compile (binaries-freebsd_amd64) / compile / Run ./.github/actions/setup-environment</td><td>Broken</td><td><pre><code>go: downloading golang.org/x/tools v0.52.0
##[error]/home/runner/go/pkg/mod/github.com/golangci/golangci-lint/v2@v2.7.2/internal/cache/cache.go:15:2: golang.org/x/tools@v0.52.0: read &#34;https://proxy.golang.org/golang.org/x/tools/@v/v0.52.0.zip&#34;: stream error: stream ID 35; INTERNAL_ERROR; received from peer
make: *** [Makefile:246: install-tools] Error 1</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264315482/job/134669639135">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>changelog / changelog / Validate ./.chloggen/*.yaml changes</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>2026-11-18T08:10:00.3170582Z chloggen validate
2026-11-18T08:10:00.3196640Z Error: specify one or more issues #&#39;s
2026-11-18T08:10:00.3197107Z Usage:
2026-11-18T08:10:00.3197582Z Error: specify one or more issues #&#39;s
2026-11-18T08:10:00.3197978Z   chloggen validate [flags]
2026-11-18T08:10:00.3200957Z
2026-11-18T08:10:00.3203216Z make: *** [Makefile:431: chlog-validate] Error 1
2026-11-18T08:10:00.3206110Z New ./.chloggen/*.yaml file failed validation.
2026-11-18T08:10:00.3222307Z ##[error]Process completed with exit code 1.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264315473/job/134664967586">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>test github.com/signalfx/splunk-otel-collector/tests/puppet.TestWindowsCustom</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>    puppet_windows_test.go:181: 
        	Error Trace:	D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:55
        	            				D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:67
        	            				D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:181
        	Error:      	[]int{0} does not contain 1
        	Test:       	TestWindowsCustom
        	Messages:   	choco upgrade -y -f puppet-agent --version 8.10.0:</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532182">Failure run</a></td><td>—</td><td><a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/4 passed (0%), 4/4 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532182">create-pull-request/update-nodejs-agent 28fbba8 failure</a></td><td>puppet-test — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2022, 8.10.0, default) — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532097">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2022, 8.10.0, default) / Test puppet deployment — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 1/1 passed (100%), 0/1 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532097">create-pull-request/update-nodejs-agent 28fbba8 pass</a><br>puppet-test / puppet-test-windows (windows-2022, 8.10.0, default) / Re-run failed tests — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532097">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2025, 8.10.0, default) — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532140">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2025, 8.10.0, default) / Test puppet deployment — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 1/1 passed (100%), 0/1 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532140">create-pull-request/update-nodejs-agent 28fbba8 pass</a><br>puppet-test / puppet-test-windows (windows-2025, 8.10.0, default) / Re-run failed tests — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532140">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2025, 8.10.0, custom_vars) — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532182">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2025, 8.10.0, custom_vars) / Test puppet deployment — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 1/1 passed (100%), 0/1 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532182">create-pull-request/update-nodejs-agent 28fbba8 pass</a><br>puppet-test / puppet-test-windows (windows-2025, 8.10.0, custom_vars) / Re-run failed tests — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532182">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2022, 8.10.0, custom_vars) — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532246">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2022, 8.10.0, custom_vars) / Test puppet deployment — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 1/1 passed (100%), 0/1 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532246">create-pull-request/update-nodejs-agent 28fbba8 pass</a><br>puppet-test / puppet-test-windows (windows-2022, 8.10.0, custom_vars) / Re-run failed tests — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532246">create-pull-request/update-nodejs-agent 28fbba8 failure</a></td></tr>
<tr><td>test github.com/signalfx/splunk-otel-collector/tests/puppet.TestWindowsDefault</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>    puppet_windows_test.go:165: 
        	Error Trace:	D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:55
        	            				D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:67
        	            				D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:165
        	Error:      	[]int{0} does not contain 1
        	Test:       	TestWindowsDefault
        	Messages:   	choco upgrade -y -f puppet-agent --version 8.10.0:</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532097">Failure run</a></td><td>—</td><td><a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/4 passed (0%), 4/4 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532097">create-pull-request/update-nodejs-agent 28fbba8 failure</a></td><td>puppet-test — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2022, 8.10.0, default) — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532097">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2022, 8.10.0, default) / Test puppet deployment — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 1/1 passed (100%), 0/1 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532097">create-pull-request/update-nodejs-agent 28fbba8 pass</a><br>puppet-test / puppet-test-windows (windows-2022, 8.10.0, default) / Re-run failed tests — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532097">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2025, 8.10.0, default) — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532140">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2025, 8.10.0, default) / Test puppet deployment — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 1/1 passed (100%), 0/1 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532140">create-pull-request/update-nodejs-agent 28fbba8 pass</a><br>puppet-test / puppet-test-windows (windows-2025, 8.10.0, default) / Re-run failed tests — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532140">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2025, 8.10.0, custom_vars) — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532182">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2025, 8.10.0, custom_vars) / Test puppet deployment — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 1/1 passed (100%), 0/1 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532182">create-pull-request/update-nodejs-agent 28fbba8 pass</a><br>puppet-test / puppet-test-windows (windows-2025, 8.10.0, custom_vars) / Re-run failed tests — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532182">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2022, 8.10.0, custom_vars) — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532246">create-pull-request/update-nodejs-agent 28fbba8 failure</a><br>puppet-test / puppet-test-windows (windows-2022, 8.10.0, custom_vars) / Test puppet deployment — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 1/1 passed (100%), 0/1 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532246">create-pull-request/update-nodejs-agent 28fbba8 pass</a><br>puppet-test / puppet-test-windows (windows-2022, 8.10.0, custom_vars) / Re-run failed tests — <a href="https://github.com/signalfx/splunk-otel-collector/commit/28fbba803bef46bb36dc6b5642feb3f1d9b2b72f">create-pull-request/update-nodejs-agent 28fbba8</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532246">create-pull-request/update-nodejs-agent 28fbba8 failure</a></td></tr>
<tr><td>puppet-test / puppet-test-windows (windows-2022, 8.10.0, custom_vars) / Re-run failed tests</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>=== RUN   TestWindowsCustom
    puppet_windows_test.go:181: 
        	Error Trace:	D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:55
        	            				D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:67
        	            				D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:181
        	Error:      	[]int{0} does not contain 1
        	Test:       	TestWindowsCustom
        	Messages:   	choco upgrade -y -f puppet-agent --version 8.10.0:</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532246">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>puppet-test / puppet-test-windows (windows-2022, 8.10.0, default) / Re-run failed tests</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>=== RUN   TestWindowsDefault
    puppet_windows_test.go:165: 
        	Error Trace:	D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:55
        	            				D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:67
        	            				D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:165
        	Error:      	[]int{0} does not contain 1
        	Test:       	TestWindowsDefault
        	Messages:   	choco upgrade -y -f puppet-agent --version 8.10.0:</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532097">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>puppet-test / puppet-test-windows (windows-2025, 8.10.0, custom_vars) / Re-run failed tests</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>=== RUN   TestWindowsCustom
    puppet_windows_test.go:181: 
        	Error Trace:	D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:55
        	            				D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:67
        	            				D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:181
        	Error:      	[]int{0} does not contain 1
        	Test:       	TestWindowsCustom
        	Messages:   	choco upgrade -y -f puppet-agent --version 8.10.0:</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532182">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>puppet-test / puppet-test-windows (windows-2025, 8.10.0, default) / Re-run failed tests</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>=== RUN   TestWindowsDefault
    puppet_windows_test.go:165: 
        	Error Trace:	D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:55
        	            				D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:67
        	            				D:/a/splunk-otel-collector/splunk-otel-collector/tests/puppet/puppet_windows_test.go:165
        	Error:      	[]int{0} does not contain 1
        	Test:       	TestWindowsDefault
        	Messages:   	choco upgrade -y -f puppet-agent --version 8.10.0:</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44264320420/job/134672532140">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
</tbody>
</table>


<a id="pr-8442"></a>
### [Bump github.com/gobwas/glob from 0.2.4-0.20181002190808-e7a84e9525fe to 1.1.0 in /internal/signalfx-agent (#8442)](https://github.com/signalfx/splunk-otel-collector/pull/8442)

Partial — the report is incomplete; findings use only available observations.
- Could not parse a failed workflow: could not parse workflow run 43659653477 attempt 1: could not read logs for workflow run 43659653477 attempt 1: could not read GitHub workflow attempt log archive for run 43659653477 attempt 1: net/http: request canceled (Client.Timeout or context cancellation while reading body)<br>workflow run 43659653477 attempt 1 job "ptp-receiver-integration-test" step "Run PTP receiver integration test" has no downloadable log<br>workflow run 43659653477 attempt 1 job "otelcol (amd64)" step "Run make binaries-linux_amd64 COVER_TESTING=true" has no downloadable log

<table>
<thead>
<tr><th>Failure</th><th>Classification</th><th>Latest failure</th><th>Main history</th><th>PR commit history</th><th>Related CI history</th></tr>
</thead>
<tbody>
<tr><td>Integration Test / otelcol (amd64) / Run make binaries-linux_amd64 COVER_TESTING=true</td><td>Broken</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653477/job/132645878537">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>Integration Test / ptp-receiver-integration-test / Run PTP receiver integration test</td><td>Broken</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653477/job/132645878373">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>auto-instrumentation / cross-compile (binaries-linux_amd64) / compile / Build Collector</td><td>Broken</td><td><pre><code>2026-11-15T00:49:18.3216907Z fatal: No names found, cannot describe anything.
2026-11-15T00:49:18.3231876Z fatal: No names found, cannot describe anything.
2026-11-15T00:49:18.3273770Z go generate ./...
2026-11-15T00:49:34.7963454Z # go.opentelemetry.io/collector/confmap/internal
2026-11-15T00:49:34.8050621Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:25:19: undefined: glob.Glob
2026-11-15T00:49:34.8143805Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:63:39: undefined: glob.Glob
2026-11-15T00:49:38.1590207Z # github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs
2026-11-15T00:49:38.1614535Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_all_matches.go:47:27: undefined: glob.Glob
2026-11-15T00:49:38.1649081Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_match.go:46:27: undefined: glob.Glob
2026-11-15T00:50:36.9321516Z # github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/translation/dpfilters
2026-11-15T00:50:36.9362745Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter@v0.165.1/internal/translation/dpfilters/matching.go:55:15: undefined: glob.Glob
2026-11-15T00:53:25.4866190Z # github.com/signalfx/signalfx-agent/pkg/utils/filter
2026-11-15T00:53:25.4944741Z ##[error]internal/signalfx-agent/pkg/utils/filter/matching.go:65:15: undefined: glob.Glob
2026-11-15T00:53:33.4414313Z make[1]: *** [Makefile:266: otelcol] Error 1
2026-11-15T00:53:33.4414681Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-15T00:53:33.4417844Z make: *** [Makefile:344: binaries-linux_amd64] Error 2
2026-11-15T00:53:33.4433097Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653683/job/132645878812">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/cmd/monitorcodegen</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/core</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/core/common/auth</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/core/common/constants</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/core/common/docker</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/core/common/dpmeta</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/core/common/ecs</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/core/common/httpclient</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/core/common/kubelet</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/core/common/kubernetes</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/core/config</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/core/config/validation</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/core/dpfilters</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/coredns</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/cpu</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/diskio</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/docker</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878844">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/ecs</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/elasticsearch/client</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/elasticsearch/stats</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/elasticsearch/utils</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/etcd</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/filesystems</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/gitlab</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/http</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/internalmetrics</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/kubernetes</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/kubernetes/apiserver</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/kubernetes/controllermanager</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/kubernetes/events</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/kubernetes/leadership</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/kubernetes/proxy</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/kubernetes/scheduler</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/kubernetes/utils</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/kubernetes/volumes</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/load</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/logstash/tcp</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/memory</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/process</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/processlist</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/go</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/nginxingress</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/nginxvts</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/node</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/postgres</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/prometheus</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/redis</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/velero</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/prometheusexporter</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/sql</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/supervisor</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/telegraf/common/accumulator</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/telegraf/common/emitter</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/telegraf/common/emitter/baseemitter</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/telegraf/common/parser</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/telegraf/monitors/exec</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/telegraf/monitors/winperfcounters</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/traefik</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/types</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/vmem</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/vsphere</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/vsphere/model</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/vsphere/service</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/windowsiis</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/monitors/zcodegen</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/neotest</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/neotest/k8s/testhelpers/fakek8s</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/utils</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/utils/filter</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/utils/hostfs</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/utils/k8sutil</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / package github.com/signalfx/signalfx-agent/pkg/utils/timeutil</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / coverage / Run Unit Tests With Coverage</td><td>Broken</td><td><pre><code>2026-11-15T01:00:26.9523918Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/go [build failed]
2026-11-15T01:00:26.9524524Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/nginxingress [build failed]
2026-11-15T01:00:26.9525169Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/nginxvts [build failed]
2026-11-15T01:00:26.9525772Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/node [build failed]
2026-11-15T01:00:26.9526403Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/postgres [build failed]
2026-11-15T01:00:26.9527012Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/prometheus [build failed]
2026-11-15T01:00:26.9527844Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/redis [build failed]
2026-11-15T01:00:26.9528519Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/prometheus/velero [build failed]
2026-11-15T01:00:26.9529343Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/prometheusexporter [build failed]
2026-11-15T01:00:26.9529977Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/sql [build failed]
2026-11-15T01:00:26.9530501Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/supervisor [build failed]
2026-11-15T01:00:26.9531109Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/telegraf/common/accumulator [build failed]
2026-11-15T01:00:26.9531760Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/telegraf/common/emitter [build failed]
2026-11-15T01:00:26.9532444Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/telegraf/common/emitter/baseemitter [build failed]
2026-11-15T01:00:26.9533129Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/telegraf/common/parser [build failed]
2026-11-15T01:00:26.9533749Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/telegraf/monitors/exec [build failed]
2026-11-15T01:00:26.9534717Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/telegraf/monitors/winperfcounters [build failed]
2026-11-15T01:00:26.9535358Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/traefik [build failed]
2026-11-15T01:00:26.9536015Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/types [build failed]
2026-11-15T01:00:26.9536516Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/vmem [build failed]
2026-11-15T01:00:26.9537013Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/vsphere [build failed]
2026-11-15T01:00:26.9537547Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/vsphere/model [build failed]
2026-11-15T01:00:26.9538114Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/vsphere/service [build failed]
2026-11-15T01:00:26.9538991Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/windowsiis [build failed]
2026-11-15T01:00:26.9539686Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/zcodegen [build failed]
2026-11-15T01:00:26.9540488Z FAIL	github.com/signalfx/signalfx-agent/pkg/neotest [build failed]
2026-11-15T01:00:26.9541074Z FAIL	github.com/signalfx/signalfx-agent/pkg/neotest/k8s/testhelpers/fakek8s [build failed]
2026-11-15T01:00:27.2670144Z FAIL	github.com/signalfx/signalfx-agent/pkg/utils [build failed]
2026-11-15T01:00:27.2670899Z FAIL	github.com/signalfx/signalfx-agent/pkg/utils/filter [build failed]
2026-11-15T01:00:27.2671432Z FAIL	github.com/signalfx/signalfx-agent/pkg/utils/hostfs [build failed]
2026-11-15T01:00:27.2671950Z FAIL	github.com/signalfx/signalfx-agent/pkg/utils/k8sutil [build failed]
2026-11-15T01:00:27.2672460Z FAIL	github.com/signalfx/signalfx-agent/pkg/utils/timeutil [build failed]
2026-11-15T01:00:27.2672827Z FAIL
2026-11-15T01:00:27.7197826Z make[2]: *** [Makefile:80: test-with-codecov] Error 1
2026-11-15T01:00:27.7200126Z make[1]: *** [Makefile:82: internal/signalfx-agent] Error 2
2026-11-15T01:00:27.7201043Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-15T01:00:27.7205075Z make: *** [Makefile:218: gotest-with-codecov] Error 2
2026-11-15T01:00:27.7229376Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878867">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / cross-compile (binaries-darwin_amd64) / compile / Build Collector</td><td>Broken</td><td><pre><code>2026-11-15T00:46:41.5996459Z go: downloading github.com/dvsekhvalnov/jose2go v1.8.0
2026-11-15T00:46:42.1365552Z fatal: No names found, cannot describe anything.
2026-11-15T00:46:42.1381063Z fatal: No names found, cannot describe anything.
2026-11-15T00:46:42.1422013Z go generate ./...
2026-11-15T00:46:58.5326436Z # go.opentelemetry.io/collector/confmap/internal
2026-11-15T00:46:58.5353961Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:25:19: undefined: glob.Glob
2026-11-15T00:46:58.5394611Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:63:39: undefined: glob.Glob
2026-11-15T00:47:01.3736520Z # github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs
2026-11-15T00:47:01.3746028Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_all_matches.go:47:27: undefined: glob.Glob
2026-11-15T00:47:01.3804370Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_match.go:46:27: undefined: glob.Glob
2026-11-15T00:47:59.3080462Z # github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/translation/dpfilters
2026-11-15T00:47:59.3098372Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter@v0.165.1/internal/translation/dpfilters/matching.go:55:15: undefined: glob.Glob
2026-11-15T00:50:39.0534098Z # github.com/signalfx/signalfx-agent/pkg/utils/filter
2026-11-15T00:50:39.0558769Z ##[error]internal/signalfx-agent/pkg/utils/filter/matching.go:65:15: undefined: glob.Glob
2026-11-15T00:50:45.9542708Z make[1]: *** [Makefile:266: otelcol] Error 1
2026-11-15T00:50:45.9543101Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-15T00:50:45.9545787Z make: *** [Makefile:332: binaries-darwin_amd64] Error 2
2026-11-15T00:50:45.9560028Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645879309">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / latest-dockerd-test / dockerd test</td><td>Broken</td><td><pre><code>2026-11-15T00:47:19.7576056Z # github.com/signalfx/signalfx-agent/pkg/utils/filter
2026-11-15T00:47:19.7663236Z ##[error]../../utils/filter/matching.go:65:15: undefined: glob.Glob
2026-11-15T00:47:28.1554538Z FAIL	github.com/signalfx/signalfx-agent/pkg/monitors/docker [build failed]
2026-11-15T00:47:28.1555310Z FAIL
2026-11-15T00:47:28.2028620Z ##[error]Process completed with exit code 1.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878844">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / lint / Lint</td><td>Flaky</td><td><pre><code>2026-11-15T00:56:09.4610695Z running misspell -error
2026-11-15T00:56:09.5047643Z golangci-lint run --allow-parallel-runners -j4
2026-11-15T00:56:21.4146727Z Check License finished successfully
2026-11-15T00:56:21.4151148Z running misspell -error
2026-11-15T00:56:21.4377542Z golangci-lint run --allow-parallel-runners -j4
2026-11-15T00:56:32.0078509Z Check License finished successfully
2026-11-15T00:56:32.0083118Z running misspell -error
2026-11-15T00:56:32.0318139Z golangci-lint run --allow-parallel-runners -j4
2026-11-15T00:56:40.8313498Z Check License finished successfully
2026-11-15T00:56:40.8323028Z running misspell -error
2026-11-15T00:56:40.8553186Z golangci-lint run --allow-parallel-runners -j4
2026-11-15T00:56:42.3403774Z Check License finished successfully
2026-11-15T00:56:42.3409252Z running misspell -error
2026-11-15T00:56:42.3678399Z golangci-lint run --allow-parallel-runners -j4
2026-11-15T00:56:45.3293993Z Check License finished successfully
2026-11-15T00:56:45.3329851Z running misspell -error
2026-11-15T00:56:45.3536330Z golangci-lint run --allow-parallel-runners -j4
2026-11-15T00:56:46.7368970Z Check License finished successfully
2026-11-15T00:56:46.7374968Z running misspell -error
2026-11-15T00:56:46.7634535Z golangci-lint run --allow-parallel-runners -j4
2026-11-15T00:57:00.1549240Z Check License finished successfully
2026-11-15T00:57:00.1555007Z running misspell -error
2026-11-15T00:57:00.1811245Z Skip linting for smartagentextension
2026-11-15T00:57:00.4933019Z Check License finished successfully
2026-11-15T00:57:00.4937834Z running misspell -error
2026-11-15T00:57:00.5254263Z Skip linting for smartagentreceiver
2026-11-15T00:57:00.5981606Z Check License finished successfully
2026-11-15T00:57:00.5991734Z running misspell -error
2026-11-15T00:57:00.6316118Z golangci-lint run --allow-parallel-runners -j4
2026-11-15T00:57:08.6211660Z Check License finished successfully
2026-11-15T00:57:08.6219991Z running misspell -error
2026-11-15T00:57:08.6548250Z golangci-lint run --allow-parallel-runners -j4
2026-11-15T00:58:00.8864322Z Check License finished successfully
2026-11-15T00:58:00.8871274Z running misspell -error
2026-11-15T00:58:00.9638874Z golangci-lint run --allow-parallel-runners -j4
2026-11-15T00:58:47.6495206Z Check License finished successfully
2026-11-15T00:58:47.6501754Z running misspell -error
2026-11-15T00:58:47.6817706Z golangci-lint run --allow-parallel-runners -j4
2026-11-15T00:59:01.3765362Z Check License finished successfully
2026-11-15T00:59:01.3771045Z running misspell -error
2026-11-15T00:59:01.3997903Z golangci-lint run --allow-parallel-runners -j4
2026-11-15T00:59:04.1116964Z Check License finished successfully
2026-11-15T00:59:04.1121196Z running misspell -error
2026-11-15T00:59:04.2477466Z golangci-lint run --allow-parallel-runners -j4
2026-11-15T01:01:06.9756886Z ##[error]internal/confmapprovider/discovery/internal/utils.go:17:8: could not import go.opentelemetry.io/collector/confmap (/home/runner/go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/confmap.go:9:2: could not import go.opentelemetry.io/collector/confmap/internal (-: # go.opentelemetry.io/collector/confmap/internal
2026-11-15T01:01:06.9929369Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:25:19: undefined: glob.Glob
2026-11-15T01:01:06.9942531Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:63:39: undefined: glob.Glob)) (typecheck)
2026-11-15T01:01:06.9945549Z import &#34;go.opentelemetry.io/collector/confmap&#34;
2026-11-15T01:01:06.9945981Z        ^
2026-11-15T01:01:07.7669980Z make[2]: *** [Makefile.Common:155: lint] Error 1
2026-11-15T01:01:07.7673242Z make[1]: *** [Makefile:82: /home/runner/work/splunk-otel-collector/splunk-otel-collector] Error 2
2026-11-15T01:01:07.7677744Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-15T01:01:07.7714768Z make: *** [Makefile:235: lint-all] Error 2
2026-11-15T01:01:07.7855427Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645878715">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / test-ubuntu-24.04-race_detection:false-coverage:true / unit-test / bundle.d</td><td>Broken</td><td><pre><code>2026-11-15T00:46:05.3165790Z # go.opentelemetry.io/collector/confmap/internal
2026-11-15T00:46:05.3337500Z ##[error]/home/runner/go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:25:19: undefined: glob.Glob
2026-11-15T00:46:05.3350332Z ##[error]/home/runner/go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:63:39: undefined: glob.Glob
2026-11-15T00:46:11.3756783Z # github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs
2026-11-15T00:46:11.3899352Z ##[error]/home/runner/go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_all_matches.go:47:27: undefined: glob.Glob
2026-11-15T00:46:11.3989146Z ##[error]/home/runner/go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_match.go:46:27: undefined: glob.Glob
2026-11-15T00:47:44.5978036Z # github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/translation/dpfilters
2026-11-15T00:47:44.6009304Z ##[error]/home/runner/go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter@v0.165.1/internal/translation/dpfilters/matching.go:55:15: undefined: glob.Glob
2026-11-15T00:51:53.1718300Z # github.com/signalfx/signalfx-agent/pkg/utils/filter
2026-11-15T00:51:53.1746693Z ##[error]../../internal/signalfx-agent/pkg/utils/filter/matching.go:65:15: undefined: glob.Glob
2026-11-15T00:52:04.3686479Z make: *** [Makefile.Common:195: bundle.d] Error 1
2026-11-15T00:52:04.3697234Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653587/job/132645879326">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>installer-script-test / cross-compile (binaries-linux_amd64) / compile / Build Collector</td><td>Broken</td><td><pre><code>2026-11-15T00:46:56.1650394Z make[1]: Entering directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-15T00:46:56.8920279Z fatal: No names found, cannot describe anything.
2026-11-15T00:46:56.8936282Z fatal: No names found, cannot describe anything.
2026-11-15T00:46:56.8980956Z go generate ./...
2026-11-15T00:47:14.7071287Z # go.opentelemetry.io/collector/confmap/internal
2026-11-15T00:47:14.7142271Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:25:19: undefined: glob.Glob
2026-11-15T00:47:14.7152304Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:63:39: undefined: glob.Glob
2026-11-15T00:47:17.8828825Z # github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs
2026-11-15T00:47:17.8843628Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_all_matches.go:47:27: undefined: glob.Glob
2026-11-15T00:47:17.8871708Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_match.go:46:27: undefined: glob.Glob
2026-11-15T00:48:20.6407485Z # github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/translation/dpfilters
2026-11-15T00:48:20.6412304Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter@v0.165.1/internal/translation/dpfilters/matching.go:55:15: undefined: glob.Glob
2026-11-15T00:51:04.8557553Z # github.com/signalfx/signalfx-agent/pkg/utils/filter
2026-11-15T00:51:04.8603452Z ##[error]internal/signalfx-agent/pkg/utils/filter/matching.go:65:15: undefined: glob.Glob
2026-11-15T00:51:12.0567218Z make[1]: *** [Makefile:266: otelcol] Error 1
2026-11-15T00:51:12.0567561Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-15T00:51:12.0570321Z make: *** [Makefile:344: binaries-linux_amd64] Error 2
2026-11-15T00:51:12.0597868Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653760/job/132645879449">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>linux-package-test / cross-compile (binaries-linux_amd64) / compile / Build Collector</td><td>Broken</td><td><pre><code>2026-11-15T00:47:06.4858154Z make[1]: Entering directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-15T00:47:07.1448359Z fatal: No names found, cannot describe anything.
2026-11-15T00:47:07.1463007Z fatal: No names found, cannot describe anything.
2026-11-15T00:47:07.1506114Z go generate ./...
2026-11-15T00:47:23.4762822Z # go.opentelemetry.io/collector/confmap/internal
2026-11-15T00:47:23.4798677Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:25:19: undefined: glob.Glob
2026-11-15T00:47:23.4809879Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:63:39: undefined: glob.Glob
2026-11-15T00:47:27.1957159Z # github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs
2026-11-15T00:47:27.1959773Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_all_matches.go:47:27: undefined: glob.Glob
2026-11-15T00:47:27.1987605Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_match.go:46:27: undefined: glob.Glob
2026-11-15T00:48:31.7657316Z # github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/translation/dpfilters
2026-11-15T00:48:31.7702785Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter@v0.165.1/internal/translation/dpfilters/matching.go:55:15: undefined: glob.Glob
2026-11-15T00:51:25.1986927Z # github.com/signalfx/signalfx-agent/pkg/utils/filter
2026-11-15T00:51:25.2093134Z ##[error]internal/signalfx-agent/pkg/utils/filter/matching.go:65:15: undefined: glob.Glob
2026-11-15T00:51:33.7742165Z make[1]: *** [Makefile:266: otelcol] Error 1
2026-11-15T00:51:33.7742605Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-15T00:51:33.7745022Z make: *** [Makefile:344: binaries-linux_amd64] Error 2
2026-11-15T00:51:33.7770815Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653655/job/132645879048">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>otelcol-fips / otelcol-fips (linux, amd64) / Run make otelcol-fips</td><td>Broken</td><td><pre><code>2026-11-15T00:47:51.7111186Z # go.opentelemetry.io/collector/confmap/internal
2026-11-15T00:47:51.7217924Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:25:19: undefined: glob.Glob
2026-11-15T00:47:51.7227052Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:63:39: undefined: glob.Glob
2026-11-15T00:47:56.1398300Z # github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs
2026-11-15T00:47:56.1525977Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_all_matches.go:47:27: undefined: glob.Glob
2026-11-15T00:47:56.1529088Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_match.go:46:27: undefined: glob.Glob
2026-11-15T00:49:09.1428232Z # github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/translation/dpfilters
2026-11-15T00:49:09.1457256Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter@v0.165.1/internal/translation/dpfilters/matching.go:55:15: undefined: glob.Glob
2026-11-15T00:52:30.4112803Z # github.com/signalfx/signalfx-agent/pkg/utils/filter
2026-11-15T00:52:30.4158366Z ##[error]internal/signalfx-agent/pkg/utils/filter/matching.go:65:15: undefined: glob.Glob
2026-11-15T00:52:40.8668953Z make: *** [Makefile:411: otelcol-fips] Error 1
2026-11-15T00:52:40.8691918Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653504/job/132645878299">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>otelcol-fips / otelcol-fips (linux, arm64) / Run make otelcol-fips</td><td>Broken</td><td><pre><code>2026-11-15T00:45:25.4742149Z # go.opentelemetry.io/collector/confmap/internal
2026-11-15T00:45:25.4908206Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:25:19: undefined: glob.Glob
2026-11-15T00:45:25.4918253Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:63:39: undefined: glob.Glob
2026-11-15T00:45:28.4266800Z # github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs
2026-11-15T00:45:28.4298654Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_all_matches.go:47:27: undefined: glob.Glob
2026-11-15T00:45:28.4302033Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_match.go:46:27: undefined: glob.Glob
2026-11-15T00:46:27.0420067Z # github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/translation/dpfilters
2026-11-15T00:46:27.0481889Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter@v0.165.1/internal/translation/dpfilters/matching.go:55:15: undefined: glob.Glob
2026-11-15T00:49:01.3447922Z # github.com/signalfx/signalfx-agent/pkg/utils/filter
2026-11-15T00:49:01.3522367Z ##[error]internal/signalfx-agent/pkg/utils/filter/matching.go:65:15: undefined: glob.Glob
2026-11-15T00:49:09.2293499Z make: *** [Makefile:411: otelcol-fips] Error 1
2026-11-15T00:49:09.2305822Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653504/job/132645878281">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>otelcol-fips / otelcol-fips (windows, amd64) / Run make otelcol-fips</td><td>Broken</td><td><pre><code>2026-11-15T00:46:12.5583611Z # go.opentelemetry.io/collector/confmap/internal
2026-11-15T00:46:12.5653145Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:25:19: undefined: glob.Glob
2026-11-15T00:46:12.5666986Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:63:39: undefined: glob.Glob
2026-11-15T00:46:18.3893436Z # github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs
2026-11-15T00:46:18.3929354Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_all_matches.go:47:27: undefined: glob.Glob
2026-11-15T00:46:18.3964048Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_match.go:46:27: undefined: glob.Glob
2026-11-15T00:47:50.4038331Z # github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/translation/dpfilters
2026-11-15T00:47:50.4042647Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter@v0.165.1/internal/translation/dpfilters/matching.go:55:15: undefined: glob.Glob
2026-11-15T00:51:56.2378663Z # github.com/signalfx/signalfx-agent/pkg/utils/filter
2026-11-15T00:51:56.2418044Z ##[error]internal/signalfx-agent/pkg/utils/filter/matching.go:65:15: undefined: glob.Glob
2026-11-15T00:52:07.5009235Z make: *** [Makefile:417: otelcol-fips] Error 1
2026-11-15T00:52:07.5023310Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653504/job/132645878009">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>otlpintput / ta-v2-build-packages / Package Splunk Connect for OTLP</td><td>Broken</td><td><pre><code>2026-11-15T00:47:42.6825983Z # go.opentelemetry.io/collector/confmap/internal
2026-11-15T00:47:42.6968785Z ##[error]/home/runner/go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:25:19: undefined: glob.Glob
2026-11-15T00:47:42.6996240Z ##[error]/home/runner/go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:63:39: undefined: glob.Glob
2026-11-15T00:47:59.2329072Z make: *** [Makefile:24: splunk-connect-for-otlp/linux_x86_64/bin/splunk-connect-for-otlp] Error 1
2026-11-15T00:47:59.2337933Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653576/job/132645878721">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>parity-test / parity-test / Build otelcol</td><td>Broken</td><td><pre><code>2026-11-15T00:46:14.5623449Z # go.opentelemetry.io/collector/confmap/internal
2026-11-15T00:46:14.5739501Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:25:19: undefined: glob.Glob
2026-11-15T00:46:14.5750933Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:63:39: undefined: glob.Glob
2026-11-15T00:46:18.1715618Z # github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs
2026-11-15T00:46:18.1751742Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_all_matches.go:47:27: undefined: glob.Glob
2026-11-15T00:46:18.1796181Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_match.go:46:27: undefined: glob.Glob
2026-11-15T00:47:20.0563310Z # github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/translation/dpfilters
2026-11-15T00:47:20.0617776Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter@v0.165.1/internal/translation/dpfilters/matching.go:55:15: undefined: glob.Glob
2026-11-15T00:50:16.0065626Z # github.com/signalfx/signalfx-agent/pkg/utils/filter
2026-11-15T00:50:16.0149121Z ##[error]internal/signalfx-agent/pkg/utils/filter/matching.go:65:15: undefined: glob.Glob
2026-11-15T00:50:24.0295335Z make: *** [Makefile:266: otelcol] Error 1
2026-11-15T00:50:24.0309167Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653676/job/132645878747">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>ta-v2 / ta-v2-build-packages / Build collector binaries (linux &amp; windows)</td><td>Broken</td><td><pre><code>2026-11-15T00:47:44.2363958Z # go.opentelemetry.io/collector/confmap/internal
2026-11-15T00:47:44.2566088Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:25:19: undefined: glob.Glob
2026-11-15T00:47:44.2612874Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:63:39: undefined: glob.Glob
2026-11-15T00:47:49.7375736Z # github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs
2026-11-15T00:47:49.7396160Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_all_matches.go:47:27: undefined: glob.Glob
2026-11-15T00:47:49.7422087Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_match.go:46:27: undefined: glob.Glob
2026-11-15T00:49:13.8960052Z # github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/translation/dpfilters
2026-11-15T00:49:13.9010357Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter@v0.165.1/internal/translation/dpfilters/matching.go:55:15: undefined: glob.Glob
2026-11-15T00:53:05.4666853Z # github.com/signalfx/signalfx-agent/pkg/utils/filter
2026-11-15T00:53:05.4738166Z ##[error]internal/signalfx-agent/pkg/utils/filter/matching.go:65:15: undefined: glob.Glob
2026-11-15T00:53:15.7311715Z make: *** [Makefile:266: otelcol] Error 1
2026-11-15T00:53:15.7342147Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653790/job/132645879027">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>win-package-test / cross-compile (binaries-windows_arm64) / compile / Build Collector</td><td>Broken</td><td><pre><code>2026-11-15T00:47:08.1883821Z go: downloading github.com/go-asn1-ber/asn1-ber v1.5.8
2026-11-15T00:47:08.8056013Z fatal: No names found, cannot describe anything.
2026-11-15T00:47:08.8071932Z fatal: No names found, cannot describe anything.
2026-11-15T00:47:08.8115854Z go generate ./...
2026-11-15T00:47:27.9928469Z # go.opentelemetry.io/collector/confmap/internal
2026-11-15T00:47:27.9998003Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:25:19: undefined: glob.Glob
2026-11-15T00:47:28.0007831Z ##[error]../../../go/pkg/mod/go.opentelemetry.io/collector/confmap@v1.71.0/internal/merge.go:63:39: undefined: glob.Glob
2026-11-15T00:47:31.9536587Z # github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs
2026-11-15T00:47:31.9540106Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_all_matches.go:47:27: undefined: glob.Glob
2026-11-15T00:47:31.9543783Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl@v1.4.1/ottlfuncs/func_replace_match.go:46:27: undefined: glob.Glob
2026-11-15T00:48:42.9613229Z # github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/translation/dpfilters
2026-11-15T00:48:42.9754712Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter@v0.165.1/internal/translation/dpfilters/matching.go:55:15: undefined: glob.Glob
2026-11-15T00:51:51.9323211Z # github.com/signalfx/signalfx-agent/pkg/utils/filter
2026-11-15T00:51:51.9352867Z ##[error]internal/signalfx-agent/pkg/utils/filter/matching.go:65:15: undefined: glob.Glob
2026-11-15T00:52:00.7148108Z make[1]: *** [Makefile:266: otelcol] Error 1
2026-11-15T00:52:00.7148812Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-15T00:52:00.7151082Z make: *** [Makefile:378: binaries-windows_arm64] Error 2
2026-11-15T00:52:00.7170370Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653632/job/132645878839">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/cmd/discoverybundler</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/cmd/otelcol</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/cmd/otelcollauncher/health</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/cmd/splunk-connect-for-otlp</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/auth</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/components</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/configconverter</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/configsource</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/configsource/envvarconfigsource</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/configsource/etcd2configsource</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/configsource/includeconfigsource</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/configsource/splunksecretconfigsource</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/configsource/vaultconfigsource</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/configsource/zookeeperconfigsource</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/confmapprovider/configsource</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/confmapprovider/discovery</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/confmapprovider/discovery/internal</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/confmapprovider/discovery/properties</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/exporter/splunkmodinputexporter</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/extension/configsourcetelemetryextension</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/extension/diskqueuestorageextension</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/extension/splunkauthextension</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/receiver/discoveryreceiver</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/receiver/gnmireceiver</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/receiver/promqlreceiver</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/receiver/ptpreceiver</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/receiver/ptpreceiver/internal/metadata</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/settings</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / windows-test (windows-2022) / Unit tests with coverage</td><td>Broken</td><td><pre><code>2026-11-15T00:54:24.9382711Z FAIL	github.com/signalfx/splunk-otel-collector/internal/configsource/zookeeperconfigsource [build failed]
2026-11-15T00:54:24.9385256Z FAIL	github.com/signalfx/splunk-otel-collector/internal/confmapprovider/configsource [build failed]
2026-11-15T00:54:24.9387361Z FAIL	github.com/signalfx/splunk-otel-collector/internal/confmapprovider/discovery [build failed]
2026-11-15T00:54:24.9389488Z FAIL	github.com/signalfx/splunk-otel-collector/internal/confmapprovider/discovery/internal [build failed]
2026-11-15T00:54:24.9391633Z FAIL	github.com/signalfx/splunk-otel-collector/internal/confmapprovider/discovery/properties [build failed]
2026-11-15T00:54:24.9393649Z FAIL	github.com/signalfx/splunk-otel-collector/internal/exporter/splunkmodinputexporter [build failed]
2026-11-15T00:54:24.9397781Z ?   	github.com/signalfx/splunk-otel-collector/internal/exporter/splunkmodinputexporter/internal/metadata	[no test files]
2026-11-15T00:54:24.9398244Z FAIL	github.com/signalfx/splunk-otel-collector/internal/extension/configsourcetelemetryextension [build failed]
2026-11-15T00:54:24.9400009Z FAIL	github.com/signalfx/splunk-otel-collector/internal/extension/diskqueuestorageextension [build failed]
2026-11-15T00:54:24.9832421Z FAIL	github.com/signalfx/splunk-otel-collector/internal/extension/splunkauthextension [build failed]
2026-11-15T00:54:24.9837243Z ?   	github.com/signalfx/splunk-otel-collector/internal/extension/splunkauthextension/internal/metadata	[no test files]
2026-11-15T00:54:25.4646634Z === RUN   TestBaseDirameFromExecutable_Error
2026-11-15T00:54:25.4756810Z 2026/11/15 00:54:25 ERROR process image path unavailable
2026-11-15T00:54:25.4757056Z --- PASS: TestBaseDirameFromExecutable_Error (0.01s)
2026-11-15T00:54:26.0411483Z ok  	github.com/signalfx/splunk-otel-collector/internal/opampsupervisor/launcher	0.061s	coverage: 95.9% of statements
2026-11-15T00:54:26.3671083Z FAIL	github.com/signalfx/splunk-otel-collector/internal/receiver/discoveryreceiver [build failed]
2026-11-15T00:54:26.9495469Z === RUN   TestLogsBuilderAppendLogRecord
2026-11-15T00:54:27.4622299Z ok  	github.com/signalfx/splunk-otel-collector/internal/receiver/discoveryreceiver/statussources	0.035s	coverage: 93.5% of statements
2026-11-15T00:54:27.4646668Z FAIL	github.com/signalfx/splunk-otel-collector/internal/receiver/gnmireceiver [build failed]
2026-11-15T00:54:27.4650674Z ?   	github.com/signalfx/splunk-otel-collector/internal/receiver/gnmireceiver/internal/metadata	[no test files]
2026-11-15T00:54:27.4651516Z FAIL	github.com/signalfx/splunk-otel-collector/internal/receiver/promqlreceiver [build failed]
2026-11-15T00:54:27.4655185Z ?   	github.com/signalfx/splunk-otel-collector/internal/receiver/promqlreceiver/internal/metadata	[no test files]
2026-11-15T00:54:27.4655662Z FAIL	github.com/signalfx/splunk-otel-collector/internal/receiver/ptpreceiver [build failed]
2026-11-15T00:54:27.5505484Z FAIL	github.com/signalfx/splunk-otel-collector/internal/receiver/ptpreceiver/internal/metadata [build failed]
2026-11-15T00:54:27.5507842Z FAIL	github.com/signalfx/splunk-otel-collector/internal/settings [build failed]
2026-11-15T00:54:27.9375750Z === RUN   TestNewClientTLSVerification
2026-11-15T00:54:29.0321807Z === RUN   TestRunProcessSampling_AllProcessesGone/exited
2026-11-15T00:54:29.0373831Z === RUN   TestRunProcessSampling_AllProcessesGone/liveness-error
2026-11-15T00:54:29.0425089Z --- PASS: TestRunProcessSampling_AllProcessesGone (0.01s)
2026-11-15T00:54:29.0425415Z     --- PASS: TestRunProcessSampling_AllProcessesGone/exited (0.01s)
2026-11-15T00:54:29.0425686Z     --- PASS: TestRunProcessSampling_AllProcessesGone/liveness-error (0.01s)
2026-11-15T00:54:29.0425928Z === RUN   TestRunProcessSampling_RetriesFailedCPUBaseline
2026-11-15T00:54:29.3647858Z ?   	github.com/signalfx/splunk-otel-collector/internal/version	[no test files]
2026-11-15T00:54:29.3648066Z FAIL
2026-11-15T00:54:36.7505245Z ##[error]Process completed with exit code 1.
2026-11-15T00:54:36.7692813Z Post job cleanup.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43659653483/job/132645878595">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
</tbody>
</table>


<a id="pr-8440"></a>
### [Bump github.com/prometheus/prometheus from 0.315.1 to 0.316.0 (#8440)](https://github.com/signalfx/splunk-otel-collector/pull/8440)

Partial — the report is incomplete; findings use only available observations.
- Could not parse a failed workflow: could not parse workflow run 44425293885 attempt 1: workflow run 44425293885 attempt 1 job "test-macos-15-race_detection:false-coverage:false / unit-test" step "bundle.d" has no downloadable log

<table>
<thead>
<tr><th>Failure</th><th>Classification</th><th>Latest failure</th><th>Main history</th><th>PR commit history</th><th>Related CI history</th></tr>
</thead>
<tbody>
<tr><td>Integration Test / otelcol (amd64) / Run make binaries-linux_amd64 COVER_TESTING=true</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td>—</td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>Integration Test / otelcol (arm64) / Run make binaries-linux_arm64 COVER_TESTING=true</td><td>Broken</td><td><pre><code>2026-11-19T00:08:18.0070094Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:121:3: cannot use appV2 (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2) as &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable value in argument to api_v1.NewAPI: &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2 does not implement &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable (missing method ExemplarQuerier)
2026-11-19T00:08:18.0079045Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:122:3: cannot use o.ExemplarStorage (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable) as func(&#34;context&#34;.Context) v1.ScrapePoolsRetriever value in argument to api_v1.NewAPI
2026-11-19T00:08:18.0085110Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:123:3: cannot use factorySPr (variable of type func(_ &#34;context&#34;.Context) v1.ScrapePoolsRetriever) as func(&#34;context&#34;.Context) v1.TargetRetriever value in argument to api_v1.NewAPI
2026-11-19T00:08:18.0251898Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:124:3: cannot use factoryTr (variable of type func(_ &#34;context&#34;.Context) v1.TargetRetriever) as func(&#34;context&#34;.Context) v1.AlertmanagerRetriever value in argument to api_v1.NewAPI
2026-11-19T00:08:18.0258195Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:125:3: cannot use factoryAr (variable of type func(_ &#34;context&#34;.Context) v1.AlertmanagerRetriever) as func() &#34;github.com/prometheus/prometheus/config&#34;.Config value in argument to api_v1.NewAPI
2026-11-19T00:08:18.0264954Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:127:3: cannot use m.promCfg.Get (value of type func() &#34;github.com/prometheus/prometheus/config&#34;.Config) as map[string]string value in argument to api_v1.NewAPI
2026-11-19T00:08:18.0270165Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:128:3: cannot use o.Flags (variable of type map[string]string) as v1.GlobalURLOptions value in argument to api_v1.NewAPI
2026-11-19T00:08:18.0406700Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: cannot use api_v1.GlobalURLOptions{…} (value of struct type v1.GlobalURLOptions) as func(&#34;net/http&#34;.HandlerFunc) &#34;net/http&#34;.HandlerFunc value in argument to api_v1.NewAPI
2026-11-19T00:08:18.0412660Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:180:5: o.STZeroIngestionEnabled undefined (type *&#34;github.com/prometheus/prometheus/web&#34;.Options has no field or method STZeroIngestionEnabled)
2026-11-19T00:08:18.0416803Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: too many errors
2026-11-19T00:09:45.0621892Z make[1]: *** [Makefile:266: otelcol] Error 1
2026-11-19T00:09:45.0622500Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-19T00:09:45.0626147Z make: *** [Makefile:350: binaries-linux_arm64] Error 2
2026-11-19T00:09:45.0638230Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293835/job/135228459647">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / cross-compile (binaries-aix_ppc64) / compile / Build Collector</td><td>Broken</td><td><pre><code>2026-11-19T00:06:16.8800051Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:121:3: cannot use appV2 (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2) as &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable value in argument to api_v1.NewAPI: &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2 does not implement &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable (missing method ExemplarQuerier)
2026-11-19T00:06:16.8804072Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:122:3: cannot use o.ExemplarStorage (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable) as func(&#34;context&#34;.Context) v1.ScrapePoolsRetriever value in argument to api_v1.NewAPI
2026-11-19T00:06:16.8807014Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:123:3: cannot use factorySPr (variable of type func(_ &#34;context&#34;.Context) v1.ScrapePoolsRetriever) as func(&#34;context&#34;.Context) v1.TargetRetriever value in argument to api_v1.NewAPI
2026-11-19T00:06:16.8809984Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:124:3: cannot use factoryTr (variable of type func(_ &#34;context&#34;.Context) v1.TargetRetriever) as func(&#34;context&#34;.Context) v1.AlertmanagerRetriever value in argument to api_v1.NewAPI
2026-11-19T00:06:16.8812856Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:125:3: cannot use factoryAr (variable of type func(_ &#34;context&#34;.Context) v1.AlertmanagerRetriever) as func() &#34;github.com/prometheus/prometheus/config&#34;.Config value in argument to api_v1.NewAPI
2026-11-19T00:06:16.8815841Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:127:3: cannot use m.promCfg.Get (value of type func() &#34;github.com/prometheus/prometheus/config&#34;.Config) as map[string]string value in argument to api_v1.NewAPI
2026-11-19T00:06:16.8829363Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:128:3: cannot use o.Flags (variable of type map[string]string) as v1.GlobalURLOptions value in argument to api_v1.NewAPI
2026-11-19T00:06:16.8833463Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: cannot use api_v1.GlobalURLOptions{…} (value of struct type v1.GlobalURLOptions) as func(&#34;net/http&#34;.HandlerFunc) &#34;net/http&#34;.HandlerFunc value in argument to api_v1.NewAPI
2026-11-19T00:06:16.8836148Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:180:5: o.STZeroIngestionEnabled undefined (type *&#34;github.com/prometheus/prometheus/web&#34;.Options has no field or method STZeroIngestionEnabled)
2026-11-19T00:06:16.8838274Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: too many errors
2026-11-19T00:07:03.2672781Z make[1]: *** [Makefile:266: otelcol] Error 1
2026-11-19T00:07:03.2673274Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-19T00:07:03.2678169Z make: *** [Makefile:328: binaries-aix_ppc64] Error 2
2026-11-19T00:07:03.2751030Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293885/job/135228459455">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / cross-compile (binaries-linux_s390x) / compile / Build Collector</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td>—</td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / lint / Lint</td><td>Flaky</td><td><pre><code>2026-11-19T00:14:24.6875138Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:121:3: cannot use appV2 (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2) as &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable value in argument to api_v1.NewAPI: &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2 does not implement &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable (missing method ExemplarQuerier)
2026-11-19T00:14:24.6885317Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:122:3: cannot use o.ExemplarStorage (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable) as func(&#34;context&#34;.Context) v1.ScrapePoolsRetriever value in argument to api_v1.NewAPI
2026-11-19T00:14:24.7610602Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:123:3: cannot use factorySPr (variable of type func(_ &#34;context&#34;.Context) v1.ScrapePoolsRetriever) as func(&#34;context&#34;.Context) v1.TargetRetriever value in argument to api_v1.NewAPI
2026-11-19T00:14:24.7619073Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:124:3: cannot use factoryTr (variable of type func(_ &#34;context&#34;.Context) v1.TargetRetriever) as func(&#34;context&#34;.Context) v1.AlertmanagerRetriever value in argument to api_v1.NewAPI
2026-11-19T00:14:24.7626833Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:125:3: cannot use factoryAr (variable of type func(_ &#34;context&#34;.Context) v1.AlertmanagerRetriever) as func() &#34;github.com/prometheus/prometheus/config&#34;.Config value in argument to api_v1.NewAPI
2026-11-19T00:14:24.7642514Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:127:3: cannot use m.promCfg.Get (value of type func() &#34;github.com/prometheus/prometheus/config&#34;.Config) as map[string]string value in argument to api_v1.NewAPI
2026-11-19T00:14:24.7649461Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:128:3: cannot use o.Flags (variable of type map[string]string) as v1.GlobalURLOptions value in argument to api_v1.NewAPI
2026-11-19T00:14:24.7676962Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: cannot use api_v1.GlobalURLOptions{…} (value of struct type v1.GlobalURLOptions) as func(&#34;net/http&#34;.HandlerFunc) &#34;net/http&#34;.HandlerFunc value in argument to api_v1.NewAPI
2026-11-19T00:14:24.7687148Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:180:5: o.STZeroIngestionEnabled undefined (type *&#34;github.com/prometheus/prometheus/web&#34;.Options has no field or method STZeroIngestionEnabled)
2026-11-19T00:14:24.7696698Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: too many errors))) (typecheck)
2026-11-19T00:14:24.7701607Z 	&#34;github.com/signalfx/splunk-otel-collector/baseline&#34;
2026-11-19T00:14:24.7706630Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-19T00:14:24.7733461Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293885/job/135228458981">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / test-macos-15-race_detection:false-coverage:false / unit-test / bundle.d</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293885/job/135228459432">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / test-macos-26-race_detection:false-coverage:false / unit-test / bundle.d</td><td>Unknown<br>There is not enough matching history to classify this failure.<br>Includes inferred workflow passes.</td><td>—</td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / test-ubuntu-24.04-race_detection:false-coverage:false / unit-test / bundle.d</td><td>Broken<br>Includes inferred workflow passes.</td><td><pre><code>2026-11-19T00:07:25.1912675Z ##[error]/home/runner/go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:121:3: cannot use appV2 (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2) as &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable value in argument to api_v1.NewAPI: &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2 does not implement &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable (missing method ExemplarQuerier)
2026-11-19T00:07:25.1951978Z ##[error]/home/runner/go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:122:3: cannot use o.ExemplarStorage (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable) as func(&#34;context&#34;.Context) v1.ScrapePoolsRetriever value in argument to api_v1.NewAPI
2026-11-19T00:07:25.1958925Z ##[error]/home/runner/go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:123:3: cannot use factorySPr (variable of type func(_ &#34;context&#34;.Context) v1.ScrapePoolsRetriever) as func(&#34;context&#34;.Context) v1.TargetRetriever value in argument to api_v1.NewAPI
2026-11-19T00:07:25.1965426Z ##[error]/home/runner/go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:124:3: cannot use factoryTr (variable of type func(_ &#34;context&#34;.Context) v1.TargetRetriever) as func(&#34;context&#34;.Context) v1.AlertmanagerRetriever value in argument to api_v1.NewAPI
2026-11-19T00:07:25.1972303Z ##[error]/home/runner/go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:125:3: cannot use factoryAr (variable of type func(_ &#34;context&#34;.Context) v1.AlertmanagerRetriever) as func() &#34;github.com/prometheus/prometheus/config&#34;.Config value in argument to api_v1.NewAPI
2026-11-19T00:07:25.1978635Z ##[error]/home/runner/go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:127:3: cannot use m.promCfg.Get (value of type func() &#34;github.com/prometheus/prometheus/config&#34;.Config) as map[string]string value in argument to api_v1.NewAPI
2026-11-19T00:07:25.2015123Z ##[error]/home/runner/go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:128:3: cannot use o.Flags (variable of type map[string]string) as v1.GlobalURLOptions value in argument to api_v1.NewAPI
2026-11-19T00:07:25.2022131Z ##[error]/home/runner/go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: cannot use api_v1.GlobalURLOptions{…} (value of struct type v1.GlobalURLOptions) as func(&#34;net/http&#34;.HandlerFunc) &#34;net/http&#34;.HandlerFunc value in argument to api_v1.NewAPI
2026-11-19T00:07:25.2028122Z ##[error]/home/runner/go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:180:5: o.STZeroIngestionEnabled undefined (type *&#34;github.com/prometheus/prometheus/web&#34;.Options has no field or method STZeroIngestionEnabled)
2026-11-19T00:07:25.2032718Z ##[error]/home/runner/go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: too many errors
2026-11-19T00:08:49.3697491Z make: *** [Makefile.Common:195: bundle.d] Error 1
2026-11-19T00:08:49.3708002Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293885/job/135228459459">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>otelcol-fips / otelcol-fips (linux, amd64) / Run make otelcol-fips</td><td>Broken</td><td><pre><code>2026-11-19T00:07:28.8892371Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:121:3: cannot use appV2 (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2) as &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable value in argument to api_v1.NewAPI: &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2 does not implement &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable (missing method ExemplarQuerier)
2026-11-19T00:07:28.8922716Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:122:3: cannot use o.ExemplarStorage (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable) as func(&#34;context&#34;.Context) v1.ScrapePoolsRetriever value in argument to api_v1.NewAPI
2026-11-19T00:07:28.8928625Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:123:3: cannot use factorySPr (variable of type func(_ &#34;context&#34;.Context) v1.ScrapePoolsRetriever) as func(&#34;context&#34;.Context) v1.TargetRetriever value in argument to api_v1.NewAPI
2026-11-19T00:07:28.8934228Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:124:3: cannot use factoryTr (variable of type func(_ &#34;context&#34;.Context) v1.TargetRetriever) as func(&#34;context&#34;.Context) v1.AlertmanagerRetriever value in argument to api_v1.NewAPI
2026-11-19T00:07:28.8939606Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:125:3: cannot use factoryAr (variable of type func(_ &#34;context&#34;.Context) v1.AlertmanagerRetriever) as func() &#34;github.com/prometheus/prometheus/config&#34;.Config value in argument to api_v1.NewAPI
2026-11-19T00:07:28.8944573Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:127:3: cannot use m.promCfg.Get (value of type func() &#34;github.com/prometheus/prometheus/config&#34;.Config) as map[string]string value in argument to api_v1.NewAPI
2026-11-19T00:07:28.9001956Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:128:3: cannot use o.Flags (variable of type map[string]string) as v1.GlobalURLOptions value in argument to api_v1.NewAPI
2026-11-19T00:07:28.9007814Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: cannot use api_v1.GlobalURLOptions{…} (value of struct type v1.GlobalURLOptions) as func(&#34;net/http&#34;.HandlerFunc) &#34;net/http&#34;.HandlerFunc value in argument to api_v1.NewAPI
2026-11-19T00:07:28.9030309Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:180:5: o.STZeroIngestionEnabled undefined (type *&#34;github.com/prometheus/prometheus/web&#34;.Options has no field or method STZeroIngestionEnabled)
2026-11-19T00:07:28.9035022Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: too many errors
2026-11-19T00:08:50.1998109Z make: *** [Makefile:411: otelcol-fips] Error 1
2026-11-19T00:08:50.2018256Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293869/job/135228458583">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>otelcol-fips / otelcol-fips (linux, arm64) / Run make otelcol-fips</td><td>Broken</td><td><pre><code>2026-11-19T00:03:54.4393129Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:121:3: cannot use appV2 (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2) as &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable value in argument to api_v1.NewAPI: &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2 does not implement &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable (missing method ExemplarQuerier)
2026-11-19T00:03:54.4398996Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:122:3: cannot use o.ExemplarStorage (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable) as func(&#34;context&#34;.Context) v1.ScrapePoolsRetriever value in argument to api_v1.NewAPI
2026-11-19T00:03:54.4403035Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:123:3: cannot use factorySPr (variable of type func(_ &#34;context&#34;.Context) v1.ScrapePoolsRetriever) as func(&#34;context&#34;.Context) v1.TargetRetriever value in argument to api_v1.NewAPI
2026-11-19T00:03:54.4406945Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:124:3: cannot use factoryTr (variable of type func(_ &#34;context&#34;.Context) v1.TargetRetriever) as func(&#34;context&#34;.Context) v1.AlertmanagerRetriever value in argument to api_v1.NewAPI
2026-11-19T00:03:54.4410871Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:125:3: cannot use factoryAr (variable of type func(_ &#34;context&#34;.Context) v1.AlertmanagerRetriever) as func() &#34;github.com/prometheus/prometheus/config&#34;.Config value in argument to api_v1.NewAPI
2026-11-19T00:03:54.4414471Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:127:3: cannot use m.promCfg.Get (value of type func() &#34;github.com/prometheus/prometheus/config&#34;.Config) as map[string]string value in argument to api_v1.NewAPI
2026-11-19T00:03:54.4417554Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:128:3: cannot use o.Flags (variable of type map[string]string) as v1.GlobalURLOptions value in argument to api_v1.NewAPI
2026-11-19T00:03:54.4513414Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: cannot use api_v1.GlobalURLOptions{…} (value of struct type v1.GlobalURLOptions) as func(&#34;net/http&#34;.HandlerFunc) &#34;net/http&#34;.HandlerFunc value in argument to api_v1.NewAPI
2026-11-19T00:03:54.4517015Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:180:5: o.STZeroIngestionEnabled undefined (type *&#34;github.com/prometheus/prometheus/web&#34;.Options has no field or method STZeroIngestionEnabled)
2026-11-19T00:03:54.4519975Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: too many errors
2026-11-19T00:04:57.2512172Z make: *** [Makefile:411: otelcol-fips] Error 1
2026-11-19T00:04:57.2522823Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293869/job/135228458532">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>otelcol-fips / otelcol-fips (windows, amd64) / Run make otelcol-fips</td><td>Broken</td><td><pre><code>2026-11-19T00:07:11.2755807Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:121:3: cannot use appV2 (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2) as &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable value in argument to api_v1.NewAPI: &#34;github.com/prometheus/prometheus/storage&#34;.AppendableV2 does not implement &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable (missing method ExemplarQuerier)
2026-11-19T00:07:11.2762365Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:122:3: cannot use o.ExemplarStorage (variable of interface type &#34;github.com/prometheus/prometheus/storage&#34;.ExemplarQueryable) as func(&#34;context&#34;.Context) v1.ScrapePoolsRetriever value in argument to api_v1.NewAPI
2026-11-19T00:07:11.2767216Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:123:3: cannot use factorySPr (variable of type func(_ &#34;context&#34;.Context) v1.ScrapePoolsRetriever) as func(&#34;context&#34;.Context) v1.TargetRetriever value in argument to api_v1.NewAPI
2026-11-19T00:07:11.2772576Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:124:3: cannot use factoryTr (variable of type func(_ &#34;context&#34;.Context) v1.TargetRetriever) as func(&#34;context&#34;.Context) v1.AlertmanagerRetriever value in argument to api_v1.NewAPI
2026-11-19T00:07:11.2777384Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:125:3: cannot use factoryAr (variable of type func(_ &#34;context&#34;.Context) v1.AlertmanagerRetriever) as func() &#34;github.com/prometheus/prometheus/config&#34;.Config value in argument to api_v1.NewAPI
2026-11-19T00:07:11.2843645Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:127:3: cannot use m.promCfg.Get (value of type func() &#34;github.com/prometheus/prometheus/config&#34;.Config) as map[string]string value in argument to api_v1.NewAPI
2026-11-19T00:07:11.2847789Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:128:3: cannot use o.Flags (variable of type map[string]string) as v1.GlobalURLOptions value in argument to api_v1.NewAPI
2026-11-19T00:07:11.2853104Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: cannot use api_v1.GlobalURLOptions{…} (value of struct type v1.GlobalURLOptions) as func(&#34;net/http&#34;.HandlerFunc) &#34;net/http&#34;.HandlerFunc value in argument to api_v1.NewAPI
2026-11-19T00:07:11.2857624Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:180:5: o.STZeroIngestionEnabled undefined (type *&#34;github.com/prometheus/prometheus/web&#34;.Options has no field or method STZeroIngestionEnabled)
2026-11-19T00:07:11.2861006Z ##[error]../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver@v0.165.1/internal/apiserver/manager.go:129:3: too many errors
2026-11-19T00:08:27.7585119Z make: *** [Makefile:417: otelcol-fips] Error 1
2026-11-19T00:08:27.7594437Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293869/job/135228458165">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/cmd/discoverybundler</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293821/job/135228458866">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/cmd/otelcol</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293821/job/135228458866">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/cmd/otelcollauncher/health</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293821/job/135228458866">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/components</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293821/job/135228458866">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/confmapprovider/discovery</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293821/job/135228458866">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / package github.com/signalfx/splunk-otel-collector/internal/settings</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293821/job/135228458866">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / windows-test (windows-2022) / Unit tests with coverage</td><td>Broken</td><td><pre><code>2026-11-19T00:16:37.6464105Z         --- PASS: TestConsumeMetrics/regexp/partial (0.00s)
2026-11-19T00:16:37.6465113Z         --- PASS: TestConsumeMetrics/regexp/failed (0.00s)
2026-11-19T00:16:37.6465856Z     --- PASS: TestConsumeMetrics/expr (0.00s)
2026-11-19T00:16:37.6467285Z         --- PASS: TestConsumeMetrics/expr/partial (0.00s)
2026-11-19T00:16:37.6467710Z         --- PASS: TestConsumeMetrics/expr/failed (0.00s)
2026-11-19T00:16:37.6468030Z === RUN   TestNewDiscoveryReceiver
2026-11-19T00:16:37.6492264Z === RUN   TestStatementEvaluation/strict/partial
2026-11-19T00:16:37.6500830Z === RUN   TestStatementEvaluation/strict/failed
2026-11-19T00:16:37.6525982Z === RUN   TestStatementEvaluation/regexp
2026-11-19T00:16:37.6551856Z === RUN   TestStatementEvaluation/regexp/partial
2026-11-19T00:16:37.6555792Z === RUN   TestStatementEvaluation/regexp/failed
2026-11-19T00:16:37.6569800Z === RUN   TestStatementEvaluation/expr
2026-11-19T00:16:37.6576348Z === RUN   TestStatementEvaluation/expr/partial
2026-11-19T00:16:37.6588897Z === RUN   TestStatementEvaluation/expr/failed
2026-11-19T00:16:37.6594529Z --- PASS: TestStatementEvaluation (0.01s)
2026-11-19T00:16:37.6596812Z         --- PASS: TestStatementEvaluation/strict/partial (0.00s)
2026-11-19T00:16:37.6597997Z         --- PASS: TestStatementEvaluation/strict/failed (0.00s)
2026-11-19T00:16:37.6598740Z     --- PASS: TestStatementEvaluation/regexp (0.00s)
2026-11-19T00:16:37.6600000Z         --- PASS: TestStatementEvaluation/regexp/partial (0.00s)
2026-11-19T00:16:37.6600562Z         --- PASS: TestStatementEvaluation/regexp/failed (0.00s)
2026-11-19T00:16:37.6601130Z     --- PASS: TestStatementEvaluation/expr (0.00s)
2026-11-19T00:16:37.6602224Z         --- PASS: TestStatementEvaluation/expr/partial (0.00s)
2026-11-19T00:16:37.6602802Z         --- PASS: TestStatementEvaluation/expr/failed (0.00s)
2026-11-19T00:16:37.6603484Z === RUN   TestStatementEvaluatorSampledLogger
2026-11-19T00:16:50.8783562Z === RUN   TestScrapeMetricsStringAndScalarResultTypesAreLoggedNotErrored/string
2026-11-19T00:16:50.8790380Z     logger.go:146: 2026-11-19T00:16:50.878Z	ERROR	PromQL response of type string is not supported:	{&#34;query&#34;: &#34;up&#34;}
2026-11-19T00:16:50.8791623Z === RUN   TestScrapeMetricsStringAndScalarResultTypesAreLoggedNotErrored/scalar
2026-11-19T00:16:50.8799661Z     logger.go:146: 2026-11-19T00:16:50.879Z	ERROR	PromQL response of type scalar is not supported:	{&#34;query&#34;: &#34;up&#34;}
2026-11-19T00:16:50.8801497Z --- PASS: TestScrapeMetricsStringAndScalarResultTypesAreLoggedNotErrored (0.00s)
2026-11-19T00:16:50.8805788Z === RUN   TestScrapeMetricsUnsupportedResultTypeIsLoggedNotErrored
2026-11-19T00:16:50.8808352Z     logger.go:146: 2026-11-19T00:16:50.880Z	ERROR	PromQL response unsupported:	{&#34;query&#34;: &#34;up&#34;, &#34;resultType&#34;: &#34;unknown&#34;}
2026-11-19T00:16:50.8810632Z --- PASS: TestScrapeMetricsUnsupportedResultTypeIsLoggedNotErrored (0.00s)
2026-11-19T00:16:53.0465099Z ok  	github.com/signalfx/splunk-otel-collector/internal/receiver/ptpreceiver/internal/metadata	0.046s	coverage: 85.2% of statements
2026-11-19T00:16:53.0525516Z FAIL	github.com/signalfx/splunk-otel-collector/internal/settings [build failed]
2026-11-19T00:16:53.5313365Z === RUN   TestNewClientTLSVerification
2026-11-19T00:16:54.9631317Z === RUN   TestRunProcessSampling_AllProcessesGone/exited
2026-11-19T00:16:54.9682903Z === RUN   TestRunProcessSampling_AllProcessesGone/liveness-error
2026-11-19T00:16:54.9739985Z --- PASS: TestRunProcessSampling_AllProcessesGone (0.01s)
2026-11-19T00:16:54.9740936Z     --- PASS: TestRunProcessSampling_AllProcessesGone/exited (0.01s)
2026-11-19T00:16:54.9741681Z     --- PASS: TestRunProcessSampling_AllProcessesGone/liveness-error (0.01s)
2026-11-19T00:16:54.9742313Z === RUN   TestRunProcessSampling_RetriesFailedCPUBaseline
2026-11-19T00:16:55.4303636Z ?   	github.com/signalfx/splunk-otel-collector/internal/version	[no test files]
2026-11-19T00:16:55.4304121Z FAIL
2026-11-19T00:17:03.5383700Z ##[error]Process completed with exit code 1.
2026-11-19T00:17:03.5701972Z Post job cleanup.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44425293821/job/135228458866">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / windows-test (windows-2025) / Unit tests with coverage</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td>—</td><td>—</td><td>—</td><td>—</td></tr>
</tbody>
</table>


<a id="pr-8397"></a>
### [Bump github.com/hashicorp/vault-plugin-auth-gcp from 0.23.2-0.20260604163449-108858b5ffea to 0.24.2 (#8397)](https://github.com/signalfx/splunk-otel-collector/pull/8397)

Partial — the report is incomplete; findings use only available observations.

<table>
<thead>
<tr><th>Failure</th><th>Classification</th><th>Latest failure</th><th>Main history</th><th>PR commit history</th><th>Related CI history</th></tr>
</thead>
<tbody>
<tr><td>Integration Test / otelcol (arm64) / Run make binaries-linux_arm64 COVER_TESTING=true</td><td>Broken</td><td><pre><code>2026-11-17T16:03:40.2291690Z go: go.mod requires go &gt;= 1.27.1 (running go 1.26.11; GOTOOLCHAIN=local)
2026-11-17T16:03:40.2297907Z make[1]: *** [Makefile:265: otelcol] Error 1
2026-11-17T16:03:40.2298458Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-17T16:03:40.2301771Z make: *** [Makefile:350: binaries-linux_arm64] Error 2
2026-11-17T16:03:40.2316782Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44120955124/job/134171831100">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / cross-compile (binaries-aix_ppc64) / compile / Build Collector</td><td>Broken</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44120955129/job/134171832249">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / generate-metrics / generate-metrics</td><td>Broken</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44120955129/job/134171830485">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / lint / Lint</td><td>Flaky</td><td><pre><code>2026-11-17T16:04:46.6370248Z Check License finished successfully
2026-11-17T16:04:46.6376558Z running misspell -error
2026-11-17T16:04:46.6629778Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:13:21.8335616Z Check License finished successfully
2026-11-17T16:13:21.8344491Z running misspell -error
2026-11-17T16:13:21.8759295Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:13:34.3041527Z Check License finished successfully
2026-11-17T16:13:34.3045271Z running misspell -error
2026-11-17T16:13:34.3279539Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:13:44.9299801Z Check License finished successfully
2026-11-17T16:13:44.9305045Z running misspell -error
2026-11-17T16:13:44.9546148Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:13:54.8620803Z Check License finished successfully
2026-11-17T16:13:54.8625663Z running misspell -error
2026-11-17T16:13:54.8856843Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:13:56.3581594Z Check License finished successfully
2026-11-17T16:13:56.3589971Z running misspell -error
2026-11-17T16:13:56.3825154Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:13:59.3625345Z Check License finished successfully
2026-11-17T16:13:59.3633121Z running misspell -error
2026-11-17T16:13:59.3871666Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:14:00.8275033Z Check License finished successfully
2026-11-17T16:14:00.8279475Z running misspell -error
2026-11-17T16:14:00.8500188Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:14:14.4673315Z Check License finished successfully
2026-11-17T16:14:14.4679906Z running misspell -error
2026-11-17T16:14:14.4926496Z Skip linting for smartagentextension
2026-11-17T16:14:14.7860814Z Check License finished successfully
2026-11-17T16:14:14.7868621Z running misspell -error
2026-11-17T16:14:14.8174403Z Skip linting for smartagentreceiver
2026-11-17T16:14:14.9169983Z Check License finished successfully
2026-11-17T16:14:14.9173350Z running misspell -error
2026-11-17T16:14:14.9472380Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:14:22.8949150Z Check License finished successfully
2026-11-17T16:14:22.8956243Z running misspell -error
2026-11-17T16:14:22.9273050Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:15:15.2560716Z Check License finished successfully
2026-11-17T16:15:15.2572676Z running misspell -error
2026-11-17T16:15:15.3288365Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:16:03.8570541Z Check License finished successfully
2026-11-17T16:16:03.8577269Z running misspell -error
2026-11-17T16:16:03.8907937Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:16:18.1149971Z Check License finished successfully
2026-11-17T16:16:18.1161386Z running misspell -error
2026-11-17T16:16:18.1403108Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:16:19.6357224Z Check License finished successfully
2026-11-17T16:16:19.6365145Z running misspell -error
2026-11-17T16:16:19.9836231Z golangci-lint run --allow-parallel-runners -j4
2026-11-17T16:16:20.0694106Z Error: can&#39;t load config: the Go language version (go1.26) used to build golangci-lint is lower than the targeted Go version (1.27.1)
2026-11-17T16:16:20.0705242Z The command is terminated due to an error: can&#39;t load config: the Go language version (go1.26) used to build golangci-lint is lower than the targeted Go version (1.27.1)
2026-11-17T16:16:20.0717729Z make[2]: *** [Makefile.Common:155: lint] Error 3
2026-11-17T16:16:20.0720848Z make[1]: *** [Makefile:82: /home/runner/work/splunk-otel-collector/splunk-otel-collector] Error 2
2026-11-17T16:16:20.0724443Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-17T16:16:20.0778482Z make: *** [Makefile:235: lint-all] Error 2
2026-11-17T16:16:20.1341827Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44120955129/job/134171830555">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / test-macos-26-race_detection:false-coverage:false / unit-test / bundle.d</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44120955129/job/134171833113">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>build-and-test / tidy / tidy</td><td>Broken</td><td><pre><code>2026-11-17T16:09:36.4758178Z go: go.mod requires go &gt;= 1.27.1 (running go 1.26.11; GOTOOLCHAIN=local)
2026-11-17T16:09:36.4765128Z make[2]: *** [Makefile.Common:160: tidy] Error 1
2026-11-17T16:09:36.4767750Z make[1]: *** [Makefile:82: /home/runner/work/splunk-otel-collector/splunk-otel-collector] Error 2
2026-11-17T16:09:36.4768956Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector&#39;
2026-11-17T16:09:36.4771800Z make: *** [Makefile:225: tidy-all] Error 2
2026-11-17T16:09:36.4788051Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44120955129/job/134171830322">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>otelcol-fips / otelcol-fips (linux, amd64) / Run make otelcol-fips</td><td>Broken</td><td><pre><code>2026-11-17T16:02:23.5921848Z go: go.mod requires go &gt;= 1.27.1 (running go 1.26.11; GOTOOLCHAIN=local)
2026-11-17T16:02:23.5926764Z make: *** [Makefile:410: otelcol-fips] Error 1
2026-11-17T16:02:23.5937513Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44120955110/job/134171830889">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>otelcol-fips / otelcol-fips (linux, arm64) / Run make otelcol-fips</td><td>Broken</td><td><pre><code>2026-11-17T16:01:43.8624874Z go: go.mod requires go &gt;= 1.27.1 (running go 1.26.11; GOTOOLCHAIN=local)
2026-11-17T16:01:43.8629307Z make: *** [Makefile:410: otelcol-fips] Error 1
2026-11-17T16:01:43.8646652Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44120955110/job/134171830875">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>otelcol-fips / otelcol-fips (windows, amd64) / Run make otelcol-fips</td><td>Broken</td><td><pre><code>2026-11-17T16:02:28.2842687Z go: go.mod requires go &gt;= 1.27.1 (running go 1.26.11; GOTOOLCHAIN=local)
2026-11-17T16:02:28.2847260Z make: *** [Makefile:416: otelcol-fips] Error 1
2026-11-17T16:02:28.2861715Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44120955110/job/134171830210">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>otlpintput / ta-v2-build-packages / Package Splunk Connect for OTLP</td><td>Broken</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44120955047/job/134171826582">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>parity-test / parity-test / Build otelcol</td><td>Broken</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44120955512/job/134171831969">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>windows-test / windows-test (windows-2022) / Unit tests with coverage</td><td>Broken</td><td><pre><code>2026-11-17T16:01:32.0660149Z go: go.mod requires go &gt;= 1.27.1 (running go 1.26.11; GOTOOLCHAIN=local)
2026-11-17T16:01:32.1390787Z ##[error]Process completed with exit code 1.
2026-11-17T16:01:32.1555079Z Post job cleanup.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44120955164/job/134171830664">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
</tbody>
</table>


<a id="pr-8076"></a>
### [Bump ruby from 3.3.8-bookworm to 4.0.0-bookworm in /deployments/puppet (#8076)](https://github.com/signalfx/splunk-otel-collector/pull/8076)

Partial — the report is incomplete; findings use only available observations.

<table>
<thead>
<tr><th>Failure</th><th>Classification</th><th>Latest failure</th><th>Main history</th><th>PR commit history</th><th>Related CI history</th></tr>
</thead>
<tbody>
<tr><td>puppet-test / puppet-lint / Lint</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>2026-11-13T18:03:39.1544758Z 	splunk-otel-connector-puppet-dev \
2026-11-13T18:03:39.1545207Z 	puppet-lint --fail-on-warnings \
2026-11-13T18:03:39.1545748Z 	--no-parameter_documentation-check --no-parameter_types-check \
2026-11-13T18:03:39.1546826Z 	.
2026-11-13T18:03:39.5386438Z Whoops! It looks like puppet-lint has encountered an error that it doesn&#39;t
2026-11-13T18:03:39.5387455Z know how to handle. Please open an issue at https://github.com/puppetlabs/puppet-lint
2026-11-13T18:03:39.5404899Z ```
2026-11-13T18:03:39.5405172Z error:
2026-11-13T18:03:39.5405458Z ```
2026-11-13T18:03:39.5418537Z ```
2026-11-13T18:03:39.6241176Z make: *** [Makefile:18: lint] Error 1
2026-11-13T18:03:39.6242327Z make: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector/deployments/puppet&#39;
2026-11-13T18:03:39.6263144Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43478614779/job/132068953553">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>puppet-test / puppet-rake-spec / Lint</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>2026-11-13T18:03:41.2387402Z
2026-11-13T18:03:41.2387853Z An error occurred while loading ./spec/classes/init_spec.rb.
2026-11-13T18:03:41.2388265Z Failure/Error: require &#39;puppetlabs_spec_helper/module_spec_helper&#39;
2026-11-13T18:03:41.2388687Z
2026-11-13T18:03:41.2397820Z Finished in 0.00004 seconds (files took 0.18007 seconds to load)
2026-11-13T18:03:41.2398140Z 0 examples, 0 failures, 1 error occurred outside of examples
2026-11-13T18:03:41.2398335Z
2026-11-13T18:03:41.2439257Z /usr/local/bin/ruby -I/usr/local/bundle/gems/rspec-support-3.13.1/lib:/usr/local/bundle/gems/rspec-core-3.13.2/lib /usr/local/bundle/gems/rspec-core-3.13.2/exe/rspec --pattern spec/\{aliases,classes,defines,functions,hosts,integration,plans,tasks,type_aliases,types,unit\}/\*\*/\*_spec.rb failed
2026-11-13T18:03:41.4245912Z make: *** [Makefile:28: rake-spec] Error 1
2026-11-13T18:03:41.4246651Z make: Leaving directory &#39;/home/runner/work/splunk-otel-collector/splunk-otel-collector/deployments/puppet&#39;
2026-11-13T18:03:41.4263119Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43478614779/job/132068953658">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
</tbody>
</table>


<a id="pr-7879"></a>
### [Bump actionshub/chef-install from 5.0.2 to 6.0.0 (#7879)](https://github.com/signalfx/splunk-otel-collector/pull/7879)

Partial — the report is incomplete; findings use only available observations.

<table>
<thead>
<tr><th>Failure</th><th>Classification</th><th>Latest failure</th><th>Main history</th><th>PR commit history</th><th>Related CI history</th></tr>
</thead>
<tbody>
<tr><td>chef-test / chef-kitchen-windows (windows-2022, custom-vars) / Install chef</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>2026-11-13T18:19:12.6033675Z [31;1m[0m[36;1m[36;1m[0m[36;1m[0m[36;1m[31;1m[31;1m[36;1m     | [31;1mInstaller download redirected to an unexpected URL: [0m
2026-11-13T18:19:12.6469467Z ##[error]Process completed with exit code 1.
2026-11-13T18:19:12.6477033Z ##[end-action id=__actionshub_chef-install.__run_2;outcome=failure;conclusion=failure;duration_ms=1221]
2026-11-13T18:19:12.6611955Z ##[group]Run actions/upload-artifact@cf430e030ddbb5b0abf93d22962f4752f3646cd9</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43478662580/job/132078472445">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>chef-test / chef-kitchen-windows (windows-2022, default) / Install chef</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>2026-11-13T18:19:19.1000222Z [31;1m[0m[36;1m[36;1m[0m[36;1m[0m[36;1m[31;1m[31;1m[36;1m     | [31;1mInstaller download redirected to an unexpected URL: [0m
2026-11-13T18:19:19.1521525Z ##[error]Process completed with exit code 1.
2026-11-13T18:19:19.1530034Z ##[end-action id=__actionshub_chef-install.__run_2;outcome=failure;conclusion=failure;duration_ms=1288]
2026-11-13T18:19:19.1700517Z ##[group]Run actions/upload-artifact@cf430e030ddbb5b0abf93d22962f4752f3646cd9</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43478662580/job/132078472258">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>chef-test / chef-kitchen-windows (windows-2022, platform-logs) / Install chef</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>2026-11-13T18:19:34.2892189Z [31;1m[0m[36;1m[36;1m[0m[36;1m[0m[36;1m[31;1m[31;1m[36;1m     | [31;1mInstaller download redirected to an unexpected URL: [0m
2026-11-13T18:19:34.3426025Z ##[error]Process completed with exit code 1.
2026-11-13T18:19:34.3434983Z ##[end-action id=__actionshub_chef-install.__run_2;outcome=failure;conclusion=failure;duration_ms=1908]
2026-11-13T18:19:34.3612292Z ##[group]Run actions/upload-artifact@cf430e030ddbb5b0abf93d22962f4752f3646cd9</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43478662580/job/132078472184">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>chef-test / chef-kitchen-windows (windows-2022, platform-logs-with-observability) / Install chef</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>2026-11-13T18:19:37.6334922Z [31;1m[0m[36;1m[36;1m[0m[36;1m[0m[36;1m[31;1m[31;1m[36;1m     | [31;1mInstaller download redirected to an unexpected URL: [0m
2026-11-13T18:19:37.6813198Z ##[error]Process completed with exit code 1.
2026-11-13T18:19:37.6821049Z ##[end-action id=__actionshub_chef-install.__run_2;outcome=failure;conclusion=failure;duration_ms=946]
2026-11-13T18:19:37.6967584Z ##[group]Run actions/upload-artifact@cf430e030ddbb5b0abf93d22962f4752f3646cd9</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43478662580/job/132078472235">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>chef-test / chef-kitchen-windows (windows-2025, custom-vars) / Install chef</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>2026-11-13T18:19:10.0256032Z [36;1m  if ($statusCode) {[0m
2026-11-13T18:19:10.0256349Z [36;1m    $message = &#34;Failed to download installer from $scriptUrl (HTTP $statusCode).&#34;[0m
2026-11-13T18:19:10.0256678Z [36;1m  } else {[0m
2026-11-13T18:19:10.0256992Z [36;1m    $message = &#34;Failed to download installer from $scriptUrl.&#34;[0m
2026-11-13T18:19:10.0257274Z [36;1m  }[0m
2026-11-13T18:19:10.0259485Z [36;1mif ($response.StatusCode -ne 200) {[0m
2026-11-13T18:19:10.0271537Z [36;1m  $message = &#34;Failed to download installer from $scriptUrl (HTTP $($response.StatusCode)).&#34;[0m
2026-11-13T18:19:10.0271867Z [36;1m  if ($env:LICENSE) {[0m
2026-11-13T18:19:10.9407213Z [31;1m[0m[36;1m[36;1m[0m[36;1m[0m[36;1m[31;1m[31;1m[36;1m     | [31;1mInstaller download redirected to an unexpected URL: [0m
2026-11-13T18:19:10.9901628Z ##[error]Process completed with exit code 1.
2026-11-13T18:19:10.9908729Z ##[end-action id=__actionshub_chef-install.__run_2;outcome=failure;conclusion=failure;duration_ms=975]</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43478662580/job/132078472197">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>chef-test / chef-kitchen-windows (windows-2025, default) / Install chef</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>2026-11-13T18:19:21.0053346Z [31;1m[0m[36;1m[36;1m[0m[36;1m[0m[36;1m[31;1m[31;1m[36;1m     | [31;1mInstaller download redirected to an unexpected URL: [0m
2026-11-13T18:19:21.0594670Z ##[error]Process completed with exit code 1.
2026-11-13T18:19:21.0602443Z ##[end-action id=__actionshub_chef-install.__run_2;outcome=failure;conclusion=failure;duration_ms=1426]
2026-11-13T18:19:21.0763509Z ##[group]Run actions/upload-artifact@cf430e030ddbb5b0abf93d22962f4752f3646cd9</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43478662580/job/132078472201">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>chef-test / chef-kitchen-windows (windows-2025, platform-logs) / Install chef</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>2026-11-13T18:19:25.0108647Z [36;1m  if ($statusCode) {[0m
2026-11-13T18:19:25.0109062Z [36;1m    $message = &#34;Failed to download installer from $scriptUrl (HTTP $statusCode).&#34;[0m
2026-11-13T18:19:25.0109494Z [36;1m  } else {[0m
2026-11-13T18:19:25.0109830Z [36;1m    $message = &#34;Failed to download installer from $scriptUrl.&#34;[0m
2026-11-13T18:19:25.0110332Z [36;1m  }[0m
2026-11-13T18:19:25.0114257Z [36;1mif ($response.StatusCode -ne 200) {[0m
2026-11-13T18:19:25.0128162Z [36;1m  $message = &#34;Failed to download installer from $scriptUrl (HTTP $($response.StatusCode)).&#34;[0m
2026-11-13T18:19:25.0128773Z [36;1m  if ($env:LICENSE) {[0m
2026-11-13T18:19:26.2473050Z [31;1m[0m[36;1m[36;1m[0m[36;1m[0m[36;1m[31;1m[31;1m[36;1m     | [31;1mInstaller download redirected to an unexpected URL: [0m
2026-11-13T18:19:26.3175102Z ##[error]Process completed with exit code 1.
2026-11-13T18:19:26.3185550Z ##[end-action id=__actionshub_chef-install.__run_2;outcome=failure;conclusion=failure;duration_ms=1324]
2026-11-13T18:19:26.3409039Z ##[group]Run actions/upload-artifact@cf430e030ddbb5b0abf93d22962f4752f3646cd9</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43478662580/job/132078472310">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>chef-test / chef-kitchen-windows (windows-2025, platform-logs-with-observability) / Install chef</td><td>Unknown<br>There is not enough matching history to classify this failure.</td><td><pre><code>2026-11-13T18:19:19.5535252Z [31;1m[0m[36;1m[36;1m[0m[36;1m[0m[36;1m[31;1m[31;1m[36;1m     | [31;1mInstaller download redirected to an unexpected URL: [0m
2026-11-13T18:19:19.6046896Z ##[error]Process completed with exit code 1.
2026-11-13T18:19:19.6055782Z ##[end-action id=__actionshub_chef-install.__run_2;outcome=failure;conclusion=failure;duration_ms=1092]
2026-11-13T18:19:19.6239201Z ##[group]Run actions/upload-artifact@cf430e030ddbb5b0abf93d22962f4752f3646cd9</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/43478662580/job/132078472309">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
</tbody>
</table>


