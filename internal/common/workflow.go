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

package common

import (
	"strings"
	"time"
)

const workflowConclusionFailure = "failure"
const workflowConclusionSuccess = "success"

type WorkflowRunAttempt struct {
	RunID        string
	Attempt      int
	Name         string
	Branch       string
	HeadSHA      string
	RunStartedAt time.Time
	Status       string
	Conclusion   string
	WorkflowURL  string
	PullRequests []PullRequestReference
}

type WorkflowRunAttemptKey struct {
	RunID   string
	Attempt int
}

type PullRequestReference struct {
	Number  int    `json:"number"`
	URL     string `json:"url,omitempty"`
	HeadSHA string `json:"headSha"`
}

func (reference *PullRequestReference) Equals(other *PullRequestReference) bool {
	if reference == nil || other == nil {
		return reference == other
	}
	return reference.Number == other.Number && reference.HeadSHA == other.HeadSHA
}

type WorkflowJob struct {
	ID         string
	Name       string
	Conclusion string
	JobURL     string
	Steps      []WorkflowStep
}

func (job *WorkflowJob) NeedsStepLogs() bool {
	conclusion := strings.ToLower(job.Conclusion)
	return len(job.Steps) > 0 && (conclusion == workflowConclusionFailure || conclusion == workflowConclusionSuccess)
}

type WorkflowStep struct {
	Number      int
	Name        string
	Conclusion  string
	CompletedAt time.Time
	StartedAt   time.Time
}
