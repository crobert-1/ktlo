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

package failingtest

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"ktlo/internal/telemetry"
	"ktlo/internal/testresults"
)

const (
	errorMissingWorkflowProvider = "workflow provider is required to parse workflow run attempts"
	errorParseWorkflowRunAttempt = "could not parse workflow run %s attempt %d: %w"
	spanParseWorkflowRunAttempts = "failingtest.parse_workflow_run_attempts"
)

const maxConcurrentWorkflowRunBuilds = 2

func ParseWorkflowRunAttempts(ctx context.Context, workflowProvider WorkflowProvider, workflowRunAttempts []WorkflowRunAttempt, testLogParsers ...TestLogParser) (workflowRuns []testresults.Workflow, resultErr error) {
	return parseWorkflowRunAttempts(ctx, workflowProvider, workflowRunAttempts, nil, testLogParsers...)
}

// ParseWorkflowRunAttemptsWithJobs parses attempts using supplied job metadata and provider logs.
func ParseWorkflowRunAttemptsWithJobs(ctx context.Context, workflowProvider WorkflowProvider, workflowRunAttempts []WorkflowRunAttempt, workflowJobsByRunID map[string][]WorkflowJob, testLogParsers ...TestLogParser) (workflowRuns []testresults.Workflow, resultErr error) {
	return parseWorkflowRunAttempts(ctx, workflowProvider, workflowRunAttempts, workflowJobsByRunID, testLogParsers...)
}

func parseWorkflowRunAttempts(ctx context.Context, workflowProvider WorkflowProvider, workflowRunAttempts []WorkflowRunAttempt, workflowJobsByRunID map[string][]WorkflowJob, testLogParsers ...TestLogParser) (workflowRuns []testresults.Workflow, resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanParseWorkflowRunAttempts,
		attribute.Int(traceAttributeWorkflowRunAttemptCount, len(workflowRunAttempts)),
	)
	defer func() {
		span.SetAttributes(attribute.Int(traceAttributeWorkflowRecordCount, len(workflowRuns)))
		telemetry.FinishSpan(ctx, span, spanParseWorkflowRunAttempts, startedAt, resultErr != nil)
	}()
	if workflowProvider == nil {
		return nil, errors.New(errorMissingWorkflowProvider)
	}
	if len(testLogParsers) == 0 {
		testLogParsers = defaultTestLogParsers()
	}
	testLogParsers = orderTestLogParsers(testLogParsers)
	buildResults := make([]workflowRunAttemptBuildResult, len(workflowRunAttempts))
	parserMutex := &sync.Mutex{}
	workerCount := min(maxConcurrentWorkflowRunBuilds, len(workflowRunAttempts))
	attemptIndexes := make(chan int)
	var workers sync.WaitGroup
	for range workerCount {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for attemptIndex := range attemptIndexes {
				workflowRunAttempt := workflowRunAttempts[attemptIndex]
				if err := ctx.Err(); err != nil {
					buildResults[attemptIndex].err = err
					continue
				}
				workflowJobs, useProvidedJobs := workflowJobsByRunID[workflowRunAttempt.RunID]
				builder := newWorkflowRecordBuilder(workflowProvider, testLogParsers)
				builder.testLogParserMutex = parserMutex
				if useProvidedJobs {
					buildResults[attemptIndex].workflow, buildResults[attemptIndex].err = builder.buildWithWorkflowJobs(ctx, workflowRunAttempt, workflowJobs, true)
				} else {
					buildResults[attemptIndex].workflow, buildResults[attemptIndex].err = builder.build(ctx, workflowRunAttempt)
				}
			}
		}()
	}
	for attemptIndex := range workflowRunAttempts {
		attemptIndexes <- attemptIndex
	}
	close(attemptIndexes)
	workers.Wait()
	workflowRuns = make([]testresults.Workflow, 0, len(workflowRunAttempts))
	collectionErrors := make([]error, 0)
	for attemptIndex, buildResult := range buildResults {
		if buildResult.err != nil {
			workflowRunAttempt := workflowRunAttempts[attemptIndex]
			collectionErrors = append(collectionErrors, fmt.Errorf(errorParseWorkflowRunAttempt, workflowRunAttempt.RunID, workflowRunAttempt.Attempt, buildResult.err))
		}
		if buildResult.workflow != nil {
			workflowRuns = append(workflowRuns, *buildResult.workflow)
		}
	}
	return workflowRuns, errors.Join(collectionErrors...)
}

type workflowRunAttemptBuildResult struct {
	err      error
	workflow *testresults.Workflow
}
