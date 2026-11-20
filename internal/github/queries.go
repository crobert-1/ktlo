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

const (
	openDependencyBumpPullRequestsQuery = `query OpenDependencyBumpPullRequests($owner: String!, $name: String!, $after: String) { repository(owner: $owner, name: $name) { pullRequests(first: 100, states: [OPEN], after: $after) { nodes { author { login } headSha: headRefOid headRef { target { ... on Commit { committedDate } } } isDraft number title url } pageInfo { endCursor hasNextPage } } } }`
	openPullRequestsQuery               = `query OpenPullRequests($owner: String!, $name: String!, $after: String) { repository(owner: $owner, name: $name) { pullRequests(first: 100, states: [OPEN], after: $after) { nodes { author { login } headSha: headRefOid headRef { target { ... on Commit { committedDate } } } isDraft number title url } pageInfo { endCursor hasNextPage } } } }`
	pullRequestCICDResultsQuery         = `query PullRequestCICDResults($owner: String!, $name: String!, $number: Int!, $after: String) { repository(owner: $owner, name: $name) { pullRequest(number: $number) { statusCheckRollup { state contexts(first: 100, after: $after) { nodes { __typename ... on CheckRun { id name status conclusion detailsUrl summary text startedAt completedAt isRequired(pullRequestNumber: $number) } ... on StatusContext { context state description targetUrl isRequired(pullRequestNumber: $number) } } pageInfo { endCursor hasNextPage } } } } } }`
	repositoryAccessQuery               = `query RepositoryAccess($owner: String!, $name: String!) { repository(owner: $owner, name: $name) { id } }`
)
