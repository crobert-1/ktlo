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
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"ktlo/internal/common"
)

const checkSuiteJobsNodeQuery = `query WorkflowRunJobs($ids: [ID!]!) { nodes(ids: $ids) { ... on CheckSuite { id branch { name } commit { oid } conclusion status workflowRun { url runAttempt createdAt displayTitle workflow { name } } checkRuns(first: 100) { nodes { name conclusion status url startedAt steps(first: 100) { nodes { name number conclusion status startedAt completedAt } pageInfo { hasNextPage } } } pageInfo { hasNextPage } } } } } }`
const checkSuiteWorkflowRunNodeQuery = `query WorkflowRunMetadata($ids: [ID!]!) { nodes(ids: $ids) { ... on CheckRun { id checkSuite { id } } } }`
const defaultCheckRunNodeBatchSize = 100
const defaultCheckSuiteNodeBatchSize = 10
const errorQueryCheckSuiteIDs = "could not query check suite IDs for failed workflow check runs: %w"
const errorQueryCheckSuiteJobs = "could not query jobs for failed workflow check suites: %w"
const githubActionsPathSegment = "actions"
const githubRunsPathSegment = "runs"

// WorkflowRunCheckSuiteDetails contains run and job metadata linked to a check suite.
type WorkflowRunCheckSuiteDetails struct {
	Attempt    common.WorkflowRunAttempt
	Jobs       []common.WorkflowJob
	IsComplete bool
}

type workflowRunMetadataPage struct {
	Nodes []*workflowRunMetadataNode `json:"nodes"`
}

type workflowRunMetadataNode struct {
	ID         string              `json:"id"`
	CheckSuite *workflowCheckSuite `json:"checkSuite"`
}

type workflowCheckSuite struct {
	ID string `json:"id"`
}

type workflowCheckSuiteJobsPage struct {
	Nodes []*workflowCheckSuiteJobsNode `json:"nodes"`
}

type workflowCheckSuiteJobsNode struct {
	ID          string                      `json:"id"`
	Branch      workflowCheckSuiteBranch    `json:"branch"`
	CheckRuns   workflowCheckRunConnection  `json:"checkRuns"`
	Commit      workflowCheckSuiteCommit    `json:"commit"`
	Conclusion  WorkflowConclusion          `json:"conclusion"`
	Status      string                      `json:"status"`
	WorkflowRun *workflowRunGraphQLMetadata `json:"workflowRun"`
}

type workflowCheckSuiteBranch struct {
	Name string `json:"name"`
}

type workflowCheckSuiteCommit struct {
	OID string `json:"oid"`
}

type workflowRunGraphQLMetadata struct {
	CreatedAt    time.Time               `json:"createdAt"`
	DisplayTitle string                  `json:"displayTitle"`
	RunAttempt   int                     `json:"runAttempt"`
	URL          string                  `json:"url"`
	Workflow     workflowGraphQLMetadata `json:"workflow"`
}

type workflowGraphQLMetadata struct {
	Name string `json:"name"`
}

type workflowCheckRunConnection struct {
	Nodes    []workflowCheckRunNode `json:"nodes"`
	PageInfo workflowPageInfo       `json:"pageInfo"`
}

type workflowCheckRunNode struct {
	Conclusion WorkflowConclusion          `json:"conclusion"`
	Name       string                      `json:"name"`
	StartedAt  time.Time                   `json:"startedAt"`
	Status     string                      `json:"status"`
	Steps      workflowCheckStepConnection `json:"steps"`
	URL        string                      `json:"url"`
}

type workflowCheckStepConnection struct {
	Nodes    []workflowCheckStepNode `json:"nodes"`
	PageInfo workflowPageInfo        `json:"pageInfo"`
}

type workflowCheckStepNode struct {
	Conclusion  WorkflowConclusion `json:"conclusion"`
	Name        string             `json:"name"`
	Number      int                `json:"number"`
	Status      string             `json:"status"`
	CompletedAt time.Time          `json:"completedAt"`
	StartedAt   time.Time          `json:"startedAt"`
}

type workflowPageInfo struct {
	HasNextPage bool `json:"hasNextPage"`
}

type workflowRunCheckSuiteMetadata struct {
	checkRunNodeIDs []string
	suite           workflowCheckSuiteJobsNode
}

// WorkflowRunDetailsForCheckRuns batches failed check run node IDs into suite and job metadata.
func (c *Client) WorkflowRunDetailsForCheckRuns(ctx context.Context, checkRunNodeIDs []string) (map[string]WorkflowRunCheckSuiteDetails, error) {
	checkRunNodeIDs = uniqueStrings(checkRunNodeIDs)
	if len(checkRunNodeIDs) == 0 {
		return map[string]WorkflowRunCheckSuiteDetails{}, nil
	}
	metadataBySuiteID := make(map[string]*workflowRunCheckSuiteMetadata)
	for start := 0; start < len(checkRunNodeIDs); start += defaultCheckRunNodeBatchSize {
		end := min(start+defaultCheckRunNodeBatchSize, len(checkRunNodeIDs))
		variables := workflowNodesVariables{IDs: checkRunNodeIDs[start:end]}
		var page workflowRunMetadataPage
		if err := c.queryWithVariables(ctx, checkSuiteWorkflowRunNodeQuery, variables, &page); err != nil {
			return nil, fmt.Errorf(errorQueryCheckSuiteIDs, err)
		}
		for _, node := range page.Nodes {
			if node == nil || node.ID == "" || node.CheckSuite == nil || node.CheckSuite.ID == "" {
				continue
			}
			metadata := metadataBySuiteID[node.CheckSuite.ID]
			if metadata == nil {
				metadata = &workflowRunCheckSuiteMetadata{}
				metadataBySuiteID[node.CheckSuite.ID] = metadata
			}
			metadata.checkRunNodeIDs = append(metadata.checkRunNodeIDs, node.ID)
		}
	}
	if len(metadataBySuiteID) == 0 {
		return map[string]WorkflowRunCheckSuiteDetails{}, nil
	}
	suiteNodeIDs := make([]string, 0, len(metadataBySuiteID))
	for suiteNodeID := range metadataBySuiteID {
		suiteNodeIDs = append(suiteNodeIDs, suiteNodeID)
	}
	sort.Strings(suiteNodeIDs)
	for start := 0; start < len(suiteNodeIDs); start += defaultCheckSuiteNodeBatchSize {
		end := min(start+defaultCheckSuiteNodeBatchSize, len(suiteNodeIDs))
		variables := workflowNodesVariables{IDs: suiteNodeIDs[start:end]}
		var page workflowCheckSuiteJobsPage
		if err := c.queryWithVariables(ctx, checkSuiteJobsNodeQuery, variables, &page); err != nil {
			return nil, fmt.Errorf(errorQueryCheckSuiteJobs, err)
		}
		for _, suite := range page.Nodes {
			if suite == nil || suite.ID == "" {
				continue
			}
			metadata := metadataBySuiteID[suite.ID]
			if metadata != nil {
				metadata.suite = *suite
			}
		}
	}
	workflowDetailsByCheckRunID := make(map[string]WorkflowRunCheckSuiteDetails)
	for _, metadata := range metadataBySuiteID {
		details, found := metadata.suite.workflowRunCheckSuiteDetails(c)
		if !found {
			continue
		}
		for _, checkRunNodeID := range metadata.checkRunNodeIDs {
			workflowDetailsByCheckRunID[checkRunNodeID] = details
		}
	}
	return workflowDetailsByCheckRunID, nil
}

type workflowNodesVariables struct {
	IDs []string `json:"ids"`
}

func (suite *workflowCheckSuiteJobsNode) workflowRunCheckSuiteDetails(client *Client) (WorkflowRunCheckSuiteDetails, bool) {
	if suite == nil || suite.WorkflowRun == nil || suite.WorkflowRun.RunAttempt < 1 || suite.Commit.OID == "" {
		return WorkflowRunCheckSuiteDetails{}, false
	}
	runID, found := client.workflowRunIDFromURL(suite.WorkflowRun.URL)
	if !found {
		return WorkflowRunCheckSuiteDetails{}, false
	}
	workflowName := suite.WorkflowRun.Workflow.Name
	if workflowName == "" {
		workflowName = suite.WorkflowRun.DisplayTitle
	}
	metadataComplete := workflowName != "" && suite.Branch.Name != ""
	attempt := WorkflowRunCheckSuiteDetails{
		Attempt: common.WorkflowRunAttempt{
			RunID:        runID,
			Attempt:      suite.WorkflowRun.RunAttempt,
			Name:         workflowName,
			Branch:       suite.Branch.Name,
			HeadSHA:      suite.Commit.OID,
			RunStartedAt: suite.WorkflowRun.CreatedAt,
			Status:       suite.Status,
			Conclusion:   string(suite.Conclusion),
			WorkflowURL:  suite.WorkflowRun.URL,
		},
		Jobs:       make([]common.WorkflowJob, 0, len(suite.CheckRuns.Nodes)),
		IsComplete: len(suite.CheckRuns.Nodes) > 0 && !suite.CheckRuns.PageInfo.HasNextPage && metadataComplete,
	}
	earliestJobStart := time.Time{}
	for _, checkRun := range suite.CheckRuns.Nodes {
		jobID, validJobID := checkRun.jobID()
		if !validJobID {
			attempt.IsComplete = false
			continue
		}
		workflowJob := common.WorkflowJob{
			ID:         jobID,
			Name:       checkRun.Name,
			Conclusion: string(checkRun.Conclusion),
			JobURL:     checkRun.URL,
			Steps:      make([]common.WorkflowStep, 0, len(checkRun.Steps.Nodes)),
		}
		for _, checkStep := range checkRun.Steps.Nodes {
			workflowJob.Steps = append(workflowJob.Steps, common.WorkflowStep{
				Number:      checkStep.Number,
				Name:        checkStep.Name,
				Conclusion:  string(checkStep.Conclusion),
				StartedAt:   checkStep.StartedAt,
				CompletedAt: checkStep.CompletedAt,
			})
		}
		if checkRun.Steps.PageInfo.HasNextPage {
			attempt.IsComplete = false
		}
		if !checkRun.StartedAt.IsZero() && (earliestJobStart.IsZero() || checkRun.StartedAt.Before(earliestJobStart)) {
			earliestJobStart = checkRun.StartedAt
		}
		attempt.Jobs = append(attempt.Jobs, workflowJob)
	}
	if !earliestJobStart.IsZero() {
		attempt.Attempt.RunStartedAt = earliestJobStart
	}
	if attempt.Attempt.RunStartedAt.IsZero() {
		attempt.IsComplete = false
	}
	return attempt, true
}

func (checkRun *workflowCheckRunNode) jobID() (string, bool) {
	parsedURL, err := url.Parse(checkRun.URL)
	if err != nil || !strings.EqualFold(parsedURL.Hostname(), githubHost) {
		return "", false
	}
	pathSegments := strings.Split(strings.Trim(parsedURL.Path, pathSeparator), pathSeparator)
	if len(pathSegments) == 0 {
		return "", false
	}
	jobID := pathSegments[len(pathSegments)-1]
	parsedJobID, err := strconv.ParseInt(jobID, 10, 64)
	return jobID, err == nil && parsedJobID > 0
}

func (c *Client) workflowRunIDFromURL(workflowURL string) (string, bool) {
	parsedURL, err := url.Parse(workflowURL)
	if err != nil || !strings.EqualFold(parsedURL.Hostname(), githubHost) {
		return "", false
	}
	pathSegments := strings.Split(strings.Trim(parsedURL.Path, pathSeparator), pathSeparator)
	if len(pathSegments) < 5 || !strings.EqualFold(pathSegments[0], c.owner) || !strings.EqualFold(pathSegments[1], c.repositoryName) || pathSegments[2] != githubActionsPathSegment || pathSegments[3] != githubRunsPathSegment {
		return "", false
	}
	// The run ID follows the `runs` segment in the GitHub workflow URL.
	runID := pathSegments[4]
	parsedRunID, err := strconv.ParseInt(runID, 10, 64)
	return runID, err == nil && parsedRunID > 0
}

func uniqueStrings(values []string) []string {
	uniqueValues := make([]string, 0, len(values))
	seenValues := make(map[string]struct{}, len(values))
	for _, value := range values {
		if _, found := seenValues[value]; found {
			continue
		}
		seenValues[value] = struct{}{}
		uniqueValues = append(uniqueValues, value)
	}
	return uniqueValues
}
