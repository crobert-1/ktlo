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

const (
	goTestJSONLogParserPassedExample = `{"Action":"pass","Package":"example.com/project/pkg","Test":"TestPassed"}
{"Action":"pass","Package":"example.com/project/pkg","Elapsed":0.01}`
	goTestJSONLogParserFailedExample = `{"Action":"run","Package":"example.com/project/pkg","Test":"TestFailed"}
{"Action":"output","Package":"example.com/project/pkg","Test":"TestFailed","Output":"    assertion failed\n"}
{"Action":"fail","Package":"example.com/project/pkg","Test":"TestFailed","Elapsed":0.01}
{"Action":"fail","Package":"example.com/project/pkg","Elapsed":0.01}`
	goTestJSONLogParserSkippedExample = `{"Action":"skip","Package":"example.com/project/pkg","Test":"TestSkipped"}
{"Action":"pass","Package":"example.com/project/pkg","Elapsed":0.01}`
	goTestJSONLogParserUnrelatedJSON           = `{"Action":"error","Package":"example.com/project/app","message":"service failed"}`
	goTestJSONLogParserTestPackage             = "example.com/project/pkg"
	goTestJSONLogParserTestPassed              = "TestPassed"
	goTestJSONLogParserTestFailed              = "TestFailed"
	goTestJSONLogParserTestSkipped             = "TestSkipped"
	goTestJSONLogParserFailureText             = "    assertion failed\n"
	goTestJSONLogParserTestNamePass            = "parses package and test pass outcomes"
	goTestJSONLogParserTestNameFail            = "parses test and package failures"
	goTestJSONLogParserTestNameSkip            = "preserves skipped test outcome"
	goTestJSONLogParserTestNameFormatSelection = "selects the most specific matching format parser"
	goTestJSONLogParserTestNameTextSelection   = "selects text parser when JSON does not match"
)

func TestGoTestJSONLogParserParsesEvents(t *testing.T) {
	testCases := []struct {
		log      string
		name     string
		outcomes map[string]testresults.Outcome
	}{
		{
			log:  goTestJSONLogParserPassedExample,
			name: goTestJSONLogParserTestNamePass,
			outcomes: map[string]testresults.Outcome{
				goTestJSONLogParserTestPackage: testresults.OutcomePassed,
				goTestJSONLogParserTestPackage + goTestJSONTestIDSeparator + goTestJSONLogParserTestPassed: testresults.OutcomePassed,
			},
		},
		{
			log:  goTestJSONLogParserFailedExample,
			name: goTestJSONLogParserTestNameFail,
			outcomes: map[string]testresults.Outcome{
				goTestJSONLogParserTestPackage: testresults.OutcomeFailed,
				goTestJSONLogParserTestPackage + goTestJSONTestIDSeparator + goTestJSONLogParserTestFailed: testresults.OutcomeFailed,
			},
		},
		{
			log:  goTestJSONLogParserSkippedExample,
			name: goTestJSONLogParserTestNameSkip,
			outcomes: map[string]testresults.Outcome{
				goTestJSONLogParserTestPackage + goTestJSONTestIDSeparator + goTestJSONLogParserTestSkipped: testresults.OutcomeSkipped,
				goTestJSONLogParserTestPackage: testresults.OutcomePassed,
			},
		},
	}
	parser := GoTestJSONLogParser{}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if !parser.CanParse(testCase.log) {
				t.Fatal("CanParse() = false, want true")
			}
			parsedTests, err := parser.Parse(testCase.log)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if len(parsedTests) != len(testCase.outcomes) {
				t.Fatalf("Parse() returned %d observations, want %d: %#v", len(parsedTests), len(testCase.outcomes), parsedTests)
			}
			for _, parsedTest := range parsedTests {
				wantOutcome, found := testCase.outcomes[parsedTest.ID]
				if !found {
					t.Errorf("unexpected result ID %q", parsedTest.ID)
					continue
				}
				if parsedTest.Outcome != wantOutcome {
					t.Errorf("result %q has outcome %q, want %q", parsedTest.ID, parsedTest.Outcome, wantOutcome)
				}
				if parsedTest.ID == goTestJSONLogParserTestPackage+goTestJSONTestIDSeparator+goTestJSONLogParserTestFailed && parsedTest.FailureMessage != goTestJSONLogParserFailureText {
					t.Errorf("failure message = %q, want %q", parsedTest.FailureMessage, goTestJSONLogParserFailureText)
				}
			}
		})
	}
}

func TestDefaultTestLogParsersAreOrderedBySpecificity(t *testing.T) {
	parsers := defaultTestLogParsers()
	if len(parsers) != 3 {
		t.Fatalf("default parser count = %d, want 3", len(parsers))
	}
	if _, isMoleculeParser := parsers[0].(AnsibleMoleculeLogParser); !isMoleculeParser {
		t.Fatalf("first parser type = %T, want AnsibleMoleculeLogParser", parsers[0])
	}
	if _, isJSONParser := parsers[1].(GoTestJSONLogParser); !isJSONParser {
		t.Fatalf("second parser type = %T, want GoTestJSONLogParser", parsers[1])
	}
	if _, isTextParser := parsers[2].(GoTestTextLogParser); !isTextParser {
		t.Fatalf("third parser type = %T, want GoTestTextLogParser", parsers[2])
	}
	for parserIndex := 1; parserIndex < len(parsers); parserIndex++ {
		if parsers[parserIndex-1].Specificity() <= parsers[parserIndex].Specificity() {
			t.Fatalf("parser specificities are not descending: %d then %d", parsers[parserIndex-1].Specificity(), parsers[parserIndex].Specificity())
		}
	}
}

func TestSelectTestLogParserPrefersHighestSpecificity(t *testing.T) {
	lessSpecificParser := testLogParserStub{specificity: goTestTextParserSpecificity, matches: true}
	moreSpecificParser := testLogParserStub{specificity: goTestJSONParserSpecificity, matches: true}
	selectedParser, err := selectTestLogParser("", []TestLogParser{lessSpecificParser, moreSpecificParser})
	if err != nil {
		t.Fatalf("selectTestLogParser() error = %v", err)
	}
	if selectedParser.Specificity() != moreSpecificParser.specificity {
		t.Fatalf("selected specificity = %d, want %d", selectedParser.Specificity(), moreSpecificParser.specificity)
	}
}

func TestSelectTestLogParserByFormatSpecificity(t *testing.T) {
	testCases := []struct {
		log      string
		name     string
		wantJSON bool
	}{
		{
			log:      goTestJSONLogParserPassedExample,
			name:     goTestJSONLogParserTestNameFormatSelection,
			wantJSON: true,
		},
		{
			log:      goTestTextLogParserWindowsIISFailureExample,
			name:     goTestJSONLogParserTestNameTextSelection,
			wantJSON: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			parser, err := selectTestLogParser(testCase.log, defaultTestLogParsers())
			if err != nil {
				t.Fatalf("selectTestLogParser() error = %v", err)
			}
			_, isJSONParser := parser.(GoTestJSONLogParser)
			if isJSONParser != testCase.wantJSON {
				t.Errorf("selected parser = %T, want GoTestJSONLogParser = %t", parser, testCase.wantJSON)
			}
		})
	}
}

func TestSelectTestLogParserRejectsEqualSpecificityMatches(t *testing.T) {
	firstParser := testLogParserStub{specificity: goTestJSONParserSpecificity, matches: true}
	secondParser := testLogParserStub{specificity: goTestJSONParserSpecificity, matches: true}
	if _, err := selectTestLogParser("", []TestLogParser{firstParser, secondParser}); err == nil {
		t.Fatal("selectTestLogParser() error = nil, want ambiguity error")
	}
}

func TestGoTestJSONLogParserDoesNotClaimUnrelatedJSON(t *testing.T) {
	if (GoTestJSONLogParser{}).CanParse(goTestJSONLogParserUnrelatedJSON) {
		t.Fatal("CanParse() = true, want false for an unrelated JSON event")
	}
}

type testLogParserStub struct {
	specificity int
	matches     bool
}

func (parser testLogParserStub) Specificity() int {
	return parser.specificity
}

func (parser testLogParserStub) CanParse(string) bool {
	return parser.matches
}

func (testLogParserStub) Parse(string) ([]testresults.Test, error) {
	return nil, nil
}
