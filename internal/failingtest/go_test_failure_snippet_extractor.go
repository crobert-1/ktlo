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
	"bufio"
	"strings"

	"ktlo/internal/testresults"
)

type goTestJSONFailureSnippetExtractor struct{}

func (goTestJSONFailureSnippetExtractor) Extract(stepLog string) string {
	parser := GoTestJSONLogParser{}
	if !parser.CanParse(stepLog) {
		return ""
	}
	tests, err := parser.Parse(stepLog)
	if err != nil {
		return ""
	}
	for testIndex := len(tests) - 1; testIndex >= 0; testIndex-- {
		test := &tests[testIndex]
		if test.Outcome != testresults.OutcomeFailed || test.FailureMessage == "" {
			continue
		}
		return limitFailureLogSnippet(test.ID + logLineSeparator + test.FailureMessage)
	}
	return ""
}

type goTestFailureSnippetExtractor struct{}

func (goTestFailureSnippetExtractor) Extract(stepLog string) string {
	parser := GoTestTextLogParser{}
	parseState := newGoTestLogParseState()
	scanner := bufio.NewScanner(strings.NewReader(stepLog))
	latestFailureSnippet := ""
	for scanner.Scan() {
		parseState.processLine(parser, parser.cleanGoTestLogLine(strings.TrimRight(scanner.Text(), "\r")))
		if parseState.failureLogTest != nil {
			latestFailureSnippet = parseState.failureLogTest.failureSnippet()
			parseState.failureLogTest = nil
		}
	}
	parseState.finalizeFailureLogTest()
	if parseState.failureLogTest != nil {
		latestFailureSnippet = parseState.failureLogTest.failureSnippet()
	}
	return latestFailureSnippet
}
