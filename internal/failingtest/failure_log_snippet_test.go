// Copyright Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package failingtest

import (
	"strings"
	"testing"
)

const failureLogSnippetWindowsDifferenceExample = `2026-11-13T19:02:06.3696503Z CoreCompile:
2026-11-13T19:02:06.3711627Z csc.exe /noconfig /nowarn:1701,1702 /reference:many-reference-files
2026-11-13T19:07:25.5360469Z diff --git a/tests/zeroconfig/windows/testdata/expected/aspnetfx.yaml b/tests/zeroconfig/windows/testdata/expected/aspnetfx.yaml
2026-11-13T19:07:25.5360470Z index f10c30e..2fde2d1 100644
2026-11-13T19:07:25.5360471Z --- a/tests/zeroconfig/windows/testdata/expected/aspnetfx.yaml
2026-11-13T19:07:25.5360472Z +++ b/tests/zeroconfig/windows/testdata/expected/aspnetfx.yaml
2026-11-13T19:07:25.5360473Z @@ -3,16 +3,16 @@ resourceSpans:
2026-11-13T19:07:25.5360474Z -                stringValue: 1.16.0
2026-11-13T19:07:25.5360475Z +                stringValue: 1.17.0
2026-11-13T19:07:25.5360476Z One or more files were changed, this is not expected.
2026-11-13T19:07:25.5360477Z Error: Process completed with exit code 1.`

const failureLogSnippetGoModuleProxyStreamFailureExample = `2026-11-09T17:00:27.7851648Z go: downloading golang.org/x/telemetry v0.0.0-20261010193011-f71dc92da45e
2026-11-09T17:00:29.6883336Z cd ./internal/tools && go install go.opentelemetry.io/build-tools/chloggen
2026-11-09T17:00:29.6950089Z go: downloading go.opentelemetry.io/build-tools/chloggen v0.26.0
2026-11-09T17:02:29.9072459Z go.opentelemetry.io/build-tools/chloggen@v0.26.0: read "https://proxy.golang.org/go.opentelemetry.io/build-tools/chloggen/@v/v0.26.0.zip": stream error: stream ID 1; INTERNAL_ERROR; received from peer
2026-11-09T17:02:29.9079994Z make: *** [Makefile:254: install-tools] Error 1
2026-11-09T17:02:29.9099665Z ##[error]Process completed with exit code 2.`

const failureLogSnippetGoModuleSetupFailureExample = `2026-11-17T16:20:00.0000000Z go: downloading go.opentelemetry.io/collector/consumer v1.70.0
2026-11-17T16:20:01.0000000Z go: downloading k8s.io/apiextensions-apiserver v0.37.2
2026-11-17T16:20:02.0000000Z Error: ../../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/xk8stest@v0.164.1/k8s_collector.go:20:2: k8s.io/api@v0.37.2: read "https://proxy.golang.org/k8s.io/api/@v/v0.37.2.zip": stream error: stream ID 21; INTERNAL_ERROR; received from peer
2026-11-17T16:20:03.0000000Z FAIL github.com/signalfx/splunk-otel-collector-chart/functional_tests/functional [setup failed]
2026-11-17T16:20:03.0000001Z FAIL
2026-11-17T16:20:03.0000002Z make: *** [Makefile:129: functionaltest] Error 1
2026-11-17T16:20:03.0000003Z Error: Process completed with exit code 2.`

const failureLogSnippetGoModuleSetupFailureMatchingDownload = "2026-11-17T16:20:00.0000000Z go: downloading k8s.io/api v0.37.2"
const failureLogSnippetGoModuleSetupFailureFirstDownload = "2026-11-17T16:20:00.0000000Z go: downloading go.opentelemetry.io/collector/consumer v1.70.0"
const failureLogSnippetGoModuleSetupFailureNoiseLine = "2026-11-17T16:20:00.0000000Z go: downloading example.com/unrelated/module v1.0.0"

const failureLogSnippetGoModuleSetupFailureExpected = `Error: ../../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/xk8stest@v0.164.1/k8s_collector.go:20:2: k8s.io/api@v0.37.2: read "https://proxy.golang.org/k8s.io/api/@v/v0.37.2.zip": stream error: stream ID 21; INTERNAL_ERROR; received from peer
FAIL github.com/signalfx/splunk-otel-collector-chart/functional_tests/functional [setup failed]
FAIL
make: *** [Makefile:129: functionaltest] Error 1`

const failureLogSnippetGoModuleSetupFailureWithDownloadExpected = `go: downloading k8s.io/api v0.37.2
Error: ../../../../go/pkg/mod/github.com/open-telemetry/opentelemetry-collector-contrib/pkg/xk8stest@v0.164.1/k8s_collector.go:20:2: k8s.io/api@v0.37.2: read "https://proxy.golang.org/k8s.io/api/@v/v0.37.2.zip": stream error: stream ID 21; INTERNAL_ERROR; received from peer
FAIL github.com/signalfx/splunk-otel-collector-chart/functional_tests/functional [setup failed]
FAIL
make: *** [Makefile:129: functionaltest] Error 1`

const failureLogSnippetGoTestJSONFailureExample = `{"Action":"run","Package":"example.com/project/pkg","Test":"TestFailed"}
{"Action":"output","Package":"example.com/project/pkg","Test":"TestFailed","Output":"    assertion failed\n"}
{"Action":"fail","Package":"example.com/project/pkg","Test":"TestFailed","Elapsed":0.01}
{"Action":"fail","Package":"example.com/project/pkg","Elapsed":0.01}`

const failureLogSnippetWindowsDifferenceExpected = `diff --git a/tests/zeroconfig/windows/testdata/expected/aspnetfx.yaml b/tests/zeroconfig/windows/testdata/expected/aspnetfx.yaml
index f10c30e..2fde2d1 100644
--- a/tests/zeroconfig/windows/testdata/expected/aspnetfx.yaml
+++ b/tests/zeroconfig/windows/testdata/expected/aspnetfx.yaml
@@ -3,16 +3,16 @@ resourceSpans:
-                stringValue: 1.16.0
+                stringValue: 1.17.0
One or more files were changed, this is not expected.
Error: Process completed with exit code 1.`

const failureLogSnippetTestNameWindowsExpectedDiff = "prioritizes the expected trace diff over unrelated compiler output"
const failureLogSnippetGoModuleProxyStreamFailureExpected = `cd ./internal/tools && go install go.opentelemetry.io/build-tools/chloggen
go: downloading go.opentelemetry.io/build-tools/chloggen v0.26.0
go.opentelemetry.io/build-tools/chloggen@v0.26.0: read "https://proxy.golang.org/go.opentelemetry.io/build-tools/chloggen/@v/v0.26.0.zip": stream error: stream ID 1; INTERNAL_ERROR; received from peer
make: *** [Makefile:254: install-tools] Error 1`
const failureLogSnippetTestNameGoModuleProxyStreamFailure = "selects the failed Go module download and make error rather than preceding downloads"
const failureLogSnippetTestNameGoModuleSetupFailure = "keeps setup failure details when the matching download line is missing"
const failureLogSnippetTestNameGoModuleSetupFailureWithDownloads = "keeps the proxy error when unrelated downloads exceed the excerpt limit"
const failureLogSnippetCmdLineFailureExpected = `=== RUN TestRunFromCmdLine/full_linux
Error: failed to get config: cannot unmarshal the configuration
'splunkhecexporter.Config' has invalid keys: max_idle_conns
main.go:142: application run finished with error`
const failureLogSnippetTestNameCmdLineFailure = "uses the failing command-line test instead of earlier passing test output"
const failureLogSnippetPackageUpgradeTimeoutExpected = `=== RUN   TestPackageUpgradeFromLibsplunk/opensuse-15/amd64
2026/11/13 18:33:10 🐳 Building image 5d4a6598-1209-444b-be9d-d80a7cedca8e:c968bc78-fb9c-4140-aaea-8be3704b44d0
panic: test timed out after 15m0s
	running tests:
		TestPackageUpgradeFromLibsplunk (5m0s)
		TestPackageUpgradeFromLibsplunk/opensuse-15/amd64 (5m0s)`
const failureLogSnippetTestNamePackageUpgradeTimeout = "keeps the latest timeout summary and drops the goroutine dump"
const failureLogSnippetTestNameGoTestJSONFailure = "uses the structured Go test failure output"
const failureLogSnippetGoTestJSONFailureExpected = "example.com/project/pkg.TestFailed\n    assertion failed\n"
const failureLogSnippetNestedGoTestFailureName = "selects the nested failing test rather than aggregate parent summaries"
const failureLogSnippetNestedGoTestFailureMarker = "=== RUN   Test_Functions/logs_and_metrics_attributes_verification/verify_cluster_receiver_metrics_attributes"

func TestExtractFailureLogSnippetRequestedFailureExamples(t *testing.T) {
	testCases := []struct {
		log  string
		name string
		want string
	}{
		{
			log:  goTestTextLogParserCmdLineFailureExample,
			name: failureLogSnippetTestNameCmdLineFailure,
			want: failureLogSnippetCmdLineFailureExpected,
		},
		{
			log:  failureLogSnippetWindowsDifferenceExample,
			name: failureLogSnippetTestNameWindowsExpectedDiff,
			want: failureLogSnippetWindowsDifferenceExpected,
		},
		{
			log:  failureLogSnippetGoModuleProxyStreamFailureExample,
			name: failureLogSnippetTestNameGoModuleProxyStreamFailure,
			want: failureLogSnippetGoModuleProxyStreamFailureExpected,
		},
		{
			log:  failureLogSnippetGoModuleSetupFailureExample,
			name: failureLogSnippetTestNameGoModuleSetupFailure,
			want: failureLogSnippetGoModuleSetupFailureExpected,
		},
		{
			log:  failureLogSnippetGoModuleProxySetupFailureWithLargeDownloadList(),
			name: failureLogSnippetTestNameGoModuleSetupFailureWithDownloads,
			want: failureLogSnippetGoModuleSetupFailureWithDownloadExpected,
		},
		{
			log:  goTestTextLogParserPackageUpgradeTimeoutExample,
			name: failureLogSnippetTestNamePackageUpgradeTimeout,
			want: failureLogSnippetPackageUpgradeTimeoutExpected,
		},
		{
			log:  failureLogSnippetGoTestJSONFailureExample,
			name: failureLogSnippetTestNameGoTestJSONFailure,
			want: failureLogSnippetGoTestJSONFailureExpected,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := extractFailureLogSnippet(testCase.log); got != testCase.want {
				t.Errorf("extractFailureLogSnippet() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func failureLogSnippetGoModuleProxySetupFailureWithLargeDownloadList() string {
	noiseLines := strings.TrimSuffix(strings.Repeat(failureLogSnippetGoModuleSetupFailureNoiseLine+logLineSeparator, 100), logLineSeparator)
	replacement := strings.Join([]string{
		failureLogSnippetGoModuleSetupFailureMatchingDownload,
		noiseLines,
		failureLogSnippetGoModuleSetupFailureFirstDownload,
	}, logLineSeparator)
	return strings.Replace(failureLogSnippetGoModuleSetupFailureExample, failureLogSnippetGoModuleSetupFailureFirstDownload, replacement, 1)
}

func TestExtractFailureLogSnippetSelectsNestedGoTestFailure(t *testing.T) {
	testCases := []struct {
		log  string
		name string
		want string
	}{
		{
			log:  goTestTextLogParserNestedFailureExample,
			name: failureLogSnippetNestedGoTestFailureName,
			want: failureLogSnippetNestedGoTestFailureMarker,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := extractFailureLogSnippet(testCase.log)
			if !strings.HasPrefix(got, testCase.want) {
				t.Errorf("extractFailureLogSnippet() = %q, want prefix %q", got, testCase.want)
			}
			if !strings.Contains(got, goTestTextLogParserNestedFailureMarker) {
				t.Errorf("extractFailureLogSnippet() = %q, want failed assertion details", got)
			}
		})
	}
}
