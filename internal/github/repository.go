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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"ktlo/internal/common"
	"ktlo/internal/telemetry"
)

const (
	graphqlRepositoryKey      = "repository"
	nullJSONValue             = "null"
	spanCheckRepositoryAccess = "github.check_repository_access"
)

type Repository = common.Repository

func ParseRepositoryURL(repositoryURL string) (Repository, error) {
	parsedRepository, err := url.Parse(repositoryURL)
	if err != nil || !parsedRepository.IsAbs() || parsedRepository.Scheme != httpsScheme || !strings.EqualFold(parsedRepository.Hostname(), githubHost) || parsedRepository.User != nil || parsedRepository.RawQuery != emptyString || parsedRepository.Fragment != emptyString {
		return Repository{}, fmt.Errorf(errorInvalidRepository, ErrInvalidRepositoryURL, repositoryURL)
	}
	repositoryParts := strings.Split(strings.Trim(parsedRepository.Path, pathSeparator), pathSeparator)
	if len(repositoryParts) != 2 || repositoryParts[0] == emptyString {
		return Repository{}, fmt.Errorf(errorInvalidRepository, ErrInvalidRepositoryURL, repositoryURL)
	}
	repositoryName := strings.TrimSuffix(repositoryParts[1], repositorySuffix)
	if repositoryName == emptyString {
		return Repository{}, fmt.Errorf(errorInvalidRepository, ErrInvalidRepositoryURL, repositoryURL)
	}
	return Repository{Owner: repositoryParts[0], Name: repositoryName}, nil
}

func (c *Client) TargetRepository() Repository {
	return Repository{
		Name:  c.repositoryName,
		Owner: c.owner,
	}
}

func (c *Client) checkAccess(ctx context.Context) (resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, githubInstrumentationName, spanCheckRepositoryAccess)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanCheckRepositoryAccess, startedAt, resultErr != nil)
	}()

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
