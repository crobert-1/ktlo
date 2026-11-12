package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"ktlo/internal/config"
	"ktlo/internal/failingtest"
	"ktlo/internal/github"
)

const commandName = "ktlo"
const emptyString = ""
const errorMissingConfig = "required flag --config is missing"
const errorUnexpectedArgs = "unexpected positional arguments: %v"
const flagConfig = "config"
const helpConfig = "Path to the KTLO YAML config file"
const usageText = "Usage: ktlo --config <config.yaml>\n"

func loadConfig(args []string, stderr io.Writer) (*config.Config, int) {
	flags := flag.NewFlagSet(commandName, flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() {
		_, _ = fmt.Fprint(stderr, usageText)
	}
	configPath := flags.String(flagConfig, emptyString, helpConfig)
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, 0
		}
		return nil, exitCodeUsage
	}
	if len(flags.Args()) > 0 {
		_, _ = fmt.Fprintf(stderr, errorUnexpectedArgs+lineEnding, flags.Args())
		return nil, exitCodeUsage
	}
	if *configPath == emptyString {
		_, _ = fmt.Fprintln(stderr, errorMissingConfig)
		_, _ = fmt.Fprint(stderr, usageText)
		return nil, exitCodeUsage
	}

	loadedConfig, err := config.Load(*configPath)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return nil, exitCodeUsage
	}
	return &loadedConfig, 0
}

func runFailingTestCheck(ctx context.Context, githubClient *github.Client, failingTestConfig config.FailingTestCheckConfig) error {
	workflowProvider, err := failingtest.NewGitHubWorkflowProvider(githubClient)
	if err != nil {
		return err
	}
	lookback, err := failingTestConfig.EffectiveLookback()
	if err != nil {
		return err
	}
	analyzer, err := failingtest.NewAnalyzerWithLookbackDuration(workflowProvider, lookback)
	if err != nil {
		return err
	}
	return analyzer.Run(ctx)
}
