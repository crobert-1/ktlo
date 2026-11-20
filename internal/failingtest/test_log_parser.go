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
	"fmt"
	"sort"

	"ktlo/internal/testresults"
)

const ansibleMoleculeLogParserSpecificity = 300
const errorAmbiguousTestLogParsers = "multiple test log parsers with specificity %d recognize the same log: %T and %T"
const goTestJSONParserSpecificity = 200
const goTestTextParserSpecificity = 100
const traceAttributeTestLogParserSpecificity = "ktlo.test_log_parser.specificity"

// TestLogParser identifies and parses one test log format.
type TestLogParser interface {
	// Specificity is a larger value for a more specific log format.
	Specificity() int
	CanParse(stepLog string) bool
	Parse(stepLog string) ([]testresults.Test, error)
}

func defaultTestLogParsers() []TestLogParser {
	return orderTestLogParsers([]TestLogParser{AnsibleMoleculeLogParser{}, GoTestJSONLogParser{}, GoTestTextLogParser{}})
}

func orderTestLogParsers(parsers []TestLogParser) []TestLogParser {
	orderedParsers := append([]TestLogParser(nil), parsers...)
	sort.SliceStable(orderedParsers, func(firstIndex int, secondIndex int) bool {
		return orderedParsers[firstIndex].Specificity() > orderedParsers[secondIndex].Specificity()
	})
	return orderedParsers
}

func selectTestLogParser(stepLog string, parsers []TestLogParser) (TestLogParser, error) {
	orderedParsers := orderTestLogParsers(parsers)
	var selectedParser TestLogParser
	for _, parser := range orderedParsers {
		if !parser.CanParse(stepLog) {
			continue
		}
		if selectedParser == nil {
			selectedParser = parser
			continue
		}
		if parser.Specificity() < selectedParser.Specificity() {
			break
		}
		return nil, fmt.Errorf(errorAmbiguousTestLogParsers, parser.Specificity(), selectedParser, parser)
	}
	if selectedParser == nil {
		return nil, nil
	}
	return selectedParser, nil
}
