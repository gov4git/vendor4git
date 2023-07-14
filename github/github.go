package github

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"github.com/google/go-github/v53/github"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/vendor4git"
	"golang.org/x/oauth2"
)

type GitHubVendor struct{}

func (GitHubVendor) CreateRepo(ctx context.Context, name string, org string, private bool) (*Repository, error) {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")}, //XXX
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repo := &github.Repository{
		Name:    github.String(name),
		Private: github.Bool(private),
	}
	repo, _, err := client.Repositories.Create(ctx, org, repo)
	must.NoError(ctx, err) //XXX

	return &vendor4git.Repository{
		HTTPSURL: repo.GetCloneURL(),
		SSHURL:   repo.GetSSHURL(),
	}, nil
}

func (GitHubVendor) RemoveRepo(name string, org string) error {
	panic("XXX")
}
