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
	"testing"

	"ktlo/internal/testresults"
)

const ansiEscape = "\x1b"

const ansibleMoleculeLogParserFailureExample = `2026-11-17T17:33:25.7131647Z ##[group]Molecule ` + ansiEscape + `[32mdefault_install_remote_version` + ansiEscape + `[0m > ` + ansiEscape + `[32midempotence` + ansiEscape + `[0m
2026-11-17T17:34:13.2262396Z ` + ansiEscape + `[31m` + ansiEscape + `[1mCRITICAL` + ansiEscape + `[0m Idempotence test failed because of the following tasks:
2026-11-17T17:34:13.2262972Z * => signalfx.splunk_otel_collector.collector : Install Splunk OpenTelemetry Collector via dnf package manager
2026-11-17T17:34:13.2264153Z ##[endgroup]
2026-11-17T17:34:13.2264952Z ` + ansiEscape + `[31mERROR   ` + ansiEscape + `[0m ` + ansiEscape + `[32mdefault_install_remote_version` + ansiEscape + `[0m ➜ ` + ansiEscape + `[33midempotence` + ansiEscape + `[0m: ` + ansiEscape + `[31mExecuted: Failed` + ansiEscape + `[0m
2026-11-17T17:34:19.7903063Z ` + ansiEscape + `[31mERROR   ` + ansiEscape + `[0m Idempotence test failed because of the following tasks:
2026-11-17T17:34:19.7903736Z * => signalfx.splunk_otel_collector.collector : Install Splunk OpenTelemetry Collector via dnf package manager`

const ansibleMoleculeLogParserMissingStatusExample = `##[group]Molecule default_install_remote_version > idempotence
Idempotence test failed because of the following tasks:
* => signalfx.splunk_otel_collector.collector : Install Splunk OpenTelemetry Collector via dnf package manager
##[endgroup]`
const ansibleMoleculeLogParserMissingTaskExample = `##[group]Molecule default_install_remote_version > idempotence
ERROR default_install_remote_version ➜ idempotence: Executed: Failed
Idempotence test failed because of the following tasks:
##[endgroup]`
const ansibleMoleculeLogParserNonIdempotenceFailureExample = `##[group]Molecule default_install_remote_version > converge
ERROR default_install_remote_version ➜ converge: Executed: Failed
Idempotence test failed because of the following tasks:
* => signalfx.splunk_otel_collector.collector : Install Splunk OpenTelemetry Collector via dnf package manager`

const ansibleMoleculeLogParserNonMoleculeFailureExample = `ERROR default_install_remote_version ➜ idempotence: Executed: Failed
Idempotence test failed because of the following tasks:
* => signalfx.splunk_otel_collector.collector : Install Splunk OpenTelemetry Collector via dnf package manager`

const ansibleMoleculeLogParserExpectedFailureID = "default_install_remote_version/signalfx.splunk_otel_collector.collector/Install Splunk OpenTelemetry Collector via dnf package manager"
const ansibleMoleculeLogParserExpectedFailureMessage = `default_install_remote_version ➜ idempotence: Executed: Failed
Idempotence test failed because of the following tasks:
* => signalfx.splunk_otel_collector.collector : Install Splunk OpenTelemetry Collector via dnf package manager`
const ansibleMoleculeLogParserTestNameFailure = "parses failed Molecule idempotence task summaries"
const ansibleMoleculeLogParserTestNameMissingGroup = "requires a matching Molecule scenario group"
const ansibleMoleculeLogParserTestNameMissingStatus = "requires an explicit failed status"
const ansibleMoleculeLogParserTestNameMissingTask = "requires a failed task entry"
const ansibleMoleculeLogParserTestNameWrongAction = "requires an idempotence failure action"

func TestAnsibleMoleculeLogParserParsesSpecificIdempotenceFailures(t *testing.T) {
	testCases := []struct {
		log             string
		name            string
		wantFailureID   string
		wantFailureText string
		wantCanParse    bool
	}{
		{
			log:             ansibleMoleculeLogParserFailureExample,
			name:            ansibleMoleculeLogParserTestNameFailure,
			wantFailureID:   ansibleMoleculeLogParserExpectedFailureID,
			wantFailureText: ansibleMoleculeLogParserExpectedFailureMessage,
			wantCanParse:    true,
		},
		{
			log:  ansibleMoleculeLogParserNonMoleculeFailureExample,
			name: ansibleMoleculeLogParserTestNameMissingGroup,
		},
		{
			log:  ansibleMoleculeLogParserNonIdempotenceFailureExample,
			name: ansibleMoleculeLogParserTestNameWrongAction,
		},
		{
			log:  ansibleMoleculeLogParserMissingStatusExample,
			name: ansibleMoleculeLogParserTestNameMissingStatus,
		},
		{
			log:  ansibleMoleculeLogParserMissingTaskExample,
			name: ansibleMoleculeLogParserTestNameMissingTask,
		},
	}
	parser := AnsibleMoleculeLogParser{}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := parser.CanParse(testCase.log); got != testCase.wantCanParse {
				t.Fatalf("CanParse() = %t, want %t", got, testCase.wantCanParse)
			}
			if !testCase.wantCanParse {
				return
			}
			selectedParser, err := selectTestLogParser(testCase.log, defaultTestLogParsers())
			if err != nil {
				t.Fatalf("selectTestLogParser() error = %v", err)
			}
			if _, selectedAnsibleParser := selectedParser.(AnsibleMoleculeLogParser); !selectedAnsibleParser {
				t.Fatalf("selected parser = %T, want AnsibleMoleculeLogParser", selectedParser)
			}
			parsedTests, err := parser.Parse(testCase.log)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if len(parsedTests) != 1 {
				t.Fatalf("Parse() returned %d tests, want 1: %#v", len(parsedTests), parsedTests)
			}
			if parsedTests[0].ID != testCase.wantFailureID {
				t.Errorf("test ID = %q, want %q", parsedTests[0].ID, testCase.wantFailureID)
			}
			if parsedTests[0].Kind != testresults.TestKindTest {
				t.Errorf("test kind = %q, want %q", parsedTests[0].Kind, testresults.TestKindTest)
			}
			if parsedTests[0].Outcome != testresults.OutcomeFailed {
				t.Errorf("test outcome = %q, want %q", parsedTests[0].Outcome, testresults.OutcomeFailed)
			}
			if parsedTests[0].FailureMessage != testCase.wantFailureText {
				t.Errorf("failure message = %q, want %q", parsedTests[0].FailureMessage, testCase.wantFailureText)
			}
		})
	}
}
