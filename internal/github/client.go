package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"ktlo/internal/graphql"
)

const (
	apiEndpoint                = "https://api.github.com/graphql"
	authorizationHeader        = "Authorization"
	bearerPrefix               = "Bearer "
	emptyString                = ""
	errorInvalidRepository     = "%w %q: expected an HTTPS GitHub.com URL in owner/repository form"
	errorInvalidRepositoryKind = "invalid repository URL"
	errorMissingTokenKind      = "GitHub token is required; set GH_TOKEN or GITHUB_TOKEN"
	errorRepositoryUnavailable = "repository %q is not available through the GitHub GraphQL API"
	githubHost                 = "github.com"
	httpsScheme                = "https"
	pathSeparator              = "/"
	redactedToken              = "[REDACTED]"
	repositorySuffix           = ".git"
)

const (
	acceptHeader              = "Accept"
	apiVersion                = "2026-03-10"
	apiVersionHeader          = "X-GitHub-Api-Version"
	contentTypeJSON           = "application/vnd.github+json"
	defaultRESTRequestTimeout = 10 * time.Second
	maxRESTErrorBodyBytes     = 8192
	maxRESTResponseBytes      = 10 << 20
)

var (
	ErrInvalidRepositoryURL = errors.New(errorInvalidRepositoryKind)
	ErrMissingToken         = errors.New(errorMissingTokenKind)
)

type Client struct {
	graphqlClient  *graphql.Client
	httpClient     *http.Client
	owner          string
	repositoryName string
	repositoryURL  string
	token          string
}

type repositoryVariables struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

func NewClient(ctx context.Context, repositoryURL string, token string, httpClient *http.Client) (*Client, error) {
	token = strings.TrimSpace(token)
	if token == emptyString {
		return nil, ErrMissingToken
	}
	client := &Client{httpClient: httpClient, repositoryURL: repositoryURL, token: token}
	if err := client.validateRepositoryURL(); err != nil {
		return nil, err
	}
	headers := make(http.Header)
	headers.Set(authorizationHeader, bearerPrefix+token)
	client.graphqlClient = graphql.NewClient(apiEndpoint, httpClient, repositoryVariables{
		Name:  client.repositoryName,
		Owner: client.owner,
	}, headers)
	if err := client.checkAccess(ctx); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) validateRepositoryURL() error {
	parsedRepository, err := url.Parse(c.repositoryURL)
	if err != nil || !parsedRepository.IsAbs() || parsedRepository.Scheme != httpsScheme || !strings.EqualFold(parsedRepository.Hostname(), githubHost) || parsedRepository.User != nil || parsedRepository.RawQuery != emptyString || parsedRepository.Fragment != emptyString {
		return fmt.Errorf(errorInvalidRepository, ErrInvalidRepositoryURL, c.repositoryURL)
	}
	repositoryParts := strings.Split(strings.Trim(parsedRepository.Path, pathSeparator), pathSeparator)
	if len(repositoryParts) != 2 || repositoryParts[0] == emptyString {
		return fmt.Errorf(errorInvalidRepository, ErrInvalidRepositoryURL, c.repositoryURL)
	}
	c.owner = repositoryParts[0]
	c.repositoryName = strings.TrimSuffix(repositoryParts[1], repositorySuffix)
	if c.repositoryName == emptyString {
		return fmt.Errorf(errorInvalidRepository, ErrInvalidRepositoryURL, c.repositoryURL)
	}
	return nil
}

func (c *Client) Query(ctx context.Context, query string, result any) error {
	return c.redactToken(c.graphqlClient.Query(ctx, query, result))
}

func (c *Client) queryWithVariables(ctx context.Context, query string, variables any, result any) error {
	return c.redactToken(c.graphqlClient.QueryWithVariables(ctx, query, variables, result))
}

func (c *Client) getRESTJSON(ctx context.Context, endpoint string, result any) error {
	var err error
	var response *http.Response
	response, err = c.getRESTResponse(ctx, endpoint)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var responseBody []byte
	responseBody, err = io.ReadAll(io.LimitReader(response.Body, maxRESTResponseBytes+1))
	if err != nil {
		return c.redactToken(fmt.Errorf("could not read GitHub API response from %s: %w", endpoint, err))
	}
	if len(responseBody) > maxRESTResponseBytes {
		return fmt.Errorf("GitHub API response from %s exceeds %d bytes", endpoint, maxRESTResponseBytes)
	}
	if err := json.Unmarshal(responseBody, result); err != nil {
		return c.redactToken(fmt.Errorf("invalid GitHub API response from %s: %w", endpoint, err))
	}
	return nil
}

func (c *Client) getRESTResponse(ctx context.Context, endpoint string) (*http.Response, error) {
	var err error
	var request *http.Request
	request, err = http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, c.redactToken(fmt.Errorf("could not build GitHub API request to %s: %w", endpoint, err))
	}
	request.Header.Set(acceptHeader, contentTypeJSON)
	request.Header.Set(authorizationHeader, bearerPrefix+c.token)
	request.Header.Set(apiVersionHeader, apiVersion)

	httpClient := c.httpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultRESTRequestTimeout}
	}
	var response *http.Response
	response, err = httpClient.Do(request)
	if err != nil {
		return nil, c.redactToken(fmt.Errorf("GitHub API request to %s failed: %w", endpoint, err))
	}
	if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices {
		return response, nil
	}
	defer response.Body.Close()
	responseBody, readErr := io.ReadAll(io.LimitReader(response.Body, maxRESTErrorBodyBytes))
	if readErr != nil {
		return nil, c.redactToken(fmt.Errorf("GitHub API request to %s returned HTTP status %s and its error response could not be read: %w", endpoint, response.Status, readErr))
	}
	return nil, c.redactToken(fmt.Errorf("GitHub API request to %s returned HTTP status %s: %s", endpoint, response.Status, strings.TrimSpace(string(responseBody))))
}

func (c *Client) redactToken(err error) error {
	if err == nil || c.token == emptyString {
		return err
	}
	return tokenRedactedError{cause: err, token: c.token}
}

type tokenRedactedError struct {
	cause error
	token string
}

func (e tokenRedactedError) Error() string {
	return strings.ReplaceAll(e.cause.Error(), e.token, redactedToken)
}

func (e tokenRedactedError) Unwrap() error {
	return e.cause
}
