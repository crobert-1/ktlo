package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	acceptHeader          = "Accept"
	contentTypeHeader     = "Content-Type"
	defaultRequestTimeout = 10 * time.Second
	emptyString           = ""
	errorBuildRequest     = "could not build GraphQL request: %w"
	errorEncodeRequest    = "could not encode GraphQL request: %w"
	errorEmptyQuery       = "GraphQL query must not be empty"
	errorEmptyResult      = "GraphQL result must not be nil"
	errorGraphQLErrors    = "GraphQL query failed: %s"
	errorInvalidResponse  = "invalid GraphQL response: %w"
	errorRequest          = "GraphQL request to %s failed: %w"
	errorResponseStatus   = "GraphQL request returned HTTP status %s: %s"
	errorResponseTooLarge = "GraphQL response exceeds %d bytes"
	errorUnexpectedData   = "GraphQL response does not contain data"
	graphqlDataKey        = "data"
	graphqlErrorsKey      = "errors"
	graphqlMessageKey     = "message"
	jsonMediaType         = "application/json"
	maxResponseBytes      = 1 << 20
	messageSeparator      = "; "
	nullJSONValue         = "null"
)

type Client struct {
	endpoint   string
	headers    http.Header
	httpClient *http.Client
	variables  any
}

func NewClient(endpoint string, httpClient *http.Client, variables any, headers http.Header) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultRequestTimeout}
	}
	return &Client{
		endpoint:   endpoint,
		headers:    cloneHeader(headers),
		httpClient: httpClient,
		variables:  variables,
	}
}

func (c *Client) Query(ctx context.Context, query string, result any) error {
	return c.query(ctx, query, c.variables, result)
}

func (c *Client) QueryWithVariables(ctx context.Context, query string, variables any, result any) error {
	return c.query(ctx, query, variables, result)
}

func (c *Client) query(ctx context.Context, query string, variables any, result any) error {
	if strings.TrimSpace(query) == emptyString {
		return errors.New(errorEmptyQuery)
	}
	if result == nil {
		return errors.New(errorEmptyResult)
	}

	encodedRequest, err := json.Marshal(queryRequest{
		Query:     query,
		Variables: variables,
	})
	if err != nil {
		return fmt.Errorf(errorEncodeRequest, err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(encodedRequest))
	if err != nil {
		return fmt.Errorf(errorBuildRequest, err)
	}
	request.Header.Set(acceptHeader, jsonMediaType)
	request.Header.Set(contentTypeHeader, jsonMediaType)
	for name, values := range c.headers {
		for _, value := range values {
			request.Header.Add(name, value)
		}
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf(errorRequest, c.endpoint, err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil {
		return fmt.Errorf(errorInvalidResponse, err)
	}
	if len(responseBody) > maxResponseBytes {
		return fmt.Errorf(errorResponseTooLarge, maxResponseBytes)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf(errorResponseStatus, response.Status, strings.TrimSpace(string(responseBody)))
	}

	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(responseBody, &envelope); err != nil {
		return fmt.Errorf(errorInvalidResponse, err)
	}
	if err := checkGraphQLErrors(envelope); err != nil {
		return err
	}

	data, ok := envelope[graphqlDataKey]
	if !ok || bytes.Equal(bytes.TrimSpace(data), []byte(nullJSONValue)) {
		return errors.New(errorUnexpectedData)
	}
	if err := json.Unmarshal(data, result); err != nil {
		return fmt.Errorf(errorInvalidResponse, err)
	}
	return nil
}

func cloneHeader(headers http.Header) http.Header {
	clonedHeaders := make(http.Header, len(headers))
	for name, values := range headers {
		clonedHeaders[name] = append([]string(nil), values...)
	}
	return clonedHeaders
}

type queryRequest struct {
	Query     string `json:"query"`
	Variables any    `json:"variables,omitempty"`
}

func checkGraphQLErrors(envelope map[string]json.RawMessage) error {
	encodedErrors, ok := envelope[graphqlErrorsKey]
	if !ok {
		return nil
	}

	var graphqlErrors []map[string]json.RawMessage
	if err := json.Unmarshal(encodedErrors, &graphqlErrors); err != nil {
		return fmt.Errorf(errorInvalidResponse, err)
	}
	if len(graphqlErrors) == 0 {
		return nil
	}

	messages := make([]string, 0, len(graphqlErrors))
	for _, graphqlError := range graphqlErrors {
		encodedMessage, ok := graphqlError[graphqlMessageKey]
		if !ok {
			continue
		}
		var message string
		if err := json.Unmarshal(encodedMessage, &message); err == nil && message != emptyString {
			messages = append(messages, message)
		}
	}
	if len(messages) == 0 {
		return fmt.Errorf(errorGraphQLErrors, string(encodedErrors))
	}
	return fmt.Errorf(errorGraphQLErrors, strings.Join(messages, messageSeparator))
}
