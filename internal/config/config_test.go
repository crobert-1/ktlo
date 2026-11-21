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

package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const newConfigSchema = `log-level: detailed
telemetry:
  enabled: true
  otlp:
    endpoint: https://collector.example
    metrics:
      insecure: "true"
dependency-pr-status:
  update:
    dependency-update-authors: [dependabot]
    stale-age: 12h
flaky-test-analysis:
  lookback: 30m
repositories:
  - url: https://github.com/owner/repository
    local-path: /tmp/repository
`

const obsoleteConfigSchema = `logLevel: detailed
repositories:
  - url: https://github.com/owner/repository
`
)

func TestLoadConfigSchema(t *testing.T) {
	tests := []struct {
		name      string
		contents  string
		expectErr bool
	}{
		{
			name:     "command-aligned YAML keys",
			contents: newConfigSchema,
		},
		{
			name:      "obsolete YAML keys",
			contents:  obsoleteConfigSchema,
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			configPath := filepath.Join(t.TempDir(), "config.yaml")
			if err := os.WriteFile(configPath, []byte(test.contents), 0o600); err != nil {
				t.Fatalf("os.WriteFile() error = %v", err)
			}

			actual, err := Load(configPath)
			if (err != nil) != test.expectErr {
				t.Fatalf("Load() error = %v, want error %t", err, test.expectErr)
			}
			if err != nil {
				return
			}
			if actual.LogLevel != "detailed" {
				t.Errorf("LogLevel = %q, want detailed", actual.LogLevel)
			}
			if actual.Telemetry.OTLP.Endpoint != "https://collector.example" {
				t.Errorf("Telemetry.OTLP.Endpoint = %q, want configured endpoint", actual.Telemetry.OTLP.Endpoint)
			}
			if actual.DependencyPRStatus.Update.StaleAge == nil || *actual.DependencyPRStatus.Update.StaleAge != "12h" {
				t.Errorf("DependencyPRStatus.Update.StaleAge = %v, want 12h", actual.DependencyPRStatus.Update.StaleAge)
			}
			if actual.FlakyTestAnalysis.Lookback == nil || *actual.FlakyTestAnalysis.Lookback != "30m" {
				t.Errorf("FlakyTestAnalysis.Lookback = %v, want 30m", actual.FlakyTestAnalysis.Lookback)
			}
			if actual.Repositories[0].LocalPath != "/tmp/repository" {
				t.Errorf("Repositories[0].LocalPath = %q, want configured local path", actual.Repositories[0].LocalPath)
			}
		})
	}
}

func TestParseLookbackDuration(t *testing.T) {
	tests := []struct {
		name      string
		lookback  string
		expected  time.Duration
		expectErr bool
	}{
		{
			name:     "minutes",
			lookback: "30m",
			expected: 30 * time.Minute,
		},
		{
			name:     "hours",
			lookback: "12h",
			expected: 12 * time.Hour,
		},
		{
			name:     "days",
			lookback: "7d",
			expected: 7 * 24 * time.Hour,
		},
		{
			name:      "empty",
			expectErr: true,
		},
		{
			name:      "zero duration",
			lookback:  "0s",
			expectErr: true,
		},
		{
			name:      "negative duration",
			lookback:  "-1h",
			expectErr: true,
		},
		{
			name:      "invalid duration",
			lookback:  "tomorrow",
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := parseLookbackDuration(test.lookback)
			if (err != nil) != test.expectErr {
				t.Fatalf("parseLookbackDuration(%q) error = %v, want error %t", test.lookback, err, test.expectErr)
			}
			if err == nil && actual != test.expected {
				t.Errorf("parseLookbackDuration(%q) = %s, want %s", test.lookback, actual, test.expected)
			}
		})
	}
}
