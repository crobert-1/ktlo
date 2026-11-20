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

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"

	"ktlo/internal/graphql"
	"ktlo/internal/telemetry"
)

const (
	apiEndpoint                     = "https://api.github.com/graphql"
	authorizationHeader             = "Authorization"
	bearerPrefix                    = "Bearer "
	emptyString                     = ""
	errorCloseRESTResponse          = "could not close GitHub API response from %s: %w"
	errorInvalidRepository          = "%w %q: expected an HTTPS GitHub.com URL in owner/repository form"
	errorInvalidRepositoryKind      = "invalid repository URL"
	errorMissingTokenKind           = "GitHub token is required; set GH_TOKEN or GITHUB_TOKEN"
	errorRepositoryUnavailable      = "repository %q is not available through the GitHub GraphQL API"
	githubHost                      = "github.com"
	httpsScheme                     = "https"
	pathSeparator                   = "/"
	redactedToken                   = "[REDACTED]"
	repositorySuffix                = ".git"
	spanRESTJSON                    = "github.get_rest_json"
	traceAttributeHTTPRequestMethod = "http.request.method"
	traceAttributeURLFull           = "url.full"
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
	client := &Client{httpClient: instrumentHTTPClient(httpClient), repositoryURL: repositoryURL, token: token}
	if err := client.validateRepositoryURL(); err != nil {
		return nil, err
	}
	headers := make(http.Header)
	headers.Set(authorizationHeader, bearerPrefix+token)
	client.graphqlClient = graphql.NewClient(apiEndpoint, client.httpClient, repositoryVariables{
		Name:  client.repositoryName,
		Owner: client.owner,
	}, headers)
	if err := client.checkAccess(ctx); err != nil {
		return nil, err
	}
	return client, nil
}

func instrumentHTTPClient(httpClient *http.Client) *http.Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultRESTRequestTimeout}
	}
	instrumentedClient := *httpClient
	instrumentedClient.Transport = otelhttp.NewTransport(httpClient.Transport, otelhttp.WithPropagators(propagation.NewCompositeTextMapPropagator()))
	return &instrumentedClient
}

func (c *Client) validateRepositoryURL() error {
	repository, err := ParseRepositoryURL(c.repositoryURL)
	if err != nil {
		return err
	}
	c.owner = repository.Owner
	c.repositoryName = repository.Name
	return nil
}

func (c *Client) Query(ctx context.Context, query string, result any) error {
	return c.redactToken(c.graphqlClient.Query(ctx, query, result))
}

func (c *Client) queryWithVariables(ctx context.Context, query string, variables any, result any) error {
	return c.redactToken(c.graphqlClient.QueryWithVariables(ctx, query, variables, result))
}

func (c *Client) getRESTJSON(ctx context.Context, endpoint string, result any) (resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, githubInstrumentationName, spanRESTJSON,
		attribute.String(traceAttributeHTTPRequestMethod, http.MethodGet),
		attribute.String(traceAttributeURLFull, endpoint),
	)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanRESTJSON, startedAt, resultErr != nil)
	}()

	var err error
	var response *http.Response
	response, err = c.getRESTResponse(ctx, endpoint)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			resultErr = errors.Join(resultErr, c.redactToken(fmt.Errorf(errorCloseRESTResponse, endpoint, closeErr)))
		}
	}()

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

func (c *Client) getRESTResponse(ctx context.Context, endpoint string) (response *http.Response, resultErr error) {
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
	response, err = httpClient.Do(request)
	if err != nil {
		return nil, c.redactToken(fmt.Errorf("GitHub API request to %s failed: %w", endpoint, err))
	}
	if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices {
		return response, nil
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			resultErr = errors.Join(resultErr, c.redactToken(fmt.Errorf(errorCloseRESTResponse, endpoint, closeErr)))
		}
	}()
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
