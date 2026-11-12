package common

const (
	dependabotLogin   = "dependabot"
	renovateLogin     = "renovate"
	srvGHO11yGDILogin = "srv-gh-o11y-gdi"
)

type PullRequestAuthor struct {
	Login string
}

func DependencyBumpPullRequestAuthors() []PullRequestAuthor {
	return []PullRequestAuthor{
		{Login: dependabotLogin},
		{Login: renovateLogin},
		{Login: srvGHO11yGDILogin},
	}
}
