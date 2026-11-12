package prstatus

import "ktlo/internal/github"

const (
	dependabotAuthorLogin   = "dependabot"
	renovateAuthorLogin     = "renovate"
	srvGHO11yGDIAuthorLogin = "srv-gh-o11y-gdi"
)

var dependencyUpdateAuthors = [...]github.PullRequestAuthor{
	{Login: dependabotAuthorLogin},
	{Login: renovateAuthorLogin},
	{Login: srvGHO11yGDIAuthorLogin},
}
