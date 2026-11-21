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

	"ktlo/internal/logging"

	"gopkg.in/yaml.v3"
)

const (
	defaultLookback              = "7d"
	emptyString                  = ""
	errorEmptyRepositories       = "config must include at least one repository"
	errorEmptyRepositoryURL      = "repository number %d must include a non-empty url"
	errorNoEnabledRepositories   = "config must include at least one enabled repository"
	errorRepositoryEntryType     = "repository entry must be an object"
	errorLegacyRepositoryEntry   = "repository number %d uses the obsolete URL-only format; use an object with a url field"
	errorRepositoryUnknownField  = "repository entry has unknown field %q"
	errorInvalidLookback         = "flaky-test-analysis.lookback must be a positive duration, such as 30m, 12h, or 7d"
	errorInvalidLogLevel         = "log-level must be none, default, or detailed"
	errorInvalidPRStatusStaleAge = "dependency-pr-status.update.stale-age must be a positive duration, such as 60m, 12h, or 7d"
	errorLoadConfig              = "could not load config file %q: %w"
	errorMultipleDocuments       = "config file %q must contain a single YAML document"
)

const defaultPRStatusStaleAge = "7d"

const defaultDependencyUpdateAuthorDependabot = "dependabot"

const defaultDependencyUpdateAuthorRenovate = "renovate"

const defaultDependencyUpdateAuthorSrvGHODiGDI = "srv-gh-o11y-gdi"

type Config struct {
	DependencyPRStatus DependencyPRStatusConfig `yaml:"dependency-pr-status"`
	FlakyTestAnalysis  FlakyTestAnalysisConfig  `yaml:"flaky-test-analysis"`
	LogLevel           logging.Level            `yaml:"log-level,omitempty"`
	Repositories       []RepositoryConfig       `yaml:"repositories"`
	Telemetry          TelemetryConfig          `yaml:"telemetry"`
}

type RepositoryConfig struct {
	Enabled   bool   `yaml:"enabled"`
	LocalPath string `yaml:"local-path,omitempty"`
	URL       string `yaml:"url"`
	legacyURL bool   `yaml:"-"`
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
		LocalPath string `yaml:"local-path"`
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

type TelemetryConfig struct {
	Enabled bool       `yaml:"enabled"`
	OTLP    OTLPConfig `yaml:"otlp"`
}

type OTLPConfig struct {
	Endpoint string                `yaml:"endpoint,omitempty"`
	Insecure string                `yaml:"insecure,omitempty"`
	Metrics  TelemetrySignalConfig `yaml:"metrics"`
	Traces   TelemetrySignalConfig `yaml:"traces"`
}

type TelemetrySignalConfig struct {
	Endpoint string `yaml:"endpoint,omitempty"`
	Insecure string `yaml:"insecure,omitempty"`
}

type DependencyPRStatusConfig struct {
	Update DependencyPRStatusUpdateConfig `yaml:"update"`
}

type DependencyPRStatusUpdateConfig struct {
	DependencyUpdateAuthors []string `yaml:"dependency-update-authors,omitempty"`
	StaleAge                *string  `yaml:"stale-age,omitempty"`
}

type FlakyTestAnalysisConfig struct {
	Lookback *string `yaml:"lookback,omitempty"`
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

func (c DependencyPRStatusUpdateConfig) EffectiveStaleAge() (time.Duration, error) {
	staleAge := defaultPRStatusStaleAge
	if c.StaleAge != nil {
		staleAge = *c.StaleAge
	}
	return parsePositiveDuration(staleAge, errorInvalidPRStatusStaleAge)
}

func (c DependencyPRStatusUpdateConfig) EffectiveDependencyUpdateAuthors() []string {
	if c.DependencyUpdateAuthors != nil {
		return append([]string(nil), c.DependencyUpdateAuthors...)
	}
	return []string{
		defaultDependencyUpdateAuthorDependabot,
		defaultDependencyUpdateAuthorRenovate,
		defaultDependencyUpdateAuthorSrvGHODiGDI,
	}
}

func (c FlakyTestAnalysisConfig) EffectiveLookback() (time.Duration, error) {
	lookback := defaultLookback
	if c.Lookback != nil {
		lookback = *c.Lookback
	}
	return parseLookbackDuration(lookback)
}

func (c Config) EffectiveLogLevel() logging.Level {
	if c.LogLevel == "" {
		return logging.LevelDefault
	}
	return c.LogLevel
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
	if c.LogLevel != "" && !c.LogLevel.IsValid() {
		return errors.New(errorInvalidLogLevel)
	}
	if _, err := c.FlakyTestAnalysis.EffectiveLookback(); err != nil {
		return err
	}
	if _, err := c.DependencyPRStatus.Update.EffectiveStaleAge(); err != nil {
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
	repositoryLocalPathKey = "local-path"
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
