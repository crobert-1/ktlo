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
	errorEmptyRepository            = "repository number %d is empty"
	errorInvalidArtifactMaxAgeHours = "failingTestCheck.maxArtifactAgeHours must not be less than zero"
	errorInvalidLookback            = "failingTestCheck.lookback must be a positive duration, such as 30m, 12h, or 7d"
	errorLoadConfig                 = "could not load config file %q: %w"
	errorMultipleDocuments          = "config file %q must contain a single YAML document"
)

type Config struct {
	FailingTestCheck FailingTestCheckConfig `yaml:"failingTestCheck"`
	PRStatus         PRStatusConfig         `yaml:"prStatus"`
	Repositories     []string               `yaml:"repositories"`
}

type FailingTestCheckConfig struct {
	Enabled             bool    `yaml:"enabled"`
	Lookback            *string `yaml:"lookback,omitempty"`
	MaxArtifactAgeHours *int    `yaml:"maxArtifactAgeHours,omitempty"`
}

type PRStatusConfig struct {
	Enabled *bool `yaml:"enabled"`
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

func (c FailingTestCheckConfig) IsEnabled() bool {
	return c.Enabled
}

func (c FailingTestCheckConfig) EffectiveLookback() (time.Duration, error) {
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

func (c Config) AnyCheckEnabled() bool {
	return c.PRStatus.IsEnabled() || c.FailingTestCheck.IsEnabled()
}

func (c Config) Validate() error {
	if c.FailingTestCheck.MaxArtifactAgeHours != nil && *c.FailingTestCheck.MaxArtifactAgeHours < 0 {
		return errors.New(errorInvalidArtifactMaxAgeHours)
	}
	if _, err := c.FailingTestCheck.EffectiveLookback(); err != nil {
		return err
	}
	if !c.AnyCheckEnabled() {
		return nil
	}
	if len(c.Repositories) == 0 {
		return errors.New(errorEmptyRepositories)
	}
	for index, repository := range c.Repositories {
		if strings.TrimSpace(repository) == emptyString {
			return fmt.Errorf(errorEmptyRepository, index+1)
		}
	}
	return nil
}

func parseLookbackDuration(lookback string) (time.Duration, error) {
	var duration time.Duration
	var err error
	lookback = strings.TrimSpace(lookback)
	if strings.HasSuffix(lookback, "d") {
		duration, err = time.ParseDuration(strings.TrimSuffix(lookback, "d") + "h")
		if err == nil && duration > 0 && duration <= time.Duration(1<<63-1)/24 {
			duration *= 24
		} else if err == nil {
			err = errors.New(errorInvalidLookback)
		}
	} else {
		duration, err = time.ParseDuration(lookback)
	}
	if err != nil || duration <= 0 {
		return 0, errors.New(errorInvalidLookback)
	}
	return duration, nil
}
