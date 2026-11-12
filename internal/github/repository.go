package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

const (
	graphqlRepositoryKey = "repository"
	nullJSONValue        = "null"
)

type Repository struct {
	Name  string
	Owner string
}

func (c *Client) TargetRepository() Repository {
	return Repository{
		Name:  c.repositoryName,
		Owner: c.owner,
	}
}

func (c *Client) checkAccess(ctx context.Context) error {
	var result map[string]json.RawMessage
	if err := c.Query(ctx, repositoryAccessQuery, &result); err != nil {
		return err
	}
	repositoryData, ok := result[graphqlRepositoryKey]
	if !ok || bytes.Equal(bytes.TrimSpace(repositoryData), []byte(nullJSONValue)) {
		return fmt.Errorf(errorRepositoryUnavailable, c.repositoryURL)
	}
	return nil
}
