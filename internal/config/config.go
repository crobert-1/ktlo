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
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	defaultArtifactMaxAgeHours      = 48
	defaultLookback                 = "7d"
	emptyString                     = ""
	errorEmptyRepositories          = "config must include at least one repository"
	errorEmptyRepositoryURL         = "repository number %d must include a non-empty url"
	errorNoEnabledRepositories      = "config must include at least one enabled repository when a check is enabled"
	errorRepositoryEntryType        = "repository entry must be an object"
	errorLegacyRepositoryEntry      = "repository number %d uses the obsolete URL-only format; use an object with a url field"
	errorRepositoryUnknownField     = "repository entry has unknown field %q"
	errorInvalidArtifactMaxAgeHours = "failingTestCheck.maxArtifactAgeHours must not be less than zero"
	errorInvalidLookback            = "failingTestCheck.jsonReport.lookback must be a positive duration, such as 30m, 12h, or 7d"
	errorLoadConfig                 = "could not load config file %q: %w"
	errorMultipleDocuments          = "config file %q must contain a single YAML document"
	errorInvalidPRStatusStaleAge    = "prStatus.staleAge must be a positive duration, such as 60m, 12h, or 7d"
)

const defaultPRStatusStaleAge = "7d"

const defaultDependencyUpdateAuthorDependabot = "dependabot"

const defaultDependencyUpdateAuthorRenovate = "renovate"

const defaultDependencyUpdateAuthorSrvGHODiGDI = "srv-gh-o11y-gdi"

type Config struct {
	FailingTestCheck      FailingTestCheckConfig      `yaml:"failingTestCheck"`
	OTelGoInstrumentation OTelGoInstrumentationConfig `yaml:"otelGoInstrumentation"`
	PRStatus              PRStatusConfig              `yaml:"prStatus"`
	Repositories          []RepositoryConfig          `yaml:"repositories"`
	VulnerabilityCheck    VulnerabilityCheckConfig    `yaml:"vulnerabilityCheck"`
}

type RepositoryConfig struct {
	Enabled   bool   `yaml:"enabled"`
	LocalPath string `yaml:"localPath,omitempty"`
	URL       string `yaml:"url"`
	legacyURL bool
}

func (repository *RepositoryConfig) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode && node.Tag == yamlStringTag {
		repository.Enabled = true
		repository.URL = node.Value
		repository.legacyURL = true
		return nil
	}
	if node.Kind != yaml.MappingNode {
		return errors.New(errorRepositoryEntryType)
	}
	for keyIndex := 0; keyIndex < len(node.Content); keyIndex += 2 {
		key := node.Content[keyIndex].Value
		if key != repositoryEnabledKey && key != repositoryLocalPathKey && key != repositoryURLKey {
			return fmt.Errorf(errorRepositoryUnknownField, key)
		}
	}
	var decoded struct {
		Enabled   *bool  `yaml:"enabled"`
		LocalPath string `yaml:"localPath"`
		URL       string `yaml:"url"`
	}
	if err := node.Decode(&decoded); err != nil {
		return err
	}
	repository.Enabled = decoded.Enabled == nil || *decoded.Enabled
	repository.LocalPath = decoded.LocalPath
	repository.URL = decoded.URL
	return nil
}

type VulnerabilityCheckConfig struct {
	ID string `yaml:"id,omitempty"`
}

type OTelGoInstrumentationConfig struct {
	Enabled                         bool   `yaml:"enabled"`
	OTELExporterOTLPEndpoint        string `yaml:"OTEL_EXPORTER_OTLP_ENDPOINT,omitempty"`
	OTELExporterOTLPInsecure        string `yaml:"OTEL_EXPORTER_OTLP_INSECURE,omitempty"`
	OTELExporterOTLPMetricsEndpoint string `yaml:"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT,omitempty"`
	OTELExporterOTLPMetricsInsecure string `yaml:"OTEL_EXPORTER_OTLP_METRICS_INSECURE,omitempty"`
	OTELExporterOTLPTracesEndpoint  string `yaml:"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT,omitempty"`
	OTELExporterOTLPTracesInsecure  string `yaml:"OTEL_EXPORTER_OTLP_TRACES_INSECURE,omitempty"`
}

type FailingTestCheckConfig struct {
	JSONReport          FailingTestJSONReportConfig `yaml:"jsonReport"`
	MaxArtifactAgeHours *int                        `yaml:"maxArtifactAgeHours,omitempty"`
	Readme              FailingTestReadmeConfig     `yaml:"readme"`
}

type FailingTestJSONReportConfig struct {
	Enabled                 bool    `yaml:"enabled"`
	Lookback                *string `yaml:"lookback,omitempty"`
	OverwriteExistingReport bool    `yaml:"overwriteExistingReport"`
}

type FailingTestReadmeConfig struct {
	Enabled bool `yaml:"enabled"`
}

type PRStatusConfig struct {
	DependencyUpdateAuthors []string `yaml:"dependencyUpdateAuthors,omitempty"`
	Enabled                 *bool    `yaml:"enabled"`
	ParseWorkflowFailures   bool     `yaml:"parseWorkflowFailures"`
	StaleAge                *string  `yaml:"staleAge,omitempty"`
}

func Load(path string) (Config, error) {
	var config Config
	file, err := os.Open(path)
	if err != nil {
		return config, fmt.Errorf(errorLoadConfig, path, err)
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(&config); err != nil {
		return config, fmt.Errorf(errorLoadConfig, path, err)
	}
	var trailingDocument any
	if err := decoder.Decode(&trailingDocument); !errors.Is(err, io.EOF) {
		if err != nil {
			return config, fmt.Errorf(errorLoadConfig, path, err)
		}
		return config, fmt.Errorf(errorMultipleDocuments, path)
	}
	if err := config.Validate(); err != nil {
		return config, err
	}
	return config, nil
}

func (c PRStatusConfig) IsEnabled() bool {
	return c.Enabled == nil || *c.Enabled
}

func (c PRStatusConfig) EffectiveStaleAge() (time.Duration, error) {
	staleAge := defaultPRStatusStaleAge
	if c.StaleAge != nil {
		staleAge = *c.StaleAge
	}
	return parsePositiveDuration(staleAge, errorInvalidPRStatusStaleAge)
}

func (c PRStatusConfig) EffectiveDependencyUpdateAuthors() []string {
	if c.DependencyUpdateAuthors != nil {
		return append([]string(nil), c.DependencyUpdateAuthors...)
	}
	return []string{
		defaultDependencyUpdateAuthorDependabot,
		defaultDependencyUpdateAuthorRenovate,
		defaultDependencyUpdateAuthorSrvGHODiGDI,
	}
}

func (c FailingTestJSONReportConfig) IsEnabled() bool {
	return c.Enabled
}

func (c FailingTestJSONReportConfig) EffectiveLookback() (time.Duration, error) {
	lookback := defaultLookback
	if c.Lookback != nil {
		lookback = *c.Lookback
	}
	return parseLookbackDuration(lookback)
}

func (c FailingTestCheckConfig) EffectiveMaxArtifactAgeHours() int {
	if c.MaxArtifactAgeHours == nil {
		return defaultArtifactMaxAgeHours
	}
	return *c.MaxArtifactAgeHours
}

func (c VulnerabilityCheckConfig) HasID() bool {
	return strings.TrimSpace(c.ID) != emptyString
}

func (c Config) AnyCheckEnabled() bool {
	return c.PRStatus.IsEnabled() || c.FailingTestCheck.JSONReport.IsEnabled() || c.FailingTestCheck.Readme.Enabled || c.VulnerabilityCheck.HasID()
}

func (c Config) RequiresGitHubToken() bool {
	return c.PRStatus.IsEnabled() || c.FailingTestCheck.JSONReport.IsEnabled() || c.VulnerabilityCheck.HasID()
}

func (c Config) EnabledRepositories() []RepositoryConfig {
	repositories := make([]RepositoryConfig, 0, len(c.Repositories))
	for _, repository := range c.Repositories {
		if repository.Enabled {
			repositories = append(repositories, repository)
		}
	}
	return repositories
}

func (c Config) Validate() error {
	if c.FailingTestCheck.MaxArtifactAgeHours != nil && *c.FailingTestCheck.MaxArtifactAgeHours < 0 {
		return errors.New(errorInvalidArtifactMaxAgeHours)
	}
	if _, err := c.FailingTestCheck.JSONReport.EffectiveLookback(); err != nil {
		return err
	}
	if _, err := c.PRStatus.EffectiveStaleAge(); err != nil {
		return err
	}
	for index, repository := range c.Repositories {
		if repository.legacyURL {
			return fmt.Errorf(errorLegacyRepositoryEntry, index+1)
		}
		if strings.TrimSpace(repository.URL) == emptyString {
			return fmt.Errorf(errorEmptyRepositoryURL, index+1)
		}
	}
	if !c.AnyCheckEnabled() {
		return nil
	}
	if len(c.Repositories) == 0 {
		return errors.New(errorEmptyRepositories)
	}
	enabledRepositories := false
	for _, repository := range c.Repositories {
		if repository.Enabled {
			enabledRepositories = true
		}
	}
	if !enabledRepositories {
		return errors.New(errorNoEnabledRepositories)
	}
	return nil
}

const (
	repositoryEnabledKey   = "enabled"
	repositoryLocalPathKey = "localPath"
	repositoryURLKey       = "url"
	yamlStringTag          = "!!str"
)

func parseLookbackDuration(lookback string) (time.Duration, error) {
	return parsePositiveDuration(lookback, errorInvalidLookback)
}

func parsePositiveDuration(value string, errorMessage string) (time.Duration, error) {
	var duration time.Duration
	var err error
	value = strings.TrimSpace(value)
	if strings.HasSuffix(value, "d") {
		duration, err = time.ParseDuration(strings.TrimSuffix(value, "d") + "h")
		if err == nil && duration > 0 && duration <= time.Duration(1<<63-1)/24 {
			duration *= 24
		} else if err == nil {
			err = errors.New(errorMessage)
		}
	} else {
		duration, err = time.ParseDuration(value)
	}
	if err != nil || duration <= 0 {
		return 0, errors.New(errorMessage)
	}
	return duration, nil
}
