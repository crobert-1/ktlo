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

const goTestTextLogParserCmdLineFailureExample = `2026-11-13T19:01:40.9872368Z time="2026-11-13T19:01:40Z" level=error msg="Could not receive ECS Task Metadata" error="failed to read ECS container data"
2026-11-13T19:01:40.9875844Z === RUN   TestDimensionToUpdate/invalid
2026-11-13T19:01:41.1278140Z === RUN   TestCommonDimensions
2026-11-13T19:01:41.1309075Z time="2026-11-13T19:01:41Z" level=info msg="Failed to parse filesystem mode"
2026-11-13T19:01:41.1318759Z --- PASS: TestCommonDimensions (0.00s)
2026-11-13T19:05:02.0000000Z === RUN   TestRunFromCmdLine/full_linux
2026-11-13T19:05:02.0000001Z Error: failed to get config: cannot unmarshal the configuration
2026-11-13T19:05:02.0000002Z 'splunkhecexporter.Config' has invalid keys: max_idle_conns
2026-11-13T19:05:02.0000003Z main.go:142: application run finished with error
2026-11-13T19:05:02.0000004Z FAIL github.com/signalfx/splunk-otel-collector/cmd/otelcol 15.439s`

const goTestTextLogParserWindowsIISFailureExample = `=== RUN   TestWindowsIISInstrumentation
windows_iis_test.go:135:
    Error Trace: windows_iis_test.go:168
    Error: Received unexpected error:
        no matching traces found, 2 collected
    Test: TestWindowsIISInstrumentation
    Messages: Failed to receive expected traces
--- FAIL: TestWindowsIISInstrumentation (0.00s)
FAIL github.com/signalfx/splunk-otel-collector/tests/zeroconfig/windows 15.439s`

const goTestTextLogParserPackageUpgradeTimeoutExample = `=== RUN   TestPackageUpgradeFromLibsplunk
=== RUN   TestPackageUpgradeFromLibsplunk/opensuse-15/amd64
2026/11/13 18:17:33 🐳 Building image old-image-id
panic: test timed out after 15m0s
	running tests:
		TestPackageUpgradeFromLibsplunk (5m0s)
		TestPackageUpgradeFromLibsplunk/opensuse-15/amd64 (5m0s)

goroutine 1 [chan receive]:
FAIL github.com/signalfx/splunk-otel-collector/tests/instrumentation 900.043s
=== RUN   TestPackageUpgradeFromLibsplunk
=== RUN   TestPackageUpgradeFromLibsplunk/opensuse-15/amd64
2026/11/13 18:33:10 🐳 Building image 5d4a6598-1209-444b-be9d-d80a7cedca8e:c968bc78-fb9c-4140-aaea-8be3704b44d0
panic: test timed out after 15m0s
	running tests:
		TestPackageUpgradeFromLibsplunk (5m0s)
		TestPackageUpgradeFromLibsplunk/opensuse-15/amd64 (5m0s)

goroutine 1 [chan receive]:
FAIL github.com/signalfx/splunk-otel-collector/tests/instrumentation 900.024s`

const goTestTextLogParserNestedFailureExample = `=== RUN   Test_Functions
=== RUN   Test_Functions/logs_and_metrics_attributes_verification
=== RUN   Test_Functions/logs_and_metrics_attributes_verification/verify_cluster_receiver_metrics_attributes
configuration_switching_test.go:348: ===> >>>> Metrics: 0
Error Trace: common.go:115
Error: Condition never satisfied
Test: Test_Functions/logs_and_metrics_attributes_verification/verify_cluster_receiver_metrics_attributes
Messages: failed to receive 5 entries, received 0 metrics in 3.000000 minutes
configuration_switching_test.go:65: Cleaning up cluster
configuration_switching_test.go:71: Running teardown
chart.go:213: Secret not found, nothing to delete
--- FAIL: Test_Functions/logs_and_metrics_attributes_verification/verify_cluster_receiver_metrics_attributes (182.05s)
--- PASS: Test_Functions/logs_and_metrics_attributes_verification/verify_cluster_receiver_attributes (27.10s)
--- PASS: Test_Functions/logs_and_metrics_attributes_verification/verify_agent_logs_attributes (43.30s)
--- PASS: Test_Functions/logs_and_metrics_attributes_verification/verify_metrics_attributes (93.06s)
--- FAIL: Test_Functions/logs_and_metrics_attributes_verification (345.52s)
--- PASS: Test_Functions/agent_logs_and_metrics_enabled_or_disabled (288.21s)
--- PASS: Test_Functions/logs_and_metrics_index_switch (189.28s)
--- PASS: Test_Functions/cluster_receiver_enabled_or_disabled (49.05s)
--- FAIL: Test_Functions (872.08s)
FAIL github.com/signalfx/splunk-otel-collector-chart/functional_tests/configuration_switching 872.105s`

const goTestTextLogParserIndependentParentFailureExample = `=== RUN   TestParent
Error: parent setup failed independently
=== RUN   TestParent/child
Error: child assertion failed
--- FAIL: TestParent/child (1.00s)
--- FAIL: TestParent (2.00s)
FAIL example.com/project/pkg 2.00s`

const goTestTextLogParserWantCmdLineFailureID = "github.com/signalfx/splunk-otel-collector/cmd/otelcol.TestRunFromCmdLine/full_linux"
const goTestTextLogParserWantWindowsIISFailureID = "github.com/signalfx/splunk-otel-collector/tests/zeroconfig/windows.TestWindowsIISInstrumentation"
const goTestTextLogParserWantPackageUpgradeTimeoutID = "github.com/signalfx/splunk-otel-collector/tests/instrumentation.TestPackageUpgradeFromLibsplunk/opensuse-15/amd64"
const goTestTextLogParserWantNestedFailureID = "github.com/signalfx/splunk-otel-collector-chart/functional_tests/configuration_switching.Test_Functions/logs_and_metrics_attributes_verification/verify_cluster_receiver_metrics_attributes"
const goTestTextLogParserNestedFailureMarker = "Condition never satisfied"
const goTestTextLogParserIndependentParentFailureID = "example.com/project/pkg.TestParent"
const goTestTextLogParserIndependentChildFailureID = "example.com/project/pkg.TestParent/child"
const goTestTextLogParserInvalidConfigKey = "invalid keys: max_idle_conns"
const goTestTextLogParserWindowsIISFailureText = "Failed to receive expected traces"
const goTestTextLogParserPackageUpgradeTimeoutText = "panic: test timed out after 15m0s"
const goTestTextLogParserTestNameCmdLineFailure = "package-level failure is attributed to the active test"
const goTestTextLogParserTestNamePackageUpgradeTimeout = "timeout panic is attributed to the active subtest"
const goTestTextLogParserTestNameWindowsIISFailure = "explicit test failure retains its test identity"
const goTestTextLogParserTestNameNestedFailure = "only the failed leaf is reported when parent summaries aggregate child failures"
const goTestTextLogParserTestNameIndependentParentFailure = "preserves a parent test with distinct failure output"
const goTestTextLogParserIncompleteRunOnly = "=== RUN TestStartedWithoutResult"
const goTestTextLogParserTestNameNoResult = "does not claim an incomplete run marker as a result"

func TestGoTestTextLogParserParsesRequestedFailureExamples(t *testing.T) {
	testCases := []struct {
		failureText string
		log         string
		name        string
		testID      string
		wantCount   int
	}{
		{
			failureText: goTestTextLogParserInvalidConfigKey,
			log:         goTestTextLogParserCmdLineFailureExample,
			name:        goTestTextLogParserTestNameCmdLineFailure,
			testID:      goTestTextLogParserWantCmdLineFailureID,
			wantCount:   1,
		},
		{
			failureText: goTestTextLogParserWindowsIISFailureText,
			log:         goTestTextLogParserWindowsIISFailureExample,
			name:        goTestTextLogParserTestNameWindowsIISFailure,
			testID:      goTestTextLogParserWantWindowsIISFailureID,
			wantCount:   1,
		},
		{
			failureText: goTestTextLogParserPackageUpgradeTimeoutText,
			log:         goTestTextLogParserPackageUpgradeTimeoutExample,
			name:        goTestTextLogParserTestNamePackageUpgradeTimeout,
			testID:      goTestTextLogParserWantPackageUpgradeTimeoutID,
			wantCount:   2,
		},
	}
	parser := GoTestTextLogParser{}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if !parser.CanParse(testCase.log) {
				t.Fatal("CanParse() = false, want true")
			}
			parsedTests, err := parser.Parse(testCase.log)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			matchedTestCount := 0
			for _, parsedTest := range parsedTests {
				if parsedTest.ID != testCase.testID {
					continue
				}
				matchedTestCount++
				if parsedTest.Kind != testresults.TestKindTest {
					t.Errorf("test kind = %q, want %q", parsedTest.Kind, testresults.TestKindTest)
				}
				if parsedTest.Outcome != testresults.OutcomeFailed {
					t.Errorf("test outcome = %q, want %q", parsedTest.Outcome, testresults.OutcomeFailed)
				}
				if !strings.Contains(parsedTest.FailureMessage, testCase.failureText) {
					t.Errorf("failure message %q does not contain %q", parsedTest.FailureMessage, testCase.failureText)
				}
			}
			if matchedTestCount != testCase.wantCount {
				t.Fatalf("Parse() returned %d observations for %q, want %d; got %#v", matchedTestCount, testCase.testID, testCase.wantCount, parsedTests)
			}
		})
	}
}

func TestGoTestTextLogParserRequiresResultEvidence(t *testing.T) {
	t.Run(goTestTextLogParserTestNameNoResult, func(t *testing.T) {
		if (GoTestTextLogParser{}).CanParse(goTestTextLogParserIncompleteRunOnly) {
			t.Fatal("CanParse() = true, want false without a result line")
		}
	})
}

func TestGoTestTextLogParserOmitsAggregateParentFailures(t *testing.T) {
	testCases := []struct {
		log             string
		name            string
		wantFailureIDs  []string
		wantFailureText string
	}{
		{
			log:             goTestTextLogParserNestedFailureExample,
			name:            goTestTextLogParserTestNameNestedFailure,
			wantFailureIDs:  []string{goTestTextLogParserWantNestedFailureID},
			wantFailureText: goTestTextLogParserNestedFailureMarker,
		},
		{
			log:  goTestTextLogParserIndependentParentFailureExample,
			name: goTestTextLogParserTestNameIndependentParentFailure,
			wantFailureIDs: []string{
				goTestTextLogParserIndependentChildFailureID,
				goTestTextLogParserIndependentParentFailureID,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			parsedTests, err := (GoTestTextLogParser{}).Parse(testCase.log)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			failedTestIDs := make([]string, 0)
			for _, parsedTest := range parsedTests {
				if parsedTest.Kind != testresults.TestKindTest || parsedTest.Outcome != testresults.OutcomeFailed {
					continue
				}
				failedTestIDs = append(failedTestIDs, parsedTest.ID)
			}
			if len(failedTestIDs) != len(testCase.wantFailureIDs) {
				t.Fatalf("failed test IDs = %v, want %v", failedTestIDs, testCase.wantFailureIDs)
			}
			for failureIndex, wantFailureID := range testCase.wantFailureIDs {
				if failedTestIDs[failureIndex] != wantFailureID {
					t.Errorf("failed test IDs = %v, want %v", failedTestIDs, testCase.wantFailureIDs)
				}
			}
			for _, parsedTest := range parsedTests {
				if parsedTest.ID != goTestTextLogParserWantNestedFailureID {
					continue
				}
				if !strings.Contains(parsedTest.FailureMessage, testCase.wantFailureText) {
					t.Errorf("failure message = %q, want it to include %q", parsedTest.FailureMessage, testCase.wantFailureText)
				}
			}
		})
	}
}
