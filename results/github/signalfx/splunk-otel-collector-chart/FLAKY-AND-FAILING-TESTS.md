# Flaky and Failing Tests

Report generated: 2026-11-19 00:13:27 UTC

Analysis status: incomplete

Main branch observations: 2026-11-17 through 2026-11-18

This report is partial; counts, rates, and classifications use the available observations. A test marked Broken may have successful runs that were not gathered.

## Flaky tests

<table>
<thead>
<tr><th>Workflow</th><th>Job</th><th>Step</th><th>Test</th><th>Classification / reason</th><th>Failures / observations</th><th>Failure rate</th><th>Recent failure</th></tr>
</thead>
<tbody>
<tr><td>functional-tests-gce</td><td>—</td><td>—</td><td>test github.com/signalfx/splunk-otel-collector-chart/functional_tests/functional.Test_Functions/control_plane_receiver_checks/receiver_started:_prometheus/kubernetes-apiserver</td><td>Flaky</td><td>1 / 5</td><td>20%</td><td><pre><code>    health_test.go:55: 
        	Error Trace:	/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/internal/component_health.go:228
        	            				/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/functional/health_test.go:55
        	Error:      	Should be true
        	Test:       	Test_Functions/control_plane_receiver_checks/receiver_started:_prometheus/kubernetes-apiserver
        	Messages:   	receiver prometheus/kubernetes-apiserver was not started by receiver_creator on any of the 1 agent pod(s)</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44130432839/job/134205738706">Failure run</a></td></tr>
<tr><td>functional-tests-gce</td><td>—</td><td>—</td><td>test github.com/signalfx/splunk-otel-collector-chart/functional_tests/functional.Test_Functions/control_plane_receiver_checks/receiver_started:_prometheus/kubernetes-scheduler</td><td>Flaky</td><td>1 / 5</td><td>20%</td><td><pre><code>    health_test.go:55: 
        	Error Trace:	/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/internal/component_health.go:228
        	            				/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/functional/health_test.go:55
        	Error:      	Should be true
        	Test:       	Test_Functions/control_plane_receiver_checks/receiver_started:_prometheus/kubernetes-scheduler
        	Messages:   	receiver prometheus/kubernetes-scheduler was not started by receiver_creator on any of the 1 agent pod(s)</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44130432839/job/134205738706">Failure run</a></td></tr>
</tbody>
</table>

## Flaky workflow and job failures

<table>
<thead>
<tr><th>Workflow</th><th>Job</th><th>Step</th><th>Classification</th><th>Failures / observations</th><th>Failure rate</th><th>Recent failure</th></tr>
</thead>
<tbody>
<tr><td>functional-tests</td><td>Test helm upgrade in kind</td><td>Run ./.github/actions/functional-test</td><td>Flaky</td><td>1 / 11</td><td>9%</td><td><pre><code>2026-11-17T17:49:44.9816912Z Deleted nodes: [&#34;kind-control-plane&#34;]
2026-11-17T17:49:44.9819457Z ERROR: failed to create cluster: command &#34;docker run --name kind-control-plane --hostname kind-control-plane --label io.x-k8s.kind.role=control-plane --privileged --security-opt seccomp=unconfined --security-opt apparmor=unconfined --tmpfs /tmp --tmpfs /run --volume /var --volume /lib/modules:/lib/modules:ro -e KIND_EXPERIMENTAL_CONTAINERD_SNAPSHOTTER --detach --tty --label io.x-k8s.kind.cluster=kind --net kind --restart=on-failure:1 --init=false --cgroupns=private --publish=0.0.0.0:80:80/TCP --publish=0.0.0.0:443:443/TCP --publish=0.0.0.0:43180:4318/TCP --publish=127.0.0.1:42513:6443/TCP -e KUBECONFIG=/etc/kubernetes/admin.conf kindest/node:v1.36.1&#34; failed with error: exit status 125
2026-11-17T17:49:44.9821232Z
2026-11-17T17:49:44.9821443Z Command Output: 1b198032662cd757fefdfb5602d4008094583672b3a1e50853f0e1bf6715c83a
2026-11-17T17:49:44.9822463Z docker: Error response from daemon: failed to set up container networking: driver failed programming external connectivity on endpoint kind-control-plane (d0c0b86295ad5daee0b4850e72d78f3f4b30025c998e2c2b6903e40d589d36fd): failed to bind host port 0.0.0.0:43180/tcp: address already in use
2026-11-17T17:49:44.9823376Z
2026-11-17T17:49:44.9823493Z Run &#39;docker run --help&#39; for more information
2026-11-17T17:49:44.9872718Z ##[end-action id=__self.__helm_kind-action;outcome=failure;conclusion=failure;duration_ms=22873]
2026-11-17T17:49:44.9876753Z ##[start-action display=Fix kubelet TLS server certificates;id=__self.__run_2]</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44141111392/job/134243834324">Failure run</a></td></tr>
<tr><td>functional-tests</td><td>kubernetes-test (1 matrix variant)</td><td>Run ./.github/actions/functional-test</td><td>Flaky</td><td>1 / 11</td><td>9%</td><td><pre><code>2026-11-18T16:03:43.7098477Z Deleted nodes: [&#34;kind-control-plane&#34;]
2026-11-18T16:03:43.7100184Z ERROR: failed to create cluster: command &#34;docker run --name kind-control-plane --hostname kind-control-plane --label io.x-k8s.kind.role=control-plane --privileged --security-opt seccomp=unconfined --security-opt apparmor=unconfined --tmpfs /tmp --tmpfs /run --volume /var --volume /lib/modules:/lib/modules:ro -e KIND_EXPERIMENTAL_CONTAINERD_SNAPSHOTTER --detach --tty --label io.x-k8s.kind.cluster=kind --net kind --restart=on-failure:1 --init=false --cgroupns=private --publish=0.0.0.0:80:80/TCP --publish=0.0.0.0:443:443/TCP --publish=0.0.0.0:43180:4318/TCP --publish=127.0.0.1:46215:6443/TCP -e KUBECONFIG=/etc/kubernetes/admin.conf kindest/node:v1.34.8&#34; failed with error: exit status 125
2026-11-18T16:03:43.7101645Z
2026-11-18T16:03:43.7101799Z Command Output: 71555ae8a9ab0afc6098783b46ed911eaf3c6a17a406004e61ff11fc4490189f
2026-11-18T16:03:43.7102483Z docker: Error response from daemon: failed to set up container networking: driver failed programming external connectivity on endpoint kind-control-plane (e69547c8a33387f15f21ab2ae62ed4159599701626a478760cf12ac04b8e8c96): failed to bind host port 0.0.0.0:43180/tcp: address already in use
2026-11-18T16:03:43.7102993Z
2026-11-18T16:03:43.7103069Z Run &#39;docker run --help&#39; for more information
2026-11-18T16:03:43.7157036Z ##[end-action id=__self.__helm_kind-action;outcome=failure;conclusion=failure;duration_ms=20332]
2026-11-18T16:03:43.7160282Z ##[start-action display=Fix kubelet TLS server certificates;id=__self.__run_2]</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44342638394/job/134936542790">Failure run</a></td></tr>
<tr><td>functional-tests-gce</td><td>Test helm install in GCE (kops) - credentials needed</td><td>Run ./.github/actions/cloud-functional-test</td><td>Flaky</td><td>1 / 5</td><td>20%</td><td><pre><code>=== RUN   Test_Functions/control_plane_receiver_checks/receiver_started:_prometheus/kubernetes-scheduler
    health_test.go:55: 
        	Error Trace:	/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/internal/component_health.go:228
        	            				/home/runner/work/splunk-otel-collector-chart/splunk-otel-collector-chart/functional_tests/functional/health_test.go:55
        	Error:      	Should be true
        	Test:       	Test_Functions/control_plane_receiver_checks/receiver_started:_prometheus/kubernetes-scheduler
        	Messages:   	receiver prometheus/kubernetes-scheduler was not started by receiver_creator on any of the 1 agent pod(s)</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44130432839/job/134205738706">Failure run</a></td></tr>
</tbody>
</table>

## Broken workflow and job failures

<table>
<thead>
<tr><th>Workflow</th><th>Job</th><th>Step</th><th>Classification</th><th>Failures / observations</th><th>Failure rate</th><th>Recent failure</th></tr>
</thead>
<tbody>
<tr><td>Check for new chart dependency updates</td><td>check_and_update (operator, operator, helm-charts/splunk-otel-collector/Chart.yaml, opentelemetry...</td><td>Check for Version Updates</td><td>Broken</td><td>1 / 1</td><td>100%</td><td><pre><code>2026-11-18T12:07:10.4992862Z Opentelemetry Operator Subchart Version: 0.129.1
2026-11-18T12:07:11.0179848Z Error: Unknown instrumentation library &#39;php&#39;. Please update get_repository_for_instrumentation function.
2026-11-18T12:07:11.0203671Z Available keys in versions.txt:
2026-11-18T12:07:11.0205731Z autoinstrumentation-nginx
2026-11-18T12:07:11.0214525Z make: *** [Makefile:198: update-chart-dep] Error 1
2026-11-18T12:07:11.0232180Z ##[error]Process completed with exit code 2.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector-chart/actions/runs/44300847748/job/134789866287">Failure run</a></td></tr>
</tbody>
</table>

