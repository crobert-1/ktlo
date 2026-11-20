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
	"encoding/json"
	"strings"

	"ktlo/internal/testresults"
)

type goTestJSONAction string

const (
	goTestJSONActionBuildFail   = "build-fail"
	goTestJSONActionBuildOutput = "build-output"
	goTestJSONActionFail        = "fail"
	goTestJSONActionOutput      = "output"
	goTestJSONActionPass        = "pass"
	goTestJSONActionRun         = "run"
	goTestJSONActionSkip        = "skip"
	goTestJSONActionField       = "Action"
	goTestJSONPackageField      = "Package"
	goTestJSONTestIDSeparator   = "."
)

// GoTestJSONLogParser parses newline-delimited events emitted by go test -json.
type GoTestJSONLogParser struct{}

func (GoTestJSONLogParser) Specificity() int {
	return goTestJSONParserSpecificity
}

func (parser GoTestJSONLogParser) CanParse(stepLog string) bool {
	for _, line := range strings.Split(stepLog, logLineSeparator) {
		event, valid := parser.eventFromLine(line)
		if valid && event.identifiesTestLog() {
			return true
		}
	}
	return false
}

func (parser GoTestJSONLogParser) Parse(stepLog string) ([]testresults.Test, error) {
	parseState := newGoTestJSONLogParseState()
	for _, line := range strings.Split(stepLog, logLineSeparator) {
		event, valid := parser.eventFromLine(line)
		if valid {
			parseState.processEvent(event)
		}
	}
	return parseState.results, nil
}

func (GoTestJSONLogParser) eventFromLine(line string) (goTestJSONEvent, bool) {
	line = strings.TrimSpace((GoTestTextLogParser{}).cleanGoTestLogLine(line))
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return goTestJSONEvent{}, false
	}
	if _, hasAction := fields[goTestJSONActionField]; !hasAction {
		return goTestJSONEvent{}, false
	}
	if _, hasPackage := fields[goTestJSONPackageField]; !hasPackage {
		return goTestJSONEvent{}, false
	}
	var event goTestJSONEvent
	if err := json.Unmarshal([]byte(line), &event); err != nil {
		return goTestJSONEvent{}, false
	}
	event.Action = goTestJSONAction(strings.ToLower(strings.TrimSpace(string(event.Action))))
	event.Package = strings.TrimSpace(event.Package)
	event.Test = strings.TrimSpace(event.Test)
	return event, event.Action != ""
}

type goTestJSONEvent struct {
	Action  goTestJSONAction `json:"Action"`
	Package string           `json:"Package"`
	Test    string           `json:"Test"`
	Output  string           `json:"Output"`
	Elapsed float64          `json:"Elapsed"`
	Cached  bool             `json:"Cached"`
}

func (event *goTestJSONEvent) identifiesTestLog() bool {
	if event.Package == "" {
		return false
	}
	switch event.Action {
	case goTestJSONActionBuildFail:
		return true
	case goTestJSONActionFail, goTestJSONActionPass, goTestJSONActionSkip:
		return true
	default:
		return false
	}
}

func (action goTestJSONAction) outcome() (testresults.Outcome, bool) {
	switch goTestJSONAction(strings.ToLower(string(action))) {
	case goTestJSONActionPass:
		return testresults.OutcomePassed, true
	case goTestJSONActionFail:
		return testresults.OutcomeFailed, true
	case goTestJSONActionSkip:
		return testresults.OutcomeSkipped, true
	default:
		return "", false
	}
}

type goTestJSONTestKey struct {
	packageName string
	testName    string
}

type goTestJSONTestState struct {
	output string
}

type goTestJSONLogParseState struct {
	packageOutput map[string]string
	results       []testresults.Test
	testsByName   map[goTestJSONTestKey]*goTestJSONTestState
}

func newGoTestJSONLogParseState() *goTestJSONLogParseState {
	return &goTestJSONLogParseState{
		packageOutput: make(map[string]string),
		results:       make([]testresults.Test, 0),
		testsByName:   make(map[goTestJSONTestKey]*goTestJSONTestState),
	}
}

func (state *goTestJSONLogParseState) processEvent(event goTestJSONEvent) {
	if event.Package == "" {
		return
	}
	if event.Test == "" {
		state.processPackageEvent(event)
		return
	}
	state.processTestEvent(event)
}

func (state *goTestJSONLogParseState) processPackageEvent(event goTestJSONEvent) {
	switch event.Action {
	case goTestJSONActionBuildOutput, goTestJSONActionOutput:
		state.appendPackageOutput(event.Package, event.Output)
	case goTestJSONActionBuildFail:
		state.results = append(state.results, testresults.Test{
			ID:             event.Package,
			Kind:           testresults.TestKindPackage,
			Outcome:        testresults.OutcomeFailed,
			FailureMessage: state.packageFailureMessage(event.Package, event.Output),
		})
		delete(state.packageOutput, event.Package)
	case goTestJSONActionFail, goTestJSONActionPass:
		outcome, valid := event.Action.outcome()
		if !valid {
			return
		}
		packageOutput := strings.TrimSpace(state.packageOutput[event.Package] + event.Output)
		cached := event.Cached || strings.HasSuffix(packageOutput, goCachedResult)
		if outcome != testresults.OutcomePassed || !cached {
			result := testresults.Test{
				ID:      event.Package,
				Kind:    testresults.TestKindPackage,
				Outcome: outcome,
			}
			if outcome == testresults.OutcomeFailed {
				result.FailureMessage = state.packageFailureMessage(event.Package, event.Output)
			}
			state.results = append(state.results, result)
		}
		delete(state.packageOutput, event.Package)
	}
}

func (state *goTestJSONLogParseState) processTestEvent(event goTestJSONEvent) {
	key := goTestJSONTestKey{packageName: event.Package, testName: event.Test}
	switch event.Action {
	case goTestJSONActionRun:
		state.testsByName[key] = &goTestJSONTestState{}
	case goTestJSONActionOutput:
		testState := state.testStateFor(key)
		testState.appendOutput(event.Output)
	case goTestJSONActionFail, goTestJSONActionPass, goTestJSONActionSkip:
		outcome, valid := event.Action.outcome()
		if !valid {
			return
		}
		testState := state.testStateFor(key)
		testState.appendOutput(event.Output)
		result := testresults.Test{
			ID:      event.Package + goTestJSONTestIDSeparator + event.Test,
			Kind:    testresults.TestKindTest,
			Outcome: outcome,
		}
		if outcome == testresults.OutcomeFailed {
			result.FailureMessage = testState.output
		}
		state.results = append(state.results, result)
		delete(state.testsByName, key)
	}
}

func (state *goTestJSONLogParseState) testStateFor(key goTestJSONTestKey) *goTestJSONTestState {
	if state.testsByName[key] == nil {
		state.testsByName[key] = &goTestJSONTestState{}
	}
	return state.testsByName[key]
}

func (state *goTestJSONLogParseState) appendPackageOutput(packageName string, output string) {
	if output == "" {
		return
	}
	state.packageOutput[packageName] = limitFailureLogSnippet(state.packageOutput[packageName] + output)
}

func (state *goTestJSONLogParseState) packageFailureMessage(packageName string, output string) string {
	return limitFailureLogSnippet(state.packageOutput[packageName] + output)
}

func (state *goTestJSONTestState) appendOutput(output string) {
	if output == "" {
		return
	}
	state.output = limitFailureLogSnippet(state.output + output)
}
