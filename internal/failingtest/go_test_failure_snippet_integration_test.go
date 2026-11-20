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

	"ktlo/internal/testresults"
)

const goTestFailureSnippetIntegrationLog = `2026-11-18T05:36:38.2535468Z === RUN   TestDefaultConfigsDoNotEnableExpvar
2026-11-18T05:36:38.2536050Z === RUN   TestDefaultConfigsDoNotEnableExpvar/gateway
2026-11-18T05:37:12.7551749Z     default_config_test.go:643:
2026-11-18T05:37:12.7553398Z         Error Trace: /home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/general/default_config_test.go:644
2026-11-18T05:37:12.7556095Z                      /opt/hostedtoolcache/go/1.26.11/x64/src/runtime/asm_amd64.s:1771
2026-11-18T05:37:12.7556960Z         Error:       Not equal:
2026-11-18T05:37:12.7557816Z                      expected: 200
2026-11-18T05:37:12.7558624Z                      actual  : 404
2026-11-18T05:37:12.7559498Z     default_config_test.go:643:
2026-11-18T05:37:12.7561252Z         Error Trace: /home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/general/default_config_test.go:644
2026-11-18T05:37:12.7564967Z         Error:       Condition never satisfied
2026-11-18T05:37:12.7566188Z         Test:        TestDefaultConfigsDoNotEnableExpvar/gateway
2026-11-18T05:37:12.7567437Z Logs:
2026-11-18T05:37:12.7567438Z {{info gateway github.com/open-telemetry/opamp-go/client/internal/httpsender.go:233 undefined } []}
2026-11-18T05:37:12.9746124Z --- FAIL: TestDefaultConfigsDoNotEnableExpvar/gateway (34.66s)
2026-11-18T05:37:12.9891473Z === RUN   TestDefaultConfigsDoNotEnableExpvar/agent
2026-11-18T05:37:47.7551749Z     default_config_test.go:643:
2026-11-18T05:37:47.7553398Z         Error Trace: /home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/general/default_config_test.go:644
2026-11-18T05:37:47.7556095Z                      /opt/hostedtoolcache/go/1.26.11/x64/src/runtime/asm_amd64.s:1771
2026-11-18T05:37:47.7556960Z         Error:       Not equal:
2026-11-18T05:37:47.7557816Z                      expected: 200
2026-11-18T05:37:47.7558624Z                      actual  : 404
2026-11-18T05:37:47.7559498Z     default_config_test.go:643:
2026-11-18T05:37:47.7561252Z         Error Trace: /home/runner/work/splunk-otel-collector/splunk-otel-collector/tests/general/default_config_test.go:644
2026-11-18T05:37:47.7564967Z         Error:       Condition never satisfied
2026-11-18T05:37:47.7566188Z         Test:        TestDefaultConfigsDoNotEnableExpvar/agent
2026-11-18T05:37:47.7567437Z 2026/11/18 05:37:47 🐳 Stopping container: 706481d8871e
2026-11-18T05:37:47.8836099Z Logs:
2026-11-18T05:37:47.8837550Z {{info 2026-11-18 05:37:17.755310862 +0000 UTC m=+66.834553837  github.com/open-telemetry/opamp-go@v0.25.0/client/internal/httpsender.go:233 undefined } []}
2026-11-18T05:37:47.8849745Z {{info 2026-11-18 05:37:17.755334197 +0000 UTC m=+66.834577176  2026-11-18T05:37:16.660Z info opampextension@v0.165.1/tracer.go:47 Internal trace telemetry disabled undefined } []}
2026-11-18T05:37:47.9746858Z --- FAIL: TestDefaultConfigsDoNotEnableExpvar/agent (34.97s)
2026-11-18T05:37:47.9747659Z --- FAIL: TestDefaultConfigsDoNotEnableExpvar (69.63s)
2026-11-18T05:39:30.9328633Z FAIL github.com/signalfx/splunk-otel-collector/tests/general 200.014s`

const goTestFailureSnippetIntegrationAgent = "agent"
const goTestFailureSnippetIntegrationExpectedStatus = "actual  : 404"
const goTestFailureSnippetIntegrationExpectedTrace = "tests/general/default_config_test.go:644"
const goTestFailureSnippetIntegrationExpectedValue = "expected: 200"
const goTestFailureSnippetIntegrationGateway = "gateway"
const goTestFailureSnippetIntegrationGatewayRunLine = "=== RUN   TestDefaultConfigsDoNotEnableExpvar/gateway"
const goTestFailureSnippetIntegrationGoTestName = "TestDefaultConfigsDoNotEnableExpvar"
const goTestFailureSnippetIntegrationGoTestPackage = "github.com/signalfx/splunk-otel-collector/tests/general"
const goTestFailureSnippetIntegrationLogMarker = "{{info"
const goTestFailureSnippetIntegrationOpAMPMarker = "opamp-go"
const goTestFailureSnippetIntegrationRunLine = "=== RUN   TestDefaultConfigsDoNotEnableExpvar/agent"
const goTestFailureSnippetIntegrationSubtestSeparator = "/"
const goTestFailureSnippetIntegrationTestCaseName = "selects Go text parsing and retains assertion failures before shutdown logs"
const goTestFailureSnippetIntegrationTestMarker = "Condition never satisfied"

func TestGoTestFailureSnippetIntegration(t *testing.T) {
	t.Run(goTestFailureSnippetIntegrationTestCaseName, func(t *testing.T) {
		selectedParser, err := selectTestLogParser(goTestFailureSnippetIntegrationLog, defaultTestLogParsers())
		if err != nil {
			t.Fatalf("selectTestLogParser() error = %v", err)
		}
		if _, isGoTextParser := selectedParser.(GoTestTextLogParser); !isGoTextParser {
			t.Fatalf("selected parser = %T, want GoTestTextLogParser", selectedParser)
		}
		parsedTests, err := selectedParser.Parse(goTestFailureSnippetIntegrationLog)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}
		testCases := []struct {
			id   string
			name string
		}{
			{
				name: goTestFailureSnippetIntegrationGateway,
				id:   goTestFailureSnippetIntegrationGoTestPackage + packageTestIDSeparator + goTestFailureSnippetIntegrationGoTestName + goTestFailureSnippetIntegrationSubtestSeparator + goTestFailureSnippetIntegrationGateway,
			},
			{
				name: goTestFailureSnippetIntegrationAgent,
				id:   goTestFailureSnippetIntegrationGoTestPackage + packageTestIDSeparator + goTestFailureSnippetIntegrationGoTestName + goTestFailureSnippetIntegrationSubtestSeparator + goTestFailureSnippetIntegrationAgent,
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				parsedTest := failureTestWithID(parsedTests, testCase.id)
				if parsedTest == nil {
					t.Fatalf("Parse() did not return failure %q", testCase.id)
				}
				if parsedTest.Kind != testresults.TestKindTest || parsedTest.Outcome != testresults.OutcomeFailed {
					t.Errorf("test kind/outcome = %q/%q, want %q/%q", parsedTest.Kind, parsedTest.Outcome, testresults.TestKindTest, testresults.OutcomeFailed)
				}
				for _, failureDetail := range []string{
					goTestFailureSnippetIntegrationExpectedTrace,
					goTestFailureSnippetIntegrationExpectedStatus,
					goTestFailureSnippetIntegrationExpectedValue,
					goTestFailureSnippetIntegrationTestMarker,
				} {
					if !strings.Contains(parsedTest.FailureMessage, failureDetail) {
						t.Errorf("failure message %q does not contain %q", parsedTest.FailureMessage, failureDetail)
					}
				}
				if strings.Contains(parsedTest.FailureMessage, goTestFailureSnippetIntegrationLogMarker) {
					t.Errorf("failure message includes the collector log dump: %q", parsedTest.FailureMessage)
				}
			})
		}
		failureSnippet := extractFailureLogSnippet(goTestFailureSnippetIntegrationLog)
		for _, failureDetail := range []string{
			goTestFailureSnippetIntegrationRunLine,
			goTestFailureSnippetIntegrationExpectedTrace,
			goTestFailureSnippetIntegrationExpectedStatus,
			goTestFailureSnippetIntegrationExpectedValue,
			goTestFailureSnippetIntegrationTestMarker,
		} {
			if !strings.Contains(failureSnippet, failureDetail) {
				t.Errorf("failure snippet %q does not contain %q", failureSnippet, failureDetail)
			}
		}
		if strings.Contains(failureSnippet, goTestFailureSnippetIntegrationLogMarker) || strings.Contains(failureSnippet, goTestFailureSnippetIntegrationOpAMPMarker) {
			t.Errorf("failure snippet includes trailing collector logs: %q", failureSnippet)
		}
		if strings.Contains(failureSnippet, goTestFailureSnippetIntegrationGatewayRunLine) {
			t.Errorf("agent failure snippet includes gateway output: %q", failureSnippet)
		}
	})
}

func failureTestWithID(tests []testresults.Test, testID string) *testresults.Test {
	for testIndex := range tests {
		if tests[testIndex].ID == testID {
			return &tests[testIndex]
		}
	}
	return nil
}
