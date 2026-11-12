package github

const (
	openDependencyBumpPullRequestsQuery = `query OpenDependencyBumpPullRequests($owner: String!, $name: String!, $after: String) { repository(owner: $owner, name: $name) { pullRequests(first: 100, states: [OPEN], after: $after) { nodes { author { login } isDraft number title url } pageInfo { endCursor hasNextPage } } } }`
	openPullRequestsQuery               = `query OpenPullRequests($owner: String!, $name: String!, $after: String) { repository(owner: $owner, name: $name) { pullRequests(first: 100, states: [OPEN], after: $after) { nodes { author { login } isDraft number title url } pageInfo { endCursor hasNextPage } } } }`
	pullRequestCICDResultsQuery         = `query PullRequestCICDResults($owner: String!, $name: String!, $number: Int!, $after: String) { repository(owner: $owner, name: $name) { pullRequest(number: $number) { statusCheckRollup { state contexts(first: 100, after: $after) { nodes { __typename ... on CheckRun { name status conclusion detailsUrl summary text startedAt completedAt isRequired(pullRequestNumber: $number) } ... on StatusContext { context state description targetUrl isRequired(pullRequestNumber: $number) } } pageInfo { endCursor hasNextPage } } } } } }`
	repositoryAccessQuery               = `query RepositoryAccess($owner: String!, $name: String!) { repository(owner: $owner, name: $name) { id } }`
)
