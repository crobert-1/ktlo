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

package main

import (
	"testing"
)

const (
	testReleaseVersion       = "v1.2.3"
	testVulnerabilityCommand = "v"
	testVulnerabilityID      = "CVE-2026-1234"
)

func TestParseVulnerabilityCheckOptions(t *testing.T) {
	tests := []struct {
		arguments             []string
		expectedError         string
		expectedRelease       string
		expectedVulnerability string
		name                  string
	}{
		{
			name:                  "without release version",
			arguments:             []string{testVulnerabilityCommand, testVulnerabilityID},
			expectedVulnerability: testVulnerabilityID,
		},
		{
			name:                  "with release version",
			arguments:             []string{testVulnerabilityCommand, dashPrefix + dashPrefix + flagReleaseVersion, testReleaseVersion, testVulnerabilityID},
			expectedRelease:       testReleaseVersion,
			expectedVulnerability: testVulnerabilityID,
		},
		{
			name:          "release version is required when flag is present",
			arguments:     []string{testVulnerabilityCommand, dashPrefix + dashPrefix + flagReleaseVersion + equalsSign, testVulnerabilityID},
			expectedError: errorMissingReleaseVersion,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			invocation, _, err := parseCLIInvocation(test.arguments)
			if test.expectedError != emptyString {
				if err == nil || err.Error() != test.expectedError {
					t.Fatalf("expected error %q, got %v", test.expectedError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}
			if invocation.releaseVersion != test.expectedRelease {
				t.Errorf("expected release version %q, got %q", test.expectedRelease, invocation.releaseVersion)
			}
			if invocation.vulnerabilityID != test.expectedVulnerability {
				t.Errorf("expected vulnerability ID %q, got %q", test.expectedVulnerability, invocation.vulnerabilityID)
			}
		})
	}
}
