# Flaky and Failing Tests

Report generated: 2026-11-19 20:47:21 UTC

Analysis status: incomplete

Main branch observations: 2026-11-18 through 2026-11-19

This report is partial; counts, rates, and classifications use the available observations. A test marked Broken may have successful runs that were not gathered.

## Flaky tests

<table>
<thead>
<tr><th>Workflow</th><th>Job</th><th>Step</th><th>Test</th><th>Classification / reason</th><th>Failures / observations</th><th>Failure rate</th><th>Recent failure</th></tr>
</thead>
<tbody>
<tr><td>Integration Test</td><td>—</td><td>—</td><td>test github.com/signalfx/splunk-otel-collector/tests/general.TestDefaultConfigsDoNotEnableExpvar/agent</td><td>Flaky</td><td>2 / 17</td><td>12%</td><td><pre><code>    default_config_test.go:643: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/general/default_config_test.go:644
        	            				/opt/hostedtoolcache/go/1.26.11/x64/src/runtime/asm_amd64.s:1771
        	Error:      	Not equal: 
        	            	expected: 200
        	            	actual  : 404
    default_config_test.go:643: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/general/default_config_test.go:643
        	Error:      	Condition never satisfied
        	Test:       	TestDefaultConfigsDoNotEnableExpvar/agent</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44242773739/job/134600997092">Failure run</a></td></tr>
<tr><td>Integration Test</td><td>—</td><td>—</td><td>test github.com/signalfx/splunk-otel-collector/tests/general.TestDefaultConfigsDoNotEnableExpvar/gateway</td><td>Flaky</td><td>2 / 17</td><td>12%</td><td><pre><code>    default_config_test.go:643: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/general/default_config_test.go:644
        	            				/opt/hostedtoolcache/go/1.26.11/x64/src/runtime/asm_amd64.s:1771
        	Error:      	Not equal: 
        	            	expected: 200
        	            	actual  : 404
    default_config_test.go:643: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/general/default_config_test.go:643
        	Error:      	Condition never satisfied
        	Test:       	TestDefaultConfigsDoNotEnableExpvar/gateway</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44242773739/job/134600997092">Failure run</a></td></tr>
<tr><td>Integration Test</td><td>—</td><td>—</td><td>test github.com/signalfx/splunk-otel-collector/tests/general/discoverymode.TestConfigDInitialAndEffectiveConfig</td><td>Flaky</td><td>1 / 17</td><td>6%</td><td><pre><code>    testcase.go:221: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/testutils/testcase.go:221
        	            				/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/testutils/testcase.go:151
        	            				/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/general/discoverymode/configd_test.go:38
        	Error:      	Received unexpected error:
        	            	start container: started hook: wait until ready: container exited with code 1
        	Test:       	TestConfigDInitialAndEffectiveConfig</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44242773739/job/134600997394">Failure run</a></td></tr>
<tr><td>Integration Test</td><td>—</td><td>—</td><td>test github.com/signalfx/splunk-otel-collector/tests/telemetry.TestInternalTelemetryPrometheusNames</td><td>Flaky</td><td>2 / 15</td><td>13%</td><td><pre><code>    golden.go:188: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/testutils/golden.go:197
        	            				/opt/hostedtoolcache/go/1.26.11/x64/src/runtime/asm_amd64.s:1771
        	Error:      	Received unexpected error:
        	            	missing expected resource: map[server.port:8888 service.instance.id/exists:true service.name/exists:true service.version/exists:true url.scheme:http]
        	            	unexpected resource: map[server.port:8888 service.instance.id:localhost:8888 service.name:internal-telemetry service_instance_id:4fe0fa5d-02fe-44ec-8270-c6f0786f1fd6 service_name:otelcol service_version:v0.165.0-8-gcb3ed37f2 url.scheme:http]
    golden.go:188: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/testutils/golden.go:188
        	            				/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/telemetry/internal_telemetry_test.go:34
        	Error:      	Condition never satisfied
        	Test:       	TestInternalTelemetryPrometheusNames</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44394053528/job/135129368711">Failure run</a></td></tr>
<tr><td>Integration Test</td><td>—</td><td>—</td><td>test github.com/signalfx/splunk-otel-collector/tests/telemetry.TestInternalTelemetryPrometheusReaderNames</td><td>Flaky</td><td>2 / 15</td><td>13%</td><td><pre><code>    golden.go:188: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/testutils/golden.go:197
        	            				/opt/hostedtoolcache/go/1.26.11/x64/src/runtime/asm_amd64.s:1771
        	Error:      	Received unexpected error:
        	            	missing expected resource: map[host.name/exists:true os.type/exists:true server.port:8889 service.instance.id/exists:true service.name/exists:true service.version/exists:true url.scheme:http]
        	            	unexpected resource: map[host.name:runnervm8axu0.ldrbiddpfe1evawpvuyaexkoyh.gx.internal.cloudapp.net os.type:linux server.port:8889 service.instance.id:localhost:8889 service.name:internal-telemetry service_instance_id:2d9c63ac-c47a-447c-b2a0-972ca5e0353c service_name:otelcol service_version:v0.165.0-8-gcb3ed37f2 url.scheme:http]
    golden.go:188: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/testutils/golden.go:188
        	            				/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/telemetry/internal_telemetry_test.go:40
        	Error:      	Condition never satisfied
        	Test:       	TestInternalTelemetryPrometheusReaderNames</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44394053528/job/135129368711">Failure run</a></td></tr>
<tr><td>auto-instrumentation</td><td>—</td><td>—</td><td>test github.com/signalfx/splunk-otel-collector/tests/instrumentation.TestExpressInstrumentation/ubuntu-jammy/amd64</td><td>Flaky</td><td>1 / 16</td><td>6%</td><td><pre><code>    instrumentation_test.go:145: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/instrumentation/instrumentation_test.go:373
        	            				/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/instrumentation/instrumentation_test.go:145
        	Error:      	Should be true
        	Test:       	TestExpressInstrumentation/ubuntu-jammy/amd64
        	Messages:   	timed out waiting for &#34;telemetry\\.sdk\\.language&#34;: Str\(nodejs\)</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44594002584/job/135814712835">Failure run</a></td></tr>
<tr><td>auto-instrumentation</td><td>—</td><td>—</td><td>test github.com/signalfx/splunk-otel-collector/tests/instrumentation.TestPackageUninstall/debian-trixie/amd64</td><td>Flaky</td><td>1 / 20</td><td>5%</td><td><pre><code>    instrumentation_test.go:185: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/instrumentation/instrumentation_test.go:329
        	            				/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/instrumentation/instrumentation_test.go:185
        	Error:      	Received unexpected error:
        	            	create container: build image: The command &#39;/bin/sh -c apt-get update &amp;&amp;     apt-get install -yq libicu76 &amp;&amp;     bash /opt/setup-dotnet.sh &#34;$TARGETARCH&#34;&#39; returned a non-zero code: 4
        	Test:       	TestPackageUninstall/debian-trixie/amd64</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44447876884/job/135309265164">Failure run</a></td></tr>
<tr><td>auto-instrumentation</td><td>—</td><td>—</td><td>test github.com/signalfx/splunk-otel-collector/tests/instrumentation.TestTomcatInstrumentation/oraclelinux-8/amd64</td><td>Flaky</td><td>2 / 16</td><td>13%</td><td><pre><code>    instrumentation_test.go:120: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/instrumentation/instrumentation_test.go:373
        	            				/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/instrumentation/instrumentation_test.go:120
        	Error:      	Should be true
        	Test:       	TestTomcatInstrumentation/oraclelinux-8/amd64
        	Messages:   	timed out waiting for &#34;deployment\\.environment\\.name&#34;: Str\(deployment_environment_from_java\)</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44594002584/job/135814705041">Failure run</a></td></tr>
<tr><td>build-and-test</td><td>—</td><td>—</td><td>test github.com/signalfx/splunk-otel-collector/internal/receiver/gnmireceiver.TestTunnelMultipleTargetsConcurrently</td><td>Flaky</td><td>1 / 54</td><td>2%</td><td><pre><code>    receiver_test.go:489: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/internal/receiver/gnmireceiver/receiver_test.go:489
        	Error:      	Should be true
        	Test:       	TestTunnelMultipleTargetsConcurrently
        	Messages:   	expected metrics from both concurrently-registered targets, got map[routerA:true]</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44402506765/job/135149269969">Failure run</a></td></tr>
</tbody>
</table>

## Flaky workflow and job failures

<table>
<thead>
<tr><th>Workflow</th><th>Job</th><th>Step</th><th>Classification</th><th>Failures / observations</th><th>Failure rate</th><th>Recent failure</th></tr>
</thead>
<tbody>
<tr><td>.github/workflows/lychee.yml</td><td>linkChecker</td><td>lychee Link Checker</td><td>Flaky</td><td>5 / 6</td><td>83%</td><td><pre><code>2026-11-19T14:45:55.6000621Z [200] https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.119.0
2026-11-19T14:45:55.6001650Z [IGNORED] passthrough:///host:port | Unsupported: Error creating request client: builder error for url (passthrough:///host:port)
2026-11-19T14:45:55.6002519Z [EXCLUDED] https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/36239
2026-11-19T14:45:57.0786860Z [EXCLUDED] https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/45352
2026-11-19T14:45:57.1605936Z [404] https://github.com/open-telemetry/opentelemetry-java-contrib/tree/main/jmx-metrics | Network error: Not Found
2026-11-19T14:45:57.1606771Z [EXCLUDED] https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/43621
2026-11-19T14:46:00.0901721Z [EXCLUDED] https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/40720
2026-11-19T14:46:00.0902558Z [IGNORED] passthrough:/// | Unsupported: Error creating request client: builder error for url (passthrough:///)
2026-11-19T14:46:00.0903452Z [EXCLUDED] https://github.com/signalfx/splunk-otel-collector/pull/7205
2026-11-19T14:46:01.5585893Z [EXCLUDED] https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/34567
2026-11-19T14:46:01.5586347Z [IGNORED] dns:/// | Unsupported: Error creating request client: builder error for url (dns:///)
2026-11-19T14:46:01.5586784Z [EXCLUDED] https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/50307
2026-11-19T14:46:02.2900981Z [EXCLUDED] https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/51139
2026-11-19T14:46:02.2901175Z [IGNORED] dns:/// | Unsupported: Error creating request client: builder error for url (dns:///)
2026-11-19T14:46:02.2901358Z [EXCLUDED] https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/22750
2026-11-19T14:46:02.2954081Z [EXCLUDED] https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/37390
2026-11-19T14:46:02.2954414Z [404] https://github.com/open-telemetry/opentelemetry-java-contrib/tree/main/jmx-metrics | Error (cached)
2026-11-19T14:46:02.2954590Z [EXCLUDED] https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/35577
2026-11-19T14:46:02.3058692Z [EXCLUDED] https://github.com/open-telemetry/opentelemetry-collector/issues/14492
2026-11-19T14:46:02.3058952Z [IGNORED] passthrough:/// | Unsupported: Error creating request client: builder error for url (passthrough:///)
2026-11-19T14:46:02.3059098Z [EXCLUDED] https://github.com/signalfx/splunk-otel-collector/pull/7078
2026-11-19T14:46:02.4306553Z [EXCLUDED] https://github.com/signalfx/splunk-otel-collector/pull/7249
2026-11-19T14:46:02.4306858Z [IGNORED] passthrough:///host:port | Unsupported: Error creating request client: builder error for url (passthrough:///host:port)
2026-11-19T14:46:02.4307052Z [EXCLUDED] https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/42605
2026-11-19T14:46:09.0321378Z
2026-11-19T14:46:09.0322900Z * [404] [https://github.com/open-telemetry/opentelemetry-java-contrib/tree/main/jmx-metrics](https://github.com/open-telemetry/opentelemetry-java-contrib/tree/main/jmx-metrics) | Error (cached)
2026-11-19T14:46:09.0324038Z * [404] [https://github.com/open-telemetry/opentelemetry-java-contrib/tree/main/jmx-metrics](https://github.com/open-telemetry/opentelemetry-java-contrib/tree/main/jmx-metrics) | Network error: Not Found
2026-11-19T14:46:09.0324626Z
2026-11-19T14:46:09.0354571Z ##[error]Process completed with exit code 2.
2026-11-19T14:46:09.0372519Z ##[end-action id=lychee.run-lychee;outcome=failure;conclusion=failure;duration_ms=21438]</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44547248768/job/135641007380">Failure run</a></td></tr>
<tr><td>Integration Test</td><td>integration-test-docker (2 matrix variants)</td><td>Run Integration Test With Cover</td><td>Flaky</td><td>2 / 15</td><td>13%</td><td><pre><code>=== RUN   TestInternalTelemetryPrometheusReaderNames
    golden.go:188: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/testutils/golden.go:197
        	            				/opt/hostedtoolcache/go/1.26.11/x64/src/runtime/asm_amd64.s:1771
        	Error:      	Received unexpected error:
        	            	missing expected resource: map[host.name/exists:true os.type/exists:true server.port:8889 service.instance.id/exists:true service.name/exists:true service.version/exists:true url.scheme:http]
        	            	unexpected resource: map[host.name:runnervm8axu0.ldrbiddpfe1evawpvuyaexkoyh.gx.internal.cloudapp.net os.type:linux server.port:8889 service.instance.id:localhost:8889 service.name:internal-telemetry service_instance_id:2d9c63ac-c47a-447c-b2a0-972ca5e0353c service_name:otelcol service_version:v0.165.0-8-gcb3ed37f2 url.scheme:http]
    golden.go:188: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/testutils/golden.go:188
        	            				/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/telemetry/internal_telemetry_test.go:40
        	Error:      	Condition never satisfied
        	Test:       	TestInternalTelemetryPrometheusReaderNames</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44394053528/job/135129368711">Failure run</a></td></tr>
<tr><td>build-and-test</td><td>cross-compile (binaries-solaris_amd64) / compile</td><td>Run ./.github/actions/setup-environment</td><td>Flaky</td><td>1 / 15</td><td>7%</td><td><pre><code>cd ./internal/tools &amp;&amp; go install go.opentelemetry.io/collector/cmd/mdatagen
go: downloading go.opentelemetry.io/collector/cmd/mdatagen v0.165.0
go.opentelemetry.io/collector/cmd/mdatagen@v0.165.0: read &#34;https://proxy.golang.org/go.opentelemetry.io/collector/cmd/mdatagen/@v/v0.165.0.zip&#34;: stream error: stream ID 1; INTERNAL_ERROR; received from peer
make: *** [Makefile:249: install-tools] Error 1</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44442621017/job/135287170062">Failure run</a></td></tr>
<tr><td>build-and-test</td><td>lint</td><td>Lint</td><td>Flaky</td><td>1 / 15</td><td>7%</td><td><pre><code>2026-11-19T03:31:31.7585896Z Check License finished successfully
2026-11-19T03:31:31.7591313Z running misspell -error
2026-11-19T03:31:31.7800812Z golangci-lint run --allow-parallel-runners -j4
2026-11-19T03:38:38.0994729Z make[1]: *** wait: No child processes.  Stop.
2026-11-19T03:38:38.1522789Z ##[error]Process completed with exit code 143.
2026-11-19T03:38:38.1527143Z ##[error]The runner has received a shutdown signal. This can happen when the runner service is stopped, or a manually started runner is canceled.
2026-11-19T03:38:38.2524491Z Cleaning up orphan processes</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44451578915/job/135316603921">Failure run</a></td></tr>
<tr><td>build-and-test</td><td>test-ubuntu-24.04-race_detection:true-coverage:false / unit-test</td><td>Unit tests</td><td>Flaky</td><td>1 / 15</td><td>7%</td><td><pre><code>=== RUN   TestTunnelMultipleTargetsConcurrently
    receiver_test.go:489: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/internal/receiver/gnmireceiver/receiver_test.go:489
        	Error:      	Should be true
        	Test:       	TestTunnelMultipleTargetsConcurrently
        	Messages:   	expected metrics from both concurrently-registered targets, got map[routerA:true]</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44402506765/job/135149269969">Failure run</a></td></tr>
<tr><td>otlpintput</td><td>otlpinput-smoke-test-windows (1 matrix variant)</td><td>Restart Splunk UF</td><td>Flaky</td><td>1 / 8</td><td>13%</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44547254151/job/135644400997">Failure run</a></td></tr>
<tr><td>ta-v2</td><td>ta-v2-smoke-test-windows (3 matrix variants)</td><td>Check collector up and running</td><td>Flaky</td><td>5 / 11</td><td>45%</td><td><pre><code>2026-11-18T21:59:05.5722546Z Timeout: expvar endpoint was not available within 120 seconds
2026-11-18T21:59:05.6010139Z ##[error]Process completed with exit code 1.
2026-11-18T21:59:05.6273128Z ##[group]Run $splunkdLog = &#34;${Env:ProgramFiles}\SplunkUniversalForwarder\var\log\splunk\splunkd.log&#34;</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44401130260/job/135160177562">Failure run</a></td></tr>
<tr><td>vulnerability scans</td><td>anchore-image-scan (1 matrix variant)</td><td>Run anchore/scan-action@27805bf3b4e84b4a5c980df22ed233c00390a439; Upload result to GitHub Code Scanning</td><td>Flaky</td><td>1 / 14</td><td>7%</td><td><pre><code>2026-11-19T00:07:18.4098678Z   image: otelcol-fips:latest
2026-11-19T00:07:18.4098820Z   fail-build: true
2026-11-19T00:07:18.4098940Z   by-cve: false
2026-11-19T00:07:19.9924982Z ##[endgroup]
2026-11-19T00:07:20.4198224Z ##[group]/opt/hostedtoolcache/grype/0.118.0/x64/grype -v -o sarif --file /tmp/grype-Uzb9ns/output --fail-on high --only-fixed --add-cpes-if-none otelcol-fips:latest
2026-11-19T00:07:20.4402325Z [0000]  INFO grype version: 0.118.0
2026-11-19T00:07:20.8276439Z
2026-11-19T00:07:22.3822735Z [0001]  WARN error updating db error=unable to update vulnerability database: unable to download db: read tcp 10.1.0.157:44610-&gt;172.66.171.139:443: read: connection reset by peer
2026-11-19T00:07:22.3823166Z
2026-11-19T00:07:27.1550808Z
2026-11-19T00:07:27.1550962Z [0006] ERROR failed to load vulnerability db: database does not exist
2026-11-19T00:07:27.1551195Z
2026-11-19T00:07:27.2514930Z ##[endgroup]
2026-11-19T00:07:27.2532180Z ##[error]grype had a non-zero exit status when running
2026-11-19T00:07:27.2698743Z ##[group]Run github/codeql-action/upload-sarif@33d7b85da1a04171e211befe2c739303b9a4c066</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44424909850/job/135232149441">Failure run</a></td></tr>
<tr><td>vulnerability scans</td><td>anchore-win-image-scan (1 matrix variant)</td><td>Build docker image; Upload result to GitHub Code Scanning</td><td>Flaky</td><td>2 / 13</td><td>15%</td><td><pre><code>2026-11-18T21:43:37.2197967Z
2026-11-18T21:43:37.2512809Z failed to connect to the docker API at npipe:////./pipe/docker_engine; check if the path is correct and if the daemon is running: open //./pipe/docker_engine: The system cannot find the file specified.
2026-11-18T21:43:37.3252796Z ##[error]Process completed with exit code 1.
2026-11-18T21:43:37.3626452Z ##[group]Run github/codeql-action/upload-sarif@33d7b85da1a04171e211befe2c739303b9a4c066</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44402506320/job/135155189546">Failure run</a></td></tr>
<tr><td>vulnerability scans</td><td>anchore-win-image-scan (2 matrix variants)</td><td>Upload result to GitHub Code Scanning</td><td>Flaky</td><td>1 / 12</td><td>8%</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44424909850/job/135233204914">Failure run</a></td></tr>
<tr><td>vulnerability scans</td><td>binaries-windows_amd64 (binaries-windows_amd64) / compile</td><td>Run ./.github/actions/setup-environment</td><td>Flaky</td><td>1 / 13</td><td>8%</td><td><pre><code>2026-11-18T05:14:42.4199530Z go: downloading github.com/prometheus/procfs v0.22.0
2026-11-18T05:14:45.0300497Z ##[error]/home/runner/go/pkg/mod/dev.gaijin.team/go/exhaustruct/v4@v4.0.0/internal/pattern/list.go:7:2: dev.gaijin.team/go/golib@v0.6.0: Get &#34;https://proxy.golang.org/dev.gaijin.team/go/golib/@v/v0.6.0.mod&#34;: read tcp 10.1.0.25:32790-&gt;142.250.115.141:443: read: connection reset by peer
2026-11-18T05:14:45.0307261Z ##[error]/home/runner/go/pkg/mod/dev.gaijin.team/go/exhaustruct/v4@v4.0.0/internal/pattern/list.go:8:2: dev.gaijin.team/go/golib@v0.6.0: Get &#34;https://proxy.golang.org/dev.gaijin.team/go/golib/@v/v0.6.0.mod&#34;: read tcp 10.1.0.25:32790-&gt;142.250.115.141:443: read: connection reset by peer
2026-11-18T05:14:45.0309772Z ##[error]/home/runner/go/pkg/mod/github.com/nunnatsa/ginkgolinter@v0.21.2/internal/expression/expression.go:11:2: github.com/go-toolsmith/astcopy@v1.1.0: Get &#34;https://proxy.golang.org/github.com/go-toolsmith/astcopy/@v/v1.1.0.mod&#34;: read tcp 10.1.0.25:32790-&gt;142.250.115.141:443: read: connection reset by peer
2026-11-18T05:14:45.0312185Z ##[error]/home/runner/go/pkg/mod/github.com/go-critic/go-critic@v0.14.2/checkers/internal/lintutil/astflow.go:8:2: github.com/go-toolsmith/astequal@v1.2.0: Get &#34;https://proxy.golang.org/github.com/go-toolsmith/astequal/@v/v1.2.0.zip&#34;: read tcp 10.1.0.25:32790-&gt;142.250.115.141:443: read: connection reset by peer
2026-11-18T05:14:45.0345858Z make: *** [Makefile:246: install-tools] Error 1
2026-11-18T05:14:45.0353961Z ##[error]Process completed with exit code 2.
2026-11-18T05:14:45.0359332Z ##[end-action id=__self.__run;outcome=failure;conclusion=failure;duration_ms=47269]
2026-11-18T05:14:45.0361800Z ##[start-action display=Run free-disk-space.sh;id=__self.__run_2]</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44242773977/job/134594385373">Failure run</a></td></tr>
</tbody>
</table>

## Broken workflow and job failures

<table>
<thead>
<tr><th>Workflow</th><th>Job</th><th>Step</th><th>Classification</th><th>Failures / observations</th><th>Failure rate</th><th>Recent failure</th></tr>
</thead>
<tbody>
<tr><td>auto-instrumentation</td><td>test-package (1 matrix variant)</td><td>Re-run failed tests</td><td>Broken</td><td>1 / 1</td><td>100%</td><td><pre><code>=== RUN   TestTomcatInstrumentation/oraclelinux-8/amd64
    instrumentation_test.go:119: 
        	Error Trace:	/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/instrumentation/instrumentation_test.go:373
        	            				/home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/instrumentation/instrumentation_test.go:119
        	Error:      	Should be true
        	Test:       	TestTomcatInstrumentation/oraclelinux-8/amd64
        	Messages:   	timed out waiting for &#34;com\\.splunk\\.sourcetype&#34;: Str\(otel\.profiling\)</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44594002584/job/135814705041">Failure run</a></td></tr>
<tr><td>vulnerability scans</td><td></td><td></td><td>Broken</td><td>15 / 15</td><td>100%</td><td><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44594002400">Failure run</a></td></tr>
<tr><td>vulnerability scans</td><td>anchore-fs-scan</td><td>Run anchore/scan-action@27805bf3b4e84b4a5c980df22ed233c00390a439</td><td>Broken</td><td>17 / 17</td><td>100%</td><td><pre><code>2026-11-19T14:45:46.8548378Z   path: .
2026-11-19T14:45:46.8549302Z   fail-build: true
2026-11-19T14:45:46.8550276Z   by-cve: false
2026-11-19T14:45:48.6581844Z ##[endgroup]
2026-11-19T14:45:49.0204780Z ##[group]/opt/hostedtoolcache/grype/0.118.0/x64/grype -v -o table --file /tmp/grype-nXcfx2/output --fail-on high --only-fixed --add-cpes-if-none dir:.
2026-11-19T14:45:49.0521770Z [0000]  INFO grype version: 0.118.0
2026-11-19T14:47:38.6277849Z
2026-11-19T14:47:38.6322879Z [0109] ERROR discovered vulnerabilities at or above the severity threshold
2026-11-19T14:47:38.6323335Z
2026-11-19T14:47:38.6420046Z
2026-11-19T14:47:38.6440618Z ##[error]Failed minimum severity level. Found vulnerabilities with level &#39;high&#39; or higher</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44547248625/job/135641007304">Failure run</a></td></tr>
<tr><td>win-package-test</td><td>choco-test (1 matrix variant)</td><td>Chocolatey install test</td><td>Broken</td><td>1 / 1</td><td>100%</td><td><pre><code>2026-11-18T21:49:10.3241491Z
2026-11-18T21:49:10.3432019Z ##[error]Process completed with exit code 1.
2026-11-18T21:49:10.3590104Z ##[group]Run Get-WinEvent -ProviderName splunk-otel-collector | Sort-Object -Property TimeCreated | Select-Object -Property Message | Format-List</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44402506936/job/135158078707">Failure run</a></td></tr>
<tr><td>win-package-test</td><td>choco-test (1 matrix variant)</td><td>Chocolatey upgrade test; splunk-otel-collector logs</td><td>Broken</td><td>1 / 1</td><td>100%</td><td><pre><code>2026-11-18T15:04:50.9393398Z By installing, you accept licenses for the packages.
2026-11-18T15:05:06.4403341Z Failed to fetch results from V2 feed at &#39;https://community.chocolatey.org/api/v2/Packages(Id=&#39;splunk-otel-collector&#39;,Version=&#39;0.165.0&#39;)&#39; with following message : Response status code does not indicate success: 504 (Gateway Timeout).
2026-11-18T15:05:06.4954076Z Need to add specific handling for exception type NuGetResolverInputException
2026-11-18T15:05:06.9576162Z
2026-11-18T15:05:07.0119239Z ##[error]Process completed with exit code 1.
2026-11-18T15:05:07.0222955Z ##[group]Run Get-WinEvent -ProviderName splunk-otel-collector | Sort-Object -Property TimeCreated | Select-Object -Property Message | Format-List
2026-11-18T15:05:07.8391811Z [31;1m[0m[36;1m[36;1m[0m[36;1m[0m[36;1m[31;1m[31;1m[36;1m     | [31;1mThere is not an event provider on the localhost computer that matches &#34;splunk-otel-collector&#34;.[0m
2026-11-18T15:05:07.8752693Z ##[error]Process completed with exit code 1.
2026-11-18T15:05:07.8899850Z Post job cleanup.</code></pre><br><a href="https://github.com/signalfx/splunk-otel-collector/actions/runs/44328009083/job/134895747111">Failure run</a></td></tr>
</tbody>
</table>

