# Dependency Bump PR Status

## Pending

<table>
<thead>
<tr><th>Author</th><th>Pull Request</th><th>CI/CD</th></tr>
</thead>
<tbody>
<tr><td>github-actions</td><td><a href="https://github.com/signalfx/splunk-otel-collector-chart/pull/2731">Bump java from v2.33.1 to v2.34.0 in helm-charts/splunk-otel-collector/values.yaml (#2731)</a></td><td><a href="#pr-2731">FAILURE</a> — 35 passed, 1 failed, 1 skipped</td></tr>
<tr><td>github-actions</td><td><a href="https://github.com/signalfx/splunk-otel-collector-chart/pull/2730">Bump java-csa from v2.33.1 to v2.34.0 in helm-charts/splunk-otel-collector/values.yaml (#2730)</a></td><td>PENDING — 36 passed, 1 skipped</td></tr>
<tr><td>github-actions</td><td><a href="https://github.com/signalfx/splunk-otel-collector-chart/pull/2728">Bump obi from 0.16.1 to 0.16.4 in helm-charts/splunk-otel-collector/Chart.yaml (#2728)</a></td><td><a href="#pr-2728">FAILURE</a> — 104 passed, 1 failed, 1 skipped</td></tr>
<tr><td>github-actions</td><td><a href="https://github.com/signalfx/splunk-otel-collector-chart/pull/2727">Bump TargetAllocator from 0.162.0 to 0.163.0 in helm-charts/splunk-otel-collector/Chart.yaml (#2727)</a></td><td><a href="#pr-2727">FAILURE</a> — 65 passed, 1 failed, 1 skipped</td></tr>
<tr><td>github-actions</td><td><a href="https://github.com/signalfx/splunk-otel-collector-chart/pull/2725">Bump nodejs from v4.13.0 to v4.14.0 in helm-charts/splunk-otel-collector/values.yaml (#2725)</a></td><td>PENDING — 66 passed, 1 skipped</td></tr>
</tbody>
</table>

## Failed

<table>
<thead>
<tr><th>Author</th><th>Pull Request</th><th>CI/CD</th></tr>
</thead>
<tbody>
<tr><td>github-actions</td><td><a href="https://github.com/signalfx/splunk-otel-collector-chart/pull/2644">Update matrix versions used for testing (#2644)</a></td><td><a href="#pr-2644">FAILURE</a> — 90 passed, 5 failed, 3 skipped, 8 cancelled</td></tr>
</tbody>
</table>

<a id="pr-2731"></a>
### [Bump java from v2.33.1 to v2.34.0 in helm-charts/splunk-otel-collector/values.yaml (#2731)](https://github.com/signalfx/splunk-otel-collector-chart/pull/2731)

Partial — the report is incomplete; findings use only available observations.
- Could not parse a failed workflow: could not parse workflow run 44527280670 attempt 2: could not read logs for workflow run 44527280670 attempt 2: GitHub API request to https://api.github.com/repos/signalfx/splunk-otel-collector-chart/actions/runs/44527280670/attempts/2/logs returned HTTP status 404 Not Found: {"message":"Not Found","documentation_url":"https://docs.github.com/rest/actions/workflow-runs#download-workflow-run-attempt-logs","status":"404"}<br>workflow run 44527280670 attempt 2 job "kubernetes-test (v1.35.5, gateway)" step "Run ./.github/actions/functional-test" has no downloadable log

<table>
<thead>
<tr><th>Failure</th><th>Classification</th><th>Latest failure</th><th>Main history</th><th>PR commit history</th><th>Related CI history</th></tr>
</thead>
<tbody>
<tr><td>functional-tests / kubernetes-test (v1.35.5, gateway) / Run ./.github/actions/functional-test</td><td>Broken</td><td><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44527280670/job/135694798050">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
</tbody>
</table>


<a id="pr-2728"></a>
### [Bump obi from 0.16.1 to 0.16.4 in helm-charts/splunk-otel-collector/Chart.yaml (#2728)](https://github.com/signalfx/splunk-otel-collector-chart/pull/2728)

Partial — the report is incomplete; findings use only available observations.
- Could not parse a failed workflow: could not parse workflow run 44520931137 attempt 2: could not read logs for workflow run 44520931137 attempt 2: GitHub API request to https://api.github.com/repos/signalfx/splunk-otel-collector-chart/actions/runs/44520931137/attempts/2/logs returned HTTP status 404 Not Found: {"message":"Not Found","documentation_url":"https://docs.github.com/rest/actions/workflow-runs#download-workflow-run-attempt-logs","status":"404"}<br>workflow run 44520931137 attempt 2 job "kubernetes-test (v1.35.5, configuration_switching)" step "Run ./.github/actions/functional-test" has no downloadable log

<table>
<thead>
<tr><th>Failure</th><th>Classification</th><th>Latest failure</th><th>Main history</th><th>PR commit history</th><th>Related CI history</th></tr>
</thead>
<tbody>
<tr><td>functional-tests / kubernetes-test (v1.35.5, configuration_switching) / Run ./.github/actions/functional-test</td><td>Broken</td><td><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44520931137/job/135681637983">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
</tbody>
</table>


<a id="pr-2727"></a>
### [Bump TargetAllocator from 0.162.0 to 0.163.0 in helm-charts/splunk-otel-collector/Chart.yaml (#2727)](https://github.com/signalfx/splunk-otel-collector-chart/pull/2727)

Partial — the report is incomplete; findings use only available observations.
- Could not parse a failed workflow: could not parse workflow run 44520942128 attempt 2: could not read logs for workflow run 44520942128 attempt 2: GitHub API request to https://api.github.com/repos/signalfx/splunk-otel-collector-chart/actions/runs/44520942128/attempts/2/logs returned HTTP status 404 Not Found: {"message":"Not Found","documentation_url":"https://docs.github.com/rest/actions/workflow-runs#download-workflow-run-attempt-logs","status":"404"}<br>workflow run 44520942128 attempt 2 job "kubernetes-test (v1.35.5, gateway)" step "Run ./.github/actions/functional-test" has no downloadable log

<table>
<thead>
<tr><th>Failure</th><th>Classification</th><th>Latest failure</th><th>Main history</th><th>PR commit history</th><th>Related CI history</th></tr>
</thead>
<tbody>
<tr><td>functional-tests / kubernetes-test (v1.35.5, gateway) / Run ./.github/actions/functional-test</td><td>Broken</td><td><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44520942128/job/135692349819">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
</tbody>
</table>


<a id="pr-2644"></a>
### [Update matrix versions used for testing (#2644)](https://github.com/signalfx/splunk-otel-collector-chart/pull/2644)

Partial — the report is incomplete; findings use only available observations.
- Could not parse a failed workflow: could not parse workflow run 43865328983 attempt 3: could not read logs for workflow run 43865328983 attempt 3: GitHub workflow attempt log archive for run 43865328983 attempt 3 has no logs for job "kubernetes-test (v1.37.0, histogram)"

<table>
<thead>
<tr><th>Failure</th><th>Classification</th><th>Latest failure</th><th>Main history</th><th>PR commit history</th><th>Related CI history</th></tr>
</thead>
<tbody>
<tr><td>build-and-test / checks / Check for vulnerabilities</td><td>Broken</td><td><pre><code>2026-11-17T16:20:38.1864709Z ##[error]      #27: internal/chart.go:256:38: internal.ChartUninstall calls action.Uninstall.Run, which eventually calls http.checkConnErrorWriter.Write
2026-11-17T16:20:38.1867041Z ##[error]      #28: internal/chart.go:256:38: internal.ChartUninstall calls action.Uninstall.Run, which eventually calls http.chunkWriter.Write
2026-11-17T16:20:38.1869070Z ##[error]      #29: internal/common.go:427:35: internal.tryGetPodLogs calls io.teeReader.Read, which eventually calls http.connReader.Read
2026-11-17T16:20:38.1871402Z ##[error]      #30: internal/common.go:427:35: internal.tryGetPodLogs calls http2.transportResponseBody.Read, which eventually calls http.expectContinueReader.Close
2026-11-17T16:20:38.1873388Z ##[error]      #31: internal/common.go:427:35: internal.tryGetPodLogs calls http.expectContinueReader.Read
2026-11-17T16:20:38.1874864Z ##[error]      #32: internal/common.go:422:2: internal.tryGetPodLogs calls http.gzipReader.Close
2026-11-17T16:20:38.1876444Z ##[error]      #33: internal/common.go:427:35: internal.tryGetPodLogs calls http.gzipReader.Read
2026-11-17T16:20:38.1878240Z ##[error]      #34: internal/common.go:427:35: internal.tryGetPodLogs calls io.teeReader.Read, which eventually calls http.loggingConn.Read
2026-11-17T16:20:38.1880259Z ##[error]      #35: internal/chart.go:256:38: internal.ChartUninstall calls action.Uninstall.Run, which eventually calls http.loggingConn.Write
2026-11-17T16:20:38.1882443Z ##[error]      #36: internal/common.go:427:35: internal.tryGetPodLogs calls http2.transportResponseBody.Read, which eventually calls http.maxBytesReader.Close
2026-11-17T16:20:38.1884157Z ##[error]      #37: internal/common.go:427:35: internal.tryGetPodLogs calls http.maxBytesReader.Read
2026-11-17T16:20:38.1887266Z ##[error]      #38: internal/sinks.go:241:31: internal.SetupOTLPLogsSinkOnPort calls sharedcomponent.Component[*go.opentelemetry.io/collector/receiver/otlpreceiver.otlpReceiver].Start, which eventually calls http.onceCloseListener.Close
2026-11-17T16:20:38.1889368Z ##[error]      #39: internal/common.go:427:35: internal.tryGetPodLogs calls http.persistConn.Read
2026-11-17T16:20:38.1891279Z ##[error]      #40: internal/chart.go:256:38: internal.ChartUninstall calls action.Uninstall.Run, which eventually calls http.persistConnWriter.ReadFrom
2026-11-17T16:20:38.1893540Z ##[error]      #41: internal/chart.go:256:38: internal.ChartUninstall calls action.Uninstall.Run, which eventually calls http.persistConnWriter.Write
2026-11-17T16:20:38.1895951Z ##[error]      #42: internal/common.go:427:35: internal.tryGetPodLogs calls http2.transportResponseBody.Read, which eventually calls http.readTrackingBody.Close
2026-11-17T16:20:38.1897810Z ##[error]      #43: internal/common.go:427:35: internal.tryGetPodLogs calls http.readTrackingBody.Read
2026-11-17T16:20:38.1899444Z ##[error]      #44: internal/common.go:427:35: internal.tryGetPodLogs calls http.readWriteCloserBody.Read
2026-11-17T16:20:38.1901481Z ##[error]      #45: internal/chart.go:87:20: internal.ChartInstallOrUpgrade calls template.Template.Execute, which eventually calls http.response.Write
2026-11-17T16:20:38.1903765Z ##[error]      #46: internal/chart.go:87:20: internal.ChartInstallOrUpgrade calls template.Template.Execute, which eventually calls http.response.WriteHeader
2026-11-17T16:20:38.1905687Z ##[error]      #47: internal/common.go:598:73: internal.DeleteObject calls http.transportReadFromServerError.Error
2026-11-17T16:20:38.1906710Z
2026-11-17T16:20:38.1908996Z Use &#39;-show verbose&#39; for more details.
2026-11-17T16:20:38.5283162Z make[2]: *** [../Makefile.common:70: govulncheck] Error 3
2026-11-17T16:20:38.5285818Z make[1]: *** [Makefile:241: functional_tests] Error 2
2026-11-17T16:20:38.5286695Z make[1]: Leaving directory &#39;/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart&#39;
2026-11-17T16:20:38.5288169Z make: *** [Makefile:260: govulncheck-all] Error 2
2026-11-17T16:20:38.5298873Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328964/job/134184195490">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>test github.com/signalfx/splunk-otel-collector-chart/functional_tests/functional.Test_Functions/kubernetes_cluster_metrics</td><td>Broken</td><td><pre><code>    kind_test.go:101: 
        	Error Trace:	/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/internal/metricassert.go:205
        	            				/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/functional/kind_test.go:101
        	Error:      	Received unexpected error:
        	            	resource map[]: scope name=&#34;&#34; version=&#34;&#34;: metric &#34;k8s.container.cpu_limit&#34;: missing datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;container.id/regex&#34;,&#34;(containerd://|cri-o://|docker://)?[0-9a-f]{64}&#34;],[&#34;container.image.name/regex&#34;,&#34;[-./:0-9a-z_]+&#34;],[&#34;container.image.tag/regex&#34;,&#34;[-.0-9A-Za-z_]+&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.container.name&#34;,&#34;kindnet-cni&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.uid/regex&#34;,&#34;[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	            	metric &#34;k8s.pod.status_reason&#34;: missing datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.qos_class&#34;,&#34;Guaranteed&#34;],[&#34;k8s.pod.uid/regex&#34;,&#34;[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	            	unexpected datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name&#34;,&#34;kind-control-plane&#34;],[&#34;k8s.pod.name&#34;,&#34;kube-scheduler-kind-control-plane&#34;],[&#34;k8s.pod.qos_class&#34;,&#34;Burstable&#34;],[&#34;k8s.pod.uid&#34;,&#34;2a52a4ed-4aa3-44e1-85df-707926028309&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	            	metric &#34;k8s.container.memory_limit&#34;: missing datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;container.id/regex&#34;,&#34;(containerd://|cri-o://|docker://)?[0-9a-f]{64}&#34;],[&#34;container.image.name/regex&#34;,&#34;[-./:0-9a-z_]+&#34;],[&#34;container.image.tag/regex&#34;,&#34;[-.0-9A-Za-z_]+&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.container.name&#34;,&#34;kindnet-cni&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.uid/regex&#34;,&#34;[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	Test:       	Test_Functions/kubernetes_cluster_metrics
        	Messages:   	Metric assertion failed for testdata/expected_kind_values/expected_cluster_receiver_assertion.yaml
=== NAME  Test_Functions
    resources_test.go:172: Deployed namespace ns-w-index
    resources_test.go:172: Deployed namespace ns-wo-index
    resources_test.go:172: Deployed namespace ns-w-exclude
    resources_test.go:196: Deployed job pod-w-index-w-ns-index
    resources_test.go:196: Deployed job pod-wo-index-w-ns-index
    resources_test.go:196: Deployed job pod-w-index-wo-ns-index
    resources_test.go:196: Deployed job pod-w-index-w-ns-exclude
    resources_test.go:196: Deployed job pod-w-exclude-wo-ns-exclude
    resources_test.go:196: Deployed job pod-wo-index-wo-ns-index</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328983/job/134248036608">Failure run</a></td><td>main — last 7 days — 15/15 passed (100%), 0/15 failed (0%); includes 9 inferred passes from successful workflows (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44342638394/job/134936542377">main pass</a></td><td><a href="https://github.com/signalfx/splunk-otel-collector-chart/commit/42178fa54cb93b72a891421d3334aa0436ba3b02">update-matrix-test-versions 42178fa</a> — 0/3 passed (0%), 3/3 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328983/job/134248036608">update-matrix-test-versions 42178fa failure</a></td><td>functional-tests — main — last 7 days — 9/11 passed (82%), 2/11 failed (18%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44200776895">functional-tests — main pass</a><br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44342638394">functional-tests — main failure</a><br>functional-tests — <a href="https://github.com/signalfx/splunk-otel-collector-chart/commit/42178fa54cb93b72a891421d3334aa0436ba3b02">update-matrix-test-versions 42178fa</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328983">update-matrix-test-versions 42178fa failure</a><br>functional-tests / kubernetes-test (v1.36.1, functional) — main — last 7 days — 11/11 passed (100%), 0/11 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44342638394/job/134936542377">functional-tests / kubernetes-test (v1.36.1, functional) — main pass</a><br>functional-tests / kubernetes-test (v1.36.1, functional) / Run ./.github/actions/functional-test — main — last 7 days — 11/11 passed (100%), 0/11 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44342638394/job/134936542377">functional-tests / kubernetes-test (v1.36.1, functional) / Run ./.github/actions/functional-test — main pass</a><br>functional-tests / kubernetes-test (v1.34.8, functional) — main — last 7 days — 11/11 passed (100%), 0/11 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44342638394/job/134936543058">functional-tests / kubernetes-test (v1.34.8, functional) — main pass</a><br>functional-tests / kubernetes-test (v1.34.8, functional) / Run ./.github/actions/functional-test — main — last 7 days — 11/11 passed (100%), 0/11 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44342638394/job/134936543058">functional-tests / kubernetes-test (v1.34.8, functional) / Run ./.github/actions/functional-test — main pass</a><br>functional-tests / kubernetes-test (v1.35.5, functional) — main — last 7 days — 11/11 passed (100%), 0/11 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44342638394/job/134936542928">functional-tests / kubernetes-test (v1.35.5, functional) — main pass</a><br>functional-tests / kubernetes-test (v1.35.5, functional) / Run ./.github/actions/functional-test — main — last 7 days — 11/11 passed (100%), 0/11 failed (0%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44342638394/job/134936542928">functional-tests / kubernetes-test (v1.35.5, functional) / Run ./.github/actions/functional-test — main pass</a><br>functional-tests / kubernetes-test (v1.36.4, functional) — main — last 7 days — 9/9 passed (100%), 0/9 failed (0%); includes 9 inferred passes from successful workflows (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44200776895">functional-tests / kubernetes-test (v1.36.4, functional) — main pass</a><br>functional-tests / kubernetes-test (v1.36.4, functional) — <a href="https://github.com/signalfx/splunk-otel-collector-chart/commit/42178fa54cb93b72a891421d3334aa0436ba3b02">update-matrix-test-versions 42178fa</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328983/job/134248036608">update-matrix-test-versions 42178fa failure</a><br>functional-tests / kubernetes-test (v1.36.4, functional) / Run ./.github/actions/functional-test — main — last 7 days — 9/9 passed (100%), 0/9 failed (0%); includes 9 inferred passes from successful workflows (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44200776895">functional-tests / kubernetes-test (v1.36.4, functional) / Run ./.github/actions/functional-test — main pass</a><br>functional-tests / kubernetes-test (v1.36.4, functional) / Run ./.github/actions/functional-test — <a href="https://github.com/signalfx/splunk-otel-collector-chart/commit/42178fa54cb93b72a891421d3334aa0436ba3b02">update-matrix-test-versions 42178fa</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328983/job/134248036608">update-matrix-test-versions 42178fa failure</a><br>functional-tests / kubernetes-test (v1.35.8, functional) — main — last 7 days — 9/9 passed (100%), 0/9 failed (0%); includes 9 inferred passes from successful workflows (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44200776895">functional-tests / kubernetes-test (v1.35.8, functional) — main pass</a><br>functional-tests / kubernetes-test (v1.35.8, functional) — <a href="https://github.com/signalfx/splunk-otel-collector-chart/commit/42178fa54cb93b72a891421d3334aa0436ba3b02">update-matrix-test-versions 42178fa</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328983/job/134248036761">update-matrix-test-versions 42178fa failure</a><br>functional-tests / kubernetes-test (v1.35.8, functional) / Run ./.github/actions/functional-test — main — last 7 days — 9/9 passed (100%), 0/9 failed (0%); includes 9 inferred passes from successful workflows (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44200776895">functional-tests / kubernetes-test (v1.35.8, functional) / Run ./.github/actions/functional-test — main pass</a><br>functional-tests / kubernetes-test (v1.35.8, functional) / Run ./.github/actions/functional-test — <a href="https://github.com/signalfx/splunk-otel-collector-chart/commit/42178fa54cb93b72a891421d3334aa0436ba3b02">update-matrix-test-versions 42178fa</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328983/job/134248036761">update-matrix-test-versions 42178fa failure</a><br>functional-tests / kubernetes-test (v1.37.0, functional) — main — last 7 days — 9/9 passed (100%), 0/9 failed (0%); includes 9 inferred passes from successful workflows (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44200776895">functional-tests / kubernetes-test (v1.37.0, functional) — main pass</a><br>functional-tests / kubernetes-test (v1.37.0, functional) — <a href="https://github.com/signalfx/splunk-otel-collector-chart/commit/42178fa54cb93b72a891421d3334aa0436ba3b02">update-matrix-test-versions 42178fa</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328983/job/134248036825">update-matrix-test-versions 42178fa failure</a><br>functional-tests / kubernetes-test (v1.37.0, functional) / Run ./.github/actions/functional-test — main — last 7 days — 9/9 passed (100%), 0/9 failed (0%); includes 9 inferred passes from successful workflows (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44200776895">functional-tests / kubernetes-test (v1.37.0, functional) / Run ./.github/actions/functional-test — main pass</a><br>functional-tests / kubernetes-test (v1.37.0, functional) / Run ./.github/actions/functional-test — <a href="https://github.com/signalfx/splunk-otel-collector-chart/commit/42178fa54cb93b72a891421d3334aa0436ba3b02">update-matrix-test-versions 42178fa</a> — 0/1 passed (0%), 1/1 failed (100%) (partial)<br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328983/job/134248036825">update-matrix-test-versions 42178fa failure</a></td></tr>
<tr><td>functional-tests / kubernetes-test (v1.35.8, functional) / Run ./.github/actions/functional-test</td><td>Broken<br>Includes inferred workflow passes.</td><td><pre><code>=== RUN   Test_Functions/kubernetes_cluster_metrics
    kind_test.go:101: 
        	Error Trace:	/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/internal/metricassert.go:205
        	            				/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/functional/kind_test.go:101
        	Error:      	Received unexpected error:
        	            	resource map[]: scope name=&#34;&#34; version=&#34;&#34;: metric &#34;k8s.pod.status_reason&#34;: missing datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.qos_class&#34;,&#34;Guaranteed&#34;],[&#34;k8s.pod.uid/regex&#34;,&#34;[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	            	unexpected datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name&#34;,&#34;kind-control-plane&#34;],[&#34;k8s.pod.name&#34;,&#34;kube-scheduler-kind-control-plane&#34;],[&#34;k8s.pod.qos_class&#34;,&#34;Burstable&#34;],[&#34;k8s.pod.uid&#34;,&#34;4580a3de-486d-4115-8e39-947377c3f534&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	            	metric &#34;k8s.container.cpu_limit&#34;: missing datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;container.id/regex&#34;,&#34;(containerd://|cri-o://|docker://)?[0-9a-f]{64}&#34;],[&#34;container.image.name/regex&#34;,&#34;[-./:0-9a-z_]+&#34;],[&#34;container.image.tag/regex&#34;,&#34;[-.0-9A-Za-z_]+&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.container.name&#34;,&#34;kindnet-cni&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.uid/regex&#34;,&#34;[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	            	metric &#34;k8s.container.memory_limit&#34;: missing datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;container.id/regex&#34;,&#34;(containerd://|cri-o://|docker://)?[0-9a-f]{64}&#34;],[&#34;container.image.name/regex&#34;,&#34;[-./:0-9a-z_]+&#34;],[&#34;container.image.tag/regex&#34;,&#34;[-.0-9A-Za-z_]+&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.container.name&#34;,&#34;kindnet-cni&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.uid/regex&#34;,&#34;[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	Test:       	Test_Functions/kubernetes_cluster_metrics
        	Messages:   	Metric assertion failed for testdata/expected_kind_values/expected_cluster_receiver_assertion.yaml
=== NAME  Test_Functions
    resources_test.go:172: Deployed namespace ns-w-index
    resources_test.go:172: Deployed namespace ns-wo-index
    resources_test.go:172: Deployed namespace ns-w-exclude
    resources_test.go:196: Deployed job pod-w-index-w-ns-index
    resources_test.go:196: Deployed job pod-wo-index-w-ns-index
    resources_test.go:196: Deployed job pod-w-index-wo-ns-index
    resources_test.go:196: Deployed job pod-w-index-w-ns-exclude
    resources_test.go:196: Deployed job pod-w-exclude-wo-ns-exclude
    resources_test.go:196: Deployed job pod-wo-index-wo-ns-index</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328983/job/134248036761">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>functional-tests / kubernetes-test (v1.36.4, functional) / Run ./.github/actions/functional-test</td><td>Broken<br>Includes inferred workflow passes.</td><td><pre><code>=== RUN   Test_Functions/kubernetes_cluster_metrics
    kind_test.go:101: 
        	Error Trace:	/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/internal/metricassert.go:205
        	            				/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/functional/kind_test.go:101
        	Error:      	Received unexpected error:
        	            	resource map[]: scope name=&#34;&#34; version=&#34;&#34;: metric &#34;k8s.container.cpu_limit&#34;: missing datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;container.id/regex&#34;,&#34;(containerd://|cri-o://|docker://)?[0-9a-f]{64}&#34;],[&#34;container.image.name/regex&#34;,&#34;[-./:0-9a-z_]+&#34;],[&#34;container.image.tag/regex&#34;,&#34;[-.0-9A-Za-z_]+&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.container.name&#34;,&#34;kindnet-cni&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.uid/regex&#34;,&#34;[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	            	metric &#34;k8s.pod.status_reason&#34;: missing datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.qos_class&#34;,&#34;Guaranteed&#34;],[&#34;k8s.pod.uid/regex&#34;,&#34;[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	            	unexpected datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name&#34;,&#34;kind-control-plane&#34;],[&#34;k8s.pod.name&#34;,&#34;kube-scheduler-kind-control-plane&#34;],[&#34;k8s.pod.qos_class&#34;,&#34;Burstable&#34;],[&#34;k8s.pod.uid&#34;,&#34;2a52a4ed-4aa3-44e1-85df-707926028309&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	            	metric &#34;k8s.container.memory_limit&#34;: missing datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;container.id/regex&#34;,&#34;(containerd://|cri-o://|docker://)?[0-9a-f]{64}&#34;],[&#34;container.image.name/regex&#34;,&#34;[-./:0-9a-z_]+&#34;],[&#34;container.image.tag/regex&#34;,&#34;[-.0-9A-Za-z_]+&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.container.name&#34;,&#34;kindnet-cni&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.uid/regex&#34;,&#34;[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	Test:       	Test_Functions/kubernetes_cluster_metrics
        	Messages:   	Metric assertion failed for testdata/expected_kind_values/expected_cluster_receiver_assertion.yaml
=== NAME  Test_Functions
    resources_test.go:172: Deployed namespace ns-w-index
    resources_test.go:172: Deployed namespace ns-wo-index
    resources_test.go:172: Deployed namespace ns-w-exclude
    resources_test.go:196: Deployed job pod-w-index-w-ns-index
    resources_test.go:196: Deployed job pod-wo-index-w-ns-index
    resources_test.go:196: Deployed job pod-w-index-wo-ns-index
    resources_test.go:196: Deployed job pod-w-index-w-ns-exclude
    resources_test.go:196: Deployed job pod-w-exclude-wo-ns-exclude
    resources_test.go:196: Deployed job pod-wo-index-wo-ns-index</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328983/job/134248036608">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>functional-tests / kubernetes-test (v1.37.0, functional) / Run ./.github/actions/functional-test</td><td>Broken<br>Includes inferred workflow passes.</td><td><pre><code>=== RUN   Test_Functions/kubernetes_cluster_metrics
    kind_test.go:101: 
        	Error Trace:	/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/internal/metricassert.go:205
        	            				/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/functional/kind_test.go:101
        	Error:      	Received unexpected error:
        	            	resource map[]: scope name=&#34;&#34; version=&#34;&#34;: metric &#34;k8s.container.memory_limit&#34;: missing datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;container.id/regex&#34;,&#34;(containerd://|cri-o://|docker://)?[0-9a-f]{64}&#34;],[&#34;container.image.name/regex&#34;,&#34;[-./:0-9a-z_]+&#34;],[&#34;container.image.tag/regex&#34;,&#34;[-.0-9A-Za-z_]+&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.container.name&#34;,&#34;kindnet-cni&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.uid/regex&#34;,&#34;[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	            	metric &#34;k8s.container.cpu_limit&#34;: missing datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;container.id/regex&#34;,&#34;(containerd://|cri-o://|docker://)?[0-9a-f]{64}&#34;],[&#34;container.image.name/regex&#34;,&#34;[-./:0-9a-z_]+&#34;],[&#34;container.image.tag/regex&#34;,&#34;[-.0-9A-Za-z_]+&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.container.name&#34;,&#34;kindnet-cni&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.uid/regex&#34;,&#34;[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	            	metric &#34;k8s.pod.status_reason&#34;: missing datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.name/regex&#34;,&#34;[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*&#34;],[&#34;k8s.pod.qos_class&#34;,&#34;Guaranteed&#34;],[&#34;k8s.pod.uid/regex&#34;,&#34;[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	            	unexpected datapoint with attributes [[&#34;cluster_name&#34;,&#34;ci-k8s-cluster&#34;],[&#34;customfield1&#34;,&#34;customvalue1&#34;],[&#34;customfield2&#34;,&#34;customvalue2&#34;],[&#34;k8s.cluster.name&#34;,&#34;dev-operator&#34;],[&#34;k8s.namespace.name&#34;,&#34;kube-system&#34;],[&#34;k8s.node.name&#34;,&#34;kind-control-plane&#34;],[&#34;k8s.pod.name&#34;,&#34;kube-scheduler-kind-control-plane&#34;],[&#34;k8s.pod.qos_class&#34;,&#34;Burstable&#34;],[&#34;k8s.pod.uid&#34;,&#34;9a4cc459-d876-43b2-96b5-359af59b3396&#34;],[&#34;metric_source&#34;,&#34;kubernetes&#34;],[&#34;receiver&#34;,&#34;k8scluster&#34;]]
        	Test:       	Test_Functions/kubernetes_cluster_metrics
        	Messages:   	Metric assertion failed for testdata/expected_kind_values/expected_cluster_receiver_assertion.yaml
=== NAME  Test_Functions
    resources_test.go:172: Deployed namespace ns-w-index
    resources_test.go:172: Deployed namespace ns-wo-index
    resources_test.go:172: Deployed namespace ns-w-exclude
    resources_test.go:196: Deployed job pod-w-index-w-ns-index
    resources_test.go:196: Deployed job pod-wo-index-w-ns-index
    resources_test.go:196: Deployed job pod-w-index-wo-ns-index
    resources_test.go:196: Deployed job pod-w-index-w-ns-exclude
    resources_test.go:196: Deployed job pod-w-exclude-wo-ns-exclude
    resources_test.go:196: Deployed job pod-wo-index-wo-ns-index</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328983/job/134248036825">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
<tr><td>functional-tests / kubernetes-test (v1.37.0, histogram)</td><td>Broken<br>Includes inferred workflow passes.</td><td><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/43865328983/job/134248036670">Failure run</a></td><td>—</td><td>—</td><td>—</td></tr>
</tbody>
</table>


