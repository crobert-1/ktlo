package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"ktlo/internal/config"
	"ktlo/internal/github"
	"ktlo/internal/prstatus"
)

const (
	errorRepositoryFailure = "repository %q failed: %v"
	exitCodeError          = 1
	exitCodeUsage          = 2
	githubTokenEnv         = "GH_TOKEN"
	githubTokenFallbackEnv = "GITHUB_TOKEN"
	lineEnding             = "\n"
)

func main() {
	// Skip the executable name so run receives only command-line arguments.
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	config, exitCode := loadConfig(args, stderr)
	if exitCode != 0 || config == nil {
		return exitCode
	}
	if !config.AnyCheckEnabled() {
		return 0
	}

	token := strings.TrimSpace(os.Getenv(githubTokenEnv))
	if token == emptyString {
		token = strings.TrimSpace(os.Getenv(githubTokenFallbackEnv))
	}
	if token == emptyString {
		_, _ = fmt.Fprintln(stderr, github.ErrMissingToken)
		return exitCodeUsage
	}

	return runChecks(context.Background(), *config, token, stdout, stderr)
}

func runChecks(ctx context.Context, config config.Config, token string, stdout io.Writer, stderr io.Writer) int {
	exitCode := 0
	for _, repositoryURL := range config.Repositories {
		githubClient, err := github.NewClient(ctx, strings.TrimSpace(repositoryURL), token, nil)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, errorRepositoryFailure+lineEnding, repositoryURL, err)
			if errors.Is(err, github.ErrInvalidRepositoryURL) {
				return exitCodeUsage
			}
			exitCode = exitCodeError
			continue
		}
		if config.PRStatus.IsEnabled() {
			if err := runPRStatus(ctx, githubClient, stdout); err != nil {
				_, _ = fmt.Fprintf(stderr, errorRepositoryFailure+lineEnding, repositoryURL, err)
				exitCode = exitCodeError
			}
		}
		if config.FailingTestCheck.IsEnabled() {
			if err := runFailingTestCheck(ctx, githubClient, config.FailingTestCheck); err != nil {
				_, _ = fmt.Fprintf(stderr, errorRepositoryFailure+lineEnding, repositoryURL, err)
				exitCode = exitCodeError
			}
		}
	}
	return exitCode
}

func runPRStatus(ctx context.Context, githubClient *github.Client, stdout io.Writer) error {
	runner, err := prstatus.NewRunner(githubClient)
	if err != nil {
		return err
	}
	return runner.Run(ctx, stdout)
}
